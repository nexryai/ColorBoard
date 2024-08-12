<script lang="ts">
    import type { PageData } from "../../../../.svelte-kit/types/src/routes/galleries/[id]/$types";
    import init, { render_blurhash } from "$lib/wasm/cb_client_wasm"

    export let data: PageData;
    const galleryId = data.id;

    interface Placeholder {
        elementId: string;
        blurhash: string;
    }

    let placeholders: Placeholder[] = [];

    for (let i = 0; i < 30; i++) {
        placeholders.push({
            elementId: `canvas-${i + 1}`,
            blurhash: "LuNIK4?DI;aL~9o{NHwMt7Seofay",
        });
    }

    init().then(() => {
        console.log("initialized");
        const startTime = Date.now();
        for (let i = 0; i < 30; i++) {
            render_blurhash(`canvas-${i + 1}`, "LuNIK4?DI;aL~9o{NHwMt7Seofay");
        }

        const endTime = Date.now();
        console.log(endTime - startTime);
    });
</script>

<div>
    <p>{galleryId}</p>
    <div class="grid gap-4 grid-cols-1 md:grid-cols-3 lg:grid-cols-5 w-[100%]">
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
