package models


const(
    FILE_PATH="./file/"
)


type  Response  struct{
    Code int	`json:"code"`
    Msg string `json:"msg"`
}


// 用于检测hash值是否相同的请求结构
type  HashRequest  struct{
    FileName string `json:"fname"`
    Hash string `json:"hash"`
}