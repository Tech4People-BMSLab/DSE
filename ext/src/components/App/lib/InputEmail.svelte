<script lang='ts'>
    // ------------------------------------------------------------
    // Imports
    // ------------------------------------------------------------
    import {
        is_empty
    } from '../utils/util'

    // ------------------------------------------------------------
    // Props
    // ------------------------------------------------------------
    export let title = ''
    export let email = ''
    export let submit
    let submitted = false

    let verified = -1

    // ------------------------------------------------------------
    // : Reactive
    // ------------------------------------------------------------
    $: { // Verify email
        const re = /^(([^<>()[\]\.,;:\s@\"]+(\.[^<>()[\]\.,;:\s@\"]+)*)|(\".+\"))@(([^<>()[\]\.,;:\s@\"]+\.)+[^<>()[\]\.,;:\s@\"]{2,})$/i
        
        if (re.test(email)) {
            verified = 1
        } else {
            verified = -1
        }
    }

    $: {
        if (submit) {
            submitted = true
            submit    = false
        }
    }
</script>

<div class='dropdown'>

    <!-- Title (question) -->
    <div class='title'>
        {#if !is_empty(title)}
            <slot name='title'>{title}</slot>
        {/if}
    </div>

    <!-- Input box -->
    <div class="input-box" class:invalid={submitted && verified == -1}>
        <div class="input-text">
            <input class='input' type="text" placeholder="john@example.com" bind:value={email}/>
        </div>

        <div class="verificator">
            <img 
                src='./images/check.svg'
                alt="check"
                class:verified     = {verified ==  1}
                class:not-verified = {verified == -1}
            />
        </div>
    </div>
</div>

<style type='text/scss'>
    .dropdown {

        .title {
            padding: 0.5rem;
        }

        .invalid {
            outline: 1px solid red;
        }
        .input-box {
            width : 100%;
            min-width: 200px;
            max-width: 300px;
            height   : 2.5rem;

            display: grid;
            grid-template-columns: auto 32px;
            gap: 0px 1px;
            grid-auto-flow: row;
            grid-template-areas:
                "l r";
            .input-text    { grid-area: l; }
            .verificator   { grid-area: r; }

            padding-left : 1rem;
            padding-right: 1rem;

            border       : 1px solid #e9e9e9;
            border-radius: 8px;
            box-shadow   : 0px 4px 15px rgba(0, 0, 0, 0.1);

            background: #FFFFFF;

            .input-text {
                width : 100%;

                input {
                    width : 100%;
                    min-width: 200px;
                    max-width: 400px;
                    height: 100%;
    
                    margin : 0;
                    padding: 0;
                    border : 0;
    
                    box-sizing: border-box;
    
                    color: black;
                    font-family: 'Montserrat';
                    font-size: 1rem;
    
                    &:focus {
                        outline:rgba(0, 0, 0, 0.1) solid 1px;
                    }
                }
            }

            > .verificator {
                width : 32px;
                height: 100%;

                display        : flex;
                justify-content: center;
                align-items    : center;

                img {
                    width : 32px;
                    height: 32px;
                }

                .verified {
                    filter: invert(47%) sepia(11%) saturate(2531%) hue-rotate(72deg) brightness(126%) contrast(77%);
                }
                .not-verified {
                    opacity: 0;
                }

            }
        }
    }
    

    
</style>