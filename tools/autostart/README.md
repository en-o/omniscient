# 自启工具
使用Systemd 服务方式构建自启管理工具


# 构建
## 克隆或下载源码
git clone <repository-url>
cd autostart

## 方法1：使用构建脚本（推荐）
chmod +x build.sh
./build.sh

## 方法2：手动构建
go build -o autostart main.go

## 安装到系统路径（可选）
sudo cp autostart /usr/local/bin/
## 或者
sudo cp build/autostart /usr/local/bin/