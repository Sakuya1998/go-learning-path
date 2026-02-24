//go:build example_buffered_channels
// +build example_buffered_channels

package main

import (
	"fmt"
	"time"
)

// bufferedChannels 专门演示缓冲Channel的特性和实验
func main() {
	fmt.Println("=== 缓冲Channel实验 ===")
	
	// 实验1: 缓冲区大小的影响
	fmt.Println("\n1. 缓冲区大小的影响实验:")
	bufferSizeExperiment()
	
	// 实验2: 阻塞与非阻塞行为
	fmt.Println("\n2. 阻塞与非阻塞行为实验:")
	blockingBehaviorExperiment()
	
	// 实验3: 性能对比实验
	fmt.Println("\n3. 性能对比实验:")
	performanceExperiment()
	
	// 实验4: 实际应用场景
	fmt.Println("\n4. 实际应用场景演示:")
	practicalUseCase()
	
	// 实验5: 缓冲区溢出实验
	fmt.Println("\n5. 缓冲区溢出实验:")
	bufferOverflowExperiment()
}

// bufferSizeExperiment 演示不同缓冲区大小的影响
func bufferSizeExperiment() {
	fmt.Println("实验目的：观察不同缓冲区大小对程序行为的影响")
	
	// 测试不同的缓冲区大小
	bufferSizes := []int{0, 1, 3, 5, 10}
	
	for _, size := range bufferSizes {
		fmt.Printf("\n--- 缓冲区大小: %d ---\n", size)
		
		ch := make(chan int, size)
		done := make(chan bool)
		
		// 生产者
		go func() {
			for i := 1; i <= 5; i++ {
				ch <- i
				fmt.Printf("  生产者: 发送 %d\n", i)
				time.Sleep(100 * time.Millisecond)
			}
			close(ch)
		}()
		
		// 消费者
		go func() {
			for value := range ch {
				fmt.Printf("  消费者: 收到 %d\n", value)
				time.Sleep(200 * time.Millisecond) // 消费者处理较慢
			}
			done <- true
		}()
		
		<-done
	}
	
	fmt.Println("\n观察结论：")
	fmt.Println("- 缓冲区大小为0（非缓冲）：发送和接收必须同时准备好，否则阻塞")
	fmt.Println("- 缓冲区大小>0：生产者可以提前发送多个数据，减少等待时间")
	fmt.Println("- 缓冲区越大，生产者和消费者的耦合度越低")
}

// blockingBehaviorExperiment 演示阻塞与非阻塞行为
func blockingBehaviorExperiment() {
	fmt.Println("实验目的：理解缓冲Channel的阻塞行为")
	
	// 创建一个小缓冲区
	ch := make(chan int, 2)
	
	fmt.Println("场景1: 缓冲区未满时的非阻塞发送")
	ch <- 1
	ch <- 2
	fmt.Println("  成功发送2个数据到缓冲区（容量2）")
	
	fmt.Println("\n场景2: 缓冲区已满时的阻塞发送")
	go func() {
		fmt.Println("  尝试发送第3个数据（将阻塞直到有空间）")
		ch <- 3
		fmt.Println("  第3个数据发送成功")
	}()
	
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  从缓冲区取出1个数据，释放空间")
	<-ch
	
	time.Sleep(500 * time.Millisecond)
	
	fmt.Println("\n场景3: 缓冲区为空时的阻塞接收")
	go func() {
		fmt.Println("  尝试从空缓冲区接收（将阻塞直到有数据）")
		value := <-ch
		fmt.Printf("  收到数据: %d\n", value)
	}()
	
	time.Sleep(500 * time.Millisecond)
	fmt.Println("  向缓冲区发送新数据")
	ch <- 4
	
	time.Sleep(500 * time.Millisecond)
	
	// 清理
	<-ch
}

// performanceExperiment 性能对比实验
func performanceExperiment() {
	fmt.Println("实验目的：对比不同缓冲区大小的性能差异")
	
	const numOperations = 100000
	
	// 测试不同缓冲区大小
	sizes := []struct {
		size int
		name string
	}{
		{0, "非缓冲"},
		{1, "缓冲1"},
		{10, "缓冲10"},
		{100, "缓冲100"},
		{1000, "缓冲1000"},
	}
	
	for _, test := range sizes {
		ch := make(chan int, test.size)
		done := make(chan bool)
		
		start := time.Now()
		
		// 生产者
		go func() {
			for i := 0; i < numOperations; i++ {
				ch <- i
			}
			close(ch)
		}()
		
		// 消费者
		go func() {
			for range ch {
				// 只是消费，不做实际工作
			}
			done <- true
		}()
		
		<-done
		elapsed := time.Since(start)
		
		fmt.Printf("%s Channel: 处理 %d 次操作耗时 %v\n", 
			test.name, numOperations, elapsed)
	}
	
	fmt.Println("\n性能分析：")
	fmt.Println("- 非缓冲Channel：每次发送都需要等待接收，上下文切换开销大")
	fmt.Println("- 小缓冲区：减少部分等待，但仍有较多上下文切换")
	fmt.Println("- 适当缓冲区：平衡内存使用和性能，减少等待时间")
	fmt.Println("- 过大缓冲区：可能隐藏性能问题，增加内存占用")
}

