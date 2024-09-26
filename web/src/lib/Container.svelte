<script lang="ts">
    // ------------------------------------------------------------
    // : Components
    // ------------------------------------------------------------
    import PopupError from './PopupMessage.svelte'

    import { lang }        from '../stores/store'
    import { storage }     from '../utils/utils'
    
    // ------------------------------------------------------------
    // : Props
    // ------------------------------------------------------------
    let year = new Date().getFullYear()

    let lang_en = false
    let lang_nl = false

    // ------------------------------------------------------------
    // : Functions
    // ------------------------------------------------------------
    function change_lang(selected) {
        if (selected == 'en') {
            lang.set('en')
            lang_en = true
            lang_nl = false
            storage.set({'language': 'en'})
        } else if (selected == 'nl') {
            lang.set('nl')
            lang_en = false
            lang_nl = true
            storage.set({'language': 'nl'})
        }
    }

    // ------------------------------------------------------------
    // : Reactive
    // ------------------------------------------------------------
    $: {
        switch ($lang) {
            case 'en':
                lang_en = true
                lang_nl = false
                break
            case 'nl':
                lang_en = false
                lang_nl = true
                break
        }
    }
</script>

<div class="container-wrapper">
    <PopupError/>

    <div class="left">
        <img src='./images/logo/dp_logo.png' alt='DP logo'>
    </div>

    <div class="center">
        <slot></slot>
    </div>

    <div class="right">
        <img src='./images/logo/ut_logo.png' alt='UT logo'>
    </div>

    <footer class="footer">
        <span>Digitale Polarisatie {year}</span>
        <div></div>
        <div class="language">
            <span on:click={() => change_lang('nl')} class:active={lang_nl}>Nederlands</span>
            <span on:click={() => change_lang('en')} class:active={lang_en}>English (UK)</span>
        </div>
    </footer>
</div>

<style type="text/scss">
    .container-wrapper {
        width : 100%;
        height: 100%;

        display: grid;
        grid-template-columns: minmax(max-content, 1fr) 2fr minmax(max-content, 1fr);
        grid-template-rows   : 1fr 4rem;
        grid-auto-flow: row;
        grid-template-areas:
            "l c r"
            "f f f";
        .left   { grid-area: l; }
        .center { grid-area: c; }
        .right  { grid-area: r; }
        .footer { grid-area: f; }

        
        // Make sure all of them are sized accordingly to the grid specification
        .left, .right {
            width : 100%;
            height: 100%;
            overflow: hidden;
        }

        .left {
            height: 10rem;

            display        : flex;
            justify-content: flex-start;
            align-items    : flex-start;

            img {
                width : 100%;
                height: 60px;
                object-fit: contain;

                padding: 1rem;
            }
        }
        
        .right {
            height: 10rem;

            display        : flex;
            justify-content: flex-end;
            align-items    : flex-start;

            img {
                width : 100%;
                height: 60px;
                object-fit: contain;

                padding: 1rem;
            }
        }

        .center {
            padding-top: 3rem;
        }

        .footer {
            display        : flex;
            justify-content: flex-start;
            align-items    : center;

            display: grid;
            grid-template-rows   : 100%;
            grid-template-columns: 5rem 1fr 1fr 1fr 5rem;
            grid-template-areas: 
                ". l c r .";

            *:nth-of-type(1) { grid-area: l;}
            *:nth-of-type(2) { grid-area: r;}

            span {
                font-size: 0.8rem;
            }

            .language {
                display        : flex;
                justify-content: flex-end;
                align-items    : center;

                cursor: pointer;

                gap: 0.5rem;

                .active {
                    font-weight: 900;
                }
            }

        }

        // Screen size
        @media only screen and (max-width: 1400px) {
            
            display: grid;
            grid-template-rows   : 1fr 4rem;
            grid-template-columns: 3rem 1fr 3rem;
            grid-auto-flow: row;
            grid-template-areas:
                ". c ."
                ". f .";
            .left   { display: none; }
            .center { grid-area: c;  }
            .right  { display: none; }
            .footer { grid-area: f;  }

            box-sizing: border-box;
        }
    }

</style>
