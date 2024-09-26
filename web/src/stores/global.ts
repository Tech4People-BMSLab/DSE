import { writable } from 'svelte/store'

export default writable({
    environment: null,
    language   : 'en',      // 'en' | 'nl'
    page       : 'consent',

    year: new Date().getFullYear(),
})
