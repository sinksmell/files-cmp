package test

import (
	"fmt"
	"github.com/sinksmell/files-cmp/models"
	"testing"
)

// 测试能否正确地计算md5值
func TestGetMd5(t *testing.T) {
	if hash, err := models.GetMd5("./test.txt"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(hash)
	}
}

// 测试能否正确计算出二进制文件的md5值
func TestGetBMd5(t *testing.T) {
	if hash, err := models.GetMd5("fbin"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(hash)
	}

}
