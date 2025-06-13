// ------------------------------------------------------------
// : Import
// ------------------------------------------------------------
import browser from 'webextension-polyfill'

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
export const is_empty     = (obj: any) => [Object, Array].includes((obj || {}).constructor) && !Object.entries((obj || {})).length
export const is_undefined = (obj: any) => obj === undefined
export const is_null      = (obj: any) => obj === null
export const is_object    = (obj: any) => obj !== null && typeof obj === 'object' && !Array.isArray(obj)
export const is_array     = (obj: any) => Array.isArray(obj)
export const is_string    = (obj: any) => typeof obj === 'string'
export const is_number    = (obj: any) => typeof obj === 'number'
export const is_boolean   = (obj: any) => typeof obj === 'boolean'
export const not_empty    = (obj: any) => !is_empty(obj)

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

export const get = (obj, path, defaultValue = undefined) => {
    const travel = regexp =>
        String.prototype.split
        .call(path, regexp)
        .filter(Boolean)
        .reduce((res, key) => (res !== null && res !== undefined ? res[key] : res), obj)
    const result = travel(/[,[\]]+?/) || travel(/[,[\].]+?/)
    return result === undefined || result === obj ? defaultValue : result
}


// ------------------------------------------------------------
// : Logging
// ------------------------------------------------------------
export class Logger {
    private id: string
    private name: string

    constructor(name: string) {
        this.id   = name.toLowerCase()
        this,name = name
    }

    async log(message: any) {

        switch (true) {
            case is_string(message): {
                message = JSON.stringify(message)
                message = `object<${message}>`
                break
            }
            
            case is_array(message): {
                message = JSON.stringify(message)
                message = `array<${message}>`
                break
            }

            case is_boolean(message): {
                message = message ? 'true' : 'false'
                break
            }

            case is_null(message): {
                message = 'null'
                break
            }

            case is_undefined(message): {
                message = 'undefined'
                break
            }
        }

        await browser.runtime.sendMessage({
            from: `logger.${this.id}`,
            data: message
        })
    }
}

export async function report(stack, ...args) {
    // const api_url  = `${await storage.get('api_url')}/log`

    // Clean the stack trace
    // const stack_lines = stack.split('\n').slice(1)

    // Get user token
    // const token = await storage.get('user.token')

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
