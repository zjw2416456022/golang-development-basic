/*
	[Goroutine]
	题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
		考察点 ： go 关键字的使用、协程的并发执行。
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			fmt.Println("奇数：", i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			fmt.Println("偶数：", i)
		}
	}()
	wg.Wait()
}

// func main() {
// 	go func() {
// 		for i := 1; i <= 10; i += 2 {
// 			fmt.Println("奇数", i)
// 		}
// 	}()
// 	go func() {
// 		for i := 2; i <= 10; i += 2 {
// 			fmt.Println("偶数", i)
// 		}
// 	}()
// 	time.Sleep(5 * time.Second)
// }
