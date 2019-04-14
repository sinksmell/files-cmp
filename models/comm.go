package models

import (
	"github.com/astaxie/beego"
	"math"
	"strings"
)

const (
	FILE_PATH   = "./file/"          // 小文件存放路径
	UPLOAD_PATH = "./static/upload/" // 上传文件保存路径
	GROUP_PATH  = "./group/"         // 分组文件保存路径
	GROUP_PRE   = "group_"           // 分组文件 文件名的前缀
	GROUP_SUF   = ".txt"             // 分组文件 文件名的后缀

	// 文件对比类型
	CMP_GROUP = "cmp_group"
	CMP_FILE  = "cmp_file"

	// 服务器响应码
	EQUAL         = 0000 // 对应的文件的md5相同
	SUCCESS       = 0001 // 文件处理成功
	NOT_EQUAL     = 1000 // 对应的文件的md5不同
	REQ_ERR       = 2000 // 请求出错
	FILE_SAVE_ERR = 3001 //文件保存出错
	FILE_DIFF_ERR = 3002 // 文件对比时出错

)

var (
	FILE_CNT int      // 文件数量
	Files    []string // 文件名集合
	GRP_SIZE int      // 分组大小
)

func InitFiles(filepath string) {
	var (
		err error
	)
	// 获取文件列表
	if Files, err = GetAllFiles(filepath); err != nil {
		beego.BeeLogger.Info(err.Error())
		return
	}
	// 获取文件数量
	FILE_CNT = len(Files)
	// 计算分组大小
	GRP_SIZE = int(math.Sqrt(float64(1.0 * FILE_CNT)))
}

type Response struct {
	Code int         `json:"code"` // 状态码
	Msg  string      `json:"msg"`  // 错误信息
	Ack  []string    `json:"ack"`  // 期望 收到的文件列表
	Diff *DiffResult `json:"diff"` // 文件对比结果
}

// 用于检测hash值是否相同的请求结构
type HashRequest struct {
	FileName string `json:"fname"` // 文件名
	Hash     string `json:"hash"`  // 文件对应的MD5值
}

// 文件对比结果
type DiffResult struct {
	FileName      string `json:"file_name"`  // 对比的文件
	ClientContent string `json:"c_content"`  // 客户端文件特有的内容
	ServerContent string `json:"s_content"`  // 服务端文件特有的内容
	Changes       string `json:"changes"`    // 文件的变化记录
	ColorText     string `json:"color_text"` // 高亮显示出两个文本不同的内容
}

func (res *DiffResult) String() string {
	builder := strings.Builder{}
	builder.WriteString("对比小文件为:\t")
	builder.WriteString(res.FileName)
	builder.WriteString("\n本地文件特有内容:\n")
	builder.WriteString(res.ClientContent)
	builder.WriteString("\n远程文件特有内容:\n")
	builder.WriteString(res.ServerContent)
	builder.WriteString("\n文件变化内容适配:\n")
	builder.WriteString(res.Changes)
	builder.WriteString("\n高亮显示不同内容:\n")
	builder.WriteString(res.ColorText)
	return builder.String()
}
