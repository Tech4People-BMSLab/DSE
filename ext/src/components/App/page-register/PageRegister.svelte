<script lang="ts">
    // ------------------------------------------------------------
    // Imports
    // ------------------------------------------------------------
    import browser     from 'webextension-polyfill'

    import { DateTime } from 'luxon'

    import { storage }       from '../utils/util'
    import { page }          from '../stores/store.js'
    import { error_message } from '../stores/store.js';

    import { has }         from '../utils/util'
    import { log }         from '../utils/util'
    import { report }      from '../utils/util'

    import { any }      from '../utils/util'
    import { is_empty } from '../utils/util'


    import { lang } from '../stores/store'
    // ------------------------------------------------------------
    // Components
    // ------------------------------------------------------------
    import Container from '../lib/Container.svelte'
    import Button    from '../lib/Button.svelte'
    
    import SectionHeading  from './SectionHeading.svelte'
    import SectionQuestion from './SectionQuestion.svelte'

    

    // ------------------------------------------------------------
    // Props
    // ------------------------------------------------------------
    const i18n = {
        'nl': {
            1: 'Bedankt voor het invullen van uw demografische gegevens. Klik op afronden om de installatie van de browser extensie te voltooien.',
            2: 'Afronden',
            3: 'Annuleren',
            4: 'Niet alle velden zijn ingevuld',
            5: 'Er is iets misgegaan met het versturen van de gegevens'
        },
        'en': {
            1: 'Thank you for filling in your demographic data. Click on finish to complete the installation of the browser extension.',
            2: 'Finish',
            3: 'Cancel',
            4: 'Not all fields are filled in',
            5: 'Something went wrong while sending the data'
        }
    }
    let questions_i18n = [ // This includes i18n (internationalization) strings
        { // 1
            'nl': {
                id        : 'q1',
                name      : 'resident',

                type      : 'dropdown',
                question  : 'Ben u woonachtend in Nederland?',
                answer    : '',
                options   : [
                    {label: 'Ja' , value: 'ja'},
                    {label: 'Nee', value: 'nee'}
                ],
            },
            'en': {
                id        : 'q1',
                name      : 'resident',

                type      : 'dropdown',
                question  : 'Are you residing in the Netherlands?',
                answer    : '',
                options   : [
                    {label: 'Yes' , value: 'ja'},
                    {label: 'No'  , value: 'nee'}
                ],
            }
        },
        { // 2
            'nl': {
                id        : 'q2',
                name      : 'sex',
                type      : 'dropdown',
                question  : 'Wat is uw geslacht?',
                answer    : '',
                options   : [
                    {label: 'Mannelijk'         , value:'mannelijk'},
                    {label: 'Vrouwelijk'        , value:'vrouwelijk'},
                    {label: 'Anders'            , value:'anders'},
                    {label: 'Zeg ik liever niet', value:'unselected'},
                ],
            },
            'en': {
                id        : 'q2',
                name      : 'sex',
                type      : 'dropdown',
                question  : 'What is your sex?',
                answer    : '',
                options   : [
                    {label: 'Male'              , value: 'mannelijk'},
                    {label: 'Female'            , value: 'vrouwelijk'},
                    {label: 'Other'             , value: 'anders'},
                    {label: 'Rather not say'    , value: 'unselected'}
                ]
            }
        },
        { // 3
            'nl': {
                id        : 'q3',
                name      : 'age',
                type      : 'dropdown',
                question  : 'Wat is uw leeftijd?',
                answer    : '',
                options   : [
                    {label: '16-24' , value: '16-24'},
                    {label: '25-34' , value: '25-34'},
                    {label: '35-44' , value: '35-44'},
                    {label: '45-54' , value: '45-54'},
                    {label: '55-64' , value: '55-64'},
                    {label: '65-74' , value: '65-74'},
                    {label: '75+'   , value: '75+'},
                    {label: 'Zeg ik liever niet', value:'unselected'},
                ],
            },
            'en': {
                id        : 'q3',
                name      : 'age',
                type      : 'dropdown',
                question  : 'What is your age?',
                answer    : '',
                options   : [
                    {label: '16-24' , value: '16-24'},
                    {label: '25-34' , value: '25-34'},
                    {label: '35-44' , value: '35-44'},
                    {label: '45-54' , value: '45-54'},
                    {label: '55-64' , value: '55-64'},
                    {label: '65-74' , value: '65-74'},
                    {label: '75+'   , value: '75+'},
                    {label: 'Rather not say', value:'unselected'},
                ],
            }
        },
        { // 4
            'nl': {
                id        : 'q4',
                name      : 'postcode',
                type      : 'postcode',
                answer    : '',
                question  : 'Wat is uw postcode? (alleen cijfers)',
            },
            'en': {
                id        : 'q4',
                name      : 'postcode',
                type      : 'postcode',
                answer    : '',
                question  : 'What is your postcode? (only numbers)',
            }
        },
        { // 5
            'nl': {
                id        : 'q5',
                name      : 'education',
                type      : 'dropdown',
                question  : 'Wat is uw hoogst genoten opleiding?',
                answer    : '',
                options   : [
                    {label: 'Geen opleiding'                           , value:'geen-opleiding'},
                    {label: 'Middelbare school (VMBO, HAVO, VWO)'      , value:'middelbare-school'},
                    {label: 'Middelbaar Beroeps Onderwijs (MBO)'       , value:'middelbaar-beroeps-onderwijs'},
                    {label: 'Hoger Beroeps Onderwijs (HBO)'            , value:'hoger-beroeps-onderwijs'},
                    {label: 'Wetenschappelijk Onderwijs (Universitair)', value:'wetenschappelijk-onderwijs'},
                    {label: 'Zeg ik liever niet', value:'unselected'},
                ],
            },
            'en': {
                id        : 'q5',
                name      : 'education',
                type      : 'dropdown',
                question  : 'What is your highest level of education?',
                answer    : '',
                options   : [
                    {label: 'No education'                           , value:'geen-opleiding'},
                    {label: 'High school (VMBO, HAVO, VWO)'          , value:'middelbare-school'},
                    {label: 'Middle Professional Education (MBO)'    , value:'middelbaar-beroeps-onderwijs'},
                    {label: 'Higher Professional Education (HBO)'    , value:'hoger-beroeps-onderwijs'},
                    {label: 'University'                             , value:'wetenschappelijk-onderwijs'},
                    {label: 'Rather not say'                         , value:'unselected'},
                ],
            }
        },
        { // 6
            'nl': {
                id        : 'q6',
                name      : 'income',
                type      : 'dropdown',
                question: 'Wat is uw persoonlijke jaarlijks netto inkomen?',
                answer    : '',
                options : [
                    {label: 'Minder dan 10.000 euro'   , value:'<10000'},
                    {label: '10.000 tot 20.000 euro'   , value:'10000-20000'},
                    {label: '20.001 tot 30.000 euro'   , value:'20001-30000'},
                    {label: '30.001 tot 40.000 euro'   , value:'30001-40000'},
                    {label: '40.001 tot 50.000 euro'   , value:'40001-50000'},
                    {label: '50.001 tot 100.000 euro'  , value:'50001-100000'},
                    {label: '100.000 of meer'          , value:'100000+'},
                    {label: 'Zeg ik liever niet'       , value:'unselected'},
                ],
            },
            'en': {
                id      : 'q6',
                name    : 'income',
                type    : 'dropdown',
                question: 'What is your personal annual net income?',
                answer  : '',
                options : [
                    {label: 'Less than 10.000 euro'    , value:'<10000'},
                    {label: '10.000 to 20.000 euro'    , value:'10000-20000'},
                    {label: '20.001 to 30.000 euro'    , value:'20001-30000'},
                    {label: '30.001 to 40.000 euro'    , value:'30001-40000'},
                    {label: '40.001 to 50.000 euro'    , value:'40001-50000'},
                    {label: '50.001 to 100.000 euro'   , value:'50001-100000'},
                    {label: '100.000 or more'          , value:'100000+'},
                    {label: 'Rather not say'           , value:'unselected'},
                ],
            }
        },
        { // 7
            'nl': {
                id        : 'q7',
                name      : 'political',
                type      : 'dropdown',
                question  : 'Welke politieke partij heeft uw voorkeur?',
                answer    : '',
                options   : [
                    {label: 'VVD'                   , value:'vvd'},
                    {label: 'D66'                   , value:'d66'},
                    {label: 'PVV'                   , value:'pvv'},
                    {label: 'CDA'                   , value:'cda'},
                    {label: 'SP'                    , value:'sp'},
                    {label: 'PvdA'                  , value:'pvda'},
                    {label: 'Groenlinks'            , value:'groenlinks'},
                    {label: 'FVD'                   , value:'fvd'},
                    {label: 'Partij voor de Dieren' , value:'pvdd'},
                    {label: 'ChristenUnie'          , value:'christenunie'},
                    {label: 'Volt'                  , value:'volt'},
                    {label: 'JA21'                  , value:'ja21'},
                    {label: 'SGP'                   , value:'sgp'},
                    {label: 'DENK'                  , value:'denk'},
                    {label: '50PLUS'                , value:'50plus'},
                    {label: 'BBB'                   , value:'bbb'},
                    {label: 'BIJ1'                  , value:'bij1'},
                    {label: 'Overige'               , value:'overige'},
                    {label: 'Zeg ik liever niet'    , value:'unselected'},
                ],
            },
            'en': {
                id      : 'q7',
                name    : 'political',
                type    : 'dropdown',
                question: 'Which political party do you prefer?',
                answer  : '',
                options : [
                    {label: 'VVD'                   , value:'vvd'},
                    {label: 'D66'                   , value:'d66'},
                    {label: 'PVV'                   , value:'pvv'},
                    {label: 'CDA'                   , value:'cda'},
                    {label: 'SP'                    , value:'sp'},
                    {label: 'PvdA'                  , value:'pvda'},
                    {label: 'Groenlinks'            , value:'groenlinks'},
                    {label: 'FVD'                   , value:'fvd'},
                    {label: 'Partij voor de Dieren' , value:'pvdd'},
                    {label: 'ChristenUnie'          , value:'christenunie'},
                    {label: 'Volt'                  , value:'volt'},
                    {label: 'JA21'                  , value:'ja21'},
                    {label: 'SGP'                   , value:'sgp'},
                    {label: 'DENK'                  , value:'denk'},
                    {label: '50PLUS'                , value:'50plus'},
                    {label: 'BBB'                   , value:'bbb'},
                    {label: 'BIJ1'                  , value:'bij1'},
                    {label: 'Other'                 , value:'overige'},
                    {label: 'Rather not say'        , value:'unselected'},
                ],
            }
        },
        { // 8
            'nl': {
                id        : 'q8',
                name      : 'employment',
                type      : 'dropdown',
                question  : 'Wat is uw huidige wersituatie?',
                answer    : '',
                options   : [
                    {label: 'Vast dienstverband'    , value:'vast-dienstverband'},
                    {label: 'Parttime dienstverband', value:'parttime-dienstverband'},
                    {label: 'Werkloos'              , value:'werkloos'},
                    {label: 'Zelfstandig'           , value:'zelfstandig'},
                    {label: 'Student'               , value:'student'},
                    {label: 'Gepensioneerd'         , value:'gepensioneerd'},
                    {label: 'Zeg ik liever niet'    , value:'unselected'},
                ],
            },
            'en': {
                id      : 'q8',
                name    : 'employment',
                type    : 'dropdown',
                question: 'What is your current employment situation?',
                answer  : '',
                options: [
                    {label: 'Full-time employment'  , value:'vast-dienstverband'},
                    {label: 'Part-time employment'  , value:'parttime-dienstverband'},
                    {label: 'Unemployed'            , value:'werkloos'},
                    {label: 'Self-employed'         , value:'zelfstandig'},
                    {label: 'Student'               , value:'student'},
                    {label: 'Retired'               , value:'gepensioneerd'},
                    {label: 'Rather not say'        , value:'unselected'},
                ],
            }
        },
        { // 9
            'nl': {
                id        : 'q9',
                name      : 'language',
                type      : 'dropdown-multi',
                question  : 'In welke taal voert u zoekopdrachten uit? (meerdere antwoorden mogelijk)',
                answer    : '',
                options   : [
                    {label: 'Nederlands'        , value:'nederlands'},
                    {label: 'Engels'            , value:'engels'},
                    {label: 'Duits'             , value:'duits'},
                    {label: 'Frans'             , value:'frans'},
                    {label: 'Spaans'            , value:'spaans'},
                    {label: 'Italiaans'         , value:'italiaans'},
                    {label: 'Zeg ik liever niet', value:'unselected'},
                ],
            },
            'en': {
                id        : 'q9',
                name      : 'language',
                type      : 'dropdown-multi',
                question  : 'In which language do you perform search queries? (multiple answers possible)',
                answer    : '',
                options   : [
                    {label: 'Dutch'             , value:'nederlands'},
                    {label: 'English'           , value:'engels'},
                    {label: 'German'            , value:'duits'},
                    {label: 'French'            , value:'frans'},
                    {label: 'Spanish'           , value:'spaans'},
                    {label: 'Italian'           , value:'italiaans'},
                    {label: 'Rather not say'    , value:'unselected'},
                ]
            }
        },
        { // 10
            'nl': {
                id        : 'q10',
                name      : 'social',
                type      : 'dropdown-multi',
                question  : 'Welke (social) media kanalen gebruik u voor nieuws en informatie? (meerdere antwoorden mogelijk)',
                answer    : '',
                options   : [
                    {label: 'TV'                , value:'tv'},
                    {label: 'De Krant'          , value:'de-krant'},
                    {label: 'Nieuwswebsites'    , value:'nieuwswebsites'},
                    {label: 'YouTube'           , value:'youtube'},
                    {label: 'Facebook'          , value:'facebook'},
                    {label: 'Instagram'         , value:'instagram'},
                    {label: 'WhatsApp'          , value:'whatsapp'},
                    {label: 'Linkedin'          , value:'linkedin'},
                    {label: 'Twitter'           , value:'twitter'},
                    {label: 'Telegram'          , value:'telegram'},
                    {label: 'Reddit'            , value:'reddit'},
                    {label: 'Radio'             , value:'radio'},
                    {label: 'Anders'            , value:'anders'},
                    {label: 'Zeg ik liever niet', value:'unselected'},
                ],
            },
            'en': {
                id      : 'q10',
                name    : 'social',
                type    : 'dropdown-multi',
                question: 'Which (social) media channels do you use for news and information? (multiple answers possible)',
                answer  : '',
                options : [
                    {label: 'TV'                , value:'tv'},
                    {label: 'The Newspaper'     , value:'de-krant'},
                    {label: 'News websites'     , value:'nieuwswebsites'},
                    {label: 'YouTube'           , value:'youtube'},
                    {label: 'Facebook'          , value:'facebook'},
                    {label: 'Instagram'         , value:'instagram'},
                    {label: 'WhatsApp'          , value:'whatsapp'},
                    {label: 'Linkedin'          , value:'linkedin'},
                    {label: 'Twitter'           , value:'twitter'},
                    {label: 'Telegram'          , value:'telegram'},
                    {label: 'Reddit'            , value:'reddit'},
                    {label: 'Radio'             , value:'radio'},
                    {label: 'Other'             , value:'anders'},
                    {label: 'Rather not say'    , value:'unselected'},
                ],
            }
        }
    ]
    let questions = []
    let submit    = false // This bind variable, is used to validate the input fields
    let completed = false


    function goto_success() {
        page.set('success')
    }

    function goto_consent() {
        error_message.set(null) // Reset the error popup
        page.set('consent')
    }

    /**
     * Send the registration data to BMS server.
     */
    async function on_submit() {
        // Notify the dropdown components that the form is submitted
        submit = true

        // Check if it's completed
        completed = questions.every(q => {
            if (q == undefined)      {return false}
            if (q['answer'] == '')   {return false}
            return true
        })

        // Show error if the form is not completed
        if (!completed) {
            error_message.set(i18n[$lang][4])
            return
        }

        let payload = {}

        // Add all the answers to the payload
        for (const answer of questions) {
            payload[answer['name']] = answer['answer']
        }

        try {
            browser.runtime.onMessage.addListener(async (payload) => {
                const target = payload.target
                const action = payload.action

                if (!target || target !== 'page.register') return

                switch (action) {
                    case 'success': 
                        goto_success()
                        break

                    case 'failed':
                        await     log(`No token was generated from the API`)
                        await     report(`No token was generated from the API`)
                        throw new Error(`No token was generated from the API`)
                }
            })

            await browser.runtime.sendMessage({
                from: 'page.register',
                data: payload,
            })
        } catch (error) {
            console.error(`Failed to register user`)
            console.error(`Reason: ${error}`)

            log(`Failed to register user`)
            log(`Reason: ${error}`)

            report(`Failed to register user`)
            report(`Reason: ${error}`)

            error_message.set(i18n[$lang][5])
        }
    }

    /**
     * On decline go to consent page.
    */
    async function on_decline() {
        goto_consent()
    }

    // ------------------------------------------------------------
    // : Reactive
    // ------------------------------------------------------------
    $: {
        if (!any(questions_i18n, is_empty)) { // Make sure it doesn't have any empty values or undefined
            questions = questions_i18n.map(q => {
                if (is_empty(q)) {
                    return q['nl']
                }
                if (is_empty(q[$lang])) { // Semi-bug, since this code is called often during render refreshes
                    return q
                }
                return q[$lang]
            })
        }
    }

    $: {
        // Close the message on completed form
        if (completed) {
            error_message.set('')
        }
    }