// practicalUseCase 实际应用场景演示
func practicalUseCase() {
	fmt.Println("实际应用：生产者-消费者模式中的缓冲Channel")
	
	// 模拟一个日志处理系统
	logChannel := make(chan string, 100) // 缓冲100条日志
	done := make(chan bool)
	
	// 日志生产者（多个goroutine）
	const numProducers = 3
	for i := 1; i <= numProducers; i++ {
		go func(producerID int) {
			for j := 1; j <= 5; j++ {
				log := fmt.Sprintf("生产者%d-日志%d", producerID, j)
				logChannel <- log
				fmt.Printf("  生产者%d: 产生日志 '%s'\n", producerID, log)
				time.Sleep(time.Duration(producerID*50) * time.Millisecond)
			}
		}(i)
	}
	
	// 日志消费者
	go func() {
		logCount := 0
		for log := range logChannel {
			logCount++
			fmt.Printf("  消费者: 处理日志 '%s' (已处理: %d)\n", log, logCount)
			time.Sleep(150 * time.Millisecond) // 模拟处理时间
			
			if logCount >= numProducers*5 {
				close(logChannel)
			}
		}
		done <- true
	}()
	
	<-done
	fmt.Println("\n应用场景总结：")
	fmt.Println("1. 缓冲Channel作为缓冲区，平衡生产者和消费者的速度差异")
	fmt.Println("2. 防止快速生产者拖慢慢速消费者")
	fmt.Println("3. 提高系统吞吐量，减少等待时间")
	fmt.Println("4. 平滑流量峰值，提高系统稳定性")
}

// bufferOverflowExperiment 缓冲区溢出实验
func bufferOverflowExperiment() {
	fmt.Println("实验目的：理解缓冲区溢出的影响和预防")
	
	fmt.Println("场景：生产者速度远大于消费者速度")
	
	// 创建一个小缓冲区
	ch := make(chan int, 3)
	stop := make(chan bool)
	
	// 快速生产者
	go func() {
		i := 1
		for {
			select {
			case <-stop:
				return
			default:
				select {
				case ch <- i:
					fmt.Printf("  生产者: 发送 %d (缓冲区占用: %d/3)\n", 
						i, len(ch))
					i++
				default:
					// 缓冲区已满，非阻塞发送失败
					fmt.Println("  ⚠️ 缓冲区已满！生产者等待...")
					time.Sleep(500 * time.Millisecond)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	// 慢速消费者
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(500 * time.Millisecond) // 消费很慢
			value := <-ch
			fmt.Printf("  消费者: 处理 %d (缓冲区占用: %d/3)\n", 
				value, len(ch))
		}
		stop <- true
	}()
	
	time.Sleep(6 * time.Second)
	
	fmt.Println("\n缓冲区管理策略：")
	fmt.Println("1. 合理设置缓冲区大小：根据生产消费速度差决定")
	fmt.Println("2. 监控缓冲区使用：使用 len(ch) 监控当前占用")
	fmt.Println("3. 非阻塞操作：使用 select+default 避免死锁")
	fmt.Println("4. 动态调整：根据负载动态调整缓冲区大小")
	fmt.Println("5. 背压机制：当缓冲区快满时通知生产者减速")
}

// ==================== 最佳实践总结 ====================

/*
缓冲Channel最佳实践：

1. **选择合适的缓冲区大小**
   - 太小：频繁阻塞，性能差
   - 太大：内存浪费，可能隐藏问题
   - 经验值：通常10-100之间，根据具体场景调整

2. **监控缓冲区使用**
   - 使用 len(ch) 获取当前元素数量
   - 使用 cap(ch) 获取缓冲区容量
   - 监控比例：len(ch)/cap(ch) > 0.8 时考虑扩容或背压

3. **处理缓冲区满的情况**
   - 使用select的default分支实现非阻塞操作
   - 实现超时机制
   - 实现背压控制

4. **避免常见陷阱**
   - 不要依赖缓冲区解决所有同步问题
   - 注意goroutine泄漏
   - 正确关闭Channel

5. **性能优化**
   - 批量处理：减少Channel操作次数
   - 对象池：复用对象减少GC压力
   - 异步处理：非关键路径使用缓冲Channel解耦
*/

// ==================== 实验练习 ====================

func bufferExperiments() {
	fmt.Println("\n=== 缓冲Channel实验练习 ===")
	
	fmt.Println("练习1: 实现一个动态调整缓冲区大小的Channel")
	fmt.Println("  要求：根据生产消费速度自动调整缓冲区大小")
	
	fmt.Println("\n练习2: 实现背压机制")
	fmt.Println("  要求：当缓冲区使用率超过阈值时通知生产者减速")
	
	fmt.Println("\n练习3: 实现多级缓冲")
	fmt.Println("  要求：创建多个不同优先级的缓冲Channel")
	
	fmt.Println("\n练习4: 实现缓冲区监控")
	fmt.Println("  要求：实时监控并报告各个Channel的缓冲区状态")
	
	fmt.Println("\n练习5: 性能测试框架")
	fmt.Println("  要求：自动化测试不同场景下的缓冲Channel性能")
	
	fmt.Println("\n提示：从修改上面的实验代码开始，逐步实现更复杂的功能。")
}