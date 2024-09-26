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
    @import url('https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap');

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