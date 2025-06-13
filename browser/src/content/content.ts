// ------------------------------------------------------------
// : Import
// ------------------------------------------------------------
import browser from 'webextension-polyfill'

import { wait_until } from '@/content/utils/utils'
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
let uploader   : Uploader
let extractor  : Extractor
let registrator: Registrator
// ------------------------------------------------------------
// : Helpers
// ------------------------------------------------------------
function is_dse() {
    const dse = new URLSearchParams(window.location.search).get('dse')
    return dse === '1'
}

function is_consent() {
    let url = new URL(window.location.href)

    switch (true) {
        case url.hostname == "localhost":
        case url.hostname == 'static.33.56.161.5.clients.your-server.de': return true

        default: return false
    }
}

async function is_ready(count = 3, interval = 333) {
    let sizes = Array(count).fill(0)

    const sleep         = ms  => new Promise(resolve => setTimeout(resolve, ms))
    const get_byte_size = str => new Blob([str]).size

    log(JSON.stringify({state: 'waiting'}))

    while (true) {
        const html_str = document.documentElement.outerHTML // Get the entire HTML as a string
        const size     = get_byte_size(html_str)            // Calculate the byte size of the HTML

        // Push the new size into the sizes array.
        sizes.push(size)

        // Keep the sizes array size consistent to 'count'
        if (sizes.length > count) {
            sizes.shift()
        }

        // Check if all sizes in the array are equal (page content is stable)
        if (sizes.every(s => s === sizes[0])) {
            break
        }
        
        await sleep(interval)
    }

    log(JSON.stringify({state: 'ready'}))
}

function log(...args) {
    browser.runtime.sendMessage({event: 'log', from: 'tab', data: args})
}
// ------------------------------------------------------------
// : Registrator
// ------------------------------------------------------------
class Registrator {

    public async init() {
        await this.wait_for_export()
        log(`DSE: Export found`)

        let data
        data = document.querySelector('#export').innerHTML
        data = JSON.parse(data)

        const raw  = data.raw
        const form = data.form

        browser.runtime.sendMessage({event: 'register', from: 'tab', data: {
            raw : raw,
            form: form,
        }})

        log(`DSE: Registered`)
    }

    async wait_for_export() {
        await wait_until(() => document.querySelector('#export') !== null)
    }
}
// ------------------------------------------------------------
// : Extractor
// ------------------------------------------------------------
class Extractor {

    public async init() {
        await is_ready()

        // Send data to server
        try       { await uploader.upload() }
        catch (e) {  }

        log(`DSE: Done`)
    }

}
// ------------------------------------------------------------
// : Uploader
// ------------------------------------------------------------
class Uploader {

    public async upload() {
        const keyword = new URLSearchParams(window.location.search).get('dse_keyword')
        const website = new URLSearchParams(window.location.search).get('dse_website')
        log(`DSE: ${website}:${keyword}`)

        const html = document.body.outerHTML
        const size = new Blob([html]).size // Get the size of html in bytes
        if (size > 10_000_000) { // 10MB
            console.error('Uploader: File too large')
            return
        }
        log(`Website: ${website}:${keyword} (${size} bytes)`)

        // Send data
        try {

            // Create payload
            const payload = {
                version  : 1,
                timestamp: new Date().toISOString(),
    
                url  : window.location.href,
                localization: window.navigator.language,
    
                keyword: keyword,
                website: website,
    
                html: html,
            }

            browser.runtime.sendMessage({event: 'upload', from: 'tab', data: payload})
            log(`Uploaded (size: ${size})`)
        } catch (e) {
            log('Error', e)
        }
    }
}
// ------------------------------------------------------------
// : Initialization
// ------------------------------------------------------------
uploader    = new Uploader()
extractor   = new Extractor()
registrator = new Registrator()

// ------------------------------------------------------------
// : Main
// ------------------------------------------------------------
async function main() {
    try {
        switch (true) {
            case is_dse()    : await extractor.init()  ; break
            case is_consent(): await registrator.init(); break
            default: return
        }
    } catch (e) {
        console.error(e)
    }

}

main()
