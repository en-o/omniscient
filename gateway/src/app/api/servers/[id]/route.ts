// app/api/servers/[id]/route.ts
import { NextResponse } from 'next/server';
import { ServerRepository } from '@utils/database';

// 删除服务器
export async function DELETE(
    request: Request,
    { params }: { params: { id: string } }
) {
    try {
        const { id } = params;

        // 检查服务器是否存在
        const server = await ServerRepository.getById(id);
        if (!server) {
            return NextResponse.json(
                { error: '服务器不存在' },
                { status: 404 }
            );
        }

        await ServerRepository.delete(id);

        return NextResponse.json({ success: true });
    } catch (error) {
        console.error('删除服务器失败:', error);
        return NextResponse.json(
            { error: '删除服务器失败' },
            { status: 500 }
        );
    }
}

// 获取单个服务器
export async function GET(
    request: Request,
    { params }: { params: { id: string } }
) {
    try {
        const { id } = params;
        const server = await ServerRepository.getById(id);

        if (!server) {
            return NextResponse.json(
                { error: '服务器不存在' },
                { status: 404 }
            );
        }

        return NextResponse.json(server);
    } catch (error) {
        console.error('获取服务器详情失败:', error);
        return NextResponse.json(
            { error: '获取服务器详情失败' },
            { status: 500 }
        );
    }
}
