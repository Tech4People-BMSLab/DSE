// ------------------------------------------------------------
// : Import
// ------------------------------------------------------------
import browser from 'webextension-polyfill'

// ------------------------------------------------------------
// : Utils
// ------------------------------------------------------------
export const is_empty     = (obj: any) => [Object, Array].includes((obj || {}).constructor) && !Object.entries((obj || {})).length
export const is_undefined = (obj: any) => obj === undefined
export const is_null      = (obj: any) => obj === null
export const is_object    = (obj: any) => obj !== null && typeof obj === 'object' && !Array.isArray(obj)
export const is_array     = (obj: any) => Array.isArray(obj)
export const is_string    = (obj: any) => typeof obj === 'string'
export const is_number    = (obj: any) => typeof obj === 'number'
export const is_boolean   = (obj: any) => typeof obj === 'boolean'

// ------------------------------------------------------------
// : Logger
// ------------------------------------------------------------
export class Logger {
    private id: string
    private name: string

    constructor(name: string) {
        this.id   = name.toLowerCase()
        this,name = name
    }

    async log(message: any) {
        try {
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
        } catch (e) {
            console.error(e)
        }
    }
}
