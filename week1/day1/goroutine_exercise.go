package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 任务1: 基础goroutine练习
// 创建10个goroutine，每个打印自己的ID后退出
// 使用sync.WaitGroup等待所有goroutine完成

func main() {
	// 练习1: 基础goroutine
	basicGoroutine()

	// 练习2: 带参数的goroutine
	goroutineWithParams()

	// 练习3: 并发安全问题演示
	concurrencyIssueDemo()

	// 扩展练习: 无锁计数器
	atomicCounterDemo()

	// 扩展练习: goroutine池
	goroutinePoolDemo()
}

// basicGoroutine 演示基本的goroutine创建和WaitGroup使用
func basicGoroutine() {
	fmt.Println("=== 练习1: 基础goroutine ===")

	var wg sync.WaitGroup
	const numGoroutines = 10

	for i := 1; i <= numGoroutines; i++ {
		wg.Add(1) // 每启动一个goroutine就增加计数

		go func(id int) {
			defer wg.Done() // goroutine结束时减少计数
			fmt.Printf("Goroutine %d 开始执行\n", id)
			// 模拟一些工作
			fmt.Printf("Goroutine %d 执行完成\n", id)
		}(i) // 注意：这里要传递i的值，而不是引用
	}

	wg.Wait() // 等待所有goroutine完成
	fmt.Println("所有goroutine执行完毕")
	fmt.Println()
}

// goroutineWithParams 演示带参数的goroutine
func goroutineWithParams() {
	fmt.Println("=== 练习2: 带参数的goroutine ===")

	var wg sync.WaitGroup
	messages := []string{"Hello", "World", "Go", "Concurrency", "Is", "Awesome"}

	for i, msg := range messages {
		wg.Add(1)

		go func(index int, message string) {
			defer wg.Done()
			fmt.Printf("Worker %d 处理消息: %s\n", index, message)
		}(i, msg) // 重要：传递当前循环变量的副本
	}

	wg.Wait()
	fmt.Println("所有消息处理完毕")
	fmt.Println()
}

// concurrencyIssueDemo 演示并发安全问题
func concurrencyIssueDemo() {
	fmt.Println("=== 练习3: 并发安全问题演示 ===")

	var wg sync.WaitGroup
	var counter int
	const numWorkers = 100

	// 有问题的版本 - 存在数据竞争
	fmt.Println("有数据竞争的版本:")
	counter = 0
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			counter++ // 这里存在数据竞争！
		}()
	}
	wg.Wait()
	fmt.Printf("期望值: %d, 实际值: %d (可能不正确)\n", numWorkers, counter)

	// 修复方案 1: 使用 sync.Mutex
	fmt.Println("\n=== 修复方案 1: 使用 sync.Mutex ===")
	var mu sync.Mutex
	counter = 0
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()
			counter++ // 受互斥锁保护
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("Mutex版本 - 期望值: %d, 实际值: %d\n", numWorkers, counter)

	// 修复方案 2: 使用 sync/atomic
	fmt.Println("\n=== 修复方案 2: 使用 sync/atomic ===")
	var atomicCounter int64
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&atomicCounter, 1) // 原子操作
		}()
	}
	wg.Wait()
	fmt.Printf("Atomic版本 - 期望值: %d, 实际值: %d\n", numWorkers, atomicCounter)
}

// 扩展练习（可选）：
// 1. 修改concurrencyIssueDemo函数，使用sync.Mutex修复数据竞争
// 2. 尝试使用atomic包实现无锁计数器
func atomicCounterDemo() {
	fmt.Println("=== 扩展练习2: 无锁计数器 ===")

	var counter int64
	const numGoroutines = 1000

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()

	fmt.Printf("最终计数器值: %d (预期: %d)\n", counter, numGoroutines)
}

// 3. 创建一个goroutine池，限制同时运行的goroutine数量为5个
func goroutinePoolDemo() {
	fmt.Println("=== 扩展练习3: goroutine池 ===")

	const numJobs = 20
	const poolSize = 5

	// 使用缓冲channel作为信号量，限制并发数
	semaphore := make(chan struct{}, poolSize)
	var wg sync.WaitGroup

	fmt.Printf("开始处理 %d 个任务，并发限制为 %d\n", numJobs, poolSize)

	for i := 1; i <= numJobs; i++ {
		wg.Add(1)

		// 获取信号量，如果满了会阻塞在这里
		semaphore <- struct{}{}

		go func(id int) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			fmt.Printf("Task %d 开始工作 (当前并发数: %d)\n", id, len(semaphore))
			// 模拟耗时操作
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Task %d 完成工作\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("所有任务处理完毕")
}
