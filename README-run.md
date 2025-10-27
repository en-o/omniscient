避免中文路径

# run omniscient
> root 运行，要不然有些命令会失效
1. 准备 [config.prod.yaml](doc/config.prod.yaml) [跟omniscient同级] 
> - 该文件是配置文件，默认会加载当前目录下的`config.prod.yaml`，如果没有则加载内置的配置文件  
> - 默认是 sqlite3数据库，如果需要使用mysql请使用`--gf.gcfg.file=./config.mysql.yaml` ，参考[config.mysql.yaml](doc/config.mysql.yaml) 
> - `omniscient sh config [file] ` 也可以设置配置文件，file是绝对路径

```yaml

2. 注册omniscient自启
    ```shell
    chmod +x ./omniscient
    ./omniscient sh install
    ./omniscient sh enable
    ```
3. 启动 `./omniscient sh start`
4. 其他命令 `./omniscient sh help`
5. 访问页面 http://127.0.0.1:7777/html/pm.html
6. test
   ```text
   1. 构建一个jar web
   2. 在ubuntu中启动 jar
   3. 访问项目管理页面
   4. 点击注册
   5. 开始使用
   ```
7. 备注
   ```text
   1. 查看进程`ps -ef | grep  omniscient`
   2. 随编译文件构建一个配置文件使用`--gf.gcfg.file`即修改内部配置文件，如下：
   > config.prod.yaml 参考[config.yaml](manifest/config/config.yaml)
   > echo 输出 PID=$(ss -tlnp | grep ":7777" | awk '{print $6}' | cut -d',' -f2 | cut -d'=' -f2)
   
   
   1. 指定配置文件`./omniscient run --gf.gcfg.file=./config.prod.yaml` or `./omniscient sh config ./config.prod.yaml`
   2. 直接运行`./omniscient run` 默认会加载当前目录下`config.prod.yaml`配置文件[如果当前没有会加载内置的配置文件]
   3. 后台运行请看启动脚本
   ```


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
> - 尽量 root启动
> - systemctl 环境检查 `systemctl --version` 我测试时使用的最低为`systemd 219`
> - 自启的service 文件存放存地方默认 `/etc/systemd/system` `autostart-开头` 
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