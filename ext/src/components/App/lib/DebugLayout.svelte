<script lang="ts">
    export let enabled = true

    // Generate a list 10 items long
    const arr = Array.from({length: 10}, (_, i) => i)
</script>

<div class="layout-wrapper" class:enabled={enabled}>
    {#each arr as item}
        <div></div>
    {/each}
</div>

<style type="text/scss">
    .layout-wrapper {
        position: fixed;
        top : 0;
        left: 0;

        width : 100vw;
        height: 100vh;

        // Grid with equal columns and rows
        display        : grid;
        grid-template-columns: repeat(auto-fill, minmax(max-content, 1fr));
        grid-template-rows   : repeat(auto-fill, minmax(max-content, 1fr));
        grid-auto-flow: column;
    
        pointer-events: none;
    }

    .enabled {
        * {
            outline: 1px solid red; // Remove this
        }
    }
</style>