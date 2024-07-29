<script lang="ts">
    import { Xterm, XtermAddon } from "@battlefieldduck/xterm-svelte"
    import "@fontsource/jetbrains-mono/400.css"
    import type {
        ITerminalOptions,
        ITerminalInitOnlyOptions,
        Terminal
    } from "@battlefieldduck/xterm-svelte"

    let options: ITerminalOptions & ITerminalInitOnlyOptions = {
        fontFamily: '"JetBrains Mono", Monaco, Menlo, Consolas, "Courier New", monospace',
    }

    if (!globalThis.terminals) {
        globalThis.terminals = []
    }

    let terminalIdx = -1
    const sessionId = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)

    async function onLoad(event: CustomEvent<{ terminal: Terminal }>) {
        console.log('Child component has loaded')
        const terminal = event.detail.terminal

        // FitAddon Usage
        const fitAddon = new (await XtermAddon.FitAddon()).FitAddon()
        terminal.loadAddon(fitAddon)
        fitAddon.fit()

        // Register terminal
        while (!globalThis.wasmInitialized) {
            await new Promise(resolve => setTimeout(resolve, 100))
        }

        terminalIdx = globalThis.terminals.length
        globalThis.terminals[terminalIdx] = terminal
        globalThis.registerTerminal(terminalIdx, sessionId, "todo")

        // Start SSH session
        globalThis.startSSH(terminalIdx)
    }

    function onData(event: CustomEvent<string>) {
        const data = event.detail;
        // console.log('onData()', data)
        globalThis.onTerminalInput(terminalIdx, data)
    }
</script>

<svelte:head>
    <script src="/wasm/term.js" />
</svelte:head>

<div>
    <p id={`message-${sessionId}`} style="display: none" />
    <div class="terminal">
        <Xterm {options} on:load={onLoad} on:data={onData} />
    </div>
</div>

<style>
    .terminal {
        border: solid black 12px;
        border-radius: 10px;
    }
</style>