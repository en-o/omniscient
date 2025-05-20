// ===== 全局变量 (for tooltips, managed within UI logic) =====
let tooltips = [];

// ===== UI渲染函数 =====

/**
 * 销毁所有工具提示实例
 */
window.destroyAllTooltips = function() {
    if (typeof bootstrap !== 'undefined' && bootstrap.Tooltip) {
        tooltips.forEach(tooltip => {
            tooltip.dispose();
        });
        tooltips = [];
    }
};

/**
 * 初始化工具提示
 */
window.initTooltips = function() {
    if (typeof bootstrap !== 'undefined' && bootstrap.Tooltip) {
        // First clear old tooltips
        if (typeof window.destroyAllTooltips === 'function') {
            window.destroyAllTooltips();
        } else {
            console.error("destroyAllTooltips function not available.");
        }


        // Create new tooltips
        const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
        tooltips = tooltipTriggerList.map(tooltipTriggerEl => {
            return new bootstrap.Tooltip(tooltipTriggerEl);
        });
    } else {
        console.error("Bootstrap Tooltip is not available.");
    }
};

/**
 * 渲染项目列表
 * @param {Array} projects - 项目列表数据
 */
window.renderProjectList = function(projects) {
    const tbody = document.getElementById('projectList');
    if (!tbody) {
        console.error("Project list tbody element not found.");
        return;
    }
    tbody.innerHTML = '';

    if (!projects || projects.length === 0) {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td colspan="9" class="text-center text-muted py-4">
                <div>
                    <i class="bi bi-inbox fs-4"></i>
                    <p class="mt-2">暂无项目数据</p>
                </div>
            </td>
        `;
        tbody.appendChild(tr);
        return;
    }

    projects.forEach(project => {
        const tr = document.createElement('tr');
        // Using window.escapeHtml to explicitly call the global function
        const escapeHtmlFunc = typeof window.escapeHtml === 'function' ? window.escapeHtml : (str) => str; // Fallback

        // 根据项目状态和运行方式准备操作菜单项
        let operationItems = '';

        if (project.status === 0) { // 已停止状态
            if (project.way === 1) { // Docker方式
                operationItems = `
                    <li><button class="dropdown-item docker-start-btn" data-pid="${project.pid}" data-reset="false">
                        <i class="bi bi-play-fill text-primary"></i> Docker启动
                    </button></li>
                    <li><hr class="dropdown-divider"></li>
                    <li><button class="dropdown-item delete-project-btn text-danger" data-id="${project.id}">
                        <i class="bi bi-trash"></i> 删除
                    </button></li>
                `;
            } else { // 原生方式
                operationItems = `
                    <li><button class="dropdown-item start-run-btn" data-pid="${project.pid}" data-background="false">
                        <i class="bi bi-play-fill text-primary"></i> 原生启动
                    </button></li>
                    <li><button class="dropdown-item start-run-btn" data-pid="${project.pid}" data-background="true">
                        <i class="bi bi-play-fill text-success"></i> 原生启动(后台)
                    </button></li>
                    <li><button class="dropdown-item start-script-btn" data-pid="${project.pid}">
                        <i class="bi bi-play-circle-fill text-success"></i> 脚本启动
                    </button></li>
                    <li><hr class="dropdown-divider"></li>
                    <li><button class="dropdown-item delete-project-btn text-danger" data-id="${project.id}">
                        <i class="bi bi-trash"></i> 删除
                    </button></li>
                `;
            }
        } else {
            // 运行中状态
            operationItems = `
                <li><button class="dropdown-item stop-project-btn" data-pid="${project.pid}">
                    <i class="bi bi-stop-fill text-danger"></i> 停止
                </button></li>
            `;
            if(project.way === 1){
                operationItems += `
                 <li><button class="dropdown-item docker-start-btn" data-pid="${project.pid}" data-reset="true">
                        <i class="bi bi-arrow-clockwise text-success"></i> Docker重启
                    </button></li>
            `;
            }

        }

        // 添加编辑选项（适用于所有项目类型）
        operationItems += `
            <li><hr class="dropdown-divider"></li>
            <li><button class="dropdown-item edit-project-btn"
                data-pid="${escapeHtmlFunc(project.pid || '')}"
                data-script="${escapeHtmlFunc(project.script || '')}"
                data-catalog="${escapeHtmlFunc(project.catalog || '')}"
                data-description="${escapeHtmlFunc(project.description || '')}">
                <i class="bi bi-pencil text-info"></i> 编辑
            </button></li>
        `;

        tr.innerHTML = `
            <td>
                <div class="code-block truncate" data-bs-toggle="tooltip" 
                     title="${project.way === 1 ?
                        '容器名: ' + escapeHtmlFunc(project.name) :
                        'JAR包: ' + escapeHtmlFunc(project.name)}"
                >${escapeHtmlFunc(project.name)}</div>
            </td>
            <td>${escapeHtmlFunc(project.ports)}</td>
            <td>${escapeHtmlFunc(project.pid)}</td>
            <td>
                <div class="code-block truncate" data-bs-toggle="tooltip" title="${escapeHtmlFunc(project.catalog)}">${escapeHtmlFunc(project.catalog)}</div>
            </td>
            <td>
                <div class="code-block truncate" data-bs-toggle="tooltip" title="${escapeHtmlFunc(project.run)}">${escapeHtmlFunc(project.run)}</div>
            </td>
            <td>
                <div class="code-block truncate" data-bs-toggle="tooltip" title="${escapeHtmlFunc(project.script || '无')}">${escapeHtmlFunc(project.script || '无')}</div>
            </td>
            <td>
                <div class="truncate" data-bs-toggle="tooltip" title="${escapeHtmlFunc(project.description || '无')}">${escapeHtmlFunc(project.description || '无')}</div>
            </td>
            <td>
                <span class="badge ${project.status === 1 ? 'bg-success' : 'bg-danger'}">
                    ${project.status === 1 ? '运行中' : '已停止'}
                </span>
            </td>
            <td>
                <span class="badge ${project.way === 1 ? 'bg-primary' : 'bg-success'}">
                    ${project.way === 1 ? 'docker' : 'jdk'}
                </span>
            </td>
            <td>
                <div class="btn-group"> <button class="btn btn-sm btn-outline-secondary dropdown-toggle"
                            type="button"
                            data-bs-toggle="dropdown"
                            data-bs-auto-close="outside"
                            aria-expanded="false">
                        操作
                    </button>
                    <ul class="dropdown-menu dropdown-menu-end">
                        ${operationItems}
                    </ul>
                </div>
            </td>
        `;
        tbody.appendChild(tr);
    });

    // Initialize tooltips after rendering the list
    if (typeof window.initTooltips === 'function') {
        window.initTooltips();
    } else {
        console.error("initTooltips function not available.");
    }
};

/**
 * 显示编辑模态框
 * @param {number} pid - 项目PID
 * @param {string} script - 脚本命令
 * @param {string} description - 项目描述
 */
window.showEditModal = function(pid, script, catalog, description) {
    const editPidInput = document.getElementById('editPid');
    const editScriptTextarea = document.getElementById('editScript');
    const editCatalogTextarea = document.getElementById('editCatalog');
    const editDescriptionTextarea = document.getElementById('editDescription');
    const editModalElement = document.getElementById('editModal');

    if (editPidInput) editPidInput.value = pid;
    if (editScriptTextarea) editScriptTextarea.value = script || '';
    if (editCatalogTextarea) editCatalogTextarea.value = catalog || '';
    if (editDescriptionTextarea) editDescriptionTextarea.value = description || '';

    // Ensure bootstrap is available globally
    if (typeof bootstrap !== 'undefined' && bootstrap.Modal && editModalElement) {
        new bootstrap.Modal(editModalElement).show();
    } else {
        console.error("Bootstrap Modal or edit modal element is not available.");
        if (typeof window.showNotification === 'function') {
            window.showNotification("无法显示编辑模态框：依赖组件未加载。", 'danger');
        }
    }
};

/**
 * 更新输出模态框底部按钮
 * @param {HTMLElement} footer - 模态框底部元素
 * @param {number} [pid] - 项目PID（用于直接运行模式）
 * @param {boolean} [isRunning=false] - 是否正在运行
 */
window.updateOutputModalFooter = function(footer, pid = null, isRunning = false) {
    if (!footer) {
        console.error("Output modal footer element not provided.");
        return;
    }
    footer.innerHTML = '';

    if (isRunning && pid) {
        // 添加停止并关闭按钮
        const stopButton = document.createElement('button');
        stopButton.type = 'button';
        stopButton.className = 'btn btn-danger';
        stopButton.textContent = '停止并关闭';
        // Make sure stopAndClose is accessible globally
        if (typeof window.stopAndClose === 'function') {
            stopButton.addEventListener('click', () => window.stopAndClose(pid));
        } else {
            console.error("stopAndClose function is not available for stop button.");
            // Add a disabled button or fallback
            stopButton.disabled = true;
            stopButton.textContent = '停止并关闭 (功能不可用)';
        }
        footer.appendChild(stopButton);
    } else {
        // 添加普通关闭按钮
        const closeButton = document.createElement('button');
        closeButton.type = 'button';
        closeButton.className = 'btn btn-secondary';
        closeButton.textContent = '关闭';
        closeButton.setAttribute('data-bs-dismiss', 'modal');
        footer.appendChild(closeButton);
    }
};
