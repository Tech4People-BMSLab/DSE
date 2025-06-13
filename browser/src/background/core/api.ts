// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------
import browser from 'webextension-polyfill'

// General
import axios        from 'axios'
import EventEmitter from 'eventemitter3'

import { pack, unpack } from 'msgpackr'
import { gzipSync }     from 'fflate'


import semver from 'semver'

// Local
import { get } from '@/background/utils/utils'

import { store }      from '@/background/core/storage'
import { crawler }    from '@/background/core/crawler'
import { Logger }     from '@/background/utils/logger'
import { wait_until } from '@/background/utils/utils'
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
const logger = new Logger('API')

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
// : API
// ------------------------------------------------------------
class API extends EventEmitter {
    private readonly BASE_URL

    private event
    private token
    private version

    private log_avoidance = false

    constructor() {
        super()
        this.BASE_URL = import.meta.env.BASE_URL
    }

    public async init() {
        await wait_until(async () => await store.get('user.token'))

        this.token   = await store.get('user.token')
        this.version = await store.get('extension.version') 

        setInterval(async () => {
            await this.send('update', await store.get())
        }, 5_000)

        store.on('update', async (key, prev, next) => {
            await this.send('update', await store.get())
        })

        this.on('update', async (data) => {
            await this.send('update', await store.get())
        })

        this.on('connected', async (info) => {
            logger.info('Connected', info)
        })

        this.on('reload', async () => {
            browser.runtime.reload()
        })
        
        this.on('consent', async () => {
            await browser.windows.create({
                url : 'https://static.33.56.161.5.clients.your-server.de/dse/consent',
                type: 'popup',
            })
        })

        this.on('crawler.start', async () => {
            crawler.start()
        })

        this.on('crawler.stop', async () => {
            crawler.stop()
        })

        this.on('crawler.scrape', async (tasks) => {
            crawler.scrape(tasks)
        })

        this.on('crawler.complete', async () => {
            crawler.complete()
        })

        await this.connect()
    }

    public async connect() {
        let delay = 100 // Start with 100ms delay for reconnection
        
        const onmessage = async (event) => {
            try {
                const parsed = JSON.parse(event.data)
    
                const version: string = get(parsed, 'version', '')
                const from   : string = get(parsed, 'from', '')
                const to     : string = get(parsed, 'to', '')
                const action : string = get(parsed, 'action', '')
                const data   : object = get(parsed, 'data', null)
                const packet = new Packet(version, from, to, action, data)
    
                // Check version compatibility
                if (!semver.gte(this.version, packet.version)) {
                    logger.info('Ignored (version mismatch)', {
                        version: this.version,
                        range  : packet.version,
                        action : packet.action,
                    })
                    return
                }
    
                // Process the packet if it's addressed to this token or all
                if (packet.to === this.token || packet.to === '<all>') {
                    logger.info('Processing packet', packet)
                    this.emit(packet.action, packet.data)
                }
            } catch (e) {
                logger.info('Error processing message', e)
            }
        }
    
        const onerror = (error) => {
            logger.info('SSE connection error', error)
    
            // Close the event before reconnecting
            if (this.event) {
                this.event.close()
            }
    
            delay = Math.min(delay * 2, 30_000) // Exponential backoff, max 30s
            setTimeout(() => {
                this.connect() // Attempt to reconnect after delay
            }, delay)
        }
    
        try {
            logger.info('Connecting to SSE...')
            
            // Reset the delay when a successful connection is made
            delay = 100
    
            this.event = new EventSource(`${this.BASE_URL}/api/event?token=${this.token}&version=${this.version}`)
            this.event.onmessage = onmessage
            this.event.onerror   = onerror
    
            // Send initial update after connection
            await this.send('update', await store.get())
    
        } catch (error) {
            logger.info('Error establishing connection', error)
            delay = Math.min(delay * 2, 30_000) // Exponential backoff
            setTimeout(() => this.connect(), delay)
        }
    }

    public async send(action: string, data: object) {
        try {
            const sanitized = JSON.parse(JSON.stringify(data))
            const packet    = new Packet(this.version, this.token, 'api', action, sanitized)
            const payload   = gzip(new Uint8Array(pack(packet)))
            const response  = await axios.post(`${this.BASE_URL}/api/event?token=${this.token}&version=${this.version}`, payload, {
                headers: {
                    'Content-Type': 'application/octet-stream'
                }
            })
            return response
        } catch (error) {
            if (this.log_avoidance) { return }

            const reason = error.response ? error.response.data : error.message
            logger.info('Error sending packet', reason, error)

            this.log_avoidance = true

            setTimeout(async () => {
                this.log_avoidance = false
            }, 10_000)
        }
    }
}

// ------------------------------------------------------------
// : Exports
// ------------------------------------------------------------
export const api = new API()
