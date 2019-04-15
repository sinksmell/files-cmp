package main

import (
	"fmt"
	"github.com/sinksmell/files-cmp/client/utils"
	"github.com/sinksmell/files-cmp/models"
)

const (
	HOST     string = "http://localhost:8080/v1/check"
	HASH_URL string = "/hash"
	FILE_URL string = "/file"
)

var (
	groups    []string //分组文件列表
	diffFiles []string // 需要对比的小文件集合
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
	diffFiles = make([]string, 0)
}

func main() {

	var (
		resp *models.Response // 文件对比结果
		msg  string           // 输出的信息
		err  error
	)

	for _, grp := range groups {
		if resp, err = utils.SendGrpMd5(grp, HOST+HASH_URL); err != nil {
			fmt.Println(err)
			continue
		}
		switch resp.Code {
		case models.EQUAL:
			// 分组文件的Md5值相同 不需要处理
			msg = setColor(grp+"\t分组文件内容相同!", 0, 0, 37)
			fmt.Println(msg)
		case models.NOT_EQUAL:
			// 如果分组文件的Md5值不同 则把该分组文件发送过去 找到需要对比的小文件
			// 红色高亮显示
			msg = setColor(grp+"\t分组文件内容不同!", 0, 0, 31)
			fmt.Println(msg)
			handleDiffGroup(grp)
			fmt.Println("")
		case models.REQ_ERR:
			fmt.Println("Request Err!")
		}
	}
}

// 处理分组文件Md5值不同
func handleDiffGroup(grpFile string) {

	var (
		res *models.Response
		err error
	)

	// 发送分组文件  文件对比类型为 cmp_group 即对比响应的分组文件
	if res, err = utils.PostFile(grpFile, HOST+FILE_URL, models.GROUP_PATH, models.CMP_GROUP); err != nil {
		fmt.Println(err)
		return
	}
	// 获取响应中的 期望对比的文件集合
	if len(res.Ack) > 0 {
		for _, fname := range res.Ack {
			// 遍历集合 把小文件发送过去 对比文件内容
			if _res, err := utils.PostFile(fname, HOST+FILE_URL, models.FILE_PATH, models.CMP_FILE); err != nil {
				fmt.Println(err)
				fmt.Println("发送错误!")
				continue
			} else {
				fmt.Println(_res.Diff)
			}

		}

	}
}

// 设定颜色打印
func setColor(msg string, conf, bg, text int) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, conf, bg, text, msg, 0x1B)
}
