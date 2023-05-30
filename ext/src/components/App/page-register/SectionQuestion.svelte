<script lang="ts">
    // ------------------------------------------------------------
    // : Components
    // ------------------------------------------------------------
    import Dropdown      from '../lib/Dropdown.svelte'
    import DropdownMulti from '../lib/DropdownMulti.svelte'
    import InputPostcode from '../lib/InputPostcode.svelte'

    import { any }      from '../utils/util'
    import { is_empty } from '../utils/util'

    // ------------------------------------------------------------
    // : Props
    // ------------------------------------------------------------
    export let questions
    export let submit    // To notify inputs when form is submitted
</script>

<div class="container">
    {#if !any(questions, is_empty)}
        {#each questions as question}
            {#if question.type === 'dropdown'}
                <div class={question.id}>
                    <Dropdown
                        title   = {question.question}
                        options = {question.options}
                        submit  = {submit} 
                        bind:answer  = {question.answer}
                    />
                </div>
                
            {:else if question.type === 'dropdown-multi'}
                <div class={question.id}>
                    <DropdownMulti
                        title   = {question.question}
                        options = {question.options}
                        submit  = {submit} 
                        bind:answer  = {question.answer}
                    />
                </div>
            {:else if question.type === 'postcode'}
                <div class={question.id}>
                    <InputPostcode 
                        title  = {question.question} 
                        submit = {submit}
                        bind:answer = {question.answer}
                    />
                </div>
            {/if}
        {/each}
    {/if}
</div>


<style type="text/scss">
    .container {  
        display: grid;
        grid-template-columns: 1fr 1fr 1fr 1fr 1fr 1fr;
        grid-template-rows   : 100px 110px 110px 110px 110px;
        gap: 0px 0px;
        grid-auto-flow: row;
        grid-template-areas:
            "q1 q1 q1 q4 q4 q4"
            "q2 q2 q2 q3 q3 q3"
            "q5 q5 q5 q6 q6 q6"
            "q7 q7 q7 q8 q8 q8"
            "q9 q9 q9 q10 q10 q10";

        .q1  { grid-area: q1;  }
        .q2  { grid-area: q2;  }
        .q3  { grid-area: q3;  }
        .q4  { grid-area: q4;  }
        .q5  { grid-area: q5;  }
        .q6  { grid-area: q6;  }
        .q7  { grid-area: q7;  }
        .q8  { grid-area: q8;  }
        .q9  { grid-area: q9;  }
        .q10 { grid-area: q10; }

        color: #000000;
        font-size: 1rem;
    }
</style>