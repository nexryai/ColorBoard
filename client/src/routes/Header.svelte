<script lang="ts">
	import { onMount } from "svelte"
    import logo from "$lib/images/logo.webp"
	import { refreshSession } from "$lib/account"
	import CloudCheck from "@tabler/icons-svelte/icons/cloud-check"
	import CloudUp from "@tabler/icons-svelte/icons/cloud-up"
	import CloudQuestion from "@tabler/icons-svelte/icons/cloud-question"
	import CloudExclamation from "@tabler/icons-svelte/icons/cloud-exclamation"
    import CloudX from "@tabler/icons-svelte/icons/cloud-x"


	interface StorageStatusAPI {
		total: number
		free: number
		used: number
	}

	enum StorageStatus {
		Unknown,
		Green,
		Yellow,
		Red
	}

	let status = StorageStatus.Unknown
	let isUploading = globalThis.isUploading || false
	let usedPercentage = 0

    onMount(() => {
        const updateUploadingStatus = () => {
            isUploading = globalThis.isUploading
        }

        // グローバル変数の変更を監視
        window.addEventListener("isUploadingChange", updateUploadingStatus)

        return () => {
            window.removeEventListener("isUploadingChange", updateUploadingStatus)
        }
    });

	// Get storage status from /api/system/storage-status
	const getStorageStatus = async () => {
		const response = await fetch("/api/system/storage-status")
		if (!response.ok) {
			if (response.status === 401) {
				refreshSession()
			}

			status = StorageStatus.Unknown
			return
		}

		const data = await response.json()
		const storageStatus = data as StorageStatusAPI

		if (storageStatus.total === 0) {
			status = StorageStatus.Unknown
		} else {
			usedPercentage = Math.round((storageStatus.used / storageStatus.total) * 100)

			if (usedPercentage < 70) {
				status = StorageStatus.Green
			} else if (usedPercentage < 90) {
				status = StorageStatus.Yellow
			} else {
				status = StorageStatus.Red
			}
		}
	}

	getStorageStatus()

	// Every 30 seconds
	setInterval(() => {
		getStorageStatus()
	}, 30000)
	
</script>

<header>
	<div class="corner">
		<a href="/">
			<img src={logo} alt="SvelteKit" />
            <span class="logo-text">ColorBoard</span>
            <span class="logo-text logo-beta">beta</span>
		</a>
	</div>

	<div class="corner">
        <div
				class="header-icons"
				class:icon-gray={status === StorageStatus.Unknown}
				class:icon-green={status === StorageStatus.Green}
				class:icon-yellow={status === StorageStatus.Yellow}
				class:icon-red={status === StorageStatus.Red}
		>
			{#if isUploading}
				<CloudUp class="mt-[11px] mr-1" />
			{:else if status === StorageStatus.Green}
				<CloudCheck class="mt-[11px] mr-1" />
			{:else if status === StorageStatus.Unknown}
				<CloudQuestion class="mt-[11px] mr-1" />
			{:else if status === StorageStatus.Yellow}
				<CloudExclamation class="mt-[11px] mr-1" />
			{:else if status === StorageStatus.Red}
				<CloudX class="mt-[11px] mr-1" />
			{/if}
            <span>{usedPercentage}%</span>
        </div>
	</div>
</header>

<style lang="scss">
	header {
		display: flex;
        justify-content: space-between;
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        z-index: 49;
        padding: 4px;
        background: #fffc;
        backdrop-filter: blur(128px);
        border-bottom: #d9d9d9 solid 1px;
	}

	.corner {
		width: 5.3em;
		height: 3em;
	}

    .logo-text {
        margin-left: 4px;
        color: #353535;
        font-family: "Outfit", sans-serif;
        font-size: 20px;
    }

    .logo-beta {
        background: #2d77fd;
        color: white;
        padding: 0 4px 0 4px;
        border-radius: 4px;
        font-family: "Ubuntu", sans-serif;
        font-size: 14px;
    }

    .header-icons {
		display: flex;
        font-size: 12px;
        color: #444444;

		span {
			margin-top: 14px;
			margin-right: 20px;
		}
    }

	.corner a {
		display: flex;
		align-items: center;
		width: 100%;
		height: 100%;
	}

    .corner a:hover {
        text-decoration: none;
    }

	.corner img {
		width: 2em;
		height: 2em;
        margin-left: 15px;
		object-fit: contain;
	}

	a:hover {
		color: var(--color-theme-1);
	}

	.icon-green {
		color: #06b000;
	}

	.icon-yellow {
		color: #ee9c00;
	}

	.icon-red {
		color: #d50000;
	}

	.icon-gray {
		color: #727272;
	}
</style>