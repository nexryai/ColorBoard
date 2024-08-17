<script lang="ts">
    import Gallery from "$lib/components/Gallery.svelte"
    import type { PageData } from "../../../../.svelte-kit/types/src/routes/galleries/[id]/$types"
    import InfiniteScroll from "svelte-infinite-scroll"

    export let data: PageData
    const galleryId = data.id
    let pageIndex = 1
    let galleries = [{ id: galleryId, page: pageIndex }]

    function loadMoreGalleries() {
        console.log("Loading more!")
        pageIndex += 1
        galleries = [...galleries, { id: galleryId, page: pageIndex }]
    }
</script>

<div>
    {#each galleries as { id, page } (page)}
        <Gallery galleryId={id} page={page} />
    {/each}

    <InfiniteScroll window on:loadMore={loadMoreGalleries} />
</div>
