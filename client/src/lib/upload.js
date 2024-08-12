import init, { upload_file } from "$lib/wasm/cb_client_wasm"

init().then(() => {
    console.log("initialized")
})

/**
 * @param {string} galleryId
 */
export function createUploadReader(galleryId) {
    const reader = new FileReader()

    reader.onload = (e) => {
        globalThis.isUploading = true
        // @ts-ignore
        const uploadRes = upload_file(galleryId, new Uint8Array(reader.result))
        console.log(uploadRes)
    }

    return reader
}