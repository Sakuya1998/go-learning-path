package main

import (
	"fmt"
	"sync"
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

	// 思考题：为什么会出现数据竞争？
	// 提示：多个goroutine同时读写同一个变量
	fmt.Println("\n思考题：如何修复这个数据竞争问题？")
	fmt.Println("提示：可以使用sync.Mutex或atomic包")
}

// 扩展练习（可选）：
// 1. 修改concurrencyIssueDemo函数，使用sync.Mutex修复数据竞争
// 2. 尝试使用atomic包实现无锁计数器
// 3. 创建一个goroutine池，限制同时运行的goroutine数量为5个