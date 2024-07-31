<script lang="ts">
    import {
        Button,
        buttonVariants
    } from "$lib/components/ui/button/index.js"
    import * as Dialog from "$lib/components/ui/dialog/index.js"
    import { Input } from "$lib/components/ui/input/index.js"
    import { Label } from "$lib/components/ui/label/index.js"
    import { Checkbox } from "$lib/components/ui/checkbox/index.js"
    import { toast } from "svelte-sonner"
    import IconPlus from "@tabler/icons-svelte/icons/plus"
    import IconAlertTriangle from "@tabler/icons-svelte/icons/alert-triangle"
    import IconLoader2 from "@tabler/icons-svelte/icons/loader-2"
    import { callApi } from "$lib/api"

    let open = false
    let isLoading = false
    let showErrorMessage = false
    let errorMessage = ""
    let galleryName = ""
    let isPublic = false

    const close = () => {
        open = false
        galleryName = ""
        isPublic = false
        showErrorMessage = false
        errorMessage = ""
    }

    const addFeed = () => {
        isLoading = true

        if (!galleryName) {
            errorMessage = "Please fill in all fields."
            showErrorMessage = true
            isLoading = false
            return
        }

        callApi(
            "post",
            "/api/gallery/add",
            {
                name: galleryName,
                isPublic: isPublic
            }
        ).catch((error) => {
            console.error(error)
            errorMessage = "Failed to add feed: " + error
            showErrorMessage = true
        }).then(() => {
            open = false
            toast.success("Feed added successfully!", {
                description:"Please reload the page to see the changes."}
            )
        }).finally(() => {
            isLoading = false
        })

        return null
    }
</script>

<div>
    <Dialog.Root bind:open>
        <Dialog.Trigger>
            <Button size="icon" class="w-[50px] h-[50px] shadow-lg rounded-full">
                <IconPlus />
            </Button>
        </Dialog.Trigger>

        <Dialog.Content>
            <Dialog.Header>
                <Dialog.Title>Add new gallery</Dialog.Title>
                <Dialog.Description>
                    <p>
                        Enter the display name to add a new gallery.
                    </p>
                    {#if showErrorMessage}
                        <div class="error-message">
                            <div class="error-icon">
                                <IconAlertTriangle size={21}/>
                            </div>
                            <span>{errorMessage}</span>
                        </div>
                    {/if}
                </Dialog.Description>
            </Dialog.Header>
            <div class="grid gap-4 py-4">
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="gallery-name" class="text-right">
                        Display Name
                    </Label>
                    <Input
                            id="gallery-name"
                            bind:value={galleryName}
                            placeholder="Real truth News"
                            class="col-span-3"
                            required
                    />
                </div>
                <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="is-public" class="text-right">
                        Public Gallery
                    </Label>
                    <Checkbox id="is-public" bind:checked={isPublic} />
                </div>
            </div>
            <Dialog.Footer>
                <Button class="mt-[10px]" variant="default" disabled={isLoading} on:click={addFeed}>
                    {#if isLoading}
                        <IconLoader2 size={21} class="animate-spin" style="margin-right: 10px"/>
                        Adding...
                    {:else}
                        Add
                    {/if}
                </Button>
                <Button class="mt-[10px]" variant="secondary" on:click={close}>Cancel</Button>
            </Dialog.Footer>
        </Dialog.Content>
    </Dialog.Root>
</div>

<style lang="scss">
    .error-message {
        margin-top: 23px;
        color: red;

        .error-icon {
            float: left;
            margin: 0 6px 0 4px;
        }
    }
</style>