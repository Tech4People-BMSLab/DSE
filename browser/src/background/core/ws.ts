// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------
import EventEmitter from 'eventemitter3'

import { pack, unpack } from 'msgpackr'
import { gzipSync }     from 'fflate'

import { store }      from '@/background/core/storage'
import { Logger }     from '@/background/utils/logger'
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
const logger = new Logger('WS')
// ------------------------------------------------------------
// : Packet
// ------------------------------------------------------------
class Packet {
    version: string
    from   : string
    to     : string
    action : string
    data   : object
    
    constructor(version: string, from: string, to: string, action: string, data: object) {
        this.version = version
        this.from    = from
        this.to      = to
        this.action  = action
        this.data    = data
    }
}
// ------------------------------------------------------------
// : Helpers
// ------------------------------------------------------------
function gzip(data: Uint8Array): Uint8Array {
    return gzipSync(data)
}
// ------------------------------------------------------------
// : WS
// ------------------------------------------------------------
class WS {
    private socket: WebSocket
    private token  : string
    private version: string

    async init() {
        this.token   = await store.get('user.token')
        this.version = await store.get('extension.version') 
        
        this.connect()

        setInterval(async () => {
            this.send('update', await store.get())
        }, 1000)
    }

    private connect() {
        this.socket = new WebSocket(import.meta.env.BASE_WS)

        this.socket.onopen = async () => {
            // logger.info('connected')
        }

        this.socket.onmessage = async (event) => {
            // logger.info(`message: ${event.data}`)
        }

        this.socket.onclose = async (event) => {
            // logger.info('connection_closed', event)
            setTimeout(() => this.connect(), 5000) // Reconnect after 5 seconds
        }

        this.socket.onerror = async (error) => {
            // logger.error('web_socket_error:', error)
        }
    }

    async send(action: string, data: object) {
        try {
            if (this.socket.readyState === WebSocket.OPEN) {
                const packet   = new Packet(this.version, this.token, 'api', action, data)
                this.socket.send(JSON.stringify(packet))
            } else {
                // logger.info('web_socket_is_not_open_ready_state:', this.socket.readyState)
            }
        } catch (e) {
            // logger.error('error_sending_message:', e)
        }
    }
}

// ------------------------------------------------------------
// : Exports
// ------------------------------------------------------------
export const ws = new WS()
