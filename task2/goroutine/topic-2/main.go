/*
	[Goroutine]
	题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
		考察点 ：协程原理、并发任务调度。
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// 定义任务结构：包含任务名称和任务执行逻辑
type Task struct {
	Name string
	Func func()
}

// 定义任务执行结构：记录任务名称、执行耗时、错误信息
type TaskResult struct {
	TaskName string
	Duration time.Duration
	Err      error
}

// 任务调度器核心函数：接收任务列表，并发执行并返回统计结果
func ScheduleTasks(tasks []Task) []TaskResult {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results = make([]TaskResult, 0, len(tasks))
	)
	wg.Add(len(tasks))
	for _, task := range tasks {
		task := task
		go func() {
			defer wg.Done()
			result := TaskResult{TaskName: task.Name}
			//记录任务开始时间
			startTime := time.Now()
			//捕获任务执行中的panic（避免单个任务崩溃影响整体）
			defer func() {
				//计算执行耗时（defer中获取时间，确保任务执行完才统计）
				result.Duration = time.Since(startTime)
				if err := recover(); err != nil {
					result.Err = fmt.Errorf("执行崩溃：%v", err)
				}
				//并发安全写入结果：加锁避免多个协程同时修改results
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}()
			task.Func()
		}()
	}
	wg.Wait()
	return results
}

func main() {
	testTasks := []Task{
		{
			Name: "睡眠100ms任务",
			Func: func() {
				time.Sleep(100 * time.Millisecond)
			},
		},
		{
			Name: "计算密集型任务",
			Func: func() {
				sum := 0
				for i := 0; i < 100000000; i++ { // 模拟计算耗时
					sum += i
				}
				fmt.Printf("计算任务结果：%d\n", sum)
			},
		},
		{
			Name: "故意panic任务",
			Func: func() {
				panic("模拟任务执行失败") // 测试错误捕获
			},
		},
	}
	fmt.Println("=== 任务调度器启动 ===")
	totalStartTime := time.Now()
	results := ScheduleTasks(testTasks)
	totalDuration := time.Since(totalStartTime)

	fmt.Printf("\n=== 所有任务执行完毕 ===")
	fmt.Printf("\n总耗时：%v\n", totalDuration)
	fmt.Println("\n任务详情：")
	for _, res := range results {
		if res.Err != nil {
			fmt.Printf("- %s：失败 | 错误：%v\n", res.TaskName, res.Err)
		} else {
			fmt.Printf("- %s：成功 | 执行时间：%v\n", res.TaskName, res.Duration)
		}
	}
}
