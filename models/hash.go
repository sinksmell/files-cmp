package models

import (
	"os"
	"crypto/md5"
	"io"
	"fmt"
)

// 计算文件的md5值
func GetMd5(filepath string)(hash string, err error){

	var(
		file *os.File
	)
	if file,err=os.Open(filepath);err!=nil{
		// 出现错误就返回
		return
	}
	defer file.Close()
	h:=md5.New()
	if _,err=io.Copy(h,file);err!=nil{
		return
	}

	return fmt.Sprintf("%x",h.Sum(nil)),nil
}
