//go:build example_pipeline_pattern
// +build example_pipeline_pattern

package main

import (
	"fmt"
	"sync"
	"time"
)

// Pipeline模式：构建数据处理流水线
// Pipeline是一种将多个处理阶段连接起来的模式，每个阶段都是一个goroutine
// 数据通过Channel在各个阶段之间流动，形成处理流水线

func main() {
	fmt.Println("=== Pipeline模式示例 ===")
	
	// 示例1: 基础Pipeline
	fmt.Println("\n1. 基础Pipeline示例:")
	basicPipeline()
	
	// 示例2: 多阶段Pipeline
	fmt.Println("\n2. 多阶段Pipeline示例:")
	multiStagePipeline()
	
	// 示例3: 并行Pipeline
	fmt.Println("\n3. 并行Pipeline示例:")
	parallelPipeline()
	
	// 示例4: 带错误处理的Pipeline
	fmt.Println("\n4. 带错误处理的Pipeline:")
	pipelineWithErrorHandling()
	
	// 示例5: 动态Pipeline
	fmt.Println("\n5. 动态Pipeline示例:")
	dynamicPipeline()
}

// ==================== 基础Pipeline ====================

// basicPipeline 演示最简单的Pipeline模式
func basicPipeline() {
	// 创建Pipeline的各个阶段
	numbers := generateNumbers(1, 10)
	squared := square(numbers)
	cubed := cube(squared)
	
	// 消费最终结果
	for result := range cubed {
		fmt.Printf("结果: %d\n", result)
	}
}

// generateNumbers 生成数字序列
func generateNumbers(start, end int) <-chan int {
	out := make(chan int)
	
	go func() {
		for i := start; i <= end; i++ {
			out <- i
		}
		close(out)
	}()
	
	return out
}

// square 计算平方
func square(in <-chan int) <-chan int {
	out := make(chan int)
	
	go func() {
		for num := range in {
			out <- num * num
		}
		close(out)
	}()
	
	return out
}

// cube 计算立方
func cube(in <-chan int) <-chan int {
	out := make(chan int)
	
	go func() {
		for num := range in {
			out <- num * num * num
		}
		close(out)
	}()
	
	return out
}

// ==================== 多阶段Pipeline ====================

// multiStagePipeline 演示多阶段Pipeline
func multiStagePipeline() {
	// 创建更复杂的Pipeline
	numbers := generateNumbers(1, 5)
	
	// 数据处理流水线：生成 → 平方 → 加10 → 乘以2
	stage1 := processStage(numbers, "平方", func(x int) int { return x * x })
	stage2 := processStage(stage1, "加10", func(x int) int { return x + 10 })
	stage3 := processStage(stage2, "乘以2", func(x int) int { return x * 2 })
	
	for result := range stage3 {
		fmt.Printf("处理结果: %d\n", result)
	}
}

// processStage 通用的处理阶段函数
func processStage(in <-chan int, stageName string, processor func(int) int) <-chan int {
	out := make(chan int)
	
	go func() {
		for num := range in {
			result := processor(num)
			fmt.Printf("[%s阶段] 输入: %d, 输出: %d\n", stageName, num, result)
			out <- result
		}
		close(out)
	}()
	
	return out
}

// ==================== 并行Pipeline ====================

// parallelPipeline 演示并行处理的Pipeline
func parallelPipeline() {
	// 生成数据
	numbers := generateNumbers(1, 20)
	
	// 创建多个并行的处理阶段
	var wg sync.WaitGroup
	results := make(chan int, 20)
	
	// 启动多个并行的处理器
	const numWorkers = 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go pipelineWorker(i, numbers, results, &wg)
	}
	
	// 等待所有worker完成
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// 收集结果
	var allResults []int
	for result := range results {
		allResults = append(allResults, result)
	}
	
	fmt.Printf("并行处理完成，共处理 %d 个数据\n", len(allResults))
	fmt.Printf("结果: %v\n", allResults)
}

// pipelineWorker 并行处理worker（专用于Pipeline模式）
func pipelineWorker(id int, in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for num := range in {
		// 模拟耗时处理
		time.Sleep(50 * time.Millisecond)
		result := num * 100
		fmt.Printf("[Worker %d] 处理 %d → %d\n", id, num, result)
		out <- result
	}
}

// ==================== 带错误处理的Pipeline ====================

// pipelineWithErrorHandling 演示带错误处理的Pipeline
func pipelineWithErrorHandling() {
	numbers := generateNumbersWithError(1, 10)
	
	// 带错误检查的处理阶段
	checkedResults := checkAndProcess(numbers)
	
	for {
		select {
		case result, ok := <-checkedResults:
			if !ok {
				fmt.Println("Pipeline处理完成")
				return
			}
			fmt.Printf("成功处理: %v\n", result)
		}
	}
}

