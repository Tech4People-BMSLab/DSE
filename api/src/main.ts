// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------
// General

// Utils
import _ from 'lodash'

// DateTime
import { DateTime } from 'luxon'

// WebSockets
import WebSocket, {WebSocketServer} from 'ws'
import msgpack                      from 'msgpack-lite'

// Events
import EventEmitter from 'eventemitter3'

// Database
import mongoose from 'mongoose'

// Mailing
import nodemailer from 'nodemailer'

// Scheduling
import Agenda from 'agenda'

// Logging
import pino        from 'pino'
import pretty      from 'pino-pretty'
import PrettyError from 'pretty-error'

// Env
import dotenv from 'dotenv'

// ------------------------------------------------------------
// : Config
// ------------------------------------------------------------
dotenv.config()
// ------------------------------------------------------------
// : Logging
// ------------------------------------------------------------
const pe = new PrettyError()
pe.skip(function(traceline) {return traceline.file != 'main.ts'})

const logger = pino({
    transport: {
        target: 'pino-pretty',
        options: {
            colorize: true,
        },
    },
    timestamp : () => `,"time":"${DateTime.now().toFormat('yyyyMMdd HH:mm:ss.SSS')}"`,
    level     : 'debug',
    levelFirst: true,
    ignore    : /(node_modules)/,
})

// ------------------------------------------------------------
// : Utils
// ------------------------------------------------------------
const is_empty     = (value: any) => _.isEmpty(value)
const is_null      = (value: any) => _.isNull(value)
const is_function  = (value: any) => _.isFunction(value)

const find   = (value: any, predicate: Function) => _.find(value, predicate)
const filter = (value: any, predicate: Function) => _.filter(value, predicate)
const sort   = (value: any, predicate: Function) => _.sortBy(value, predicate)

const first = (value: any) => _.first(value)
const last  = (value: any) => _.last(value)

const sleep = (ms: number) => {
    return new Promise(resolve => setTimeout(resolve, ms))
}

