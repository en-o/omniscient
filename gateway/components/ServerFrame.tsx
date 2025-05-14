'use client'

import {useServer} from "@components/ServerContext"
import {useState, useEffect} from "react"

export default function ServerFrame() {
    const {selectedServerUrl} = useServer()
    const [iframeError, setIframeError] = useState(false)
    const [currentUrl, setCurrentUrl] = useState<string | null>(null)

    // 处理 URL
    const getFormattedUrl = (url: string) => {
        if (!url) return ''
        const suffix = '/html/pm.html'
        return url.endsWith(suffix) ? url : `${url}${suffix}`
    }

    // 验证URL并更新当前URL
    useEffect(() => {
        if (!selectedServerUrl) {
            setCurrentUrl(null)
            setIframeError(false)
            return
        }

        // 验证URL格式
        try {
            new URL(getFormattedUrl(selectedServerUrl))
            setCurrentUrl(getFormattedUrl(selectedServerUrl))
            setIframeError(false)
        } catch (error) {
            setIframeError(true)
            setCurrentUrl(null)
        }
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

    if (iframeError) {
        return (
            <div
                className="w-full h-[calc(100vh-8rem)] flex items-center justify-center text-red-500 dark:text-red-400 text-lg bg-white dark:bg-gray-800 rounded-lg shadow-sm">
                无法加载页面，请重新设置服务器 URL
            </div>
        )
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
