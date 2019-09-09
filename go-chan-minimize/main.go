package main

import (
	"fmt"
)

func Count(ch chan int, s int) {
	fmt.Println("counting: ", s)
	ch <- 1
}

func main() {
	num := 10000
	chs := make([]chan int, num)

	for i := 0; i < num; i++ {
		chs[i] = make(chan int)
		go Count(chs[i], i)
	}

	for _, ch := range chs {
		<-ch
	}
}
