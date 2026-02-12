# Go Learning Path - 4周专业开发训练计划

## 📅 训练计划

### 第1周：并发编程实战
- Day1: Goroutine基础与WaitGroup
- Day2: Channel通信
- Day3: Select与超时控制  
- Day4: Worker Pool实现

### 第2周：接口与测试驱动
- Day1: 接口设计
- Day2: 表格驱动测试
- Day3: Mock测试
- Day4: Benchmark性能分析

### 第3周：标准库深度
- Day1: net/http源码阅读
- Day2: 中间件链实现
- Day3: io.Reader/Writer
- Day4: context深入

### 第4周：完整项目实战
- 项目三选一：URL短链/任务队列/Markdown文档服务器

## 🏗️ 项目结构

```
.
├── cmd/                 # 可执行程序入口
├── internal/            # 私有应用程序代码
├── pkg/                 # 公共库代码
├── week1/               # 第1周学习代码
│   ├── day1/           # 第1天任务
│   ├── day2/           # 第2天任务
│   ├── day3/           # 第3天任务
│   └── day4/           # 第4天任务
├── week2/               # 第2周学习代码
├── week3/               # 第3周学习代码
├── week4/               # 第4周项目
├── scripts/             # 构建脚本
├── deploy/              # 部署配置
├── go.mod              # Go模块定义
└── README.md           # 项目说明
```

## 🚀 开始学习

1. 每周按顺序完成day1-day4的任务
2. 每天提交代码到GitHub
3. 周五进行代码审查
4. 遇到问题及时提问

## 📚 学习资源

- [Go官方文档](https://go.dev/doc/)
- [Go by Example](https://gobyexample.com/)
- 《Go语言圣经》
- 《Concurrency in Go》