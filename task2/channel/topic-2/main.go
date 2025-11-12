/*
	[Channel]
	题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
		考察点 ：通道的缓冲机制。
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	// 创建带缓冲通道：缓冲区大小设为10
	ch := make(chan int, 10)
	var wg sync.WaitGroup
	// 注册2个协程（生产者+消费者），主协程等待两者完成
	wg.Add(2)
	// 生产者协程：向通道发送1-100的整数
	go func() {
		defer wg.Done()
		defer close(ch)

		for i := 1; i <= 100; i++ {
			ch <- i
			fmt.Printf("生产者发送：%d（当前缓冲区剩余容量：%d）\n", i, cap(ch)-len(ch))
		}
	}()
	// 消费者协程：从通道接收数据并打印
	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Printf("消费者接收：%d\n", num)
		}
	}()

	wg.Wait()
}
