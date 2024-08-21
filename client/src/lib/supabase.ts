import { createClient, SupabaseClient } from "@supabase/supabase-js"

interface Meta {
    supabaseUrl: string
    supabaseAnonKey: string
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

export async function getSupabaseClient(): Promise<SupabaseClient<any, "public", any>> {
    const meta = await fetchMeta()
    const supabaseUrl = meta.supabaseUrl
    const supabaseAnonKey = meta.supabaseAnonKey

    return createClient(supabaseUrl, supabaseAnonKey)   
}