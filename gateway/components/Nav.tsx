'use client'

import { useState } from 'react'
import PmLogo from "@components/logos/Pm"
import ServerManager from './ServerManager'
import { useServer } from './ServerContext'

export default function Nav() {
    const [selectedServer, setSelectedServer] = useState<string>('')
    const [showModal, setShowModal] = useState(false)
    const {
        servers,
        isLoading,
        error,
        setSelectedServerUrl,
        refreshIframe, // 从ServerContext获取刷新函数
        isRefreshing  // 刷新状态
    } = useServer()

    // 处理服务器选择变化
    const handleServerChange = (serverId: string) => {
        setSelectedServer(serverId)
        const selectedServer = servers.find(server => server.id === serverId)
        // 设置全局数据
        setSelectedServerUrl(selectedServer ? selectedServer.url : '')
    }

    // 处理刷新操作
    const handleRefresh = () => {
        refreshIframe()
    }

    return (
        <>
            {/* Nav bar with padding */}
            <nav className="w-full bg-white dark:bg-gray-800 shadow-md p-4">
                {/* Full-width flex container to push items to edges */}
                <div className="flex items-center justify-between w-full">
                    {/* Left side content: Logo and Select, kept together */}
                    <div className="flex items-center gap-4">
                        <PmLogo />
                        <div className="flex items-center">
                            <i className="bi bi-hdd text-gray-500 text-xl mr-2"></i>
                            <select
                                value={selectedServer}
                                onChange={(e) => handleServerChange(e.target.value)}
                                className="block w-64 px-4 py-2 border rounded-md dark:bg-gray-700 dark:text-white"
                                disabled={isLoading}
                            >
                                <option value="">选择服务器</option>
                                {servers.map(server => (
                                    <option key={server.id} value={server.id}>
                                        {server.description}
                                    </option>
                                ))}
                            </select>
                            {isLoading && (
                                <span className="ml-2 text-gray-500">
                                    <i className="bi bi-arrow-repeat animate-spin"></i>
                                </span>
                            )}
                        </div>

                        {/* 刷新按钮 - 位于服务器选择旁边 */}
                        {selectedServer && (
                            <button
                                onClick={handleRefresh}
                                disabled={isRefreshing || isLoading}
                                className={`px-4 py-2 rounded-md flex items-center gap-2 ${
                                    isRefreshing || isLoading
                                        ? 'bg-gray-400 text-gray-200 cursor-not-allowed'
                                        : 'bg-green-500 text-white hover:bg-green-600'
                                }`}
                                title="刷新服务器页面(清除缓存)"
                            >
                                <i className={`bi ${isRefreshing ? 'bi-arrow-repeat animate-spin' : 'bi-arrow-clockwise'}`}></i>
                                {isRefreshing ? '刷新中...' : '刷新页面'}
                            </button>
                        )}
                    </div>

                    {/* Right side content: Server Management button */}
                    <button
                        onClick={() => setShowModal(true)}
                        className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 flex items-center gap-2"
                    >
                        <i className="bi bi-gear-fill"></i>
                        服务器管理
                    </button>
                </div>
            </nav>

            {/* 错误提示 */}
            {error && (
                <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mt-2 mx-4" role="alert">
                    <strong className="font-bold">错误: </strong>
                    <span className="block sm:inline">{error}</span>
                </div>
            )}

            {showModal && (
                <ServerManager
                    onClose={() => setShowModal(false)}
                />
            )}
        </>
    )
}