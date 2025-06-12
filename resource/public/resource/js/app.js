// ===== 全局变量 =====
let contextMenuTarget = null;
let autoRegisterInterval = null;


// 添加一个全局变量来存储项目列表
window.projectsData = [];


// ===== 自动注册功能 =====

/**
 * 开始自动注册
 */
window.startAutoRegister = function () {
    // 清除现有定时器
    if (typeof window.stopAutoRegister === 'function') {
        window.stopAutoRegister();
    } else {
        console.error("stopAutoRegister function not available.");
    }


    // 立即执行一次
    if (typeof window.registerOnline === 'function') {
        window.registerOnline();
    } else {
        console.error("registerOnline function not available.");
    }


    // 设置定时器
    // Access AUTO_REGISTER_INTERVAL globally
    if (typeof AUTO_REGISTER_INTERVAL !== 'undefined') {
        autoRegisterInterval = setInterval(async () => {
            try {
                if (typeof window.registerOnline === 'function') {
                    await window.registerOnline();
                } else {
                    console.error("registerOnline function not available during interval.");
                }
            } catch (error) {
                console.error('Automatic registration failed:', error);
                if (typeof window.stopAutoRegister === 'function') {
                    window.stopAutoRegister();
                } else {
                    console.error("stopAutoRegister function not available after interval error.");
                }
            }
        }, AUTO_REGISTER_INTERVAL);
    } else {
        console.error("AUTO_REGISTER_INTERVAL is not defined globally. Auto-registration will not start.");
    }
};


/**
 * 停止自动注册
 */
window.stopAutoRegister = function () {
    if (autoRegisterInterval) {
        clearInterval(autoRegisterInterval);
        autoRegisterInterval = null;
    }
};

// ===== 事件监听器 =====

/**
 * 设置所有事件监听器
 */
