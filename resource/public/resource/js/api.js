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

// ===== API调用函数 =====

/**
 * 执行API请求
 * @param {string} url - API地址
 * @param {string} method - 请求方法
 * @param {Object} [body] - 请求体数据
 * @returns {Promise<Object>} - 响应数据
 */
async function apiRequest(url, method = 'GET', body = null) {
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
}

/**
 * 获取项目列表
 */
async function fetchProjects() {
    try {
        const data = await apiRequest(API_ENDPOINTS.LIST);
        renderProjectList(data.data.list || []);

        // 获取第一个项目的worker信息来更新标题
        if (data.data.list && data.data.list.length > 0) {
            const firstProject = data.data.list[0];
            if (firstProject.worker) {
                updatePageTitle(firstProject.worker);
            }
        }
    } catch (error) {
        showNotification(`获取项目列表失败：${error.message}`, 'danger');
    }
}

/**
 * 注册在线项目
 */
async function registerOnline() {
    const registerButton = document.getElementById('registerButton');

    try {
        // 禁用按钮防止重复点击
        registerButton.disabled = true;
        registerButton.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> 注册中...';

        const result = await apiRequest(API_ENDPOINTS.REGISTER);
        showNotification(result.data.message);
        await fetchProjects();
    } catch (error) {
        showNotification(`注册失败：${error.message}`, 'danger');
    } finally {
        // 恢复按钮状态
        registerButton.disabled = false;
        registerButton.innerHTML = '<i class="bi bi-plus-circle" aria-hidden="true"></i> 注册在线项目';
    }
}

/**
 * 更新项目
 * @param {number} pid - 项目PID
 * @param {string} script - 脚本命令
 * @param {string} description - 项目描述
 */
async function updateProject(pid, script, description) {
    try {
        const result = await apiRequest(
            `${API_ENDPOINTS.UPDATE}${pid}`,
            'POST',
            { script, description }
        );

        showNotification(result.message);
        bootstrap.Modal.getInstance(document.getElementById('editModal')).hide();
        await fetchProjects();
    } catch (error) {
        showNotification(`更新失败：${error.message}`, 'danger');
    }
}

/**
 * 删除项目
 * @param {number} id - 项目ID
 */
async function deleteProject(id) {
    if (!confirm('确定要删除该项目吗？此操作不可恢复。')) {
        return;
    }

    try {
        const result = await apiRequest(`${API_ENDPOINTS.DELETE}${id}`, 'DELETE');
        showNotification(result.message);
        await fetchProjects();
    } catch (error) {
        showNotification(`删除失败：${error.message}`, 'danger');
    }
}

/**
 * 停止项目
 * @param {number} pid - 项目PID
 */
async function stopProject(pid) {
    try {
        const result = await apiRequest(`${API_ENDPOINTS.STOP}${pid}`, 'POST');
        showNotification(result.message);
        await fetchProjects();
    } catch (error) {
        showNotification(`停止失败：${error.message}`, 'danger');
    }
}

/**
 * 停止项目并关闭输出窗口
 * @param {number} pid - 项目PID
 */
async function stopAndClose(pid) {
    if (!confirm('确定要停止运行吗？')) {
        return;
    }

    try {
        await apiRequest(`${API_ENDPOINTS.STOP}${pid}`, 'POST');

        // 等待进程终止
        await new Promise(resolve => setTimeout(resolve, 500));

        // 关闭模态框
        bootstrap.Modal.getInstance(document.getElementById('outputModal')).hide();

        // 刷新项目列表
        await fetchProjects();
    } catch (error) {
        showNotification(`停止进程失败：${error.message}`, 'danger');
    }
}

/**
 * 处理运行请求并显示输出
 * @param {string} url - API地址
 * @param {string} title - 模态框标题
 */
function handleRunRequest(url, title = "运行输出") {
    // 判断是否是直接运行模式（非后台运行）
    const isDirectRun = url.includes('/start/run') && !url.includes('background=true');

    // 显示输出模态框
    const outputModal = new bootstrap.Modal(document.getElementById('outputModal'));
    document.getElementById('outputModalLabel').textContent = title;

    // 清空并初始化输出内容和底部按钮
    const outputContent = document.getElementById('outputContent');
    const outputLoading = document.getElementById('outputLoading');
    const outputModalFooter = document.getElementById('outputModalFooter');

    outputContent.innerHTML = '';
    outputLoading.style.display = 'block';
    outputModalFooter.innerHTML = '';

    outputModal.show();

    // 创建一个变量存储新的PID
    let newPid = null;

    // 创建SSE连接
    const eventSource = new EventSource(url);

    // 处理输出消息
    eventSource.addEventListener('output', (e) => {
        const line = e.data;

        // 处理带颜色的终端输出
        if (line.includes('\x1b[')) {
            const coloredLine = line
                .replace(/\x1b\[1;31m/g, '<span style="color: #ff5555; font-weight: bold;">') // 红色
                .replace(/\x1b\[1;32m/g, '<span style="color: #50fa7b; font-weight: bold;">') // 绿色
                .replace(/\x1b\[1;33m/g, '<span style="color: #f1fa8c; font-weight: bold;">') // 黄色
                .replace(/\x1b\[1;34m/g, '<span style="color: #bd93f9; font-weight: bold;">') // 蓝色
                .replace(/\x1b\[0m/g, '</span>');

            outputContent.innerHTML += coloredLine + '\n';

            // 检查并提取新的PID
            const pidMatch = line.match(/==> 获取到Java进程 PID: (\d+)/);
            if (pidMatch) {
                newPid = parseInt(pidMatch[1]);

                // 如果是直接运行模式，更新关闭按钮的PID
                if (isDirectRun) {
                    updateOutputModalFooter(outputModalFooter, newPid, true);
                }
            }
        } else {
            outputContent.innerHTML += escapeHtml(line) + '\n';
        }

        // 自动滚动到底部
        outputContent.scrollTop = outputContent.scrollHeight;
    });

    // 处理错误
    eventSource.onerror = () => {
        eventSource.close();

        // 如果已经收到了complete事件，就不显示错误
        if (!outputContent.innerHTML.includes('执行完成')) {
            outputContent.innerHTML += '<span style="color: #ff5555; font-weight: bold;">错误：连接已断开</span>\n';
        }

        outputContent.scrollTop = outputContent.scrollHeight;
        outputLoading.style.display = 'none';

        // 显示关闭按钮
        updateOutputModalFooter(outputModalFooter);
    };

    // 处理完成消息
    eventSource.addEventListener('complete', () => {
        eventSource.close();
        outputLoading.style.display = 'none';
        outputContent.scrollTop = outputContent.scrollHeight;

        // 只在非直接运行模式下添加普通关闭按钮
        if (!isDirectRun) {
            updateOutputModalFooter(outputModalFooter);
        }

        // 刷新项目列表
        fetchProjects();
    });
}