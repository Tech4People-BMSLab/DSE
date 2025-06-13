// ------------------------------------------------------------
// : Imports    
// ------------------------------------------------------------
import localforage  from 'localforage'

import Bowser  from 'bowser'
import browser from 'webextension-polyfill'

import EventEmitter from 'eventemitter3'
import { DateTime } from 'luxon'

import { api }       from '@/background/core/api'

import { has }            from '@/background/utils/utils'
import { get }            from '@/background/utils/utils'
import { set }            from '@/background/utils/utils'
import { unset }          from '@/background/utils/utils'
import { is_empty }       from '@/background/utils/utils'
import { default_to}      from '@/background/utils/utils'
import { generate_token } from '@/background/utils/utils'
// ------------------------------------------------------------
// : Store
// ------------------------------------------------------------
export class Storage extends EventEmitter {
    private readonly storage = browser.storage.local
    public ready             = false

    constructor() {
        super()

        localforage.config({
            name     : 'dse',
            storeName: 'dse',
        })
    }
    
    public async init() {

        // Migration for v3.0.0
        const has_old_storage = has(await this.storage.get(), 'user.token')
        const has_new_storage = has(await this.get(), 'user.token')
        if (has_old_storage && !has_new_storage) {
            await this.migrate()
        }

        // Set defaults
        await this.set('user.token',           default_to(await this.get('user.token'), generate_token()))
        await this.set('crawler.state'       , default_to(await this.get('crawler.state')       , 'idle'))
        await this.set('crawler.started_at'  , default_to(await this.get('crawler.started_at')  , ''))
        await this.set('crawler.completed_at', default_to(await this.get('crawler.completed_at'), ''))

        const properties = Bowser.getParser(navigator.userAgent)
        await this.set('browser', {
            name      : properties.getBrowserName(),
            version   : properties.getBrowserVersion(),
            os        : properties.getOSName(),
            os_version: properties.getOSVersion(),
        })
        await this.set('extension.language', navigator.language)
        await this.set('extension.version' , browser.runtime.getManifest().version)

        browser.idle.onStateChanged.addListener(async (state) => {
            await this.set('extension.state', state)
        })
        browser.idle.queryState(60).then(async (state) => {
            await this.set('extension.state', state)
        })

        this.ready = true
    }

    public async migrate() {
        let current: any = await this.storage.get()

        await this.set('browser'           , get(current, 'browser', {}))
        await this.set('user.form'         , get(current, 'user.form', {}))
        await this.set('user.token'        , get(current, 'user.token', null))
        await this.set('user.popup'        , get(current, 'user.popup', false))
        await this.set('extension.language', get(current, 'extension.language', 'en'))
        await this.set('extension.version' , get(current, 'extension.version', '3.0.0'))
    }

    public async has(path: string): Promise<boolean> {
        const state = await this.get()
        return has(state, path)
    }

    public async set(path: string, value: any): Promise<void> {
        const state = await this.get()

        set(state, path, value)
        set(state, 'updated_at', DateTime.utc().toISO())

        await localforage.setItem('state', state)

        this.emit('update', path, value)
    }

    public async get(path?: string): Promise<any> {
        const state = default_to(await localforage.getItem('state'), {})
        if (!path) {
            return state
        } else {
            return get(state, path)
        }
    }

    public async remove(path: string): Promise<void> {
        const state = await this.get()

        unset(state, path)
        set(state, 'updated_at', DateTime.utc().toISO())

        await localforage.setItem('state', state)
    }

    public async clear(): Promise<void> {
        await localforage.clear()
    }
}
// ------------------------------------------------------------
// : Instance
// ------------------------------------------------------------
export const store = new Storage()
