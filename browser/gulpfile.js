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

import webExt from 'web-ext'

import axios from 'axios'

import svelte           from 'esbuild-svelte'
import sveltePreprocess from 'svelte-preprocess'

import { readFile, writeFile } from 'fs/promises'
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
const env = {
    mode       : 'production',
    base_url   : 'https://static.33.56.161.5.clients.your-server.de/dse',
    base_ws    : 'wss://static.33.56.161.5.clients.your-server.de/dse/ws'
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

    // for (const path of paths) {
    //     const manifest = JSON.parse(await readFile(path, 'utf-8'))

    //     const parts             = manifest.version.split('.')
    //     parts[parts.length - 1] = (parseInt(parts[parts.length - 1], 10) + 1).toString()
    //     manifest.version        = parts.join('.')

    //     await writeFile(path, JSON.stringify(manifest, null, 4))
    // }
})

gulp.task('background', function () {
    return gulp.src('src/background/background.ts')
        .pipe(esbuild({
            bundle   : true,
            outfile  : 'background.js',
            platform : 'browser',
            target   : 'es2020',
            define   : {
                'import.meta.env.MODE'    : JSON.stringify(env.mode),
                'import.meta.env.BASE_URL': JSON.stringify(env.base_url),
                'import.meta.env.BASE_WS' : JSON.stringify(env.base_ws),
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
            try       { await axios.get('http://localhost:5000/api/debug/reload') }
            catch (e) { }
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
                'import.meta.env.MODE'    : JSON.stringify(env.mode),
                'import.meta.env.BASE_URL': JSON.stringify(env.base_url),
                'import.meta.env.BASE_WS' : JSON.stringify(env.base_ws),
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
            try       { await axios.get('http://localhost:5000/api/debug/reload') }
            catch (e) { }
            
        })
})

gulp.task('content', function () {
    return gulp.src('src/content/content.ts')
        .pipe(esbuild({
            bundle   : true,
            minify   : false, 
            sourcemap: false,
            outfile : 'content.js',
            platform: 'browser',
            target  : 'es2020',
            define  : {
                'import.meta.env.MODE'         : JSON.stringify(env.mode),
                'import.meta.env.BASE_URL'     : JSON.stringify(env.base_url),
                'import.meta.env.BASE_WS'      : JSON.stringify(env.base_ws),
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
            try       { await axios.get('http://localhost:5000/api/debug/reload') }
            catch (e) { }
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
gulp.task('run-prod', 
    gulp.series('background', 'popup', 'content', 'assets', 'manifest', 
        function () {
            env.mode      = 'production'
            env.base_url  = 'https://static.33.56.161.5.clients.your-server.de/dse'
            env.base_ws   = 'wss://static.33.56.161.5.clients.your-server.de/dse/ws'

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

gulp.task('run-dev', 
    gulp.series('background', 'popup', 'content', 'assets', 'manifest', 
        function () {
            env.mode      = 'development'
            env.base_url  = 'http://localhost:5000'
            env.base_ws   = 'ws://localhost:5000/ws'

            gulp.watch('src/background/**/*.ts', gulp.series('background'))
            gulp.watch('src/content/**/*.ts'   , gulp.series('content'))
            gulp.watch(['src/popup/**/*.ts', 'src/popup/**/*.svelte'], gulp.series('popup'))

            gulp.watch('public/**/*'      , gulp.series('assets'))
            gulp.watch('public/**/*.json' , gulp.series('manifest'))
    })
)

gulp.task('run-firefox', function (done) {
    env.mode     = 'development'
    env.base_url = 'http://localhost:5000'
    env.base_ws  = 'ws://localhost:5000/ws'

    // Start watching files and rebuild when they change
    gulp.watch('src/background/**/*.ts', gulp.series('background'))
    gulp.watch('src/content/**/*.ts'   , gulp.series('content'))

    gulp.watch(['src/popup/**/*.ts', 'src/popup/**/*.svelte'], gulp.series('popup')
    )

    gulp.watch('public/**/*'     , gulp.series('assets'))
    gulp.watch('public/**/*.json', gulp.series('manifest'))

    // Start web-ext in watch mode
    webExt.cmd.run({
        sourceDir : 'dist/firefox',
        watchFiles: ['dist/firefox/**/*'],
        noInput   : true, // Disable interactive prompts
        firefox   : 'C:/Users/Derwin/scoop/apps/firefox/current/firefox.exe'
    }).then((extensionRunner) => {
        // Extension runner is now running
        done()
    }).catch((err) => {
        console.error('Failed to run web-ext:', err)
        done(err)
    })
})


// ------------------------------------------------------------
// : Task - Build
// ------------------------------------------------------------
gulp.task('build', gulp.series('clean', 'versioning', 'background', 'popup', 'content', 'assets', 'manifest', 'zip'))
