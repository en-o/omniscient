// ===== 常量定义 =====
const API_ENDPOINTS = {
    LIST: '/jpid',
    REGISTER: '/jpid/auto/register',
    STOP: '/jpid/stop/',
    START_RUN: '/jpid/start/run/',
    START_SCRIPT: '/jpid/start/script/',
    DELETE: '/jpid/delete/',
    UPDATE: '/jpid/update/'
};

const AUTO_REGISTER_INTERVAL = 60000; // 60秒

// 将常量暴露到全局作用域，以便从其他文件 (e.g., app.js)
window.API_ENDPOINTS = API_ENDPOINTS;
window.AUTO_REGISTER_INTERVAL = AUTO_REGISTER_INTERVAL;


// ===== API调用函数 =====

/**
 * 执行API请求
 * @param {string} url - API地址
 * @param {string} method - 请求方法
 * @param {Object} [body] - 请求体数据
 * @returns {Promise<Object>} - 响应数据
 */
window.apiRequest = async function(url, method = 'GET', body = null) {
    const options = {
        method,
        headers: { 'Content-Type': 'application/json' }
    };

    if (body && (method === 'POST' || method === 'PUT')) {
        options.body = JSON.stringify(body);
    }

    const response = await fetch(url, options);
    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || `请求失败 (${response.status})`);
    }

    return data;
};

/**
 * 获取项目列表
 */
window.fetchProjects = async function() {
    try {
        const data = await window.apiRequest(window.API_ENDPOINTS.LIST); // Using window.apiRequest and window.API_ENDPOINTS
        if (typeof window.renderProjectList === 'function') {
            window.renderProjectList(data.data.list || []); // Using window.renderProjectList
        } else {
            console.error("renderProjectList function not available.");
        }


        // 获取第一个项目的worker信息来更新标题
        if (data.data.list && data.data.list.length > 0) {
            const firstProject = data.data.list[0];
            if (firstProject.worker) {
                if (typeof window.updatePageTitle === 'function') {
                    window.updatePageTitle(firstProject.worker); // Using window.updatePageTitle
                } else {
                    console.error("updatePageTitle function not available.");
                }
            }
        }
    } catch (error) {
        if (typeof window.showNotification === 'function') {
            window.showNotification(`获取项目列表失败：${error.message}`, 'danger'); // Using window.showNotification
        } else {
            console.error("showNotification function not available.");
        }
    }
};

/**
 * 注册在线项目
 */
window.registerOnline = async function() {
    const registerButton = document.getElementById('registerButton');
    if (!registerButton) {
        console.error("Register button not found.");
        return;
    }


    try {
        // 禁用按钮防止重复点击
        registerButton.disabled = true;
        registerButton.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> 注册中...';

        const result = await window.apiRequest(window.API_ENDPOINTS.REGISTER); // Using window.apiRequest and window.API_ENDPOINTS
        if (typeof window.showNotification === 'function') {
            window.showNotification(result.data.message); // Using window.showNotification
        } else {
            console.error("showNotification function not available.");
        }

        if (typeof window.fetchProjects === 'function') {
            await window.fetchProjects(); // Using window.fetchProjects
        } else {
            console.error("fetchProjects function not available.");
        }
    } catch (error) {
        if (typeof window.showNotification === 'function') {
            window.showNotification(`注册失败：${error.message}`, 'danger'); // Using window.showNotification
        } else {
            console.error("showNotification function not available.");
        }
    } finally {
        // 恢复按钮状态
        registerButton.disabled = false;
        registerButton.innerHTML = '<i class="bi bi-plus-circle" aria-hidden="true"></i> 注册在线项目';
    }
};

/**
 * 更新项目
 * @param {number} pid - 项目PID
 * @param {string} script - 脚本命令
 * @param {string} description - 项目描述
 */
window.updateProject = async function(pid, script, description) {
    try {
        const result = await window.apiRequest( // Using window.apiRequest
            `${window.API_ENDPOINTS.UPDATE}${pid}`, // Using window.API_ENDPOINTS
            'POST',
            { script, description }
        );

        if (typeof window.showNotification === 'function') {
            window.showNotification(result.message); // Using window.showNotification
        } else {
            console.error("showNotification function not available.");
        }

        // Ensure bootstrap is available globally
        const editModal = document.getElementById('editModal');
        if (typeof bootstrap !== 'undefined' && bootstrap.Modal && editModal) {
            bootstrap.Modal.getInstance(editModal).hide();
        } else {
            console.error("Bootstrap Modal or edit modal element is not available.");
        }
        if (typeof window.fetchProjects === 'function') {
            await window.fetchProjects(); // Using window.fetchProjects
        } else {
            console.error("fetchProjects function not available.");
        }
    } catch (error) {
        if (typeof window.showNotification === 'function') {
            window.showNotification(`更新失败：${error.message}`, 'danger'); // Using window.showNotification
        } else {
            console.error("showNotification function not available.");
        }
    }
};

