'use client'

import { useServer } from "@components/ServerContext"

export default function ServerFrame() {
    const { selectedServerUrl } = useServer()
    return (
        <>
            {selectedServerUrl ? (
                <iframe
                    src={`${selectedServerUrl}/html/pm.html`}
                    title="服务器内容"
                    className="w-full h-[calc(100vh-8rem)] border-0 rounded-lg bg-white shadow-sm"
                >
                    您的浏览器不支持 iframe。
                </iframe>
            ) : (
                <div className="w-full h-[calc(100vh-8rem)] flex items-center justify-center text-gray-500 dark:text-gray-400 text-lg bg-white dark:bg-gray-800 rounded-lg shadow-sm">
                    请在导航栏选择一个服务器
                </div>
            )}
        </>
    )
}
