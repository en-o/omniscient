'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
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

    // 首次加载时获取服务器列表
    useEffect(() => {
        loadServers()
    }, [])

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
