# syntax=docker.io/docker/dockerfile:1

FROM node:20.9.0-slim AS base

# 1. Dependencies stage
FROM base AS deps
WORKDIR /app

# Combine apt commands and cleanup in single layer
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    python3 \
    make \
    g++ && \
    rm -rf /var/lib/apt/lists/*

# Install dependencies
COPY package*.json ./
RUN npm install  --build-from-source

# 2. Builder stage
FROM base AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# 3. Runner stage
FROM base AS runner
WORKDIR /app

# Create non-root user more efficiently
RUN groupadd --gid 1001 nodejs && \
    useradd --uid 1001 --gid nodejs --shell /bin/bash --create-home nextjs

# Set proper permissions
COPY --from=builder --chown=nextjs:nodejs /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000
ENV PORT=3000 \
    NODE_ENV=production \
    HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
