import init, { generate_thumbnail, upload_file } from "../../static/wasm/cb_client_wasm"

init().then(() => {
    console.log("initialized")
})

export function createUploadReader() {
    const reader = new FileReader()

    reader.onload = (e) => {
        globalThis.isUploading = true

        /*
        // @ts-ignore
        const thumb = generate_thumbnail(new Uint8Array(reader.result))
        console.log("thumbnail generated")
        const img = document.createElement("img")
        img.src = URL.createObjectURL(new Blob([thumb], { type: "image/jpeg" }))
        document.body.appendChild(img)*/

        // @ts-ignore
        const uploadRes = upload_file("f", new Uint8Array(reader.result))
        console.log(uploadRes)
    }

    return reader
}