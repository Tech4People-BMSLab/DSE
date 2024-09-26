// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------

// DateTime
import { DateTime } from 'luxon'

// ------------------------------------------------------------
// : Class
// ------------------------------------------------------------
export class Logger {

    private readonly name: string

    constructor(name: string) {
        this.name = name
    }

    private get_prefix() {
        return `${DateTime.local().toFormat('[yyyy-LL-dd HH:mm:ss]')} ${this.name}:`
    }

    public log(...args) {
        console.log(this.get_prefix(), ...args)
    }

    public error(...args) {
        console.error(this.get_prefix(), ...args)
    }
}
