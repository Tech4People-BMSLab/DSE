// ------------------------------------------------------------
// : Imports
// ------------------------------------------------------------
import fs   from 'fs'

import gulp       from 'gulp'
import terser     from 'gulp-terser'
import rename     from 'gulp-rename'
import zip        from 'gulp-zip'
import babel      from 'gulp-babel'
import esbuild    from 'gulp-esbuild'
import gulpif     from 'gulp-if'
import sourcemaps from 'gulp-sourcemaps'
import merge      from 'merge-stream'

import svelte           from 'esbuild-svelte'
import sveltePreprocess from 'svelte-preprocess'

import { readFile, writeFile } from 'fs/promises'

import { connect } from 'nats'

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
let nc    = null
const env = {
    mode       : 'production',
    base_url   : 'http://localhost:5000',
    nats_url   : 'wss://dev.bmslab.utwente.nl:9222',
    nats_user  : 'dev', // prod, dev, local
    nats_pass  : 'dev', // prod, dev, local
}

// ------------------------------------------------------------
// : Subtasks
// ------------------------------------------------------------
gulp.task('clean', async function () {
    const directory = 'dist'
    fs.rmdirSync(directory, { recursive: true, force: true})
    fs.mkdirSync(directory, { recursive: true, force: true})

    fs.mkdirSync('dist/chrome',  { recursive: true, force: true})
    fs.mkdirSync('dist/firefox', { recursive: true, force: true})
})


gulp.task('versioning', async function () {
    const paths = [
        './public/firefox_manifest.json', 
        './public/chrome_manifest.json'
    ]

    for (const path of paths) {
        const manifest = JSON.parse(await readFile(path, 'utf-8'))

        const parts             = manifest.version.split('.')
        parts[parts.length - 1] = (parseInt(parts[parts.length - 1], 10) + 1).toString()
        manifest.version        = parts.join('.')

        await writeFile(path, JSON.stringify(manifest, null, 4))
    }
})

gulp.task('background', function () {
    return gulp.src('src/background/background.ts')
        .pipe(esbuild({
            bundle   : true,
            outfile  : 'background.js',
            platform : 'browser',
            target   : 'es2020',
            define   : {
                'import.meta.env.MODE'         : env.mode ? `"${env.mode}"` : 'null',
                'import.meta.env.BASE_URL'     : env.base_url ? `"${env.base_url}"` : 'null',
                'import.meta.env.NATS_URL'     : env.nats_url ? `"${env.nats_url}"` : 'null',
                'import.meta.env.NATS_USER'    : env.nats_user ? `"${env.nats_user}"` : 'null',
                'import.meta.env.NATS_PASS'    : env.nats_pass ? `"${env.nats_pass}"` : 'null',
                'import.meta.env.IPC_VERSION'  : '"v1.5"',
            }
        }))
        .pipe(gulpif(env.mode === 'production', terser()))
        .on('error', function (err) {
            console.error(err)
            this.emit('end')
        })
        .pipe(gulp.dest('dist/chrome'))
        .pipe(gulp.dest('dist/firefox'))
        .on('end', async function () {
            if (!nc) return
            nc.publish('bms.dse.background.reload')
        })
})

gulp.task('popup', function () {
    return gulp.src('src/popup/popup.ts')
        .pipe(esbuild({
            bundle  : true,
            outfile : 'popup.js',
            platform: 'browser',
            target  : 'es2020',
            define: {
                'import.meta.env.MODE'         : env.mode ? `"${env.mode}"` : 'null',
                'import.meta.env.BASE_URL'     : env.base_url ? `"${env.base_url}"` : 'null',
                'import.meta.env.NATS_URL'     : env.nats_url ? `"${env.nats_url}"` : 'null',
                'import.meta.env.NATS_USER'    : env.nats_user ? `"${env.nats_user}"` : 'null',
                'import.meta.env.NATS_PASS'    : env.nats_pass ? `"${env.nats_pass}"` : 'null',
                'import.meta.env.IPC_VERSION'  : '"v1.5"',
            },
            plugins : [svelte({
                preprocess: sveltePreprocess(),
            })],
            loader: {
                '.ts': 'ts',
            },
        }))
        .on('error', function (err) {
            console.error(err)
            this.emit('end')
        })
        .pipe(gulp.dest('dist/chrome'))
        .pipe(gulp.dest('dist/firefox'))
        .on('end', async function () {
            if (!nc) return
            nc.publish('bms.dse.popup.reload')
        })
})

gulp.task('content', function () {
    return gulp.src('src/content/content.ts')
        .pipe(esbuild({
            bundle  : true,
            minify   : false, 
            sourcemap: false,
            outfile : 'content.js',
            platform: 'browser',
            target  : 'es2020',
            define  : {
                'import.meta.env.MODE'         : env.mode ? `"${env.mode}"` : 'null',
                'import.meta.env.BASE_URL'     : env.base_url ? `"${env.base_url}"` : 'null',
                'import.meta.env.NATS_URL'     : env.nats_url ? `"${env.nats_url}"` : 'null',
                'import.meta.env.NATS_USER'    : env.nats_user ? `"${env.nats_user}"` : 'null',
                'import.meta.env.NATS_PASS'    : env.nats_pass ? `"${env.nats_pass}"` : 'null',
                'import.meta.env.IPC_VERSION'  : '"v1.5"',
            }
        }))
        // .pipe(babel({
        //     presets: ['@babel/preset-env']
        // }))
        .pipe(gulpif(env.mode === 'production', terser()))
        .on('error', function (err) {
            console.error(err)
            this.emit('end')
        })
        .pipe(gulp.dest('dist/chrome'))
        .pipe(gulp.dest('dist/firefox'))
        .on('end', async function () {
            if (!nc) return
            nc.publish('bms.dse.background.reload')
            nc.publish('bms.dse.content.reload')
        })
})

