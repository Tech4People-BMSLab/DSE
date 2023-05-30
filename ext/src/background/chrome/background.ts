/// <reference types="chrome"/>
// ------------------------------------------------------------
// :: Import
// ------------------------------------------------------------

// Polyfill
import browser                from 'webextension-polyfill'

// DateTime
import { DateTime }           from 'luxon'

// Async
import { Mutex }              from 'async-mutex'
import { parallelLimit }      from 'async'

import msgpack                from 'msgpack-lite'

// Lodash
import { get }                       from 'lodash'
import { has }                       from 'lodash'
import { map }                       from 'lodash'
import { flatten }                   from 'lodash'
import { filter }                    from 'lodash'
import { isString    as is_string}   from 'lodash'
import { isNumber    as is_number}   from 'lodash'
import { isEmpty     as is_empty }   from 'lodash'
import { defaultTo   as default_to } from 'lodash'

// Events
import EventEmitter from 'eventemitter3'

// XState
import { createMachine }   from '@xstate/fsm'
import { interpret }       from '@xstate/fsm'

import Bowser from 'bowser'

import { local_config } from './config'

// ------------------------------------------------------------
// :: Util
// ------------------------------------------------------------
// https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore
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

const format_bytes = (bytes, decimals = 2) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ['Bytes', 'KB', 'MB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

// ------------------------------------------------------------
// :: Logger
// ------------------------------------------------------------
class Logger {
    private readonly MAX_LOGS = 50

    private logs: string[] = []
    private excluded = [
        // `^AlarmManager`,
        // `^SearchProcess`,
        // `^Content:.+MouseCapture:`
    ]

    constructor() {
        this.log.bind(console)
        this.warn.bind(console)
        this.error.bind(console)
    }

    private processArgs(args: any[]): any[] {
        const timestamp = DateTime.local().toFormat('[yyyy-LL-dd HH:mm:ss]')
        args.unshift(timestamp)

        // Add to logs
        this.logs.push(args.join(' '))

        // Keep the last 50 logs (MAX_LOGS)
        if (this.logs.length > this.MAX_LOGS) {
            this.logs = this.logs.slice(this.logs.length - this.MAX_LOGS)
        }

        return args
    }

    log(...args) {
        args = this.processArgs(args)
        
        console.log.apply(console, args)
    }

    warn(...args) {
        args = this.processArgs(args)
        console.warn.apply(console, args)
    }

    error(...args) {
        args = this.processArgs(args)
        console.error.apply(console, args)
    }
}

// ------------------------------------------------------------
// :: Types
// ------------------------------------------------------------
type Window     = browser.Windows.Window
type WindowOpts = browser.Windows.CreateCreateDataType
type Tab        = browser.Tabs.Tab
type TabOpts    = browser.Tabs.CreateCreatePropertiesType

// ------------------------------------------------------------
// :: Classes
// ------------------------------------------------------------

// ------------------------------------------------------------
// : Storage
// ------------------------------------------------------------
/**
 * Wrapper class for the browser storage.
 * The wrapper class provides a bit more simpler interface for the storage to fetch data in a single function call.
 */
class Storage {
    /**
     * The storage object
     */
    private storage  = browser.storage.local

    /**
     * The listener for the storage changes
     */
    public onChanged = browser.storage.local.onChanged

    /**
     * Get a value from the storage
     * @param key The key to get the value from
     * @returns The value of the key
     */
    public async get(key: string) {
        try {
            const result = await this.storage.get(key)
            if (!has(result, key)) {
                return undefined
            }
            return result[key]
        } catch (e) {
            console.error(e)
        }
    }

    // Define functions to get all
    public async get_all() {
        try {
            return await this.storage.get()
        } catch (e) {
            console.error(e)
        }
    }
    
    /**
     * Set a value in the storage
     * @param obj The object to set the values from
     */
    public async set(obj: object) {
        try {
            await this.storage.set(obj)
        } catch (e) {
            console.error(e)
        }
    }
    
    /**
     * Listen to a value in the storage
     * @param event The event to listen to
     * @param callback The callback to call when the event is triggered
     */
    public on(event: string, callback: Function) {
        storage.onChanged.addListener(async(changes: object) => {
            if (!has(changes, event))               return
            if (!has(changes, `${event}.newValue`)) return

            const value = changes[event].newValue
            callback(value)
        })
    }

    /**
     * Remove a value from the storage
     * @param key The key to remove from the storage
     */
    public async remove(key: string) {
        try {
            await this.storage.remove(key)
        } catch (e) {
            console.error(e)
        }
    }

    /**
     * clear the storage
     */
    public async clear() {
        try {
            await this.storage.clear()
        } catch (e) {
            console.error(e)
        }
    }
}

// ------------------------------------------------------------
// : Window Manager
// ------------------------------------------------------------
class WindowManager {
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
            logger.log(`WindowManager: Window already closed <${window}>`)
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
        await this.mutex.waitForUnlock()

        try {
            if (is_number(tab)) { tab = await browser.tabs.get(tab) }
            if (is_string(tab)) { tab = await browser.tabs.get(parseInt(tab)) }

            await browser.tabs.remove(tab.id)
        } catch (error) {
            logger.log(`WindowManager: Tab already closed <${tab}>`)
        }

        let tabs: Tab[]
        tabs = await storage.get('window_manager.tabs')
        tabs = filter(tabs, t => t.id !== tab['id'])
        await storage.set({'window_manager.tabs': tabs})
    }

    public async close_all_tabs(): Promise<void> {
        await this.mutex.waitForUnlock()

        const tabs = await storage.get('window_manager.tabs')
        tabs.forEach(tab => {
            this.close_tab(tab.id)
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
            this.close_window(window.id)
        })
    }

    public async open_popup(opts: {minimized: boolean}): Promise<void> {
        await this.mutex.waitForUnlock()

        switch (opts.minimized) {
            case false: await this.new_tab   ({url: 'popup.html'})                     ; break
            case true : await this.new_window({url: 'popup.html', state: 'minimized',}); break

            default: logger.error(`WindowManager: Invalid option for open_popup: ${opts.minimized}`)
        }
    }
}

