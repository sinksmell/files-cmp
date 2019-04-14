package controllers

import (
	"github.com/astaxie/beego"
	"github.com/sinksmell/files-cmp/models"
	"encoding/json"
	"fmt"
	"mime/multipart"
)

// 用于检测文件及md5的Controller
type CheckController struct {
	beego.Controller
}

func (c *CheckController) URLMapping() {
	c.Mapping("GET", c.Get)
	c.Mapping("Hash", c.Hash)
	c.Mapping("File", c.File)
}

// @Title Test CheckController
// @Description get all objects
// @Success 200
// @Failure 403
// @router / [get]
func (c *CheckController) Get() {
	resp := &models.Response{}
	resp.Code = models.SUCCESS
	resp.Msg = "OK"
	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title Check Hash
// @Description 用于检测组文件的hash值是否相同
// @Param body  body  models.HashRequest true  "body for Check Content"
// @Success 200
// @Failure 403
// @router /hash [post]
func (c *CheckController) Hash() {
	var req models.HashRequest
	resp := &models.Response{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		fmt.Println(req)
		resp.Code = models.REQ_ERR
		resp.Msg = err.Error()
	} else {
		if equal, _ := models.CmpMd5(models.GROUP_PATH+req.FileName, req.Hash); equal {
			resp.Code = models.EQUAL
			resp.Msg = "OK"
		} else {
			resp.Code = models.NOT_EQUAL
			resp.Msg = "Not Equal"
		}

	}
	c.Data["json"] = resp
	c.ServeJSON()
}

// Update ...
// @Title Update
// @Description 上传对应的文件 检测是否相同
// @Success 200
// @Failure 403 body is empty
// @router /file [post]
func (c *CheckController) File() {
	var (
		resp *models.Response
		f    multipart.File
		h    *multipart.FileHeader
		err  error
	)

	resp = &models.Response{}
	checkType := c.GetString("type") // 获取文件比较类型 分组文件对比 还是小文件对比
	f, h, err = c.GetFile("file")    // 获取待比较的文件
	if err != nil {
		// 为了使错误处理看起来比较简洁 使用goto+label统一处理
		resp.Code = models.REQ_ERR
		goto ERR
	}
	if f != nil {
		defer f.Close()
	}

	// 保存位置在 static/upload/, 没有文件夹要先创建
	err = c.SaveToFile("file", models.UPLOAD_PATH+h.Filename)
	if err != nil {
		resp.Code = models.FILE_SAVE_ERR
		goto ERR
	}

	switch checkType {
	case models.CMP_GROUP:
		// 比较分组文件
		if files, err := models.CmpGroup(h.Filename); err != nil {
			resp.Code = models.FILE_DIFF_ERR
			goto ERR
		} else {
			resp.Code = models.SUCCESS
			resp.Ack = files // 期望客户端发送 列表内的文件 进行比对
		}

	case models.CMP_FILE:
		// 比较小文件
		resp.Code=models.SUCCESS
		resp.Diff=models.CmpFile(h.Filename)
	}

	c.Data["json"] = resp
	c.ServeJSON()
	return

ERR:
	resp.Msg = err.Error()
	c.Data["json"] = resp
	c.ServeJSON()
}
