<script lang="ts">
    // ------------------------------------------------------------
    // : Imports
    // ------------------------------------------------------------
    import { fly }                       from 'svelte/transition';
    import { elasticInOut, elasticOut }  from 'svelte/easing'
    import { error_message }             from '../stores/store'
    import { success_message }           from '../stores/store'

    const delay  = 3000;
    let show     = false
    let error    = false
    let success  = false
    let interval = null

    // ------------------------------------------------------------
    // : Reactive
    // ------------------------------------------------------------
    $: {
        if ($error_message) {
            show    = true
            error   = true
            success = false
            clearInterval(interval)
            interval = setTimeout(() => {
                show = false
                error_message.set(null)
            }, delay)
        } else if ($success_message) {
            show    = true
            success = true
            error   = false
            clearInterval(interval)
            interval = setTimeout(() => {
                show = false
                success_message.set(null)
            }, delay)
        } else {
            show = false
        }
    }
</script>

<div class="screen-container">
    {#if show}
        <div class="popup-container">
            <div 
                class        = "popup"
                class:error  = {error} 
                class:success= {success} 
                in:fly  = {{y:-100,delay:-500, duration:1000, easing:elasticInOut}}
                out:fly = {{y:-100,delay:0,    duration:750,  easing:elasticOut}}
                >

                {#if success}
                    {$success_message}
                {/if}

                {#if error}
                    {$error_message}
                {/if}
            </div>
        </div>
    {/if}
</div>

<style type="text/scss">
    .screen-container {
        position: fixed;
        top : 0;
        left: 0;

        width : 100vw;
        height: 100vh;

        pointer-events: none;

        .popup-container {
            height: 4rem;

            display        : flex;
            justify-content: center;
            align-items    : center;

            .popup {
                min-width: 20rem;
                width: max-content;
                max-width: 50rem;
                height: 1.8rem;

                display        : flex;
                justify-content: center;
                align-items    : center;

                border-radius: 8px;
                padding-left : 1rem;
                padding-right: 1rem;

                color: white;
                font-size: 14px;
                font-weight: 600;

                background: #ff5252;
            }

            .error {
                background: #ff5252;
            }

            .success {
                background: #1EB135;
            }
        }

    }


</style>