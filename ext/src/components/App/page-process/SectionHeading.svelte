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
            1: 'NIET AFSLUITEN',
            2: 'Zoekopdrachten binnenhalen',
            3: 'Op het moment is de browser extensie bezig met het uitvoeren van zoekopdrachten. Het is belangrijk om dit venster niet af te sluiten Het venster zal automatisch gesloten worden wanneer het proces voltooid is.',
        },
        'en': {
            1: 'DO NOT CLOSE',
            2: 'Fetching search queries',
            3: 'At the moment the browser extension is busy with executing search queries. It is important not to close this window. The window will automatically close when the process is completed.',
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

        // Zoekopdracht binnenhalen
        h1 {
            width : 100%;

            display        : flex;
            justify-content: center;
            align-items    : center;
        }

        // Long text (nothing to do with it yet)
        // p { }
    }
</style>