// ------------------------------------------------------------
// : Search Process
// ------------------------------------------------------------
class SearchProcess {
    readonly MAX_TABS  = 1
    readonly WAIT_TIME = 4000 // 4 seconds


    readonly websites = [
        //@todo: remove comments
        {name: 'Google'        ,url: 'https://www.google.com/search?q='},
        {name: 'Google News'   ,url: 'https://www.google.com/search?tbm=nws&q='},
        {name: 'Google Videos' ,url: 'https://www.google.com/search?tbm=vid&q='},
        {name: 'DuckDuckGo'    ,url: 'https://duckduckgo.com/?q='},
        {name: 'YouTube'       ,url: 'https://www.youtube.com/results?search_query='},
        // {name: 'Twitter'       ,url: 'https://twitter.com/search?q='},
        {name: 'Bing'          ,url: 'https://www.bing.com/search?q='},
        {name: 'Yahoo'         ,url: 'https://search.yahoo.com/search?p='},
    ]
    
    private window  : Window       // To keep track of window reference
    private tabs    : Tab[]  = []  // To keep track of all the tabs created

    private mutex = new Mutex()

    private state  : any
    private service: any

    constructor() {
        this.state = createMachine({
            id     : 'search_process',
            initial: 'inactive',
            context: {
                is_connected : false,
                is_registered: false,
            },

            states: {
                inactive: { 
                    on: { 
                        START:  {
                            target: 'active',
                            cond: (context, event: any) => {
                                return event.is_connected && event.is_registered
                            },
                        }
                    }
                },

                active: { 
                    entry: (context, event) => {this.on_start()},
                    on   : { STOP:  {target: 'cancelled'}}
                },

                cancelled: { 
                    entry: (context, event) => {this.on_cancel()},
                    on   : { DONE:  {target: 'inactive'}}
                },
            }
        })

        this.service = interpret(this.state).start()
        this.service.subscribe(async (state) => { // Update state in storage on change
            await storage.set({'search_process.state': state.value})

            try {
                wsclient.send({
                    token: await storage.get('user.token') || await storage.get('token'),
                    action: 'event',
                    data: {
                        type: `search_process.state.${state.value}`,
                    }
                })
            } catch (e) { /* Ignore */ }

            try {
                // This can trigger error if the popup is not open (ignore it)
                await browser.runtime.sendMessage({
                    target: 'popup',
                    state : state
                })
            } catch (e) { /* Ignore */ }
        })

        // Set handler for 'search process'
        browser.runtime.onMessage.addListener(async(payload: object) => {
            const target = get(payload, 'target')
            const action = get(payload, 'action')

            // Check if the message is for this process
            if (!target || target !== 'search_process') return

            switch (action) {
                case 'start':
                    logger.log(`SearchProcess: Received start request`)
                    this.start()
                    break
                case 'stop':
                    logger.log(`SearchProcess: Received stop request`)
                    this.stop()
                    break
                case 'get_state':
                    logger.log(`SearchProcess: Received get_state request`)
                    const id    = browser.runtime.id
                    const state = this.get_state()

                    browser.runtime.sendMessage(id, {
                        target: 'popup',
                        state : state
                    })
                    break
            }
        })

        // Listen for window close event
        browser.windows.onRemoved.addListener(this.on_window_close.bind(this))
        return
    }

