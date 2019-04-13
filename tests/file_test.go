package test

import (
	"testing"
	"github.com/sinksmell/files-cmp/models"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// 测试获取指定目录下的文件列表
func TestGetFiles(t *testing.T) {

	if names, err := models.GetAllFiles("../file/"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(len(names))
		for _, name := range names {
			fmt.Println(name)
		}
	}

}

// 测试对文件进行分组
func TestDivide(t *testing.T) {

	models.InitFiles("../file/")

	models.Divide("../file/", "../group/")

}

// 测试比较两个文件 获取文件的不同内容
func TestFileCmp(t *testing.T) {

	dmp := diffmatchpatch.New()
	if diffs, err := models.FileDiff("../group/group_9.txt", "../static/upload/group_9.txt"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(dmp.DiffPrettyText(diffs))
		//fmt.Println(dmp.DiffToDelta(diffs))
	}

	/*

	diffs := dmp.DiffMain(text1, text2, true)
	patches := dmp.PatchMake(text1, diffs)
	patchText := dmp.PatchToText(patches)
	patchesFromText, _ := dmp.PatchFromText(patchText)
	text2FromPatches, _ := dmp.PatchApply(patchesFromText, text1)
	*/

}

// 测试获取文件的动态变化
func TestParch(t *testing.T) {
	var (
		dmp   *diffmatchpatch.DiffMatchPatch
		diffs []diffmatchpatch.Diff
		patches []diffmatchpatch.Patch
		//err   error
		text1 = "hello world go!"
		text2 = "hello world java!"
	)

	dmp = diffmatchpatch.New()
	diffs = dmp.DiffMain(text1, text2, false)
	patches = dmp.PatchMake(text1, diffs)
	patchText:=dmp.PatchToText(patches)
	fmt.Println(patchText)
	patchesFromText, _ := dmp.PatchFromText(patchText)
	fmt.Println(patchesFromText)
	text2FromPatches, _ := dmp.PatchApply(patchesFromText, text1)
	fmt.Println(text2FromPatches)
}

// 测试对比两个文件的MD5值
func TestCmpMd5(t*testing.T){
	res,_:= models.CmpMd5("../group/group_1.txt","9190eab0e70157e598c8ad3247aab38d")
	fmt.Println(res)
}