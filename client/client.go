package main

import (
	"github.com/sinksmell/files-cmp/models"
	"github.com/sinksmell/files-cmp/client/utils"
	"fmt"
)

var(
	groups []string
)

func init() {

	// 初始化http请求端
	utils.Init()
	// 初始化文件列表
	models.InitFiles(models.FILE_PATH)
	// 文件分组 组文件中存放 文件名 MD5值
	models.Divide(models.FILE_PATH,models.GROUP_PATH)
	// 获取组文件列表
	groups,_=models.GetAllFiles(models.GROUP_PATH)
}

func main(){
	for _, grp := range groups {
		if resp, err := utils.SendGroups(grp);err!=nil{
			fmt.Println(err)
			continue
		}else{
			fmt.Println(resp)
		}
	}


}

