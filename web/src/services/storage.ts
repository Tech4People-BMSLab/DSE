// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------
import localforage from 'localforage'

import { set } from '@/utils/utils'
// ------------------------------------------------------------
// : Cookie Manager
// ------------------------------------------------------------
export class StorageManager {
    private lf: LocalForage

    constructor() {
        localforage.config({
            version     : 1.0,
            driver      : localforage.INDEXEDDB,
            name        : 'opsreb',
            storeName   : 'opsreb',
            description : 'OPSREB'
        })

        this.lf = localforage.createInstance({
            name        : 'opsreb',
            storeName   : 'opsreb',
        })
    }

    public async get(key: string): Promise<any> {
        return await this.lf.getItem(key)
    }

    public async all(): Promise<any> {
        const keys = await this.lf.keys()
        const data = {}
        for (const key of keys) {
            set(data, key, await this.lf.getItem(key))
        }
        return data
    }

    public async set(key: string, value: any): Promise<void> {
        return await this.lf.setItem(key, value)
    }

    public async has(key: string): Promise<boolean> {
        return await this.lf.getItem(key) !== null
    }

    public async remove(key: string): Promise<void> {
        return await this.lf.removeItem(key)
    }

    public async clear(): Promise<void> {
        return await this.lf.clear()
    }
}

// ------------------------------------------------------------
// : Instance
// ------------------------------------------------------------
const storage_manager = new StorageManager()
// ------------------------------------------------------------
// : Exports
// ------------------------------------------------------------
export default storage_manager
