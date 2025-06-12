
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


## 数据库
1. run [schema.sql](doc/schema.sql)
2. run [table.sql](doc/table.sql)

## run
- `gf run main.go`
- `go run main.go`

## web
1. [项目管理页面](http://127.0.0.1:8000/html/pm.html)
> title 规则 ： Hostname+IP地址的最后一段
![project_ui.png](doc/images/project_ui.png)

2. [接口文档](http://127.0.0.1:8000/swagger/#tag/Java)




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