gulp.task('assets', function () {
    return gulp.src(['public/index.html', 'src/popup/popup.html', 'public/images/**/*'], { allowEmpty: true })
        .pipe(gulp.dest(function(file) {
            return file.base.includes('images') ? 'dist/chrome/images' : 'dist/chrome'
        }))
        .pipe(gulp.dest(function(file) {
            return file.base.includes('images') ? 'dist/firefox/images' : 'dist/firefox'
        }))
})

gulp.task('manifest', function () {
    return merge(
        gulp.src('public/chrome_manifest.json')
            .pipe(rename('manifest.json'))
            .pipe(gulp.dest('dist/chrome')),

        gulp.src('public/firefox_manifest.json')
            .pipe(rename('manifest.json'))
            .pipe(gulp.dest('dist/firefox'))
    )
})

gulp.task('zip', function () {
    return merge(
        gulp.src('dist/chrome/**/*') // Chrome
            .pipe(zip('extension.zip'))
            .pipe(gulp.dest('dist/chrome')),
        gulp.src('src/**/*')
            .pipe(zip('source.zip'))
            .pipe(gulp.dest('dist/chrome')),

        gulp.src('dist/firefox/**/*') // Firefox
            .pipe(zip('extension.zip'))
            .pipe(gulp.dest('dist/firefox')),
        gulp.src('src/**/*')
            .pipe(zip('source.zip'))
            .pipe(gulp.dest('dist/firefox'))
    )
})

// ------------------------------------------------------------
// : Task - Development
// ------------------------------------------------------------
gulp.task('dev-remote', 
    gulp.series('background', 'popup', 'content', 'assets', 'manifest', 
        function () {
            env.mode      = 'development'
            env.nats_url  = 'wss://dev.bmslab.utwente.nl:9222'
            env.nats_user = 'dev' // Local testing (default : local)
            env.nats_pass = 'dev' // Local testing (default : local)

            new Promise(async () => {
                nc = await connect({
                    servers: ['nats://dev.bmslab.utwente.nl:4222'],
                    user    : 'dev', // default : local
                    pass    : 'dev', // default : local
                    timeout             : 1000,
                    reconnect           : true,
                    reconnectTimeWait   : 1000,
                    maxReconnectAttempts: -1,
                })
            })

            gulp.watch('src/background/**/*.ts', gulp.series('background'))
            gulp.watch('src/content/**/*.ts'   , gulp.series('content'))
            gulp.watch([
                'src/popup/**/*.ts', 
                'src/popup/**/*.svelte'
            ], gulp.series('popup'))

            gulp.watch('public/**/*'      , gulp.series('assets'))
            gulp.watch('public/**/*.json' , gulp.series('manifest'))
    })
)

gulp.task('dev-demo', 
    gulp.series('background', 'popup', 'content', 'assets', 'manifest', 
        function () {
            env.mode      = 'development'
            env.base_url  = 'http://localhost:5000'
            env.nats_url  = 'wss://demo.nats.io:8443'
            env.nats_user = null
            env.nats_pass = null

            new Promise(async () => {
                nc = await connect({
                    servers: ['nats://demo.nats.io:4222'],
                    timeout             : 1000,
                    reconnect           : true,
                    reconnectTimeWait   : 1000,
                    maxReconnectAttempts: -1,
                })
            })

            gulp.watch('src/background/**/*.ts', gulp.series('background'))
            gulp.watch('src/content/**/*.ts'   , gulp.series('content'))
            gulp.watch([
                'src/popup/**/*.ts', 
                'src/popup/**/*.svelte'
            ], gulp.series('popup'))

            gulp.watch('public/**/*'      , gulp.series('assets'))
            gulp.watch('public/**/*.json' , gulp.series('manifest'))
    })
)


gulp.task('dev-local', 
    gulp.series('background', 'popup', 'content', 'assets', 'manifest', 
        function () {
            env.mode      = 'development'
            env.base_url  = 'http://localhost:5000'
            env.nats_url  = 'wss://10.0.0.50:9222'
            env.nats_user = 'local' // Local testing (default : local)
            env.nats_pass = 'local' // Local testing (default : local)

            new Promise(async () => {
                nc = await connect({
                    servers: ['nats://10.0.0.50:4222'],
                    user    : 'local', // default : local
                    pass    : 'local', // default : local
                    timeout             : 1000,
                    reconnect           : true,
                    reconnectTimeWait   : 1000,
                    maxReconnectAttempts: -1,
                })
            })

            gulp.watch('src/background/**/*.ts', gulp.series('background'))
            gulp.watch('src/content/**/*.ts'   , gulp.series('content'))
            gulp.watch([
                'src/popup/**/*.ts', 
                'src/popup/**/*.svelte'
            ], gulp.series('popup'))

            gulp.watch('public/**/*'      , gulp.series('assets'))
            gulp.watch('public/**/*.json' , gulp.series('manifest'))
    })
)

// ------------------------------------------------------------
// : Task - Build
// ------------------------------------------------------------
gulp.task('build', gulp.series('clean', 'versioning', 'background', 'popup', 'content', 'assets', 'manifest', 'zip'))
