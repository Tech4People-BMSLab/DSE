// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------
// Browser (polyfill)
import browser from 'webextension-polyfill'

// General
import _ from 'lodash'
import EventEmitter from 'eventemitter3'

// NATS
import { connect, NatsConnection, Events, DebugEvents, NatsError } from 'nats.ws'

// Logger
import { Logger }         from '@/background/utils/logger'
import { window_manager } from '@/background/utils/window'
import { storage }        from '@/background/utils/storage'
import { wait_until }     from '@/background/utils/utils'
import { sleep }          from '@/background/utils/utils'
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
const logger = new Logger('IPC')

const nats_url  = import.meta.env.NATS_URL
const nats_user = import.meta.env.NATS_USER
const nats_pass = import.meta.env.NATS_PASS

// ------------------------------------------------------------
// : Error Handlers
// ------------------------------------------------------------
async function on_disconnect() {
    logger.log('Disconnected')
    await sleep(1000)
    logger.log('Reloading...')
    browser.runtime.reload()
}

async function on_reconnecting() {
    logger.log('Disconnected')
    await sleep(1000)
    logger.log('Reloading...')
    browser.runtime.reload()
}

async function on_error() {
    logger.log('Disconnected')
    await sleep(1000)
    logger.log('Reloading...')
    browser.runtime.reload()
}

// ------------------------------------------------------------
// : IPC (Interprocess Communication)
// ------------------------------------------------------------
/**
 * Communicate with the DSE server using NATS
 */
export class IPC {
    private token : string

    public nc       : NatsConnection
    public connected: boolean = false

    async init() {
        await this.connect()
    }

    private async connect(): Promise<void> {
        try {
            storage.set({'ipc.connected': false})

            let config = {
                servers: [nats_url],
                user    : nats_user,
                pass    : nats_pass,

                timeout             : 1000,
                reconnect           : true,
                reconnectTimeWait   : 1000,
                maxReconnectAttempts: Number.MAX_VALUE,
            }

            if (config.user == 'null') config.user = undefined
            if (config.pass == 'null') config.pass = undefined

            this.nc = await connect(config)

            new Promise(async () => {
                for await (const s of this.nc.status()) {
                    switch (s.type) {
                        case Events.Disconnect          : on_disconnect(); break;
                        case DebugEvents.Reconnecting   : on_reconnecting(); break;
                        case Events.Error               : on_error(); break;

                        case Events.LDM                 : logger.log("Requested reconnect")       ; break;
                        case Events.Update              : logger.log(`Cluster update`)            ; break;
                        case Events.Reconnect           : logger.log(`Client reconnected`)        ; break;
                        case DebugEvents.StaleConnection: logger.log("Stale connection")          ; break;
                        default                         : break;
                    }
                }
            })

            this.connected = true
            storage.set({'ipc.connected': true})
            logger.log('Connected')
        } catch (e) {
            /** Ignore */
        }
    }

    public async register(): Promise<void> {

        // Wait until connected and wait for token
        await wait_until(() => this.connected)
        await wait_until(() => storage.get('user.token'))

        this.token = await storage.get('user.token')

        // Define direct topics
        const topic_user_connected = `bms.dse.users.${this.token}.connected` 
        const topic_user_ping      = `bms.dse.users.${this.token}.ping`
        const topic_user_consent   = `bms.dse.users.${this.token}.consent`
        const topic_user_reload    = `bms.dse.users.${this.token}.background.reload`

        // Define broadcast topics
        const topic_users_ping     = `bms.dse.users.ping`
        const topic_users_pong     = `bms.dse.users.pong`

        // Define local/debug topics
        const topic_reload   = `bms.dse.background.reload`

        { /// Publish
            // Publish connected message
            await this.nc.publish(topic_user_connected, JSON.stringify({
                token   : await storage.get('user.token'),
                form    : await storage.get('user.form'),
                storage : await storage.get_all(),
                manifest: browser.runtime.getManifest(),
            }))
        }

        { /// Subscribe
            // Listen for broadcast ping
            await this.nc.subscribe(topic_users_ping, {
                callback: async (err, msg) => {
                    if (err) {
                        logger.error(err)
                        return
                    }

                    logger.log('Received ping')

                    await this.nc.publish(topic_users_pong, JSON.stringify({
                        token   : await storage.get('user.token'),
                        form    : await storage.get('user.form'),
                        storage : await storage.get_all(),
                        manifest: browser.runtime.getManifest(), 
                    }))
                }
            })

            // Listen for direct ping
            await this.nc.subscribe(topic_user_ping, {
                callback: async (err, msg) => {
                    if (err) {
                        logger.error(err)
                        return
                    }

                    await msg.respond(JSON.stringify({
                        token   : await storage.get('user.token'),
                        form    : await storage.get('user.form'),
                        storage : await storage.get_all(),
                        manifest: browser.runtime.getManifest(),
                    }))
                }
            })

            // Re-popup consent form
            await this.nc.subscribe(topic_user_consent, {
                callback: async () => {
                    browser.tabs.create({
                        url: 'http://dev.bmslab.utwente.nl/dse/consent'
                    })
                }
            })
    
            // Reload the extension
            this.nc.subscribe(topic_user_reload, {
                callback: async () => {
                    browser.runtime.reload()
                }
            })
        }

        { /// Debugging
            const is_development = import.meta.env.MODE === 'development'

            if (is_development == false) return

            await this.nc.subscribe(topic_reload, {
                callback: async () => {
                    logger.log('Reloading extension...')
                    browser.runtime.reload()
                }
            })
        }
    }

    public async publish(topic: string, data: any): Promise<void> {
        await wait_until(() => this.connected)
        await this.nc.publish(topic, JSON.stringify(data))
    }

    public async subscribe(topic: string, callback: (err: NatsError, msg: any) => void): Promise<void> {
        await wait_until(() => this.connected)
        await this.nc.subscribe(topic, { callback })
    }
}

/**
 * Communicate with the local background script and the content script
 */
export class LocalIPC extends EventEmitter {
    private logger = new Logger('LocalIPC')
    private logger_tab   = new Logger('Tab')
    private logger_popup = new Logger('Popup')
    

    constructor() {
        super()
        browser.runtime.onMessage.addListener(this.on_message.bind(this))
    }

    async on_message(message: any, sender: any): Promise<any> {
        try {
            const type   = message.type
            const from   = message.from
            const action = message.action
            const data   = message.data

            switch (true) {
                case type == 'log' && from == 'tab'   : this.logger_tab.log(data.join(' '))  ; break
                case type == 'log' && from == 'popup' : this.logger_popup.log(data.join(' ')); break

                case type == 'action' && from == 'tab'   : this.emit(`tab.${sender.tab.id}.${action}`, data); break
                case type == 'action' && from == 'popup' : this.emit(`popup.${action}`, data)               ; break

                default: 
                    this.logger.error('Unknown message', message)
                    break
            }
        } catch (e) {
            this.logger.error(e)
        }
    }
    
    public async send(tab_id: number, message: any): Promise<any> {
        return browser.tabs.sendMessage(tab_id, message)
    }

    public async broadcast(message: any): Promise<any> {
        return browser.runtime.sendMessage(message)
    }
}

// ------------------------------------------------------------
// : Instance
// ------------------------------------------------------------
export const ipc       = new IPC()
export const local_ipc = new LocalIPC()
