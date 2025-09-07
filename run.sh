#!/bin/bash

# Docker本地测试启动脚本
set -ex

# 设置默认值
IMAGE_NAME="ai-news-api"
CONTAINER_NAME="ai-news-container"

# 构建Docker镜像
echo "构建Docker镜像..."
docker build -t $IMAGE_NAME .

# 停止并删除现有容器（如果存在）
echo "清理现有容器..."
docker stop $CONTAINER_NAME 2>/dev/null || true
docker rm $CONTAINER_NAME 2>/dev/null || true

# 运行容器
echo "启动容器..."
docker run -d \
  --name $CONTAINER_NAME \
  -p 18080:18080 \
  -p 18081:18081 \
  $IMAGE_NAME

# 显示容器状态
echo "容器启动成功！"
echo "容器名称: $CONTAINER_NAME"
echo "前端页面访问地址: http://localhost:18080"
echo "API端点: http://localhost:18081/api/news"
echo "MySQL端口: localhost:3306"

# 显示日志
echo "正在显示容器日志..."
docker logs -f $CONTAINER_NAME