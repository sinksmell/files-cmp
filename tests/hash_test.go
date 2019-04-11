package test

import (
	"testing"
	"github.com/sinksmell/files-cmp/models"
	"fmt"
)


// 测试能否正确地计算md5值
func TestGetMd5(t *testing.T) {
	if hash, err := models.GetMd5("./test.txt"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(hash)
	}
}
