package common

import (
	"encoding/base64"
)

func EncodeBase64(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

func DecodeBase64(src string) string {
	code, err := base64.StdEncoding.DecodeString(src)
	CheckPanic(err)
	return string(code)
}
