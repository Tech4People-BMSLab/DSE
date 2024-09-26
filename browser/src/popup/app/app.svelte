<script lang="ts">
    // ------------------------------------------------------------
    // : Import
    // ------------------------------------------------------------
    import browser     from 'webextension-polyfill'
    import { onMount } from 'svelte'

    import { storage }   from '@popup/utils/utils'
    import { page }      from '@popup/stores/store'
    import { lang }      from '@popup/stores/store'
    import { ipc }       from '@popup/stores/store'
    import { connected } from '@popup/stores/store'

    // ------------------------------------------------------------
    // : Components
    // ------------------------------------------------------------
    import PageProcess   from '@popup/page-process/Index.svelte'           // Process page for the keyword searches
    import PagePopup     from '@popup/page-popup/Popup.svelte'             // Pop-up (top right)

    // ------------------------------------------------------------
    // : Helpers
    // ------------------------------------------------------------
    const sleep = (ms) => {
        return new Promise(resolve => setTimeout(resolve, ms))
    }

    const wait_until = async (predicate, timeout=-1, interval=100) => {
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

    // ------------------------------------------------------------
    // : Functions
    // ------------------------------------------------------------
    /**
     * This is a hack to detect whether this window is the popup, or a regular tab/window.
     * If it throws an error than it's the popup, otherwise it's a window/tab.
     */
     async function is_popup() {
        try {
            const tab = await browser.tabs.getCurrent()
            tab.width
            tab.height
            return false
        } catch (error) {
            return true
        }
    }

    /**
     * Check whether the search process is busy.
     */
    async function is_busy() {
        const data  = await storage.get_all()
        const state = data['search_proccess.state']
        log(data)
        log(state)
        return state === 'busy'
    }

    function log(...args) {
        browser.runtime.sendMessage({type: 'log', from: 'popup', data: args})
    }

    onMount(async () => {
        setTimeout(async () => {
            log($page)
        }, 1000)

        // Get the language from storage
        const language = await storage.get('language')

        // Set the language based on the user system language
        switch (language) {
            case 'nl':
            case 'nl-BE':
                lang.set('nl')
                break

            case 'uk':
            case 'en':
            case 'en-US':
            case 'en-GB':
            case 'en-CA':
            default:
                lang.set('en')
                break
        }

        if (!await is_popup() && await is_busy()) {
            page.set('process')
            return
        }

        const has_form       = await storage.get('user.form')
        const not_registered = !has_form

        if (not_registered) {
            if (await is_popup()) {
                browser.tabs.create({url: 'https://dev.bmslab.utwente.nl/dse/consent'})
            } else {
                page.set('process')
            }
        } else {
            if (await is_popup()) {
                page.set('popup')
            } else {
                page.set('process')
            }
        }
    })

    onMount(async () => {   
        await wait_until(() => $connected)

        $ipc.subscribe('bms.dse.popup.reload', {
            callback: () => {
                window.location.reload()
            },
        })
    })

</script>

<!-- Process Page (during search process) -->
{#if $page != 'popup'}
    <div class='window-wrapper'>
        {#if $page === 'process'}
            <PageProcess/>
        {/if}
    </div>
{/if}

<!-- Popup (Top right icon in the browser) -->
{#if $page === 'popup'}
    <div class="popup-wrapper">
        <PagePopup/>
    </div>
{/if}

<style type='text/scss'>
    .window-wrapper {
        width : 100vw;
        height: 100vh;
    }
    
    :root {
        @import url('https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap');
        @import url('https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap');
        @import url('https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100;0,300;0,400;0,500;0,700;0,900;1,100;1,300;1,400;1,500;1,700;1,900&display=swap');
    
        font-family: 'Roboto';
        font-style : normal;
    
        box-sizing: border-box;
    
        --black       : #000;
        --gray        : #848484;
        --purple-dark : 193498;
        --purple-light: #113cfc;
        
        --red-400: #ff5252;
        --red-500: #f44336;
    }
    
    :global(body, html) {
        margin : 0;
        padding: 0;
        
        overflow-x: hidden;
    }
    
    :global(div) {
        box-sizing: border-box;
    }
    
    // Width
    :global(::-webkit-scrollbar) {
        width : 0.5rem;
        height: 0.5rem;
    }
    // Track
    :global(::-webkit-scrollbar-track) {
        background-color: white;
    }
    
    // Handle
    :global(::-webkit-scrollbar-thumb) {
        background-color: var(--gray);
        border-radius: 3px;
    }
    // Handle on hover
    :global(::-webkit-scrollbar-thumb:hover) {
        background-color: var(--gray);
        border-radius: 3px;
    }
    
    :global(.bar) {
            height: 45px;
            width : 100%;
            
            display        : flex;
            justify-content: flex-start;
            align-items    : center;
    
            border-radius: 50px;
            padding-left : 1rem;
            padding-right: 1rem;
    
            overflow: hidden;
            box-sizing: border-box;
            
            // Font
            color: white; 
            font-family: 'Poppins';
            font-style : normal;
            font-weight: 700;
            font-size  : 18px;
            line-height: 60px;
    
            background: #193498; // 6377EE
    }
    
    :global(a) {
        color: #6377EE;
        font-family   : 'Roboto';
        font-style    : normal;
        font-weight   : 400;
        font-size     : 17px;
        line-height   : 32px;
        letter-spacing: 0.1em;
    }
    
    // Big bold text
    :global(h1) {
        width : 100%;
    
        margin : 0;
        
        color: #000000;
        font-family: 'Poppins';
        font-style: normal;
        font-weight: 700;
        font-size: 40px;
        line-height: 60px;
        text-align: center;
    }
    
    // Purple text (usually on top of heading)
    :global(h2) {
        width : 100%;
    
        margin: 0;
    
        color: #6377EE;
        font-style    : normal;
        font-weight   : 400;
        font-size     : 17px;
        line-height   : 32px;
        letter-spacing: 0.1em;
    }
    
    // Any paragraph (especially in the heading)
    :global(p) {
        display        : flex;
        justify-content: center;
        align-items    : center;
    
        color: #848484;
        font-style : normal;
        font-weight: 400;
        font-size  : 17px;
        text-align : center;   
        word-spacing: 0.1rem;
        line-height: 1.7rem;
    }
</style>
