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

export function ServerProvider({ children }: { children: ReactNode }) {
    const [selectedServerUrl, setSelectedServerUrl] = useState('')

    return (
        <ServerContext.Provider value={{ selectedServerUrl, setSelectedServerUrl }}>
            {children}
        </ServerContext.Provider>
    )
}

export const useServer = () => useContext(ServerContext)