    //// Public Methods
    public async start() {
        const is_connected  = await storage.get('server.connected')
        const is_registered = await storage.get('user.registered')

        this.service.send({
            type: 'START',
            is_connected : is_connected,
            is_registered: is_registered,
        })
    }
    
    public async stop() {
        switch (this.service.state.value) {
            case 'stopped'  : logger.log(`SearchProcess: Already stopped`) ; break
            case 'started'  : logger.log(`SearchProcess: Stopping`)        ; break
            case 'cancelled': logger.log(`SearchProcess: Already stopping`); break
        }

        // Close all tabs related to this process
        for (const tab of this.tabs) {
            try       { await browser.tabs.remove(tab.id) }
            catch (e) { /* Ignore */ }
        }

        this.service.send('STOP')
    }

    public get_state() {
        return this.service.state.value
    }

    //// State Handlers
    private async on_start() {
        // Acquire the lock
        if (!await this.mutex.acquire()) {
            return
        }

        // Change state to BUSY
        logger.log(`SearchProcess: Starting search process...`)

        // Get keywords from server
        const config   = await storage.get('config')
        const keywords = config['keywords']

        // Create a new window
        this.window = await window_manager.new_window({
            url  : 'popup.html',
            state: 'minimized',
        })
        logger.log(`SearchProcess: Window <${this.window.id}> created`)


        // Map keyword with their unique website
        let visitations = flatten(map(keywords, (keyword) => {
            return map(this.websites, (website) => {
                return {
                    'keyword': keyword,
                    'website': website,
                    'url'    : website.url + keyword,
                    'state'  : 'unfinished'
                }
            })
        }))

        logger.log(`SearchProcess: Has to do ${visitations.length} visits`)

        // Define function to iterate through each visitation
        const tasks = map(visitations, (visitation) => {
            return (cb) => {                
                const fn = async () => {
                    const keyword = visitation.keyword
                    const website = visitation.website
                    const url     = `${website.url}${keyword.replace(/\s/g, '+')}`

                    // Condition to stop the process
                    if (this.get_state() === 'cancelled') {
                        logger.log(`SearchProcess: Cancelled, aborting...`)
                        cb(null, -1)
                        return
                    }

                    // Stop the process if user has closed the search window
                    if (!this.window) {
                        logger.log(`SearchProcess: Window closed, aborting...`)
                        cb(null, -1)
                        return
                    }

                    // Generate a hash for the tab
                    const tab_hash = (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1)
                    
                    // Create new tab
                    const tab = await window_manager.new_tab({
                        windowId: this.window.id,
                        url     : `${url}&dse=${tab_hash}&dse_keyword=${keyword}&dse_website=${website.name}`,
                        active  : false,
                    })
                    logger.log(`SearchProcess: Tab <${tab.id}> created for <${keyword}> on <${website.name}>`)

                    // Add tab to the list 
                    this.tabs.push(tab)

                    // Get current time
                    const time_start = DateTime.now()

                    await new Promise((resolve) => {
                        // Wait for timeout
                        const timeout = setTimeout(() => {
                            resolve(true)
                            return
                        }, 30_000)

                        const interval = setInterval(async() => {
                            const state = this.get_state()
                            if (state === 'cancelled') {
                                logger.log(`SearchProcess: Tab <${tab.id}> cancelled`)
                                clearTimeout(timeout)
                                clearInterval(interval)
                                resolve(true)
                                return
                            }
                        }, 500)
                        
                        // Wait until receives event from the tab (content script)
                        const receive_handler = async (payload) => {
                            const target = get(payload, 'target')
                            const from   = get(payload, 'from')
                            const state  = get(payload, 'state')
                            const error  = get(payload, 'error')

                            if (target !== 'search_process' || from !== 'extractor') {
                                return
                            }

                            switch (state) {
                                case 'success':
                                    logger.log(`SearchProcess: Tab <${tab.id}> finished`)
                                    break
                                case 'failed':
                                    logger.error(`SearchProcess: Tab <${tab.id}> error <${JSON.stringify(error)}>`)
                                    break
                            }
                            clearTimeout(timeout)
                            clearInterval(interval)
                            resolve(true)  
                            
                            browser.runtime.onMessage.removeListener(receive_handler)
                            return
                        }
                        browser.runtime.onMessage.addListener(receive_handler)
                    }).catch(async (err) => {
                        logger.error(`SearchProcess: Tab <${tab.id}> error <${err.message}>`)
                    })

                    // Close the tab
                    await window_manager.close_tab(tab.id)

                    // Remove from the list
                    this.tabs = this.tabs.filter((t) => t.id !== tab.id)

                    // Get the current time and calculate the delta
                    const time_end   = DateTime.now()
                    const time_delta = time_end.diff(time_start, 'milliseconds').milliseconds

                    // Wait the remaining time
                    await sleep(Math.max(this.WAIT_TIME - time_delta))
                    
                    // Continue to the next visitation
                    cb(null, 1)
                }

                // Run the function
                fn()
            }
        })

        parallelLimit(tasks, this.MAX_TABS, async (err, results) => {
            if (err) {
                logger.error(`SearchProcess: An error occurred in one of the tasks: ${err}`)
            } else {
                logger.log(`SearchProcess: Finished`)
                logger.log(results)
            }

            // Close the tabs (eventhough they might have been closed already)
            for (const tab of this.tabs) {
                window_manager.close_tab(tab.id)
            }

            // Close window (eventhough it might have been closed already)
            await window_manager.close_window(this.window.id)

            // Reset state
            this.stop()

            // Release the lock
            this.mutex.release()
        })
        return
    }

