// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------
// Browser (polyfill)
import browser from 'webextension-polyfill'

// General
import axios from 'axios'
import { pack } from 'msgpackr'

// NATS
import type { NatsConnection, NatsError, Msg } from 'nats.ws'

// Utils
import { storage }        from '@/background/utils/storage'
import { Logger }         from '@/background/utils/logger'
import { window_manager } from '@/background/utils/window'
import { ipc }            from '@/background/modules/ipc'
import { local_ipc }      from '@/background/modules/ipc'
// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type Window     = browser.Windows.Window
type WindowOpts = browser.Windows.CreateCreateDataType
type Tab        = browser.Tabs.Tab
type TabOpts    = browser.Tabs.CreateCreatePropertiesType

type Visitation = {
    keyword: string
    website: string
    url    : string
}

type Result = {
    ok   : boolean
    error: string
}

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
const logger  = new Logger('Crawler')
const loggert = new Logger('Tab')
// ------------------------------------------------------------
// : Helpers    
// ------------------------------------------------------------
const sleep = (ms) => {
    return new Promise(resolve => setTimeout(resolve, ms))
}

const wait_until = async (predicate, timeout=-1, interval=100) => {
    const time_start = Date.now()
    while (true) {
        const result = await predicate()
        if (result) {
            break
        }
        if (timeout > 0 && Date.now() - time_start > timeout) {
            throw new Error('TimeoutError')
        }
        await sleep(interval)
    }
}

// ------------------------------------------------------------
// : Crawler
// ------------------------------------------------------------
export class Crawler {
    private token : string

    private window: browser.Windows.Window

    private timeout: number

    public ready : boolean = false
    public active: boolean = false

    public async init() {
        // Set default values
        this.token = await storage.get('user.token')
        this.window = null
        this.active = false

        await storage.set({'crawler.active': false})

        // Close all tabs with dse=1
        const tabs = await browser.tabs.query({})
        for (const tab of tabs) {
            if (tab && tab.url && tab.url.includes('dse=1')) {
                await browser.tabs.remove(tab.id)
            }
        }

        // TODO: Define topics here as constant

        // Monitor changes on the active state
        new Promise(async () => {
            while (true) {
                if (!this.window) {
                    this.active = false
                    await storage.set({'crawler.active': false})
                }
                await sleep(1000)
            }
        })

        // Set listeners
        local_ipc.on('popup.start_crawler', async () => {
            ipc.publish(`bms.dse.v1.5.crawler.start`, {
                token: this.token
            })
        })

        local_ipc.on('popup.stop_crawler', async () => {
            ipc.publish(`bms.dse.v1.5.crawler.stop`, {
                token: this.token
            })
        })


        logger.log(`Registering bms.dse.${this.token}.crawler.start`)
        ipc.subscribe(`bms.dse.${this.token}.crawler.start`, async (err: NatsError, msg: Msg) => {
            logger.log(`Starting`)
            await this.start()
            msg.respond(JSON.stringify(""))
        })

        logger.log(`Registering bms.dse.${this.token}.crawler.stop`)
        ipc.subscribe(`bms.dse.${this.token}.crawler.stop`, async (err: NatsError, msg: Msg) => {
            logger.log(`Stopping`)
            await this.stop()
            msg.respond(JSON.stringify(""))
        })

        logger.log(`Registering bms.dse.${this.token}.crawler.scrape`)
        ipc.subscribe(`bms.dse.${this.token}.crawler.scrape`, async (err: NatsError, msg: Msg) => {
            const decoder = new TextDecoder()
            const data    = JSON.parse(decoder.decode(msg.data))

            logger.log('Scraping', data)
    
            const visitation = data as Visitation
            
            await this.crawl(visitation)
                .then((result) => {
                    msg.respond(JSON.stringify(result))
                })
                .catch((error) => {
                    msg.respond(JSON.stringify({ok: false, error: error}))
                })
        })

        this.ready = true
    }

    public async start() {
        await wait_until(() => this.ready == true)
        if (this.active) {
            console.warn('Crawler already started')
            return
        }
        
        // Close all tabs
        try       { window_manager.close_all_tabs() }
        catch (e) {  }

        // Close all windows
        try       { window_manager.close_all_windows() }
        catch (e) {  }

        // Creata a new window
        this.window = await window_manager.new_window({
            url  : 'popup.html',
            state: 'minimized'
        })

        this.active = true
        await storage.set({'crawler.active': true})

        setTimeout(async () => {


        }, 60 * 1000) // 1 minute
    }

    public async stop() {
        if (!this.active) {
            console.warn('Crawler already stopped')
            return
        }

        // Close the window
        await window_manager.close_window(this.window)

        // Remove instance from memory
        this.window = null
        this.active = false
        await storage.set({'crawler.active': false})
    }

    public async crawl(visitation: Visitation): Promise<Result> {
        await wait_until(() => this.ready == true)

        if (this.active == false) return {ok: false, error: 'Crawler is not active'}

        // If window is not open, discard
        if (!this.window) {
            logger.error('Window is not open, opening a new one')
            this.window = await window_manager.new_window({
                url  : 'popup.html',
                state: 'minimized'
            })
        }

        try {
            // Open a new tab
            const tab = await window_manager.new_tab({ 
                windowId: this.window.id,
                url     : visitation.url,
                active  : false,
            })

            // Wait for payload
            const payload = await new Promise<Buffer>((resolve) => {
                local_ipc.on(`tab.${tab.id}.upload`, (data: object) => { 
                    const payload = pack(data)
                    const size    = new Blob([payload]).size

                    logger.log(`Payload size: ${size} bytes`)

                    resolve(payload)
                })
            })

            // Upload data
            await axios({
                method : 'post',
                baseURL: `${import.meta.env.BASE_URL}`,
                url    : '/api/upload',
                data   : payload,
                timeout: 1 * 50 * 1000, // 50 seconds
                headers: {'Content-Type': 'application/json'},
            }).then(resp => {
                const data = resp.data['data']
                logger.log('Uploaded data')
            }).catch(err => {
                logger.error(err)
            })

            // Close tab
            await window_manager.close_tab(tab.id)

            return {ok: true, error: null}
        } catch (e) {
            // Close window
            try       { await window_manager.close_window(this.window) }
            catch (e) { logger.error(e) }
            
            // Close all tabs
            try       { await window_manager.close_all_tabs() }
            catch (e) { logger.error(e) }

            this.window = null
            this.active = false
            await storage.set({'crawler.active': false})

            return {ok: false, error: e.message}
        }
    }
}

// ------------------------------------------------------------
// : Instance
// ------------------------------------------------------------
export const crawler = new Crawler()
