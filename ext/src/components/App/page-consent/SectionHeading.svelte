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
            2: 'Toestemmingsformulier',
            3: 'Digitale Polarisatie',
            4: 'BMS Etische Commissie goedkeuringsnummer 220261',
        },
        'en': {
            1: 'Browser Extension',
            2: 'Consent Form',
            3: 'Digital Polarization',
            4: 'BMS Ethical Committee approval number 220261',
        },
    }
</script>

<section class='section'>
    <h2>{i18n[$lang][1]}</h2>
    <h1>{i18n[$lang][2]}</h1>
    <div>{i18n[$lang][3]}</div>
    <div>{i18n[$lang][4]}</div>
</section>


<style type="text/scss">
    .section {
        width : 100%;

        //@todo: Meeting 2021-03-01: I changed the values from "auto auto auto auto" to those, and it looks worse, since the form is bigger vertically (causes scrolling)
        display: grid;
        grid-template-rows   : auto auto auto auto;
        grid-template-columns: 1fr;
        gap: 0.5rem;

        padding-bottom: 1rem;

        // Browser Extensie
        h2 { 
            width : 100%;

            display        : flex;
            justify-content: center;
            align-items    : center;

            color: #193498;
        }

        // Toestemmingsformulier
        h1 {
            width: 100%;

            display        : flex;
            justify-content: center;
            align-items    : center;
        }

        // Digitale Polarisatie
        div:nth-of-type(1) {
            width : 100%;
            height: 40px;

            display        : flex;
            justify-content: center;
            align-items    : center;

            color: #000000;
            font-family: 'Roboto';
            font-style : normal;
            font-weight: 700;
            font-size  : 17px;
            text-align : center;
        }
        
        // BMS Etische Commissie goedkeuringsnummer 220261
        div:nth-of-type(2) {
            width : 100%;
            height: 2rem;

            display        : flex;
            justify-content: center;
            align-items    : center;

            color: #000000;
            font-family: 'Roboto';
            font-style : normal;
            font-weight: 400;
            font-size  : 17px;
            text-align : center;   

            word-spacing: 2;
        }
    }
</style>