# 聚合html
1. 在不干扰原本单个服务使用的同时将多个`omniscient` 服务聚合在一起显示

# use
## script
```bash
# dir
cd gateway
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
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
│   └── pages/  # 页面路由器
└── package.json
```