    private async on_cancel() {
        await this.mutex.release()
        this.service.send('DONE')
    }

    //// Event Handlers
    public async on_window_close(id: number) {
        logger.log(`SearchProcess: Window <${id}> closed`)
        if (!is_empty(this.window)) {
            if (this.window.id == id) {
                // Stop the process
                this.stop()
            }
        }
        return
    }
}

// ------------------------------------------------------------
// : WebSocket Client
// ------------------------------------------------------------
type Packet = {
    token    : string | null
    action   : string | null
    data     : object | null
}

class WSClient {
    private socket: WebSocket
    
    private url  : string
    private mutex: Mutex = new Mutex()

    public status    : 'connected' | 'connecting' | 'disconnected' = 'disconnected'
    public registered: 'registered' | 'unregistered'               = 'unregistered'

    private delay: number = 1000 // in ms

    constructor() {
        this.url = WS_URL
    }
    
    //// Initializer
    public async init() {
        try {
            if (this.status === 'connected') {
                return
            }

            // Set status
            this.status = 'connecting'
            await storage.set({'server.connected': false})

            // Create the socket
            this.socket = new WebSocket(this.url)

            // Add the listeners
            this.socket.onopen    = this.on_connect.bind(this)
            this.socket.onmessage = this.on_receive.bind(this)
            this.socket.onclose   = this.on_disconnect.bind(this)

            this.socket.onerror   = (e) => {
                // console.error(e)
            }
        } catch (e) { /* Ignore */ }
    }

    //// Event Handlers
    private async on_connect(event) {
        console.log(`Connected to ${this.url}`)

        // Set the status
        this.status = 'connected'
        await storage.set({'server.connected': true})

        // Reset the delay
        this.delay = 50 

        // Check-in with the server
        await this.socket.send(msgpack.encode({
            'timestamp': DateTime.now().toISO(),
            'action'   : 'connect',
            'data'     : {
                'token': await storage.get('user.token') || await storage.get('token')
            }
        }))
    }

