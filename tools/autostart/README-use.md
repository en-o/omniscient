# Java 应用
sudo autostart add myapp "java -jar /path/to/app.jar --server.port=8080" --workdir=/path/to --user=myuser

# Python 应用，带环境变量
sudo autostart add pyapp "python3 /path/to/app.py" --user=www-data --env=PYTHONPATH=/path/to --env=DEBUG=true

# Node.js 应用，自定义重启策略
sudo autostart add nodeapp "node /path/to/app.js" --restart=on-failure --restart-sec=10

# Go 应用，带服务依赖
sudo autostart add webapp "./webapp" --after=network.target --after=mysql.service --requires=mysql.service

# Shell 脚本
sudo autostart add backup "bash /path/to/backup.sh" --user=backup --workdir=/backups

# 查看服务日志
autostart logs myapp 100

# 查看服务配置
sudo autostart edit myapp