window.setupEventListeners = function () {
    // 注册按钮点击事件
    const registerButton = document.getElementById('registerButton');
    if (registerButton && typeof window.registerOnline === 'function') {
        registerButton.addEventListener('click', window.registerOnline);
    } else {
        console.error("Register button or registerOnline function not available.");
    }


    // 保存编辑按钮点击事件
    const saveEditButton = document.getElementById('saveEditButton');
    if (saveEditButton) {
        saveEditButton.addEventListener('click', () => {
            const pid = document.getElementById('editPid').value;
            const script = document.getElementById('editScript').value;
            const catalog = document.getElementById('editCatalog').value;
            const description = document.getElementById('editDescription').value;

            if (!script.trim() || !description.trim()) {
                if (typeof window.showNotification === 'function') {
                    window.showNotification('请填写完整信息', 'warning');
                } else {
                    console.error("showNotification function not available.");
                }
                return;
            }

            if (typeof window.updateProject === 'function') {
                window.updateProject(pid, script, catalog, description);
            } else {
                console.error("updateProject function not available.");
            }
        });
    } else {
        console.error("Save edit button not found.");
    }


    // 复制上下文菜单按钮点击事件
    const copyContextButton = document.getElementById('copyContextButton');
    if (copyContextButton) {
        copyContextButton.addEventListener('click', () => {
            if (contextMenuTarget) {
                if (typeof window.copyText === 'function') {
                    window.copyText(contextMenuTarget.textContent);
                } else {
                    console.error("copyText function not available.");
                }
            }
            if (typeof window.hideContextMenu === 'function') {
                window.hideContextMenu();
            } else {
                console.error("hideContextMenu function not available.");
            }
        });
    } else {
        console.error("Copy context menu button not found.");
    }


    // 项目列表中的操作按钮事件委托
    const projectList = document.getElementById('projectList');
    if (projectList) {
        projectList.addEventListener('click', (e) => {
            // 启动项目（原生）
            if (e.target.closest('.start-run-btn')) {
                const button = e.target.closest('.start-run-btn');
                const pid = button.getAttribute('data-pid');
                const background = button.getAttribute('data-background') === 'true';
                const title = background ? "原生启动(后台运行)" : "原生启动";
                // 启动项目（脚本）
                if (typeof window.handleRunRequest === 'function' && typeof API_ENDPOINTS !== 'undefined') {
                    window.handleRunRequest(`${API_ENDPOINTS.START_RUN}${pid}?background=${background}`, title);
                } else {
                    console.error("handleRunRequest or API_ENDPOINTS not available.");
                }
            }

            // 启动项目（脚本）
            if (e.target.closest('.start-script-btn')) {
                const button = e.target.closest('.start-script-btn');
                const pid = button.getAttribute('data-pid');
                // Access API_ENDPOINTS globally
                if (typeof window.handleRunRequest === 'function' && typeof API_ENDPOINTS !== 'undefined') {
                    window.handleRunRequest(`${API_ENDPOINTS.START_SCRIPT}${pid}`, "脚本启动");
                } else {
                    console.error("handleRunRequest or API_ENDPOINTS not available.");
                }
            }

            // Docker 启动项目
            if (e.target.closest('.docker-start-btn')) {
                const button = e.target.closest('.docker-start-btn');
                const pid = button.getAttribute('data-pid');
                const reset = button.getAttribute('data-reset') === 'true';
                if (typeof window.handleDockerRequest === 'function') {
                    window.handleDockerRequest(parseInt(pid), reset);
                } else {
                    console.error("handleDockerRequest function not available.");
                }
            }

            // 停止项目
            if (e.target.closest('.stop-project-btn')) {
                const button = e.target.closest('.stop-project-btn');
                const pid = button.getAttribute('data-pid');
                if (typeof window.stopProject === 'function') {
                    window.stopProject(pid);
                } else {
                    console.error("stopProject function not available.");
                }
            }

            // 删除项目
            if (e.target.closest('.delete-project-btn')) {
                const button = e.target.closest('.delete-project-btn');
                const id = button.getAttribute('data-id');
                if (typeof window.deleteProject === 'function') {
                    window.deleteProject(parseInt(id)).then(r => {}); // Ensure ID is a number
                } else {
                    console.error("deleteProject function not available.");
                }
            }

            // 编辑项目
            if (e.target.closest('.edit-project-btn')) {
                const button = e.target.closest('.edit-project-btn');
                const pid = button.getAttribute('data-pid');
                const script = button.getAttribute('data-script');
                const catalog = button.getAttribute('data-catalog');
                const description = button.getAttribute('data-description');
                if (typeof window.showEditModal === 'function') {
                    window.showEditModal(pid, script, catalog, description);
                } else {
                    console.error("showEditModal function not available.");
                }
            }

            // 在项目列表事件委托中替换自启状态点击处理
            if (e.target.closest('.autostart-status-btn')) {
                const button = e.target.closest('.autostart-status-btn');
                const id = parseInt(button.getAttribute('data-id'));
                const currentAutostart = parseInt(button.getAttribute('data-autostart'));
                const newAutostart = currentAutostart === 1 ? 0 : 1;

                // 从项目数据中获取项目信息
                const project = window.projectsData.find(p => p.id === id);
                if (!project) {
                    if (typeof window.showNotification === 'function') {
                        window.showNotification('项目信息未找到', 'danger');
                    }
                    return;
                }

                if (typeof window.showAutostartConfirmModal === 'function') {
                    window.showAutostartConfirmModal(id, newAutostart, project.name, project.description);
                } else {
                    console.error("showAutostartConfirmModal function not available.");
                }
            }

        });
    } else {
        console.error("Project list element not found.");
    }

    // 上下文菜单相关事件
    if (typeof window.handleContextMenu === 'function') {
        document.addEventListener('contextmenu', window.handleContextMenu);
    } else {
        console.error("handleContextMenu function not available.");
    }

    if (typeof window.hideContextMenu === 'function') {
        document.addEventListener('click', window.hideContextMenu);
    } else {
        console.error("hideContextMenu function not available.");
    }

    if (typeof window.handleContextMenuKeyboard === 'function') {
        document.addEventListener('keydown', window.handleContextMenuKeyboard);
    } else {
        console.error("handleContextMenuKeyboard function not available.");
    }

    // 在 setupEventListeners 函数末尾添加确认按钮的事件监听
    const confirmAutostartButton = document.getElementById('confirmAutostartButton');
    if (confirmAutostartButton) {
        confirmAutostartButton.addEventListener('click', () => {
            const id = parseInt(document.getElementById('autostartProjectId').value);
            const newAutostart = parseInt(document.getElementById('autostartNewValue').value);
            const actionText = newAutostart === 1 ? '注册' : '卸载';

            // 关闭模态框
            const modal = bootstrap.Modal.getInstance(document.getElementById('autostartConfirmModal'));
            if (modal) {
                modal.hide();
            }

            // 显示操作提示
            if (typeof window.showNotification === 'function') {
                window.showNotification(`正在${actionText}自启...`, 'info');
            }

            if (typeof window.updateAutostart === 'function') {
                window.updateAutostart(id, newAutostart).then(success => {
                    if (success) {
                        if (typeof window.showNotification === 'function') {
                            window.showNotification(`${actionText}自启成功`, 'success');
                        }
                    } else {
                        if (typeof window.showNotification === 'function') {
                            window.showNotification(`${actionText}自启失败`, 'danger');
                        }
                    }
                });
            } else {
                console.error("updateAutostart function not available.");
                if (typeof window.showNotification === 'function') {
                    window.showNotification('自启功能不可用', 'danger');
                }
            }
        });
    } else {
        console.error("Confirm autostart button not found.");
    }

    // 清理页面卸载时的资源
    window.addEventListener('beforeunload', () => {
        if (typeof window.stopAutoRegister === 'function') {
            window.stopAutoRegister();
        } else {
            console.error("stopAutoRegister function not available during unload.");
        }
        // Assuming destroyAllTooltips is in ui.js and globally accessible or called correctly
        if (typeof window.destroyAllTooltips === 'function') {
            window.destroyAllTooltips();
        } else {
            console.error("destroyAllTooltips function not available during unload.");
        }
    });
};

