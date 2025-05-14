'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { ServerEntity } from "@typesss/serverEntity"
import {GET} from "@/app/api/servers/export/route";

// 定义服务器上下文接口
interface ServerContextType {
    servers: ServerEntity[];
    isLoading: boolean;
    error: string | null;
    selectedServerUrl: string;
    setSelectedServerUrl: (url: string) => void;
    addServer: (server: Omit<ServerEntity, 'id'>) => Promise<ServerEntity>;
    deleteServer: (id: string) => Promise<boolean>;
    exportServers: () => Promise<void>;
    importServers: (data: any) => Promise<{ imported: number; failed: number }>;
    resetDatabase: () => Promise<boolean>;
    loadServers: () => Promise<void>;
}

// 创建上下文
const ServerContext = createContext<ServerContextType | undefined>(undefined);

// 数据库名称和对象仓库名称
const DB_NAME = 'pm-gateway';
const DB_VERSION = 1;
const STORE_NAME = 'servers';

// 导出数据文件名
const EXPORT_FILENAME = 'pm-gateway-servers.json';

// 服务器提供者Props接口
interface ServerProviderProps {
    children: ReactNode;
}

// 服务器提供者组件
export function ServerProvider({ children }: ServerProviderProps) {
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
            if (servers.length === 0) {
                throw new Error('没有可导出的服务器数据');
            }

            const response = await fetch(`/api/servers/export`, {
                method: 'GET',
            })

            // 将数据转换为JSON字符串
            const dataStr = JSON.stringify(response.json(), null, 2);

            // 创建Blob对象
            const blob = new Blob([dataStr], { type: 'application/json' });

            // 创建下载链接
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = EXPORT_FILENAME;

            // 触发下载
            document.body.appendChild(a);
            a.click();

            // 清理
            setTimeout(() => {
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
            }, 100);
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : '导出服务器数据时出错';
            setError(errorMessage);
            throw err;
        }
    };

    /**
     * 导入服务器数据
     * 根据后端API优化的导入函数
     */
    const importServers = async (importData: any[]): Promise<{ total: number; imported: number; failed: number }> => {
        try {
            // 验证导入数据格式
            if (!Array.isArray(importData)) {
                throw new Error('无效的导入数据格式');
            }

            // 发送到后端API处理导入
            const response = await fetch('/api/servers/import', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(importData),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || '导入服务器数据失败');
            }

            // 获取后端返回的导入结果
            const result = await response.json();

            // 导入完成后重新加载服务器列表
            await loadServers();

            return {
                total: result.results.total,
                imported: result.results.imported,
                failed: result.results.failed
            };
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : '导入服务器数据时出错';
            setError(errorMessage);
            throw err;
        }
    };

    // 重置数据库
    const resetDatabase = async (): Promise<boolean> => {
        try {

            const response = await fetch(`/api/servers/reset`, {
                method: 'POST',
            });
            await loadServers();
            return response.ok
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : '重置数据库时出错';
            setError(errorMessage);
            throw err;
        }
    };

    // 提供上下文值
    const contextValue: ServerContextType = {
        servers,
        isLoading,
        error,
        selectedServerUrl,
        setSelectedServerUrl,
        addServer,
        deleteServer,
        exportServers,
        importServers,
        resetDatabase,
        loadServers
    };

    return (
        <ServerContext.Provider value={contextValue}>
            {children}
        </ServerContext.Provider>
    );
}

// 自定义钩子以使用服务器上下文
export function useServer(): ServerContextType {
    const context = useContext(ServerContext);
    if (context === undefined) {
        throw new Error('useServer must be used within a ServerProvider');
    }
    return context;
}
