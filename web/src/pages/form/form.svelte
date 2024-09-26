<script lang="ts">
    // ------------------------------------------------------------
    // : Imports
    // ------------------------------------------------------------
    import { onMount }     from 'svelte'
    import { navigate }    from 'svelte-routing'
    import { map, filter } from '@/utils/utils'

    import Header from '@/components/header.svelte'
    import Footer from '@/components/footer.svelte'

    import FieldPostcode      from '@/components/fields/postcode.svelte'
    import FieldDropdown      from '@/components/fields/dropdown.svelte'
    import FieldDropdownMulti from '@/components/fields/dropdown-multi.svelte'

    import { every } from '@/utils/utils'

    import storage from '@/services/storage'

    import form   from '@/stores/form'
    import global from '@/stores/global'
    // ------------------------------------------------------------
    // : Init
    // ------------------------------------------------------------
    const base_path = import.meta.env.BASE_URL || '';

    let questions: any[]

    async function on_finish() {
        if (!$form.completed) { return } // If form is not completed

        await storage.set('form', $form)
        navigate(`${base_path}/complete`)
    }

    async function on_decline() {

    }

    $: { // Handle language change
        questions = map($form.questions, (q: { [key: string]: any }) => {
            return {
                ...q,
                question: q['question'][$global.language],
            }
        })
    }

    $: { // Handle validation
        $form.completed = every($form.questions, (q: { [key: string]: any }) => {
            return q.answer.valid == true
        })
    }

</script>

<!-- svelte-ignore empty-block -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<main class="page-demographic">
    
    <div class="header">
        <Header/>
    </div>

    <div class="content">
        <h1>Details</h1>
        <p>To compare the data we ask you to fill in some demographic details. These details cannot be traced back to you. If you do not want to fill in certain details you can select this.</p>

        <div class="form">
            {#each questions as question}

                {#if question['type'] == 'postcode'}
                    <FieldPostcode bind:ref={question}/>
                {/if}

                {#if question['type'] == 'dropdown'}
                    <FieldDropdown bind:ref={question}/>
                {/if}

                {#if question['type'] == 'dropdown-multi'}
                    <FieldDropdownMulti bind:ref={question}/>
                {/if}
                
            {/each}
        </div>

        <p>Thank you for filling in your demographic data. Click on finish to complete the installation of the browser extension.</p>

        <div class="btn-group">
            <button class="primary" class:active={$form.completed} on:click={() => on_finish()}>Finish</button>
            <button class="secondary"   on:click={() => on_decline()}>Decline</button>
        </div>

    </div>

    <div class="footer">
        <Footer/>
    </div>
</main>


<style lang="scss" scoped>
.page-demographic {
    position: absolute;
    top : 0;
    left: 0;

    width : 100%;
    height: 100%;

    display: grid;
    grid-template:
        "h" auto
        "c" 1fr
        "f" auto / 1fr;

    justify-items: center;

    font-family: 'Poppins';

    .header {
        grid-area: h;
    }

    .content {
        grid-area: c;

        display: grid;
        grid-template:
            "r1" min-content        // Title
            "r2" min-content        // Description
            "r3" min-content        // Form
            "r4" min-content        // Thank you
            "r5" min-content / 1fr; // Button group
        
        h1 {
            grid-area: r1;

            display        : flex;
            justify-content: center;
            align-items    : center;
        }

        p:nth-of-type(1) {
            grid-area: r2;

            display        : flex;
            justify-content: center;
            align-items    : center;
            flex-wrap      : wrap;
        }

        .form {
            grid-area: r3;

            display        : grid;
            grid-template-columns: repeat(2, 1fr);
            grid-auto-rows: 100px;
            grid-gap: 1rem;
        }

        p:nth-of-type(2) {
            grid-area: r4;

            display        : flex;
            justify-content: center;
            align-items    : center;

            padding-top   : 3rem;
            padding-bottom: 2rem;
        }

        .btn-group {
            grid-area: r5;

            display        : flex;
            justify-content: center;
            align-items    : center;

            padding-top   : 1rem;
            padding-bottom: 1rem;

            button {
                width : 180px;
                height: 50px;
                
                display        : flex;
                justify-content: center;
                align-items    : center;

                border-radius: 50px;
                border : none;
                outline: none;
                
                cursor: pointer;
                
                font-family: 'Poppins';
                font-style : normal;
                font-weight: 500;
                font-size  : 18px;
                line-height: 60px;
            }

            .unclickable {
                cursor: not-allowed;

                color     : #ffffff;
                background: #606060; // 6377EE

                &:hover {
                    color: #848484;
                }
            }

            
            .primary {
                cursor: not-allowed;
                
                color     : #ffffff;
                background: #606060; // 6377EE
                
                
                &:hover {
                    color: #113CFC; // 6377EE
                    background: #ffffff;
                }

                &.active {
                    cursor: pointer;

                    color     : #ffffff;
                    background: #113CFC; // 6377EE
                }
            }

            .secondary {
                cursor: pointer;

                color     : #848484;
                background: #ffffff;

                &:hover {
                    color: #000000;
                    background: #c7c7c7;
                }
            }
        }
    }

    .footer {
        grid-area: f;
    }
    
}
</style>
