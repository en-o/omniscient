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
                    className="w-full h-full border-0 shadow-md rounded-lg"
                >
                    您的浏览器不支持 iframe。
                </iframe>
            ) : (
                <div className="w-full h-full flex items-center justify-center text-gray-500 dark:text-gray-400 text-lg">
                    请在导航栏选择一个服务器
                </div>
            )}
        </>
    )
}
