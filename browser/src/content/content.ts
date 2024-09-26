// ------------------------------------------------------------
// : Import
// ------------------------------------------------------------
import browser from 'webextension-polyfill'

import { Storage } from '@/content/utils/storage'

import { sleep }      from '@/content/utils/utils'
import { wait_until } from '@/content/utils/utils'

import      { connect }                        from 'nats.ws'
import type { NatsConnection, NatsError, Msg } from 'nats.ws'
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
let storage    : Storage
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

function is_registration() {
    let url = new URL(window.location.href)

    switch (true) {
        case url.hostname == "localhost":
        case url.hostname == 'dev.bmslab.utwente.nl': return true

        default: return false
    }
}

function has_clean() {
    const clean = new URLSearchParams(window.location.search).get('dse_clean')
    return clean === '1'
}

async function is_ready(count = 3, interval = 333) {
    let sizes = Array(count).fill(0)

    const sleep         = ms  => new Promise(resolve => setTimeout(resolve, ms))
    const get_byte_size = str => new Blob([str]).size

    while (true) {
        const html_str = document.documentElement.outerHTML // Get the entire HTML as a string
        const size     = get_byte_size(html_str)            // Calculate the byte size of the HTML

        // Push the new size into the sizes array.
        sizes.push(size)
        log(`DSE: Waiting for page ready... (${sizes})`)

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

    log(`DSE: Ready`)
}

function clean() {
    try {
        log(`DSE: Cleaning`)
        const body_clone = document.body.cloneNode(true) as HTMLElement
            
        // Remove all unwanted elements and attributes
        const remove_elements = ['svg', 'script', 'iframe', 'style', 'img']
        remove_elements.forEach((tag : string) => {
            const elements = body_clone.querySelectorAll(tag)
            elements.forEach(el => el.remove())
        })
    
        // Remove all jsdata and jsaction attributes
        const all_elements = body_clone.querySelectorAll('*')
        all_elements.forEach(el => {
            // el.removeAttribute('jsdata') 
            // el.removeAttribute('jsaction')
        })

        document.documentElement.replaceChild(body_clone, document.body) // Set the HTML to the cleaned, cloned body
        log(`DSE: Cleaned`)
    } catch (e) {
        // Ignore
    }
}

function log(...args) {
    browser.runtime.sendMessage({type: 'log', from: 'tab', data: args})
}

// ------------------------------------------------------------
// : Search Capture
// ------------------------------------------------------------

// ------------------------------------------------------------
// : Extractor
// ------------------------------------------------------------
class Extractor {

    public async init() {
        await is_ready()

        // Clean the page if requested
        if (has_clean()) {
            log(`DSE: Cleaning`)
            clean()
        }

        // Send data to server
        try       { await uploader.upload() }
        catch (e) {  }

        log(`DSE: Done`)
    }

}

// ------------------------------------------------------------
// : Registrator
// ------------------------------------------------------------
class Registrator {

    public async init() {
        log(`DSE: Registration loaded`)
        await this.wait_for_export()
        log(`DSE: Export found`)

        const element = document.querySelector('.export')

        // Get form data from page
        let data
        data = element.innerHTML
        data = JSON.parse(data)

        log(data)

        // Check if backup exists
        const backup = await storage.get('user.form_backup')

        if (!backup) {
            await storage.set({'user.form_backup': data})
        }

        // Set current data in user.form
        await storage.set({'user.update': new Date().toISOString()})
        await storage.set({'user.form'  : data})

        // Close page
        // window.close()
    }

    async wait_for_export() {
        await wait_until(() => document.querySelector('.export') !== null)
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

        const html = document.documentElement.outerHTML
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
                timestamp: new Date().toISOString(),
    
                token: await storage.get('user.token'),
                url  : window.location.href,
    
                browser     : await storage.get('browser'),
                localization: window.navigator.language,
    
                keyword: keyword,
                website: website,
    
                html: html,
            }

            browser.runtime.sendMessage({type: 'action', from: 'tab', action: 'upload', data: payload})

            // // Convert to binary using msgpack
            // const data = pack(payload)
            // const size = data.length

            log(`Uploaded (size: ${size})`)
        } catch (e) {
            log('Error', e)
        }
    }
}

// ------------------------------------------------------------
// : Initialization
// ------------------------------------------------------------
storage     = new Storage()
uploader    = new Uploader()
extractor   = new Extractor()
registrator = new Registrator()

// ------------------------------------------------------------
// : Main
// ------------------------------------------------------------
async function main() {
    try {
        switch (true) {
            case is_dse()         : await extractor.init()  ; break
            case is_registration(): await registrator.init(); break
            default: return
        }
    } catch (e) {
        console.error(e)
    }

}

main()
