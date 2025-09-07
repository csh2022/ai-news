#!/bin/bash

# Docker本地测试启动脚本
set -e

# 设置默认值
IMAGE_NAME="ai-news-api"
CONTAINER_NAME="ai-news-container"
PORT="8080"

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
  -p $PORT:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=3306 \
  -e DB_USER=root \
  -e DB_PASSWORD=123456 \
  -e DB_NAME=ai_news_db \
  $IMAGE_NAME

# 显示容器状态
echo "容器启动成功！"
echo "容器名称: $CONTAINER_NAME"
echo "本地访问地址: http://localhost:$PORT"
echo "API端点: http://localhost:$PORT/api/news"

# 显示日志
echo "正在显示容器日志..."
docker logs -f $CONTAINER_NAME
