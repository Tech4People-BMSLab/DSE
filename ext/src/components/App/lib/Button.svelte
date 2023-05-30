<script lang="ts">
    // ------------------------------------------------------------
    // Props
    // ------------------------------------------------------------
    export let primary   = false
    export let secondary = false
    export let onclick   = () => {}
</script>

<button 
    on:click={onclick}
    class:btn-primary   = {primary}
    class:btn-secondary = {secondary}>
    <slot></slot>
</button>

<style type="text/scss">
    // Can't use .btn because of tailwind
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

    .btn-primary {
        color     : #ffffff;
        background: #113CFC; // 6377EE

        &:hover {
            color: #113CFC; // 6377EE
            background: #ffffff;
        }
    }

    .btn-secondary {
        color     : #848484;
        background: #ffffff;

        &:hover {
            color: #000000;
            background: #c7c7c7;
        }
    }
    



</style>