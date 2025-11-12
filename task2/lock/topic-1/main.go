/*
	[锁机制]
	题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
		考察点 ： sync.Mutex 的使用、并发数据安全。
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int // 定义共享资源：计数器（多个协程会并发修改）
	var mu sync.Mutex
	var wg sync.WaitGroup

	const (
		numGoroutines = 10   // 协程数量
		numIncrements = 1000 // 每个协程的递增次数
	)

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			// 每个协程执行1000次递增
			for j := 0; j < numIncrements; j++ {
				// 核心：加锁保护临界区（对counter的读写操作）
				mu.Lock()
				// 临界区：同一时间只有一个协程能执行
				counter++
				// 解锁：释放锁，让其他协程可以进入临界区
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("最终计数器值：%d\n", counter)
}
