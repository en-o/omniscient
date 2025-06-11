# 自启工具
使用Systemd 服务方式构建自启管理工具

# 已测环境
> 如果 `systemctl --version` 版本比我的低可能会用不起
systemd 255 (255.4-1ubuntu8.8)
systemd 219 (CentOS Linux release 7.9.2009 (Core))

# 构建
## 克隆或下载源码
git clone <repository-url>
cd autostart

## 方法1：使用构建脚本（推荐）
构建完成后，二进制文件将位于：
- AMD64: build/<version>/amd64/autostart
- ARM64: build/<version>/arm64/autostart
```shell
chmod +x build.sh
./build.sh
```

## 方法2：手动构建
### linux 
```shell
# AMD64 架构
GOOS=linux GOARCH=amd64 go build -o build/amd64/autostart main.go
# ARM64 架构
GOOS=linux GOARCH=arm64 go build -o build/arm64/autostart main.go
```
### windows
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
