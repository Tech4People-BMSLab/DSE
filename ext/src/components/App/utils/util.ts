// ------------------------------------------------------------
// : Import
// ------------------------------------------------------------
import browser           from 'webextension-polyfill'
import EventEmitter      from 'eventemitter3'

import { DateTime }      from 'luxon'

import { get }                         from 'lodash'
import { set }                         from 'lodash'
import { find }                        from 'lodash'
import { isUndefined as is_undefined } from 'lodash'
import { isNull      as is_null}       from 'lodash'
import { isObject    as is_object}     from 'lodash'
import { isArray     as is_array}      from 'lodash'
import { isString    as is_string}     from 'lodash'
import { isNumber    as is_number}     from 'lodash'
import { isBoolean   as is_boolean}    from 'lodash'
import { isFunction  as is_function}   from 'lodash'


import { Mutex }     from 'async-mutex'


// ------------------------------------------------------------
// : Enum
// ------------------------------------------------------------
export enum State {
    ERROR = 'error',

    IDLE = 'idle',
    BUSY = 'busy',

    CANCELLED = 'cancelled',
}

// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
export const sleep = async (ms: number) => {
    return new Promise(resolve => setTimeout(resolve, ms))
}

export const wait_until = async (predicate, timeout=-1, interval=100) => {
    const time_start = Date.now()
    while (true) {
        const result = await predicate()
        if (result) {
            break
        }
        if (timeout > 0 && Date.now() - time_start > timeout) {
            break
        }
        await sleep(interval)
    }
}

export const all = (arr: Array<any>, predicate: (item: any) => boolean) => {
    return arr.filter(predicate).length === arr.length
}

export const any = (arr: Array<any>, predicate: (item: any) => boolean) => {
    return arr.filter(predicate).length > 0
}

// ------------------------------------------------------------
// : Utils
// ------------------------------------------------------------
// From https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore
export const is_empty = obj => [Object, Array].includes((obj || {}).constructor) && !Object.entries((obj || {})).length;
export const contains = (arr, value) => arr.includes(value)
export const has      = (obj, key) => {
    var keyParts = key.split('.');
    const result = !!obj && (
        keyParts.length > 1
            ? has(obj[key.split('.')[0]], keyParts.slice(1).join('.'))
            : Object.hasOwnProperty.call(obj, key)
    )
    return result
}


// ------------------------------------------------------------
// : Logging
// ------------------------------------------------------------
export async function log(message: any) {

    if (is_object(message)) {
        message = JSON.stringify(message)
        message = `object<${message}>`
    }

    if (is_array(message)) {
        message = JSON.stringify(message)
        message = `array<${message}>`
    }

    if (is_boolean(message))   {message = message ? 'true' : 'false'}
    if (is_null(message))      {message = 'null'                    }
    if (is_undefined(message)) {message = 'undefined'               }

    await browser.runtime.sendMessage({
        from: 'logger',
        data: message
    })
}

export async function report(stack, ...args) {
    const api_url  = `${await storage.get('api_url')}/log`

    // Clean the stack trace
    const stack_lines = stack.split('\n').slice(1)

    // Get user token
    const token = await storage.get('user.token')

    // fetch(api_url, {
    //     method: 'POST',
    //     headers: {
    //         'content-type': 'application/json'
    //     },
    //     body: JSON.stringify({
    //         user : token,
    //         stack: stack_lines,
    //         log  : args
    //     })
    // }).then(response => {
    //     if (!response.ok) {
    //         const reason = response.text()
    //         console.error(`Failed to report to server (reason: ${reason})`)
    //     }
    //     return
    // }).catch(error => {
    //     console.error(`Failed to report to server (reason: ${error})`)
    // })
}
// ------------------------------------------------------------
// : Storage
// ------------------------------------------------------------
export class Storage {
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
        return get(await this.storage.get(), key)
    }
    
    /**
     * Set a value in the storage
     * @param obj The object to set the values from
     */
    public async set(obj: object) {
        await this.storage.set(obj)
    }

    /**
     * Remove a value from the storage
     * @param key The key to remove from the storage
     */
    public async remove(key: string) {
        await this.storage.remove(key)
    }

    public async get_all() {
        return await this.storage.get()
    }
}

// ------------------------------------------------------------
// : API Functions
// ------------------------------------------------------------
export async function is_registered() {
    return await storage.get('user.registered') === true
}

export async function verify_user(): Promise<boolean> {
    return await storage.get('user.registered') === true
}

export async function get_search_state(): Promise<State> {
    const result = await storage.get('search_state')
    const state  = result.search_state
    return state
}

// ------------------------------------------------------------
// : Global
// ------------------------------------------------------------
export let storage = new Storage()

// ------------------------------------------------------------
// : Debug
// ------------------------------------------------------------