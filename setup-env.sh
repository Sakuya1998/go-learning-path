#!/bin/bash
# Go学习环境设置脚本

echo "设置Go 1.26开发环境..."

# 设置环境变量
export GOROOT=/home/node/.go
export GOPATH=/home/node/go
export PATH=$GOROOT/bin:$PATH

# 创建必要的目录
mkdir -p $GOPATH/{src,bin,pkg}

# 验证安装
echo "Go版本:"
go version

echo "环境变量:"
echo "GOROOT=$GOROOT"
echo "GOPATH=$GOPATH"
echo "PATH中包含Go: $(echo $PATH | grep -o '.go/bin' || echo '未找到')"

# 进入项目目录
cd /home/node/.openclaw/workspace/go-learning-path

echo "环境设置完成！"
echo "运行以下命令开始学习:"
echo "  cd week1/day1"
echo "  go run goroutine_exercise.go"