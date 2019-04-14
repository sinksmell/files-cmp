package main

import (
	_ "github.com/sinksmell/files-cmp/routers"

	"github.com/astaxie/beego"
	"github.com/sinksmell/files-cmp/models"
)

func init(){
	// 初始化文件列表
	models.InitFiles(models.FILE_PATH)
	// 文件分组 组文件中存放 文件名 MD5值
	models.Divide(models.FILE_PATH, models.GROUP_PATH)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
