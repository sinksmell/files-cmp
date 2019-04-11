package models

import (
	"os/exec"
	"strings"
)

// 根据文件存放目录,获取目录下的所有文件名
func GetAllFiles(path string)(fnames[]string,err error) {
	var(
		files string
		bytes []byte
		cmd *exec.Cmd
	)
	// 通过ls命令 获取指定目录下的所有文件名  自动根据文件名排序
	cmd=exec.Command("/bin/bash","-c","ls "+path)

	if bytes,err=cmd.CombinedOutput();err!=nil{
		return
	}
	files=string(bytes)

	// ls命令输出结果是一行一行的 根据\n 对字符串进行切割
	names:=strings.Split(files,"\n")
	// 多次实验发现 最后一个值是""(空串)  丢弃最后一个值
	end:=len(names)-1
	if end>0{
		fnames=names[:end]
	}
	return
}
