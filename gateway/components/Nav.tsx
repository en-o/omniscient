// gateway/components/Nav.tsx
'use client'

import { useState, useEffect } from 'react'
import PmLogo from "@components/logos/Pm";
import ServerManager from './ServerManager'
import {generateId} from "@utils/uuid";
import { useServer } from './ServerContext'

interface Server {
    id: string
    url: string
    description: string
}


export default function Nav() {
    const [servers, setServers] = useState<Server[]>([])
    const [selectedServer, setSelectedServer] = useState<string>('')
    const [showModal, setShowModal] = useState(false)
    const { setSelectedServerUrl } = useServer()

    // 处理服务器选择变化
    const handleServerChange = (serverId: string) => {
        setSelectedServer(serverId)
        const selectedServer = servers.find(server => server.id === serverId)
        // 设置全局数据
        setSelectedServerUrl(selectedServer ? selectedServer.url : '')
    }

    // 从 localStorage 加载服务器列表
    useEffect(() => {
        const savedServers = localStorage.getItem('omniscient-servers')
        if (savedServers) {
            setServers(JSON.parse(savedServers))
        }
    }, [])

    // 添加列表
    const handleAddServer = (newServer: Omit<Server, 'id'>) => {
        const serverWithId = {
            id: generateId(),
            ...newServer
        }
        const updatedServers = [...servers, serverWithId]
        setServers(updatedServers)
        localStorage.setItem('servers', JSON.stringify(updatedServers))
    }

    // 删除列表
    const handleDeleteServer = (id: string) => {
        if (!confirm('确定要移除此服务器吗？')) return

        const updatedServers = servers.filter(server => server.id !== id)
        setServers(updatedServers)
        localStorage.setItem('servers', JSON.stringify(updatedServers))
        if (selectedServer === id) {
            setSelectedServer('')
        }
    }


    return (
        <>
            {/* Nav bar with padding */}
            <nav className="w-full bg-white dark:bg-gray-800 shadow-md p-4">
                {/* Full-width flex container to push items to edges */}
                {/* max-w-7xl removed here to allow full width flex positioning */}
                <div className="flex items-center justify-between w-full">
                    {/* Left side content: Logo and Select, kept together */}
                    <div className="flex items-center gap-4"> {/* Added gap-4 for space */}
                        <PmLogo />
                        <div className="flex items-center">
                            <i className="bi bi-hdd text-gray-500 text-xl mr-2"></i>
                            <select
                                value={selectedServer}
                                onChange={(e) => handleServerChange(e.target.value)}
                                className="block w-64 px-4 py-2 border rounded-md dark:bg-gray-700 dark:text-white"
                            >
                                <option value="">选择服务器</option>
                                {servers.map(server => (
                                    <option key={server.id} value={server.id}>
                                        {server.description}
                                    </option>
                                ))}
                            </select>
                        </div>
                    </div>
                    {/* Right side content: Server Management button */}
                    {/* justify-between on the parent pushes this to the right */}
                    <button
                        onClick={() => setShowModal(true)}
                        className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 flex items-center gap-2"
                    >
                        <i className="bi bi-gear-fill"></i>
                        服务器管理
                    </button>
                </div>
            </nav>

            {showModal && (
                <ServerManager
                    servers={servers}
                    onServerAdd={handleAddServer}
                    onServerDelete={handleDeleteServer}
                    onClose={() => setShowModal(false)}
                />
            )}
        </>
    )
}
