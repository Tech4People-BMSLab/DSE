<script lang="ts">
    import global from '@/stores/global'

    function change_lang(selected: string) {
        switch (selected) {
            case 'en':
                $global.language = 'en'
                break
            case 'nl':
                $global.language = 'nl'
                break 
        }
    }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<footer class="footer">
    <div class="left">
        <span>Digitale Polarisatie {$global.year}</span>
    </div>

    <div class="right">
        <div class="language">
            <span on:click={() => change_lang('nl')} class:active={$global.language == 'nl'}>Nederlands</span>
            <span on:click={() => change_lang('en')} class:active={$global.language == 'en'}>English (UK)</span>
        </div>
    </div>
</footer>


<style lang="scss" scoped>

    .footer {
        grid-area: f;

        width : 100vw;
        height: 100px;

        display: grid;
        grid-template:
            "l c r" 1fr / 1fr;

        font-family: 'Segoe UI';

        .left {
            grid-area: l;

            display        : flex;
            justify-content: flex-start;
            align-items    : flex-end;

            padding-left  : 1rem;
            padding-bottom: 1rem;
        }
        
        .right {
            grid-area: r;

            display        : flex;
            justify-content: flex-end;
            align-items    : flex-end;

            padding-right : 1rem;
            padding-bottom: 1rem;

            .active {
                color: black;
                font-weight: 800;
            }
        }
    }
</style>
