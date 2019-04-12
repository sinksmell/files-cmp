package models

import (
	"fmt"
	"math"
)

const (
	FILE_PATH = "./file/"
	GROUP_PATH="./group/"
	GROUP_PRE = "group_" // 分组文件 文件名的前缀
	GROUP_SUF = ".txt"   // 分组文件 文件名的后缀
)

var(
	FILE_CNT int
	Files []string
	GRP_SIZE int
)


func init(){
	InitFiles(FILE_PATH)
}



func InitFiles(filepath string){
	var(
		err error
	)
	if Files,err=GetAllFiles(filepath);err!=nil{
		fmt.Println(err)
		return
	}
	FILE_CNT=len(Files)
	GRP_SIZE=int(math.Sqrt(float64(1.0*FILE_CNT)))
}


type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 用于检测hash值是否相同的请求结构
type HashRequest struct {
	FileName string `json:"fname"`
	Hash     string `json:"hash"`
}
