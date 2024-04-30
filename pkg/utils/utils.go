package utils

import (
	"bytes"
	"context"
	"go.elastic.co/apm"
	"strings"
)

func RandomString(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(rune(RandInt(65, 90))) != temp {
			temp = string(rune(RandInt(65, 90)))
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}

func CopyContext(ctx context.Context) context.Context {
	return apm.DetachedContext(ctx)
}

func GetLastName(path string) string {
	var splicedList = strings.Split(path, "/")
	//if len(splicedList) == 0 {
	//	return ""
	//}
	return splicedList[len(splicedList)-1]
}
