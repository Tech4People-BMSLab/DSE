<script lang="ts">
    // ------------------------------------------------------------
    // : Initialization
    // ------------------------------------------------------------
    export let ref: any

    $: { // Validation on the postcode
        const re     = /^[1-9]?\d{1,4}$/
        const length = ref.answer.postcode.length 

        if (re.test(ref.answer.postcode) && length == 4) {
            ref.answer.valid = true
        } else {
            ref.answer.valid = false
        }
    }
    
</script>

<div class="component">
    <div class="title">
        {ref.question}
    </div>
    
    <div class="fields">
        <div class="field-number">
            <input type="text" placeholder="1234" maxlength="4" bind:value={ref.answer.postcode}>
        </div>

        <div class="field-verificator">
            <img 
                src="./images/check.svg"
                alt="check"
                class:valid   = {ref.answer.valid == true}
                class:invalid = {ref.answer.valid == false}
            />
        </div>
    </div>


</div>


<style lang="scss" scoped>

    .component {
        display: grid;
        grid-template: 
            "t" 1fr
            "w" 1fr / 1fr;

        padding-left  : 1rem;
        padding-right : 1rem;
        padding-top   : 1rem;
        padding-bottom: 1rem;

        box-sizing: border-box;

        .title {
            grid-area: t;

            display        : flex;
            justify-content: flex-start;
            align-items    : center;
        }

        .fields {
            grid-area: w;

            width : 338px;
            height: 40px;

            display: grid;
            grid-template:
                "l r" 1fr / auto auto;

            gap: 1rem;

            border       : 1px solid #e9e9e9;
            border-radius: 8px;
            box-shadow   : 0px 4px 15px rgba(0,0,0, 0.1);

            .field-number {
                grid-area: l;

                height: 2.5rem;
                
                box-sizing: border-box;

                input {
                    height: 100%;
                    width : 100%; 

                    margin : 0;
                    padding: 0;
                    border : 0;
                    outline: 0;
                    border-radius: 8px;

                    padding-left: 1rem;
                    margin-right: 1rem;

                    box-sizing: border-box;

                    font-size: 1.2rem;
                    color: black;

                    &:focus {
                        outline:rgba(0, 0, 0, 0.1) solid 1px;
                    }
                }
            }

            .field-verificator {
                grid-area: r;

                width : 35px;
                height: 2.5rem;

                display        : flex;
                justify-content: center;
                align-items    : center;

                box-sizing: border-box;

                .valid {
                    filter: invert(47%) sepia(11%) saturate(2531%) hue-rotate(72deg) brightness(126%) contrast(77%);
                }
            
                .invalid {
                    display: none;
                }

            }

        }
    }

</style>
