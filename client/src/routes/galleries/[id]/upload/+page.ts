/** @type {import('./$types').PageLoad} */
import type { PageLoad } from "../../../../../.svelte-kit/types/src/routes/galleries/[id]/upload/$types"

export const prerender = false
export const ssr = false

export const load: PageLoad = ({ params }) => {
    return {
        id: params.id
    }
}

