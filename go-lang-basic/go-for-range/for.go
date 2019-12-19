package main

import "fmt"

/**
# Go For
for循环是Go语言唯一的循环结构

## 1.单一条件循环
	类似其他语言的while循环
## 2.初始/条件/循环
	初始化/条件判断/循环后条件变化
## 3.无条件for循环
	死循环，除非break跳出，或return返回
*/
func main() {

	//1.单一条件循环
	i := 1
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}

	//2.初始/条件/循环
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}

	//3.无条件for循环
	for {
		fmt.Println("loop")
		break
	}
}
