// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------
import browser      from 'webextension-polyfill'
import EventEmitter from 'eventemitter3'

import { api }    from '@/background/core/api'
import { store }  from '@/background/core/storage'
import { Logger } from '@/background/utils/logger'

import { default_to } from '../utils/utils'
// ------------------------------------------------------------
// : IPC
// ------------------------------------------------------------
class IPC extends EventEmitter {

    private loggers = {
        'tab'  : new Logger('Tab'),
        'popup': new Logger('Popup')
    }

    async init() {
        browser.runtime.onMessage.addListener((packet: any, sender: any) => {
            if (packet.event === 'log')      { this.emit('log'     , packet.from, packet.data); return }
            if (packet.event === 'upload')   { this.emit('upload'  , packet.from, packet.data); return }
            if (packet.event === 'register') { this.emit('register', packet.from, packet.data); return }
            if (packet.event === 'consent')  { this.emit('consent' , packet.from, packet.data); return }
        })


        this.on('log', (from, data) => {
            let logger
            if (from === 'tab')   logger = this.loggers['tab']
            if (from === 'popup') logger = this.loggers['popup']

            logger.info(data.join(' '))
        })

        this.on('upload', async (from, data) => {
            api.send('upload', {
                timestamp: new Date().toISOString(),
                token    : await store.get('user.token'),

                url: data.url,
                browser: await store.get('browser'),
                keyword: data.keyword,
                website: data.website,
                localization: data.localization,

                html: data.html,
            })
        })

        this.on('register', async (from, data) => {
            const timestamp = new Date().toISOString()
            const history = default_to(await store.get('user.consent.history'), [])
            const current = default_to(await store.get('user.consent.current'), {
                timestamp: timestamp,
                raw : data.raw,
                data: data.form,
            })

            history.push(current)

            await store.set(`user.form`, data.form)
            await store.set('user.update', timestamp)
            await store.set(`user.consent.current`, current)
            await store.set(`user.consent.history`, history)
        })

        this.on('consent', async (from, data) => {
            await browser.windows.create({
                url : 'https://static.33.56.161.5.clients.your-server.de/dse/consent',
                type: 'popup',
            })
        })
    }

}

// ------------------------------------------------------------
// : Exports
// ------------------------------------------------------------
export const ipc = new IPC()
