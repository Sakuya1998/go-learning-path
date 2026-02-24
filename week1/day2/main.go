// 第1周第2天：Channel通信示例选择器
// 这个文件提供了一个统一的入口点来运行各个示例

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("第1周第2天：Channel通信示例选择器")
	fmt.Println("==================================")
	fmt.Println()
	fmt.Println("请选择要运行的示例：")
	fmt.Println("1. Channel基础操作 (channel_basics.go)")
	fmt.Println("2. 缓冲Channel实验 (buffered_channels.go)")
	fmt.Println("3. select多路复用 (select_demo.go)")
	fmt.Println("4. 生产者-消费者模式 (producer_consumer.go)")
	fmt.Println("5. Pipeline模式 (pipeline_pattern.go)")
	fmt.Println("6. 练习题目 (exercises.go)")
	fmt.Println("7. 运行所有测试")
	fmt.Println("0. 退出")
	fmt.Println()

	for {
		fmt.Print("请输入选择 (0-7): ")
		
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		switch input {
		case "1":
			fmt.Println("\n运行Channel基础示例...")
			fmt.Println("使用命令: go run -tags example_channel_basics channel_basics.go")
			fmt.Println()
			runExample("example_channel_basics", "channel_basics.go")
			
		case "2":
			fmt.Println("\n运行缓冲Channel实验...")
			fmt.Println("使用命令: go run -tags example_buffered_channels buffered_channels.go")
			fmt.Println()
			runExample("example_buffered_channels", "buffered_channels.go")
			
		case "3":
			fmt.Println("\n运行select多路复用示例...")
			fmt.Println("使用命令: go run -tags example_select_demo select_demo.go")
			fmt.Println()
			runExample("example_select_demo", "select_demo.go")
			
		case "4":
			fmt.Println("\n运行生产者-消费者模式示例...")
			fmt.Println("使用命令: go run -tags example_producer_consumer producer_consumer.go")
			fmt.Println()
			runExample("example_producer_consumer", "producer_consumer.go")
			
		case "5":
			fmt.Println("\n运行Pipeline模式示例...")
			fmt.Println("使用命令: go run -tags example_pipeline_pattern pipeline_pattern.go")
			fmt.Println()
			runExample("example_pipeline_pattern", "pipeline_pattern.go")
			
		case "6":
			fmt.Println("\n运行练习题目...")
			fmt.Println("使用命令: go run exercises.go")
			fmt.Println()
			runExercises()
			
		case "7":
			fmt.Println("\n运行所有测试...")
			fmt.Println("使用命令: go test -v")
			fmt.Println()
			runTests()
			
		case "0":
			fmt.Println("\n感谢使用，再见！")
			return
			
		default:
			fmt.Println("无效的选择，请重新输入")
		}
		
		fmt.Println()
		fmt.Println("----------------------------")
		fmt.Println()
	}
}

func runExample(tag string, filename string) {
	fmt.Printf("执行: go run -tags %s %s\n", tag, filename)
	fmt.Println("按Enter键继续...")
	bufio.NewReader(os.Stdin).ReadString('\n')
	
	// 在实际环境中，这里可以调用exec.Command来运行
	// 但为了简单起见，我们只显示命令
	fmt.Println("（在实际环境中，这里会自动运行示例）")
	fmt.Println("要手动运行，请使用上面的命令")
}

func runExercises() {
	fmt.Println("执行: go run exercises.go")
	fmt.Println("按Enter键继续...")
	bufio.NewReader(os.Stdin).ReadString('\n')
	
	fmt.Println("（在实际环境中，这里会自动运行练习）")
	fmt.Println("要手动运行，请使用: go run exercises.go")
}

func runTests() {
	fmt.Println("执行: go test -v")
	fmt.Println("按Enter键继续...")
	bufio.NewReader(os.Stdin).ReadString('\n')
	
	fmt.Println("（在实际环境中，这里会自动运行测试）")
	fmt.Println("要手动运行，请使用: go test -v")
}