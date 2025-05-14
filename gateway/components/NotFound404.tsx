'use client'

import Image from "next/image"

export default function NotFound404() {
    return (
        <div className="w-full h-[calc(100vh-8rem)] flex flex-col items-center justify-center bg-white dark:bg-gray-800 rounded-lg shadow-sm">
            <div className="text-center">
                <Image
                    src="/404.svg"
                    alt="404 错误图标"
                    width={200}
                    height={200}
                    className="mx-auto mb-6"
                    priority
                />
                <h2 className="text-2xl font-bold text-gray-800 dark:text-gray-200 mb-2">
                    找不到服务器
                </h2>
                <p className="text-gray-600 dark:text-gray-400 mb-6">
                    您尝试访问的服务器地址无法连接或不存在
                </p>
                <div className="p-4 bg-gray-100 dark:bg-gray-700 rounded-lg mb-4 max-w-md mx-auto">
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                        请检查服务器 URL 是否正确，或者服务器是否处于运行状态。
                    </p>
                </div>
            </div>
        </div>
    )
}
