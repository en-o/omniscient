'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { ServerEntity } from "@typesss/serverEntity"

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
    const [servers, setServers] = useState<ServerEntity[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedServerUrl, setSelectedServerUrl] = useState('');

    // 初始化加载服务器
    useEffect(() => {
        loadServers();
    }, []);

    // 打开数据库连接
    const openDB = (): Promise<IDBDatabase> => {
        return new Promise((resolve, reject) => {
            const request = indexedDB.open(DB_NAME, DB_VERSION);

            request.onerror = (event) => {
                const error = (event.target as IDBOpenDBRequest).error;
                reject(new Error(`数据库错误: ${error?.message || '未知错误'}`));
            };

            request.onupgradeneeded = (event) => {
                const db = (event.target as IDBOpenDBRequest).result;
                if (!db.objectStoreNames.contains(STORE_NAME)) {
                    // 创建服务器存储，使用id作为键
                    db.createObjectStore(STORE_NAME, { keyPath: 'id', autoIncrement: true });
                }
            };

            request.onsuccess = (event) => {
                const db = (event.target as IDBOpenDBRequest).result;
                resolve(db);
            };
        });
    };

    // 加载所有服务器
    const loadServers = async (): Promise<void> => {
        setIsLoading(true);
        setError(null);

        try {
            const db = await openDB();
            const transaction = db.transaction(STORE_NAME, 'readonly');
            const store = transaction.objectStore(STORE_NAME);
            const request = store.getAll();

            return new Promise((resolve, reject) => {
                request.onsuccess = () => {
                    setServers(request.result);
                    setIsLoading(false);
                    resolve();
                };

                request.onerror = () => {
                    const err = new Error('加载服务器失败');
                    setError(err.message);
                    setIsLoading(false);
                    reject(err);
                };
            });
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : '加载服务器时出错';
            setError(errorMessage);
            setIsLoading(false);
            throw err;
        }
    };

    // 添加服务器
    const addServer = async (serverData: Omit<ServerEntity, 'id'>): Promise<ServerEntity> => {
        try {
            // 验证URL格式
            try {
                new URL(serverData.url);
            } catch (e) {
                throw new Error('无效的URL格式');
            }

            // 检查URL是否重复
            const existingServer = servers.find(s => s.url === serverData.url);
            if (existingServer) {
                throw new Error('服务器URL已存在');
            }

            // 生成唯一ID
            const newServer: ServerEntity = {
                ...serverData,
                id: `server_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
            };

            const db = await openDB();
            const transaction = db.transaction(STORE_NAME, 'readwrite');
            const store = transaction.objectStore(STORE_NAME);

            return new Promise((resolve, reject) => {
                const request = store.add(newServer);

                request.onsuccess = () => {
                    // 更新本地状态
                    setServers(prev => [...prev, newServer]);
                    resolve(newServer);
                };

                request.onerror = () => {
                    reject(new Error('添加服务器失败'));
                };
            });
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : '添加服务器时出错';
            setError(errorMessage);
            throw err;
        }
    };

    // 删除服务器
    const deleteServer = async (id: string): Promise<boolean> => {
        try {
            const db = await openDB();
            const transaction = db.transaction(STORE_NAME, 'readwrite');
            const store = transaction.objectStore(STORE_NAME);

            return new Promise((resolve, reject) => {
                const request = store.delete(id);

                request.onsuccess = () => {
                    // 更新本地状态
                    setServers(prev => prev.filter(server => server.id !== id));

                    // 如果删除了当前选中的服务器，清除选中状态
                    const deletedServer = servers.find(s => s.id === id);
                    if (deletedServer && deletedServer.url === selectedServerUrl) {
                        setSelectedServerUrl('');
                    }

                    resolve(true);
                };

                request.onerror = () => {
                    reject(new Error('删除服务器失败'));
                };
            });
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : '删除服务器时出错';
            setError(errorMessage);
            throw err;
        }
    };

    // 导出服务器数据
    const exportServers = async (): Promise<void> => {
        try {
            if (servers.length === 0) {
                throw new Error('没有可导出的服务器数据');
            }

            // 创建一个包含导出数据和元数据的对象
            const exportData = {
                version: DB_VERSION,
                timestamp: new Date().toISOString(),
                data: servers
            };

            // 将数据转换为JSON字符串
            const dataStr = JSON.stringify(exportData, null, 2);

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

    // 导入服务器数据
    const importServers = async (importData: any): Promise<{ imported: number; failed: number }> => {
        try {
            // 验证导入数据格式
            if (!importData || !Array.isArray(importData.data)) {
                throw new Error('无效的导入数据格式');
            }

            const serversToImport = importData.data;
            const result = { imported: 0, failed: 0 };

            const db = await openDB();
            const transaction = db.transaction(STORE_NAME, 'readwrite');
            const store = transaction.objectStore(STORE_NAME);

            // 获取当前所有URL，用于检查重复
            const existingUrls = servers.map(s => s.url);

            // 处理每个服务器
            const importPromises = serversToImport.map(async (server: any) => {
                // 基本验证
                if (!server.url || !server.description || !server.id) {
                    result.failed++;
                    return;
                }

                // 检查URL是否重复
                if (existingUrls.includes(server.url)) {
                    result.failed++;
                    return;
                }

                // 添加到数据库
                try {
                    return new Promise<void>((resolve, reject) => {
                        const request = store.add(server);

                        request.onsuccess = () => {
                            result.imported++;
                            existingUrls.push(server.url); // 更新URL列表防止重复导入
                            resolve();
                        };

                        request.onerror = () => {
                            result.failed++;
                            resolve(); // 继续处理其他记录
                        };
                    });
                } catch (e) {
                    result.failed++;
                }
            });

            // 等待所有导入操作完成
            await Promise.all(importPromises);

            // 重新加载服务器列表以反映导入结果
            await loadServers();

            return result;
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : '导入服务器数据时出错';
            setError(errorMessage);
            throw err;
        }
    };

    // 重置数据库
    const resetDatabase = async (): Promise<boolean> => {
        try {
            const db = await openDB();
            const transaction = db.transaction(STORE_NAME, 'readwrite');
            const store = transaction.objectStore(STORE_NAME);

            return new Promise((resolve, reject) => {
                const request = store.clear();

                request.onsuccess = async () => {
                    // 清空本地状态
                    setServers([]);
                    setSelectedServerUrl('');

                    // 重新加载空数据
                    await loadServers();
                    resolve(true);
                };

                request.onerror = () => {
                    reject(new Error('重置数据库失败'));
                };
            });
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
