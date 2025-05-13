import Image from "next/image";
import { useServer } from "@components/ServerContext"



export default function Home() {
    const { selectedServerUrl } = useServer()

    return (
        // 调整网格布局以更好地适应 iframe 内容
        <div
            className="grid grid-rows-[auto_1fr_auto] items-center justify-items-center min-h-screen p-8 pb-20 gap-8 sm:p-20 font-[family-name:var(--font-geist-sans)]"> {/* Adjusted grid-rows and gap */}
            {/* main 标签在第二行，占据剩余空间 */}
            <main
                className="flex flex-col gap-4 row-start-2 w-full h-full"> {/* flex-col for stacked content, w-full h-full to fill grid cell */}
                {/* 根据是否有 serverUrl 来条件渲染 iframe 或提示 */}
                {selectedServerUrl ? (
                    <iframe
                        // 使用传入的 serverUrl 作为 iframe 的 src
                        src={selectedServerUrl}
                        title="服务器内容"
                        // 使 iframe 填充其父容器 (main)
                        className="w-full h-full border-0 shadow-md rounded-lg" // Added styling for better appearance
                    >
                        您的浏览器不支持 iframe。
                    </iframe>
                ) : (
                    // 如果没有选中服务器，显示提示信息
                    <div
                        className="w-full h-full flex items-center justify-center text-gray-500 dark:text-gray-400 text-lg">
                        请在导航栏选择一个服务器
                    </div>
                )}
            </main>
            <footer
                className="row-start-3 flex gap-6 flex-wrap items-center justify-center mt-8"> {/* Adjusted gap and added margin-top */}
                <a
                    className="flex items-center gap-2 hover:underline hover:underline-offset-4"
                    href="https://gitee.com/etn/omniscient"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    <Image
                        aria-hidden
                        src="/globe.svg"
                        alt="Globe icon"
                        width={16}
                        height={16}
                    />
                    Go to gitee →
                </a>
            </footer>
        </div>
    );
}
