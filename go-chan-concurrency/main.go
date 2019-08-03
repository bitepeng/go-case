package main

import (
	"fmt"
	"math/rand"
	"time"
)

/**
原生支持并发（go协程+通道）
*/
func main() {

	channel := make(chan string)
	go producer("cat:", channel)
	go producer("dog:", channel)
	customer(channel)

}

/**
数据生产者
*/
func producer(header string, channel chan<- string) {
	for {
		//将字符串发给通道
		channel <- fmt.Sprintf("%s: %v", header, rand.Int31())
		//等待一秒
		time.Sleep(time.Second)
	}
}

/**
数据消费者
*/
func customer(channel <-chan string) {
	for {
		//从通道获取字符串
		message := <-channel
		//打印输出
		fmt.Println(message)
	}
}
