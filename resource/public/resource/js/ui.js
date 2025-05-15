// ===== UI渲染函数 =====

/**
 * 渲染项目列表
 * @param {Array} projects - 项目列表数据
 */
function renderProjectList(projects) {
    const tbody = document.getElementById('projectList');
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
        tr.innerHTML = `
            <td>
                <div class="code-block truncate" data-bs-toggle="tooltip" title="${escapeHtml(project.name)}">${escapeHtml(project.name)}</div>
            </td>
            <td>${escapeHtml(project.ports)}</td>
            <td>${escapeHtml(project.pid)}</td>
            <td>
                <div class="code-block truncate" data-bs-toggle="tooltip" title="${escapeHtml(project.catalog)}">${escapeHtml(project.catalog)}</div>
            </td>
            <td>
                <div class="code-block truncate" data-bs-toggle="tooltip" title="${escapeHtml(project.run)}">${escapeHtml(project.run)}</div>
            </td>
            <td>
                <div class="code-block truncate" data-bs-toggle="tooltip" title="${escapeHtml(project.script || '无')}">${escapeHtml(project.script || '无')}</div>
            </td>
            <td>
                <div class="truncate" data-bs-toggle="tooltip" title="${escapeHtml(project.description || '无')}">${escapeHtml(project.description || '无')}</div>
            </td>
            <td>
                <span class="badge ${project.status === 1 ? 'bg-success' : 'bg-danger'}">
                    ${project.status === 1 ? '运行中' : '已停止'}
                </span>
            </td>
            <td>
                <div class="btn-group">
                    <button class="btn btn-sm btn-outline-secondary dropdown-toggle"
                            type="button"
                            data-bs-toggle="dropdown"
                            data-bs-auto-close="outside"
                            aria-expanded="false">
                        操作
                    </button>
                    <ul class="dropdown-menu dropdown-menu-end">
                        ${project.status === 0 ? `
                            <li><button class="dropdown-item start-run-btn" data-pid="<span class="math-inline">\{project\.pid\}" data\-background\="false"\>
<i class\="bi bi\-play\-fill text\-primary"\></i\> 原生启动
</button\></li\>
<li\><button class\="dropdown\-item start\-run\-btn" data\-pid\="</span>{project.pid}" data-background="true">
                                <i class="bi bi-play-fill text-success"></i> 原生启动(后台)
                            </button></li>
                            <li><button class="dropdown-item start-script-btn" data-pid="<span class="math-inline">\{project\.pid\}"\>
<i class\="bi bi\-play\-circle\-fill text\-success"\></i\> 脚本启动
</button\></li\>
<li\><hr class\="dropdown\-divider"\></li\>
<li\><button class\="dropdown\-item delete\-project\-btn" data\-id\="</span>{project.id}">
                                <i class="bi bi-trash"></i> 删除
                            </button></li>
                        ` : `
                            <li><button class="dropdown-item stop-project-btn" data-pid="${project.pid}">
                                <i class="bi bi-stop-fill text-danger"></i> 停止
                            </button></li>
                        `}
                        <li><hr class="dropdown-divider"></li>
                        <li><button class="dropdown-item edit-project-btn"
                            data-pid="${project.pid}"
                            data-script="${escapeHtml(project.script || '')}"
                            data-description="${escapeHtml(project.description || '')}">
                            <i class="bi bi-pencil text-info"></i> 编辑
                        </button></li>
                    </ul>
                </div>
            </td>
        `;
        tbody.appendChild(tr);
    });

    // 初始化工具提示
    initTooltips();
}

/**
 * 显示编辑模态框
 * @param {number} pid - 项目PID
 * @param {string} script - 脚本命令
 * @param {string} description - 项目描述
 */
function showEditModal(pid, script, description) {
    document.getElementById('editPid').value = pid;
    document.getElementById('editScript').value = script || '';
    document.getElementById('editDescription').value = description || '';
    new bootstrap.Modal(document.getElementById('editModal')).show();
}

/**
 * 更新输出模态框底部按钮
 * @param {HTMLElement} footer - 模态框底部元素
 * @param {number} [pid] - 项目PID（用于直接运行模式）
 * @param {boolean} [isRunning=false] - 是否正在运行
 */
function updateOutputModalFooter(footer, pid = null, isRunning = false) {
    footer.innerHTML = '';

    if (isRunning && pid) {
        // 添加停止并关闭按钮
        const stopButton = document.createElement('button');
        stopButton.type = 'button';
        stopButton.className = 'btn btn-danger';
        stopButton.textContent = '停止并关闭';
        stopButton.addEventListener('click', () => stopAndClose(pid));
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
}