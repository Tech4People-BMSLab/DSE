// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------
import browser from 'webextension-polyfill'
import EventEmitter from 'eventemitter3'

import { api }    from '@/background/core/api'
import { store }  from '@/background/core/storage'
import { Logger } from '@/background/utils/logger'
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
        })
    }

    async send(event: string, data: any) {
        browser.runtime.sendMessage({event, from: 'tab', data})
    }

    async log(...args) {
        browser.runtime.sendMessage({event: 'log', from: 'tab', data: args})
    }
}
// ------------------------------------------------------------
// : Exports
// ------------------------------------------------------------
export const ipc = new IPC()
