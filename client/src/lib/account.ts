import { browser } from "$app/environment"
import { getSupabaseClient } from "./supabase"

export function isLoggedIn(): boolean {
    // isLoginのCookieがあればログインしていると判断
    if (!browser) {
        console.error("This function is only available in the browser")
        return false
    }
    console.log(document.cookie)
    return document.cookie.includes("auth_uid=google%") || document.cookie.includes("auth_uid=azuread%")
}
export async function refreshSession() {
    if (!browser) {
        console.error("This function is only available in the browser")
        return
    }
    
    if (isLoggedIn()) {
        const supabase = await getSupabaseClient()
        const { data: { session }, error: error } = await supabase.auth.getSession()
        if (error) {
            console.log(error)
            throw new Error("Failed to refresh session")
        }
        
        console.log(session)
    }

}