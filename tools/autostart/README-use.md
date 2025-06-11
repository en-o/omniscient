# Test
```shell
sudo ./autostart add myapp "java -jar /mnt/c/Test/omniscient_test-0.0.1-SNAPSHOT.jar" --workdir=/mnt/c/Test
```


# 案例
## Java 应用
```shell
sudo autostart add myapp "java -jar /path/to/app.jar --server.port=8080" --workdir=/path/to --user=myuser
```

## Python 应用，带环境变量
```shell
sudo autostart add pyapp "python3 /path/to/app.py" --user=www-data --env=PYTHONPATH=/path/to --env=DEBUG=true
```

## Node.js 应用，自定义重启策略
```shell
sudo autostart add nodeapp "node /path/to/app.js" --restart=on-failure --restart-sec=10
```

## Go 应用，带服务依赖
```shell
sudo autostart add webapp "./webapp" --after=network.target --after=mysql.service --requires=mysql.service
```

## Shell 脚本
```shell
sudo autostart add backup "bash /path/to/backup.sh" --user=backup --workdir=/backups
```

## 查看服务日志
```shell
autostart logs myapp 100
```

## 查看服务配置
```shell
sudo autostart edit myapp
```