<script lang='ts'>
    // ------------------------------------------------------------
    // : Imports
    // ------------------------------------------------------------
    import { onMount } from 'svelte'

    import { lang } from '../stores/store'

    // ------------------------------------------------------------
    // : Props
    // ------------------------------------------------------------
    // Input variables
    export let title  = '<Title>'     // The title above the dropdown (usually a question)
    export let options: {label: string, value: string|number}[] = []

    // Bind variables
    export let answer       // The "value" of the selected option
    export let answers = {} // All the checkbox values

    // Trigger variables
    export let submit       // This variable is changed when the form is submitted

    // Local variable
    const i18n = {
        'nl': {
            1: 'Selecteer'
        },
        'en': {
            1: 'Select'
        }
    }
    let selected     = i18n[$lang][1]   // The selected option's label ('Selecteer' is the default)
    let opened       = null             // Whether the dropdown is open or not
    let submitted    = false            // Whether the form has been submitted or not (useful for invalid class on div)

    // ------------------------------------------------------------
    // : Functions
    // ------------------------------------------------------------
    /**
     * This handler listens for clicks outside of the div.
     * E.g. this is used for Dropdown to close when user clicks outside of its div.
     */
    function handler(node) {
        const handle_click = (event) => {
            if (!node.contains(event.target)) {
                const ce = new CustomEvent('outclick') // Triggers on:outclick
                node.dispatchEvent(ce)
            }
        }

        // Add listener to document
        document.addEventListener("click", handle_click, true)

        return {
            destroy() {
                // Remove listener from document on component removal from the DOM
                document.removeEventListener("click", handle_click, true)
            },
        }
    }

    function on_open() {
        opened = !opened
    }

    function on_select(option) {
        // Toggle the selected option
        answers[option.value] = !answers[option.value]

        // Falsify all the other options if 'Zeg ik liever niet' is selected
        if (option.value === 'unselected') {
            options.forEach(o => {
                if (o.value === 'unselected') {
                    answers[o.value] = true
                } else {
                    answers[o.value] = false
                }
            })
        } else {
            // Deselect the 'unselected'
            answers['unselected'] = false
        }

        // Generate the label on the selector (e.g. Option 1, Option 2, etc.)
        let str = ''
        for (const option of options) {
            let answer = answers[option.value]
            if (answer) {
                str += option.label + ', '
            }
        }

        // Remove last , 
        str = str.slice(0, -2)
        selected = str

        // Remove all options except 'Zeg ik liever niet' when it's selected

        // Use 'Selecteer' or 'Select' if it's empty
        if (str === '') {
            selected = i18n[$lang][1]
        }
    }

    function on_outclick() {
        opened = false // Close the dropdown if user clicked outside of it
    }

    onMount(() => {
        // Fill the answers array with false (same size as options)
        for (let i = 0; i < options.length; i++) {
            let option = options[i]
            answers[option.value] = false
        }
    })

    // ------------------------------------------------------------
    // : Reactive
    // ------------------------------------------------------------
    $: {
        answer = answers
    }

    $: {
        selected = i18n[$lang][1]   // The selected option's label ('Selecteer' is the default)
    }

    $: { // Check if the form has been submitted
        if (submit) {
            submitted = true
            submit    = false
        }
    }
</script>