</script>


<main>
    <Container>
        <div class="grid-wrapper">
            <div class="heading">
                <SectionHeading/>
            </div>

            <div class="content">
                <SectionQuestion 
                    bind:questions = {questions}
                    bind:submit    = {submit}
                />
            </div>

            <div class='thanks'>{i18n[$lang][1]}</div>

            <div class="btn-group">
                <Button primary={true}   onclick={on_submit}>{i18n[$lang][2]}</Button>
                <Button secondary={true} onclick={on_decline}>{i18n[$lang][3]}</Button>
            </div>
        </div>
    </Container>
</main>


<style type="text/scss">
    main {
        height: 100vh;

        .grid-wrapper {
            max-height: 100vh;
            min-width : 1024px;
            
            display: grid;
            grid-template-rows   : auto 1fr auto auto;
            gap: 0px 0px;
            grid-auto-flow: row;
            grid-template-areas:
                "heading"
                "content"
                "thanks"
                "buttons";
            .heading   { grid-area: heading;}
            .content   { grid-area: content;}
            .thanks    { grid-area: thanks; }
            .btn-group { grid-area: buttons;}
        }

        .thanks {
            /* Bedankt voor het invullen van uw demografische gegevens. Klik op “Afronden” om de installatie van de browser extensie te voltooien. */
            height: 60px;

            color: #000000; // #848484
            font-family: 'Roboto';
            font-style : normal;
            font-weight: 400;
            font-size  : 17px;
            line-height: 32px;
            text-align: center;
        }

        .btn-group {
            display        : flex;
            justify-content: center;
            align-items    : center;
            flex-direction : row;
            gap: 0 0.5rem;

            padding-top   : 1rem;
            padding-bottom: 1rem;
        }
    }
</style>