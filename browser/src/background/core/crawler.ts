    // ------------------------------------------------------------
    // : Imports    
    // ------------------------------------------------------------
    import axios        from 'axios'
    import browser      from 'webextension-polyfill'
    import EventEmitter from 'eventemitter3'
    import { pack }     from 'msgpackr'

    import { createMachine, createActor } from 'xstate'

    import { api }            from '@/background/core/api'
    import { ipc }            from '@/background/core/ipc'
    import { sleep }          from '@/background/utils/utils'
    import { wait_until }     from '@/background/utils/utils'
    import { store }          from '@/background/core/storage'
    import { Logger }         from '@/background/utils/logger'
    // ------------------------------------------------------------
    // : Locals
    // ------------------------------------------------------------
    const logger = new Logger('Crawler')

    // ------------------------------------------------------------
    // : Types
    // ------------------------------------------------------------
    type Task = {
        keyword: string
        website: string
        url    : string
    }

    // ------------------------------------------------------------
    // : Crawler
    // ------------------------------------------------------------
    class Crawler extends EventEmitter {
        public token: string

        private machine: ReturnType<typeof createMachine>
        private actor  : ReturnType<typeof createActor>

        private timer_interval = null
        private timer_timeout  = null

        public async init() {
            this.token = await store.get('user.token')

            this.machine = createMachine({
                id     : 'crawler',
                initial: 'init',

                states: {
                    // Initial setup state
                    'init': {
                        entry: (ctx) => { this.emit('init', ctx) },
                        on   : { 'next': { target: 'idle' } }
                    },
                
                    // Main idle state before starting the window
                    'idle': {
                        entry: (ctx) => { this.emit('idle', ctx) },
                        on   : { 'start': { target: 'setup' } }
                    },
                
                    // Prepare the crawler for scraping (by opening the window)
                    'setup': {
                        entry: (ctx) => { this.emit('setup', ctx) },
                        on   : { 'next': { target: 'ready' } }
                    },
                    
                    'ready': {
                        entry: (ctx) => { this.emit('ready', ctx) },
                        on   : { 
                            'scrape'  : { target: 'scraping' }, // Start scraping
                            'stop'    : { target: 'stop'     }, // Stop scraping
                            'error'   : { target: 'error'    },  // Error during scraping
                            'complete': { target: 'complete' }  // Complete scraping
                        }
                    },
                
                    'scraping': {
                        entry: (ctx) => { this.emit('scraping', ctx) },
                        on   : { 
                            'next'    : { target: 'ready'    }, // Scraping task done, back to idle for next
                            'stop'    : { target: 'stop'     }, // Stop scraping
                            'error'   : { target: 'error'    }, // Error during scraping
                            'complete': { target: 'complete' }  // Complete scraping
                        }
                    },
                
                    'stop': {
                        entry: (ctx) => { this.emit('stop', ctx) },
                        on   : { 'next': { target: 'complete' } } // Stop the process
                    },
                    
                    // Cleanup and return to idle
                    'complete': {
                        entry: (ctx) => { this.emit('complete', ctx) },
                        on   : { 'next': { target: 'idle' } }  // Return to idle after completing the process
                    },

                    'error': {
                        entry: (ctx) => { this.emit('error', ctx) },
                        on   : { 
                            'next' : { target: 'idle'  } // Go back to idle
                        }
                    },
                },
            })

            this.on('init', async (ctx) => { 
                await store.set('crawler.state', 'init')
                logger.info({state: 'init'})

                // Close window belonging to the crawler
                try       { await browser.windows.remove(await store.get('crawler.window')) }
                catch (e) {  }
                await store.set('crawler.window', null)

                await sleep(1_000)
                this.actor.send({ type: 'next' })
            })

            this.on('idle', async (ctx) => {
                await sleep(1_000)
                await store.set('crawler.state', 'idle')
                logger.info({state: 'idle'})
            })

            this.on('setup', async (ctx) => {
                await store.set('crawler.state', 'setup')
                logger.info({state: 'setup'})

                // Open window
                const window = await browser.windows.create({
                    url  : 'popup.html',
                    state: 'minimized',
                })

                await store.set('crawler.window', window.id)

                await sleep(1_000)
                this.actor.send({ type: 'next' })
            })

            this.on('ready', async (ctx) => {
                await store.set('crawler.state', 'ready')
                logger.info({state: 'ready'})

                // Set search_start
                await store.set('crawler.start_at', new Date().toISOString())

                // Monitor if window is closed
                this.timer_interval = setInterval(async (ctx) => {
                    try {
                        await browser.windows.get(await store.get('crawler.window'))
                    } catch (error) {
                        this.actor.send({ type: 'error', error: 'window_closed' })
                    }
                }, 1_000)

                this.timer_timeout = setTimeout(async () => {
                    this.actor.send({ type: 'error', error: 'timeout' })
                }, 15_000)
            })

            this.on('scraping', async (ctx) => {
                await store.set('crawler.state', 'scraping')
                logger.info({state: 'scraping'})
            
                clearTimeout(this.timer_timeout)
            
                try {
                    const tasks = ctx.event.tasks // Now expecting multiple tasks
                    const tab_promises = tasks.map(async (task) => {
                        const website = task.website
                        const keyword = task.keyword
            
                        const url = new URL(website.url)
                        url.searchParams.append(website.query, keyword)
                        url.searchParams.append('dse', '1')
                        url.searchParams.append('dse_keyword', keyword)
                        url.searchParams.append('dse_website', website.name)
            
                        const tab = await browser.tabs.create({ 
                            windowId: await store.get('crawler.window'),
                            url     : url.toString(),
                            active  : false, // Open in the background
                        })
            
                        // Wait for the 'upload' event to signal that the task is done
                        await new Promise<void>((resolve) => {
                            ipc.on('upload', async (ctx, data) => {
                                if (keyword != data.keyword) return
                                if (website.name != data.website) return
                                resolve()
                            })
                        })
            
                        try {
                            await browser.tabs.remove(tab.id)
                        } catch (e) {
                            logger.error('Failed to close tab', e)
                        }
                    })
            
                    await Promise.all(tab_promises)
            
                } catch (e) {
                    this.actor.send({ type: 'error', error: e })
                }
            
                await sleep(500)
                this.actor.send({ type: 'next' }) // Proceed to the next state
            })

            this.on('stop', async (ctx) => {
                await store.set('crawler.state', 'stop')
                logger.info({state: 'stop'})

                try       { await browser.windows.remove(await store.get('crawler.window')) }
                catch (e) {  }
                await store.set('crawler.window', null)
                
                await sleep(1_000)
                this.actor.send({ type: 'next' }) // Goes to complete
            })


            this.on('complete', async (ctx) => {
                await store.set('crawler.state', 'complete')
                logger.info({state: 'complete'})

                await store.set('crawler.completed_at', new Date().toISOString())
                clearInterval(this.timer_interval)
                clearTimeout(this.timer_timeout)

                try       { await browser.windows.remove(await store.get('crawler.window')) }
                catch (e) {  }

                await sleep(1_000)
                this.actor.send({ type: 'next' })
            })

            this.on('error', async (ctx) => {
                await store.set('crawler.state', 'error')
                logger.info({state: 'error'})
                logger.error(ctx.event.error)

                clearInterval(this.timer_interval)
                clearTimeout(this.timer_timeout)

                switch (ctx.event.error) {
                    case 'window_closed': 
                        logger.info('Interrupted')
                        break
                }

                try       { await browser.windows.remove(await store.get('crawler.window')) }
                catch (e) {  }
                await store.set('crawler.window', null)
                
                await sleep(1_000)
                this.actor.send({ type: 'next' })
            })

            await store.set('crawler.state', null)
            this.actor = createActor(this.machine)
            this.actor.start()
        }

        private async until(state: string): Promise<void> {
            await wait_until(async() => await store.get('crawler.state') === state)
        }

        // ------------------------------------------------------------
        // : Methods
        // ------------------------------------------------------------
        public async start() {
            this.actor.send({ type: 'start' })
            await this.until('ready')
        }

        public async stop() {
            this.actor.send({ type: 'stop' })
            await this.until('idle')
        }

        public async scrape(tasks: Task[]) {
            await this.until('ready')
            this.actor.send({ type: 'scrape', tasks: tasks })
            await this.until('scraping')
            await this.until('ready')
        }

        public async complete() {
            await this.until('ready')
            this.actor.send({ type: 'complete' })
            await this.until('idle')
        }
    }

    // ------------------------------------------------------------
    // : Exports
    // ------------------------------------------------------------
    export const crawler = new Crawler()
