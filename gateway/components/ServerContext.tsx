'use client'

import { createContext, useContext, useState, ReactNode } from 'react'

interface ServerContextType {
    selectedServerUrl: string
    setSelectedServerUrl: (url: string) => void
}

const ServerContext = createContext<ServerContextType>({
    selectedServerUrl: '',
    setSelectedServerUrl: () => {}
})

// https://zh-hans.react.dev/learn/passing-data-deeply-with-context#step-3-provide-the-context
export function ServerProvider({ children }: { children: ReactNode }) {
    const [selectedServerUrl, setSelectedServerUrl] = useState('')

    return (
        <ServerContext.Provider value={{ selectedServerUrl, setSelectedServerUrl }}>
            {children}
        </ServerContext.Provider>
    )
}

export const useServer = () => useContext(ServerContext)
