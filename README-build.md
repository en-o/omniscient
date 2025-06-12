# build omniscient
> 1. https://www.bilibili.com/video/BV1Uu4y1u7kX?spm_id_from=333.788.videopod.episodes&vd_source=6a1f4a95d77312275ea86329958a172f&p=46
> 2. https://goframe.org.cn/docs/cli/build

打包配置信息在[config.yaml](hack/config.yaml)
> - system = linux,darwin,windows
> - arch = 386,amd64,arm,arm64 [uname -s/uname -m]
    >   - ps: amd64 = x86_64
          >   https://juejin.cn/post/7097032561092165640
```shell
cd omniscient
gf build
```

# build gateway  for docker
> 1. https://github.com/vercel/next.js/tree/canary/examples/with-docker-multi-env
> 2. 我这个项目不能用aline镜像，slim也不行，sqlite安装会出问题 （😔
> 3. 目前镜像有点大，1.6g 但是 load 下来只有400mb
```bash
cd gateway
# build的时候注意 package-lock.json，我换个环境重新生成就出问题了[如果重新生成，请把node_modules先删除]
# --no-cache 禁止缓存
# docker builder prune # 清理缓存
#docker build --no-cache  -t  tannnn/omniscient-gateway:0.0.1 .
#docker build  -t  tannnn/omniscient-gateway:0.0.1 .
#注意不支持 linux/arm/v7 ， 请自己适配
docker  build --platform linux/amd64,linux/arm64/v8 -t tannnn/omniscient-gateway:0.0.1 .
docker run -p 3000:3000 tannnn/omniscient-gateway:0.0.1
```

# build autostart for 手动构建
cd tools/autostart
1. linux
```shell
# AMD64 架构
GOOS=linux GOARCH=amd64 go build -o build/amd64/autostart main.go
# ARM64 架构
GOOS=linux GOARCH=arm64 go build -o build/arm64/autostart main.go
```
2. windows
```shell 
# PowerShell
## AMD64 架构
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o build/amd64/autostart main.go

## ARM64 架构
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o build/arm64/autostart main.go

#  CMD
## AMD64 架构
set GOOS=linux && set GOARCH=amd64 && go build -o build/amd64/autostart main.go

## ARM64 架构
set GOOS=linux && set GOARCH=arm64 && go build -o build/arm64/autostart main.go
```