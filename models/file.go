package models

import (
	"bufio"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
		grp     *os.File
		grpId   = 0    // 分组文件的id
		remain  = 0    // 小组文件剩余量
		grpName string // 分组文件名
		content string // 文件内容
		hash    string // 文件md5值
		err     error
	)
	// Files是包内的一个全局变量 存放小文件目录
	for _, fname := range Files {
		if hash, err = GetMd5(filePath + fname); err != nil {
			fmt.Println(err)
			continue
		}
		if remain == 0 {
			// 剩余数量为0时再创建一个分组
			remain = GRP_SIZE
			grpId++
			grpName = fmt.Sprintf("%s%d%s", GROUP_PRE, grpId, GROUP_SUF)
			if grp != nil {
				// 关闭上一个分组文件
				grp.Close()
			}
			// 新建一个分组文件
			if grp, err = os.OpenFile(grpPath+grpName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666); err != nil {
				//fmt.Println(err)
				return
			}
		}

		// 写入记录
		content = fmt.Sprintf("%s : %s\n", fname, hash)
		grp.WriteString(content)
		remain--

	}

	// 关闭最后一个分组文件
	grp.Close()

}

// Client上的小文件对比Server上的小文件
func TextFileDiff(cFile, sFile string) (diffResult *DiffResult) {

	var (
		text1     string                         //  Server的文件内容
		text2     string                         //  Client的文件内容
		bytes1    []byte                         // Server文件内容
		bytes2    []byte                         // Client文件内容
		content1  strings.Builder                // Server的文件特有内容
		content2  strings.Builder                // Client的文件特有内容
		dmp       *diffmatchpatch.DiffMatchPatch // 文本对比器
		diffs     []diffmatchpatch.Diff          // 对比结果
		patchText string                         // 记录文件变化
		err       error
	)
	// 初始化
	diffResult = &DiffResult{}
	dmp = diffmatchpatch.New()

	// 读取Server文件 以Server文件为基准 进行对比
	if bytes1, err = ioutil.ReadFile(sFile); err != nil {
		// 文件打开失败或者不存在该文件
		bytes1 = []byte{}
	}
	text1 = string(bytes1)

	// 读取Client文件
	if bytes2, err = ioutil.ReadFile(cFile); err != nil {
		// 文件打开失败或者不存在该文件
		bytes2 = []byte{}
	}
	text2 = string(bytes2)
	//// 获取两个文件内容比较结果
	//diffs = dmp.DiffMain(text1, text2, true)
	//// 将Server端文件内容作为patch基准
	//patches = dmp.PatchMake(text1, diffs)
	//patchText := dmp.PatchToText(patches)
	// 获取文件变化 及patch
	diffs, patchText = diffAndPatch(text1, text2, true)

	// 记录Server与Client特有的文件内容
	for _, diff := range diffs {
		if diff.Type == diffmatchpatch.DiffInsert {
			content2.WriteString(diff.Text)
			content2.WriteString("\n")
		} else if diff.Type == diffmatchpatch.DiffDelete {
			content1.WriteString(diff.Text)
			content1.WriteString("\n")
		}
	}

	// 构造对比结果
	diffResult.FileName = strings.TrimPrefix(sFile, FILE_PATH)
	diffResult.Changes = patchText
	diffResult.ServerContent = content1.String()
	diffResult.ClientContent = content2.String()
	diffResult.ColorText = dmp.DiffPrettyText(diffs)
	return
}

// 对比对应文件的 md5
func CmpMd5(filepath string, hash string) (equal bool, err error) {
	var (
		// 本地文件的MD5值
		_hash string
	)
	// 计算本地文件的md5值
	if _hash, err = GetMd5(filepath); err != nil {
		return
	} else {
		if _hash == hash {
			equal = true
		}
	}
	return
}

