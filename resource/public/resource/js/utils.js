// ===== 工具函数 =====


/**
 * 显示通知消息
 * @param {string} message - 消息内容
 * @param {string} type - 消息类型 (success, danger, warning, info)
 */
window.showNotification = function(message, type = 'success') {
    const notificationArea = document.getElementById('notificationArea');
    const notificationMessage = document.getElementById('notificationMessage');

    if (!notificationArea || !notificationMessage) {
        console.error("Notification area elements not found.");
        return;
    }

    notificationArea.classList.remove('alert-success', 'alert-danger', 'alert-warning', 'alert-info');
    notificationArea.classList.add(`alert-${type}`);
    notificationArea.classList.add('show');

    notificationMessage.textContent = message;

    // 自动隐藏通知（3秒后）
    setTimeout(() => {
        notificationArea.classList.remove('show');
    }, 3000);
};

/**
 * HTML转义防止XSS攻击
 * @param {string} str - 需要转义的字符串
 * @returns {string} - 转义后的字符串
 */
window.escapeHtml = function(str) {
    if (!str) return '';

    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
};

/**
 * 复制文本到剪贴板
 * @param {string} text - 需要复制的文本
 */
window.copyText = async function(text) {
    if (!text) return;

    text = text.trim();

    try {
        await navigator.clipboard.writeText(text);
        if (typeof window.showNotification === 'function') {
            window.showNotification('复制成功！'); // showNotification is now window.showNotification, accessible globally
        }
    } catch (err) {
        console.warn("Clipboard API writeText failed, attempting fallback copy.", err);
        // 降级处理方案
        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'fixed';
        textArea.style.left = '-9999px';
        textArea.style.top = '-9999px';
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();

        try {
            const successful = document.execCommand('copy');
            if (successful) {
                if (typeof window.showNotification === 'function') {
                    window.showNotification('复制成功！'); // showNotification is now window.showNotification
                }
            } else {
                throw new Error('Copy command execution failed');
            }
        } catch (err) {
            console.error("Manual copy fallback failed.", err);
            if (typeof window.showNotification === 'function') {
                window.showNotification('复制失败：请手动复制', 'danger'); // showNotification is now window.showNotification
            }
        }

        document.body.removeChild(textArea);
    }
};

// 销毁所有工具提示实例 and initTooltips are moved to ui.js


/**
 * 更新页面标题
 * @param {string} worker - 工作节点名称
 */
window.updatePageTitle = function(worker) {
    const serverNameElement = document.getElementById('serverName');
    if (worker && serverNameElement) {
        const title = `${worker} - Java项目管理`;
        serverNameElement.textContent = title;
        document.title = title;
    } else if (serverNameElement) {
        serverNameElement.textContent = 'Java项目管理'; // Reset if worker is empty
        document.title = '项目管理';
    }
};