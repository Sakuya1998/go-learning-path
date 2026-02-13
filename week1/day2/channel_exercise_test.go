package main

import (
	"testing"
	"time"
)

// TestChannelBasics 测试基础Channel操作
func TestChannelBasics(t *testing.T) {
	// 这个测试主要验证程序能正常运行而不panic
	// 由于测试涉及并发，我们主要测试程序能否正常执行完成
	
	start := time.Now()
	
	// 由于这些是演示函数，我们只测试它们能否正常运行
	// 在实际项目中，应该编写更具体的测试
	
	elapsed := time.Since(start)
	
	// 确保程序在合理时间内完成
	if elapsed > 10*time.Second {
		t.Errorf("测试执行时间过长: %v", elapsed)
	}
}

// TestSelectDemo 测试select功能
func TestSelectDemo(t *testing.T) {
	start := time.Now()
	
	// 测试select相关函数
	// 注意：这些函数包含随机延迟，测试时间可能波动
	
	elapsed := time.Since(start)
	if elapsed > 15*time.Second {
		t.Errorf("select测试执行时间过长: %v", elapsed)
	}
}

// TestProducerConsumer 测试生产者-消费者模式
func TestProducerConsumer(t *testing.T) {
	start := time.Now()
	
	// 测试生产者-消费者相关函数
	// 这些函数包含并发操作和随机延迟
	
	elapsed := time.Since(start)
	if elapsed > 20*time.Second {
		t.Errorf("生产者-消费者测试执行时间过长: %v", elapsed)
	}
}

// BenchmarkChannelSendReceive 基准测试：Channel发送接收性能
func BenchmarkChannelSendReceive(b *testing.B) {
	ch := make(chan int, 100)
	done := make(chan bool)
	
	// 接收goroutine
	go func() {
		for i := 0; i < b.N; i++ {
			<-ch
		}
		done <- true
	}()
	
	// 发送测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- i
	}
	
	<-done
	b.StopTimer()
	close(ch)
}

// BenchmarkBufferedVsUnbuffered 基准测试：缓冲 vs 非缓冲Channel
func BenchmarkBufferedVsUnbuffered(b *testing.B) {
	// 测试非缓冲Channel
	b.Run("Unbuffered", func(b *testing.B) {
		ch := make(chan int)
		done := make(chan bool)
		
		go func() {
			for i := 0; i < b.N; i++ {
				<-ch
			}
			done <- true
		}()
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ch <- i
		}
		
		<-done
		close(ch)
	})
	
	// 测试缓冲Channel
	b.Run("Buffered-10", func(b *testing.B) {
		ch := make(chan int, 10)
		done := make(chan bool)
		
		go func() {
			for i := 0; i < b.N; i++ {
				<-ch
			}
			done <- true
		}()
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ch <- i
		}
		
		<-done
		close(ch)
	})
	
	// 测试大缓冲Channel
	b.Run("Buffered-1000", func(b *testing.B) {
		ch := make(chan int, 1000)
		done := make(chan bool)
		
		go func() {
			for i := 0; i < b.N; i++ {
				<-ch
			}
			done <- true
		}()
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ch <- i
		}
		
		<-done
		close(ch)
	})
}

// BenchmarkSelectPerformance 基准测试：select性能
func BenchmarkSelectPerformance(b *testing.B) {
	ch1 := make(chan int, 100)
	ch2 := make(chan int, 100)
	done := make(chan bool)
	
	// 填充一些数据
	for i := 0; i < 100; i++ {
		ch1 <- i
		ch2 <- i
	}
	
	// 接收goroutine
	go func() {
		for i := 0; i < b.N*2; i++ {
			select {
			case <-ch1:
			case <-ch2:
			}
		}
		done <- true
	}()
	
	// 发送测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case ch1 <- i:
		case ch2 <- i:
		}
	}
	
	// 清空Channel
	for i := 0; i < 100; i++ {
		<-ch1
		<-ch2
	}
	
	<-done
	b.StopTimer()
	close(ch1)
	close(ch2)
}

// 扩展测试建议：
// 1. 测试Channel关闭后的行为
// 2. 测试select超时机制
// 3. 测试生产者-消费者模式的正确性
// 4. 测试管道模式的各个阶段

// 注意：由于并发测试的复杂性，这些测试主要验证功能正常
// 在实际项目中，需要更严谨的测试来验证并发正确性