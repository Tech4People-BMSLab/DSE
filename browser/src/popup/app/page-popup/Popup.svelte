<script lang="ts">
    // ------------------------------------------------------------
    // Imports
    // ------------------------------------------------------------
    import { onMount } from 'svelte'

    import browser           from 'webextension-polyfill'

    import { storage }    from '@popup/utils/utils'
    import { sleep }      from '@popup/utils/utils'
    import { wait_until } from '@popup/utils/utils'

    import { ipc  }       from '@popup/stores/store'
    import { lang }       from '@popup/stores/store'
    import { connected }  from '@popup/stores/store'

    // ------------------------------------------------------------
    // Components
    // ------------------------------------------------------------
    import Button    from '@popup/lib/Button.svelte'

    // ------------------------------------------------------------
    // Props
    // ------------------------------------------------------------
    const i18n = {
        'nl': {
            1: 'De browser extensie zal automatisch zoekopdrachten uitvoeren op de achtergrond. Dit gebeurt wanneer jouw browser actief is.',
            2: 'Mocht je problemen ervaren of vragen hebben dan kan je altijd contact met ons opnemen.',

            3: 'Test',
            4: 'Annuleren',
            5: 'Even geduld',

            6: 'Probleem melden',

            7: 'Verbinden',
            8: 'Verbonden',
            9: 'Verbinden mislukt, opnieuw proberen',
            10: 'Bezig...',
            11: 'Zoeken afgebroken'
        },
        'en': {
            1: 'The browser extension will automatically perform search queries in the background. This happens when your browser is active.',
            2: 'If you experience problems or have questions, you can always contact us.',

            3: 'Test',
            4: 'Cancel',
            5: 'Please wait',

            6: 'Report problem',

            7: 'Connecting',
            8: 'Connected',
            9: 'Connection failed, try again',
            10: 'Busy...',
            11: 'Search aborted'
        }
    }
    let state    = undefined
    let interval = undefined
    let light    = 'yellow'
    let status   = i18n[$lang][7]

    let lang_en = false
    let lang_nl = false

    let version = ''
    let token   = ''
    // ------------------------------------------------------------
    // : Functions
    // ------------------------------------------------------------
    function log(...args) {
        browser.runtime.sendMessage({type: 'log', from: 'popup', data: args})
    }

    function change_lang(selected) {
        if (selected == 'en') {
            lang.set('en')
            lang_en = true
            lang_nl = false
            storage.set({'language': 'en'})
            log('Changed language to English')
        } else if (selected == 'nl') {
            lang.set('nl')
            lang_en = false
            lang_nl = true
            storage.set({'language': 'nl'})
            log('Changed language to Dutch')
        }
    }

    function visit_contact() {
        browser.tabs.create({url: 'https://digitalepolarisatie.nl/'})
    }

    async function crawler_start() { 
        log('Starting crawler')
        browser.runtime.sendMessage({type: 'action', from: 'popup', action: 'start_crawler'})
    }

    async function crawler_stop() {
        log('Stopping crawler')
        browser.runtime.sendMessage({type: 'action', from: 'popup', action: 'stop_crawler'})
    }

    onMount(async () => {
        token   = await storage.get('user.token')
        version = import.meta.env.IPC_VERSION

        // Get the current language
        switch ($lang) {
            case 'en':
                lang_en = true
                lang_nl = false
                break
            case 'nl':
                lang_en = false
                lang_nl = true
                break
        }
    })

    onMount(async () => {
        // Connecting (animation)
        let dots = '.'
        interval =  setInterval(() => {
            dots += '.'
            if (dots.length > 3) {
                dots = '.'
            }
            status = `${i18n[$lang][7]}${dots}`
        }, 300)

        await wait_until(async () => await storage.get('ipc.connected'))

        clearInterval(interval)
        
        light  = 'green'
        status = i18n[$lang][8]
    })

    onMount(async () => {
        await wait_until(async () => await storage.get('ipc.connected'))

        while (true) {
            const active = await storage.get('crawler.active')

            switch (true) {
                case active: {
                    state  = 'active'
                    status = i18n[$lang][10]
                    break
                }
                default: {
                    state  = 'inactive'
                    status = i18n[$lang][8]
                    break
                }
            }
            
            await sleep(1000)
        }
    })

    // ------------------------------------------------------------
    // : Reactive
    // ------------------------------------------------------------
    $: {
        switch ($lang) {
            case 'en':
                lang_en = true
                lang_nl = false
                break
            case 'nl':
                lang_en = false
                lang_nl = true
                break
        }
    }
</script>

