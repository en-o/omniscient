// ===== 常量定义 =====
const API_ENDPOINTS = {
    LIST: '/jpid',
    REGISTER: '/jpid/auto/register',
    STOP: '/jpid/stop/',
    START_RUN: '/jpid/start/run/',
    START_SCRIPT: '/jpid/start/script/',
    START_DOCKER: '/jpid/start/docker/',
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
window.apiRequest = async function (url, method = 'GET', body = null) {
    const options = {
        method,
        headers: {'Content-Type': 'application/json'}
    };

    if (body && (method === 'POST' || method === 'PUT')) {
        options.body = JSON.stringify(body);
    }

    const response = await fetch(url, options);
    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || `Request failed (${response.status})`);
    }

    return data;
};

/**
 * 获取项目列表
 */
window.fetchProjects = async function () {
    try {
        const data = await window.apiRequest(window.API_ENDPOINTS.LIST); // Using window.apiRequest and window.API_ENDPOINTS
        const projectList = data.data.list || [];
        window.projectsData = projectList; // Store the fetched data in the global variable

        if (typeof window.renderProjectList === 'function') {
            window.renderProjectList(data.data.list || []); // Using window.renderProjectList
        } else {
            console.error("renderProjectList function not available.");
        }


        // 获取第一个项目的worker信息来更新标题
        if (projectList.length > 0) {
            const firstProject = projectList[0];
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
window.registerOnline = async function () {
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
 * @param { string } catalog - jar目录， 设置也会被更新，用于临时设置后自动更新
 * @param {string} description - 项目描述
 */
window.updateProject = async function (pid, script, catalog, description) {
    try {
        const result = await window.apiRequest( // Using window.apiRequest
            `${window.API_ENDPOINTS.UPDATE}${pid}`, // Using window.API_ENDPOINTS
            'POST',
            {script, catalog, description}
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
window.deleteProject = async function (id) {
    //从存储的数据中查找项目详细信息
    //假设projectsData全局可用
    if (!window.projectsData || window.projectsData.length === 0) {
        console.error("Project data is not loaded. Cannot confirm deletion.");
        if (typeof window.showNotification === 'function') {
            window.showNotification('项目数据未加载，无法执行删除确认。', 'danger');
        }
        return;
    }
    const projectToDelete = window.projectsData.find(project => project.id === id);

    if (!projectToDelete) {
        console.error(`Project with ID ${id} not found in local data.`);
        if (typeof window.showNotification === 'function') {
            window.showNotification('无法找到项目信息进行删除确认。', 'danger');
        }
        //没有项目详细信息无法继续
        return;
    }

    // 填充并显示删除确认模式
    const deleteConfirmModalElement = document.getElementById('deleteConfirmModal');
    const deleteProjectNameSpan = document.getElementById('deleteProjectName');
    const deleteProjectDescriptionSpan = document.getElementById('deleteProjectDescription');
    const deleteProjectIdInput = document.getElementById('deleteProjectId');
    const confirmDeleteButton = document.getElementById('confirmDeleteButton');

    if (!deleteConfirmModalElement || !deleteProjectNameSpan || !deleteProjectDescriptionSpan || !deleteProjectIdInput || !confirmDeleteButton) {
        console.error("Delete confirmation modal elements not found. Check index.html.");
        if (typeof window.showNotification === 'function') {
            window.showNotification('删除确认模态框组件未加载。', 'danger');
        }
        //缺少元素时不能显示模态
        return;
    }

    const deleteConfirmModal = new bootstrap.Modal(deleteConfirmModalElement);


    const escapeHtmlFunc = typeof window.escapeHtml === 'function' ? window.escapeHtml : (str) => str; // Fallback
    deleteProjectNameSpan.textContent = escapeHtmlFunc(projectToDelete.name);
    deleteProjectDescriptionSpan.textContent = escapeHtmlFunc(projectToDelete.description || '无');
    deleteProjectIdInput.value = id;

    //确保每次只有一个点击监听器连接到确认按钮
    //创建一个命名函数，以便于删除
    const confirmDeleteHandler = () => {
        // Hide the modal first
        deleteConfirmModal.hide();

        // Proceed with the actual deletion API call
        executeDelete(id);

        // Remove the event listener after execution or cancellation
        confirmDeleteButton.removeEventListener('click', confirmDeleteHandler);
    };

    // Remove any previous click listeners
    confirmDeleteButton.removeEventListener('click', confirmDeleteHandler); // Defensive removal

    // Add the click listener for the confirm button
    confirmDeleteButton.addEventListener('click', confirmDeleteHandler);

    // Add a listener to clean up when the modal is hidden (e.g., by clicking cancel or backdrop)
    deleteConfirmModalElement.addEventListener('hidden.bs.modal', () => {
        confirmDeleteButton.removeEventListener('click', confirmDeleteHandler);
    }, {once: true}); // Use { once: true } to automatically remove the listener after it's triggered once

    // Show the modal
    deleteConfirmModal.show();
};

/**
 * 执行实际的删除API调用
 * @param {number} id - Project ID
 */
async function executeDelete(id) {
    try {
        const result = await window.apiRequest(`${window.API_ENDPOINTS.DELETE}${id}`, 'DELETE');
        if (typeof window.showNotification === 'function') {
            window.showNotification(result.message);
        }
        if (typeof window.fetchProjects === 'function') {
            await window.fetchProjects(); // Refresh list after successful deletion
        } else {
            console.error("fetchProjects function not available after delete.");
        }
    } catch (error) {
        if (typeof window.showNotification === 'function') {
            window.showNotification(`删除失败：${error.message}`, 'danger');
        } else {
            console.error("showNotification function not available after delete error.");
        }
    }
}

/**
 * 停止项目
 * @param {number} pid - 项目PID
 */
window.stopProject = async function (pid) {
    // 使用 Bootstrap Modal 替代 confirm
    const stopConfirmModalElement = document.getElementById('stopConfirmModal');
    if (!stopConfirmModalElement) {
        // 动态创建确认模态框
        const modalTemplate = `
            <div class="modal fade" id="stopConfirmModal" tabindex="-1" aria-hidden="true">
                <div class="modal-dialog modal-dialog-centered">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title">确认停止</h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                        </div>
                        <div class="modal-body">
                            <p>确定要停止运行吗？</p>
                            <div id="stopProgress" class="progress d-none">
                                <div class="progress-bar progress-bar-striped progress-bar-animated" 
                                     role="progressbar" style="width: 100%"></div>
                            </div>
                        </div>
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">取消</button>
                            <button type="button" class="btn btn-danger" id="confirmStopBtn">停止运行</button>
                        </div>
                    </div>
                </div>
            </div>
        `;
        document.body.insertAdjacentHTML('beforeend', modalTemplate);
    }

    const stopConfirmModal = new bootstrap.Modal(document.getElementById('stopConfirmModal'));
    const confirmStopBtn = document.getElementById('confirmStopBtn');
    const progressBar = document.getElementById('stopProgress');

    // 创建Promise以处理用户确认
    return new Promise((resolve) => {
        const handleConfirm = async () => {
            try {
                // 禁用确认按钮并显示进度条
                confirmStopBtn.disabled = true;
                progressBar.classList.remove('d-none');

                // 执行停止操作
                const result = await window.apiRequest(`${window.API_ENDPOINTS.STOP}${pid}`, 'POST');

                if (typeof window.showNotification === 'function') {
                    window.showNotification(result.message);
                }

                // 刷新项目列表
                if (typeof window.fetchProjects === 'function') {
                    await window.fetchProjects();
                }

                // 关闭模态框
                stopConfirmModal.hide();
                resolve(true);
            } catch (error) {
                if (typeof window.showNotification === 'function') {
                    window.showNotification(`停止失败：${error.message}`, 'danger');
                }
                resolve(false);
            } finally {
                // 重置按钮和进度条状态
                confirmStopBtn.disabled = false;
                progressBar.classList.add('d-none');
            }
        };

        // 绑定确认事件
        confirmStopBtn.onclick = handleConfirm;

        // 显示确认对话框
        stopConfirmModal.show();

        // 模态框关闭时清理事件监听器
        stopConfirmModal._element.addEventListener('hidden.bs.modal', () => {
            confirmStopBtn.onclick = null;
            resolve(false);
        }, { once: true });
    });
};

/**
 * 停止项目并关闭输出窗口
 * @param {number} pid - 项目PID
 */
window.stopAndClose = async function (pid) {
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
window.handleRunRequest = function (url, title = "运行输出") {
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



/**
 * 处理Docker启动请求
 * @param {number} pid - 项目PID
 * @param {boolean} reset - 是否为重启操作
 */
window.handleDockerRequest = function(pid, reset=false) {
    if (!pid) {
        console.error("无效的PID");
        if (typeof window.showNotification === 'function') {
            window.showNotification("无效的项目PID", 'danger');
        }
        return;
    }

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

    document.getElementById('outputModalLabel').textContent = reset ? 'Docker重启输出' : 'Docker启动输出';;

    // 清空并初始化输出内容和底部按钮
    const outputContent = document.getElementById('outputContent');
    const outputLoading = document.getElementById('outputLoading');
    const outputModalFooter = document.getElementById('outputModalFooter');

    if (outputContent) outputContent.innerHTML = '';
    if (outputLoading) outputLoading.style.display = 'block';
    if (outputModalFooter) outputModalFooter.innerHTML = '';

    outputModal.show();


    // Create SSE connection
    const eventSource = new EventSource(`${API_ENDPOINTS.START_DOCKER}/${pid}?reset=${reset}`);

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

        if (typeof window.updateOutputModalFooter === 'function' && outputModalFooter) {
            window.updateOutputModalFooter(outputModalFooter);
        } else {
            console.error("updateOutputModalFooter function or footer not available on error.");
        }
    };

    // Handle complete message
    eventSource.addEventListener('complete', () => {
        eventSource.close();
        if (outputLoading) outputLoading.style.display = 'none';
        if (outputContent) outputContent.scrollTop = outputContent.scrollHeight;

        // 添加这行代码来显示关闭按钮
        if (typeof window.updateOutputModalFooter === 'function' && outputModalFooter) {
            window.updateOutputModalFooter(outputModalFooter);
        } else {
            console.error("updateOutputModalFooter function or footer not available on complete.");
            // 添加一个默认的关闭按钮作为后备方案
            outputModalFooter.innerHTML = `
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">关闭</button>
        `;
        }
        // 刷新项目列表
        if (typeof window.fetchProjects === 'function') {
            window.fetchProjects(); // Using window.fetchProjects
        } else {
            console.error("fetchProjects function not available on complete.");
        }
    });
};


