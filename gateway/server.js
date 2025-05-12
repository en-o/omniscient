const express = require('express');
const axios = require('axios');
const bodyParser = require('body-parser');
const cors = require('cors');
const SSE = require('express-sse'); // For Server-Sent Events proxying
const path = require('path');

const app = express();
const port = 3000; // Port for this Node.js backend

// --- In-Memory Server Registry ---
// Structure: { name: "Unique Server Name", baseUrl: "http://java-server-host:port" }
let registeredServers = [
    // Example initial servers (replace or add dynamically)
    { name: "wsl", baseUrl: "http://172.20.3.114:8000" },
];
// ---------------------------------

// --- Middleware ---
app.use(cors()); // Allow requests from frontend (especially if served on different port)
app.use(bodyParser.json()); // Parse JSON bodies
app.use(express.static(path.join(__dirname, 'public'))); // Serve static files (HTML, CSS, JS)
// ------------------

// --- Helper Function ---
function findServerUrl(serverName) {
    const server = registeredServers.find(s => s.name === serverName);
    return server ? server.baseUrl : null;
}
// ---------------------

// --- API Endpoints ---

// GET /api/servers - List registered servers
app.get('/api/servers', (req, res) => {
    res.json(registeredServers.map(s => ({ name: s.name, baseUrl: s.baseUrl }))); // Send only name/baseUrl
});

// POST /api/servers - Register a new server
app.post('/api/servers', (req, res) => {
    const { name, baseUrl } = req.body;
    if (!name || !baseUrl) {
        return res.status(400).json({ error: 'Server name and baseUrl are required.' });
    }
    if (registeredServers.some(s => s.name === name)) {
        return res.status(409).json({ error: `Server with name "${name}" already exists.` });
    }
    // Basic URL validation (optional but recommended)
    try {
        new URL(baseUrl);
    } catch (e) {
        return res.status(400).json({ error: 'Invalid baseUrl format.' });
    }

    const newServer = { name, baseUrl };
    registeredServers.push(newServer);
    console.log(`Registered new server: ${name} (${baseUrl})`);
    res.status(201).json({ message: `Server "${name}" registered successfully.`, server: newServer });
});

// DELETE /api/servers/:name - Unregister a server
app.delete('/api/servers/:name', (req, res) => {
    const serverName = decodeURIComponent(req.params.name); // Handle spaces, etc.
    const initialLength = registeredServers.length;
    registeredServers = registeredServers.filter(s => s.name !== serverName);

    if (registeredServers.length < initialLength) {
        console.log(`Unregistered server: ${serverName}`);
        res.json({ message: `Server "${serverName}" unregistered.` });
    } else {
        res.status(404).json({ error: `Server "${serverName}" not found.` });
    }
});


// --- Proxy Endpoints for Project Management ---

// GET /api/projects/:serverName - Fetch projects from a specific server
app.get('/api/projects/:serverName', async (req, res) => {
    const serverName = decodeURIComponent(req.params.serverName);
    const baseUrl = findServerUrl(serverName);
    if (!baseUrl) {
        return res.status(404).json({ error: `Server "${serverName}" not found.` });
    }

    try {
        const response = await axios.get(`${baseUrl}/jpid`);
        // Add the server name to the response data so frontend knows origin
        const responseData = response.data;
        if (responseData.data && responseData.data.list) {
            responseData.data.list = responseData.data.list.map(p => ({ ...p, serverName: serverName }));
            // Set the worker based on the selected server name
            if (responseData.data.list.length > 0) {
                responseData.data.list[0].worker = serverName;
            }
        } else if (responseData.data) {
            // Handle case where list might be empty or null but data object exists
            responseData.data.worker = serverName;
        }
        res.json(responseData);
    } catch (error) {
        console.error(`Error fetching projects from ${serverName} (${baseUrl}):`, error.message);
        res.status(error.response?.status || 500).json({
            error: `Failed to fetch projects from ${serverName}`,
            details: error.message
        });
    }
});

