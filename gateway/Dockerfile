# syntax=docker.io/docker/dockerfile:1

FROM node:20.9.0 AS base

# Step 1. Rebuild the source code only when needed
FROM base AS builder
# 安装 SQLite 运行时依赖 - 在 Alpine 中正确的包名

WORKDIR /app
# 安装构建时依赖
COPY package.json package-lock.json* ./
RUN npm ci

# 先复制配置文件
COPY next.config.js* ./
COPY tsconfig.json* ./
COPY tailwind.config.js* ./
COPY postcss.config.js* ./

# 然后再复制源代码和公共资源
COPY src ./src
COPY public ./public
COPY components ./components
#COPY data ./data
COPY utils ./utils
COPY types ./types
COPY scripts ./scripts

# 运行数据库设置脚本
RUN npm run db:setup

# Environment variables must be present at build time
ARG ENV_VARIABLE
ENV ENV_VARIABLE=${ENV_VARIABLE}
ARG NEXT_PUBLIC_ENV_VARIABLE
ENV NEXT_PUBLIC_ENV_VARIABLE=${NEXT_PUBLIC_ENV_VARIABLE}

# 禁用 Next.js 遥测数据收集
ENV NEXT_TELEMETRY_DISABLED=1

# 使用 turbopack 进行生产构建（如果支持）
RUN npm run build

# Step 2. Production image, copy all the files and run next
FROM base AS runner

WORKDIR /app



# Don't run production as root
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# 复制数据库文件和脚本（确保有适当的权限）
COPY --from=builder --chown=nextjs:nodejs /app/scripts ./scripts
COPY --from=builder --chown=nextjs:nodejs /app/data ./data

# 创建数据库目录并设置适当的权限
RUN mkdir -p /app/data/db && chown -R nextjs:nodejs /app/data

USER nextjs

COPY --from=builder /app/public ./public

# Automatically leverage output traces to reduce image size
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Environment variables must be redefined at run time
ARG ENV_VARIABLE
ENV ENV_VARIABLE=${ENV_VARIABLE}
ARG NEXT_PUBLIC_ENV_VARIABLE
ENV NEXT_PUBLIC_ENV_VARIABLE=${NEXT_PUBLIC_ENV_VARIABLE}

# 禁用 Next.js 遥测数据收集
ENV NEXT_TELEMETRY_DISABLED=1

# 暴露默认的 Next.js 端口
EXPOSE 3000

# 使用 Node.js 启动服务器
CMD ["node", "server.js"]