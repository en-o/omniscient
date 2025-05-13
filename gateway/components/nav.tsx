// gateway/components/nav.tsx
'use client'

import { useState, useEffect } from 'react'
import PmLogo from "@components/logos/Pm";

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
    const [formData, setFormData] = useState({
        url: '',
        description: ''
    })

    // 从 localStorage 加载服务器列表
    useEffect(() => {
        const savedServers = localStorage.getItem('servers')
        if (savedServers) {
            setServers(JSON.parse(savedServers))
        }
    }, [])

    // 提交表单
    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        if (!formData.url.trim() || !formData.description.trim()) {
            alert('请填写完整信息')
            return
        }

        const newServer = {
            id: generateId(),
            ...formData
        }

        const updatedServers = [...servers, newServer]
        setServers(updatedServers)
        localStorage.setItem('servers', JSON.stringify(updatedServers))
        setShowModal(false)
        setFormData({ url: '', description: '' })
    }

    // 删除服务器
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
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-lg w-full max-w-2xl">
                        <div className="p-6">
                            <div className="flex justify-between items-center mb-6">
                                <h5 className="text-xl font-bold">服务器管理</h5>
                                <button
                                    onClick={() => setShowModal(false)}
                                    className="text-gray-500 hover:text-gray-700"
                                >
                                    <span className="text-2xl">&times;</span>
                                </button>
                            </div>

                            <div className="mb-6">
                                <div className="flex justify-between items-center mb-3">
                                    <h6 className="font-bold">已注册服务器</h6>
                                    {servers.length > 0 && (
                                        <button
                                            onClick={() => {
                                                if (confirm('确定要清空所有服务器吗？')) {
                                                    setServers([]);
                                                    localStorage.removeItem('servers');
                                                    setSelectedServer('');
                                                }
                                            }}
                                            className="text-red-500 hover:text-red-600 text-sm flex items-center gap-1"
                                        >
                                            <i className="bi bi-trash"></i>
                                            清空全部
                                        </button>
                                    )}
                                </div>
                                <ul className="space-y-2">
                                    {servers.map(server => (
                                        <li key={server.id}
                                            className="flex justify-between items-center p-3 bg-gray-50 dark:bg-gray-700 rounded"
                                        >
                                            <span>
                                                {server.url}
                                                <small className="text-gray-500 ml-2">({server.description})</small>
                                            </span>
                                            <button
                                                onClick={() => handleDeleteServer(server.id)}
                                                className="text-red-500 hover:text-red-600"
                                            >
                                                <i className="bi bi-trash"></i>
                                            </button>
                                        </li>
                                    ))}
                                    {servers.length === 0 && (
                                        <li className="text-gray-500 text-center p-3">暂无注册服务器</li>
                                    )}
                                </ul>
                            </div>

                            <div>
                                <h6 className="font-bold mb-3">添加新服务器</h6>
                                <form onSubmit={handleSubmit} className="space-y-4">
                                    <div className="grid grid-cols-1 gap-4">
                                        <input
                                            type="url"
                                            placeholder="服务器 Base URL (例如: http://localhost:8080)"
                                            required
                                            value={formData.url}
                                            onChange={(e) => setFormData({...formData, url: e.target.value})}
                                            className="w-full px-4 py-2 border rounded-md dark:bg-gray-600 dark:text-white"
                                        />
                                        <input
                                            type="text"
                                            placeholder="服务器描述"
                                            required
                                            value={formData.description}
                                            onChange={(e) => setFormData({...formData, description: e.target.value})}
                                            className="w-full px-4 py-2 border rounded-md dark:bg-gray-600 dark:text-white"
                                        />
                                        <button
                                            type="submit"
                                            className="w-full px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600"
                                        >
                                            添加服务器
                                        </button>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </>
    )
}