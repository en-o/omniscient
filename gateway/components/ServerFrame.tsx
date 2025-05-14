'use client'

import {useServer} from "@components/ServerContext"
import {useState, useEffect} from "react"
import NotFound404 from "@components/NotFound404"

export default function ServerFrame() {
    const {selectedServerUrl} = useServer()
    const [iframeError, setIframeError] = useState(false)
    const [currentUrl, setCurrentUrl] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(false)

    // 处理 URL
    const getFormattedUrl = (url: string) => {
        if (!url) return ''
        const suffix = '/html/pm.html'
        return url.endsWith(suffix) ? url : `${url}${suffix}`
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

        // 尝试访问URL检查其是否可用
        fetch(formattedUrl, { method: 'HEAD', mode: 'no-cors' })
            .then(() => {
                // 因为no-cors模式总是成功，我们认为请求至少能发出去
                setCurrentUrl(formattedUrl)
                setIframeError(false)
            })
            .catch(() => {
                // 请求失败，URL不可访问
                setIframeError(true)
                setCurrentUrl(null)
            })
            .finally(() => {
                setIsLoading(false)
            })
    }, [selectedServerUrl])

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
            <div
                className="w-full h-[calc(100vh-8rem)] flex items-center justify-center text-gray-500 dark:text-gray-400 text-lg bg-white dark:bg-gray-800 rounded-lg shadow-sm">
                请在导航栏选择一个服务器
            </div>
        )
    }

    if (isLoading) {
        return (
            <div className="w-full h-[calc(100vh-8rem)] flex flex-col items-center justify-center bg-white dark:bg-gray-800 rounded-lg shadow-sm">
                <div className="flex items-center space-x-2">
                    <div className="w-8 h-8 border-4 border-t-blue-500 border-b-blue-500 border-l-transparent border-r-transparent rounded-full animate-spin"></div>
                    <span className="text-gray-600 dark:text-gray-300">正在验证服务器连接...</span>
                </div>
            </div>
        )
    }

    if (iframeError) {
        return <NotFound404 />
    }

    return (
        <iframe
            src={currentUrl}
            title="服务器内容"
            className="w-full h-[calc(100vh-8rem)] border-0 rounded-lg bg-white shadow-sm"
            onError={handleIframeError}
            onLoad={handleIframeLoad}
        >
            您的浏览器不支持 iframe。
        </iframe>
    )
}
