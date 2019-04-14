package main

import (
	"testing"
	"github.com/sinksmell/files-cmp/client/utils"
	"github.com/sinksmell/files-cmp/models"
	"fmt"
)


var(
	HOST     string = "http://localhost:8080/v1/check"
	HASH_URL string = "/hash"
	FILE_URL string = "/file"
)

// 测试是否能正确post json数据
func TestPostHash(t *testing.T) {

	req := &models.HashRequest{FileName: "hello.txt", Hash: "qwertyuiop"}
	if resp, err := utils.PostHash(req,HOST+HASH_URL); err != nil {
		t.Fatal(err)
		return
	} else {
		fmt.Println(resp)
	}

}




// 测试json解析是否正常
func TestParseResp(t*testing.T){
	if res,err:=utils.ParseResp([]byte(`{"code": 0,"msg": "OK"}`));err!=nil{
	//	fmt.Println(err)
		t.Fatal(err)
	}else{
		t.Log(res)
	}
}

// 测试能否正确地提交文件
func TestPostFile(t*testing.T){

	target:="http://localhost:8080/v1/check/file"
	fileName:="dog.png"

	if res,err:=utils.PostFile(fileName,target,models.FILE_PATH,models.CMP_FILE);err!=nil{
		fmt.Println(err)
		return
	}else {
		fmt.Println(res)
	}

}



// 测试能否正确地计算分组文件的MD5 并发送
func TestSendGroupMd5(t *testing.T){
	utils.SendGrpMd5("group_2.txt",HOST+HASH_URL)

}