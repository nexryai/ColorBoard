<script lang="ts">
    import { afterUpdate, onMount } from "svelte"
    import type { PageData } from "../../../../.svelte-kit/types/src/routes/galleries/[id]/$types"
    import init, { render_blurhash } from "$lib/wasm/cb_client_wasm"
    import { fetchGallery } from "$lib/api"
    import { Button } from "$lib/components/ui/button"
    import { Skeleton } from "$lib/components/ui/skeleton"
    import CloudUpload from "@tabler/icons-svelte/icons/cloud-upload"

    import PhotoSwipeLightbox from "photoswipe/lightbox"
    import "photoswipe/style.css"

    let isLoading = true
    let galleryName = "loading..."

    export let data: PageData
    const galleryId = data.id

    interface Placeholder {
        elementId:    string
        blurhash:     string
        // Blurhashが描画されたか
        rendered:     boolean
        thumbnailUrl: string
        imageUrl:     string
        // サムネイルが読み込まれたか
        loaded:       boolean
        width:       number
        height:      number
    }

    let placeholdersAreReady = false
    let placeholders: Placeholder[] = []

    async function initAndFetch() {
        try {
            const [_, gallery] = await Promise.all([
                init(),
                fetchGallery(galleryId),
            ])

            galleryName = gallery.name
            console.log("initialized")
            isLoading = false

            if (gallery && gallery.images.length > 0) {
                for (let i = 0; i < gallery.images.length; i++) {
                    const image = gallery.images[i]
                    // PushだとSvelte側で再描画されない
                    placeholders = [...placeholders, {
                        elementId: `ph-canvas-${i}`,
                        blurhash: image.blurhash,
                        rendered: false,
                        thumbnailUrl: `/api/files/${image.thumbnailKey}`,
                        imageUrl: `/api/files/${image.storageKey}`,
                        loaded: false,
                        width: image.width,
                        height: image.height
                    }]
                }
            }

            placeholdersAreReady = true

            if (gallery) {
                console.log(gallery)
            }
        } catch (error) {
            console.error("An error occurred:", error)
        }
    }

    function handleImageLoad(index: number) {
        placeholders = placeholders.map((placeholder, i) => {
            if (i === index) {
                return {
                    ...placeholder, 
                    loaded: true
                }
            }
            return placeholder
        })
    }

    afterUpdate(() => {
        // placeholdersAreReadyがtrueになった後に実行される
        if (placeholdersAreReady) {
            placeholders = placeholders.map(placeholder => {
                if (!placeholder.rendered && !placeholder.loaded) {
                    console.log("Rendering blurhash...")
                    render_blurhash(placeholder.elementId, placeholder.blurhash)
                    return { ...placeholder, rendered: true }
                }
                return placeholder
            })
        }
    })

    onMount(() => {
        let lightbox = new PhotoSwipeLightbox({
            gallery: '#' + galleryId,
            children: 'a',
            pswpModule: () => import("photoswipe"),
        })
        lightbox.init()
    })

    initAndFetch()
</script>

<div>
    {#if isLoading}
        <div class="skeleton space-y-2 mt-12">
            <Skeleton class="h-4 w-[250px]" />
            <Skeleton class="h-4 w-[200px]" />
        </div>
    {:else}
    <div class="flex justify-between items-center mb-8">
        <p class="text-2xl">{galleryName}</p>
        <Button variant="outline" href="/galleries/{galleryId}/upload">
            <CloudUpload class="mr-2" />
            Upload
        </Button>
    </div>
    <div id={galleryId} class="grid place-items-center gap-4 grid-cols-1 md:grid-cols-3 lg:grid-cols-5 w-[100%] pswp-gallery" class:hidden={!placeholdersAreReady}>
        {#each placeholders as placeholder, index}
            <div class="w-[150px] h-[150px] overflow-hidden border border-slate-200 transition hover:shadow-md">
                <canvas
                    id={placeholder.elementId}
                    class="w-[300px] h-[150px]"
                    class:hidden={placeholder.loaded}
                />
                <a
                    href={placeholder.imageUrl}
                    data-pswp-width={placeholder.width}
                    data-pswp-height={placeholder.height}
                    target="_blank"
                    rel="noreferrer"
                >
                    <img 
                        class="h-[150px] aspect-auto object-cover"
                        src={placeholder.thumbnailUrl} 
                        on:load={() => handleImageLoad(index)}
                        class:hidden={!placeholder.loaded}
                        alt=""
                    />
                </a>
            </div>
        {/each}
    </div>
    {/if}
</div>
