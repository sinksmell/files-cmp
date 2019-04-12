package models

import (
	"os/exec"
	"strings"
	"os"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
)

// 根据文件存放目录,获取目录下的所有文件名
func GetAllFiles(path string) (fnames []string, err error) {
	var (
		files string
		bytes []byte
		cmd   *exec.Cmd
	)
	// 通过ls命令 获取指定目录下的所有文件名  自动根据文件名排序
	cmd = exec.Command("/bin/bash", "-c", "ls "+path)

	if bytes, err = cmd.CombinedOutput(); err != nil {
		return
	}
	files = string(bytes)

	// ls命令输出结果是一行一行的 根据\n 对字符串进行切割
	names := strings.Split(files, "\n")
	// 多次实验发现 最后一个值是""(空串)  丢弃最后一个值
	end := len(names) - 1
	if end > 0 {
		fnames = names[:end]
	}
	return
}

// 对文件进行分组 并把文件名:md5值 写到 分组文件中
func Divide(filePath, grpPath string) {

	var (
		grp *os.File
		grpId  = 0 // 分组文件的id
		remain = 0 // 小组文件剩余量
		hash string
		grpName string
		content string
		err error
	)

	for _, fname := range Files {
		if hash,err=GetMd5(filePath+fname);err!=nil{
			fmt.Println(err)
			continue
		}

		if remain==0{
			// 剩余数量为0时再创建一个分组
			remain=GRP_SIZE
			grpId++
			grpName=fmt.Sprintf("%s%d%s",GROUP_PRE,grpId,GROUP_SUF)
			if grp!=nil{
				// 关闭上一个分组文件
				grp.Close()
			}
			// 新建一个分组文件
			if grp,err=os.OpenFile(grpPath+grpName,os.O_RDWR|os.O_TRUNC|os.O_CREATE,0666);err!=nil{

				fmt.Println(err)
				return
			}
		}

		// 写入记录
		content=fmt.Sprintf("%s : %s\n",fname,hash)
		grp.WriteString(content)
		remain--

	}

	// 关闭最后一个分组文件
	grp.Close()

}

// 文件对比
func FileDiff(file1 , file2 string)(diffs []diffmatchpatch.Diff, err error){
	var(
		bytes1 []byte
		bytes2 []byte
	)

	dmp:=diffmatchpatch.New()
	if bytes1,err=ioutil.ReadFile(file1);err!=nil{
		fmt.Println(err)
		return
	}

	if bytes2,err=ioutil.ReadFile(file2);err!=nil{
		fmt.Println(err)
		return
	}
	diffs=dmp.DiffMain(string(bytes1),string(bytes2),false)
	return
}

// 对比对应文件的 md5
func CmpMd5(filepath string,hash string)(equal bool,err error){
	var(
		// 本地文件的MD5值
		_hash string
	)
	// 计算本地文件的md5值
	if _hash, err = GetMd5(filepath);err!=nil{
		return
	}else{
		if _hash ==hash{
			equal=true
		}
	}
	return
}