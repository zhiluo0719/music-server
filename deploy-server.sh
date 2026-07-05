#!/bin/bash
set -e
echo "=== Music Server 部署 (无 Docker 版) ==="

# 1. 安装 Go
echo "[1/6] 安装 Go..."
if ! command -v go &>/dev/null; then
    yum install -y wget git
    wget -q https://go.dev/dl/go1.24.0.linux-amd64.tar.gz -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    export PATH=$PATH:/usr/local/go/bin
    rm /tmp/go.tar.gz
fi
echo "Go: $(go version)"

# 2. 配置 Go 国内代理
echo "[2/6] 配置 Go 代理..."
go env -w GOPROXY=https://goproxy.cn,direct

# 3. 创建目录
echo "[3/6] 创建目录..."
mkdir -p /data/music/{uploads/audio,uploads/covers,data,logs}

# 4. 拉取代码
echo "[4/6] 拉取代码..."
if [ -d "/data/music/music-server" ]; then
    cd /data/music/music-server
    git pull
    git checkout main
else
    cd /data/music
    git clone https://github.com/zhiluo0719/music-server.git
    cd music-server
fi

# 5. 编译
echo "[5/6] 编译..."
export PATH=$PATH:/usr/local/go/bin
go mod tidy
CGO_ENABLED=0 go build -o server ./cmd/server/
echo "编译完成: $(ls -lh server)"

# 6. 创建服务并启动
echo "[6/6] 部署服务..."

# 停止旧进程
pkill -f "./server" 2>/dev/null || true

# 复制到运行目录
cp server /data/music/
cp -r public /data/music/ 2>/dev/null || true

# 创建 systemd 服务
cat > /etc/systemd/system/music-server.service << 'EOF'
[Unit]
Description=Music Server
After=network.target

[Service]
Type=simple
WorkingDirectory=/data/music
ExecStart=/data/music/server
Restart=always
RestartSec=5
Environment=PORT=3001
Environment=TZ=Asia/Shanghai

[Install]
WantedBy=multi-user.target
EOF

# 开放端口
firewall-cmd --add-port=3001/tcp --permanent 2>/dev/null && firewall-cmd --reload 2>/dev/null || true

# 启动
systemctl daemon-reload
systemctl enable music-server
systemctl restart music-server

echo ""
echo "=== 部署完成! ==="
echo "访问: http://119.45.205.187:3001"
echo "状态: systemctl status music-server"
echo "日志: journalctl -u music-server -f"