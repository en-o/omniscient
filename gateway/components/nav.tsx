// gateway/components/nav.tsx
'use client'

import { useState, useEffect } from 'react'
import PmLogo from "@components/logos/Pm";
import ServerManager from './ServerManager'


interface Server {
    id: string
    url: string
    description: string
}
// 生成唯一ID的函数
const generateId = () => {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = Math.random() * 16 | 0
        const v = c === 'x' ? r : (r & 0x3 | 0x8)
        return v.toString(16)
    })
}

export default function Nav() {
    const [servers, setServers] = useState<Server[]>([])
    const [selectedServer, setSelectedServer] = useState<string>('')
    const [showModal, setShowModal] = useState(false)


    // 从 localStorage 加载服务器列表
    useEffect(() => {
        const savedServers = localStorage.getItem('servers')
        if (savedServers) {
            setServers(JSON.parse(savedServers))
        }
    }, [])

    const handleAddServer = (newServer: Omit<Server, 'id'>) => {
        const serverWithId = {
            id: generateId(),
            ...newServer
        }
        const updatedServers = [...servers, serverWithId]
        setServers(updatedServers)
        localStorage.setItem('servers', JSON.stringify(updatedServers))
    }

    const handleDeleteServer = (id: string) => {
        if (!confirm('确定要移除此服务器吗？')) return

        const updatedServers = servers.filter(server => server.id !== id)
        setServers(updatedServers)
        localStorage.setItem('servers', JSON.stringify(updatedServers))
        if (selectedServer === id) {
            setSelectedServer('')
        }
    }


    const handleClearServers = () => {
        setServers([])
        localStorage.removeItem('servers')
        setSelectedServer('')
    }

    return (
        <>
            <nav className="w-full bg-white dark:bg-gray-800 shadow-md p-4">
                <div className="max-w-7xl mx-auto flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <PmLogo />
                        <div className="flex items-center">
                            <i className="bi bi-hdd text-gray-500 text-xl mr-2"></i>
                            <select
                                value={selectedServer}
                                onChange={(e) => setSelectedServer(e.target.value)}
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
                    onServersClear={handleClearServers}
                    onClose={() => setShowModal(false)}
                />
            )}
        </>
    )
}