// GET /api/register/:serverName - Trigger auto-register on a specific server
app.get('/api/register/:serverName', async (req, res) => {
    const serverName = decodeURIComponent(req.params.serverName);
    const baseUrl = findServerUrl(serverName);
    if (!baseUrl) return res.status(404).json({ error: `Server "${serverName}" not found.` });

    try {
        const response = await axios.get(`${baseUrl}/jpid/auto/register`);
        res.status(response.status).json(response.data);
    } catch (error) {
        console.error(`Error triggering auto-register on ${serverName}:`, error.message);
        res.status(error.response?.status || 500).json({ error: `Failed on ${serverName}`, details: error.message });
    }
});

// POST /api/stop/:serverName/:pid - Stop a project
app.post('/api/stop/:serverName/:pid', async (req, res) => {
    const serverName = decodeURIComponent(req.params.serverName);
    const pid = req.params.pid;
    const baseUrl = findServerUrl(serverName);
    if (!baseUrl) return res.status(404).json({ error: `Server "${serverName}" not found.` });

    try {
        const response = await axios.post(`${baseUrl}/jpid/stop/${pid}`);
        res.status(response.status).json(response.data);
    } catch (error) {
        console.error(`Error stopping project ${pid} on ${serverName}:`, error.message);
        res.status(error.response?.status || 500).json({ error: `Failed on ${serverName}`, details: error.message });
    }
});

// POST /api/update/:serverName/:pid - Update a project
app.post('/api/update/:serverName/:pid', async (req, res) => {
    const serverName = decodeURIComponent(req.params.serverName);
    const pid = req.params.pid;
    const baseUrl = findServerUrl(serverName);
    if (!baseUrl) return res.status(404).json({ error: `Server "${serverName}" not found.` });

    try {
        // Forward the request body
        const response = await axios.post(`${baseUrl}/jpid/update/${pid}`, req.body, {
            headers: { 'Content-Type': 'application/json' }
        });
        res.status(response.status).json(response.data);
    } catch (error) {
        console.error(`Error updating project ${pid} on ${serverName}:`, error.message);
        res.status(error.response?.status || 500).json({ error: `Failed on ${serverName}`, details: error.message });
    }
});

// DELETE /api/delete/:serverName/:id - Delete a project
app.delete('/api/delete/:serverName/:id', async (req, res) => {
    const serverName = decodeURIComponent(req.params.serverName);
    const id = req.params.id;
    const baseUrl = findServerUrl(serverName);
    if (!baseUrl) return res.status(404).json({ error: `Server "${serverName}" not found.` });

    try {
        const response = await axios.delete(`${baseUrl}/jpid/delete/${id}`);
        res.status(response.status).json(response.data);
    } catch (error) {
        console.error(`Error deleting project ${id} on ${serverName}:`, error.message);
        res.status(error.response?.status || 500).json({ error: `Failed on ${serverName}`, details: error.message });
    }
});


// GET /api/start/:serverName/:type/:pid - Start a project (Handles SSE)
// :type should be 'run' or 'script'
const sse = new SSE(); // Initialize SSE middleware globally (or per request if needed)

