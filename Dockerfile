# 使用官方的Go镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o main .

# 使用更小的运行时镜像
FROM python:3.11-alpine

# 安装必要的运行时依赖：Go运行环境（不再需要安装Python3）
RUN apk --no-cache add \
    ca-certificates

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件和所有源代码
COPY --from=builder /app/main .
COPY --from=builder /app/index.html .
COPY --from=builder /app/server.py .

# 暴露端口（Go服务端口和Python服务器端口）
EXPOSE 18081 18080

# 设置环境变量默认值
ENV DB_HOST=localhost
ENV DB_PORT=3306
ENV DB_USER=root
ENV DB_PASSWORD=123456
ENV DB_NAME=ai_news_db

# 创建启动脚本
RUN echo '#!/bin/sh' > start.sh && \
    echo 'echo "启动Go API服务..."' >> start.sh && \
    echo './main &' >> start.sh && \
    echo 'echo "启动Python HTTP服务器..."' >> start.sh && \
    echo 'python3 server.py' >> start.sh && \
    chmod +x start.sh

# 设置容器启动命令
CMD ["./start.sh"]