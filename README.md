```text
├─omniscient # 主项目 java进程管理
│  ├─api # 接口
│  └─manifest # 配置文件
│  └─internal # 进程主逻辑
├─tools #  其他项目- 跟主项目的逻辑无关只是放在了一起开发
│  └─autostart # 自启工具 - 跟主项目的逻辑无关只是放在了一起开发
└─gateway # 进行网关 - 跟主项目的逻辑无关只是放在了一起开发
└─release # 当前发行版

```


#  java进程管理

## 环境
> https://goframe.org.cn/docs/cli/install
```shell
# go version 
go version go1.24.2 windows/amd64
# gf version
v2.9.0
```
## 已经测过的环境
- ubuntu 22.04
- CentOS Linux release 7.9.2009 (Core)

## 配置注释
1. 默认数据库使用sqlite3, 文件在 `./data/` 想改路径自己在config中设置
2. 如果需要使用mysql，请修改[config.yaml](manifest/config/config.yaml)中的数据库配置
3. 二进制文件需要使用mysql的时候请参考[config.mysql.yaml](doc/config.mysql.yaml)


## 数据库
1. mysql的需要自行创建数据库，run [schema.sql](doc/schema.sql)， 至于表结构他会自动创建
2. sqlite的会自动创建数据库文件和表结构

## run
- `gf run main.go`
- `go run main.go`

## web
1. [项目管理页面](http://127.0.0.1:8000/html/pm.html)
> title 规则 ： Hostname+IP地址的最后一段
![project_ui.png](doc/images/project_ui.png)

2. [接口文档](http://127.0.0.1:8000/swagger/#tag/Java)

# 自启工具
使用Systemd 服务方式构建自启管理工具

## 已测环境
> 如果 `systemctl --version` 版本比我的低可能会用不起
systemd 255 (255.4-1ubuntu8.8)
systemd 219 (CentOS Linux release 7.9.2009 (Core))

## 查看使用文档
autostart help

# 聚合网关
> 进行管理工具的前端集成
## script
```bash
# dir
cd gateway
# install
npm install
# run
npm run dev
```
## access url
http://127.0.0.1:3000

## 环境
```shell
# react
^19.0.0
# next
15.3.2
# node -v
v20.9.0  or  v22.14.0  
# npm -v
10.1.0
```


# 问题注意
## docker 启动的java会造成干扰
会抓取到docker启动jar程序，但是获取的目录是docker内部的，这里需要注意一下，docker的进程不要手贱去操作了
```shell
# 使用这个查看抓取的 java进程 是不是就是docker的，
docker top <容器名称或ID>
# 查看映射
docker inspect --format '{{json .Mounts}}' <容器名或ID>
```

# 自启备注
[omniscient.service; enabled; vendor preset: disabled](https://www.yuque.com/tanning/mbquef/zi21spxc6l5nwazh)