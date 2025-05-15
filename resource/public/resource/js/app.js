// ===== 全局变量 =====
let contextMenuTarget = null;
let autoRegisterInterval = null;
let tooltips = [];

// ===== 自动注册功能 =====

/**
 * 开始自动注册
 */
function startAutoRegister() {
    // 清除现有定时器
    stopAutoRegister();

    // 立即执行一次
    registerOnline();

    // 设置定时器
    autoRegisterInterval = setInterval(async () => {
        try {
            await registerOnline();
        } catch (error) {
            console.error('自动注册失败:', error);
            stopAutoRegister();
        }
    }, AUTO_REGISTER_INTERVAL);
}

/**
 * 停止自动注册
 */
function stopAutoRegister() {
    if (autoRegisterInterval) {
        clearInterval(autoRegisterInterval);
        autoRegisterInterval = null;
    }
}

// ===== 事件监听器 =====

/**
 * 设置所有事件监听器
 */
function setupEventListeners() {
    // 注册按钮点击事件
    document.getElementById('registerButton').addEventListener('click', registerOnline);

    // 保存编辑按钮点击事件
    document.getElementById('saveEditButton').addEventListener('click', () => {
        const pid = document.getElementById('editPid').value;
        const script = document.getElementById('editScript').value;
        const description = document.getElementById('editDescription').value;

        if (!script.trim() || !description.trim()) {
            showNotification('请填写完整信息', 'warning');
            return;
        }

        updateProject(pid, script, description);
    });

    // 复制上下文菜单按钮点击事件
    document.getElementById('copyContextButton').addEventListener('click', () => {
        if (contextMenuTarget) {
            copyText(contextMenuTarget.textContent);
        }
        hideContextMenu();
    });

    // 项目列表中的操作按钮事件委托
    document.getElementById('projectList').addEventListener('click', (e) => {
        // 启动项目（原生）
        if (e.target.closest('.start-run-btn')) {
            const button = e.target.closest('.start-run-btn');
            const pid = button.getAttribute('data-pid');
            const background = button.getAttribute('data-background') === 'true';
            const title = background ? "原生启动(后台运行)" : "原生启动";
            handleRunRequest(`${API_ENDPOINTS.START_RUN}${pid}?background=${background}`, title);
        }

        // 启动项目（脚本）
        if (e.target.closest('.start-script-btn')) {
            const button = e.target.closest('.start-script-btn');
            const pid = button.getAttribute('data-pid');
            handleRunRequest(`${API_ENDPOINTS.START_SCRIPT}${pid}`, "脚本启动");
        }

        // 停止项目
        if (e.target.closest('.stop-project-btn')) {
            const button = e.target.closest('.stop-project-btn');
            const pid = button.getAttribute('data-pid');
            stopProject(pid);
        }

        // 删除项目
        if (e.target.closest('.delete-project-btn')) {
            const button = e.target.closest('.delete-project-btn');
            const id = button.getAttribute('data-id');
            deleteProject(id);
        }

        // 编辑项目
        if (e.target.closest('.edit-project-btn')) {
            const button = e.target.closest('.edit-project-btn');
            const pid = button.getAttribute('data-pid');
            const script = button.getAttribute('data-script');
            const description = button.getAttribute('data-description');
            showEditModal(pid, script, description);
        }
    });

    // 右键菜单相关事件
    document.addEventListener('contextmenu', handleContextMenu);
    document.addEventListener('click', hideContextMenu);
    document.addEventListener('keydown', handleContextMenuKeyboard);

    // 在页面卸载时清理资源
    window.addEventListener('beforeunload', () => {
        stopAutoRegister();
        destroyAllTooltips();
    });
}

/**
 * 处理右键菜单事件
 * @param {Event} e - 事件对象
 */
function handleContextMenu(e) {
    if (e.target.classList.contains('code-block')) {
        e.preventDefault();
        const menu = document.getElementById('contextMenu');
        contextMenuTarget = e.target;

        // 设置菜单位置，确保不会超出窗口边界
        const x = e.pageX;
        const y = e.pageY;
        const menuWidth = menu.offsetWidth;
        const menuHeight = menu.offsetHeight;
        const windowWidth = window.innerWidth;
        const windowHeight = window.innerHeight;

        menu.style.left = (x + menuWidth > windowWidth ? x - menuWidth : x) + 'px';
        menu.style.top = (y + menuHeight > windowHeight ? y - menuHeight : y) + 'px';
        menu.style.display = 'block';

        // 自动聚焦第一个菜单项
        const firstItem = menu.querySelector('.context-menu-item');
        if (firstItem) {
            firstItem.focus();
        }
    }
}

/**
 * 隐藏右键菜单
 */
function hideContextMenu() {
    document.getElementById('contextMenu').style.display = 'none';
    contextMenuTarget = null;
}

/**
 * 处理右键菜单键盘事件
 * @param {KeyboardEvent} e - 键盘事件对象
 */
function handleContextMenuKeyboard(e) {
    const menu = document.getElementById('contextMenu');
    if (menu.style.display === 'block') {
        if (e.key === 'Escape') {
            hideContextMenu();
        } else if (e.key === 'Enter') {
            if (contextMenuTarget) {
                copyText(contextMenuTarget.textContent);
            }
            hideContextMenu();
        }
    }
}


// ===== 初始化 =====

/**
 * 应用初始化
 */
async function initApp() {
    // 首次加载项目列表
    await fetchProjects();

    // 设置事件监听器
    setupEventListeners();

    // 启动自动注册
    startAutoRegister();

    // 隐藏通知区域
    document.getElementById('notificationArea').classList.remove('show');
}

// 页面加载完成后初始化应用
document.addEventListener('DOMContentLoaded', initApp);