export default function StyledInput(props: {
    id?: string,
    type?: string,
    name?: string,
    placeholder?: string,
    required: boolean
}) {
    const id: string = (props.id) ? props.id : '';
    const type: string = (props.type) ? props.type : 'text';
    const name: string = (props.name) ? props.name : '';
    const placeholder: string = (props.placeholder) ? props.placeholder : '';

    if (props.required) {
        return (
            <div className="relative">
                <input
                    className="
                bg-neutral-700
                peer
                block
                w-full
                h-10
                px-3
                border
                border-neutral-500
                rounded-md
                outline-2
                placeholder:text-neutral-500"
                    id={id}
                    type={type}
                    name={name}
                    placeholder={placeholder}
                    required
                />
            </div>
        )
    }
    return (
        <input
            className="
                bg-neutral-700
                peer
                block
                w-full
                h-10
                px-3
                border
                border-neutral-500
                rounded-md
                outline-2
                placeholder:text-neutral-500"
            id={id}
            type={type}
            name={name}
            placeholder={placeholder}
        />
    )
}