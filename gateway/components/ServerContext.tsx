'use client'

import { createContext, useContext, useState, ReactNode } from 'react'
import { ServerEntity } from "@typesss/serverEntity"

// 服务器上下文接口
interface ServerContextType {
    servers: ServerEntity[]
    isLoading: boolean
    error: string | null
    selectedServerUrl: string
    setSelectedServerUrl: (url: string) => void
    addServer: (server: Omit<ServerEntity, 'id'>) => Promise<ServerEntity>
    deleteServer: (id: string) => Promise<boolean>
    loadServers: () => Promise<void>
    exportServers: () => Promise<Blob | null>
    importServers: (data: any[]) => Promise<{ imported: number, failed: number, errors?: string[] }>
    resetDatabase: () => Promise<boolean>
    validateImportData: (data: any) => { valid: boolean, message?: string }
}

// 创建上下文
const ServerContext = createContext<ServerContextType | undefined>(undefined)

// 服务器提供者组件
export function ServerProvider({ children }: { children: ReactNode }) {
    const [servers, setServers] = useState<ServerEntity[]>([])
    const [selectedServerUrl, setSelectedServerUrl] = useState<string>('')
    const [isLoading, setIsLoading] = useState<boolean>(true)
    const [error, setError] = useState<string | null>(null)

    // 加载服务器列表
    const loadServers = async (): Promise<void> => {
        try {
            setIsLoading(true)
            setError(null)

            const response = await fetch('/api/servers')

            if (!response.ok) {
                throw new Error(`获取服务器列表失败 (${response.status}: ${response.statusText})`)
            }

            const data = await response.json()
            setServers(data)
        } catch (err) {
            setError(err instanceof Error ? err.message : '加载服务器时出错')
            console.error('加载服务器失败:', err)
        } finally {
            setIsLoading(false)
        }
    }

    // 添加新服务器
    const addServer = async (server: Omit<ServerEntity, 'id'>): Promise<ServerEntity> => {
        try {
            const response = await fetch('/api/servers', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(server),
            })

            if (!response.ok) {
                const errorData = await response.json().catch(() => null)
                throw new Error(errorData?.message || `添加服务器失败 (${response.status})`)
            }

            const newServer = await response.json()

            // 更新本地状态
            setServers(prev => [...prev, newServer])

            return newServer
        } catch (err) {
            setError(err instanceof Error ? err.message : '添加服务器时出错')
            console.error('添加服务器失败:', err)
            throw err
        }
    }

    // 删除服务器
    const deleteServer = async (id: string): Promise<boolean> => {
        try {
            const response = await fetch(`/api/servers/${id}`, {
                method: 'DELETE',
            })

            if (!response.ok) {
                throw new Error(`删除服务器失败 (${response.status}: ${response.statusText})`)
            }

            // 更新本地状态
            setServers(prev => prev.filter(server => server.id !== id))

            // 如果删除的是当前选中的服务器，清除选择
            const deletedServer = servers.find(server => server.id === id)
            if (deletedServer && deletedServer.url === selectedServerUrl) {
                setSelectedServerUrl('')
            }

            return true
        } catch (err) {
            setError(err instanceof Error ? err.message : '删除服务器时出错')
            console.error('删除服务器失败:', err)
            return false
        }
    }

    // 验证导入数据
    const validateImportData = (data: any): { valid: boolean, message?: string } => {
        // 检查数据是否为数组
        if (!Array.isArray(data)) {
            return { valid: false, message: '导入数据格式错误: 应为数组' }
        }

        // 检查数组是否为空
        if (data.length === 0) {
            return { valid: false, message: '导入数据为空' }
        }

        // 检查每个项是否包含必要字段
        for (let i = 0; i < data.length; i++) {
            const item = data[i]
            if (!item.url || typeof item.url !== 'string') {
                return { valid: false, message: `第 ${i + 1} 项缺少有效的 URL` }
            }
            if (!item.description || typeof item.description !== 'string') {
                return { valid: false, message: `第 ${i + 1} 项缺少有效的描述` }
            }
        }

        return { valid: true }
    }

    // 导出服务器数据
    const exportServers = async (): Promise<Blob | null> => {
        try {
            setIsLoading(true)
            setError(null)

            // 检查是否有服务器可导出
            if (servers.length === 0) {
                setError('没有服务器数据可导出')
                return null
            }

            // 调用导出API
            const response = await fetch('/api/servers/export')

            if (!response.ok) {
                throw new Error(`导出服务器数据失败 (${response.status}: ${response.statusText})`)
            }

            // 获取 blob 数据
            const blob = await response.blob()
            return blob
        } catch (err) {
            setError(err instanceof Error ? err.message : '导出服务器数据失败')
            console.error('导出服务器数据失败:', err)
            return null
        } finally {
            setIsLoading(false)
        }
    }

    // 导入服务器数据
    const importServers = async (data: any[]): Promise<{ imported: number, failed: number, errors?: string[] }> => {
        try {
            setIsLoading(true)
            setError(null)

            // 验证数据格式
            const validation = validateImportData(data)
            if (!validation.valid) {
                throw new Error(validation.message || '导入数据格式无效')
            }

            // 调用导入API
            const response = await fetch('/api/servers/import', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            })

            if (!response.ok) {
                const errorData = await response.json().catch(() => null)
                throw new Error(errorData?.message || `导入服务器数据失败 (${response.status})`)
            }

            const result = await response.json()

            // 重新加载服务器列表
            await loadServers()

            return {
                imported: result.results.imported,
                failed: result.results.failed,
                errors: result.results.errors
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : '导入服务器数据失败')
            console.error('导入服务器数据失败:', err)
            throw err
        } finally {
            setIsLoading(false)
        }
    }

    // 重置数据库
    const resetDatabase = async (): Promise<boolean> => {
        try {
            setIsLoading(true)
            setError(null)

            // 调用重置API
            const response = await fetch('/api/servers/reset', {
                method: 'POST',
            })

            if (!response.ok) {
                throw new Error(`重置数据库失败 (${response.status}: ${response.statusText})`)
            }

            // 清除选中服务器
            setSelectedServerUrl('')

            // 重新加载服务器列表
            await loadServers()

            return true
        } catch (err) {
            setError(err instanceof Error ? err.message : '重置数据库失败')
            console.error('重置数据库失败:', err)
            return false
        } finally {
            setIsLoading(false)
        }
    }

    // 提供上下文值
    const value = {
        servers,
        isLoading,
        error,
        selectedServerUrl,
        setSelectedServerUrl,
        addServer,
        deleteServer,
        loadServers,
        exportServers,
        importServers,
        resetDatabase,
        validateImportData
    }

    return (
        <ServerContext.Provider value={value}>
            {children}
        </ServerContext.Provider>
    )
}

// 自定义钩子便于使用上下文
export function useServer() {
    const context = useContext(ServerContext)

    if (context === undefined) {
        throw new Error('useServer must be used within a ServerProvider')
    }

    return context
}
