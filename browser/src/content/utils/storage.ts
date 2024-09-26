// ------------------------------------------------------------
// : Import
// ------------------------------------------------------------
import browser from 'webextension-polyfill'

// ------------------------------------------------------------
// : Utils
// ------------------------------------------------------------
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
