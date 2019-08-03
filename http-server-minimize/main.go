package main

import (
	"fmt"
	"net/http"
)

/**
三行代码搭建http的web服务
http://127.0.0.1:8080
 */
func main() {
	http.Handle("/",http.FileServer(http.Dir(".")))
	if err:=http.ListenAndServe(":8080",nil);err!=nil{
		fmt.Println(err)
	}
}