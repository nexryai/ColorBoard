<script lang="ts">
    import type { PageData } from "../../../../../.svelte-kit/types/src/routes/galleries/[id]/$types"
    export let data: PageData
    const galleryId = data.id
</script>
<div>
    <p>{galleryId}</p>
</div>