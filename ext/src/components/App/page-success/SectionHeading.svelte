<script lang="ts">
    // ------------------------------------------------------------
    // : Import
    // ------------------------------------------------------------
    import { lang } from '../stores/store.js'

    // ------------------------------------------------------------
    // : Props
    // ------------------------------------------------------------
    const i18n = {
        'nl': {
            1: 'Browser Extensie',
            2: 'Succesvol geregistreerd',
            3: 'Bedankt voor het installeren en registeren van de browser extensie. Mocht je problemen ervaren of vragen hebben, dan kan je altijd contact met ons opnemen. Je browser kan nu normaal gebruikt worden.',
        },
        'en': {
            1: 'Browser Extension',
            2: 'Successfully registered',
            3: 'Thank you for installing and registering the browser extension. If you experience problems or have questions, you can always contact us. You can now use your browser normally.',
        },
    }
</script>

<section class='section'>
    <h2>{i18n[$lang][1]}</h2>
    <h1>{i18n[$lang][2]}</h1>
    <div>{i18n[$lang][3]}</div>
</section>


<style type="text/scss">
    .section {
        width : 100%;

        display: grid;
        grid-template-rows   : auto auto auto auto;
        grid-template-columns: 1fr;
        gap: 0.5rem;

        // Browser Extensie
        h2 {
            width : 100%;

            display        : flex;
            justify-content: center;
            align-items    : center;

            color: #193498;
        }

        // Succesvol geregistreerd
        h1 {
            width : 100%;

            display        : flex;
            justify-content: center;
            align-items    : center;
        }

        // Long text
        div:nth-of-type(1) {
            display        : flex;
            justify-content: center;
            align-items    : center;

            color: #000000;
            font-family: 'Roboto';
            font-style : normal;
            font-weight: 400;
            font-size  : 17px;
            text-align : center;   
        }
    }
</style>