<script lang="ts">
    import "../app.css"
    import "./styles.css"
    import Header from "./Header.svelte"

    import { fly } from "svelte/transition"
    import { cubicIn, cubicOut } from "svelte/easing"

    // @ts-ignore
    import NProgress from "nprogress"
    import "nprogress/nprogress.css"
    import { beforeNavigate, afterNavigate } from "$app/navigation"

    export let data

    let isLoading = false

    NProgress.configure({
        showSpinner: false // スピナーを表示しない
    })

    beforeNavigate(() => (isLoading = true))
    afterNavigate(() => (isLoading = false))

    $: {
        if (isLoading) {
            NProgress.start()
        } else {
            NProgress.done()
        }
    }

</script>

<svelte:head>
    <title>ColorBoard</title>
    <meta name="description" content="Svelte demo app" />
    <link rel="preconnect" href="https://fonts.bunny.net">
    <link rel="preconnect" href="https://cdn.jsdelivr.net">
    <link rel="stylesheet" href="https://fonts.bunny.net/css?family=noto-sans-jp:400|noto-sans-sc:400|noto-sans-kr:400|outfit:400|poppins:400|ubuntu:400">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/misans@4.0.0/lib/Latin/MiSansLatin-Medium.min.css">
</svelte:head>

<div class="app">
	<Header />

	<main>
        {#key data.pathname}
            <div
                in:fly={{ easing: cubicOut, y: 10, duration: 200, delay: 300 }}
                out:fly={{ easing: cubicIn, y: -10, duration: 200 }}
            >
		        <slot></slot>
            </div>
        {/key}
	</main>
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	main {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 1rem;
		width: 100%;
        max-width: 85vw;
		margin: 70px auto 0 auto;
		box-sizing: border-box;
	}
</style>