<div class="popup-container">
    <div class="popup">

        <div class="logo">
            <img src='./images/logo/dp_logo.png' alt='DP logo'>
        </div>

        <div class="intro">
            <p>{i18n[$lang][1]}</p>
            <p>{i18n[$lang][2]}</p>
        </div>
    
        <div class="control">
            
            <div class="identifier">
                <label for="identifier">ID:</label>
                <input type="text" id="identifier" name="identifier" bind:value={token} disabled>
            </div>
    
            <div class="btn-wrapper">
                <!-- Idle -->
                {#if state == 'inactive'}
                    <Button primary={true} onclick={crawler_start}>{i18n[$lang][3]}</Button> 
                {/if}
        
                <!-- Busy -->
                {#if state == 'active'}
                    <Button secondary={true} onclick={crawler_stop}>{i18n[$lang][4]}</Button> 
                {/if}

                <!--  Cancelled -->
                {#if state == 'cancelled'}
                    <Button secondary={true}>{i18n[$lang][5]}</Button> 
                {/if}
            </div>

            <div class="link" on:click={visit_contact}>{i18n[$lang][6]}</div>
        </div>
    
        <div class="status">
            <div class="indicator">
                <div
                    class="circle"
                    class:yellow = {light === 'yellow'}
                    class:red    = {light === 'red'   }
                    class:green  = {light === 'green' }
                />
            </div>
            <div class="text">{status}</div>
            <div class="language">
                <span on:click={() => change_lang('nl')} class:active={lang_nl}>Nederlands</span>
                <span on:click={() => change_lang('en')} class:active={lang_en}>English (UK)</span>
            </div>
        </div>
    </div>
</div>


<style type="text/scss">
    .popup {
        width : 450px;

        display: grid;
        grid-template-columns: auto;
        grid-template-rows   : auto auto auto auto;
        gap: 0px 0px;
        grid-auto-flow: row;
        grid-template-areas:
            "logo"
            "intro"
            "control"
            "status";

        .logo    { grid-area: logo;    }
        .intro   { grid-area: intro;   }
        .control { grid-area: control; }
        .status  { grid-area: status;  }

        .logo {
            display        : flex;
            justify-content: center;
            align-items    : center;

            padding-top   : 1rem;
            padding-bottom: 1rem;

            img {
                width: 200px;
            }
        }

        .intro {
            width : 100%;

            display        : flex;
            justify-content: flex-start;
            align-items    : flex-start;
            flex-direction : column;

            color: #848484;
            font-family   : 'Roboto';
            font-style : normal;
            font-weight: 400;
            font-size  : 17px;
            line-height: 32px;

            p {
                max-width: 100%;

                padding: 0;
                margin: 0;

                padding-left : 3rem;
                padding-right: 3rem;
                padding-top  : 1rem;

                text-align: left;
                white-space: pre-wrap;
                word-break: break-word;
            }
        }

        .control {
            display        : flex;
            justify-content: center;
            align-items    : center;
            flex-direction : column;

            display: grid;
            gap    : 0px 0px;
            grid-template-columns: repeat(1, 1fr);
            grid-template-rows   : 1fr 4rem;
            grid-auto-flow: row;
            grid-template-areas:
                "i"
                "t"
                "b";
            .identifier  { grid-area: i;}
            .btn-wrapper { grid-area: t;}
            .link        { grid-area: b;}

            color: #848484;
            font-family: 'Poppins';
            font-style : normal;
            font-weight: 600;
            font-size  : 18px;
            line-height: 60px;

            .btn-wrapper {
                height: 100%;
                width : 100%;

                display        : flex;
                justify-content: center;
                align-items    : flex-end;

                padding-top   : 1.5rem;
                box-sizing: border-box;
            }

            .identifier {
                display        : flex;
                justify-content: center;
                align-items    : center;

                border-top   : 1px #eee solid;
                border-bottom: 1px #eee solid;

                margin-top   : 1rem;

                background-color: #f9f9f9;

                label {
                    height: 100%;

                    font-weight: bold;
                    margin-right: 10px;
                    color: #333;
                }

                input {
                    height: 100%;
                    width : 30%;

                    padding      : 0.5rem;
                    border       : 1px solid #ddd;
                    border-radius: 3px;

                    font-size: 1em;

                    background: white;

                    cursor: text;
                }
            }

            .link {
                text-align : center;
                font-weight: 500;
                font-size  : 16px;

                cursor: pointer;

                &:hover {
                    color: #000000;
                }
            }
        }

        // The footer (ext)
        .status {
            width : 100%;
            height: 2rem;
            
            display: grid;
            grid-template-columns: 2rem 1fr auto;
            gap: 0px 0px;
            grid-auto-flow: row;
            grid-template-areas:
                "l c r";

            .indicator {
                width : 100%;
                height: 100%;

                display        : flex;
                justify-content: center;
                align-items    : center;

                .circle {
                    width : 15px;
                    height: 15px;
    
                    content: '';
                    border-radius: 50%;
                }
                .green  { background: #1DCE00; }
                .yellow { background: #FFB800; }
                .red    { background: #CE0000; }
            }
            
            .text {
                display        : flex;
                justify-content: flex-start;
                align-items    : center;
            }

            .language {
                display        : flex;
                justify-content: flex-end;
                align-items    : center;

                padding-right: 0.5rem;

                cursor: pointer;

                gap: 0.5rem;

                .active {
                    font-weight: 900;
                }
            }
        }
    }
</style>
    