/**
 * 删除项目
 * @param {number} id - 项目ID
 */
window.deleteProject = async function(id) {
    if (!confirm('确定要删除该项目吗？此操作不可恢复。')) {
        return;
    }

    try {
        const result = await window.apiRequest(`${window.API_ENDPOINTS.DELETE}${id}`, 'DELETE'); // Using window.apiRequest and window.API_ENDPOINTS
        if (typeof window.showNotification === 'function') {
            window.showNotification(result.message); // Using window.showNotification
        } else {
            console.error("showNotification function not available.");
        }
        if (typeof window.fetchProjects === 'function') {
            await window.fetchProjects(); // Using window.fetchProjects
        } else {
            console.error("fetchProjects function not available.");
        }
    } catch (error) {
        if (typeof window.showNotification === 'function') {
            window.showNotification(`删除失败：${error.message}`, 'danger'); // Using window.showNotification
        } else {
            console.error("showNotification function not available.");
        }
    }
};

/**
 * 停止项目
 * @param {number} pid - 项目PID
 */
window.stopProject = async function(pid) {
    try {
        const result = await window.apiRequest(`${window.API_ENDPOINTS.STOP}${pid}`, 'POST'); // Using window.apiRequest and window.API_ENDPOINTS
        if (typeof window.showNotification === 'function') {
            window.showNotification(result.message); // Using window.showNotification
        } else {
            console.error("showNotification function not available.");
        }
        if (typeof window.fetchProjects === 'function') {
            await window.fetchProjects(); // Using window.fetchProjects
        } else {
            console.error("fetchProjects function not available.");
        }
    } catch (error) {
        if (typeof window.showNotification === 'function') {
            window.showNotification(`停止失败：${error.message}`, 'danger'); // Using window.showNotification
        } else {
            console.error("showNotification function not available.");
        }
    }
};

/**
 * 停止项目并关闭输出窗口
 * @param {number} pid - 项目PID
 */
window.stopAndClose = async function(pid) {
    if (!confirm('确定要停止运行吗？')) {
        return;
    }

    try {
        await window.apiRequest(`${window.API_ENDPOINTS.STOP}${pid}`, 'POST'); // Using window.apiRequest and window.API_ENDPOINTS

        // 等待进程终止
        await new Promise(resolve => setTimeout(resolve, 500));

        // 关闭模态框
        // Ensure bootstrap is available globally
        const outputModal = document.getElementById('outputModal');
        if (typeof bootstrap !== 'undefined' && bootstrap.Modal && outputModal) {
            bootstrap.Modal.getInstance(outputModal).hide();
        } else {
            console.error("Bootstrap Modal or output modal element is not available.");
        }

        // 刷新项目列表
        if (typeof window.fetchProjects === 'function') {
            await window.fetchProjects(); // Using window.fetchProjects
        } else {
            console.error("fetchProjects function not available.");
        }
    } catch (error) {
        if (typeof window.showNotification === 'function') {
            window.showNotification(`停止进程失败：${error.message}`, 'danger'); // Using window.showNotification
        } else {
            console.error("showNotification function not available.");
        }
    }
};

/**
 * 处理运行请求并显示输出
 * @param {string} url - API地址
 * @param {string} title - 模态框标题
 */
