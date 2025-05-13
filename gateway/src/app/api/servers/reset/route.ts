// app/api/servers/reset/route.ts
import { NextResponse } from 'next/server';

// 重置数据库
export async function POST() {
    try {
        const db = await getDb();

        // 清空服务器表并重新创建
        await db.exec(`
      DROP TABLE IF EXISTS servers;
      CREATE TABLE servers (
        id TEXT PRIMARY KEY,
        url TEXT NOT NULL,
        description TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      );
    `);

        return NextResponse.json({
            success: true,
            message: '数据库已成功重置'
        });
    } catch (error) {
        console.error('重置数据库失败:', error);
        return NextResponse.json(
            { error: '重置数据库失败' },
            { status: 500 }
        );
    }
}
