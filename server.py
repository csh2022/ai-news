#!/usr/bin/env python3
import http.server
import socketserver
import os
import time

# 使用18080端口（与Dockerfile中暴露的端口一致）
PORT = 18080

class MyHTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        # 设置CORS头
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        
        # 根据文件类型设置不同的缓存控制策略
        if self.path.endswith('.xml') or self.path == '/robots.txt':
            # 对于sitemap.xml和robots.txt设置较短的缓存时间，确保搜索引擎能获取最新版本
            self.send_header('Cache-Control', 'public, max-age=86400')
            self.send_header('Expires', self.date_time_string(time.time() + 86400))
        else:
            # 对于其他静态文件，防止缓存
            self.send_header('Cache-Control', 'no-store, no-cache, must-revalidate, max-age=0')
            self.send_header('Pragma', 'no-cache')
            self.send_header('Expires', '0')
        super().end_headers()

# 切换到项目根目录
os.chdir(os.path.dirname(os.path.abspath(__file__)))

with socketserver.TCPServer(("", PORT), MyHTTPRequestHandler) as httpd:
    print(f"Python HTTP服务器启动在端口 {PORT}")
    print(f"前端页面访问地址: http://localhost:{PORT}")
    print(f"API服务访问地址: http://localhost:18081/api/news")
    print(f"外网访问地址: http://<您的IP地址>:{PORT}")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\n服务器已停止")