package common

import (
	"bytes"
	"fmt"
	"strings"
)

//利用buf拼接字符串
func JoinStringBuf(split string, code ...string) string {
	buf := bytes.Buffer{}
	for _, s := range code {
		buf.WriteString(s)
		buf.WriteString(split)
	}
	return buf.String()
}

//最高效的拼接字符串
func JoinString(split string, code ...string) string {
	return strings.Join(code, split)
}

func StripString(str string) string {
	return strings.TrimSpace(str)
}

func FormatString(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
