#!/bin/sh

echo "创建MySQL所需目录..."
mkdir -p /var/run/mysqld
chown mysql:mysql /var/run/mysqld
mkdir -p /var/lib/mysql
chown mysql:mysql /var/lib/mysql

echo "初始化MySQL数据库..."
mysql_install_db --user=mysql --datadir=/var/lib/mysql

echo "启动MySQL服务..."
mysqld --user=mysql --datadir=/var/lib/mysql --socket=/var/run/mysqld/mysqld.sock --port=3306 --bind-address=0.0.0.0 --skip-networking=0 &

echo "等待MySQL完全启动..."
while ! mysqladmin -S /var/run/mysqld/mysqld.sock -uroot ping --silent; do
    sleep 1
done

echo "设置MySQL root用户权限..."
mariadb -S /var/run/mysqld/mysqld.sock -uroot -e "ALTER USER 'root'@'localhost' IDENTIFIED VIA mysql_native_password USING PASSWORD('123456');"
mariadb -S /var/run/mysqld/mysqld.sock -uroot -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;"
mariadb -S /var/run/mysqld/mysqld.sock -uroot -e "FLUSH PRIVILEGES;"

echo "创建数据库和表结构..."
mariadb -S /var/run/mysqld/mysqld.sock -uroot -p123456 -e "CREATE DATABASE IF NOT EXISTS ai_news_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mariadb -S /var/run/mysqld/mysqld.sock -uroot -p123456 ai_news_db < /docker-entrypoint-initdb.d/init-mysql.sql

echo "启动Go API服务..."
./main &

echo "启动Python HTTP服务器..."
python3 server.py