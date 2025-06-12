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
                <span class="badge ${project.autostart === 1 ? 'bg-success' : 'bg-danger'} autostart-status-btn" 
                      data-id="${project.id}" 
                      data-autostart="${project.autostart}"
                      style="cursor: pointer;"
                      data-bs-toggle="tooltip" 
                      title="点击${project.autostart === 1 ? '卸载' : '注册'}自启">
                    <i class="bi bi-${project.autostart === 1 ? 'check-circle' : 'dash-circle'} me-1"></i>
                    ${project.autostart === 1 ? '自启中' : '待自启'}
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


/**
 * 显示自启确认模态框
 * @param {number} id - 项目ID
 * @param {number} newAutostart - 新的自启状态 (0或1)
 * @param {string} projectName - 项目名称
 * @param {string} projectDescription - 项目描述
 */
window.showAutostartConfirmModal = function(id, newAutostart, projectName, projectDescription) {
    const modal = document.getElementById('autostartConfirmModal');
    const titleElement = document.getElementById('autostartConfirmModalLabel');
    const iconElement = document.getElementById('autostartConfirmIcon');
    const actionElement = document.getElementById('autostartConfirmAction');
    const nameElement = document.getElementById('autostartProjectName');
    const descriptionElement = document.getElementById('autostartProjectDescription');
    const noteElement = document.getElementById('autostartConfirmNote');
    const confirmButton = document.getElementById('confirmAutostartButton');
    const projectIdInput = document.getElementById('autostartProjectId');
    const newValueInput = document.getElementById('autostartNewValue');

    if (!modal) {
        console.error("Autostart confirm modal not found.");
        return;
    }

    // 设置模态框内容
    const isRegister = newAutostart === 1;
    const actionText = isRegister ? '注册' : '卸载';
    const actionColor = isRegister ? 'text-success' : 'text-warning';
    const buttonClass = isRegister ? 'btn-success' : 'btn-warning';
    const iconClass = isRegister ? 'bi-check-circle text-success' : 'bi-dash-circle text-warning';

    if (titleElement) titleElement.textContent = `${actionText}自启确认`;
    if (iconElement) {
        iconElement.className = `bi ${iconClass} fs-1 mb-2`;
    }
    if (actionElement) {
        actionElement.textContent = actionText;
        actionElement.className = actionColor;
    }
    if (nameElement) nameElement.textContent = projectName;
    if (descriptionElement) descriptionElement.textContent = projectDescription || '无描述';
    if (noteElement) {
        const baseMessage = isRegister
            ? '注册自启后，项目将在系统启动时自动运行。'
            : '卸载自启后，项目将不再在系统启动时自动运行。';
        const additionalInfo = isRegister
            ? '<br>(请确保autostart环境正确安装 <a href="https://gitee.com/tanoo/omniscient/blob/master/tools/autostart/README.md#%E5%85%A8%E5%B1%80%E7%8E%AF%E5%A2%83%E8%AE%BE%E7%BD%AE" target="_blank">查看安装指南</a>)'
            : '';
        noteElement.innerHTML = `${baseMessage}${additionalInfo}`;
    }
    if (confirmButton) {
        confirmButton.textContent = `确认${actionText}`;
        confirmButton.className = `btn ${buttonClass}`;
    }
    if (projectIdInput) projectIdInput.value = id;
    if (newValueInput) newValueInput.value = newAutostart;

    // 显示模态框
    if (typeof bootstrap !== 'undefined' && bootstrap.Modal) {
        new bootstrap.Modal(modal).show();
    } else {
        console.error("Bootstrap Modal is not available.");
        if (typeof window.showNotification === 'function') {
            window.showNotification("无法显示确认对话框：依赖组件未加载。", 'danger');
        }
    }
};
