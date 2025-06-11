# 自启工具
使用Systemd 服务方式构建自启管理工具

# 已测环境
> 如果 `systemctl --version` 版本比我的低可能会用不起
systemd 255 (255.4-1ubuntu8.8)
systemd 219 (CentOS Linux release 7.9.2009 (Core))

# 注意事项
1. 必须使用 sudo
2. 可执行文件写绝对路径
3. 可执行文件执行命令不允许后台启动，例如：不要使用 `&` 或 `nohup`
4. 必须设置工作目录，请使用 `--workdir` 参数，即可执行文件所在目录

```shell
# 安装命令：  add 项目名  可执行文件路径和执行方式  --workdir=执行文件所在目录
sudo ./autostart add myapp "java -jar /mnt/c/Test/omniscient_test-0.0.1-SNAPSHOT.jar" --workdir=/mnt/c/Test
```
# 全局环境设置
```shell
# chmod +x autostart
# 安装到全局环境
sudo ./autostart install-global

# 卸载全局环境
sudo ./autostart uninstall-global

# 查看帮助（包含新命令）
autostart help
```

# 构建
## 克隆或下载源码
git clone <repository-url>
cd autostart

## 方法1：使用构建脚本（推荐）
构建完成后，二进制文件将位于：
- AMD64: build/<version>/amd64/autostart
- ARM64: build/<version>/arm64/autostart
```shell
chmod +x autostart
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
