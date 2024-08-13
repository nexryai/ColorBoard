import { browser } from "$app/environment"

export interface Image {
    id: string
    storageKey: string
    thumbnailKey: string
    blurhash: string
    createdAt: string
    updatedAt: string
    userId: string
    galleryId: string
}

export interface Gallery {
    id: string
    name: string
    createdAt: string
    updatedAt: string
    isPublic: boolean
    userId: string
    images: Image[]
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
            throw new Error(res.statusText)
        }
        return res.json()
    })
}

export function fetchMyGalleries(): Promise<Gallery[]> {
    return callApi<Gallery[]>("GET", "/api/gallery/list")
}

export async function fetchGallery(id: String): Promise<Gallery> {
    try {
        const response = await callApi("GET", `/api/gallery/${id}`)
        return response as Gallery
    } catch (error) {
        console.error("Failed to fetch gallery data:", error)
        throw error
    }
}
