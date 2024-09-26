<script lang='ts'>
    // ------------------------------------------------------------
    // : Imports
    // ------------------------------------------------------------
    import { lang } from '../stores/store'
    
    // ------------------------------------------------------------
    // : Props
    // ------------------------------------------------------------
    // Message of the user
    export let text   = ''

    // Bind variable
    export let submit
    let submitted = false

    // Text properties
    export let length
    export let min_length
    export let max_length

    // State variable
    let verified = -1

    const i18n = {
        'nl': {
            1: 'Type hier je tekst...'
        },
        'en': {
            1: 'Type your text here...'
        }
    }
 
    // : Reactive
    // ------------------------------------------------------------
    $: {
        if (text.length < min_length) {
            verified = 0
        } else {
            verified = 1
        }
    }

    $: {
        length = text.length
    }

    $: {
        if (submit) {
            submitted = true
            submit    = false
        }
    }
</script>

<div class='textarea'>
    <div class="input-box" class:invalid={submitted && verified == 0}>
        <textarea class='input' placeholder={i18n[$lang][1]} maxlength={max_length} bind:value={text}></textarea>
    </div>
</div>

<style type='text/scss'>
    .textarea {
        width : 100%;
        height: 100%;

        // Selector (with arrow)
        .invalid { outline: 1px solid red; }
        .input-box {
            width : 100%;
            min-height: 200px;
            max-height: 500px;

            display        : flex;
            justify-content: flex-start;
            align-items    : flex-start;

            border       : 1px solid #e9e9e9;
            border-radius: 8px;
            box-shadow   : 0px 4px 15px rgba(0, 0, 0, 0.1);
            padding: 1rem;

            background: #FFFFFF;

            textarea {
                width : 100%;
                min-height: 200px;
                max-height: 300px;
                
                margin : 0;
                padding: 0;
                border : 0;
                // outline: 0;


                padding: 0.5rem;
                box-sizing: border-box;

                color: black;
                font-family: 'Montserrat'; 
                font-size: 1rem;

                resize: none;

                &:focus {
                    outline:rgba(0, 0, 0, 0.1) solid 1px;
                }
            }
        }
    }
    

    
</style>