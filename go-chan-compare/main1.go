package main

import (
	"fmt"
	"time"
)

func main() {
	num := 10000

	t1 := time.Now()
	execute1(num)
	t1s := time.Since(t1)

	fmt.Println("serial func", t1s)
}

// 模拟耗时任务
func job1(s int) {
	time.Sleep(1 * time.Millisecond)
	fmt.Println("serial execute job  : ", s)
}

// 串行循环执行
func execute1(num int) {
	for i := 0; i < num; i++ {
		job1(i)
	}
}
