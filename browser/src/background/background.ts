// ------------------------------------------------------------
// : Import
// ------------------------------------------------------------

// Browser (polyfill)
import browser from 'webextension-polyfill'
import Bowser  from 'bowser'

// General
import axios   from 'axios'
import { DateTime }      from 'luxon'
import { createMachine } from '@xstate/fsm'
import { interpret }     from '@xstate/fsm'
import _ from 'lodash'

// Modules
import { Logger }         from '@/background/utils/logger'
import { window_manager } from '@/background/utils/window'
import { crawler }        from '@/background/crawler/crawler'
import { ipc }            from '@/background/modules/ipc'

// Utils
import { storage }        from '@/background/utils/storage'
import { generate_token } from '@/background/utils/utils'
import { wait_until }     from '@/background/utils/utils'
// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
declare global {
    interface ImportMeta {
        env: Record<string, string>
    }
}

declare var self: any
// ------------------------------------------------------------
// : Utils
// ------------------------------------------------------------
const is_empty  = _.isEmpty
const not_empty = _.negate(is_empty)

const pick = _.pick
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
const logger = new Logger('Background')
// ------------------------------------------------------------
// : Background
// ------------------------------------------------------------
class Background {
    private state  : ReturnType<typeof createMachine>
    private service: any

    /**
     * Initial point of the script.
     */
    public async init() {
        this.state = createMachine({
            id     : 'background',
            initial: 'debug',

            states: {
                'debug': {
                    entry: () => { this.on_debug() },
                    on   : {'next': {target: 'init'}}
                },

                'init': {
                    entry: () => { this.on_init() },
                    on   : {'next': {target: 'startup'}}
                },

                'startup': {
                    entry: () => { this.on_startup() },
                    on   : {'verify': { target: 'verifying' }},
                },

                'verifying': {
                    entry: () => { this.on_verifying() },
                    on: {
                        'not_registered': {target: 'registering'},
                        'registered'    : {target: 'ready'      },
                    },
                },
                
                'registering': {
                    entry: () => { this.on_registering() },
                    on   : {'ready' : {target: 'verifying'}},
                },

                'ready': {
                    entry: () => { this.on_ready() }
                },
          },
        })

        this.service = interpret(this.state).start()
        this.service.subscribe((state) => {
            
        })
    }

