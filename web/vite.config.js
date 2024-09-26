import { defineConfig } from 'vite'
import { svelte }       from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig(async({command, mode}) => {
    const is_production  = mode === 'production'
    const is_development = mode === 'development'
    const base           = is_production ? '/dse' : '/'

    switch (true) {
        case is_production: {
            return {
                plugins: [svelte()],
                define: {
                    'import.meta.env.MODE'     : JSON.stringify('production'),
                    'import.meta.env.BASE_URL' : JSON.stringify(base),
                },
                base   : base,
                resolve: {
                    alias: {
                        '@': new URL('./src', import.meta.url).pathname
                    }
                },
                build: {
                    minify: false,
                    rollupOptions: {
                        treeshake: true
                    }
                }
            }
        }

        case is_development: {
            return {
                plugins: [svelte()],
                define : {},
                base   : base,
                resolve: {
                    alias: {
                        '@': new URL('./src', import.meta.url).pathname
                    }
                },
                build: {
                    minify: false,
                    rollupOptions: {
                        treeshake: false
                    }
                }
            }
        }   
    }
})
