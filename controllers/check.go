package controllers

import (
	"github.com/astaxie/beego"
	"github.com/sinksmell/files-cmp/models"
	"encoding/json"
	"fmt"
)


// 用于检测文件及md5的Controller
type  CheckController  struct{
    beego.Controller
}

func(c*CheckController)URLMapping(){
	c.Mapping("GET",c.Get)
	c.Mapping("Hash",c.Hash)
	c.Mapping("File",c.File)
}


// @Title Test CheckController
// @Description get all objects
// @Success 200
// @Failure 403
// @router / [get]
func(c*CheckController)Get(){
	resp:=&models.Response{0,"OK"}
	c.Data["json"]=resp
	c.ServeJSON()
}



// @Title Check Hash
// @Description 用于检测组文件的hash值是否相同
// @Param body  body  models.HashRequest true  "body for Check Content"
// @Success 200
// @Failure 403
// @router /hash [post]
func (c*CheckController)Hash(){
	var req models.HashRequest
	resp:=&models.Response{}
	if err:=json.Unmarshal(c.Ctx.Input.RequestBody,&req);err==nil{
		fmt.Println(req)
		resp.Code=0
		resp.Msg="OK"
	}
	c.Data["json"]=resp
	c.ServeJSON()
}

// Update ...
// @Title Update
// @Description 上传对应的文件 检测是否相同
// @Success 200
// @Failure 403 body is empty
// @router /file [post]
func (c*CheckController)File(){
	resp:=&models.Response{}
	f,h,err:=c.GetFile("file")
	if err!=nil{
		resp.Code=100
		resp.Msg=err.Error()
	}

	defer f.Close()
	err = c.SaveToFile("file", "static/upload/"+h.Filename) // 保存位置在 static/upload, 没有文件夹要先创建
	if err != nil {
		resp.Msg=err.Error()
		resp.Code=200
	}else{
		resp.Code=0
		resp.Msg="OK"
	}
	c.Data["json"]=resp
	c.ServeJSON()
}