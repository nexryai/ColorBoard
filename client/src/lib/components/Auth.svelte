<script lang="ts">
    import { Button } from "$lib/components/ui/button"
    import { Label } from "$lib/components/ui/label"
    import { Input } from "$lib/components/ui/input"
    import { getSupabaseClient } from "$lib/supabase"

    let loading = false
    let email = ""
    let password = ""

    const handleLogin = async () => {
        try {
            loading = true
            const supabase = await getSupabaseClient()
            const { data, error } = await supabase.auth.signInWithPassword({ email, password })
            if (error) throw error
        } catch (error) {
            if (error instanceof Error) {
                alert(error.message)
            }
        } finally {
            loading = false
        }
    }
</script>

<div>
    <div aria-live="polite">
        <form on:submit|preventDefault={handleLogin}>
            <div>
                <Label for="email">Email</Label>
                <Input
                    id="email"
                    class="inputField"
                    type="email"
                    placeholder="Your email"
                    bind:value={email}
                />
            </div>
            <div class="mt-[6px]">
                <Label for="password">Password</Label>
                <Input
                    id="password"
                    class="inputField"
                    type="password"
                    placeholder="Your strong password"
                    bind:value={password}
                />
            </div>
            <div class="mt-[19px]">
                <Button
                    type="submit"
                    class="button block m-auto w-[60%]"
                    aria-live="polite"
                    disabled={loading}
                >
                    <span>{loading ? "Loading" : "Login"}</span>
                </Button>
            </div>
        </form>
    </div>
</div>
