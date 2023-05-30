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
            1: 'Toestemmingsverklaring',

            2: 'Door op akkoord te klikken geeft u aan dat:',

            3: 'U begrijpt dat u bij vragen contact kan opnemen met het onderzoeksteam via info@digitalepolarisatie.nl.',
            4: 'U begrijpt dat u de extensie ten alle tijden weer van uw computer kan verwijderen.',
            5: 'U begrijpt dat u bij vragen over de ethische uitvoering contact kan opnemen met de ethische commissie van de Universiteit Twente ethicscommittee-cis@utwente.nl.',
            6: 'U instemt met deelname aan het onderzoeksproject.',
        },
        'en': {
            1: 'Consent form',

            2: 'By clicking on agree, you indicate that:',

            3: 'You understand that you can contact the research team via info@digitalepolarisatie.nl.',
            4: 'You understand that you can remove the extension from your computer at any time.',
            5: 'You understand that you can contact the ethical committee of the University of Twente via ethicscommittee-cis@utwente.nl.',
            6: 'You agree to participate in the research project.',
        },
    }
</script>

<section class='section'>
    <div class='bar'><span>{i18n[$lang][1]}</span></div>

    <div class="content">
        <div>{i18n[$lang][2]}</div>
        <ul class='list'>
            <li>{i18n[$lang][3]}</li>
            <li>{i18n[$lang][4]}</li>
            <li>{i18n[$lang][5]}</li>
            <li>{i18n[$lang][6]}</li>
        </ul>
    </div>
</section>



<style type="text/scss">
    .section {
        width : 100%;

        .content {
            padding-top   : 1rem;
            padding-bottom: 1rem;
            padding-left  : 1rem;
            padding-right : 1rem;
            margin-left   : 1rem;
            margin-right  : 1rem;

            color: #000000;
            font-size: 1rem;

            div {
                padding-bottom: 0.5rem;

                line-height: 20px;
            }

            ul {
                list-style-type: disc;
                margin: 0;

                li {
                    padding-top   : 0.15rem;
                    padding-bottom: 0.15rem;

                    line-height: 25px;
                    
                }
            }
        }
    }
</style>