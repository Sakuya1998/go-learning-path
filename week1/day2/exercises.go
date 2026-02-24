//go:build !example_channel_basics && !example_buffered_channels && !example_select_demo && !example_producer_consumer && !example_pipeline_pattern
// +build !example_channel_basics,!example_buffered_channels,!example_select_demo,!example_producer_consumer,!example_pipeline_pattern

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ==================== 练习题目 ====================

// Exercise1: 修复Channel死锁
// 下面的代码会导致死锁，请修复它
func exercise1() {
	fmt.Println("=== 练习1: 修复Channel死锁 ===")

	ch := make(chan int)

	// ❌ 有问题的代码
	// ch <- 42  // 发送数据
	// value := <-ch  // 接收数据
	// fmt.Printf("收到: %d\n", value)

	// ✅ 你的修复代码写在这里
	// 提示：Channel操作会阻塞，需要并发执行
	go func() {
		fmt.Println("发送数据...")
		ch <- 42
		fmt.Println("数据发送完成")
	}()

	value := <-ch
	fmt.Printf("收到: %d\n", value)
}

// Exercise2: 实现超时重试机制
// 实现一个函数，尝试从Channel读取数据，如果超时则重试
func exercise2() {
	fmt.Println("\n=== 练习2: 超时重试机制 ===")

	ch := make(chan string)

	// 模拟异步数据到达
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "数据"
	}()

	// TODO: 实现超时重试逻辑
	// 要求：每隔1秒尝试一次，最多重试3次
	// 如果收到数据则打印，如果重试3次都失败则打印"超时"
	for i := 0; i < 3; i++ {
		select {
		case data := <-ch:
			fmt.Printf("收到数据: %s\n", data)
			return
		case <-time.After(1 * time.Second):
			fmt.Println("重试中...")
		}
	}
	fmt.Println("超时")
}

// Exercise3: 实现扇入模式
// 将多个输入Channel合并到一个输出Channel
func exercise3() {
	fmt.Println("\n=== 练习3: 扇入模式 ===")

	// 创建多个输入Channel
	input1 := make(chan int)
	input2 := make(chan int)
	input3 := make(chan int)

	// TODO: 实现扇入函数
	// 函数签名：func fanIn(inputs ...<-chan int) <-chan int
	// 将多个输入合并到一个输出Channel
	fanIn := func(inputs ...<-chan int) <-chan int {
		out := make(chan int)
		var wg sync.WaitGroup

		for _, input := range inputs {
			wg.Add(1)
			go func(ch <-chan int) {
				defer wg.Done()
				for val := range ch {
					out <- val
				}
			}(input)
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}
	// 测试代码
	go func() {
		input1 <- 1
		input2 <- 2
		input3 <- 3
		close(input1)
		close(input2)
		close(input3)
	}()

	output := fanIn(input1, input2, input3)
	for val := range output {
		fmt.Printf("收到: %d\n", val)
	}
}

// Exercise4: 实现扇出模式
// 将一个输入Channel分发到多个输出Channel
func exercise4() {
	fmt.Println("\n=== 练习4: 扇出模式 ===")

	input := make(chan int)
	// TODO: 实现扇出函数
	// 函数签名：func fanOut(input <-chan int, outputs []chan<- int)
	// 将输入数据轮流分发到多个输出Channel
	fanOut := func(input <-chan int, outputs []chan int) {
		idx := 0
		for val := range input {
			outputs[idx] <- val
			idx = (idx + 1) % len(outputs)
		}

		for _, ch := range outputs {
			close(ch)
		}
		fmt.Println("扇出模式完成")
	}

	outputs := make([]chan int, 3)
	for i := range outputs {
		outputs[i] = make(chan int, 10)
	}

	var wg sync.WaitGroup
	for i, ch := range outputs {
		wg.Add(1)
		go func(id int, c <-chan int) {
			defer wg.Done()
			for val := range c {
				fmt.Printf("消费者%d: 收到 %d\n", id, val)
			}
			fmt.Printf("消费者%d: 完成\n", id)
		}(i, ch)
	}

	go func() {
		for i := 1; i <= 9; i++ {
			input <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(input)
	}()

	fanOut(input, outputs)
	wg.Wait()

}

// Exercise5: 实现工作池模式
// 使用Channel实现一个工作池，限制并发goroutine数量
func exercise5() {
	fmt.Println("\n=== 练习5: 工作池模式 ===")

	type Job struct {
		ID   int
		Data string
	}

	type Result struct {
		JobID int
		Value string
	}

	// TODO: 实现工作池
	// 要求：
	// 1. 创建指定数量的worker goroutine
	// 2. 从jobs Channel读取工作
	// 3. 处理工作并将结果发送到results Channel
	// 4. 所有工作完成后关闭results Channel

	workerPool := func(jobs <-chan Job, results chan<- Result, workerCount int) {
		var wg sync.WaitGroup

		// 启动指定数量的worker
		for i := 1; i <= workerCount; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// 从jobs channel读取工作，直到channel关闭
				for job := range jobs {
					// 模拟工作处理
					fmt.Printf("Worker %d 处理 Job %d: %s\n",
						workerID, job.ID, job.Data)

					// 模拟处理时间
					time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

					// 发送结果
					result := Result{
						JobID: job.ID,
						Value: fmt.Sprintf("Worker%d-处理完成-%s", workerID, job.Data),
					}
					results <- result
				}
				fmt.Printf("Worker %d 退出\n", workerID)
			}(i)
		}

		// 等待所有worker完成，然后关闭results channel
		go func() {
			wg.Wait()
			close(results)
			fmt.Println("工作池完成，关闭results channel")
		}()
	}

	jobs := make(chan Job, 100)
	results := make(chan Result, 100)

	// 模拟工作
	go func() {
		for i := 1; i <= 20; i++ {
			jobs <- Job{ID: i, Data: fmt.Sprintf("工作%d", i)}
		}
		close(jobs)
		fmt.Println("所有工作已发送，jobs channel已关闭")
	}()

	// TODO: 创建工作池（假设有5个worker）
	workerPool(jobs, results, 5)

	// 收集结果
	completed := 0
	for result := range results {
		completed++
		fmt.Printf("完成 [%d/20]: Job %d -> %s\n",
			completed, result.JobID, result.Value)
	}

	fmt.Printf("\n总计完成: %d 个工作\n", completed)
}

