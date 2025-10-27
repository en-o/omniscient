#!/bin/bash

# 定义端口号, 跟配置文件中的端口一致
PORT=7777

# 根据端口号查找进程并获取 PID
PID=$(ss -tlnp | grep ":${PORT}" | awk '{print $6}' | cut -d',' -f2 | cut -d'=' -f2 | sort -u)

if [ -n "$PID" ]; then
    echo "Killing process $PID"
    kill -9 $PID
else
    echo "No process found on port $PORT"
fi

# 赋予 ./omniscient 文件可执行权限
chmod +x ./omniscient

# 执行 nohup ./omniscient --gf.gcfg.file=./config.prod.yaml > nohup.log 2>&1 & 进行部署
nohup ./omniscient --gf.gcfg.file=./config.prod.yaml > nohup.log 2>&1 &

# 获取新启动的进程 PID
NEW_PID=$!

echo "Deployment completed! PID: $NEW_PID, PORT: $PORT"
