package main

import (
	"fmt"
	"os"
)

/**
# Go Defer
Derfer 用来保证一个函数在程序执行的最后被调用
通常用于清理工作
注意：
- derfer方法通常在当前main函数的最后触发执行
- 即使在定义了defer后触发了错误，也可以被执行
- 如果在defer定义前就触发错误，则不能执行

## 文件操作的例子
- 假设我们要创建一个文件 createFile
- 然后写入数据 writeFile
- 最后关闭文件 使用defer调用closeFile
*/
func main() {

	//文件操作流程
	f := createFile("/tmp/defer.txt")
	//panic("error start")
	defer closeFile(f)
	writeFile(f, "data")

	//即使在defer后触发错误，也可以执行
	//panic("error")

	//main函数结尾
	fmt.Println("main() end")
}

// 创建文件
func createFile(p string) *os.File {
	fmt.Println("creating")
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	return f
}

// 写入文件
func writeFile(f *os.File, d string) {
	fmt.Println("writing")
	_, err := fmt.Fprintln(f, d)
	if err != nil {
		panic(err)
	}
}

// 关闭文件
func closeFile(f *os.File) {
	fmt.Println("closing")
	err := f.Close()
	if err != nil {
		panic(err)
	}
}
