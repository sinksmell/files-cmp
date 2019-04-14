package utils

import (
	"github.com/sinksmell/files-cmp/models"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
	"mime/multipart"
	"io"
)

// 发送 文件名 及其 md5值 检查是否一致
func PostHash(hashreq *models.HashRequest,targetURL string) (res *models.Response, err error) {
	var (
		jsonStr []byte
		req     *http.Request
		resp    *http.Response
		body    []byte
	)
	// 序列化为json
	if jsonStr, err = json.Marshal(hashreq); err != nil {
		return
	}
	// 创建请求
	if req, err = http.NewRequest("POST", targetURL, bytes.NewBuffer(jsonStr)); err != nil {
		return
	}
	// 获取响应
	if resp, err = CLIENT.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	// 读取响应中的数据 并解析json
	res = &models.Response{}
	body, _ = ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, res); err != nil {
		return
	}
	return
}

// 遍历各个组文件 提交文件名 和对应文件的hash值
func SendGrpMd5(grpFile ,targetURL string) (resp *models.Response, err error) {
	var (
		req  models.HashRequest
		hash string
	)
	// 计算分组文件的hash值
	if hash, err = models.GetMd5(models.GROUP_PATH + grpFile); err != nil {
		fmt.Println(err)
		return
	}
	// 构造请求
	req.FileName = grpFile
	req.Hash = hash
	// 发送 使用http发送json数据,比较分组的md5值
	resp, err = PostHash(&req,targetURL)
	return
}

// 提交文件
func PostFile(filename ,targetURL ,filepath,cmpType string) (res *models.Response,err error) {

	var (
		bodyBuf    *bytes.Buffer
		bodyWriter *multipart.Writer
		fileWriter io.Writer
		content []byte
		jsonData []byte  // http response body中的json数据
	)
	// 初始化
	bodyBuf = &bytes.Buffer{}
	bodyWriter = multipart.NewWriter(bodyBuf)

	// 创建文件上传表单
	if fileWriter, err = bodyWriter.CreateFormFile("file", filename); err != nil {
		return
	}

	// 打开文件
	if content, err = ioutil.ReadFile(filepath+filename);err!=nil{
		return
	}

	// 拷贝文件内容
	if _,err=io.Copy(fileWriter,bytes.NewReader(content));err!=nil{
		return
	}
	// 设置请求头 发送请求
	contentType:=bodyWriter.FormDataContentType()
	// 设置文件对比类型
	bodyWriter.WriteField("type",cmpType)
	// 关闭Writer 使用defer会出错
	bodyWriter.Close()

	resp,err:=http.Post(targetURL,contentType,bodyBuf)
	if err!=nil{
		return
	}
	defer  resp.Body.Close()

	// 获取返回的json数据
	if jsonData,err=ioutil.ReadAll(resp.Body);err!=nil{
		return
	}

	// 解析返回的json数据
	res,err=ParseResp(jsonData)
	return
}




// 根据HTTP response 解析到 models.Response
func ParseResp(bytes []byte)(res *models.Response,err error){
	res = &models.Response{}
	if err = json.Unmarshal(bytes, res); err != nil {
		return
	}
	return
}


//

func PostFile2(filename string, targetUrl string) (res *models.Response , err error) {
	var(
		bodyBuf *bytes.Buffer
		bodyWriter *multipart.Writer
		fileWriter io.Writer
		content []byte
	)

	bodyBuf = &bytes.Buffer{}
	bodyWriter = multipart.NewWriter(bodyBuf)

	//构造一个请求表单
	if fileWriter, err = bodyWriter.CreateFormFile("file", filename);err!=nil{
		fmt.Println("error writing to buffer")
		return
	}
	// 读取文件内容
	if content, err=ioutil.ReadFile(filename);err != nil {
		fmt.Println("error opening file")
		return
	}

	//拷贝文件内容到表单
	_, err = io.Copy(fileWriter, bytes.NewReader(content))
	if err != nil {
		return
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	res,err=ParseResp(resp_body)
	return
}