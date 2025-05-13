'use client'

import { createContext, useContext, useState,  ReactNode } from 'react'
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
    exportServers: () => Promise<void>
    importServers: (data: any[]) => Promise<{ imported: number, failed: number }>
    resetDatabase: () => Promise<boolean>
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
                throw new Error('获取服务器列表失败')
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
                throw new Error('添加服务器失败')
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
                throw new Error('删除服务器失败')
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

    // 导出服务器数据
    const exportServers = async (): Promise<void> => {
        try {
            setIsLoading(true)
            setError(null)

            // 调用导出API
            const response = await fetch('/api/servers/export')

            if (!response.ok) {
                throw new Error('导出服务器数据失败')
            }

            // 获取 blob 数据
            const blob = await response.blob()

            // 创建下载链接
            const url = window.URL.createObjectURL(blob)
            const a = document.createElement('a')
            a.href = url
            a.download = `servers_backup_${new Date().toISOString().split('T')[0]}.json`
            document.body.appendChild(a)
            a.click()
            a.remove()

            // 释放 URL 对象
            window.URL.revokeObjectURL(url)
        } catch (err) {
            setError(err instanceof Error ? err.message : '导出服务器数据失败')
            console.error('导出服务器数据失败:', err)
        } finally {
            setIsLoading(false)
        }
    }

    // 导入服务器数据
    const importServers = async (data: any[]): Promise<{ imported: number, failed: number }> => {
        try {
            setIsLoading(true)
            setError(null)

            // 调用导入API
            const response = await fetch('/api/servers/import', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            })

            if (!response.ok) {
                throw new Error('导入服务器数据失败')
            }

            const result = await response.json()

            // 重新加载服务器列表
            await loadServers()

            return {
                imported: result.results.imported,
                failed: result.results.failed
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
                throw new Error('重置数据库失败')
            }

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
        exportServers, importServers, resetDatabase
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
