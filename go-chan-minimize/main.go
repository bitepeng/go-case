package main

import "fmt"

func Count(ch chan int, s int) {
	fmt.Println("counting: ", s)
	ch <- 1
}

func main() {
	chs := make([]chan int, 1000)

	for i := 0; i < 1000; i++ {
		chs[i] = make(chan int)
		go Count(chs[i], i)
	}

	for _, ch := range chs {
		<-ch
	}
}
