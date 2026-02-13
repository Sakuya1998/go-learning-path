package main

import (
	"sync"
	"testing"
	"time"
)

// TestBasicGoroutine 测试基础goroutine功能
func TestBasicGoroutine(t *testing.T) {
	// 这个测试主要验证程序能正常运行而不panic
	// 由于goroutine是并发执行的，输出顺序不确定
	// 所以我们只测试程序能否正常执行完成

	start := time.Now()
	basicGoroutine()
	elapsed := time.Since(start)

	// 确保程序在合理时间内完成（应该很快）
	if elapsed > 2*time.Second {
		t.Errorf("basicGoroutine执行时间过长: %v", elapsed)
	}
}

// TestGoroutineWithParams 测试带参数的goroutine
func TestGoroutineWithParams(t *testing.T) {
	start := time.Now()
	goroutineWithParams()
	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Errorf("goroutineWithParams执行时间过长: %v", elapsed)
	}
}

// TestConcurrencyIssueDemo 测试并发问题演示
func TestConcurrencyIssueDemo(t *testing.T) {
	start := time.Now()
	concurrencyIssueDemo()
	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Errorf("concurrencyIssueDemo执行时间过长: %v", elapsed)
	}
}

// BenchmarkGoroutineCreation 基准测试：测量goroutine创建开销
func BenchmarkGoroutineCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 空操作，只测量创建开销
		}()
		wg.Wait()
	}
}

// 扩展测试（可选）：
// 1. 编写测试验证修复后的concurrencyIssueDemo能正确计数
// 2. 测试goroutine池的功能正确性
// 3. 测试带缓冲channel的通信
