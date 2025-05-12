const express = require('express');
const sqlite3 = require('sqlite3').verbose();
const path = require('path');
const cors = require('cors');

const app = express();
const port = 3000;

// 初始化 SQLite 数据库
const db = new sqlite3.Database('servers.db', (err) => {
    if (err) {
        console.error('数据库连接失败:', err);
    } else {
        console.log('已连接到 SQLite 数据库');
        // 创建服务器表
        db.run(`CREATE TABLE IF NOT EXISTS servers (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE NOT NULL,
            baseUrl TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`);
    }
});

app.use(cors());
app.use(express.json());
app.use(express.static(path.join(__dirname, 'public')));

// 获取所有注册的服务器
app.get('/api/servers', (req, res) => {
    db.all('SELECT name, baseUrl FROM servers', [], (err, rows) => {
        if (err) {
            res.status(500).json({ error: '获取服务器列表失败', details: err.message });
            return;
        }
        res.json(rows);
    });
});


// 注册新服务器
app.post('/api/servers', (req, res) => {
    const { name, baseUrl } = req.body;
    if (!name || !baseUrl) {
        return res.status(400).json({ error: '服务器名称和URL都是必需的' });
    }

    try {
        new URL(baseUrl);
    } catch (e) {
        return res.status(400).json({ error: '无效的URL格式' });
    }

    db.run('INSERT INTO servers (name, baseUrl) VALUES (?, ?)',
        [name, baseUrl],
        function(err) {
            if (err) {
                if (err.message.includes('UNIQUE constraint failed')) {
                    res.status(409).json({ error: `服务器 "${name}" 已存在` });
                } else {
                    res.status(500).json({ error: '注册服务器失败', details: err.message });
                }
                return;
            }
            res.status(201).json({
                message: `服务器 "${name}" 注册成功`,
                server: { name, baseUrl }
            });
        }
    );
});

// 删除服务器
app.delete('/api/servers/:name', (req, res) => {
    const serverName = decodeURIComponent(req.params.name);
    db.run('DELETE FROM servers WHERE name = ?', [serverName], function(err) {
        if (err) {
            res.status(500).json({ error: '删除服务器失败', details: err.message });
            return;
        }
        if (this.changes === 0) {
            res.status(404).json({ error: `服务器 "${serverName}" 未找到` });
            return;
        }
        res.json({ message: `服务器 "${serverName}" 已移除` });
    });
});


// 获取服务器
app.get('/api/servers/:name', (req, res) => {
    const serverName = decodeURIComponent(req.params.name);
    db.get('SELECT name, baseUrl FROM servers WHERE name = ?', [serverName], (err, row) => {
        if (err) {
            console.error('数据库查询错误:', err);
            res.status(500).json({ error: '获取服务器失败', details: err.message });
            return;
        }
        if (!row) {
            res.status(404).json({ error: `服务器 "${serverName}" 未找到` });
            return;
        }
        res.json(row);
    });
});

app.listen(port, () => {
    console.log(`服务注册中心运行在 http://localhost:${port}`);
});

process.on('SIGINT', () => {
    db.close((err) => {
        if (err) {
            console.error('关闭数据库时出错:', err);
        } else {
            console.log('数据库连接已关闭');
        }
        process.exit(0);
    });
});
