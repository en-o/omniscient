// app/api/servers/route.ts
import { NextResponse } from 'next/server';
import { ServerRepository } from '@utils/database';
import { generateId } from "@utils/uuid";

// 获取所有服务器
export async function GET() {
    try {
        const servers = await ServerRepository.getAll();
        return NextResponse.json(servers);
    } catch (error) {
        console.error('获取服务器列表失败:', error);
        return NextResponse.json(
            { error: '获取服务器列表失败' },
            { status: 500 }
        );
    }
}

// 添加新服务器
export async function POST(request: Request) {
    try {
        const body = await request.json();

        // 验证请求数据
        if (!body.url || !body.description) {
            return NextResponse.json(
                { error: '服务器URL和描述为必填项' },
                { status: 400 }
            );
        }

        // 创建新服务器记录
        const newServer = {
            id: generateId(),
            url: body.url,
            description: body.description
        };

        await ServerRepository.add(newServer);

        return NextResponse.json(newServer, { status: 201 });
    } catch (error) {
        console.error('添加服务器失败:', error);
        return NextResponse.json(
            { error: '添加服务器失败' },
            { status: 500 }
        );
    }
}
