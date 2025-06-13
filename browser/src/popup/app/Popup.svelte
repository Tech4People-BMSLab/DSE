<script lang="ts">
    import axios          from 'axios'
    import { onMount }    from 'svelte'

    import { not_empty } from '@popup/utils/utils'

    import { store } from '@/background/core/storage'
    import { ipc }   from '@popup/utils/ipc'

    import Image from '@popup/lib/Image.svelte'
    import Flex  from '@popup/lib/Flex.svelte'
    import Grid  from '@popup/lib/Grid.svelte'
    import Stack from '@popup/lib/Stack.svelte'
    import Layer from '@popup/lib/Layer.svelte'

    import Status    from '@popup/lib/Status.svelte'
    import StatusBar from '@popup/lib/StatusBar.svelte'

    let state = $state({
        version  : '',
        language : 'nl',
        token    : '',
        connected: null, // null = unknown, true = connected, false = disconnected
        year     : new Date().getFullYear(),

        ready     : false,
        registered: false,
    })

    const select_language = (lang: string) => {
        state.language = lang
    }

    const register = () => {
        ipc.send('consent')
    }
    
    onMount(async () => {
        ipc.init()

        state.version    = await store.get('extension.version')
        state.token      = await store.get('user.token')
        state.registered = not_empty(await store.get('user.form'))

        await new Promise(async (resolve, reject) => {
            const url = 'https://static.33.56.161.5.clients.your-server.de/dse/api'

            axios({
                method: 'get',
                url   : url,
            })
            .then(response => {
                if (response.status == 200) {
                    state.connected = true
                } else {
                    state.connected = false
                }
            })
            .catch(error => {
                state.connected = false
            })
        })
    })
    
</script>

<main class="page">
    <Grid>
        <Flex cls="block-0" position="tr" direction="row">
            <Grid cols="1fr 1fr">
                <Flex position="tl" direction="row">Version: <strong>{state.version}</strong></Flex>
                <Flex position="tr" direction="row">
                    <button class:active={state.language == "nl"} onclick={() => select_language('nl')}>NL</button>
                    <div>|</div>
                    <button class:active={state.language == "en"} onclick={() => select_language('en')}>EN</button>
                </Flex>
            </Grid>
        </Flex>
        <Flex cls="block-1" position="cc"><img src="./images/logo/dp_logo.png" alt="DP Logo"/></Flex>
        <Flex cls="block-2" position="cc">
            {#if state.language == "nl"}De browser extensie zal automatisch zoekopdrachten uitvoeren op de achtergrond. Dit gebeurt wanneer jouw browser actief is.{/if}
            {#if state.language == "en"}The browser extension will automatically perform search queries in the background. This happens when your browser is active.{/if}
        </Flex>
        <Flex cls="block-3">
            {#if state.language == "nl"}Mocht je problemen ervaren of vragen hebben dan kan je altijd contact met ons opnemen.{/if}
            {#if state.language == "en"}If you experience problems or have questions, you can always contact us.{/if}
        </Flex>
        <Flex cls="block-4" position="cc" direction="row">
            <input type="text" value={`ID: ${state.token}`} readonly/>
        </Flex>
        <Flex cls="block-5" position="cc" direction="row">
            {#if state.language == 'nl'}<button class:disabled={true}>Testen</button>{/if}
            {#if state.language == 'en'}<button class:disabled={true}>Test</button>{/if}

            {#if state.language == 'nl'}<button class:primary={state.registered == false} onclick={register}>Registreren</button>{/if}
            {#if state.language == 'en'}<button class:primary={state.registered == false} onclick={register}>Register</button>{/if}
        </Flex>
        <Flex cls="block-6" position="tl">
            <Grid cols="1fr 1fr">
                <Status state={state}/>
                <Flex position="br">Digitale Polarisatie {state.year}</Flex>
            </Grid>
        </Flex>
        <Flex cls="block-7"><StatusBar state={state}/></Flex>
    </Grid>
</main>

<style lang="scss" global>
    @import url('https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap');

    .page {
        width: 25rem;
        border-radius: 5px;

        font-size  : 12px;
        font-family: 'Poppins';

        overflow: hidden;

        .block-0 {
            padding-top: 0.5rem;
            padding-left : 1rem;
            padding-right: 1rem;

            strong {
                padding-left : 0.25rem;
            }

            button {
                display        : flex;
                justify-content: center;
                align-items    : center;

                border : none;
                outline: none;
                
                cursor: pointer;
                
                font-size: 12px;
                background: transparent;

                &.active {
                    font-weight: 600;
                }
            }
        }

        .block-1 {
            padding-top   : 1rem;
            padding-bottom: 1rem;

            img {
                width : 100%;
                height: 5rem;
                object-fit: contain;
            }
        }

        .block-2 {
            padding-top   : 0.5rem;
            padding-bottom: 0.5rem;
            padding-left : 2rem;
            padding-right: 2rem;
        }

        .block-3 {
            padding-top   : 0.5rem;
            padding-bottom: 0.5rem;
            padding-left : 2rem;
            padding-right: 2rem;
        }

        .block-4 {
            border-top   : 1px solid #f0f0f0;
            border-bottom: 1px solid #f0f0f0;
            padding-left : 2rem;
            padding-right: 2rem;

            input[type="text"] {
                width : 100%;
                height: 2.5rem;
                
                border : none;
                outline: none;
                
                color: #606060;
                font-weight: 500;
                font-size  : 18px;
                line-height: 60px;
                text-align: center;
            }
        }

        .block-5 {
            padding-top   : 1rem;
            padding-bottom: 1rem;

            gap: 1rem;

            button {
                width : 10rem;
                height: 2.5rem;
                
                display        : flex;
                justify-content: center;
                align-items    : center;

                border-radius: 50px;
                border : none;
                outline: none;

                cursor: pointer;
                
                font-weight: 500;
                font-size  : 18px;
                line-height: 60px;

                &.primary {
                    color     : #ffffff;
                    background: #113CFC; // 6377EE

                    cursor: pointer;

                    &:hover {
                        color: #113CFC; // 6377EE
                        background: #ffffff;
                    }
                }

                &.disabled {
                    color     : #848484;
                    background: #fcfcfc;

                    cursor: not-allowed;

                    &:hover {
                        border: 1px solid #848484;
                    }
                }
            }
        }

        .block-6 {
            padding-top   : 0.5rem;
            padding-bottom: 0.5rem;
            padding-left : 1rem;
            padding-right: 1rem;
        }

        .block-7 {
            height: auto;
        }
    }
</style>
