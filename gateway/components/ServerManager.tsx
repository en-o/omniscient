// gateway/components/ServerManager.tsx
'use client'

import { useState } from 'react'
import { useServer } from './ServerContext'

interface ServerManagerProps {
    onClose: () => void
}

export default function ServerManager({ onClose }: ServerManagerProps) {
    const { servers, addServer, deleteServer, isLoading, error } = useServer()

    const [formData, setFormData] = useState({
        url: '',
        description: ''
    })
    const [formError, setFormError] = useState<string | null>(null)
    const [isSubmitting, setIsSubmitting] = useState(false)

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        // 表单验证
        if (!formData.url.trim() || !formData.description.trim()) {
            setFormError('请填写完整信息')
            return
        }

        try {
            setIsSubmitting(true)
            setFormError(null)

            // 添加服务器
            await addServer(formData)

            // 重置表单
            setFormData({ url: '', description: '' })
        } catch (err) {
            setFormError(err instanceof Error ? err.message : '添加服务器失败')
        } finally {
            setIsSubmitting(false)
        }
    }

    const handleDeleteServer = async (id: string) => {
        if (!confirm('确定要移除此服务器吗？')) return

        await deleteServer(id)
    }

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white dark:bg-gray-800 rounded-lg w-full max-w-2xl max-h-[90vh] overflow-y-auto">
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

                    {/* 全局错误提示 */}
                    {error && (
                        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
                            <span className="block sm:inline">{error}</span>
                        </div>
                    )}

                    <div className="mb-6">
                        <div className="flex justify-between items-center mb-3">
                            <h6 className="font-bold">已注册服务器</h6>
                            {isLoading && (
                                <span className="text-gray-500">
                                    <i className="bi bi-arrow-repeat animate-spin mr-1"></i>
                                    加载中...
                                </span>
                            )}
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
                                        onClick={() => handleDeleteServer(server.id)}
                                        className="text-red-500 hover:text-red-600 dark:hover:text-red-400"
                                        aria-label={`删除服务器 ${server.description}`}
                                    >
                                        <i className="bi bi-trash"> 删除 </i>
                                    </button>
                                </li>
                            ))}
                            {servers.length === 0 && !isLoading && (
                                <li className="text-gray-500 dark:text-gray-400 text-center p-3">暂无注册服务器</li>
                            )}
                        </ul>
                    </div>

                    <div>
                        <h6 className="font-bold mb-3">添加新服务器</h6>
                        {formError && (
                            <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-2 rounded relative mb-4" role="alert">
                                <span className="block sm:inline">{formError}</span>
                            </div>
                        )}
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div className="grid grid-cols-1 gap-4">
                                <input
                                    type="url"
                                    placeholder="服务器 URL (例如: http://localhost:8080/html/pm.html)"
                                    required
                                    value={formData.url}
                                    onChange={(e) => setFormData({...formData, url: e.target.value})}
                                    className="w-full px-4 py-2 border rounded-md dark:bg-gray-600 dark:text-white dark:border-gray-700"
                                    disabled={isSubmitting}
                                />
                                <input
                                    type="text"
                                    placeholder="服务器描述"
                                    required
                                    value={formData.description}
                                    onChange={(e) => setFormData({...formData, description: e.target.value})}
                                    className="w-full px-4 py-2 border rounded-md dark:bg-gray-600 dark:text-white dark:border-gray-700"
                                    disabled={isSubmitting}
                                />
                                <button
                                    type="submit"
                                    className={`w-full px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 dark:hover:bg-green-700 ${
                                        isSubmitting ? 'opacity-70 cursor-not-allowed' : ''
                                    }`}
                                    disabled={isSubmitting}
                                >
                                    {isSubmitting ? (
                                        <>
                                            <i className="bi bi-arrow-repeat animate-spin mr-1"></i>
                                            添加中...
                                        </>
                                    ) : '添加服务器'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    )
}
