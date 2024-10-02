import { browser } from "$app/environment"

interface RegisterSessionReq {
    token: string
}

export function isLoggedIn(): boolean {
    // isLoginのCookieがあればログインしていると判断
    if (!browser) {
        console.error("This function is only available in the browser")
        return false
    }
    console.log(document.cookie)
    return document.cookie.includes("isLoggedIn=1")
}

export async function refreshSession() {
    if (!browser) {
        console.error("This function is only available in the browser")
        return
    }

    if (isLoggedIn()) {
        // ToDo
    }
}

export async function registerSession(token: string) {
    const url = "/auth/register-session"

    const payload: RegisterSessionReq = {
        token,
    }

    try {
        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        })

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`)
        }

        const responseData = await response.json()
        console.log("Response:", responseData)
    } catch (error) {
        console.error("Error:", error)
    }
}
