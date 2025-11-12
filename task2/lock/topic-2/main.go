/*
	[锁机制]
	题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
		考察点 ：原子操作、并发数据安全。
*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {

	var counter int32 // 定义原子计数器：必须使用atomic支持的类型（int32/int64/uint32等）
	var wg sync.WaitGroup

	const (
		numGoroutines = 10   // 协程数量
		numIncrements = 1000 // 每个协程递增次数
	)

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIncrements; j++ {
				atomic.AddInt32(&counter, 1) // 参数1：计数器指针（必须传指针，直接操作原始内存），参数2：递增步长（这里是1）
			}
		}()
	}

	wg.Wait()
	// 原子读取最终值（确保拿到最新的内存值，避免缓存可见性问题）
	sum := atomic.LoadInt32(&counter)
	fmt.Printf("最终计数器值：%d\n", sum)
}
