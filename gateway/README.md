# 聚合html
ps: 用的浏览器缓存，缓存被删就没了更换浏览器也没了

1. 在不干扰原本单个服务使用的同时将多个`omniscient` 服务聚合在一起显示

# init 

## 运行数据库脚本
```shell
cd gateway
node scripts/setup-db.js 
// or
npm run db:setup
```

## sqlite3 安装问题
> npm install 失败
1. [安装 Visual Studio 及 C++ 工具链](https://visualstudio.microsoft.com/zh-hans/visual-cpp-build-tools/)
```text
1. 选择“使用 C++ 的桌面开发”工作负载。
2. 确保安装了以下组件：
2.1 Visual C++ 工具集
2.2 Windows SDK
2.3 CMake 工具
2.4 适用于 C++ 的测试工具
```
2. 安装 Python
```shell
# 不想手动安装 Visual Studio 和 Python，可以使用 windows-build-tools 来自动安装
npm install --global windows-build-tools
```
4. 清理 npm 缓存并重新安装
```shell
npm cache clean --force
npm install
```
5. 如果编译仍然失败，可以尝试安装预编译的 sqlite3 二进制文件
```shell
npm install sqlite3 --build-from-source=false
```

# use

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

# 环境
```shell
# react
^19.0.0
# next
15.3.2
# node -v
v20.9.0
# npm -v
10.1.0
```

# 文档
1. [Deploy now](https://vercel.com/new?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app)
2. [Read our docs](https://nextjs.org/docs?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app)
3. [Learn](https://nextjs.org/learn?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app)
4. [Examples](https://vercel.com/templates?framework=next.js&utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app)

# 项目结构
> https://nextjs.org/docs/app/getting-started/project-structure
```text
gateway/
├── public/  # 静态资源
│   ├── css/
│   ├── js/
├── src/  # 应用程序源文件夹
│   ├── app/  # 应用路由器
└── package.json
```

# 样式参考
https://github.com/HumeAI/hume-evi-next-js-starter