// 从两个分组文件中找出不同文件的文件名集合
func GetDiffFiles(cFile, sFile string) (flist []string, err error) {

	var (
		file1     *os.File          // server端文件指针
		file2     *os.File          // client端文件指针
		sc        *bufio.Scanner    // 带缓冲的Scanner效率较高
		serverMap map[string]string // 用于寻找Md5值不同的文件
	)
	// 记录server端 文件名:md5值
	serverMap = make(map[string]string)
	flist = make([]string, 0)
	// 读取Server端的文件
	if file1, err = os.Open(sFile); err != nil {
		// 说明服务端可能没有这个文件
		// 直接不处理 只分析Client传过来的文件
		//return
	}

	if file1 != nil {
		// 文件存在而且能正常读取
		defer file1.Close()
		sc = bufio.NewScanner(file1)
		for sc.Scan() {
			record := sc.Text()
			if len(record) == 0 {
				break
			}
			vals := strings.Split(record, " : ")
			if len(vals) == 2 {
				// 保存 文件名:md5值
				serverMap[vals[0]] = vals[1]
			}
		}
	}
	//	读取 Client传过来的文件
	if file2, err = os.Open(cFile); err != nil {
		//	fmt.Println(err)
		return
	}
	defer file2.Close()

	sc = bufio.NewScanner(file2)
	for sc.Scan() {
		record := sc.Text()
		if len(record) == 0 {
			break
		}

		vals := strings.Split(record, " : ")
		if len(vals) < 2 {
			continue
		}
		// 如果对应的md5值不同或者 无记录则加入结果集中
		if hash, ok := serverMap[vals[0]]; !ok || hash != vals[1] {
			flist = append(flist, vals[0])
		}

	}

	return
}

// 比较对应的两个分组文件
func CmpGroup(fileName string) (files []string, err error) {
	files, err = GetDiffFiles(UPLOAD_PATH+fileName, GROUP_PATH+fileName)
	return
}

// 比较对应的两个小文件
func CmpFile(fileName string) (diff *DiffResult) {

	isText := isTextFile(UPLOAD_PATH + fileName)
	if isText {
		return TextFileDiff(UPLOAD_PATH+fileName, FILE_PATH+fileName)
	}

	return NonTextFileDiff(UPLOAD_PATH+fileName, FILE_PATH+fileName)
}

// 获取文件类型 判断文件是否为文本文件
func isTextFile(fileName string) (flag bool) {
	var (
		cmd    *exec.Cmd
		output []byte
	)
	cmd = exec.Command("file", fileName)
	output, _ = cmd.CombinedOutput()
	// fmt.Println(string(output))
	flag = strings.Contains(string(output), "text")
	return
}

// 比较非文本文件 为了减少输出 只比较前512个字节
func NonTextFileDiff(cFile, sFile string) (diffResult *DiffResult) {

	var (
		// Server
		file1    *os.File
		buff1    []byte
		text1    string
		content1 strings.Builder // Server的文件特有内容
		// Client
		file2     *os.File
		buff2     []byte
		text2     string
		content2  strings.Builder                // Client的文件特有内容
		dmp       *diffmatchpatch.DiffMatchPatch // 文本对比器
		diffs     []diffmatchpatch.Diff          // 对比结果
		patchText string                         // 记录文件变化
	)
	// 初始化
	diffResult = &DiffResult{}
	dmp = diffmatchpatch.New()
	buff1 = make([]byte, CMP_SIZE)
	buff2 = make([]byte, CMP_SIZE)

	// 打开Server端文件
	file1, _ = os.Open(sFile)
	if file1 != nil {
		defer file1.Close()
		file1.Read(buff1)
		text1 = string(buff1)
	} else {
		text1 = ""
	}
	// 打开Client端文件
	file2, _ = os.Open(cFile)
	if file2 != nil {
		defer file2.Close()
		file2.Read(buff2)
		text2 = string(buff2)
	} else {
		text2 = ""
	}

	//	beego.BeeLogger.Debug("%+v\n",text1)
	//	beego.BeeLogger.Debug("%+v\n",text2)

	// 对比文本内容
	diffs, patchText = diffAndPatch(text1, text2, false)

	// 记录Server与Client特有的文件内容
	for _, diff := range diffs {
		if diff.Type == diffmatchpatch.DiffInsert {
			content2.WriteString(diff.Text)
			content2.WriteString("\n")
		} else if diff.Type == diffmatchpatch.DiffDelete {
			content1.WriteString(diff.Text)
			content1.WriteString("\n")
		}
	}

	// 构造对比结果
	diffResult.FileName = strings.TrimPrefix(sFile, FILE_PATH)
	diffResult.Changes = patchText
	diffResult.ServerContent = content1.String()
	diffResult.ClientContent = content2.String()
	diffResult.ColorText = dmp.DiffPrettyText(diffs)

	return
}

// 字符串diff与patch 只关注于字符串处理
func diffAndPatch(text1, text2 string, isText bool) (diffs []diffmatchpatch.Diff, patchText string) {
	// 初始化比较器
	dmp := diffmatchpatch.New()
	// 获取两个文件内容比较结果
	diffs = dmp.DiffMain(text1, text2, true)
	// 只对文本数据处理patch变化
	if isText {
		// 将text1作为patch基准
		patches := dmp.PatchMake(text1, diffs)
		patchText = dmp.PatchToText(patches)
	}
	return
}
