# 第1周第1天：Goroutine基础与WaitGroup

## 学习目标
1. 理解goroutine的基本概念和创建方式
2. 掌握sync.WaitGroup的使用方法
3. 认识并发编程中的数据竞争问题
4. 学习goroutine的参数传递注意事项

## 任务说明

### 主要任务
完成 `goroutine_exercise.go` 文件中的三个练习：

1. **basicGoroutine()** - 基础goroutine练习
   - 创建10个goroutine
   - 每个goroutine打印自己的ID
   - 使用WaitGroup等待所有goroutine完成

2. **goroutineWithParams()** - 带参数的goroutine
   - 处理字符串数组中的消息
   - 注意循环变量捕获的问题
   - 正确传递参数到goroutine

3. **concurrencyIssueDemo()** - 并发安全问题演示
   - 演示数据竞争问题
   - 理解为什么counter++不是原子操作
   - 思考如何修复

### 运行代码
```bash
cd week1/day1
go run goroutine_exercise.go
```

### 运行测试
```bash
go test -v
```

## 关键知识点

### 1. goroutine创建
```go
go func() {
    // 并发执行的代码
}()
```

### 2. WaitGroup使用模式
```go
var wg sync.WaitGroup
wg.Add(1)  // 启动前增加计数
go func() {
    defer wg.Done()  // 完成后减少计数
    // 工作代码
}()
wg.Wait()  // 等待所有完成
```

### 3. 参数传递陷阱
```go
// 错误：所有goroutine共享同一个i
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i)  // 可能都打印10
    }()
}

// 正确：传递参数副本
for i := 0; i < 10; i++ {
    go func(id int) {
        fmt.Println(id)  // 正确打印0-9
    }(i)
}
```

### 4. 数据竞争（Data Race）
多个goroutine同时读写同一个变量，且至少有一个是写操作。

## 扩展练习

### 必做
1. 修复 `concurrencyIssueDemo()` 中的数据竞争问题
   - 使用 `sync.Mutex` 实现线程安全计数器
   - 使用 `sync/atomic` 包实现原子计数器

### 选做
1. 实现一个简单的goroutine池
   - 限制同时运行的goroutine数量为5个
   - 使用channel控制并发度
2. 测量goroutine创建和销毁的开销
3. 比较Mutex和atomic的性能差异

## 学习检查
完成今天的学习后，你应该能够：
- [ ] 独立创建和管理多个goroutine
- [ ] 正确使用WaitGroup进行同步
- [ ] 理解并避免常见的goroutine参数传递错误
- [ ] 识别基本的数据竞争问题
- [ ] 知道如何使用Mutex保护共享数据

## 常见问题

### Q: goroutine和线程有什么区别？
A: goroutine是用户态的轻量级线程，由Go运行时调度。创建开销小（几KB），可以创建成千上万个。OS线程是内核态的，创建开销大（MB级）。

### Q: defer wg.Done()放在哪里？
A: 通常放在goroutine函数开头，使用defer确保无论函数如何退出都会执行。

### Q: 如何控制goroutine的数量？
A: 可以使用带缓冲的channel作为信号量，或者使用worker pool模式。

## 下一步
完成今天的练习后，提交代码到GitHub，并准备第2天的Channel学习。