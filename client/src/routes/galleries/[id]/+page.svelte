<script lang="ts">
    import type { PageData } from "../../../../.svelte-kit/types/src/routes/galleries/[id]/$types"
    import init, { render_blurhash } from "$lib/wasm/cb_client_wasm"
    import { fetchGallery } from "$lib/api"
    import { afterUpdate } from "svelte"

    export let data: PageData
    const galleryId = data.id

    interface Placeholder {
        elementId: string
        blurhash:  string
        rendered: boolean
    }

    let placeholdersAreReady = false
    let placeholders: Placeholder[] = []

    async function initAndFetch() {
        try {
            const [_, gallery] = await Promise.all([
                init(),
                fetchGallery(galleryId),
            ])
            console.log("initialized")

            const startTime = Date.now()
            if (gallery && gallery.images.length > 0) {
                for (let i = 0; i < gallery.images.length; i++) {
                    const image = gallery.images[i]
                    // PushだとSvelte側で再描画されない
                    placeholders = [...placeholders, {
                        elementId: `ph-canvas-${i}`,
                        blurhash: image.blurhash,
                        rendered: false
                    }]
                }
            }

            placeholdersAreReady = true

            const endTime = Date.now()
            console.log(endTime - startTime)

            if (gallery) {
                console.log(gallery)
            }
        } catch (error) {
            console.error("An error occurred:", error)
        }
    }

    afterUpdate(() => {
        // placeholdersAreReadyがtrueになった後に実行される
        if (placeholdersAreReady) {
            placeholders = placeholders.map(placeholder => {
                if (!placeholder.rendered) {
                    console.log("Rendering blurhash...")
                    render_blurhash(placeholder.elementId, placeholder.blurhash);
                    return { ...placeholder, rendered: true };
                }
                return placeholder;
            });
        }
    })

    initAndFetch()
</script>

<div>
    <p>{galleryId}</p>
    <div class="grid gap-4 grid-cols-1 md:grid-cols-3 lg:grid-cols-5 w-[100%]" class:hidden={!placeholdersAreReady}>
        {#each placeholders as placeholder}
            <div class="w-[150px] h-[150px] overflow-hidden">
                <canvas
                    id={placeholder.elementId}
                    class="w-[300px] h-[150px]"
                />
            </div>
        {/each}
    </div>
</div>
