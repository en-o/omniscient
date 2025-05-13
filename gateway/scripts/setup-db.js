// scripts/setup-db.js
const fs = require('fs');
const path = require('path');
const sqlite3 = require('sqlite3').verbose();

// 确保数据目录存在
const dataDir = path.join(__dirname, '..', 'data');
if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir, { recursive: true });
}

// 数据库文件路径
const dbPath = path.join(dataDir, 'servers.db');

// 创建/打开数据库连接
const db = new sqlite3.Database(dbPath, (err) => {
    if (err) {
        console.error('创建数据库连接失败:', err.message);
        process.exit(1);
    }
    console.log('已连接到 SQLite 数据库');
});

// 创建服务器表
db.serialize(() => {
    db.run(`
    CREATE TABLE IF NOT EXISTS servers (
      id TEXT PRIMARY KEY,
      url TEXT NOT NULL,
      description TEXT NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
  `, (err) => {
        if (err) {
            console.error('创建服务器表失败:', err.message);
        } else {
            console.log('服务器表创建成功或已存在');
        }
    });
});

// 关闭数据库连接
db.close((err) => {
    if (err) {
        console.error('关闭数据库连接失败:', err.message);
    } else {
        console.log('数据库连接已关闭');
        console.log('数据库设置完成!');
    }
});
