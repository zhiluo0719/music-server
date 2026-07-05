#!/bin/bash
set -e
echo "=== Music Server 部署 ==="

# 配置 Docker 国内镜像加速
echo "[0/5] 配置 Docker 镜像加速..."
mkdir -p /etc/docker
if [ ! -f /etc/docker/daemon.json ]; then
    cat > /etc/docker/daemon.json << 'EOF'
{
  "registry-mirrors": [
    "https://docker.m.daocloud.io",
    "https://docker.1ms.run"
  ]
}
EOF
    systemctl restart docker 2>/dev/null || true
    echo "镜像加速已配置"
fi

# 安装 Git
echo "[1/5] 安装 Git..."
yum install -y git 2>/dev/null || true

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
echo "[4/5] 构建镜像(首次较慢，请耐心等待)..."
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

# 开放端口
firewall-cmd --add-port=3001/tcp --permanent 2>/dev/null && firewall-cmd --reload 2>/dev/null || true

echo ""
echo "部署完成! http://119.45.205.187:3001"
echo "日志: docker logs -f music-server"