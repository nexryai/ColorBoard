<script lang="ts">
    import { fly } from "svelte/transition"
    import type { PageData } from "../../../../../.svelte-kit/types/src/routes/galleries/[id]/$types"
    import Dropzone from "svelte-file-dropzone"
    import LoaderCircle from "lucide-svelte/icons/loader-circle"
    import PhotoX from "@tabler/icons-svelte/icons/photo-x"
    import Check from "@tabler/icons-svelte/icons/check"
    import init, { upload_file } from "$lib/wasm/cb_client_wasm"

    init().then(() => {
        console.log("initialized")
    })

    export let data: PageData
    const galleryId = data.id

    let uploadQueue: { name: string, done: boolean }[] = []
    let uploadedFiles: { name: string }[] = []
    let failedFiles: { name: string, reason?: string }[] = []
    
    async function waitForDone(item: { done: boolean }) {
        return new Promise<void>((resolve) => {
            const checkDone = () => {
                if (item.done) {
                    resolve()
                } else {
                    // 100ミリ秒ごとにチェック
                    setTimeout(checkDone, 100)
                }
            }
            checkDone()
        })
    }

    function createUploadReader(fileIndex: number) {
        const reader = new FileReader()

        reader.onload = async (e) => {
            globalThis.isUploading = true
            const filename = uploadQueue[fileIndex].name
            let uploadRes: number

            try {
                // @ts-ignore
                uploadRes = await upload_file(galleryId, new Uint8Array(reader.result))
                console.log(`Upload resp: ${uploadRes}`)
            } catch(e) {
                failedFiles = [...failedFiles, {
                    name: filename,
                    reason: `Failed to process image: ${e}`
                }]

                return
            } finally {
                uploadQueue[fileIndex].done = true
            }

            if (uploadRes != 200) {
                if (uploadRes == 1) {
                    failedFiles = [...failedFiles, {
                        name: filename,
                        reason: "Failed to upload image."
                    }]
                }else if (uploadRes == 409){
                    failedFiles = [...failedFiles, {
                        name: filename,
                        reason: "The same file already exists."
                    }]
                } else {
                    failedFiles = [...failedFiles, {
                        name: filename,
                        reason: `Server response code was not 200 (${uploadRes})`
                    }]
                }
            } else {
                uploadedFiles = [...uploadedFiles, {
                    name: filename
                }]
            }
        }

        return reader
    }

    async function handleFilesSelect(e: any) {
        const { acceptedFiles, fileRejections } = e.detail
        acceptedFiles.done = false

        // Add to queue
        uploadQueue = [...acceptedFiles]
        failedFiles = [...failedFiles, ...fileRejections]

        console.log(uploadQueue)
        for (let i = 0; i < uploadQueue.length; i++) {
            const reader = createUploadReader(i)
            console.log(reader.readAsArrayBuffer(acceptedFiles[i]))

            await waitForDone(uploadQueue[i])
            console.log(`Item ${i} is done.`)
        }

        // リセットする
        uploadQueue = []
    }
</script>

<div>
    <p>Queue</p>
    <div class="mt-8 mb-8">
        <Dropzone on:drop={handleFilesSelect} />
    </div>

    {#each uploadQueue as item}
        {#if !item.done}
            <div class="flex m-4 p-4 border border-gray-300 rounded-[10px]" in:fly="{{y: 70}}" out:fly="{{x: 70}}">
                <LoaderCircle class="mt-[4px] mr-2 h-4 w-4 animate-spin" />
                {item.name}
            </div>
        {/if}
    {/each}

    <p class="mt-8">Errors</p>
    {#each failedFiles as failed}
        <div class="text-orange-700 flex m-4 p-4 border border-orange-700 rounded-[10px]" in:fly="{{y: 70}}" out:fly="{{x: 70}}">
            <PhotoX class="mt-[4px] mr-2 h-4 w-4" />
            {failed.name}: {failed.reason || "Failed to process image"}
        </div>
    {/each}

    <p class="mt-8">Uploaded files</p>
    {#each uploadedFiles as upload}
        <div class="text-green-700 flex m-4 p-4 border border-green-700 rounded-[10px]" in:fly="{{y: 70}}" out:fly="{{x: 70}}">
            <Check class="mt-[4px] mr-2 h-4 w-4" />
            {upload.name}
        </div>
    {/each}
</div>
