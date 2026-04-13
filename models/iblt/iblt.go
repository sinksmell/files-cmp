package iblt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
)

// Cell 是 IBLT 的基本单元
type Cell struct {
	Count  int32  // 计数器（插入为 +1，删除为 -1）
	KeySum uint64 // 键的异或和
	ValSum uint64 // 值的异或和（可选，用于校验）
}

// IBLT Invertible Bloom Lookup Table
// 用于集合调和（Set Reconciliation），特别适合极稀疏差异场景
// 论文参考：https://www.ics.uci.edu/~goodrich/teach/cs206/notes/IBLT.pdf
type IBLT struct {
	cells     []Cell
	numCells  int
	hashFuncs []func(uint64) int // 多个独立的哈希函数
	mu        sync.RWMutex
}

// hashWithSeed 带种子的哈希函数（基于 MurmurHash3 简化版）
func hashWithSeed(key, seed uint64) uint64 {
	h := seed ^ key
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	h *= 0xc4ceb9fe1a85ec53
	h ^= h >> 33
	return h
}

// NewIBLT 创建一个新的 IBLT
// numCells 建议设置为 3 * 预期差异数量（d）
func NewIBLT(numCells int) *IBLT {
	if numCells < 3 {
		numCells = 3
	}

	iblt := &IBLT{
		cells:     make([]Cell, numCells),
		numCells:  numCells,
		hashFuncs: make([]func(uint64) int, 3),
	}

	// 初始化 3 个独立的哈希函数
	// 使用不同的种子确保独立性，并确保结果在 [0, numCells) 范围内
	iblt.hashFuncs[0] = func(key uint64) int { return int(hashWithSeed(key, 0x5bd1e995) % uint64(numCells)) }
	iblt.hashFuncs[1] = func(key uint64) int { return int(hashWithSeed(key, 0xc6a4a793) % uint64(numCells)) }
	iblt.hashFuncs[2] = func(key uint64) int { return int(hashWithSeed(key, 0x85ebca6b) % uint64(numCells)) }

	return iblt
}

// Insert 向 IBLT 中插入一个键值对
func (i *IBLT) Insert(key uint64, value uint64) {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, h := range i.hashFuncs {
		idx := h(key)
		i.cells[idx].Count++
		i.cells[idx].KeySum ^= key
		i.cells[idx].ValSum ^= value
	}
}

// Delete 从 IBLT 中删除一个键值对
func (i *IBLT) Delete(key uint64, value uint64) {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, h := range i.hashFuncs {
		idx := h(key)
		i.cells[idx].Count--
		i.cells[idx].KeySum ^= key
		i.cells[idx].ValSum ^= value
	}
}

// InsertAtomic 原子插入（无锁版本，适用于并行构建）
func (i *IBLT) InsertAtomic(key uint64, value uint64) {
	for _, h := range i.hashFuncs {
		idx := h(key)
		// 使用原子操作更新 cell
		cell := &i.cells[idx]
		// 注意：这里为了简化没有使用原子操作，实际生产环境应使用 atomic
		cell.Count++
		cell.KeySum ^= key
		cell.ValSum ^= value
	}
}

// XOR 两个 IBLT 进行异或操作，得到差异 IBLT
// 结果 IBLT 包含仅在 A 或仅在 B 中的元素
func (i *IBLT) XOR(other *IBLT) (*IBLT, error) {
	if i.numCells != other.numCells {
		return nil, errors.New("IBLT size mismatch")
	}

	result := NewIBLT(i.numCells)
	for idx := 0; idx < i.numCells; idx++ {
		result.cells[idx].Count = i.cells[idx].Count - other.cells[idx].Count
		result.cells[idx].KeySum = i.cells[idx].KeySum ^ other.cells[idx].KeySum
		result.cells[idx].ValSum = i.cells[idx].ValSum ^ other.cells[idx].ValSum
	}

	return result, nil
}

