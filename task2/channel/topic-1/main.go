/*
	[Channel]
	题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
		考察点 ：通道的基本使用、协程间通信。
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	// 生产者协程：生成1-10的整数并发送到通道
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	// 消费者协程：从通道接收数据并打印
	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Printf("接收到数据：%d\n", num)
		}
	}()

	wg.Wait()
}
