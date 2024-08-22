import { type FirebaseOptions } from "firebase/app"

interface Meta {
    firebaseApiKey: string
    firebaseAuthDomain: string
    firebaseProjectId: string
    firebaseStorageBucket: string
    firebaseMEssagingSenderId: string
    firebaseAppId: string
}

async function fetchMeta(): Promise<Meta> {
    try {
        const response = await fetch("/meta")
        if (!response.ok) {
            throw new Error("Failed to fetch meta data")
        }

        const data: Meta = await response.json()
        return data
    } catch (error) {
        console.error("Error fetching meta data:", error)
        throw error
    }
}

export async function getFirebaseConfig(): Promise<FirebaseOptions> {
    const meta = await fetchMeta()

    return {
        apiKey: meta.firebaseApiKey,
        authDomain: meta.firebaseAuthDomain,
        projectId: meta.firebaseProjectId,
        storageBucket: meta.firebaseStorageBucket,
        messagingSenderId: meta.firebaseMEssagingSenderId,
        appId: meta.firebaseAppId
    }
}

