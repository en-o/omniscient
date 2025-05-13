import sqlite3 from 'sqlite3'
import { open, Database } from 'sqlite'
import path from 'path'

// 全局数据库连接实例
let db: Database | null = null;

/**
 * 初始化数据库连接
 */
export async function getDb(): Promise<Database> {
    if (db) return db;

    // 确保在服务器端运行
    if (typeof window !== 'undefined') {
        throw new Error('数据库操作只能在服务器端执行');
    }

    const dbPath = path.join(process.cwd(), 'data', 'servers.db');

    db = await open({
        filename: dbPath,
        driver: sqlite3.Database
    });

    // 创建服务器表（如果不存在）
    await db.exec(`
    CREATE TABLE IF NOT EXISTS servers (
      id TEXT PRIMARY KEY,
      url TEXT NOT NULL,
      description TEXT NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
  `);

    return db;
}

// 服务器实体模型
export interface ServerDbEntity {
    id: string;
    url: string;
    description: string;
    created_at?: string;
}

// 服务器数据存储库
export class ServerRepository {
    // 获取所有服务器
    static async getAll(): Promise<ServerDbEntity[]> {
        const db = await getDb();
        return db.all<ServerDbEntity[]>('SELECT * FROM servers ORDER BY created_at DESC');
    }

    // 添加新服务器
    static async add(server: ServerDbEntity): Promise<void> {
        const db = await getDb();
        await db.run(
            'INSERT INTO servers (id, url, description) VALUES (?, ?, ?)',
            [server.id, server.url, server.description]
        );
    }

    // 删除服务器
    static async delete(id: string): Promise<void> {
        const db = await getDb();
        await db.run('DELETE FROM servers WHERE id = ?', id);
    }

    // 根据ID获取服务器
    static async getById(id: string): Promise<ServerDbEntity | undefined> {
        const db = await getDb();
        return db.get<ServerDbEntity>('SELECT * FROM servers WHERE id = ?', id);
    }
}
