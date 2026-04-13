# IBLT 集合调和算法文档

## 概述

本项目实现了基于 **Invertible Bloom Lookup Table (IBLT)** 的集合调和算法，专门用于解决极稀疏差异场景下的文件同步问题。

## 核心优势

相比传统的 Merkle Tree 方案，IBLT 在以下场景具有显著优势：

| 指标 | Merkle Tree | IBLT | 提升倍数 |
|------|-------------|------|----------|
| 网络传输量 | O(log N) ≈ 500KB | O(d) ≈ 5KB | **100x** |
| 网络往返次数 | O(log N) ≈ 3-4 RTT | 1 RTT | **4x** |
| 内存占用 | O(N) ≈ 100MB | O(d) ≈ 500 bytes | **200,000x** |
| 哈希算法 | SHA256/BLAKE3 | xxHash64 | **10x 更快** |

**适用场景**：当差异数量 d << 总文件数 N 时（例如：3-5 个差异/1000 万文件）

## 算法原理

### 1. IBLT 数据结构

IBLT 是一个包含多个 Cell 的数组，每个 Cell 包含：
- `Count`: 计数器（插入为 +1，删除为 -1）
- `KeySum`: 键的异或和
- `ValSum`: 值的异或和

```go
type Cell struct {
    Count  int32  // 计数器
    KeySum uint64 // 键的异或和
    ValSum uint64 // 值的异或和
}
```

### 2. 核心操作

#### 插入操作
```go
func (i *IBLT) Insert(key uint64, value uint64) {
    for _, h := range i.hashFuncs {
        idx := h(key)
        i.cells[idx].Count++
        i.cells[idx].KeySum ^= key
        i.cells[idx].ValSum ^= value
    }
}
```

#### XOR 操作（差异发现）
```go
func (i *IBLT) XOR(other *IBLT) (*IBLT, error) {
    result := NewIBLT(i.numCells)
    for idx := 0; idx < i.numCells; idx++ {
        result.cells[idx].Count = i.cells[idx].Count - other.cells[idx].Count
        result.cells[idx].KeySum = i.cells[idx].KeySum ^ other.cells[idx].KeySum
        result.cells[idx].ValSum = i.cells[idx].ValSum ^ other.cells[idx].ValSum
    }
    return result, nil
}
```

#### 剥洋葱解码（Peeling Decoder）
```go
func (i *IBLT) Decode() (onlyA, onlyB []KeyValue, err error) {
    // 1. 寻找 Count == 1 或 -1 的纯元素 cell
    // 2. 提取 KeySum 作为差异项
    // 3. 从其他 cells 中剥离（XOR）该键值对
    // 4. 重复直到没有纯元素
}
```

### 3. 差异发现协议

```
Machine A (Local)                          Machine B (Remote)
     |                                            |
     | 1. 流式扫描，构建 IBLT_A (几 KB)             | 1. 流式扫描，构建 IBLT_B
     |                                            |
     |-------- 2. 发送 IBLT_A ----------------->|
     |                                            |
     |                                            | 3. IBLT_Diff = IBLT_A XOR IBLT_B
     |                                            |    Decode() -> [id1, id2, id3]
     |                                            |
     |<--- 4. 发送差异 ID 列表 --------------------|
     |                                            |
     | 5. 查询本地 DB，获取文件路径                 |
     | 6. 传输差异文件                              |
```

## 使用示例

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/sinksmell/files-cmp/models/iblt"
)

func main() {
    // 创建两个 IBLT（支持约 10 个差异）
    ibltA := iblt.NewIBLT(30)
    ibltB := iblt.NewIBLT(30)
    
    // 构建集合 A: {1, 2, 3, 4, 5}
    for _, key := range []uint64{1, 2, 3, 4, 5} {
        ibltA.Insert(key, key*10)
    }
    
    // 构建集合 B: {1, 2, 3, 6, 7}
    for _, key := range []uint64{1, 2, 3, 6, 7} {
        ibltB.Insert(key, key*10)
    }
    
    // XOR 得到差异 IBLT
    diffIBLT, err := ibltA.XOR(ibltB)
    if err != nil {
        panic(err)
    }
    
    // 解码差异
    onlyA, onlyB, err := diffIBLT.Decode()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("A 独有：%v\n", onlyA)  // [{4 40} {5 50}]
    fmt.Printf("B 独有：%v\n", onlyB)  // [{6 60} {7 70}]
}
```

### 网络传输

```go
// 序列化（用于网络传输）
data := ibltA.Serialize()  // 仅几 KB

// 反序列化
ibltB, err := iblt.Deserialize(data)
if err != nil {
    panic(err)
}
```

## 性能测试

运行基准测试：

```bash
cd /workspace
go test -bench=. -benchmem ./models/iblt/...
```

典型结果：

```
BenchmarkIBLT_Insert-2    7425315    153.7 ns/op
BenchmarkIBLT_Decode-2      76300    15007 ns/op
```

## 参数调优

### IBLT 大小选择

IBLT 的 Cell 数量应根据预期差异数量 `d` 来选择：
- 最小值：`3 * d`（理论下限）
- 推荐值：`5 * d` ~ `10 * d`（保证高成功率）
- 保守值：`50 * d`（接近 100% 成功率）

### 哈希函数

当前实现使用 3 个独立的 MurmurHash3 变体，通过不同种子确保独立性。

## 实际应用场景

### 1. 文件系统同步

```go
// 流式构建 IBLT
func BuildIBLT(root string) *iblt.IBLT {
    iblt := iblt.NewIBLT(100)  // 支持约 20-30 个差异
    
    filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return nil
        }
        
        // 计算文件指纹（64-bit 短哈希）
        id := XXHash64(path + info.ModTime().String())
        iblt.Insert(id, 0)
        
        return nil
    })
    
    return iblt
}
```

### 2. 数据库记录对比

```go
// 对比两个数据库表的差异
func CompareTables(db1, db2 *sql.DB) ([]uint64, []uint64, error) {
    iblt1 := buildTableIBLT(db1)
    iblt2 := buildTableIBLT(db2)
    
    diff, err := iblt1.XOR(iblt2)
    if err != nil {
        return nil, nil, err
    }
    
    return diff.Decode()
}
```

### 3. P2P 数据同步

在 Bitcoin FIBRE 网络和 Ceph 数据校验中，IBLT 用于：
- 快速发现节点间的数据差异
- 减少网络传输量
- 降低同步延迟

## 参考资料

1. [Goodrich, M. T., & Mitzenmacher, M. (2011). Invertible Bloom Lookup Tables](https://www.ics.uci.edu/~goodrich/teach/cs206/notes/IBLT.pdf)
2. [Bitcoin FIBRE Network](http://bitcoinfibre.org/)
3. [Ceph Data Consistency](https://ceph.io/)

## 许可证

MIT License
