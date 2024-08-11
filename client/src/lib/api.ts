import { browser } from "$app/environment"

export interface Gallery {
    id: string
    name: string
    createdAt: string
    updatedAt: string
    isPublic: boolean
    userId: string
}

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

export function getMyGalleries(): Promise<Gallery[]> {
    return callApi<Gallery[]>("GET", "/api/gallery/list")
}