    private async on_receive(event) {
        try {
            if (event.data instanceof Blob) {
                // const payload = event.data as Blob
                const payload = await new Response(event.data).arrayBuffer()
                const decoded = msgpack.decode(new Uint8Array(payload))

                const action = decoded.action
                const data   = decoded.data

                const id = browser.runtime.id // Get the extension id to send message to other scripts

                console.log(`Received message: ${JSON.stringify(decoded)}`)

                switch (action) {
                    case 'register.success':
                        await storage.set({'user.token'     : data.token})
                        await storage.set({'token'          : data.token})
                        await storage.set({'user.registered': true})

                        browser.runtime.sendMessage(id, {
                            target : 'page.register',
                            action : 'success'
                        })

                        // Check-in with the server
                        await this.socket.send(msgpack.encode({
                            'timestamp': DateTime.now().toISO(),
                            'action'   : 'connect',
                            'data'     : {
                                'token': await storage.get('user.token') || await storage.get('token')
                            }
                        }))
                        break

                    case 'register.failed':
                        browser.runtime.sendMessage(id, {
                            target : 'page.register',
                            action : 'failed'
                        })
                        break

                    case 'token.invalid':
                        await storage.set({'user.registered': false})
                        break

                    case 'token.valid':
                        await storage.set({'user.registered': true })
                        break

                    case 'search_process.start':
                        search_process.start()
                        break

                    case 'search_process.stop':
                        search_process.stop()
                        break

                    default:
                        console.log(`WSClient: Unknown action: ${action}`)
                        break
                }
            }
        } catch (e) {
            console.error(e)
        }

    }

    private async on_disconnect(event) {
        console.error(`Disconnected, reason: ${(event as CloseEvent).code}`)

        // Set status
        this.status = 'disconnected'
        await storage.set({'server.connected': false})

        if (this.mutex.isLocked()) return
        await this.mutex.acquire()

        while (this.status === 'disconnected') {
            await sleep(this.delay)

            console.log(`Reconnecting...`)
            await this.init()

            // Increase the delay
            if (this.delay < 10_000) this.delay *= 2
            if (this.delay > 10_000) this.delay = 10_000
        }

        await this.mutex.release()
    }

    //// Public Methods
    public async send(data: Packet) {
        try {
            const payload = {
                timestamp: DateTime.now().toISO(),
                token    : await storage.get('user.token') || await storage.get('token'),
                ...data
            }

            const buffer = msgpack.encode(payload)
            this.socket.send(buffer)
        } catch (e) {
            console.error(e)
        }
    }
}


// ------------------------------------------------------------
// :: Background
// ------------------------------------------------------------
class Background {
    private state  : any
    private service: any

    /**
     * Initial point of the script.
     */
    public async init() {
        this.state = createMachine({
            id     : 'background',
            initial: 'startup',

            states: {
                startup: {
                    entry: (context, event) => { this.on_startup() },
                    on: {
                        VERIFY: { target: 'verifying' },
                    },
                },

                verifying: {
                    entry: (context, event) => { this.on_verifying() },
                    on: {
                        REGISTERED    : {target: 'ready'      },
                        NOT_REGISTERED: {target: 'registering'},
                    },
                },
                
                registering: {
                    entry: (context, event) => { this.on_registering() },
                    on: {
                        READY: {target: 'ready'},   // If the user registers, we will be ready to go
                        CLOSE: {target: 'on_hold'}, // If the user closes the registration window, we will wait for them to open it again
                    },
                },

                on_hold: {
                    on: {
                        NEXT: {target: 'verifying'},
                    }
                },

                ready: {
                    entry: (context, event) => { this.on_ready() }
                },
            },
        })

        this.service = interpret(this.state).start()
        this.service.subscribe((state) => {
            logger.log(`Background: State changed to ${state.value}`)
        })
    }

