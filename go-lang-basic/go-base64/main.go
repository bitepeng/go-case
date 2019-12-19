package main

import (
	"encoding/base64"
	"fmt"
)

/**
# base64编码解码
Go支持标准和兼容URL两种base64编码
兼容URL模式主要针对URL传递进行优化

## 1.标准模式
- StdEncoding.EncodeToString([]byte)
- StdEncoding.DecodeString([]byte)

## 2.兼容URL模式
- URLEncoding.EncodeToString([]byte)
- URLEncoding.DecodeString([]byte)
*/
func main() {

	// 原始字符串
	data := "abc123!?$*&()-=@'~"

	// 1.标准模式
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	sDec := base64.StdEncoding.EncodeToString([]byte(sEnc))
	fmt.Println("base64.StdEncoding: ", data, sEnc, sDec)

	// 2.兼容URL模式
	uEnc := base64.URLEncoding.EncodeToString([]byte(data))
	uDec := base64.URLEncoding.EncodeToString([]byte(uEnc))
	fmt.Println("base64.URLEncoding: ", data, uEnc, uDec)

}
