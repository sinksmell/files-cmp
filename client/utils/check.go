package utils

import (
	"github.com/sinksmell/files-cmp/models"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
)

// 发送 文件名 及其 md5值 检查是否一致
func PostHash(hashreq *models.HashRequest) (res *models.Response, err error) {
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
	if req, err = http.NewRequest("POST", HOST+HASH_URL, bytes.NewBuffer(jsonStr)); err != nil {
		return
	}
	// 获取响应
	if resp, err = CLIENT.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	res = &models.Response{}

	body, _ = ioutil.ReadAll(resp.Body)

	if err = json.Unmarshal(body, res); err != nil {
		return
	}
	return
}

// 遍历各个组文件 提交文件名 和对应文件的hash值
func SendGroups(grpFile string) (resp *models.Response, err error) {

	var (
		req  models.HashRequest
		hash string
	)

	if hash, err = models.GetMd5(models.GROUP_PATH + grpFile); err != nil {
		fmt.Println(err)
		return
	}

	req.FileName = grpFile
	req.Hash = hash
	fmt.Println(req)
	resp,err=PostHash(&req)
	return
}
