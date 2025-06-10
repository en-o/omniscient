#!/bin/bash

# 部署脚本 - 支持 Systemd 服务管理 , 只是临时文件，从来没有测试过
# 使用方法: ./deploy.sh [port] [config_file]

# 默认配置
DEFAULT_PORT=8001
DEFAULT_CONFIG="config.prod.yaml"
SERVICE_NAME="omniscient"

# 获取参数
PORT=${1:-$DEFAULT_PORT}
CONFIG_FILE=${2:-$DEFAULT_CONFIG}

echo "=== Omniscient Service Deployment ==="
echo "Port: $PORT"
echo "Config: $CONFIG_FILE"
echo "Service: $SERVICE_NAME"
echo "====================================="

# 检查配置文件是否存在
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Warning: Config file $CONFIG_FILE not found!"
    echo "Please ensure the config file exists before starting the service."
fi

# 停止现有服务（如果存在）
echo "Stopping existing service..."
sudo systemctl stop $SERVICE_NAME 2>/dev/null || true

# 传统方式停止进程（兼容性）
echo "Checking for processes on port $PORT..."
PID=$(ss -tlnp | grep ":${PORT}" | awk '{print $6}' | cut -d',' -f2 | cut -d'=' -f2 | sort -u)

if [ -n "$PID" ]; then
    echo "Killing process $PID on port $PORT"
    kill -9 $PID 2>/dev/null || true
    sleep 2
else
    echo "No process found on port $PORT"
fi

# 赋予可执行权限
echo "Setting executable permissions..."
chmod +x ./omniscient

# 检查 omniscient 文件是否存在
if [ ! -f "./omniscient" ]; then
    echo "Error: omniscient executable not found!"
    echo "Please build the project first: gf build"
    exit 1
fi

# 选择部署方式
echo ""
echo "Choose deployment method:"
echo "1) Systemd service (recommended)"
echo "2) Traditional nohup"
echo "3) Install service only (don't start)"
read -p "Enter choice (1-3): " choice

case $choice in
    1)
        echo "Deploying as systemd service..."

        # 安装服务
        echo "Installing systemd service..."
        sudo ./omniscient sh install

        if [ $? -eq 0 ]; then
            echo "Service installed successfully!"

            # 启用开机自启
            echo "Enabling auto-start..."
            sudo ./omniscient sh enable

            # 启动服务
            echo "Starting service..."
            sudo ./omniscient sh start

            # 检查状态
            echo "Service status:"
            sudo ./omniscient sh status

            echo ""
            echo "=== Service Management Commands ==="
            echo "Check status:    ./omniscient sh status"
            echo "Start service:   ./omniscient sh start"
            echo "Stop service:    ./omniscient sh stop"
            echo "Restart service: ./omniscient sh restart"
            echo "Enable auto-start: ./omniscient sh enable"
            echo "Disable auto-start: ./omniscient sh disable"
            echo "View logs:       journalctl -u $SERVICE_NAME -f"
            echo "=================================="
        else
            echo "Failed to install service, falling back to nohup..."
            choice=2
        fi
        ;;
    2)
        echo "Deploying with nohup..."
        nohup ./omniscient run --gf.gcfg.file=./$CONFIG_FILE > nohup.log 2>&1 &
        NEW_PID=$!
        echo "Deployment completed! PID: $NEW_PID, PORT: $PORT"
        echo "Log file: nohup.log"
        echo "Stop with: kill $NEW_PID"
        ;;
    3)
        echo "Installing service only..."
        sudo ./omniscient sh install
        echo "Service installed but not started."
        echo "Use './omniscient sh start' to start the service."
        ;;
    *)
        echo "Invalid choice, exiting..."
        exit 1
        ;;
esac

echo ""
echo "Deployment completed!"
echo "Access the service at: http://localhost:$PORT"
echo "Project management: http://localhost:$PORT/html/pm.html"
echo "API documentation: http://localhost:$PORT/swagger/"