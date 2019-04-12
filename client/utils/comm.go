package utils

import "net/http"

var(
	CLIENT *http.Client
	HOST string
	HASH_URL string
	FILE_URL string
)

func Init(){
	CLIENT =&http.Client{}
	HOST="http://localhost:8080/v1/"
	HASH_URL="check/hash"
	FILE_URL="check/file"
}
