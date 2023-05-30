<script lang='ts'>
    // ------------------------------------------------------------
    // Props
    // ------------------------------------------------------------
    export let title  // The question to ask the user
    export let answer // Answer selected by the user
    export let submit // A bind variable that is used from parent component to notify when the form is submitted

    let value_numbers = '' // Contains the numbers section of the postcode
    let value_letters = '' // Contains the letter section of the postcode

    let verified_numbers = -1
    let verified_letters = -1
    let verified         = -1

    let submitted = false

    // ------------------------------------------------------------
    // : Reactive
    // ------------------------------------------------------------
    $: { // Check for the 4 digits in postcode
        const re = /^\d{1,4}$/

        if (re.test(value_numbers)) {
            if (value_numbers.length == 4) {
                verified_numbers = 1

                // Move the focus to letter
                const el = document.querySelector('.input-letters input') as HTMLElement
                el.focus()
            } else {
                verified_numbers = 0
            }
        } else {
            value_numbers = value_numbers.slice(0, -1)
            verified_numbers = 0 
        }
    }

    $: { // Check for the 2 alphabetic characters in postcode
        const re = /^[a-zA-Z]{1,2}$/

        if (re.test(value_letters)) {
            value_letters = value_letters.toUpperCase() // Capitalize

            if (value_letters.length == 2) {
                verified_letters = 1
            } else {
                verified_letters = 0
            }
        } else {
            value_letters = value_letters.slice(0, -1)
            verified_letters = 0
        }
    }

    $: { // Check if both digits and letters are filled, and set the answer
        verified = (verified_numbers == 1 && verified_letters == 1) ? 1 : 0
        if (verified) {
            answer = {value: value_numbers + value_letters}
        }
    }

    $: { // Check if submit has been triggered to change the submitted state
        if (submit) {
            submitted = true
            submit    = false
        }
    }
</script>

<div class='wrapper'>

    <!-- Title (question) -->
    <div class='title'>
        <slot name='title'>{title}</slot>
    </div>

    <!-- Input box -->
    <div class="input-box">
        <div class="fields" class:invalid={submitted && !answer}>
            <div class="input-numbers">
                <input class='input' type="text" placeholder='1234' maxlength="4" bind:value={value_numbers}/>
            </div>
    
            <div class="input-letters">
                <input class='input' type="text" placeholder='AB' maxlength="2" bind:value={value_letters}/>
            </div>
    
            <div class="verificator">
                <img 
                    src='./images/check.svg'
                    alt="check"
                    class:verified     = {verified == 1}
                    class:not-verified = {verified == 0}
                />
            </div>
        </div>
    </div>
</div>

<style type='text/scss'>
    .wrapper {
        .title {
            padding: 0.5rem;
        }

        // ------------------------------------------------------------
        // Post code input box
        // ------------------------------------------------------------
        .input-box {
            .invalid { outline: 1px solid #ff5252 }
            .fields {
                min-width: 150px;
                max-width: 350px;
                height   : 2.5rem;

                display: grid;
                grid-template-columns: minmax(100px, 175px) minmax(100px, 175px) 2rem 1fr;
                gap: 0px 1px;
                grid-auto-flow: row;
                grid-template-areas:
                    "l c r .";
                .input-numbers { grid-area: l; }
                .input-letters { grid-area: c; }
                .verificator   { grid-area: r; }

                border       : 1px solid #e9e9e9;
                border-radius: 8px;
                box-shadow   : 0px 4px 15px rgba(0, 0, 0, 0.1);

                background: #FFFFFF;

                .input {
                    height: 100%;
                    width : 100%;

                    margin : 0;
                    padding: 0;
                    border : 0;
                    outline: 0;
                    border-radius: 8px;

                    font-size: 1.2rem;
                    color: black;

                    &:focus {
                        outline:rgba(0, 0, 0, 0.1) solid 1px;
                    }
                }

                > .input-numbers {
                    .input {
                        text-align: center;
                    }
                }

                > .input-letters {
                    .input {
                        text-align: center;
                    }
                }

                > .verificator {
                    display        : flex;
                    justify-content: center;
                    align-items    : center;

                    .verified {
                        filter: invert(47%) sepia(11%) saturate(2531%) hue-rotate(72deg) brightness(126%) contrast(77%);
                    }
                    .not-verified {
                        display: none;
                    }

                }
            }
        }
    }
    
</style>