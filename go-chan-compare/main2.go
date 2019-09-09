package main

import (
	"fmt"
	"time"
)

func main() {
	num := 10000

	t1 := time.Now()
	execute2(num)
	t1s := time.Since(t1)

	fmt.Println("concurrent func: ", t1s)
}

// 模拟耗时任务
func job2(ch chan int, s int) {
	time.Sleep(1 * time.Millisecond)
	fmt.Println("concurrent execute job : ", s)
	ch <- 1
}

// 协程并发执行
func execute2(num int) {
	chs := make([]chan int, num)
	for i := 0; i < num; i++ {
		chs[i] = make(chan int)
		go job2(chs[i], i)
	}
	for _, ch := range chs {
		<-ch
	}
}
