//go:build example_producer_consumer
// +build example_producer_consumer

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// basicProducerConsumer 基础生产者-消费者模式
func basicProducerConsumer() {
	fmt.Println("=== 练习1: 基础生产者-消费者 ===")

	const bufferSize = 5
	const numItems = 10

	// 共享缓冲区
	buffer := make(chan int, bufferSize)
	done := make(chan bool)

	// 生产者
	go func() {
		for i := 1; i <= numItems; i++ {
			// 模拟生产时间
			time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
			
			buffer <- i
			fmt.Printf("生产者: 生产了物品 %d (缓冲区大小: %d/%d)\n", 
				i, len(buffer), bufferSize)
		}
		close(buffer) // 生产完成后关闭Channel
		fmt.Println("生产者: 所有物品已生产完成")
	}()

	// 消费者
	go func() {
		for item := range buffer {
			// 模拟消费时间
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			
			fmt.Printf("消费者: 消费了物品 %d (缓冲区剩余: %d/%d)\n", 
				item, len(buffer), bufferSize)
		}
		done <- true
		fmt.Println("消费者: 所有物品已消费完成")
	}()

	// 等待消费完成
	<-done
}

// multipleProducersConsumers 多个生产者和消费者
func multipleProducersConsumers() {
	fmt.Println("\n=== 练习2: 多个生产者-消费者 ===")

	const numProducers = 3
	const numConsumers = 2
	const totalItems = 15

	buffer := make(chan int, 5)
	var wg sync.WaitGroup

	// 启动多个生产者
	for p := 1; p <= numProducers; p++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()
			
			itemsPerProducer := totalItems / numProducers
			for i := 1; i <= itemsPerProducer; i++ {
				item := producerID*100 + i // 生成唯一ID
				time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
				
				buffer <- item
				fmt.Printf("生产者 %d: 生产了 %d\n", producerID, item)
			}
			fmt.Printf("生产者 %d: 完成\n", producerID)
		}(p)
	}

	// 等待所有生产者完成，然后关闭buffer
	go func() {
		wg.Wait()
		close(buffer)
		fmt.Println("所有生产者已完成，关闭缓冲区")
	}()

	// 启动多个消费者
	var consumerWg sync.WaitGroup
	for c := 1; c <= numConsumers; c++ {
		consumerWg.Add(1)
		go func(consumerID int) {
			defer consumerWg.Done()
			
			for item := range buffer {
				time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
				fmt.Printf("消费者 %d: 消费了 %d\n", consumerID, item)
			}
			fmt.Printf("消费者 %d: 完成\n", consumerID)
		}(c)
	}

	// 等待所有消费者完成
	consumerWg.Wait()
	fmt.Println("所有消费者已完成")
}

// producerConsumerWithControl 带控制信号的生产者-消费者
func producerConsumerWithControl() {
	fmt.Println("\n=== 练习3: 带控制信号的生产者-消费者 ===")

	buffer := make(chan string, 3)
	stop := make(chan bool)  // 停止信号
	stopped := make(chan bool) // 确认停止

	// 生产者：持续生产直到收到停止信号
	go func() {
		counter := 1
		for {
			select {
			case <-stop:
				fmt.Println("生产者: 收到停止信号")
				close(buffer)
				stopped <- true
				return
			default:
				item := fmt.Sprintf("产品-%d", counter)
				time.Sleep(200 * time.Millisecond)
				
				select {
				case buffer <- item:
					fmt.Printf("生产者: 生产了 %s\n", item)
					counter++
				case <-stop:
					continue // 立即检查停止信号
				}
			}
		}
	}()

	// 消费者
	go func() {
		for item := range buffer {
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("消费者: 消费了 %s\n", item)
		}
		fmt.Println("消费者: 缓冲区已关闭，停止消费")
	}()

	// 主程序控制：5秒后停止生产
	time.Sleep(5 * time.Second)
	fmt.Println("\n主程序: 发送停止信号...")
	stop <- true
	
	// 等待生产者确认停止
	<-stopped
	time.Sleep(1 * time.Second) // 给消费者时间处理剩余项目
	fmt.Println("主程序: 系统已停止")
}

// pipelinePattern 管道模式：多个处理阶段
func pipelinePattern() {
	fmt.Println("\n=== 练习4: 管道模式（Pipeline） ===")

	// 第一阶段：数据生成
	generate := func(done <-chan bool) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 1; i <= 10; i++ {
				select {
				case out <- i:
					fmt.Printf("生成: %d\n", i)
					time.Sleep(100 * time.Millisecond)
				case <-done:
					return
				}
			}
		}()
		return out
	}

	// 第二阶段：数据处理（平方）
	square := func(done <-chan bool, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for num := range in {
				select {
				case out <- num * num:
					fmt.Printf("平方: %d -> %d\n", num, num*num)
					time.Sleep(150 * time.Millisecond)
				case <-done:
					return
				}
			}
		}()
		return out
	}

	// 第三阶段：数据输出
	print := func(done <-chan bool, in <-chan int) {
		for result := range in {
			select {
			case <-done:
				return
			default:
				fmt.Printf("输出: %d\n", result)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	// 构建管道
	done := make(chan bool)
	defer close(done)

	numbers := generate(done)
	squares := square(done, numbers)
	
	fmt.Println("开始管道处理...")
	print(done, squares)
	fmt.Println("管道处理完成")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Println("第1周第2天: 生产者-消费者模式")
	fmt.Println("============================\n")

	basicProducerConsumer()
	multipleProducersConsumers()
	producerConsumerWithControl()
	pipelinePattern()

	fmt.Println("\n=== 模式总结 ===")
	fmt.Println("1. 基础模式: 单个生产者和消费者，使用缓冲Channel协调")
	fmt.Println("2. 多对多模式: 多个生产者和消费者，需要同步机制")
	fmt.Println("3. 控制模式: 使用信号Channel实现优雅停止")
	fmt.Println("4. 管道模式: 多个处理阶段串联，每个阶段独立并发")
	fmt.Println("5. 扇入扇出: 多个输入合并或一个输入分发到多个输出")
	
	fmt.Println("\n=== 最佳实践 ===")
	fmt.Println("• 由生产者关闭Channel")
	fmt.Println("• 使用sync.WaitGroup等待goroutine完成")
	fmt.Println("• 合理设置缓冲区大小避免死锁")
	fmt.Println("• 使用select实现超时和取消")
	fmt.Println("• 考虑背压(backpressure)机制")
}