package test

import (
	"testing"
	"github.com/sinksmell/files-cmp/models"
	"fmt"
)

func TestGetFiles(t*testing.T){

	if names,err:=models.GetAllFiles("../file/");err!=nil{
		t.Fatal(err)
	}else{
		fmt.Println(len(names))
		for _, name := range names {
			fmt.Println(name)
		}
	}

}
