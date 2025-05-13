import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    async headers() {
        return [
            {
                source: '/:path*',
                headers: [
                    {
                        key: 'Access-Control-Allow-Origin',
                        value: '*',
                    },
                ],
            },
        ]
    },
    // 允许所有主机的请求
    experimental: {
        // 移除不支持的 allowedDevOrigins 配置
    },
};

export default nextConfig;
