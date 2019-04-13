package utils

import "net/http"

var(
	CLIENT *http.Client

)

func init(){
	// 初始化单例
	CLIENT=&http.Client{}
}
