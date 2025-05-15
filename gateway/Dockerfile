# syntax=docker.io/docker/dockerfile:1

FROM node:20.9.0 AS base

# 1. 安装依赖
FROM base AS deps

WORKDIR /app

COPY package.json package-lock.json* ./
# 确保使用正确的平台构建原生模块
RUN npm ci


# 2. 打包
FROM base AS builder
WORKDIR /app
# 从 deps 中复制 node_modules
COPY --from=deps /app/node_modules ./node_modules
# 复制源代码（ 会自动过滤掉.dockerignore中的内容）
COPY . .

# **添加这一行：创建目标目录**
RUN mkdir -p /app/datatemp
# 运行数据库设置脚本
RUN npm run db:docker


RUN ls -l /app/scripts/setup-db-docker.js
RUN ls -l /app/datatemp/servers.db


RUN npm run build

# 3. 运行
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production

RUN groupadd -r -g 1001 nodejs && \
    useradd -r -u 1001 -g nodejs nextjs



RUN mkdir -p /app/data && chown nextjs:nodejs /app/data

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/datatemp ./datatemp
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

# 暴露默认的 Next.js 端口
EXPOSE 3000

ENV PORT=3000

# 使用 Node.js 启动服务器
#CMD ["node", "server.js"]
# 修改启动命令，添加文件检查和移动逻辑
CMD ["sh", "-c", "if [ ! -f /app/data/servers.db ]; then mv /app/datatemp/servers.db /app/data/; fi && node server.js"]


