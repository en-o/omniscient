'use client'
import { useServer } from "@components/ServerContext"
import { useState, useEffect, useRef } from "react"
import NotFound404 from "@components/NotFound404"

export default function ServerFrame() {
    const { selectedServerUrl, refreshKey, isRefreshing } = useServer()
    const [iframeError, setIframeError] = useState(false)
    const [currentUrl, setCurrentUrl] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(false)
    const iframeRef = useRef<HTMLIFrameElement>(null)

    // 处理 URL
    const getFormattedUrl = (url: string) => {
        if (!url) return ''
        const suffix = '/html/pm.html'
        return url.endsWith(suffix) ? url : `${url}${suffix}`
    }

    // 添加时间戳或随机参数来防止缓存
    const getUrlWithCacheBuster = (url: string) => {
        if (!url) return ''
        const separator = url.includes('?') ? '&' : '?'
        return `${url}${separator}_cache=${refreshKey}`
    }

    // 验证URL是否可访问并更新当前URL
    useEffect(() => {
        if (!selectedServerUrl) {
            setCurrentUrl(null)
            setIframeError(false)
            return
        }

        // 先验证URL格式
        let formattedUrl
        try {
            formattedUrl = getFormattedUrl(selectedServerUrl)
            new URL(formattedUrl)
        } catch (error) {
            setIframeError(true)
            setCurrentUrl(null)
            return
        }

        // 显示加载状态
        setIsLoading(true)
        setIframeError(false)

        // 设置30秒超时
        const timeoutId = setTimeout(() => {
            setIsLoading(false)
            setIframeError(true)
            setCurrentUrl(null)
        }, 30000) // 30秒超时

        // 尝试访问URL检查其是否可用
        fetch(formattedUrl, { method: 'HEAD', mode: 'no-cors' })
            .then(() => {
                // 清除超时计时器
                clearTimeout(timeoutId)
                // 因为no-cors模式总是成功，我们认为请求至少能发出去
                setCurrentUrl(formattedUrl)
                setIframeError(false)
            })
            .catch(() => {
                // 清除超时计时器
                clearTimeout(timeoutId)
                // 请求失败，URL不可访问
                setIframeError(true)
                setCurrentUrl(null)
            })
            .finally(() => {
                setIsLoading(false)
            })

        // 清理函数，组件卸载或者依赖项变化时执行
        return () => {
            clearTimeout(timeoutId)
        }
    }, [selectedServerUrl])

    // 监听refreshKey的变化，当它更新时重载iframe内容
    useEffect(() => {
        if (refreshKey > 0 && iframeRef.current && currentUrl) {
            try {
                // 尝试以编程方式重新加载iframe
                const iframeDoc = iframeRef.current.contentDocument ||
                    (iframeRef.current.contentWindow?.document)
                if (iframeDoc) {
                    iframeDoc.location.reload() //重新加载
                }
            } catch (e) {
                // 如果因为跨域问题无法访问contentDocument，使用src刷新方法
                // 先临时清空URL
                const tempUrl = currentUrl
                setCurrentUrl(null)
                // 使用setTimeout确保DOM有时间更新
                setTimeout(() => {
                    setCurrentUrl(tempUrl)
                }, 50)
            }
        }
    }, [refreshKey, currentUrl])

    // 处理 iframe 加载错误
    const handleIframeError = () => {
        setIframeError(true)
    }

    // 处理 iframe 加载成功
    const handleIframeLoad = () => {
        setIframeError(false)
    }

    if (!selectedServerUrl) {
        return (
            <div className="w-full h-[calc(100vh-8rem)] flex items-center justify-center text-gray-500 dark:text-gray-400 text-lg bg-white dark:bg-gray-800 rounded-lg shadow-sm">
                请在导航栏选择一个服务器
            </div>
        )
    }

    if (isLoading) {
        return (
            <div className="w-full h-[calc(100vh-8rem)] flex flex-col items-center justify-center bg-white dark:bg-gray-800 rounded-lg shadow-sm">
                <div className="flex flex-col items-center space-y-4">
                    <div className="flex items-center space-x-2">
                        <div className="w-8 h-8 border-4 border-t-blue-500 border-b-blue-500 border-l-transparent border-r-transparent rounded-full animate-spin"></div>
                        <span className="text-gray-600 dark:text-gray-300">正在验证服务器连接...</span>
                    </div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                        如果30秒内未响应，将视为连接失败
                    </p>
                </div>
            </div>
        )
    }

    if (iframeError) {
        return <NotFound404 />
    }

    return (
        <div className="w-full h-[calc(100vh-8rem)] flex flex-col bg-white dark:bg-gray-800 rounded-lg shadow-sm">
            {isRefreshing && (
                <div className="absolute inset-0 bg-black bg-opacity-20 flex items-center justify-center z-10">
                    <div className="bg-white p-4 rounded-lg shadow-lg flex items-center space-x-3">
                        <div className="w-6 h-6 border-4 border-t-blue-500 border-b-blue-500 border-l-transparent border-r-transparent rounded-full animate-spin"></div>
                        <span>正在刷新页面...</span>
                    </div>
                </div>
            )}

            <iframe
                ref={iframeRef}
                key={refreshKey} // 使用key强制重新创建iframe元素
                src={currentUrl ? getUrlWithCacheBuster(currentUrl) : undefined}
                title="服务器内容"
                className="w-full h-full border-0 rounded-lg"
                onError={handleIframeError}
                onLoad={handleIframeLoad}
            >
                您的浏览器不支持 iframe。
            </iframe>
        </div>
    )
}