    //// State Handlers
    private async on_startup() {
        if (is_debugging()) { 
            // Debugging block

            //// Clear token (new user)
            // storage.clear()

            //// Registered user (Home server)
            // storage.set({'user.token': 'fc5f3908-d57c-4e00-9f38-554695e97b07'})

            //// Registered user (Fake token)
            // storage.set({'user.token': '12345'})

            //// Registered user (BMS server)
            // storage.set({'user.token': 'e20eba0f02fa023a3bbe2eb8ec9bb9a06436c3398d2006f4a9817967607af7a5'})

            //// Load popups
            // setTimeout(() => {
            //     browser.tabs.create({
            //         url: 'popup.html'
            //     })
            // }, 100)


            //// Force start search process
            // setTimeout(() => {
            //     search_process.start()
            // }, 150)

            //// Opens the pop up in a new tab (for fast debugging)
            // const url = browser.runtime.getURL('/popup.html/')
            // const tab = await browser.tabs.create({url: url})

            //// Reload extension periodically
            // setInterval(() => {
            //     browser.runtime.reload()
            // }, 2500)
        }

        // Connect to server immediately
        await wsclient.init()

        // Close all windows and tabs that belong to this extension (incase a restart happens to the extension)
        await window_manager.close_all_tabs()
        await window_manager.close_all_windows()

        // Set the extension version
        await storage.set({'extension.version': browser.runtime.getManifest().version})

        // Set the configuration
        await storage.set({'config': local_config})

        // Get and set the default language (based on the system language)
        const language = navigator.language
        await storage.set({'language': language})

        // Set the browser properties
        const browser_properties = Bowser.getParser(navigator.userAgent)
        await storage.set({'browser': {
            name       : browser_properties.getBrowserName(),
            version    : browser_properties.getBrowserVersion(),
            os         : browser_properties.getOSName(),
            os_version : browser_properties.getOSVersion(),
        }})

        // Log the start of the extension
        const now   = DateTime.now()
        const token = await storage.get('user.token') || await storage.get('token') 

        logger.log(`Background: Started at ${now.toFormat('yyyy-MM-dd HH:mm:ss')}`)
        logger.log(`Background: Version    : ${browser.runtime.getManifest().version}`)
        logger.log(`Background: Browser    : ${browser_properties.getBrowserName()} ${browser_properties.getBrowserVersion()}`)
        logger.log(`Background: Token      : ${token}`)

        // Set handler for idle state changes
        browser.idle.onStateChanged.addListener(async(state) => {
            logger.log(`Background: Idle state: ${state}`)
            await storage.set({'extension.state': state})
        })

        // Set idle detection 
        browser.idle.queryState(60).then(async (state) => {
            logger.log(`Background: Idle state: ${state}`)
            await storage.set({'extension.state': state})
        })

        // Set handler for logs from 'content' and 'popup' script
        browser.runtime.onMessage.addListener(async(payload: any) => {
            if (!payload.from) return
            if (!payload.data) return

            const from = payload.from
            const data = payload.data

            const token = await storage.get('user.token') || await storage.get('token')

            switch (from) {
                case 'logger':
                    logger.log(`Content: ${data}`)
                    break

                case 'content.extractor':
                    try {
                        wsclient.send({
                            token : token,
                            action: 'content.extractor',
                            data  : data,
                        })
                        logger.log(`Background: Sent 'extractor' data`)
                        logger.log(`Background: Payload Size: ${format_bytes(JSON.stringify(data).length)}`)
                    } catch (e) {
                        logger.log(`Background: Error sending 'extractor' data (${e})`)
                    }
                    
                    break

                case 'content.click':
                    try {
                        wsclient.send({
                            token : token,
                            action: 'content.click',
                            data  : data,
                        })
                        logger.log(`Background: Sent 'click' data`)
                        logger.log(`Background: Payload Size: ${format_bytes(JSON.stringify(data).length)}`)
                    } catch (e) {
                        logger.log(`Background: Error sending 'click' data (${e})`)
                    }
                    break
                
                case 'page.register': 
                    try {
                        wsclient.send({
                            token : null,
                            action: 'register',
                            data  : data,
                        })
                        logger.log(`Background: Sent data from registration process to server`)
                    } catch (e) {
                        logger.log(`Background: Error sending data from registration process to server (${e})`)
                    }
                    break
            }
        })

        

        this.service.send('VERIFY')
    }

    private async on_verifying() {
        if (await storage.get('user.registered')) {
            this.service.send('REGISTERED')
        } else {
            this.service.send('NOT_REGISTERED')
        }
    }

    private async on_registering() {
        if (!await storage.get('user.registered_at')) {
            // Open the consent form
            await window_manager.open_popup({minimized: false})
        }
    }

    private async on_ready() {
        // write to storage the registered time
        if (!await storage.get('user.registered_at')) {
            await storage.set({'user.registered_at': DateTime.now().toISO()})
        }
    }
}
// ------------------------------------------------------------
// :: Global
// ------------------------------------------------------------
const WS_URL  = `ws://dev.bmslab.utwente.nl:5004`
// const WS_URL  = `wss://dev.bmslab.utwente.nl/dse/ws`
// const WS_URL  = `wss://dev.bmslab.utwente.nl/dse_test/ws`
// const WS_URL  = `ws://10.0.0.10:5000/ws`
// const WS_URL  = `wss://apps.bmslab.utwente.nl/dse/ws`
// const WS_URL  = `ws://apps.bmslab.utwente.nl:5004/ws`

const logger       = new Logger()
const emitter      = new EventEmitter()
const storage      = new Storage()

const background = new Background()

const window_manager = new WindowManager()
const search_process = new SearchProcess()
const wsclient       = new WSClient()

// ------------------------------------------------------------
// :: Main
// ------------------------------------------------------------
background.init()