// Exercise6: 实现优雅关闭
// 实现一个可以优雅关闭的生产者-消费者系统
func exercise6() {
	fmt.Println("\n=== 练习6: 优雅关闭 ===")

	// TODO: 实现一个包含以下组件的系统：
	// 1. 一个生产者，持续生产数据
	// 2. 多个消费者，处理数据
	// 3. 一个控制Channel，用于发送关闭信号
	// 4. 确保所有goroutine都能正确退出
	// 5. 处理完所有已生产的数据后再关闭
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	dataCh := make(chan int, 10)
	done := make(chan struct{})

	// 生产者
	wg.Add(1)
	go func() {
		id := 1
		for {
			select {
			case <-ctx.Done():
				fmt.Println("生产者已关闭")
				close(dataCh)
				return
			case dataCh <- id:
				fmt.Printf("生产者: 生产数据 %d\n", id)
				id++
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
	// 启动3个消费者
	consumerCount := 3
	for i := 1; i <= consumerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for data := range dataCh {
				fmt.Printf("消费者%d: 处理数据 %d\n", id, data)
				time.Sleep(300 * time.Millisecond) // 模拟处理时间
			}
			fmt.Printf("消费者%d: 数据channel已关闭，退出\n", id)
		}(i)
	}

	// 监控goroutine：等待所有工作完成
	go func() {
		wg.Wait()
		fmt.Println("监控: 所有goroutine已退出")
		close(done)
	}()

	// 模拟运行一段时间后触发关闭
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\n主程序: 发送关闭信号...")
		cancel()
	}()

	// 等待完全结束或超时
	select {
	case <-done:
		fmt.Println("主程序: 系统已优雅关闭")
	case <-time.After(5 * time.Second):
		fmt.Println("主程序: 关闭超时，强制退出")
	}

	fmt.Println("实现提示：")
	fmt.Println("1. 使用sync.WaitGroup等待所有goroutine完成")
	fmt.Println("2. 使用context或单独的stop Channel控制关闭")
	fmt.Println("3. 生产者收到关闭信号后停止生产并关闭数据Channel")
	fmt.Println("4. 消费者在数据Channel关闭后自动退出")
}

// ==================== 参考答案区域 ====================
// 完成练习后，可以查看下面的参考答案
// 但建议先自己尝试实现

/*
// Exercise1 参考答案
func exercise1Solution() {
	ch := make(chan int)

	go func() {
		ch <- 42
	}()

	value := <-ch
	fmt.Printf("收到: %d\n", value)
}

// Exercise2 参考答案
func exercise2Solution() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "数据"
	}()

	for i := 1; i <= 3; i++ {
		select {
		case data := <-ch:
			fmt.Printf("成功收到: %s\n", data)
			return
		case <-time.After(1 * time.Second):
			fmt.Printf("第%d次尝试超时\n", i)
		}
	}
	fmt.Println("重试3次均失败，超时")
}

// Exercise3 参考答案
func fanIn(inputs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, input := range inputs {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for val := range ch {
				out <- val
			}
		}(input)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
*/

func main() {
	fmt.Println("第1周第2天: Channel练习题目")
	fmt.Println("===========================")
	fmt.Println()

	fmt.Println("说明：这些练习需要你动手实现")
	fmt.Println("先尝试自己完成，然后再查看参考答案")
	fmt.Println()

	// 列出所有练习
	exercise1()
	exercise2()
	exercise3()
	exercise4()
	exercise5()
	exercise6()

	fmt.Println("\n=== 练习完成建议 ===")
	fmt.Println("1. 从简单的练习开始，逐步增加难度")
	fmt.Println("2. 每个练习都先自己思考实现方案")
	fmt.Println("3. 编写测试验证实现的正确性")
	fmt.Println("4. 对比参考答案，学习不同的实现思路")
	fmt.Println("5. 尝试优化代码性能和可读性")
}