<div class='wrapper' use:handler on:outclick={on_outclick}>

    <!-- Title (question) -->
    <div class='title'>
        <slot name='title'>{title}</slot>
    </div>

    <!-- Dropdown button -->
    <div class="selector-box">

        <!-- Left of arrow -->
        <div class="selector" 
            class:invalid = {submitted && !answer}
            on:click={() => on_open()}
        >
            <div class="selected">{selected}</div>
            <div class="arrow">
                <img 
                    src                = 'images/arrow.svg'
                    alt                = "arrow"
                    class:arrow-closed = {opened == false} 
                    class:arrow-opened = {opened == true}
                    >
            </div>
        </div>
    </div>

    <!-- Dropdown menu -->
    <div class="options"
        class:options-closed = {!opened} 
        class:options-opened = {opened}>

        {#if opened == true}
        {#each options as option}
            <div class="option" on:click={() => on_select(option)}>
                <div><input type="checkbox" bind:checked={answers[option.value]} /></div><div>{option.label}</div>
            </div>
        {/each}
        {/if}
    </div>

</div>

<style type='text/scss'>
    .wrapper {
        .title {
            width : 100%;
            padding: 0.5rem;
            box-sizing: border-box;
        }

        // ------------------------------------------------------------
        // : Dropdown box (with arrow)
        // ------------------------------------------------------------
        .selector-box {

            // Selector (left side of arrow)
            // TODO: Uncomment below
            // .invalid { outline: 1px solid #ff5252; }
            .selector {
                min-width: 150px;
                max-width: 350px;
                height   : 2.5rem;

                display: grid;
                grid-template-columns: 1fr auto;
                gap: 0px 0px;
                grid-auto-flow: row;
                grid-template-areas:
                    "l r";
                .selected  { grid-area: l; }
                .arrow     { grid-area: r; }

                border: 1px solid #e9e9e9;
                border-radius: 8px;
                box-shadow: 0px 4px 15px rgba(0, 0, 0, 0.1);

                background: #FFFFFF;

                .selected {
                    display        : flex;
                    justify-content: flex-start;
                    align-items    : center;

                    padding: 0.25rem;
                    padding-left : 0.5rem;
                }

                .arrow {
                    z-index: 1;
                    display        : flex;
                    justify-content: center;
                    align-items    : center;

                    padding: 0.25rem;

                    animation-duration : 100ms;
                    animation-fill-mode: forwards;

                    img {
                        width : 20px;
                        height: 20px;
                    }
                }

                &:hover {
                    cursor: pointer;
                    background: #F5F5F5;
                }
            }
        }


        // ------------------------------------------------------------
        // : Dropdown options (menu)
        // ------------------------------------------------------------
        .options {
            z-index: 100;
            position: relative;

            min-width: 150px;
            max-width: 350px;
            width    : 100%;

            display        : none;
            justify-content: flex-start;
            align-items    : flex-start;
            flex-direction : column;

            margin-top   : 0.5rem;
            margin-bottom: 0.5rem;

            border       : 1px solid #e9e9e9;
            border-radius: 8px;
            box-shadow   : 0px 4px 14px rgba(0, 0, 0, 0.1);

            overflow  : hidden;
            background: white;

            .option {
                width : 100%;
                min-height: 2.5rem;
                height    : 2.5rem;

                display: grid;
                grid-template-columns: max-content 1fr;
                grid-template-rows   : 100%;

                // Label
                div:nth-of-type(1) {
                    display        : flex;
                    justify-content: flex-start;
                    align-items    : center;

                    padding-left : 0.5rem;
                    padding-right: 0.5rem;

                    box-sizing: border-box;
                }

                // Checkmark
                div:nth-of-type(2) {
                    display        : flex;
                    justify-content: flex-start;
                    align-items    : center;

                    padding-right: 0.5rem;

                    box-sizing: border-box;


                    img {
                        height: 28px;
                    }
                }

                &:hover {
                    color: #FFFFFF;
                    background: #113CFC; // 6377EE
                }
            }
        }

        .options-opened {
            display: flex;
            animation-name     : slide-down;
            animation-duration : 100ms;
            animation-direction: normal;
            animation-fill-mode: forwards;
            animation-timing-function: linear;
            overflow-y: scroll;

            > * {
                display: flex;
            }
        }

        .options-closed {
            display: flex;
            animation-name     : slide-up;
            animation-duration : 100ms;
            animation-direction: normal;
            animation-fill-mode: forwards;
            animation-timing-function: linear;
            overflow-y: scroll;

            > * {
                display: none;
            }
        }

        // ------------------------------------------------------------
        // : Scroll
        // ------------------------------------------------------------
        ::-webkit-scrollbar-thumb {
            background-color: rgb(145, 145, 145);
        }

        ::-webkit-scrollbar-thumb:hover {
            background-color: rgb(128, 128, 128);
        }


        // ------------------------------------------------------------
        // : Dropdown animation
        // ------------------------------------------------------------
        @keyframes slide-up {
            0% { 
                opacity   : 1;
                max-height: 200px;
            }
            100% {
                opacity   : 0;
                max-height: 0;
            }
        }

        @keyframes slide-down {
            0% {
                opacity   : 0;
                max-height: 0;
            }
            100% {
                opacity   : 1;
                max-height: 200px;
            }
        }

         // ------------------------------------------------------------
         // : Arrow animation
         // ------------------------------------------------------------
        .arrow-opened {
            animation-name: rotate-arrow;
            animation-duration: 100ms;
            animation-direction: reverse;
            animation-fill-mode: forwards;
        }

        .arrow-closed {
            animation-name: rotate-arrow-reverse;
            animation-duration: 100ms;
            animation-direction: normal;
            animation-fill-mode: forwards;
        }

        // @keyframes for arrow-up and arrow-down
        @keyframes rotate-arrow {
            0%   { 
                opacity: 1;
                transform: rotate(180deg); 
            }
            100% { 
                opacity: 1;
                transform: rotate(0deg); 
            }
        }
        
        @keyframes rotate-arrow-reverse {
            0%   { 
                opacity: 1;
                transform: rotate(180deg); 
            }
            100% { 
                opacity: 1;
                transform: rotate(0deg); 
            }
        }
    }
    

    
</style>