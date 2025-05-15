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

# 运行数据库设置脚本
RUN npm run db:setup

RUN npm run build

# 3. 运行
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production

RUN groupadd -r -g 1001 nodejs && \
    useradd -r -u 1001 -g nodejs nextjs

COPY --from=builder /app/public ./public
#
## Automatically leverage output traces to reduce image size
## https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

# 暴露默认的 Next.js 端口
EXPOSE 3000

ENV PORT=3000

# 使用 Node.js 启动服务器
CMD ["node", "server.js"]