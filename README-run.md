# run omniscient
> root 运行，要不然有些命令会失效
1. 准备 config.prod.yaml [跟omniscient同级]
2. 注册omniscient自启
    ```shell
    chmod +x ./omniscient
    ./omniscient sh install
    ./omniscient sh enable
    ```
3. 启动 `./omniscient sh start`
4. 其他命令 `./omniscient sh help`
5. 访问页面 http://127.0.0.1:7777/html/pm.html


# run gateway  for docker
> 详见：https://gitee.com/tanoo/docker-compose/tree/tn/omniscient

ACCESS: http://127.0.0.1:3000
```shell
docker run -d \
  --name omniscient-gateway \
  --privileged \
  --restart unless-stopped \
  -p 3000:3000 \
  -v ./data:/app/data \
  -e TZ=Asia/Shanghai \
  -e LANG=en_US.UTF-8 \
  tannnn/omniscient-gateway:0.0.2
```

# run autostart
> 尽量 root启动
> systemctl 环境检查 `systemctl --version` 我测试时使用的最低为`systemd 219`
1. 设置全局变量
    ```shell
    chmod +x autostart
    
    # 安装到全局环境
    sudo ./autostart install-global
    
    # 卸载全局环境
    sudo ./autostart uninstall-global
    
    # 查看帮助（包含新命令）
    autostart help
    ```
2. 注册项目自启
> - 必须使用 sudo(root可不用)
> - 可执行文件写绝对路径
> - 可执行文件执行命令不允许后台启动，例如：不要使用 `&` 或 `nohup`
> - 必须设置工作目录，请使用 `--workdir` 参数，即可执行文件所在目录
> - rm会先停止项目的运行，然后在删除自启 
> - 目前我就测试了java 和 sh 命令
```shell
sudo autostart add myapp "java -jar /home/software/jar/omniscient_test/omniscient_test-0.0.1-SNAPSHOT.jar" --workdir=/home/software/jar/omniscient_test

sudo autostart add myapp "sh /home/software/jar/omniscient_test/run2.sh start -b false" --workdir=/home/software/jar/omniscient_test
```
3. 启用自启
> 其他命令请 `autostart -h`
```shell
sudo autostart enable myapp
```