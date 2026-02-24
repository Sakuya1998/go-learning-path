//go:build example_select_demo
// +build example_select_demo

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// selectBasics 演示select基本用法
func selectBasics() {
	fmt.Println("=== 练习1: select基础 ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// 启动两个goroutine，分别向不同的Channel发送数据
	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		ch1 <- "来自ch1的消息"
	}()

	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		ch2 <- "来自ch2的消息"
	}()

	// 使用select等待第一个可用的Channel
	select {
	case msg1 := <-ch1:
		fmt.Printf("收到: %s\n", msg1)
	case msg2 := <-ch2:
		fmt.Printf("收到: %s\n", msg2)
	}
}

// selectWithTimeout 演示select超时控制
func selectWithTimeout() {
	fmt.Println("\n=== 练习2: select超时控制 ===")

	ch := make(chan string)

	go func() {
		// 模拟耗时操作（2秒）
		time.Sleep(2 * time.Second)
		ch <- "操作结果"
	}()

	select {
	case result := <-ch:
		fmt.Printf("成功: %s\n", result)
	case <-time.After(1 * time.Second):
		fmt.Println("错误: 操作超时（1秒）")
	}
}

// selectMultiple 演示同时处理多个Channel
func selectMultiple() {
	fmt.Println("\n=== 练习3: 处理多个Channel ===")

	ch1 := make(chan int)
	ch2 := make(chan string)
	done := make(chan bool)

	// 数据生产者1
	go func() {
		for i := 1; i <= 3; i++ {
			time.Sleep(400 * time.Millisecond)
			ch1 <- i
		}
	}()

	// 数据生产者2
	go func() {
		messages := []string{"A", "B", "C", "D"}
		for _, msg := range messages {
			time.Sleep(300 * time.Millisecond)
			ch2 <- msg
		}
	}()

	// 定时器：3秒后结束
	go func() {
		time.Sleep(3 * time.Second)
		done <- true
	}()

	fmt.Println("开始监听多个Channel（3秒）:")
	count1, count2 := 0, 0

	for {
		select {
		case num := <-ch1:
			count1++
			fmt.Printf("  从ch1收到数字: %d\n", num)
		case str := <-ch2:
			count2++
			fmt.Printf("  从ch2收到字符: %s\n", str)
		case <-done:
			fmt.Printf("\n结束。共收到 %d 个数字，%d 个字符\n", count1, count2)
			return
		}
	}
}

// selectDefault 演示非阻塞操作
func selectDefault() {
	fmt.Println("\n=== 练习4: 非阻塞操作（default） ===")

	ch := make(chan string)

	// 尝试非阻塞接收
	select {
	case msg := <-ch:
		fmt.Printf("收到: %s\n", msg)
	default:
		fmt.Println("没有数据可接收，继续执行其他操作")
	}

	// 尝试非阻塞发送
	select {
	case ch <- "测试消息":
		fmt.Println("发送成功")
	default:
		fmt.Println("Channel未准备好接收，发送失败")
	}

	// 现在启动接收者
	go func() {
		time.Sleep(100 * time.Millisecond)
		if msg, ok := <-ch; ok {
			fmt.Printf("Goroutine收到: %s\n", msg)
		}
	}()

	// 再次尝试发送
	time.Sleep(50 * time.Millisecond)
	select {
	case ch <- "Hello":
		fmt.Println("第二次发送成功")
	default:
		fmt.Println("第二次发送失败")
	}

	time.Sleep(200 * time.Millisecond) // 等待goroutine完成
	close(ch)
}

// workerPoolWithSelect 演示使用select的工作池模式
func workerPoolWithSelect() {
	fmt.Println("\n=== 练习5: 工作池与select ===")

	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	done := make(chan bool)

	// 创建工作池
	for w := 1; w <= numWorkers; w++ {
		go func(workerID int) {
			for job := range jobs {
				fmt.Printf("Worker %d 开始处理工作 %d\n", workerID, job)
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				result := job * 2 // 简单处理：乘以2
				results <- result
				fmt.Printf("Worker %d 完成工作 %d -> %d\n", workerID, job, result)
			}
		}(w)
	}

	// 发送工作
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
			fmt.Printf("发送工作: %d\n", j)
			time.Sleep(100 * time.Millisecond)
		}
		close(jobs)
	}()

	// 收集结果
	go func() {
		for i := 0; i < numJobs; i++ {
			result := <-results
			fmt.Printf("收到结果: %d\n", result)
		}
		done <- true
	}()

	// 等待完成
	<-done
	fmt.Println("所有工作完成")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("第1周第2天: select多路复用")
	fmt.Println("==========================\n")

	selectBasics()
	selectWithTimeout()
	selectMultiple()
	selectDefault()
	workerPoolWithSelect()

	fmt.Println("\n=== select使用场景总结 ===")
	fmt.Println("1. 超时控制: 使用 time.After() 防止操作永久阻塞")
	fmt.Println("2. 非阻塞操作: 使用 default 分支避免阻塞")
	fmt.Println("3. 多路复用: 同时监听多个Channel，处理最先就绪的")
	fmt.Println("4. 工作池模式: 结合Channel和select构建并发系统")
	fmt.Println("5. 优雅关闭: 使用select实现程序的优雅退出")
}