package main

import (
	"fmt"
	"os"
)

/**
# Go Exit
使用os.Exit可以给定一个状态
然后立刻退出程序运行
*/
func main() {

	// 当使用`os.Exit`的时候defer操作不会被运行
	defer fmt.Println("main() end")

	// 退出程序并设置退出状态值
	os.Exit(100)
}
