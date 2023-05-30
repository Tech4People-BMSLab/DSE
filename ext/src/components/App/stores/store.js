import { writable } from 'svelte/store'

export const page          = writable('consent')
export const process_state = writable(0) // 0: idle, 1: processing

export const error_message   = writable(null)
export const success_message = writable(null)

export const lang = writable('nl')