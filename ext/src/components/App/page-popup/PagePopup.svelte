<script lang="ts">
    // ------------------------------------------------------------
    // Imports
    // ------------------------------------------------------------
    import { onMount } from 'svelte'

    import browser           from 'webextension-polyfill'

    import { State }                from '../utils/util'

    import { isEmpty     as is_empty }     from 'lodash'

    import { get } from 'lodash'

    import { lang }        from '../stores/store'
    import { storage }     from '../utils/util'
    import { log }         from '../utils/util'

    // ------------------------------------------------------------
    // Components
    // ------------------------------------------------------------
    import Button    from '../lib/Button.svelte'

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

            7: 'Verbinden...',
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

            7: 'Connecting...',
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
    
    // ------------------------------------------------------------
    // : Functions
    // ------------------------------------------------------------
    function change_lang(selected) {
        if (selected == 'en') {
            lang.set('en')
            lang_en = true
            lang_nl = false
            storage.set({'language': 'en'})
        } else if (selected == 'nl') {
            lang.set('nl')
            lang_en = false
            lang_nl = true
            storage.set({'language': 'nl'})
        }
    }


    function visit_contact() {
        browser.tabs.create({url: 'https://digitalepolarisatie.nl/'})
    }

    /**
     * Start search process (manually)
    */
    async function process_start() { 
        browser.runtime.sendMessage({target: 'search_process', action: 'start'})
    }

    /**
     * Stop search process (manually)
    */
    async function process_stop() {
        browser.runtime.sendMessage({target: 'search_process', action: 'stop'})
    }

    onMount(async () => {
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
        
        // Do dots animation (bottom left)
        let dots = '.'
        interval =  setInterval(() => {
            dots += '.'
            if (dots.length > 3) {
                dots = '.'
            }
            status = `Verbinden${dots}`
        }, 300)

        browser.runtime.onMessage.addListener(async(payload: object) => {
            const target                = get(payload, 'target')
            const search_process_state  = get(payload, 'state')

            if (!target || target !== 'popup') return

            switch (search_process_state) {
                case 'inactive':
                    state  = 'inactive'
                    break
                case 'active':
                    state  = 'active'
                    status = i18n[$lang][10]
                    break
                case 'cancelled':
                    state  = 'cancelled'
                    status = i18n[$lang][11]
                    break
            }
        })

        // Check if the server is online
        const is_connected = await storage.get('server.connected')
        if (is_connected) {
            light  = 'green'
            status = i18n[$lang][8]
        } else {
            light  = 'red'
            status = i18n[$lang][9]
        }

        clearInterval(interval)

        storage.onChanged.addListener(async () => {
            const is_connected         = await storage.get('server.connected')
            const search_process_state = await storage.get('search_process.state')
            if (is_connected) {
                state  = search_process_state
                light  = 'green'
                status = i18n[$lang][8]
            } else {
                state  = search_process_state === 'inactive' && !is_connected  ?  undefined : 'inactive'
                light = 'red'
                status = i18n[$lang][9]
            }
        })

        // Check if the search process is running
        {
            const search_process_state = await storage.get('search_process.state')
            switch (search_process_state) {
                case 'inactive':
                    state = search_process_state === 'inactive' && !is_connected  ?  undefined : 'inactive'
                    break
                case 'active':
                    state  = 'active'
                    status = i18n[$lang][10]
                    break
                case 'cancelled':
                    state  = 'cancelled'
                    status = i18n[$lang][11]
                    break
            }
        }
    })

    // ------------------------------------------------------------
    // : Reactive
    // ------------------------------------------------------------
    $: {
        log(state)
    }
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
    
            <div class="btn-wrapper">
                <!-- Idle -->
                {#if state == 'inactive'}
                    <Button primary={true} onclick={process_start}>{i18n[$lang][3]}</Button> 
                {/if}
        
                <!-- Busy -->
                {#if state == 'active'}
                    <Button secondary={true} onclick={process_stop}>{i18n[$lang][4]}</Button> 
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
                "t"
                "b";
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

            //@todo: remove below
            .subtitle {
                display        : flex;
                justify-content: flex-end;
                align-items    : center;

                padding-right: 1rem;
            }
        }
    }
</style>
    