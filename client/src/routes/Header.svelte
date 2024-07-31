<script lang="ts">
    import logo from "$lib/images/hummingbird.webp"
	import AntennaBarsOff from "@tabler/icons-svelte/icons/antenna-bars-off"
	import AntennaBars_2 from "@tabler/icons-svelte/icons/antenna-bars-2"
	import AntennaBars_3 from "@tabler/icons-svelte/icons/antenna-bars-3"
	import AntennaBars_4 from "@tabler/icons-svelte/icons/antenna-bars-4"
    import AntennaBars_5 from "@tabler/icons-svelte/icons/antenna-bars-5"
	import {browser} from "$app/environment";

	enum ConnectionStatus {
		Unknown,
		Green,
		Yellow,
		Red
	}
	let status = ConnectionStatus.Unknown
	let pingIsTooHigh = false
	let isOffline = false
	let ping = 0

	// @ts-ignore
	let ws = null
	if (!browser) {
		ws = null;
	} else {
		// 現在のURLから生成
		const url = window.location.href
		let protocol = "wss://"
		if (url.startsWith("http://")) {
			console.log("Unsecure connection. Do not use in production.")
			protocol = "ws://"
		}

		const host = new URL(url).host
		ws = new WebSocket(protocol + host + "/ping")
	}

	ws?.addEventListener("open", function open() {
		console.log("connected");
		// @ts-ignore
		ws.send(Date.now());
	})

	ws?.addEventListener("error", function error() {
		status = ConnectionStatus.Red;
		ping = 999;
		isOffline = true;
	})

	ws?.addEventListener("close", function close() {
		status = ConnectionStatus.Red;
		ping = 999;
		isOffline = true;
	})

	ws?.addEventListener("message", function incoming(event) {
		const message = event.data;
		const pong = Date.now() - parseInt(message);
		if (pong <= 999) {
			ping = pong;
		} else {
			// 3桁を超える場合は999にする
			ping = 999;
		}

		if (pong < 100) {
			status = ConnectionStatus.Green;
		} else if (pong < 200) {
			status = ConnectionStatus.Yellow;
		} else {
			status = ConnectionStatus.Red;
		}

		if (pong > 300) {
			pingIsTooHigh = true;
		} else {
			pingIsTooHigh = false;
		}

		setTimeout(() => {
			// @ts-ignore
			ws.send(Date.now());
		}, 3000);
	})
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
				class:icon-gray={status === ConnectionStatus.Unknown}
				class:icon-green={status === ConnectionStatus.Green}
				class:icon-yellow={status === ConnectionStatus.Yellow}
				class:icon-red={status === ConnectionStatus.Red}
		>
			{#if status === ConnectionStatus.Unknown || isOffline}
				<AntennaBarsOff class="antenna-icon" />
			{:else if status === ConnectionStatus.Green}
				<AntennaBars_5 class="antenna-icon" />
			{:else if status === ConnectionStatus.Yellow}
				<AntennaBars_4 style="margin-top: 11px;" />
			{:else if status === ConnectionStatus.Red && !pingIsTooHigh}
				<AntennaBars_3 style="margin-top: 11px;" />
			{:else if pingIsTooHigh}
				<AntennaBars_2 style="margin-top: 11px;" />
			{/if}
            <span>{ping}ms</span>
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

	svg {
		width: 2em;
		height: 3em;
		display: block;
	}

	path {
		fill: var(--background);
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