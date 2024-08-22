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
    return {
        apiKey: "AIzaSyDbF9fLPkYaNh3A84vw74vi-VIacZN7FXo",
        authDomain: "colorboard-lab.firebaseapp.com",
        projectId: "colorboard-lab",
        storageBucket: "colorboard-lab.appspot.com",
        messagingSenderId: "367082909949",
        appId: "1:367082909949:web:c4b3d25649ed717f5f0e94"
    }
}

