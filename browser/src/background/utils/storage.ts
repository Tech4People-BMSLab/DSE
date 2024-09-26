// ------------------------------------------------------------
// : Imports    
// ------------------------------------------------------------

// Polyfill
import browser from 'webextension-polyfill'

// Lodash
import _ from 'lodash'

// ------------------------------------------------------------
// : Utils
// ------------------------------------------------------------
const has    = _.has
const get    = _.get
const pick   = _.pick
const values = _.values

// ------------------------------------------------------------
// : Storage
// ------------------------------------------------------------
/**
 * Wrapper class for the browser storage.
 * The wrapper class provides a bit more simpler interface for the storage to fetch data in a single function call.
 */
export class Storage {
    private storage  = browser.storage.local           // The storage object
    public onChanged = browser.storage.local.onChanged // The listener for the storage changes

    /**
     * Get a value from the storage
     * @param key The key to get the value from
     * @returns The value of the key
     */
    public async get(key: string) {
        try {
            const data = await this.get_all()
            return get(data, [key])
        } catch (e) {
            console.error(e)
        }
    }

    // Define functions to get all
    public async get_all() {
        try {
            return await this.storage.get()
        } catch (e) {
            console.error(e)
        }
    }
    
    /**
     * Set a value in the storage
     * @param obj The object to set the values from
     */
    public async set(obj: object) {
        try {
            await this.storage.set(obj)
        } catch (e) {
            console.error(e)
        }
    }
    
    /**
     * Listen to a value in the storage
     * @param event The event to listen to
     * @param callback The callback to call when the event is triggered
     */
    public on(event: string, callback: Function) {
        this.onChanged.addListener(async(changes: object) => {
            if (!has(changes, event))               return
            if (!has(changes, `${event}.newValue`)) return

            const value = changes[event].newValue
            callback(value)
        })
    }

    /**
     * Remove a value from the storage
     * @param key The key to remove from the storage
     */
    public async remove(key: string) {
        try {
            await this.storage.remove(key)
        } catch (e) {
            console.error(e)
        }
    }

    /**
     * clear the storage
     */
    public async clear() {
        try {
            await this.storage.clear()
        } catch (e) {
            console.error(e)
        }
    }
}

// ------------------------------------------------------------
// : Instance
// ------------------------------------------------------------
export const storage = new Storage()