# 使用官方的Go镜像作为构建环境
FROM golang:1.23-alpine AS builder

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

# 使用包含MySQL的运行时镜像
FROM python:3.11

# 安装MariaDB服务器和必要的依赖
RUN apt-get update && apt-get install -y --no-install-recommends \
    mariadb-server mariadb-client \
    && rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件和所有源代码
COPY --from=builder /app/main .
COPY --from=builder /app/index.html .
COPY --from=builder /app/server.py .

# 复制MySQL初始化脚本
COPY init-mysql.sql /docker-entrypoint-initdb.d/

# 复制启动脚本
COPY start.sh .

# 暴露端口（Go服务端口和Python服务器端口）
EXPOSE 18081 18080 3306

# 设置MySQL环境变量
ENV MYSQL_ROOT_PASSWORD=123456
ENV MYSQL_DATABASE=ai_news_db
ENV MYSQL_USER=root
ENV MYSQL_PASSWORD=123456

# 设置数据库连接为本地（容器内）
ENV DB_HOST=localhost
ENV DB_PORT=3306
ENV DB_USER=root
ENV DB_PASSWORD=123456
ENV DB_NAME=ai_news_db

# 设置容器启动命令
CMD ["./start.sh"]