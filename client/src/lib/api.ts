import { browser } from "$app/environment"

export function callApi<T>(method: string, url: string, data?: any): Promise<T> {
    console.log("API Called:", method, url)
    return fetch(url, {
        method,
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    }).then((res) => {
        if (!res.ok) {
            if (res.status === 401) {
                return
            }

            throw new Error(res.statusText)
        }
        return res.json()
    })
}