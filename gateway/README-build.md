# 默认部署
> 不建议node_modules太大了
## 打包
```shell
# 打包完成后，需要拷贝.next、node_modules、public、package.json四个文件到服务器
npm run build
```
## 运行
```shell
npm run start
```

# 静态打包
> [使用这个分支打包，这个分支没有使用 sqlLite](https://gitee.com/etn/omniscient/tree/gateway%E9%9D%99%E6%80%81%E9%A1%B5%E9%9D%A2/)
## 添加配置
配置文件next.config.mjs中添加`output: 'export'`
ps: `output: 'export'`, 移除 headers 配置，因为静态导出不支持
```ts
const nextConfig: NextConfig = {
    output: 'export',
    images: {
        unoptimized: true, // 静态导出时需要禁用图片优化
    },
    // output: 'export', 移除 headers 配置，因为静态导出不支持
    // async headers() {
    //     return [
    //         {
    //             source: '/:path*',
    //             headers: [
    //                 {
    //                     key: 'Access-Control-Allow-Origin',
    //                     value: '*',
    //                 },
    //             ],
    //         },
    //     ]
    // },
    // 允许所有主机的请求
    experimental: {
        // 移除不支持的 allowedDevOrigins 配置
    },
};
```
## 打包
```shell
# 打包完成后，文件在out文件夹下，拷贝out到服务器
npm run build
```
## 运行
通过nginx等代理
```nginx configuration
server{
    listen 80;
    server_name your.domain.com;
    location / {
        index index.html index.htm;
        try_files $uri $uri/ /index.html;
    }
}
```

# standalone单独部署
## 添加配置
配置文件next.config.mjs中添加`output: 'standalone'`
## 打包
```shell
# 打包完成后，在.next文件夹中生成了standalone和static文件夹
# standalone相当于一个express应用，包含server.js和少量的服务端依赖node_modules
# static是前端包
npm run build
```
## 运行
1. 需要拷贝`static`到`standalone/.next`文件夹下，拷贝外层`public`文件夹到`standalone`下
2. 进入standalone目录，执行node server.js


## Building for docker
> https://github.com/vercel/next.js/tree/canary/examples/with-docker-multi-env
```bash
docker build  -t  tannnn/omniscient-gateway:0.0.1 .
docker run -p 3000:3000 tannnn/omniscient-gateway:0.0.1
```
