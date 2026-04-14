package iblt

import (
	"fmt"
	"testing"
)

func TestIBLT_InsertAndDecode(t *testing.T) {
	// 创建一个支持约 10 个差异的 IBLT
	iblt := NewIBLT(30)

	// 插入一些键值对
	keys := []uint64{1, 2, 3, 4, 5}
	for _, key := range keys {
		iblt.Insert(key, key*10)
	}

	// 解码应该返回所有插入的键
	onlyA, onlyB, err := iblt.Decode()
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(onlyA) != len(keys) {
		t.Errorf("Expected %d keys in onlyA, got %d", len(keys), len(onlyA))
	}

	if len(onlyB) != 0 {
		t.Errorf("Expected 0 keys in onlyB, got %d", len(onlyB))
	}

	fmt.Printf("Successfully decoded %d keys\n", len(onlyA))
}

func TestIBLT_XORAndDecode(t *testing.T) {
	// 模拟两个集合的差异发现
	// 集合 A: {1, 2, 3, 4, 5}
	// 集合 B: {1, 2, 3, 6, 7}
	// 差异：A 独有 {4, 5}, B 独有 {6, 7}

	ibltA := NewIBLT(30)
	ibltB := NewIBLT(30)

	// 构建 A 的 IBLT
	for _, key := range []uint64{1, 2, 3, 4, 5} {
		ibltA.Insert(key, key*10)
	}

	// 构建 B 的 IBLT
	for _, key := range []uint64{1, 2, 3, 6, 7} {
		ibltB.Insert(key, key*10)
	}

	// XOR 得到差异 IBLT
	diffIBLT, err := ibltA.XOR(ibltB)
	if err != nil {
		t.Fatalf("XOR failed: %v", err)
	}

	// 解码差异
	onlyA, onlyB, err := diffIBLT.Decode()
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// 验证结果
	if len(onlyA) != 2 {
		t.Errorf("Expected 2 keys in onlyA, got %d", len(onlyA))
	}

	if len(onlyB) != 2 {
		t.Errorf("Expected 2 keys in onlyB, got %d", len(onlyB))
	}

	// 检查具体键值
	aKeys := make(map[uint64]bool)
	for _, kv := range onlyA {
		aKeys[kv.Key] = true
	}

	bKeys := make(map[uint64]bool)
	for _, kv := range onlyB {
		bKeys[kv.Key] = true
	}

	expectedA := map[uint64]bool{4: true, 5: true}
	expectedB := map[uint64]bool{6: true, 7: true}

	for key := range expectedA {
		if !aKeys[key] {
			t.Errorf("Expected key %d in onlyA, but not found", key)
		}
	}

	for key := range expectedB {
		if !bKeys[key] {
			t.Errorf("Expected key %d in onlyB, but not found", key)
		}
	}

	fmt.Printf("Successfully found differences: A-only=%v, B-only=%v\n", onlyA, onlyB)
}

func TestIBLT_SerializeDeserialize(t *testing.T) {
	iblt := NewIBLT(20)

	// 插入一些数据
	for i := uint64(1); i <= 5; i++ {
		iblt.Insert(i, i*100)
	}

	// 序列化
	data := iblt.Serialize()
	fmt.Printf("Serialized size: %d bytes\n", len(data))

	// 反序列化
	iblt2, err := Deserialize(data)
	if err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	// 验证反序列化后的 IBLT 与原 IBLT 相同
	diffIBLT, err := iblt.XOR(iblt2)
	if err != nil {
		t.Fatalf("XOR after deserialize failed: %v", err)
	}

	onlyA, onlyB, err := diffIBLT.Decode()
	if err != nil {
		t.Fatalf("Decode after deserialize failed: %v", err)
	}

	if len(onlyA) != 0 || len(onlyB) != 0 {
		t.Error("Deserialized IBLT should be identical to original")
	}

	fmt.Println("Serialize/Deserialize test passed")
}

func TestIBLT_Size(t *testing.T) {
	iblt := NewIBLT(50)
	size := iblt.Size()
	expected := 50 * 20 // 每个 cell 20 字节

	if size != expected {
		t.Errorf("Expected size %d, got %d", expected, size)
	}

	fmt.Printf("IBLT with 50 cells occupies %d bytes\n", size)
}

func TestIBLT_LargeSet(t *testing.T) {
	// 测试大规模数据集（模拟 1000 万文件场景）
	const totalFiles = 100000
	const diffCount = 5

	// IBLT 大小应该根据差异数量来设置，建议为 3*d ~ 5*d
	// 这里我们使用更大的 IBLT 以确保解码成功
	ibltA := NewIBLT(diffCount * 50) // 增加 IBLT 大小以提高成功率
	ibltB := NewIBLT(diffCount * 50)

	// 插入大量相同的数据
	for i := uint64(1); i <= totalFiles; i++ {
		ibltA.Insert(i, 0)
		ibltB.Insert(i, 0)
	}

	// 添加一些差异
	for i := uint64(1); i <= diffCount; i++ {
		ibltA.Insert(totalFiles+i, 0) // A 独有的文件
	}
	for i := uint64(1); i <= diffCount; i++ {
		ibltB.Insert(totalFiles+diffCount+i, 0) // B 独有的文件
	}

	// XOR 并解码
	diffIBLT, err := ibltA.XOR(ibltB)
	if err != nil {
		t.Fatalf("XOR failed: %v", err)
	}

	onlyA, onlyB, err := diffIBLT.Decode()
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(onlyA) != diffCount {
		t.Errorf("Expected %d differences in onlyA, got %d", diffCount, len(onlyA))
	}

	if len(onlyB) != diffCount {
		t.Errorf("Expected %d differences in onlyB, got %d", diffCount, len(onlyB))
	}

	fmt.Printf("Large set test passed: found %d differences in each direction\n", len(onlyA))
	fmt.Printf("Memory usage: IBLT-A=%d bytes, IBLT-B=%d bytes\n", ibltA.Size(), ibltB.Size())
}

func BenchmarkIBLT_Insert(b *testing.B) {
	iblt := NewIBLT(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		iblt.Insert(uint64(i), uint64(i*10))
	}
}

func BenchmarkIBLT_Decode(b *testing.B) {
	iblt := NewIBLT(100)
	for i := 0; i < 30; i++ {
		iblt.Insert(uint64(i), uint64(i*10))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = iblt.Decode()
	}
}
