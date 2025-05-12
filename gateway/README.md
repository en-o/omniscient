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
http://127.0.0.1:3000/index.html


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
