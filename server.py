#!/usr/bin/env python3
import http.server
import socketserver
import os

# 使用18080端口（与Dockerfile中暴露的端口一致）
PORT = 18080

class MyHTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        # 添加CORS头，允许跨域访问
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        # 添加缓存控制头，防止浏览器缓存静态文件
        self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
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