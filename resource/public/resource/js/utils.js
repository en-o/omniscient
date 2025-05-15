// ===== 工具函数 =====

/**
 * 显示通知消息
 * @param {string} message - 消息内容
 * @param {string} type - 消息类型 (success, danger, warning, info)
 */
function showNotification(message, type = 'success') {
    const notificationArea = document.getElementById('notificationArea');
    const notificationMessage = document.getElementById('notificationMessage');

    notificationArea.classList.remove('alert-success', 'alert-danger', 'alert-warning', 'alert-info');
    notificationArea.classList.add(`alert-${type}`);
    notificationArea.classList.add('show');

    notificationMessage.textContent = message;

    // 自动隐藏通知（3秒后）
    setTimeout(() => {
        notificationArea.classList.remove('show');
    }, 3000);
}

/**
 * HTML转义防止XSS攻击
 * @param {string} str - 需要转义的字符串
 * @returns {string} - 转义后的字符串
 */
function escapeHtml(str) {
    if (!str) return '';

    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

/**
 * 复制文本到剪贴板
 * @param {string} text - 需要复制的文本
 */
async function copyText(text) {
    if (!text) return;

    text = text.trim();

    try {
        await navigator.clipboard.writeText(text);
        showNotification('复制成功！');
    } catch (err) {
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
                showNotification('复制成功！');
            } else {
                throw new Error('复制命令执行失败');
            }
        } catch (err) {
            showNotification('复制失败：请手动复制', 'danger');
        }

        document.body.removeChild(textArea);
    }
}

/**
 * 销毁所有工具提示实例
 */
function destroyAllTooltips() {
    tooltips.forEach(tooltip => {
        tooltip.dispose();
    });
    tooltips = [];
}

/**
 * 初始化工具提示
 */
function initTooltips() {
    // 先清除旧的工具提示
    destroyAllTooltips();

    // 创建新的工具提示
    const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltips = tooltipTriggerList.map(tooltipTriggerEl => {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });
}

/**
 * 更新页面标题
 * @param {string} worker - 工作节点名称
 */
function updatePageTitle(worker) {
    if (worker) {
        const title = `${worker} - Java项目管理`;
        document.getElementById('serverName').textContent = title;
        document.title = title;
    }
}