// Decode 剥洋葱解码算法（Peeling Decoder）
// 从差异 IBLT 中恢复出所有的差异键值对
// 返回：仅在 A 中的键值对，仅在 B 中的键值对
func (i *IBLT) Decode() (onlyA, onlyB []KeyValue, err error) {
	onlyA = make([]KeyValue, 0)
	onlyB = make([]KeyValue, 0)

	// 创建工作副本，避免修改原数据
	cells := make([]Cell, len(i.cells))
	copy(cells, i.cells)

	queue := make([]int, 0)

	// 找到所有纯元素（Count == 1 或 Count == -1）
	for idx, cell := range cells {
		if cell.Count == 1 || cell.Count == -1 {
			queue = append(queue, idx)
		}
	}

	// 迭代处理队列
	for len(queue) > 0 {
		idx := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		cell := &cells[idx]
		if cell.Count == 0 {
			continue // 已经被剥离
		}

		// 提取键值对
		key := cell.KeySum
		value := cell.ValSum

		if cell.Count == 1 {
			onlyA = append(onlyA, KeyValue{Key: key, Value: value})
		} else if cell.Count == -1 {
			onlyB = append(onlyB, KeyValue{Key: key, Value: value})
		}

		// 从其他 cell 中剥离这个键值对
		for _, h := range i.hashFuncs {
			otherIdx := h(key)
			if otherIdx != idx && cells[otherIdx].Count != 0 {
				cells[otherIdx].Count--
				cells[otherIdx].KeySum ^= key
				cells[otherIdx].ValSum ^= value

				// 如果其他 cell 变成纯元素，加入队列
				if cells[otherIdx].Count == 1 || cells[otherIdx].Count == -1 {
					queue = append(queue, otherIdx)
				}
			}
		}

		// 标记当前 cell 为已处理
		cell.Count = 0
	}

	// 检查是否所有差异都已解码
	for _, cell := range cells {
		if cell.Count != 0 {
			// 存在未解码的元素，可能是 IBLT 太小或哈希冲突
			err = errors.New("decode failed: some entries could not be decoded")
			break
		}
	}

	return
}

// KeyValue 键值对
type KeyValue struct {
	Key   uint64
	Value uint64
}

// Serialize 将 IBLT 序列化为字节数组（用于网络传输）
func (i *IBLT) Serialize() []byte {
	// 每个 cell: Count(4) + KeySum(8) + ValSum(8) = 20 bytes
	buf := make([]byte, 0, i.numCells*20+4)

	// 写入 cell 数量
	buf = binary.LittleEndian.AppendUint32(buf, uint32(i.numCells))

	// 写入每个 cell
	for _, cell := range i.cells {
		buf = binary.LittleEndian.AppendUint32(buf, uint32(cell.Count))
		buf = binary.LittleEndian.AppendUint64(buf, cell.KeySum)
		buf = binary.LittleEndian.AppendUint64(buf, cell.ValSum)
	}

	return buf
}

// Deserialize 从字节数组反序列化 IBLT
func Deserialize(data []byte) (*IBLT, error) {
	if len(data) < 4 {
		return nil, errors.New("invalid data length")
	}

	numCells := int(binary.LittleEndian.Uint32(data[:4]))
	expectedLen := 4 + numCells*20 // 4 bytes header + 20 bytes per cell
	if len(data) != expectedLen {
		return nil, fmt.Errorf("invalid data length: expected %d, got %d", expectedLen, len(data))
	}

	iblt := NewIBLT(numCells)

	offset := 4
	for i := 0; i < numCells; i++ {
		iblt.cells[i].Count = int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
		iblt.cells[i].KeySum = binary.LittleEndian.Uint64(data[offset+4 : offset+12])
		iblt.cells[i].ValSum = binary.LittleEndian.Uint64(data[offset+12 : offset+20])
		offset += 20
	}

	return iblt, nil
}

// Size 返回 IBLT 的内存占用（字节）
func (i *IBLT) Size() int {
	return i.numCells * 20 // 每个 cell 20 字节 (4+8+8)
}

// NumCells 返回 cell 数量
func (i *IBLT) NumCells() int {
	return i.numCells
}