// generateNumbersWithError 生成带随机错误的数字序列
func generateNumbersWithError(start, end int) <-chan interface{} {
	out := make(chan interface{})
	
	go func() {
		for i := start; i <= end; i++ {
			// 模拟随机错误（每3个数字有一个错误）
			if i%3 == 0 {
				out <- fmt.Errorf("处理数字 %d 时发生错误", i)
			} else {
				out <- i
			}
		}
		close(out)
	}()
	
	return out
}

// checkAndProcess 检查并处理数据，过滤错误
func checkAndProcess(in <-chan interface{}) <-chan int {
	out := make(chan int)
	
	go func() {
		for item := range in {
			switch v := item.(type) {
			case int:
				// 正常处理
				result := v * 2
				out <- result
			case error:
				// 处理错误
				fmt.Printf("⚠️ 错误: %v\n", v)
			default:
				fmt.Printf("⚠️ 未知类型: %v\n", v)
			}
		}
		close(out)
	}()
	
	return out
}

// ==================== 动态Pipeline ====================

// dynamicPipeline 演示动态构建的Pipeline
func dynamicPipeline() {
	fmt.Println("构建动态Pipeline...")
	
	// 初始数据源
	source := make(chan int)
	
	// 动态添加处理阶段
	var currentStage <-chan int = source
	
	// 定义可用的处理函数
	processors := []struct {
		name string
		fn   func(int) int
	}{
		{"加倍", func(x int) int { return x * 2 }},
		{"平方", func(x int) int { return x * x }},
		{"加5", func(x int) int { return x + 5 }},
		{"乘以3", func(x int) int { return x * 3 }},
	}
	
	// 动态选择3个处理阶段
	fmt.Println("选择处理阶段: 加倍 → 平方 → 加5")
	for i := 0; i < 3; i++ {
		processor := processors[i]
		currentStage = dynamicProcessStage(currentStage, processor.name, processor.fn)
	}
	
	// 启动数据源
	go func() {
		for i := 1; i <= 5; i++ {
			source <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(source)
	}()
	
	// 收集结果
	for result := range currentStage {
		fmt.Printf("动态Pipeline结果: %d\n", result)
	}
}

// dynamicProcessStage 动态处理阶段
func dynamicProcessStage(in <-chan int, stageName string, processor func(int) int) <-chan int {
	out := make(chan int)
	
	go func() {
		for num := range in {
			result := processor(num)
			fmt.Printf("[动态阶段: %s] %d → %d\n", stageName, num, result)
			out <- result
		}
		close(out)
	}()
	
	return out
}

// ==================== Pipeline模式总结 ====================

/*
Pipeline模式的核心优势：

1. **模块化**: 每个处理阶段独立，易于测试和维护
2. **可组合性**: 可以灵活组合不同的处理阶段
3. **并发性**: 各个阶段可以并行执行，提高吞吐量
4. **可扩展性**: 容易添加新的处理阶段或调整现有阶段
5. **错误隔离**: 错误可以在各个阶段独立处理

常见应用场景：
1. 数据ETL（提取、转换、加载）
2. 图像/视频处理流水线
3. 日志处理和分析
4. 实时数据处理
5. 工作流引擎

最佳实践：
1. 明确每个阶段的输入输出格式
2. 使用缓冲Channel提高性能
3. 实现优雅的关闭机制
4. 添加超时和错误处理
5. 监控各个阶段的性能指标
*/

// ==================== 练习题目 ====================

/*
扩展练习：

1. **实现文件处理Pipeline**
   创建一个Pipeline，依次执行：
   - 读取文件
   - 过滤无效行
   - 转换数据格式
   - 写入新文件

2. **实现可配置的Pipeline**
   创建一个可以从配置文件动态加载处理阶段的Pipeline系统

3. **实现带限流的Pipeline**
   为Pipeline添加限流功能，控制每个阶段的最大处理速率

4. **实现监控Pipeline**
   添加监控功能，实时显示各个阶段的处理状态和性能指标

5. **实现容错Pipeline**
   添加重试机制和故障转移功能，提高Pipeline的可靠性
*/

func pipelineExercises() {
	fmt.Println("\n=== Pipeline扩展练习 ===")
	fmt.Println("1. 实现文件处理Pipeline")
	fmt.Println("2. 实现可配置的Pipeline")
	fmt.Println("3. 实现带限流的Pipeline")
	fmt.Println("4. 实现监控Pipeline")
	fmt.Println("5. 实现容错Pipeline")
	
	fmt.Println("\n提示：可以从修改上面的示例代码开始，逐步实现更复杂的功能。")
}