const wait_until = async (predicate: Function, timeout=-1, interval=100) => {
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

const generate_uid = async () => {
    const chars = '0123456789abcdef'
    let id      = ''

    while (true) {
        // Generate a random 12-character ID
        for (let i = 0; i < 12; i++) {
            const random_index = Math.floor(Math.random() * chars.length)
            id += chars.charAt(random_index)
        }

        // Check if the ID is valid and unique
        const existingUser = await User.findOne({ uid: id })
        if (!existingUser) {
            // Return the unique ID if it passes the validation check
            return id
        } else {
            // If the ID is not unique, reset the ID string and generate a new ID
            id = ''
        }
    }
}

// ------------------------------------------------------------
// : Schemas
// ------------------------------------------------------------
const user_schema = new mongoose.Schema({
    token     : {type: String, required: true},
    timestamp : {type: Date, default: DateTime.now().toJSDate()},

    resident  : {type: String, required: true},
    sex       : {type: String, required: true},
    age       : {type: String, required: true},
    postcode  : {
        value : {type: String, required: true},
    },
    education : {type: String, required: true},
    language  : [{
        nederlands: {type: Boolean, required: true},
        engels    : {type: Boolean, required: true},
        duits     : {type: Boolean, required: true},
        frans     : {type: Boolean, required: true},
        spaans    : {type: Boolean, required: true},
        italiaans : {type: Boolean, required: true},
    }],
    political : {type: String, required: true},
    employment: {type: String, required: true},
    income    : {type: String, required: true},
})

/**
 * Similar to users, but this is actual user who are connected to the websocket server
 */
const client_schema = new mongoose.Schema({
    token: {type: String, required: true},

    status: {type: String, required: true, default: 'offline', enum: ['online', 'offline']},
    events: [{
        timestamp: {type: Date, default: DateTime.now().toJSDate(), required: false},
        type     : {type: String, required: false},
        data     : {type: Object, required: false},
    }]
})


const search_schema = new mongoose.Schema({
    timestamp: {type: Date, default: DateTime.now().toJSDate(), required: true},

    user: {type: mongoose.Schema.Types.ObjectId, ref: 'User', required: true},

    version      : {type: String, required: false},
    browser      : {type: Object, required: false},
    localization : {type: String, required: false},
    url          : {type: String, required: false},
    keyword      : {type: String, required: false},
    website      : {type: String, required: false},
    results      : {type: Object, required: false},
})

const click_schema = new mongoose.Schema({
    timestamp: {type: Date, default: DateTime.now().toJSDate(), required: true},

    user: {type: mongoose.Schema.Types.ObjectId, ref: 'User', required: true},

    version      : {type: String, required: false},
    browser      : {type: Object, required: false},
    localizations: {type: String, required: false},
    url          : {type: String, required: false},
    query        : {type: String, required: false},
    results      : {type: Object, required: false},
    element      : {type: String, required: false},
    element_text : {type: String, required: false},
    element_href : {type: String, required: false},
})

// ------------------------------------------------------------
// : Models
// ------------------------------------------------------------
const User    = mongoose.model('User'   , user_schema)
const Client  = mongoose.model('Client' , client_schema)
const Search  = mongoose.model('Search' , search_schema)
const Click   = mongoose.model('Click'  , click_schema)

// ------------------------------------------------------------
// : Database
// ------------------------------------------------------------
class Database {
    public db  : any | null = null
    public status: 'connected' | 'disconnected' = 'disconnected'

    private host: string
    private port: string
    private name: string

    private user: string
    private pass: string

    private url : string

    constructor() {
        if (is_empty(DB_HOST)) throw new Error('DB_HOST is empty')
        if (is_empty(DB_PORT)) throw new Error('DB_PORT is empty')
        if (is_empty(DB_NAME)) throw new Error('DB_NAME is empty')
        if (is_empty(DB_USER)) throw new Error('DB_USER is empty')
        if (is_empty(DB_PASS)) throw new Error('DB_PASS is empty')

        this.host = DB_HOST
        this.port = DB_PORT
        this.name = DB_NAME
        
        this.user = DB_USER
        this.pass = DB_PASS

        this.url = `mongodb://${this.user}:${this.pass}@${this.host}:${this.port}/dse?authSource=admin`
    }

    async init() {
        try {
            this.db  = await mongoose.connect(this.url)

            logger.debug(`DB: Connected with ${this.host}:${this.port}/${this.name}`)

            // Set status
            this.status = 'connected'
        } catch (e) {
            logger.error(pe.render(e as any))
        }
    }
}

// ------------------------------------------------------------
// : Scheduler
// ------------------------------------------------------------
class UserOfflineError extends Error { constructor(token: string) { super(`User ${token} is offline`) }}

class Scheduler {
    private url   : string 

    public agenda: Agenda | null = null

    private host: string
    private port: string
    private user: string
    private pass: string

    public status: string = 'inactive'

    constructor() {
        this.host = DB_HOST
        this.port = DB_PORT
        
        this.user = DB_USER
        this.pass = DB_PASS

        this.url = `mongodb://${this.user}:${this.pass}@${this.host}:${this.port}/agenda?authSource=admin`
    }

    public async init() {
        try {
            this.agenda = new Agenda({db: {address: this.url, collection: 'agenda'}})

            // Event listeners
            this.agenda.on('fail', (e: any, job: any) => {
                logger.error(`Scheduler failed job ${job.attrs.name}`)
                logger.error(pe.render(e))
            })
            this.agenda.on('start', (job: any) => {
                logger.info(`Scheduler started job ${job.attrs.name}`)
            })
            this.agenda.on('complete', (job: any) => {
                logger.info(`Scheduler completed job ${job.attrs.name}`)
            })

            // Define jobs
            this.agenda.define('search_process', async (job: any, done: any) => {
                try {
                    const {token} = job.attrs.data

                    if (!await ws.is_online(token)) {
                        throw new UserOfflineError(token)
                    }

                    await ws.send(token, {action: 'search_process.start'})
                    done()
                } catch (e) {
                    if (e instanceof UserOfflineError || e instanceof SendError) {
                        job.attrs.lockedAt   = null // Unlock job
                        job.attrs.failReason = e.message
                        job.attrs.failedAt   = DateTime.now()
                        job.attrs.nextRunAt  = DateTime.now().plus({minute: 1})
                        job.save()
                        done(e)
                    } else {
                        logger.error(pe.render(e as any))
                        done(e)
                    }
                }
            })

            this.agenda.define('cleanup', async (job: any, done: any) => {
                try {
                    const jobs = await this.agenda?.jobs({name: 'search_process'}) as any

                    // Remove duplicate search_jobs for each client
                    const jobs_by_token = new Map()
                    jobs.forEach((job: any) => {
                        const token = job.attrs.data.token
                        const job_next_run = DateTime.fromJSDate(job.attrs.nextRunAt)
                        const job_existing = jobs_by_token.get(token)
                        if (job_existing) {
                            const job_existing_next_run = DateTime.fromJSDate(job_existing.attrs.nextRunAt)
                            if (job_next_run > job_existing_next_run) {
                                jobs_by_token.set(token, job)
                            } else {
                                job.remove()
                            }
                        } else {
                            jobs_by_token.set(token, job)
                        }
                    })              

                    // Reschedule search_jobs with invalid nextRunAt
                    for (const job of jobs_by_token.values()) {
                        const job_next_run = DateTime.fromJSDate(job.attrs.nextRunAt)
                        if (job_next_run.hour < 11 || job_next_run.hour > 23) {
                            job.attrs.nextRunAt = job_next_run.plus({ days: 1 }).set({ hour: 11, minute: 0, second: 0, millisecond: 0 }).toJSDate()
                            await job.save()
                        }
                    }

                    // Reschedule if client already done a search this week
                    const clients = await Client.find()

                    for (const client of clients) {
                         const events = client.events.filter((event: any) => event.name === 'search_process.start')
                         
                        if (events.length > 0) {
                            const last_event      = events[events.length - 1] as any
                            const last_event_time = DateTime.fromJSDate(last_event.time)

                            const last_monday = DateTime.local().startOf('week').minus({ days: 7 })
                            const next_monday = DateTime.local().startOf('week').plus({ weeks: 1 })

                            if (last_event_time >= last_monday && last_event_time < next_monday) {
                                const job = await this.agenda?.jobs({name: 'search_process', 'data.token': client.token}) as any
                                if (job) {
                                    job.attrs.lockedAt  = null
                                    job.attrs.nextRunAt = next_monday.set({ hour: 11, minute: 0, second: 0, millisecond: 0 }).toJSDate()
                                    await job.save()
                                }
                            }
                        }
                    }
    
                    done()
                } catch (e) {
                    logger.error(pe.render(e as any))
                }
            })

            // Event Reschedulers
            emitter.on('client_connected', async (token: string) => {
                try {
                    if (this.agenda === null) throw new Error('Agenda is null')

                    // Get jobs for this user
                    let jobs = await this.agenda.jobs({name: 'search_process', 'data.token': token})
    
                    // If no jobs are found, create one
                    if (is_empty(jobs)) {
                        await this.agenda?.every('0 11 * * 1', 'search_process', {token: token})
                    } else {
                        const job  = jobs[0]

                        const time_now = DateTime.now()
                        const time_job = DateTime.fromJSDate(job.attrs.nextRunAt as any)

                        const is_in_past          = time_job < time_now
                        const is_within_timeframe = time_job.hour >= 11 && time_job.hour <= 23 && time_job.day === time_now.day

                        if (is_in_past) {
                            if (is_within_timeframe) {
                                // Execute the job
                                await job.run()
                            } else {
                                // Reschedule for tomorrow 11am
                                job.attrs.lockedAt  = null
                                job.attrs.nextRunAt = time_now.plus({ days: 1 }).set({ hour: 11, minute: 0, second: 0, millisecond: 0 }).toJSDate()
                                await job.save()
                            }
                        } else if (!is_within_timeframe) {
                                // Reschedule for tomorrow 11am
                                job.attrs.lockedAt = null
                                job.attrs.nextRunAt = time_now.plus({ days: 1 }).set({ hour: 11, minute: 0, second: 0, millisecond: 0 }).toJSDate()
                                await job.save()
                        }
                    }
                } catch (e) {
                    logger.error(pe.render(e as any))
                }
            })

            // Start scheduler
            await this.agenda.maxConcurrency(10_000)
            await this.agenda.defaultLockLifetime(1000)
            await this.agenda.processEvery('1 second')

            await this.agenda.start()
            await this.agenda.every('0 * * * *', 'cleanup')

            this.status = 'active'

            logger.debug('Scheduler: Started')
        } catch (e) {
            logger.error(pe.render(e as any))
        }
    }

    public async schedule(name: string, data: object, when: string) {
        try {
            await this.agenda?.schedule(when, name, data)
        } catch (e) {
            logger.error(pe.render(e as any))
        }
    }

    public async has(name: string) {
        try {
            return await this.agenda?.jobs({name: name}) !== undefined
        } catch (e) {
            logger.error(pe.render(e as any))
        }
    }
}

// ------------------------------------------------------------
// : WebSocket Server
// ------------------------------------------------------------
class SendError extends Error { constructor(token: string) { super(`Failed to send message to ${token}`) }}

class WSServer {
    private server: WebSocketServer | null = null
    private host  : string
    private port  : string

    private clients: object[] = []

    constructor() {
        this.host = WS_HOST
        this.port = WS_PORT
    }

    public async init() {
        try {
            this.server = new WebSocketServer({
                host: this.host,          // 10.0.0.10
                port: parseInt(this.port) // 6000
            })

            logger.debug(`WS: Listening on ${this.host}:${this.port}`)

            this.server.on('connection', (client: WebSocket) => { this.on_connect(client) })
            this.server.on('error', (e: any) => { logger.error(pe.render(e as any)) })
        } catch (e) {
            logger.error(pe.render(e as any))
        }
    }

    //// Event Handlers
    private async on_connect(client: WebSocket) {
        // Wait until database is actually connected
        await wait_until(() => db.status === 'connected')
        await wait_until(() => scheduler.status === 'active')

        const ip   = (client as any)._socket.remoteAddress
        const port = (client as any)._socket.remotePort
        logger.info(`WS: [${ip}|xxxxxxxxxxxx]: connected`)

        client.on('message', (buffer: any) => { this.on_receive(client, buffer)   })
        client.on('close'  , (event : any) => { this.on_disconnect(client, event) })

        client.on('error'  , (e: any) => { logger.error(pe.render(e as any)) })
    }

    private async on_receive(client: WebSocket, buffer: any) {
        try {
            // Get metadata
            const ip   = (client as any)._socket.remoteAddress
            const port = (client as any)._socket.remotePort

            // Get the payload
            const payload = msgpack.decode(buffer)

            const action    = payload.action
            const data      = payload.data   as any

            let token = undefined
            try       { token = payload.token  ?? payload.data.token }
            catch (e) { /* Ignore */ }

            switch (action) { // Handles request without token
                //// Register
                case 'register': {
                    const token = await generate_uid()

                    await new User({
                        token: token,

                        resident  : data.resident,
                        sex       : data.sex,
                        age       : data.age,
                        postcode  : data.postcode,
                        education : data.education,
                        language  : data.language,
                        political : data.political,
                        employment: data.employment,
                        income    : data.income,
                    }).save()

                    client.send(msgpack.encode({
                        action: 'register.success',
                        data  : { token: token }
                    }))

                    logger.info(`WS: [${ip}|${token}]: registered`)
                    return
                }
            }


            if (token) {
                const user = await User.findOne({ token: token })
                if (!user) {
                    client.send(msgpack.encode({action: 'token.invalid'}))
                    return
                }
            } else {
                client.send(msgpack.encode({action: 'token.invalid'}))
                return
            }

            switch (action) { // Handles request with token
                //// Connect
                case 'connect': {
                    if (await this.verify_token(token) === false) 

                    // Add client to local list
                    this.clients.push({token: token, client: client})

                    // Notify client that token is valid
                    await client.send(msgpack.encode({action: 'token.valid'}))

                    // Check if client already exists in database
                    const doc = await Client.findOneAndUpdate(
                        {token: token}, 
                        {status: 'online'}, 
                        {upsert: true, new: true}
                    )

                    if (doc) {
                        doc.events.push({type: 'connected'})
                        await doc.save()
                    }

                    emitter.emit('client_connected', token)

                    logger.info(`WS: [${ip}|${token}]: connected`)
                    break
                }

                //// Search Process Result
                case 'content.extractor': {
                    if (is_empty(token)) return

                    const user   = await User.findOne({ token: token })
                    const search = await new Search({
                        token: token,
                        user : user,

                        version     : data.version,
                        browser     : data.browser,
                        localization: data.localization,
                        url         : data.url,
                        keyword     : data.keyword,
                        website     : data.website,
                        results     : data.results,
                    }).save()

                    logger.info(`WS: [${ip}|${token}]: Added extractor data`)
                    break
                }

                //// Click Result
                case 'content.click': {
                    if (is_empty(token)) return

                    const user  = await User.findOne({ token: token })
                    const click = await new Click({
                        token: token,
                        user : user,

                        version      : data.version,
                        browser      : data.browser,
                        localizations: data.localizations,
                        url          : data.url,
                        query        : data.query,
                        results      : data.results,
                        element      : data.element,
                        element_text : data.element_text,
                        element_href : data.element_href,
                    }).save()

                    logger.info(`WS: [${ip}|${token}]: Added click data`)
                    break
                }

                //// Events
                case 'event': {
                    if (is_empty(token)) return

                    const client = await Client.findOne({ token: token })
                    if (!client) return

                    client.events.push({type: data.type})
                    await client.save()

                    logger.info(`WS: [${ip}|${token}]: Added event ${data.type}`)
                    break
                }

                //// Email
                case 'email': {
                    
                    break
                }

                default:
                    throw new Error(`Unknown action: ${action}`)
            }
        } catch (e) {
            logger.error(e)
        }
    }

    private async on_disconnect(client: WebSocket, event: any) {
        try {
            let token = find(this.clients, (c: any) => c.client === client)?.token
            if (is_empty(token)) return

            const ip = (client as any)._socket.remoteAddress
            logger.info(`WS: [${ip}|${token}]: disconnected (code: ${event})`)

            // Remove client from local list
            this.clients = this.clients.filter((c: any) => c.token !== token)

            // Update client in database
            const doc = await Client.findOne({token: token})

            if (doc) {
                doc.status = 'offline'
                doc.events.push({type: 'disconnected'})
                await doc?.save()
            }
        } catch (e) {
            logger.error(pe.render(e as any))
        }
    }

    //// Public Methods
    public async send(token: string, data: object) {
        if (is_empty(token)) throw new Error('Token is empty')
        if (is_empty(data))  throw new Error('Data is empty')
            
        try {
            const client = find(this.clients, (c: any) => c.token === token)?.client
            if (!client) throw new Error(`Client (${token}) not found`)

            const payload = msgpack.encode(data)
            client.send(payload)
        } catch (e) {
            logger.error(pe.render(e as any))
            throw new SendError(token)
        }
    }

    public async is_online(token: string) {
        if (is_empty(token)) throw new Error('Token is empty')

        try {
            const client = find(this.clients, (c: any) => c.token === token)?.client
            if (!client) throw new Error(`Client (${token}) not found`)

            return true
        } catch (e) {
            logger.error(pe.render(e as any))
            return false
        }
    }

    //// Private Methods
    private async verify_token(token: string) {
        const client = await User.findOne({ token })
        if (!client) { return false }
        return true
    }
}

// ------------------------------------------------------------
// : Email
// ------------------------------------------------------------
class Email {
    private transporter: any

    constructor () {
        if (!SMTP_HOST) throw new Error(`SMTP_HOST is empty (${SMTP_HOST})`)
        if (!SMTP_PORT) throw new Error(`SMTP_PORT is empty (${SMTP_PORT})`)
        if (!SMTP_USER) throw new Error(`SMTP_USER is empty (${SMTP_USER})`)
        if (!SMTP_PASS) throw new Error(`SMTP_PASS is empty (${SMTP_PASS})`)

        this.transporter = nodemailer.createTransport({
            host: SMTP_HOST,
            port: parseInt(SMTP_PORT),
            secure: false,
            auth: {
                user: SMTP_USER,
                pass: SMTP_PASS
            }
        })
    }

    public async send_email(to: string, subject: string, text: string, html: string) {
        const opts = {
            from   : SMTP_USER,
            to     : to,
            subject: subject,
            text   : text,
            html   : html
        }

        try       { await this.transporter.sendMail(opts) } 
        catch (e) { logger.error(pe.render(e as any)) }
    }
}

// ------------------------------------------------------------
// : Global
// ------------------------------------------------------------
const DB_HOST = process.env.DB_HOST   ?? 'localhost'
const DB_PORT = process.env.DB_PORT   ?? '27017'
const DB_NAME = process.env.DB_NAME   ?? 'dse'
const DB_USER = process.env.DB_USER   ?? 'root'
const DB_PASS = process.env.DB_PASS   ?? 'password'

const SMTP_HOST = process.env.SMTP_HOST 
const SMTP_PORT = process.env.SMTP_PORT ?? '25'
const SMTP_USER = process.env.SMTP_USER
const SMTP_PASS = process.env.SMTP_PASS

const WS_HOST = process.env.WS_HOST   ?? '0.0.0.0'
const WS_PORT = process.env.WS_PORT   ?? '5000'

const emitter = new EventEmitter()

let db        : Database
let ws        : WSServer
let scheduler : Scheduler
let email     : Email

// ------------------------------------------------------------
// : Main
// ------------------------------------------------------------
async function main() {
    try {
        db        = new Database()
        ws        = new WSServer()
        scheduler = new Scheduler()

        try       { email = new Email() }
        catch (e) { logger.error(e)     }
        

        Promise.allSettled([
            db.init(), // Database
            ws.init(), // WebSocket
            scheduler.init(), // Scheduler
        ])
    } catch (error) {
        logger.error(pe.render(error as any))
    }
}

main()