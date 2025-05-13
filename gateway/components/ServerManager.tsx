'use client'

import { useState } from 'react'

interface Server {
    id: string
    url: string
    description: string
}

interface ServerManagerProps {
    servers: Server[]
    onServerAdd: (server: Omit<Server, 'id'>) => void
    onServerDelete: (id: string) => void
    onClose: () => void
}

export default function ServerManager({
                                          servers,
                                          onServerAdd,
                                          onServerDelete,
                                          onClose
                                      }: ServerManagerProps) {
    const [formData, setFormData] = useState({
        url: '',
        description: ''
    })

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        if (!formData.url.trim() || !formData.description.trim()) {
            alert('请填写完整信息')
            return
        }

        onServerAdd(formData)
        setFormData({ url: '', description: '' })
    }

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"> {/* Added z-50 for modal stacking */}
            <div className="bg-white dark:bg-gray-800 rounded-lg w-full max-w-2xl max-h-[90vh] overflow-y-auto"> {/* Added max-h and overflow for long lists */}
                <div className="p-6">
                    <div className="flex justify-between items-center mb-6">
                        <h5 className="text-xl font-bold">服务器管理</h5>
                        <button
                            onClick={onClose}
                            className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
                        >
                            <span className="text-2xl">&times;</span>
                        </button>
                    </div>

                    <div className="mb-6">
                        <div className="flex justify-between items-center mb-3">
                            <h6 className="font-bold">已注册服务器</h6>
                        </div>
                        <ul className="space-y-2">
                            {servers.map(server => (
                                <li key={server.id}
                                    className="flex justify-between items-center p-3 bg-gray-50 dark:bg-gray-700 rounded"
                                >
                                    <span>
                                        {server.url}
                                        <small className="text-gray-500 dark:text-gray-400 ml-2">({server.description})</small>
                                    </span>
                                    <button
                                        onClick={() => onServerDelete(server.id)}
                                        className="text-red-500 hover:text-red-600 dark:hover:text-red-400"
                                        aria-label={`删除服务器 ${server.description}`}
                                    >
                                        <i className="bi bi-trash"> 删除 </i>
                                    </button>
                                </li>
                            ))}
                            {servers.length === 0 && (
                                <li className="text-gray-500 dark:text-gray-400 text-center p-3">暂无注册服务器</li>
                            )}
                        </ul>
                    </div>

                    <div>
                        <h6 className="font-bold mb-3">添加新服务器</h6>
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div className="grid grid-cols-1 gap-4">
                                <input
                                    type="url"
                                    placeholder="服务器 URL (例如: http://localhost:8080/html/pm.html)"
                                    required
                                    value={formData.url}
                                    onChange={(e) => setFormData({...formData, url: e.target.value})}
                                    className="w-full px-4 py-2 border rounded-md dark:bg-gray-600 dark:text-white dark:border-gray-700"
                                />
                                <input
                                    type="text"
                                    placeholder="服务器描述"
                                    required
                                    value={formData.description}
                                    onChange={(e) => setFormData({...formData, description: e.target.value})}
                                    className="w-full px-4 py-2 border rounded-md dark:bg-gray-600 dark:text-white dark:border-gray-700"
                                />
                                <button
                                    type="submit"
                                    className="w-full px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 dark:hover:bg-green-700"
                                >
                                    添加服务器
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    )
}
