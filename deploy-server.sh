#!/bin/bash
set -e
echo "=== Music Server 部署 ==="

# 安装依赖
echo "[0/5] 安装依赖..."
yum install -y git 2>/dev/null || apt-get install -y git 2>/dev/null

# 安装 Docker
if ! command -v docker &>/dev/null; then
    echo "[1/5] 安装 Docker..."
    yum install -y yum-utils
    yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    yum install -y docker-ce docker-ce-cli containerd.io
    systemctl start docker
    systemctl enable docker
    echo "Docker: $(docker --version)"
else
    echo "[1/5] Docker 已安装: $(docker --version)"
fi

# 创建目录
echo "[2/5] 准备目录..."
mkdir -p /data/music/{uploads/audio,uploads/covers,data,logs}

# 克隆代码
echo "[3/5] 拉取代码..."
if [ -d "/data/music/music-server" ]; then
    cd /data/music/music-server && git pull
else
    cd /data/music && git clone https://github.com/zhiluo0719/music-server.git
fi

# 构建
echo "[4/5] 构建镜像..."
cd /data/music/music-server
docker build -t music-server .

# 启动
echo "[5/5] 启动服务..."
docker stop music-server 2>/dev/null || true
docker rm music-server 2>/dev/null || true
docker run -d --name music-server --restart always \
    -p 3001:3001 \
    -v /data/music/uploads:/app/uploads \
    -v /data/music/data:/app/data \
    -v /data/music/logs:/app/logs \
    -e PORT=3001 -e TZ=Asia/Shanghai \
    music-server

echo ""
echo "部署完成! http://119.45.205.187:3001"
echo "日志: docker logs -f music-server"