app.get('/api/start/:serverName/:type/:pid', async (req, res) => {
    const serverName = decodeURIComponent(req.params.serverName);
    const type = req.params.type; // 'run' or 'script'
    const pid = req.params.pid;
    const background = req.query.background === 'true'; // Get background query param
    const baseUrl = findServerUrl(serverName);

    if (!baseUrl) {
        res.status(404).json({ error: `Server "${serverName}" not found.` });
        return; // Ensure we don't proceed
    }
    if (type !== 'run' && type !== 'script') {
        res.status(400).json({ error: 'Invalid start type. Must be "run" or "script".' });
        return;
    }

    // Construct the target URL on the Java server
    let targetUrl = `${baseUrl}/jpid/start/${type}/${pid}`;
    if (type === 'run') {
        targetUrl += `?background=${background}`;
    }

    console.log(`Proxying SSE request for ${pid} on ${serverName} (${type}, background: ${background}) to ${targetUrl}`);

    // --- Manual SSE Proxying using Axios stream ---
    try {
        const response = await axios({
            method: 'get',
            url: targetUrl,
            responseType: 'stream'
        });

        res.setHeader('Content-Type', 'text/event-stream');
        res.setHeader('Cache-Control', 'no-cache');
        res.setHeader('Connection', 'keep-alive');
        res.flushHeaders();

        let buffer = '';

        response.data.on('data', (chunk) => {
            buffer += chunk.toString();

            // 处理多行消息
            while (buffer.includes('\n\n')) {
                const messageEnd = buffer.indexOf('\n\n');
                const message = buffer.substring(0, messageEnd);
                buffer = buffer.substring(messageEnd + 2);

                // 解析和格式化消息
                try {
                    const lines = message.split('\n');
                    const eventLine = lines.find(line => line.startsWith('event:'));
                    const dataLine = lines.find(line => line.startsWith('data:'));

                    if (dataLine) {
                        let eventType = 'message';
                        if (eventLine) {
                            eventType = eventLine.substring(6).trim();
                        }

                        let data = dataLine.substring(5).trim();

                        // 尝试解析数据为 JSON，如果失败则作为纯文本发送
                        try {
                            JSON.parse(data); // 验证是否为有效的 JSON
                        } catch (e) {
                            // 如果不是有效的 JSON，将其封装为 JSON 格式
                            data = JSON.stringify({ text: data });
                        }

                        const formattedMessage = `event: ${eventType}\ndata: ${data}\n\n`;
                        res.write(formattedMessage);
                    }
                } catch (parseError) {
                    console.warn('消息解析错误:', parseError);
                    // 发送格式化的错误消息
                    const errorMessage = {
                        error: 'Message parsing error',
                        details: parseError.message,
                        originalMessage: message
                    };
                    res.write(`event: error\ndata: ${JSON.stringify(errorMessage)}\n\n`);
                }
            }
        });

        response.data.on('end', () => {
            // 处理剩余的缓冲区数据
            if (buffer.length > 0) {
                try {
                    const formattedMessage = {
                        text: buffer.trim()
                    };
                    res.write(`event: message\ndata: ${JSON.stringify(formattedMessage)}\n\n`);
                } catch (e) {
                    console.warn('处理最终缓冲区数据时出错:', e);
                }
            }

            // 发送完成事件
            const completeMessage = {
                status: 'complete',
                message: `Stream completed for PID ${pid}`
            };
            res.write(`event: complete\ndata: ${JSON.stringify(completeMessage)}\n\n`);
            res.end();
        });

        response.data.on('error', (err) => {
            console.error(`SSE 流错误 (${serverName}, PID ${pid}):`, err);
            const errorMessage = {
                error: 'Stream error',
                details: err.message,
                pid: pid,
                server: serverName
            };
            res.write(`event: error\ndata: ${JSON.stringify(errorMessage)}\n\n`);
            res.end();
        });

        req.on('close', () => {
            console.log(`客户端断开 SSE 连接 (${serverName}, PID ${pid})`);
            if (response.request) {
                try {
                    response.request.abort();
                } catch (e) {
                    console.warn('中止请求时出错:', e);
                }
            }
        });

    } catch (error) {
        const errorResponse = {
            error: `与服务器 ${serverName} 建立连接失败`,
            details: error.message,
            status: error.response?.status
        };

        if (!res.headersSent) {
            res.status(error.response?.status || 500).json(errorResponse);
        } else {
            res.write(`event: error\ndata: ${JSON.stringify(errorResponse)}\n\n`);
            res.end();
        }
    }
});

// --- Catch-all for serving index.html (for SPA-like behavior if needed) ---
// app.get('*', (req, res) => {
//   res.sendFile(path.join(__dirname, 'public', 'index.html'));
// });
// -------------------------------------------------------------------


// --- Start Server ---
app.listen(port, () => {
    console.log(`Project Manager Node backend listening at http://localhost:${port}`);
    console.log("Registered Servers:", registeredServers);
});
// --------------------
