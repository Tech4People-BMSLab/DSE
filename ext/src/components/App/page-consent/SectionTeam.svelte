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
            1: 'Onderzoeksteam',
        },
        'en': {
            1: 'Research team',
        },
    }
    const contacts = [
        {name: 'Dr. Shenja van der Graaf'              , emails: ['shenja.vandergraaf@utwente.nl']},
        {name: 'Dr. Sikke Jansma'                      , emails: ['s.r.jansma@utwente.nl']},
        {name: 'Dr. Maryam Amir Heari'                 , emails: ['m.amirhaeri@utwente.nl']},
        {name: 'Prof. Dr. Ing. Alexander van Deursen'  , emails: ['a.j.a.m.vandeursen@utwente.nl']},
        {name: 'Andre Bester'                          , emails: ['a.n.bester@utwente.nl', 'bmslab@utwente.nl']},
        {name: 'Derwin Tromp'                          , emails: ['d.e.tromp@utwente.nl']},
        {name: 'Kars Snijders'                         , emails: ['k.j.snijders-1@student.utwente.nl']},
        {name: 'BMS Lab'                               , emails: ['bmslab@utwente.nl']},
    ]
</script>

<section class='section'>
    <div class='bar'><span>{i18n[$lang][1]}</span></div>

    <div class="content">
        <table class='table' cellpadding=1 cellspacing=0>
            {#each contacts as contact}
                <tr>
                    <td>{contact.name}</td>
                    <td>
                        {#each contact.emails as email}
                            <a href="mailto:{email}">{email}</a>
                            <span></span>
                        {/each}
                    </td>
                </tr>
            {/each}
        </table>
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
    
            color: #848484;
            font-size: 1rem;

            .table {
                width : 100%;

                color: #848484;

                tr {
                    td {
                        padding-top   : 0.20rem;
                        padding-bottom: 0.20rem;
                    }
                }
            }
        }
    }
</style>