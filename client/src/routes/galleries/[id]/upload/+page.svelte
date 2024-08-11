<script lang="ts">
    import { onMount } from "svelte";
    import type { PageData } from "../../../../../.svelte-kit/types/src/routes/galleries/[id]/$types";
    import Dropzone from "svelte-file-dropzone";

    import { createUploadReader } from "$lib/upload"

    export let data: PageData;
    const galleryId = data.id;

    let uploadQueue: { name: string }[] = [];
    let failedFiles: { name: string }[] = [];

    function handleFilesSelect(e: any) {
        const { acceptedFiles, fileRejections } = e.detail;
        console.log(e.detail);
        uploadQueue = [...uploadQueue, ...acceptedFiles];
        failedFiles = [...failedFiles, ...fileRejections];

        console.log(uploadQueue);
        for (let i = 0; i < uploadQueue.length; i++) {
            const reader = createUploadReader()
            console.log(reader.readAsArrayBuffer(acceptedFiles[i]));
        }
    }

</script>

<div>
    <div class="mt-8 mb-8">
        <Dropzone on:drop={handleFilesSelect} />
    </div>

    {#each uploadQueue as item}
        <li>{item}</li>
    {/each}
</div>
