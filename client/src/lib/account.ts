import { browser } from "$app/environment"

export function isLoggedIn(): boolean {
    // isLoginのCookieがあればログインしていると判断
    if (!browser) {
        console.error("This function is only available in the browser")
        return false
    }
    console.log(document.cookie)
    return document.cookie.includes("auth_uid=google%") || document.cookie.includes("auth_uid=azuread%")
}
export function refreshSession() {
    if (!browser) {
        console.error("This function is only available in the browser")
        return
    }
    
    if (document.cookie.includes("auth_uid=google%")) {
        location.href = "/auth/google"
    } else if (document.cookie.includes("auth_uid=azuread%")) {
        location.href = "/auth/azuread"
    }

}