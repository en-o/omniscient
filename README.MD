# run
- `gf run main.go`
- `go run main.go`

## API address
1. [项目管理页面](http://127.0.0.1:8000/html/pm.html)
1. [接口文档](http://127.0.0.1:8000/swagger/#tag/Java)

# build 
> https://goframe.org.cn/docs/cli/build
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

# 文档
- [GoFrame github](https://github.com/gogf/gf)
- [GoFrame 中文doc](https://goframe.org.cn/quick/install)
- [代码生成](https://goframe.org.cn/quick/scaffold-api-gen-dao#%E6%89%A7%E8%A1%8C%E4%BB%A3%E7%A0%81%E7%94%9F%E6%88%90)
```shell
# dao
gf gen dao
# 接口
gf gen ctrl
```

# 目录
https://goframe.org.cn/docs/design/project-structure
# 项目目录结构说明

| 目录/文件名称       | 说明                | 描述                                                                 |
|---------------------|---------------------|----------------------------------------------------------------------|
| **api**             | 对外接口            | 对外提供服务的输入/输出数据结构定义。考虑到版本管理需要，往往以 `api/xxx/v1...` 存在。 |
| **hack**            | 工具脚本            | 存放项目开发工具、脚本等内容（如 CLI 工具配置、shell/bat 脚本等）。               |
| **internal**        | 内部逻辑            | 业务逻辑存放目录（通过 Golang internal 特性对外部隐藏可见性）。                  |
| ├─ cmd              | 入口指令            | 命令行管理目录，可维护多个命令行。                                               |
| ├─ consts           | 常量定义            | 项目所有常量定义。                                                             |
| ├─ controller       | 接口处理            | 接收/解析用户输入参数的入口/接口层。                                            |
| ├─ dao              | 数据访问            | 数据访问对象（抽象层），仅包含最基础的 CRUD 方法。                                |
| ├─ model            | 结构模型            | 数据结构管理模块，管理数据实体对象及输入/输出数据结构。                            |
| │  ├─ do            | 领域对象            | 用于 dao 数据操作中的业务模型与实例模型转换（由工具维护，用户不可修改）。            |
| │  └─ entity        | 数据模型            | 数据模型与数据集合的一对一关系（由工具维护，用户不可修改）。                        |
| └─ service          | 业务实现            | 业务逻辑的具体实现和封装（项目中最复杂的部分）。                                   |
| **manifest**        | 交付清单            | 程序编译、部署、运行、配置的相关文件。                                           |
| ├─ config           | 配置管理            | 配置文件存放目录。                                                            |
| ├─ docker           | 镜像文件            | Docker 镜像相关依赖文件及脚本。                                                |
| ├─ deploy           | 部署文件            | 部署相关文件（默认提供 Kubernetes 的 Yaml 模板，通过 kustomize 管理）。          |
| └─ protobuf         | 协议文件            | GRPC 使用的 protobuf 协议定义文件（编译后生成 go 文件到 api 目录）。              |
| **resource**        | 静态资源            | 静态资源文件（可通过资源打包/镜像编译注入到发布文件中）。                           |
| **go.mod**          | 依赖管理            | Go Module 包管理的依赖描述文件。                                               |
| **main.go**         | 入口文件            | 程序入口文件。                                                                |


# todo
- [ ] html优化
  - [x] 下拉选择变成列表
  - [x] 输出脚本日志
  - [x] 新增注册在线项目按钮
  - [x] 操作完后刷新列表
- [x] 运行启动后,更新的pid
- [x] 运行run失败
- [ ] 原生运行停止延迟问题
- [ ] 原生后台运行需要输出一段时间的日志后结束，并且需要正确更新pid
- [ ] 新增结构
  - [ ] 新增表(服务器mac,服务器名(可修改，默认系统名+最后一位IP),服务器ip,备注)
  - [ ] 新增选择栏目，选择服务器
  - [ ] 通过选择的服务器加载项目列表
- [ ] 打包脱机运行