window.handleRunRequest = function(url, title = "运行输出") {
    // 判断是否是直接运行模式（非后台运行）
    const isDirectRun = url.includes('/start/run') && !url.includes('background=true');

    // 显示输出模态框
    // Ensure bootstrap is available globally
    const outputModalElement = document.getElementById('outputModal');
    if (typeof bootstrap === 'undefined' || !bootstrap.Modal || !outputModalElement) {
        console.error("Bootstrap Modal or output modal element is not available. Cannot show output.");
        if (typeof window.showNotification === 'function') {
            window.showNotification("无法显示运行输出：依赖组件未加载。", 'danger');
        }
        return;
    }
    const outputModal = new bootstrap.Modal(outputModalElement);
    document.getElementById('outputModalLabel').textContent = title;

    // 清空并初始化输出内容和底部按钮
    const outputContent = document.getElementById('outputContent');
    const outputLoading = document.getElementById('outputLoading');
    const outputModalFooter = document.getElementById('outputModalFooter');

    if (outputContent) outputContent.innerHTML = '';
    if (outputLoading) outputLoading.style.display = 'block';
    if (outputModalFooter) outputModalFooter.innerHTML = '';

    outputModal.show();

    // 创建一个变量存储新的 PID
    let newPid = null;

    // Create SSE connection
    const eventSource = new EventSource(url);

    // Handle output messages
    eventSource.addEventListener('output', (e) => {
        const line = e.data;

        if (outputContent) {
            // Handle colored terminal output
            if (line.includes('\x1b[')) {
                const coloredLine = line
                    .replace(/\x1b\[1;31m/g, '<span style="color: #ff5555; font-weight: bold;">') // Red
                    .replace(/\x1b\[1;32m/g, '<span style="color: #50fa7b; font-weight: bold;">') // Green
                    .replace(/\x1b\[1;33m/g, '<span style="color: #f1fa8c; font-weight: bold;">') // Yellow
                    .replace(/\x1b\[1;34m/g, '<span style="color: #bd93f9; font-weight: bold;">') // Blue
                    .replace(/\x1b\[0m/g, '</span>');

                outputContent.innerHTML += coloredLine + '\n';

                // Check and extract new PID
                const pidMatch = line.match(/==> 获取到Java进程 PID: (\d+)/);
                if (pidMatch) {
                    newPid = parseInt(pidMatch[1]);

                    // If in direct run mode, update the close button PID
                    if (isDirectRun) {
                        // Using window.updateOutputModalFooter
                        if (typeof window.updateOutputModalFooter === 'function' && outputModalFooter) {
                            window.updateOutputModalFooter(outputModalFooter, newPid, true);
                        } else {
                            console.error("updateOutputModalFooter function or footer not available.");
                        }
                    }
                }
            } else {
                // Using window.escapeHtml
                if (typeof window.escapeHtml === 'function') {
                    outputContent.innerHTML += window.escapeHtml(line) + '\n';
                } else {
                    console.error("escapeHtml function not available.");
                    outputContent.innerHTML += line + '\n'; // Fallback if escapeHtml is missing
                }
            }

            // Auto-scroll to bottom
            outputContent.scrollTop = outputContent.scrollHeight;
        } else {
            console.error("outputContent element not found.");
        }
    });

    // Handle errors
    eventSource.onerror = () => {
        eventSource.close();

        if (outputContent) {
            // If a complete event has already been received, don't show error
            if (!outputContent.innerHTML.includes('执行完成')) {
                outputContent.innerHTML += '<span style="color: #ff5555; font-weight: bold;">错误：连接已断开</span>\n';
            }

            outputContent.scrollTop = outputContent.scrollHeight;
        } else {
            console.error("outputContent element not found on error.");
        }


        if (outputLoading) outputLoading.style.display = 'none';

        // Display close button
        // Using window.updateOutputModalFooter
        if (typeof window.updateOutputModalFooter === 'function' && outputModalFooter) {
            window.updateOutputModalFooter(outputModalFooter);
        } else {
            console.error("updateOutputModalFooter function or footer not available on error.");
            // Optionally add a default close button here
        }
    };

    // Handle complete message
    eventSource.addEventListener('complete', () => {
        eventSource.close();
        if (outputLoading) outputLoading.style.display = 'none';
        if (outputContent) outputContent.scrollTop = outputContent.scrollHeight;

        // Only add a regular close button in non-direct run mode
        if (!isDirectRun) {
            // Using window.updateOutputModalFooter
            if (typeof window.updateOutputModalFooter === 'function' && outputModalFooter) {
                window.updateOutputModalFooter(outputModalFooter);
            } else {
                console.error("updateOutputModalFooter function or footer not available on complete.");
                // Optionally add a default close button here
            }
        }

        // Refresh project list
        if (typeof window.fetchProjects === 'function') {
            window.fetchProjects(); // Using window.fetchProjects
        } else {
            console.error("fetchProjects function not available on complete.");
        }
    });
};