/**
 * 处理右键菜单事件
 * @param {Event} e - 事件对象
 */
window.handleContextMenu = function (e) {
    if (e.target.classList.contains('code-block')) {
        e.preventDefault();
        const menu = document.getElementById('contextMenu');
        contextMenuTarget = e.target; // contextMenuTarget is managed within app.js

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
};

/**
 * 隐藏右键菜单
 */
window.hideContextMenu = function () {
    document.getElementById('contextMenu').style.display = 'none';
    contextMenuTarget = null; // contextMenuTarget is managed within app.js
};

/**
 * 处理右键菜单键盘事件
 * @param {KeyboardEvent} e - 键盘事件对象
 */
window.handleContextMenuKeyboard = function (e) {
    const menu = document.getElementById('contextMenu');
    if (menu.style.display === 'block') {
        if (e.key === 'Escape') {
            if (typeof window.hideContextMenu === 'function') {
                window.hideContextMenu();
            } else {
                console.error("hideContextMenu function not available during keyboard event.");
            }
        } else if (e.key === 'Enter') {
            if (contextMenuTarget) { // contextMenuTarget is managed within app.js
                if (typeof window.copyText === 'function') {
                    window.copyText(contextMenuTarget.textContent);
                } else {
                    console.error("copyText function not available during keyboard event.");
                }
            }
            if (typeof window.hideContextMenu === 'function') {
                window.hideContextMenu();
            } else {
                console.error("hideContextMenu function not available during keyboard event.");
            }
        }
    }
};



// ===== 初始化 =====

/**
 * 应用初始化
 */
async function initApp() {
    // 首次加载项目列表
    if (typeof window.fetchProjects === 'function') {
        await window.fetchProjects();
    } else {
        console.error("fetchProjects function not available. Cannot initialize app.");
        return; // Stop initialization if critical function is missing
    }

    // 设置事件监听器
    if (typeof window.setupEventListeners === 'function') {
        window.setupEventListeners();
    } else {
        console.error("setupEventListeners function not available.");
    }


    // 启动自动注册
    if (typeof window.startAutoRegister === 'function') {
        window.startAutoRegister();
    } else {
        console.error("startAutoRegister function not available.");
    }


    // 隐藏通知区域
    const notificationArea = document.getElementById('notificationArea');
    if (notificationArea) {
        notificationArea.classList.remove('show');
    } else {
        console.warn("Notification area element not found.");
    }
}

// 页面加载完成后初始化应用
document.addEventListener('DOMContentLoaded', initApp);
