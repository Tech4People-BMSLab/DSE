// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------

// Browser (polyfill)
import browser from 'webextension-polyfill'

// Async
import { Mutex } from 'async-mutex'

// Lodash
import _ from 'lodash'

// Utils
import { storage }   from '../utils/storage'
import { Logger }    from '../utils/logger'

// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type Window     = browser.Windows.Window
type WindowOpts = browser.Windows.CreateCreateDataType
type Tab        = browser.Tabs.Tab
type TabOpts    = browser.Tabs.CreateCreatePropertiesType

// ------------------------------------------------------------
// : Utils
// ------------------------------------------------------------
const filter       = _.filter
const map          = _.map
const default_to   = _.defaultTo

const is_empty     = _.isEmpty
const is_number    = _.isNumber
const is_string    = _.isString

// ------------------------------------------------------------
// : Locasl
// ------------------------------------------------------------
const logger = new Logger('WindowManager')

// ------------------------------------------------------------
// : Class
// ------------------------------------------------------------
export class WindowManager {
    private mutex = new Mutex()

    constructor() {
        // Load all the windows and tabs, some are undefined because they can be closed when this script is asleep/idle
        this.mutex.acquire().then(async() => {
            try {
                let windows: Window[]
                let tabs   : Tab[]

                windows = await storage.get('window_manager.windows')
                tabs    = await storage.get('window_manager.tabs')

                if (is_empty(windows)) {
                    windows = []
                } else { 
                    const promises = map(windows, async(window) => {
                        try           { return await browser.windows.get(window.id) } 
                        catch (error) { return undefined }
                    })
                    windows = await Promise.all(promises)
                    windows = filter(windows, (window) => window !== undefined)
                }
                await storage.set({'window_manager.windows': windows})

                if (is_empty(tabs)) {
                    tabs = []
                } else {
                    const promises = map(tabs, async(tab) => {
                        try           { return await browser.tabs.get(tab.id) } 
                        catch (error) { return undefined }
                    })
                    tabs = await Promise.all(promises)
                    tabs = filter(tabs, (tab) => tab !== undefined)
                }
                await storage.set({'window_manager.tabs': tabs})
            } finally {
                // Release the mutex so other functions can use it
                this.mutex.release()
            }
        })

        browser.tabs.onRemoved.addListener(this.on_tab_closed.bind(this))
        browser.windows.onRemoved.addListener(this.on_window_closed.bind(this))
    }

    private async on_tab_closed(id: number) {
        await this.mutex.waitForUnlock()
        
        let tabs: Tab[]
        tabs = await storage.get('window_manager.tabs')
        tabs = filter(tabs, (tab) => tab.id !== id)
        await storage.set({'window_manager.tabs': tabs})
    }

    private async on_window_closed(id: number) {
        await this.mutex.waitForUnlock()

        let windows: Window[]
        windows = await storage.get('window_manager.windows')
        windows = filter(windows, (window) => window.id !== id)
        await storage.set({'window_manager.windows': windows})
    }

    public async new_window(opts: WindowOpts): Promise<Window> {
        await this.mutex.waitForUnlock()

        const windows = await storage.get('window_manager.windows')
        const window  = await browser.windows.create(opts)

        windows.push(window)
        await storage.set({'window_manager.windows': windows})
        return window
    }

    public async close_window(window: Window | number): Promise<void> {
        await this.mutex.waitForUnlock()

        try {
            if (is_number(window)) { window = await browser.windows.get(window) }
            if (is_string(window)) { window = await browser.windows.get(parseInt(window)) }

            await browser.windows.remove(window.id)
        } catch (error) {
            logger.log(`Window already closed <${window}>`)
        }

        let windows: Window[]
        windows = await storage.get('window_manager.windows')
        windows = filter(windows, w => w.id !== window['id'])
        await storage.set({'window_manager.windows': windows})
    }

    public async get_window(window_id: number): Promise<Window> {
        return await browser.windows.get(window_id)
    }

    public async new_tab(opts: TabOpts): Promise<Tab> {
        await this.mutex.waitForUnlock()

        const tabs = default_to(await storage.get('window_manager.tabs'), [])
        const tab  = await browser.tabs.create(opts)

        tabs.push(tab)
        await storage.set({'window_manager.tabs': tabs})
        return tab
    }

    public async close_tab(tab: Tab | number): Promise<void> {
        try {
            await this.mutex.waitForUnlock()
    
            try {
                if (is_number(tab)) { tab = await browser.tabs.get(tab) }
                if (is_string(tab)) { tab = await browser.tabs.get(parseInt(tab)) }
    
                await browser.tabs.remove(tab.id)
            } catch (error) {
                logger.log(`Tab already closed <${tab}>`)
            }
    
            let tabs: Tab[]
            tabs = await storage.get('window_manager.tabs')
            tabs = filter(tabs, t => t.id !== tab['id'])
            await storage.set({'window_manager.tabs': tabs})
        } catch (e) {
            logger.error(e)
        }
    }

    public async close_all_tabs(): Promise<void> {
        await this.mutex.waitForUnlock()

        const tabs = await storage.get('window_manager.tabs')
        tabs.forEach(tab => {
            try       { this.close_tab(tab.id) }
            catch (e) { /* Do nothing */       }
        })
    }

    public async is_window_open(id: number) {
        try { 
            await browser.windows.get(id)
            return true
        } catch (error) {
            return false
        }
    }

    public async is_tab_open(id: number) {
        try {
            await browser.tabs.get(id)
            return true
        } catch (error) {
            return false
        }
    }

    public async close_all_windows(): Promise<void> {
        await this.mutex.waitForUnlock()

        const windows = await storage.get('window_manager.windows')
        windows.forEach(window => {
            try       { this.close_window(window.id) }
            catch (e) { /* Do nothing */             }
        })
    }

    public async open_popup(opts: {minimized: boolean}): Promise<void> {
        await this.mutex.waitForUnlock()

        switch (opts.minimized) {
            case false: await this.new_tab   ({url: 'https://dev.bmslab.utwente.nl/dse/consent'})                     ; break
            case true : await this.new_window({url: 'https://dev.bmslab.utwente.nl/dse/consent', state: 'minimized',}); break

            default: logger.error(`Invalid option for open_popup: ${opts.minimized}`)
        }
    }
}

// ------------------------------------------------------------
// : Instance
// ------------------------------------------------------------
export const window_manager = new WindowManager()
