// ------------------------------------------------------------
// : Import
// ------------------------------------------------------------
import browser from 'webextension-polyfill'

import { isEmpty as is_empty } from 'lodash'

import { log }     from '../../components/App/utils/util'
import { Storage } from '../../components/App/utils/util'

// ------------------------------------------------------------
// : Utilities
// ------------------------------------------------------------
const intersection = (a) => {
    return a.reduce((a, b) => a.filter(c => b.includes(c)))
}

async function sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms))
}
async function wait_until(predicate, timeout=-1, interval=100) {
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

/**
 * Extract elements using XPath format instead of CSS selector.
 * @param xpath The XPath query
 * @param parent The document to search in
 * @returns {Array} The elements found
 */
function query(xpath, parent) {
    let results = []
    let query = document.evaluate(xpath, parent || document,
        null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null)
    for (let i = 0, length = query.snapshotLength; i < length; ++i) {
        results.push(query.snapshotItem(i))
    }
    return results;
}

function query_text(xpath, parent) {
    try {
        const items = query(xpath, parent)

        if (items.length == 0) {
            return ''
        } else if (items.length == 1) {
            return items[0].textContent
        } else if (items.length > 1) {
            return items.map(item => item.textContent).join(' ')
        }
    } catch (error) {
        log(`${xpath} not found (reason: ${error})`)
        return null
    }
}

async function query_all(xpath, parent = null): Promise<HTMLElement[]> {
    parent = parent || document.querySelector('html')

    function query(xpath, parent) {
        let results = []
        let query = document.evaluate(xpath, parent || document,
            null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null)
        for (let i = 0, length = query.snapshotLength;i < length;++i) {
            results.push(query.snapshotItem(i))
        }
        return results
    }

    return await new Promise(resolve => {
        const elements = query(xpath, parent) as HTMLElement[]

        // Direct method
        if (!is_empty(elements)) {
            return resolve(elements)
        }

        // Observation method
        const observer = new MutationObserver(async (mutations) => {
            const elements = await query(xpath, parent)

            if (!is_empty(elements)) {
                clearInterval(interval)
                clearTimeout(timeout)
                observer.disconnect()
                return resolve(elements)
            }
        })
        observer.observe(document, {
            attributes: true,
            childList: true,
            subtree: true,
        })

        // Interval method
        const interval = setInterval(async () => {
            const elements = query(xpath, parent)

            if (!is_empty(elements)) {
                clearTimeout(timeout)
                clearInterval(interval)
                observer.disconnect()
                return resolve(elements)
            }
        }, 100)

        // Timeout method
        const timeout = setTimeout(() => {
            clearTimeout(timeout)
            clearInterval(interval)
            observer.disconnect()

            const elements = query(xpath, parent)
            if (!is_empty(elements)) {
                return resolve(elements[0])
            }
            return resolve(null)
        }, 5_000)
    })
}

async function query_one(xpath, parent = null): Promise<HTMLElement> {
    function query(xpath, parent) {
        let results = []
        //@note: 2
        const query = document.evaluate(
            xpath,
            parent || document,
            null,
            XPathResult.ORDERED_NODE_SNAPSHOT_TYPE,
            null,
        )
        // let query = document.evaluate(xpath, parent,
        //     null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null)
        for (let i = 0, length = query.snapshotLength;i < length;++i) {
            results.push(query.snapshotItem(i))
            break
        }
        return results
    }

    return await new Promise(resolve => {
        const elements = query(xpath, parent) as HTMLElement[]

        // Direct method
        if (!is_empty(elements)) {
            return resolve(elements[0])
        }

        // Interval method
        const interval = setInterval(async () => {
            const elements = query(xpath, parent)

            if (!is_empty(elements)) {
                clearTimeout(timeout)
                clearInterval(interval)
                observer.disconnect()

                return resolve(elements[0])
            }
        }, 100)

        // Timeout method
        const timeout = setTimeout(() => {
            clearTimeout(timeout)
            clearInterval(interval)
            observer.disconnect()

            const elements = query(xpath, parent)
            if (!is_empty(elements)) {
                return resolve(elements[0])
            }
        }, 5_000)

        // Observation method
        const observer = new MutationObserver(async (mutations) => {
            const elements = await query(xpath, parent)

            if (!is_empty(elements)) {
                clearInterval(interval)
                clearTimeout(timeout)
                observer.disconnect()

                return resolve(elements[0])
            }
        })
        observer.observe(document, {
            attributes: true,
            childList: true,
            subtree: true,
        })
    })
}


// ------------------------------------------------------------
// : Service : Extractor
// ------------------------------------------------------------
class Extractor {
    private config           // The injection loaded from the server
    private config_extractor // Contains the allowed domains 

    public async init() {
        const tab_hash = new URLSearchParams(window.location.search).get('dse')
        if (is_empty(tab_hash)) {
            return
        }

        // Wait for page to stop loading
        await this.wait_until_ready() 

        // Get the extractor configuration
        try {
            this.config           = await storage.get('config')
            this.config_extractor = this.config['extractor_config']
        } catch {
            log(`Extractor: Failed to load configuration`)
            return
        }

        // Check if the current url is in the configuration (as key)
        // Extract the data from search results
        let extracted_data = null
        for (const value in this.config_extractor) {
            const value_url     = new URL(`https://${value}`) // Act as if it's a url
            const value_domain  = value_url.hostname
            const value_search  = value_url.search
            const value_queries = Object.keys(Object.fromEntries(new URLSearchParams(value_search)))

            const window_url     = new URL(window.location.href)
            const window_domain  = window_url.hostname
            const window_search  = window_url.search
            const window_queries = Object.keys(Object.fromEntries(new URLSearchParams(window_search)))

            // Skip if window domain does not match with the value (config) domain
            if (!window_domain.includes(value_domain)) continue

            // Check if both value and window have intersected queries
            const intersected_queries  = intersection([value_queries, window_queries])
            const has_matching_queries = intersected_queries.length == value_queries.length

            if (has_matching_queries) {
                log(`Extractor: Found matching domain: ${value_domain}`)
                extracted_data = await this.extract(value)
                break
            }
        }

        if (is_empty(extracted_data)) {
            return
        }

        const keyword = new URLSearchParams(window.location.search).get('dse_keyword')
        const website = new URLSearchParams(window.location.search).get('dse_website')


        const html    = document.querySelector('html').innerHTML
        const payload = {
            'version'     : await storage.get('version'),
            'user.token'    : await storage.get('user.token'),
            'browser'     : await storage.get('browser'),
            'localization': window.navigator.language,
            'url'         : window.location.href,
            'keyword'     : keyword,
            'website'     : website,

            'html'    : html,
            'results' : extracted_data,
        }

        try {
            await browser.runtime.sendMessage({
                from   : 'content.extractor',
                data   : payload,
            })
        } catch (e) {
            log(`Extractor: Failed to send data to background script (${e})`)
        }

        // Clsoe current tab
        const tab = await browser.tabs.getCurrent()
        await browser.tabs.remove(tab.id)
    }

    private async extract(url: string): Promise<object> {
        try {
            // Define the result object
            const result = {}

            // Get the selectors
            const selectors = this.config_extractor[url]['selectors']

            // Go through each selector and extract the data
            for (const selector of selectors) {
                const key      = selector['name']
                const elements = query(selector['xpath'], document)

                // Add key to result object
                result[key] = []

                // Go through each xpath and extract the data and add it to the metadata
                for (const element of elements) {
                    const data = {}
                    for (const xpath in selector['xpaths']) {
                        const value  = query_text(selector['xpaths'][xpath], element)

                        if (value) {
                            data[xpath] = value
                        } else {
                            console.warn(`Key  : ${xpath}\nXPath: ${selector['xpaths'][xpath]}`)
                            
                            data[xpath] = null
                        }
                    }
                    result[key].push(data)
                }
            }
            return result
        } catch (error) {
            log(`Extractor: Failed to extract data (reason: ${error})`)
        }
    }

    /**
     * Wait until the page is completely loaded.
     * @returns {Promise<void>}
     */
    private async wait_until_ready() {
        // Define the size of the lengths array (bigger is slower)
        const count    = 2
        const interval = 333

        // Create array with {count} amount of elements
        const lengths = Array(count).fill(0)

        while (true) {
            const html   = document.querySelector('html') // Get the html element
            const length = html.innerText.length          // Get the length of html

            if (length <= 1000) {
                await sleep(interval)
                continue
            }

            lengths.push(length)                          // Add the length to the array

            if (lengths.length > count) {
                lengths.shift()                           // Remove the first element of the array
            }

            log(`Extractor: ${JSON.stringify(lengths)}`)                 

            // Check if the length of the html is changing
            if (lengths.every(length => length === lengths[0])) {
                break
            }
            
            await sleep(interval)
        }
    }
}

// ------------------------------------------------------------
// : Service : MouseCapture
// ------------------------------------------------------------
class ClickCapture {
    private config              // The injection loaded from the server
    private config_content      // The content of the configuration
    private config_mousecapture // The mouse capture configuration
    private config_extractor    // The extractor configurations (in order to fetch metadata from current page)

    private keywords            // The keywords to check relatability

    private allowed_domains     // Contains the domain this class is allowed to capture from
    private query_parameters    // Contains the query parameters to be captured (e.g. google has "q", youtube has "search_query")
    
    private selectors = 'a' // The selectors to capture from

    public async init() {
        log(`MouseCapture: Initializing`)

        // Check if user is registered
        let is_registered = await storage.get('user.registered')
        if (!is_registered) {
            log(`MouseCapture: User is not registered`)
            return
        }
        
        // Wait until the configuration is loaded from the server
        log(`MouseCapture: Getting content config`)
        wait_until(async() => {
            return !is_empty(await storage.get('config'))
        })


        // Get the content config
        this.config = await storage.get('config')

        // Throw error if config is empty
        if (is_empty(this.config)) {
            log(`MouseCapture: No configuration found`)
            return
        }

        this.config_content      = this.config['content_config']
        this.config_mousecapture = this.config['mousecapture_config']
        this.config_extractor    = this.config['extractor_config']
        
        this.keywords = this.config['keywords']

        this.allowed_domains     = this.config_content['allowed_domains']
        this.query_parameters    = this.config_mousecapture['query_parameters']
        
        // Exit if the current domain is not in the allowed domains (using regex)
        let allowed = false
        for (const allowed_domain of this.allowed_domains) {
            const regex_allowed_domain = new RegExp(allowed_domain)
            if (window.location.href.match(regex_allowed_domain)) {
                allowed = true
                break
            }
        }
        if (!allowed) {
            log(`MouseCapture: Current domain is not in the allowed domains`)
            log(`MouseCapture: Exiting`)
            return
        }

        // Get the query of the current url
        let query_name = null
        for (const query_parameter of this.query_parameters) {
            if (window.location.href.includes(query_parameter.name)) {
                query_name = query_parameter.value
                break
            }
        }
        const query_value = new URLSearchParams(window.location.search).get(query_name)
        if (!query_value) {
            log(`MouseCapture: No query parameter found`)
            log(`MouseCapture: Exiting`)
            return
        }

        // Wait for page to load
        log(`MouseCapture: Waiting for page load`)
        await this.wait_until_ready()
        log(`MouseCapture: Ready`)

        // Extract the data from search results
        let extracted_data = null
        for (const value in this.config_extractor) {
            const value_url     = new URL(`https://${value}`) // Act as if it's a url
            const value_domain  = value_url.hostname
            const value_search  = value_url.search
            const value_queries = Object.keys(Object.fromEntries(new URLSearchParams(value_search)))

            const window_url     = new URL(window.location.href)
            const window_domain  = window_url.hostname
            const window_search  = window_url.search
            const window_queries = Object.keys(Object.fromEntries(new URLSearchParams(window_search)))

            // Check if value domain matches window domain
            if (!window_domain.includes(value_domain)) continue

            // Check if both value and window have intersected queries
            const intersected_queries = intersection([value_queries, window_queries])
            const has_matching_queries = intersected_queries.length == value_queries.length

            if (has_matching_queries) {
                extracted_data = await this.extract_data(value)
            }
        }

        if (!is_empty(extracted_data)) {
            let relatable = false
            const data = JSON.stringify(extracted_data)

            // Check if keyword is in the extracted data
            for (const keyword of this.keywords) {
                if (data.includes(keyword)) {
                    relatable = true
                    break
                }
            }

            if (!relatable) {
                log(`MouseCapture: Extracted data is not relatable to the keywords`)
                return
            }
        }

        document.querySelectorAll(this.selectors).forEach((element: HTMLElement) => {
            element.addEventListener('click', async (e) => {
                const payload = {
                    'version'     : await storage.get('version'),
                    'user.token'  : await storage.get('user.token'),
                    'browser'     : await storage.get('browser'),
                    'localization': window.navigator.language,

                    'url'         : window.location.href,
                    'query'       : query_value,
                    'results'     : extracted_data,
                    'element'     : element.innerHTML,
                    'element_text': element.textContent.trim(),
                    'element_href': element.getAttribute('href'),
                }

                try {
                    log(`MouseCapture: Sending payload`)
                    await browser.runtime.sendMessage({
                        from: 'content.click',
                        data: payload,
                    })
                    log(`MouseCapture: Sent data to background script`)
                } catch (e) {
                    log(`MouseCapture: Failed to send data to background script (${e})`)
                }
            })
        })
    }

    private async extract_data(url: string) {
        try {
            const result = {}

            // Get the selectors
            const selectors = this.config_extractor[url]['selectors']

            // Go through each selector and extract the data
            for (const selector of selectors) {
                const key      = selector['name']
                const elements = query(selector['xpath'], document)

                // Add key to result object
                result[key] = []

                // Go through each xpath and extract the data and add it to the metadata
                for (const element of elements) {
                    const data = {}
                    for (const xpath in selector['xpaths']) {
                        const value = query_text(selector['xpaths'][xpath], element)
                        if (value) {
                            data[xpath] = value
                        } else {
                            data[xpath] = null
                            console.warn(`Key  : ${xpath}\nXPath: ${selector['xpaths'][xpath]}`)
                        }
                    }
                    result[key].push(data)
                }
            }
            return result
        } catch (error) {
            log(`MouseCapture: Failed to extract data (reason: ${error})`)
            throw error
        }
    }

    /**
     * Waits until the selector is found in the page.
     */
    private async wait_until_ready() {
        // Define the size of the lengths array (bigger is slower, but more accurate)
        const count    = 3
        const interval = 100

        // Create array with {count} amount of elements
        const lengths = Array(count).fill(0)

        while (true) {
            //@todo: Implement a more accurate and faster way using Observers
            const elements = document.querySelectorAll(this.selectors) // Get the html element
            const length   = elements.length                           // Get the length of html

            lengths.push(length) // Add the length to the array

            if (lengths.length > count) {
                lengths.shift()  // Remove the first element of the array
            }

            log(`MouseCapture: ${JSON.stringify(lengths)}`)

            // Check if the length of the html is changing
            if (lengths.every(length => length === lengths[0])) {
                break
            }
            
            await sleep(interval)
        }
    }
}

// ------------------------------------------------------------
// : Global
// ------------------------------------------------------------
let storage: Storage

const extractor     = new Extractor()
const mouse_capture = new ClickCapture()

// ------------------------------------------------------------
// : Start
// ------------------------------------------------------------
storage = new Storage()

setTimeout(() => {
    try {
        extractor.init()
    } catch (error) {
        log(`MouseCapture: Failed to initialize (reason: ${error})`)
    }

    try {
        mouse_capture.init()
    } catch (error) {
        log(`MouseCapture: Failed to initialize (reason: ${error})`)
    }
}, 10)