    private async on_debug() {
        const is_development = import.meta.env.MODE === 'development'
        await storage.set({'environment': is_development ? 'dev' : 'prod'})

        /// Define debugging functions
        {
            self.restart = async function() {
                browser.runtime.reload()
            }
    
            self.register = async function() {
                browser.tabs.create({
                    url: 'http://dev.bmslab.utwente.nl/dse/consent'
                })
            }

            self.form = async function() {
                console.log(await storage.get('user.form'))
            }

            self.token = async function() {
                console.log(await storage.get('user.token'))
            }

            self.storage = async function() {
                console.log(await storage.get_all())
            }

            self.start = async function() {
                const token = await storage.get('user.token')

                await axios({
                    method : 'get',
                    baseURL: import.meta.env.BASE_URL,
                    url    : `/api/users/${token}/start`,
                    params : {token: 'dse2023'}
                }).catch((e) => {
                    logger.error(e)
                })
            }

            self.stop = async function() {
                const token = await storage.get('user.token')

                await axios({
                    method : 'get',
                    baseURL: import.meta.env.BASE_URL,
                    url    : `/api/users/${token}/stop`,
                    params : {token: 'dse2023'}
                }).catch((e) => {
                    logger.error(e)
                })
            }
        }

        if (is_development) {
            logger.log('Development mode enabled')

            self.register = async function() {
                browser.tabs.create({
                    url: 'http://localhost:3000/consent'
                })
            }

            //// Clear token (new user)
            // await storage.clear()

            //// Clear form (old versions)  
            // await storage.remove('user.token')
            // await storage.remove('user.form')
            // await storage.set({'user.form': {}})

            //// Registered user (Home server)
            // await storage.set({'user.token': '87edb4dbf4da'})

            //// Registered user (Fake token)
            // await storage.set({'user.token': '12345'})

            //// Registered user (BMS server)
            // await storage.set({'user.token': 'e20eba0f02fa023a3bbe2eb8ec9bb9a06436c3398d2006f4a9817967607af7a5'})

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

        // Wait until the service is ready
        await wait_until(() => this.service)

        // Go to the next state
        this.service.send('next')
    }

    //// State Handlers
    private async on_init() {
        // Generate token first thing
        let has_token = !!await storage.get('user.token')

        if (!has_token) { 
            await storage.set({'user.token': generate_token()})
        }

        // Initialize modules
        await ipc    .init()
        await crawler.init()



        // Listen for changes in storage
        let last_update = DateTime.local()
        storage.onChanged.addListener(async () => {
            if (last_update.plus({seconds: 5}) > DateTime.local()) {
                return
            }

            ipc.publish('bms.dse.v1.5.client.update', JSON.stringify({
                token   : await storage.get('user.token'),
                form    : await storage.get('user.form'),
                storage : await storage.get_all(),
                manifest: browser.runtime.getManifest(),
            }))
        })

        this.service.send('next')
    }

    private async on_startup() {
        // Register IPC handlers
        ipc.register()

        // Filter out unnecessary data (due to previous extension versions)
        {
            const data = await storage.get_all()
            const filtered = pick(data, [
                'browser',
                
                'connected',
                'language',

                'user.form',
                'user.registered',
                'user.registered_timestamp',
                'user.token',

                'window_manager.tabs',
                'window_manager.windows',

                'extension.popup_timestamp',
                'extension.state',
                'extension.version',
                'extension.installation_time',
            ])
            
            await storage.set(filtered)
        }

        // Remove config from storage (latest version)
        try       { await storage.remove('config') }
        catch (e) {  }

        await storage.set({'extension.version': browser.runtime.getManifest().version}) // Set the extension version

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
        const now   = DateTime.local()
        logger.log(`Started at ${now.toFormat('yyyy-MM-dd HH:mm:ss')}`)
        logger.log(`Version: ${browser.runtime.getManifest().version}`)
        logger.log(`Browser: ${browser_properties.getBrowserName()} ${browser_properties.getBrowserVersion()}`)
        logger.log(`Token  : ${await storage.get('user.token')}`)
        logger.log(`IPC    : ${import.meta.env.NATS_URL} (${import.meta.env.NATS_USER})`)


        // Set handler for idle state changes
        browser.idle.onStateChanged.addListener(async(state) => {
            logger.log(`Idle state: ${state}`)
            await storage.set({'extension.state': state})
        })

        // Set idle detection 
        browser.idle.queryState(60).then(async (state) => {
            logger.log(`Idle state: ${state}`)
            await storage.set({'extension.state': state})
        })

        this.service.send('verify')
    }

    private async on_verifying() {
        try {
            const token = await storage.get('user.token')
            const form  = await storage.get('user.form')

            const has_token = not_empty(token)
            const has_form  = not_empty(form)

            switch (true) {
                case has_token && has_form: {
                    ipc.publish('bms.dse.user.register', JSON.stringify({
                        token: token,
                        form : form,
                    }))

                    // Jumps to on_ready
                    this.service.send('registered') 
                    break
                }

                default: {
                    // Jumps to on_registering
                    this.service.send('not_registered') 
                    break
                }
            }
        } catch (e) {
            logger.error(e)
            
            // Jumps to on_registering
            this.service.send('not_registered') 
        }
    }

    private async on_registering() {    
        try {
            await wait_until(() => ipc.connected) // Wait until connected to the server

            // Skip if there is registered_timestamp
            if (await storage.get('user.registered_timestamp')) {
                this.service.send('ready')
                return
            }

            // Open the registration popup
            window_manager.open_popup({minimized: false})

            // Set current time in storage
            const now = DateTime.now()
            await storage.set({'user.registered_timestamp': now.toISO()})
            
            await wait_until(() => storage.get('user.form')) // Wait until the user registers
            this.service.send('ready')                       // Jumps to on_ready
        } catch (e) {
            logger.error(e)
        }
    }

    private async on_ready() {
        logger.log('Ready')
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
