package main

import (
	"testing"
	"github.com/sinksmell/files-cmp/client/utils"
	"github.com/sinksmell/files-cmp/models"
	"fmt"
)

// 测试是否能正确post json数据
func TestPostHash(t *testing.T) {

	utils.Init()
	req := &models.HashRequest{FileName: "hello.txt", Hash: "qwertyuiop"}
	if resp, err := utils.PostHash(req); err != nil {
		t.Fatal(err)
		return
	} else {
		fmt.Println(resp)
	}

}

// 测试能否正确地计算分组文件的MD5 并发送
func TestSendGroup(t *testing.T){
	utils.SendGroups("group_2.txt")

}