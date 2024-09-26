// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------

// Lodash
import _ from 'lodash'

// Crypto
import cryptoRandomString from 'crypto-random-string'

// ------------------------------------------------------------
// : Lodash
// ------------------------------------------------------------
export const is_empty  = _.isEmpty
export const not_empty = _.negate(is_empty)

// ------------------------------------------------------------
// : Utils
// ------------------------------------------------------------
// https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore
export async function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms))
}

export async function wait_until(predicate, timeout=-1, interval=100) {
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

export function format_bytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes'
    const k     = 1024
    const dm    = decimals < 0 ? 0 : decimals
    const sizes = ['Bytes', 'KB', 'MB']
    const i     = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]

}

export function generate_token() {
    return cryptoRandomString({length: 12, characters: '0123456789abcdef'})
}
