import Image from "next/image";
import ServerFrame from "@components/ServerFrame"


export default function Home() {

    return (
        // 调整网格布局以更好地适应 iframe 内容
        <div
            className="grid grid-rows-[auto_1fr_auto] items-center justify-items-center min-h-screen p-8 pb-20 gap-8 sm:p-20 font-[family-name:var(--font-geist-sans)]"> {/* Adjusted grid-rows and gap */}
            <main className="flex flex-col gap-4 row-start-2 w-full h-full">
                <ServerFrame />
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
