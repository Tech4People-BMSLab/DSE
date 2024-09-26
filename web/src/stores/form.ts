import { writable } from 'svelte/store'

export default writable({
    questions: [
        {
            id        : 'q1',
            name      : 'resident',
            type      : 'dropdown',
            question  : {
                'nl': 'Bent u woonachtig in Nederland?*',
                'en': 'Are you residing in the Netherlands?*'
            },
            answer    : {selected: {value: '', label: ''}, valid: 'invalid'},
            options   : [
                {label: 'Ja' , value: 'ja'},
                {label: 'Nee', value: 'nee'}
            ],
        },
        {
            id        : 'q2',
            name      : 'sex',
            type      : 'dropdown',
            question  : {
                'nl': 'Wat is uw geslacht?*',
                'en': 'What is your sex?*'
            },
            answer    : {selected: {value: '', label: ''}, valid: 'invalid'},
            options   : [
                {label: 'Mannelijk'         , value:'mannelijk'},
                {label: 'Vrouwelijk'        , value:'vrouwelijk'},
                {label: 'Anders'            , value:'anders'},
                {label: 'Zeg ik liever niet', value:'unselected'},
            ],
        },
        {
            id        : 'q3',
            name      : 'age',
            type      : 'dropdown',
            question  : {
                'nl': 'Wat is uw leeftijd?',
                'en': 'What is your age?'
            },
            answer    : {selected: {value: '', label: ''}, valid: 'invalid'},
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
        {
            id        : 'q4',
            name      : 'postcode',
            type      : 'postcode',
            question  : {
                'nl': 'Wat is uw postcode?* (alleen cijfers)',
                'en': 'What is your postcode?* (only numbers)'
            },
            answer    : {postcode: '', valid: 'invalid'},
        },
        {
            id        : 'q5',
            name      : 'education',
            type      : 'dropdown',
            question  : {
                'nl': 'Wat is uw hoogst genoten opleiding?',
                'en': 'What is your highest level of education?'
            },
            answer    : {selected: {value: '', label: ''}, valid: 'invalid'},
            options   : [
                {label: 'Geen opleiding'                           , value:'geen-opleiding'},
                {label: 'Middelbare school (VMBO, HAVO, VWO)'      , value:'middelbare-school'},
                {label: 'Middelbaar Beroeps Onderwijs (MBO)'       , value:'middelbaar-beroeps-onderwijs'},
                {label: 'Hoger Beroeps Onderwijs (HBO)'            , value:'hoger-beroeps-onderwijs'},
                {label: 'Wetenschappelijk Onderwijs (Universitair)', value:'wetenschappelijk-onderwijs'},
                {label: 'Zeg ik liever niet', value:'unselected'},
            ],
        },
        {
            id        : 'q6',
            name      : 'income',
            type      : 'dropdown',
            question  : {
                'nl': 'Wat is uw persoonlijke netto jaarinkomen?',
                'en': 'What is your personal annual net income?'
            },
            answer    : {selected: {value: '', label: ''}, valid: 'invalid'},
            options : [
                {label: 'Minder dan 10.000 euro'   , value:'<10000'},
                {label: '10.000 tot 20.000 euro'   , value:'10000-20000'},
                {label: '20.001 tot 30.000 euro'   , value:'20001-30000'},
                {label: '30.001 tot 40.000 euro'   , value:'30001-40000'},
                {label: '40.001 tot 50.000 euro'   , value:'40001-50000'},
                {label: '50.001 tot 100.000 euro'  , value:'50001-100000'},
                {label: '100.001 of meer'          , value:'100000+'},
                {label: 'Zeg ik liever niet'       , value:'unselected'},
            ],
        },
        {
            id        : 'q7',
            name      : 'political',
            type      : 'dropdown',
            question  : {
                'nl': 'Welke politieke partij heeft uw voorkeur?',
                'en': 'Which political party do you prefer?'
            },
            answer    : {selected: {value: '', label: ''}, valid: 'invalid'},
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
        {
            id        : 'q8',
            name      : 'employment',
            type      : 'dropdown',
            question  : {
                'nl': 'Wat is uw huidige werksituatie?',
                'en': 'What is your current employment situation?'
            },
            answer    : {selected: {value: '', label: ''}, valid: 'invalid'},
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
        {
            id        : 'q9',
            name      : 'language',
            type      : 'dropdown-multi',
            question  : {
                'nl': 'In welke taal voert u zoekopdrachten uit? (meerdere antwoorden mogelijk)',
                'en': 'In which language do you perform search queries? (multiple answers possible)'
            },
            answer    : {selected: [], label: '', valid: 'invalid'},
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
        {
            id        : 'q10',
            name      : 'social',
            type      : 'dropdown-multi',
            question  : {
                'nl': 'Welke (social) media kanalen gebruikt u voor nieuws en informatie? (meerdere antwoorden mogelijk)',
                'en': 'Which (social) media channels do you use for news and information? (multiple answers possible)'
            },
            answer    : {selected: [], label: '', valid: 'invalid'},
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
        {
            id        : 'q11',
            name      : 'browser',
            type      : 'dropdown-multi',
            question  : {
                'nl': 'Op welke browser voert u voornamelijk zoekopdrachten uit? (meerdere antwoorden mogelijk)',
                'en': 'On which browser do you mainly perform search queries? (Multiple answers possible)'
            },
            answer    : {selected: [], label: '', valid: 'invalid'},
            options   : [
                {label: 'Firefox'           , value:'firefox'},
                {label: 'Chrome'            , value:'chrome'},
                {label: 'Microsoft Edge'    , value:'microsoft-edge'},
                {label: 'Opera'             , value:'opera'},
                {label: 'Safari'            , value:'safari'},
                {label: 'Brave'             , value:'brave'},
                {label: 'Zeg ik liever niet', value:'unselected'},
            ],
        },
        {
            id        : 'q12',
            name      : 'search_engine',
            type      : 'dropdown-multi',
            question  : {
                'nl': 'Op welke zoekmachine voert u voornamelijk zoekopdrachten uit? (meerdere antwoorden mogelijk)',
                'en': 'On which search engine do you mainly perform search queries? (Multiple answers possible)'
            },
            answer    : {selected: [], label: '', valid: 'invalid'},
            options   : [
                {label: 'Google'            , value:'google'},
                {label: 'DuckDuckGo'        , value:'duckduckgo'},
                {label: 'Bing'              , value:'bing'},
                {label: 'Yahoo'             , value:'yahoo'},
                {label: 'StartPage'         , value:'startpage'},
                {label: 'Ecosia'            , value:'ecosia'},
                {label: 'Anders'            , value:'anders'},
                {label: 'Zeg ik liever niet', value:'unselected'},
            ],
        },
    ],
    completed: false,
    submitted: false,
})
