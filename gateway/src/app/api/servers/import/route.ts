// app/api/servers/import/route.ts
import { NextResponse } from 'next/server';
import { ServerRepository } from '@utils/database';
import { generateId } from "@utils/uuid";

// 服务器导入API
export async function POST(request: Request) {
    try {
        // 获取上传的文件数据
        const data = await request.json();

        // 验证数据格式
        if (!Array.isArray(data)) {
            return NextResponse.json(
                { error: '无效的数据格式' },
                { status: 400 }
            );
        }

        // 处理导入数据
        const results = {
            total: data.length,
            imported: 0,
            failed: 0
        };

        // 导入服务器数据
        for (const server of data) {
            try {
                // 验证必要字段
                if (!server.url || !server.description) {
                    results.failed++;
                    continue;
                }

                // 创建新记录（使用新ID避免冲突）
                const newServer = {
                    id: generateId(),
                    url: server.url,
                    description: server.description
                };

                await ServerRepository.add(newServer);
                results.imported++;
            } catch (err) {
                console.error('导入单个服务器失败:', err);
                results.failed++;
            }
        }

        return NextResponse.json({
            message: `成功导入 ${results.imported} 个服务器，失败 ${results.failed} 个`,
            results
        });
    } catch (error) {
        console.error('导入服务器数据失败:', error);
        return NextResponse.json(
            { error: '导入服务器数据失败' },
            { status: 500 }
        );
    }
}
