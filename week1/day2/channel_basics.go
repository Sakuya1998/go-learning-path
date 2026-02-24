//go:build example_channel_basics
// +build example_channel_basics

package main

import (
	"fmt"
	"time"
)

// channelBasics 演示Channel的基本操作
func channelBasics() {
	fmt.Println("=== 练习1: Channel基础操作 ===")

	// 1. 创建非缓冲Channel
	ch := make(chan string)

	// 启动一个goroutine发送数据
	go func() {
		fmt.Println("Goroutine: 准备发送数据...")
		time.Sleep(500 * time.Millisecond) // 模拟耗时操作
		ch <- "Hello from goroutine!"
		fmt.Println("Goroutine: 数据发送完成")
	}()

	// 主goroutine接收数据
	fmt.Println("Main: 等待接收数据...")
	msg := <-ch
	fmt.Printf("Main: 收到消息: %s\n", msg)

	// 2. 带缓冲的Channel
	fmt.Println("\n=== 练习2: 缓冲Channel ===")
	bufferedCh := make(chan int, 3) // 缓冲区大小为3

	// 可以连续发送多个数据而不阻塞
	bufferedCh <- 1
	bufferedCh <- 2
	bufferedCh <- 3
	fmt.Println("已发送3个数据到缓冲Channel")

	// 接收数据
	fmt.Println("开始接收数据:")
	for i := 0; i < 3; i++ {
		value := <-bufferedCh
		fmt.Printf("  收到: %d\n", value)
	}

	// 3. Channel关闭和遍历
	fmt.Println("\n=== 练习3: Channel关闭与遍历 ===")
	numbers := make(chan int)

	// 生产者goroutine
	go func() {
		for i := 1; i <= 5; i++ {
			numbers <- i
			time.Sleep(200 * time.Millisecond)
		}
		close(numbers) // 发送完成后关闭Channel
		fmt.Println("生产者: 所有数据已发送，Channel已关闭")
	}()

	// 消费者：使用range遍历Channel
	fmt.Println("消费者: 开始接收数据")
	for num := range numbers {
		fmt.Printf("  收到数字: %d\n", num)
	}
	fmt.Println("消费者: 所有数据已接收")

	// 4. 检查Channel是否关闭
	fmt.Println("\n=== 练习4: 检查Channel状态 ===")
	statusCh := make(chan string)

	go func() {
		statusCh <- "运行中"
		time.Sleep(300 * time.Millisecond)
		statusCh <- "完成"
		close(statusCh)
	}()

	// 接收并检查Channel状态
	for {
		status, ok := <-statusCh
		if !ok {
			fmt.Println("Channel已关闭")
			break
		}
		fmt.Printf("状态: %s\n", status)
	}
}

// channelDirection 演示Channel方向（只读/只写）
func channelDirection() {
	fmt.Println("\n=== 练习5: Channel方向 ===")

	// 创建一个双向Channel
	ch := make(chan int)

	// 只写Channel参数
	sendOnly := func(ch chan<- int, value int) {
		ch <- value
		fmt.Printf("已发送: %d\n", value)
	}

	// 只读Channel参数
	receiveOnly := func(ch <-chan int) {
		value := <-ch
		fmt.Printf("已接收: %d\n", value)
	}

	// 使用只写Channel发送数据
	go sendOnly(ch, 42)

	// 使用只读Channel接收数据
	receiveOnly(ch)
}

// channelBlocking 演示Channel的阻塞行为
func channelBlocking() {
	fmt.Println("\n=== 练习6: Channel阻塞行为 ===")

	// 非缓冲Channel的阻塞特性
	ch := make(chan string)

	go func() {
		fmt.Println("Goroutine: 等待3秒后发送数据")
		time.Sleep(3 * time.Second)
		ch <- "延迟消息"
	}()

	fmt.Println("Main: 等待接收数据（会阻塞直到有数据）")
	start := time.Now()
	msg := <-ch
	elapsed := time.Since(start)
	fmt.Printf("Main: 收到 '%s'，等待了 %v\n", msg, elapsed)
}

func main() {
	fmt.Println("第1周第2天: Channel通信基础")
	fmt.Println("============================\n")

	// 运行所有练习
	channelBasics()
	channelDirection()
	channelBlocking()

	fmt.Println("\n=== 总结 ===")
	fmt.Println("1. 非缓冲Channel: 发送和接收必须同时准备好，否则会阻塞")
	fmt.Println("2. 缓冲Channel: 可以暂存多个数据，发送时缓冲区满才会阻塞")
	fmt.Println("3. Channel关闭: 由发送方关闭，接收方可以使用range遍历")
	fmt.Println("4. Channel方向: 可以指定只读或只写，增加类型安全性")
	fmt.Println("5. 阻塞行为: Channel操作是同步的，理解阻塞对并发设计很重要")
}