import {ReactNode} from "react";

export default function Card(props: {children: ReactNode, height?: string, width?: string}) {
    const height: string = (props.height) ? props.height : "h-auto"
    const width: string = (props.width) ? props.width : "rounded-w-100vh"
    const divContentStyle: string = `p-2 ${width} ${height} overflow-auto px-20 py-15`;

    return (
        <div data-testid="card"
             className="
             flex
             justify-center
             items-center
             shadow-[0_3px_10px_-1px_rgba(0,0,0,1)]
             rounded-xl py-2">
            <div className="flex flex-col">
                <div className={divContentStyle}>
                    {props.children}
                </div>
            </div>
        </div>
    )
}