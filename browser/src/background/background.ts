// ------------------------------------------------------------
// : Import
// ------------------------------------------------------------
// Browser (polyfill)
import browser from 'webextension-polyfill'

// General
import EventEmitter                   from 'eventemitter3'
import { DateTime }                   from 'luxon'
import { createMachine, createActor } from 'xstate'

// Modules
import { Logger }  from '@/background/utils/logger'
import { api }     from '@/background/core/api'
import { ws }      from '@/background/core/ws'
import { ipc }     from '@/background/core/ipc'
import { crawler } from '@/background/core/crawler'

// Utils
import { store }          from '@/background/core/storage'
import { not_empty }      from '@/background/utils/utils'
import { wait_until }     from '@/background/utils/utils'
// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
declare global {
    interface ImportMeta {
        env: Record<string, string>
    }
}
declare var self  : any
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
const logger = new Logger('Background')
// ------------------------------------------------------------
// : Background
// ------------------------------------------------------------
class Background extends EventEmitter {
    private state: ReturnType<typeof createMachine>
    private actor: ReturnType<typeof createActor>

    /**
     * Initial point of the script.
     */
    public async init() {

        this.state = createMachine({
            /** @xstate-layout N4IgpgJg5mDOIC5QCMCGBjA1lATgewFcA7CAOljABcCAHAYgIpwFoAzASx1koG0AGALqJQNPLHaV2eIsJAAPRACYAjH1IBOAGwAOAOwBmdeoCsmvsb67dAGhABPRKoAspJ0fXbty47qfbNipoAvkG2aFi4hCTkVLQMTMzEOGBQ7NxgyRD8QkggouKS0rIKCCpqWnqGJmYWVrYOpU66pIpuHup8iuaWyoohYRjY+MRkFNT0jBnMyanpmdmy+RJSMrklZRo6Bkam3XX2SurNuhZ83rqayk4B1-0g4UNRZJMsHFyUdERgcryCi2LLIprRCmFzGfTKZTucHaYzGRT1Rx8FxtDxeHx+AL6O4PSIjUgvRJEGZpSgZSCfb6-HIiAGFVagEraJzGUhXcwGPj+ZkGREISHqUjGbT6fROcX6cytUw4wZ46KEklzClfH4LXJLenFRyqZp+Jz6RQqE66bR85Ri0iaa2aU3mbwQvwhUIgIh4CBwWS44Ykf4FFbahDMTR84OkPgRyNRqO6WURH2jWI0P2AhnyRBOBEHfmdUh6U6qLQssXqOOPfGEt7cFNa4EIJzmIU3YzuO3ivleIXR61uA18YIu71PAkJJIpUnkiA1gN1zQtoXXLzWo4Wzx8rRCkWizOwkwR7GDuUJkdTJVkzLToGMxDqBts0zKbQeSUWvjqc2Qq0220dTTVI5lvKZDJKgEANLS-pXum-JdIorjuPoT5Qv2+g2NmkLKF+Nq6L+-6xs6QA */
            id     : 'background',
            initial: 'init',

            states: {
                'init': {
                    entry: (ctx) => { this.emit('init', ctx) },
                    on   : {
                        'next': { target: 'setup' }
                    }
                },

                'setup': {
                    entry: (ctx) => { this.emit('setup', ctx) },
                    on   : {
                        'user-0': { target: 'user-0' },
                        'user-1': { target: 'user-1' },
                        'user-2': { target: 'user-2' },
                    }
                },

                'user-0': { // Newcomer
                    entry: (ctx) => { this.emit('user-0', ctx) },
                    on   : { 'next': { target: 'user-1'} }
                },
                'user-1': { // Unregistered user
                    entry: (ctx) => { this.emit('user-unregistered', ctx) },
                    on   : { 'next': { target: 'user-2'} }
                },
                'user-2': { // Registered user
                    entry: (ctx) => { this.emit('user-2', ctx) },
                    on   : { 'next': { target: 'ready'} }
                },

                'ready': {
                    entry: (ctx) => { this.emit('ready', ctx) }
                },
            },
        })

        this.on('init', async() => {
            await store.set('background.state', 'init')
            const is_development = import.meta.env.MODE === 'development'

            // Initialize storage
            await store.init()

            // Log the start of the extension
            const now = DateTime.local()
            logger.info(`Started at ${now.toFormat('yyyy-MM-dd HH:mm:ss')}`)
            logger.info(`Version: ${await store.get('extension.version')}`)
            logger.info(`Browser: ${await store.get('browser.name')} ${await store.get('browser.version')}`)
            logger.info(`Token  : ${await store.get('user.token')}`)
            logger.info(`API    : ${import.meta.env.BASE_URL}`)
            logger.info(`WS     : ${import.meta.env.BASE_WS}`)

            self.reload = async function () {
                browser.runtime.reload()
            }

            self.reset = async function () {
                await api.send('reset', {})
            }

            self.register = async function () {
                browser.runtime.getURL('/index.html')
            }

            self.state = async function () {
                console.log(await store.get())
            }

            if (is_development) {
                logger.info('Development mode enabled')

                // await browser.windows.create({
                //     // url : 'https://static.33.56.161.5.clients.your-server.de/dse/consent',
                //     url : 'http://localhost:3000/consent',
                //     type: 'popup',
                // })

                self.register = async function () {
                    await browser.windows.create({
                        url : 'https://static.33.56.161.5.clients.your-server.de/dse/consent',
                        type: 'popup',
                    })
                }

                self.update = async function () {
                    api.emit('update')
                }

                //// Load popups
                // setTimeout(() => {
                //     browser.tabs.create({
                //         url: 'http://google.com'
                //     })
                // }, 1000)

                //// Opens the pop up in a new tab (for fast debugging)
                // const url = browser.runtime.getURL('/popup.html')
                // const tab = await browser.tabs.create({url: url})

                //// Reload extension periodically
                // setInterval(() => {
                //     browser.runtime.reload()
                // }, 2500)
            }

            this.actor.send({ type: 'next' })
        })

        this.on('setup', async() => {
            await store.set('background.state', 'setup')

            // Initialize modules
            crawler.init()
            ws     .init()
            ipc    .init()
            api    .init()

            switch (true) {
                case await store.has('user.form') : this.actor.send({ type: 'user-2' }); break
                case await store.has('user.popup'): this.actor.send({ type: 'user-1' }); break
                default                           : this.actor.send({ type: 'user-0' }); break
            }
        })

        // Newcomer
        this.on('user-0', async() => {
            await store.set('background.state', 'user-0')

            await browser.windows.create({
                url : 'https://static.33.56.161.5.clients.your-server.de/dse/consent',
                type: 'popup',
            })

            await store.set('user.popup', true)
            await store.set('user.type' , 'newcomer')
            this.actor.send({ type: 'next' })
        })

        // Unregistered user
        this.on('user-1', async() => {
            await store.set('background.state', 'user-1')

            await wait_until(async() => await store.get('user.form'))
            await store.set('user.popup', true)
            await store.set('user.type' , 'unregistered')
            this.actor.send({ type: 'next' })
        })

        // Registered user
        this.on('user-2', async() => {
            await store.set('background.state', 'user-2')

            await store.set('user.popup', true)
            await store.set('user.type' , 'registered')
            this.actor.send({ type: 'next' })
        })

        this.on('ready', async() => {
            logger.info('Ready')
            await store.set('background.state', 'ready')
        })

        this.actor = createActor(this.state)
        this.actor.start()

        setTimeout(() => {
            // TODO: Put this back
            browser.runtime.reload()
        }, 3600_000)
    }
}
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
const background = new Background()

// ------------------------------------------------------------
// : Main
// ------------------------------------------------------------
background.init()
