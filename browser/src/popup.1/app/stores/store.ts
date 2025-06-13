// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------
import      { connect }        from 'nats.ws'
import type { NatsConnection } from 'nats.ws'

import      { writable } from 'svelte/store'
import type { Writable } from 'svelte/store'

// ------------------------------------------------------------
// : IPC
// ------------------------------------------------------------
async function init() {
    ipc.set(await connect({
        servers: [import.meta.env.NATS_URL],
        user   : import.meta.env.NATS_USER,
        pass   : import.meta.env.NATS_PASS,
    }))

    connected.set(true)
}

// ------------------------------------------------------------
// : Stores
// ------------------------------------------------------------
export const ipc       = writable(null)  as Writable<NatsConnection | null>
export const connected = writable(false) as Writable<boolean>

export const page          = writable('consent') as Writable<string>
export const process_state = writable(0) as Writable<number> // 0: idle, 1: processing

export const error_message   = writable(null) as Writable<string | null>
export const success_message = writable(null) as Writable<string | null>

export const lang = writable('nl') as Writable<string>

// ------------------------------------------------------------
// : Init
// ------------------------------------------------------------
init()

