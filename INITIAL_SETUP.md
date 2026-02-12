# 初始化设置指南

## 第一步：克隆仓库到本地
```bash
git clone https://github.com/Sakuya1998/go-learning-path.git
cd go-learning-path
```

## 第二步：初始化Git（如果仓库是空的）
```bash
# 如果仓库是新建的，需要设置上游
git add .
git commit -m "初始提交: 项目结构和第1周第1天学习代码"
git push -u origin main
```

## 第三步：安装依赖
```bash
go mod download
```

## 第四步：运行第1天练习
```bash
# 方法1: 使用Makefile
make week1-day1

# 方法2: 直接运行
cd week1/day1
go run goroutine_exercise.go
```

## 第五步：运行测试
```bash
make week1-day1-test
```

## 开发工具推荐安装
```bash
# 代码格式化
go install golang.org/x/tools/cmd/goimports@latest

# 代码检查
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 测试覆盖率
go install golang.org/x/tools/cmd/cover@latest
```

## 每日学习流程
1. 阅读当天的README.md了解学习目标
2. 完成代码练习
3. 运行测试确保代码正确
4. 提交代码到GitHub
5. 如有问题，查看常见问题部分或提问

## 代码提交规范
```bash
# 每天学习完成后提交
git add week1/day1/
git commit -m "第1周第1天: 完成goroutine基础练习"
git push
```

## 获取帮助
如果遇到问题：
1. 先查看当天的README.md中的常见问题
2. 运行 `make help` 查看可用命令
3. 随时向我提问