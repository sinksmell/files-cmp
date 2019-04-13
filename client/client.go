package main

import (
	"github.com/sinksmell/files-cmp/models"
	"github.com/sinksmell/files-cmp/client/utils"
	"fmt"
)

const (
	HOST     string = "http://localhost:8080/v1/check"
	HASH_URL string = "/hash"
	FILE_URL string = "/file"
)

var (
	groups []string
)

func init() {

	// 初始化http请求端
	//utils.Init()
	// 初始化文件列表
	models.InitFiles(models.FILE_PATH)
	// 文件分组 组文件中存放 文件名 MD5值
	models.Divide(models.FILE_PATH, models.GROUP_PATH)
	// 获取组文件列表
	groups, _ = models.GetAllFiles(models.GROUP_PATH)
}

func main() {

	var (
		resp *models.Response // 文件对比结果
		err  error
	)

	for _, grp := range groups {
		if resp, err = utils.SendGrpMd5(grp,HOST+HASH_URL); err != nil {
			fmt.Println(err)
			continue
		}
		switch resp.Code {
		case models.EQUAL:
			fmt.Println(resp)
		case models.NOT_EQUAL:
			fmt.Println("NOT EQUAL")
		handleDiff(grp)
		case models.REQ_ERR:
		}
	}
}

func handleDiff(grpFile string) {
	if res, err := utils.PostFile(grpFile, HOST+FILE_URL,models.GROUP_PATH);err!=nil {
		fmt.Println(err)
		return
	}else{
		fmt.Println(res)
	}
}
