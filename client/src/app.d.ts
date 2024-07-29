import type {
	Terminal
} from "@battlefieldduck/xterm-svelte"

// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}

	var wasmInitialized: boolean
	var terminals: Terminal[]
	function updateElement(id: string, text: string): void

	function onTerminalInput(index: number, event: string): void
	function registerTerminal(index: number, sessionId: string, sshHostId: string): void

	function startSSH(index: number): void

}

export {}