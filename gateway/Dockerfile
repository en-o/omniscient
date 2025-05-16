# syntax=docker.io/docker/dockerfile:1

FROM node:20.9.0 AS base

# 1. 安装依赖
FROM base AS deps

WORKDIR /app

COPY package.json package-lock.json* ./

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


# 创建数据目录并设置适当的权限
RUN mkdir -p /app/data /app/datatemp && \
    chown -R nextjs:nodejs /app/data /app/datatemp

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/datatemp ./datatemp
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

# 暴露默认的 Next.js 端口
EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"
# 禁用 Next.js 遥测数据收集
ENV NEXT_TELEMETRY_DISABLED=1

# 使用 Node.js 启动服务器
#CMD ["node", "server.js"]
# 修改启动命令，添加文件检查和移动逻辑
#CMD ["sh", "-c", "if [ ! -f /app/data/servers.db ]; then mv /app/datatemp/servers.db /app/data/; fi && node server.js"]
CMD ["sh", "-c", "if [ ! -f /app/data/servers.db ]; then cp /app/datatemp/servers.db /app/data/ && echo 'Database file copied successfully'; fi && node server.js"]

