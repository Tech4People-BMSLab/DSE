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
            2: 'Gegevens',
            3: 'Om de data te kunnen vergelijken vragen wij een aantal demografische gegevens in te vullen. Deze gegevens kunnen niet herleid worden. Mocht u bepaalde gegevens niet willen invullen dan kunt u dit selecteren.',
        },
        'en': {
            1: 'Browser Extension',
            2: 'Details',
            3: 'To compare the data we ask you to fill in some demographic details. These details cannot be traced back to you. If you do not want to fill in certain details you can select this.',
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

        // Gegevens
        h1 {
            width: 100%;

            display        : flex;
            justify-content: center;
            align-items    : center;
        }

        // Long text
        div:nth-of-type(1) {
            width : 100%;
            height: 60px;
            
            display        : flex;
            justify-content: center;
            align-items    : center;

            color: #000000;
            font-family: 'Roboto';
            font-style : normal;
            font-weight: 400;
            font-size  : 17px;
            text-align : left;
        }
    }
</style>