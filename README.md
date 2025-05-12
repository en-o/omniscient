# 环境
> https://goframe.org.cn/docs/cli/install
```shell
# go version 
go version go1.24.2 windows/amd64
# gf version
v2.9.0
```
# 已经测过的环境
- ubuntu 22.04

# 数据库
1. run [schema.sql](doc/schema.sql)
2. run [table.sql](doc/table.sql)

# run
- `gf run main.go`
- `go run main.go`

# web
1. [项目管理页面](http://127.0.0.1:8000/html/pm.html)
> title 规则 ： Hostname+IP地址的最后一段
![project_ui.png](doc/images/project_ui.png)

2. [接口文档](http://127.0.0.1:8000/swagger/#tag/Java)



# test
1. 构建一个jar web 
2. 在ubuntu中启动 jar
3. 访问项目管理页面
4. 点击注册
5. 开始使用

# build
> 1. https://www.bilibili.com/video/BV1Uu4y1u7kX?spm_id_from=333.788.videopod.episodes&vd_source=6a1f4a95d77312275ea86329958a172f&p=46
> 2. https://goframe.org.cn/docs/cli/build
## 打包命令
打包配置信息在[config.yaml](hack/config.yaml)
> - system = linux,darwin,windows
> - arch = 386,amd64,arm,arm64 [uname -s/uname -m]
>   - ps: amd64 = x86_64
>   https://juejin.cn/post/7097032561092165640
```shell
gf build
```
## 启动脚本
1. 查看进程`ps -ef | grep  omniscient`
2. 随编译文件构建一个配置文件使用`--gf.gcfg.file`即修改内部配置文件，如下：
```shell
#!/bin/bash

# 定义端口号
PORT=6001

# 根据端口杀死进程
PID=$(netstat -nlp | grep :$PORT | awk '{print $7}' | cut -d/ -f1)
if [ -n "$PID" ]; then
    echo "Killing process $PID"
    kill -9 $PID
fi

# 赋予 ./reURL 文件可执行权限
chmod +x ./omniscient

# 执行 nohup ./reURl --gf.gcfg.file=./config.pro.yaml > nohup.log 2>&1 & 进行部署
nohup ./omniscient --gf.gcfg.file=./config.pro.yaml > nohup.log 2>&1 &

echo "Deployment completed"
```