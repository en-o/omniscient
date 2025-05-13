import { NextResponse } from 'next/server';
import { ServerRepository } from '@utils/database';



// 导出服务器数据
export async function GET() {
    try {
        // 获取所有服务器数据
        const servers = await ServerRepository.getAll();

        // 设置响应头，让浏览器下载文件
        const headers = new Headers();
        headers.append('Content-Disposition', 'attachment; filename=servers_backup.json');
        headers.append('Content-Type', 'application/json');

        // 返回 JSON 数据
        return new NextResponse(JSON.stringify(servers, null, 2), {
            status: 200,
            headers
        });
    } catch (error) {
        console.error('导出服务器数据失败:', error);
        return NextResponse.json(
            { error: '导出服务器数据失败' },
            { status: 500 }
        );
    }
}
