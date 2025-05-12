# 聚合html
1. 在不干扰原本单个服务使用的同时将多个`omniscient` 服务聚合在一起显示

# env
```shell
# node -v
v20.9.0
# npm -v
10.1.0
```

# use 
```shell
# Start
node server.js
```


# import
```shell
npm install express axios body-parser cors express-sse
# axios：从 Node.js 后端向后端发送 HTTP 请求
# body-parser：用于解析 JSON 请求体
# cors：用于处理跨域资源共享（在开发过程中很有用）
# express-sse：用于轻松处理服务器发送事件代理
```



# create 
npm init -y


## 目录结构
```text
project-manager-node/
├── public/
│   ├── css/
│   │   ├── bootstrap.min.css
│   │   ├── bootstrap-icons.css
│   │   └── pm.css
│   ├── js/
│   │   └── bootstrap.bundle.min.js
│   └── index.html       # Modified frontend HTML
├── server.js            # Node.js backend logic
└── package.json
```
