import Image from "next/image";
import ServerFrame from "@components/ServerFrame"


export default function Home() {

    return (
        // 调整网格布局以更好地适应 iframe 内容
        <div className="grid grid-rows-[1fr_auto] min-h-screen p-4 gap-4">
            <main className="w-full h-full">
                <ServerFrame/>
            </main>
            <footer className="flex gap-6 flex-wrap items-center justify-center">
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
                    Go to gitee by tan →
                </a>
            </footer>
        </div>
    );
}
