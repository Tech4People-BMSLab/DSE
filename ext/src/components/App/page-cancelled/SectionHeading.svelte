<script lang="ts">

    // ------------------------------------------------------------
    // : Imports
    // ------------------------------------------------------------
    import { lang } from '../stores/store.js'

    // ------------------------------------------------------------
    // : Props
    // ------------------------------------------------------------
    const i18n = {
        'nl': {
            1: 'Browser Extensie',
            2: 'Aanmeldproces geannuleerd',
            3: 'Het aanmeldproces is geannuleerd en de browser extensie zal hierdoor niet werken. Mocht je problemen ervaren of vragen hebben, dan kan je altijd contact met ons opnemen.  Wij proberen de browser extensie dagelijks te verbeteren voor een goede gebruikerservaring. Daarom wordt ook jou feedback gewaardeerd.',
        },
        'en': {
            1: 'Browser Extension',
            2: 'Registration process cancelled',
            3: 'The registration process has been cancelled and the browser extension will not work as a result. If you experience problems or have questions, you can always contact us. We try to improve the browser extension daily for a good user experience. That is why your feedback is also appreciated.',
        },
    }
</script>

<section class='section'>
    <h2>{i18n[$lang][1]}</h2>
    <h1>{i18n[$lang][2]}</h1>
    <p>{i18n[$lang][3]}</p>
</section>


<style type="text/scss">
    .section {
        width: 100%;

        display: grid;
        grid-template-rows   : auto auto auto auto;
        grid-template-columns: 1fr;
        gap: 0.5rem;

        // Browser Extensie
        h2 {
            width: 100%;

            display        : flex;
            justify-content: center;
            align-items    : center;

            color: #193498;
        }

        // Gegevens
        h1 {
            width: 100%;

            display        : flex;
            justify-content: center;
            align-items    : center;
        }

        // Long text
        p { 
            display        : flex;
            justify-content: flex-start;
            align-items    : flex-start;

            text-align : left;
        }
    }
</style>