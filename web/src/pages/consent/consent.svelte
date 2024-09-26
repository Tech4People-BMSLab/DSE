<script lang="ts">
    // ------------------------------------------------------------
    // : Imports
    // ------------------------------------------------------------
    import { onMount }  from 'svelte'
    import { navigate } from 'svelte-routing'

    import Header from '@/components/header.svelte'

    import global from '@/stores/global'
    // ------------------------------------------------------------
    // : Init
    // ------------------------------------------------------------
    const base_path = import.meta.env.BASE_URL || '';
    
    function change_lang(selected: string) {
        switch (selected) {
            case 'en':
                $global.language = 'en'
                break
            case 'nl':
                $global.language = 'nl'
                break
        }
    }

    function visit(page: string) {
        switch (page) {
            case 'agree'  : navigate(`${base_path}/form`,     { replace: true }); break
            case 'decline': navigate(`${base_path}/declined`, { replace: true }); break
        }
    }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<main class="page-consent">
    <div class="header">
        <Header/>
    </div>

    <div class="content">
        <div class="center">
            <h1>Consent Form</h1>
            <h2>Digital Polarization</h2>
            <h3>BMS Ethical Committee approval number 220261</h3>

            <div class="separator">Research team</div>
            <div class="team">
                <table class="table" cellpadding="1" cellspacing="0">
                    <tr>
                        <td>Dr. Shenja van der Graaf</td>
                        <td><a href="mailto:shenja.vandergraaf@utwente.nl">shenja.vandergraaf@utwente.nl</a></td>
                    </tr>
                    <tr>
                        <td>Dr. Sikke Jansma</td>
                        <td><a href="mailto:s.r.jansma@utwente.nl">s.r.jansma@utwente.nl</a></td>
                    </tr>
                    <tr>
                        <td>Dr. Maryam Amir Heari</td>
                        <td><a href="mailto:m.amirhaeri@utwente.nl">m.amirhaeri@utwente.nl</a></td>
                    </tr>
                    <tr>
                        <td>Prof. Dr. Ing. Alexander van Deursen</td>
                        <td><a href="mailto:a.j.a.m.vandeursen@utwente.nl">a.j.a.m.vandeursen@utwente.nl</a></td>
                    </tr>
                    <tr>
                        <td>Andre Bester</td>
                        <td><a href="mailto:a.n.bester@utwente.nl">a.n.bester@utwente.nl</a> <a href="mailto:bmslab@utwente.nl">bmslab@utwente.nl</a></td>
                    </tr>
                    <tr>
                        <td>Derwin Tromp</td>
                        <td><a href="mailto:d.e.tromp@utwente.nl">d.e.tromp@utwente.nl</a></td>
                    </tr>
                    <tr>
                        <td>Kars Snijders</td>
                        <td><a href="mailto:k.j.snijders-1@student.utwente.nl">k.j.snijders-1@student.utwente.nl</a></td>
                    </tr>
                    <tr>
                        <td>BMS Lab</td>
                        <td><a href="mailto:bmslab@utwente.nl">bmslab@utwente.nl</a></td>
                    </tr>
                </table>
            </div>

            <div class="separator">Consent form</div>
            <div class="form">
                <span>By clicking on agree, you indicate that:</span>

                <ul>
                    <li>You understand that you can contact the research team via <a href="mailto:info@digitalepolarisatie.nl">info@digitalepolarisatie.nl</a>.</li>
                    <li>You understand that you can remove the extension from your computer at any time.</li>
                    <li>You understand that you can contact the ethical committee of the University of Twente via <a href="mailto:ethicscommittee-cis@utwente.nl">ethicscommittee-cis@utwente.nl</a>.</li>
                    <li>You agree to participate in the research project.</li>
                </ul>
            </div>

            <div class="btn-group">
                <button class="primary" on:click={() => visit('agree')}>Agree</button>
                <button class="secondary" on:click={() => visit('decline')}>Decline</button>
            </div>

        </div>
    </div>

    <div class="footer">
        <div class="left">
            <span>Digitale Polarisatie {$global.year}</span>
        </div>

        <div class="right">
            <div class="language">
                <span on:click={() => change_lang('nl')} class:active={$global.language == 'nl'}>Nederlands</span>
                <span on:click={() => change_lang('en')} class:active={$global.language == 'en'}>English (UK)</span>
            </div>
        </div>
    </div>
</main>


<style lang="scss" scoped>

.page-consent {
    position: absolute;
    top : 0;
    left: 0;

    width : 100%;
    height: 100%;

    display: grid;
    grid-template:
        "h" auto
        "c" 1fr
        "f" auto / 1fr;

    font-family: 'Poppins';

    .header {
        grid-area: h;
    }

    .content {
        grid-area: c;

        display: grid;
        grid-template:
            ". c ." 1fr / 200px 1fr 200px;

        font-family: 'Poppins';

        .center {
            grid-area: c;

            display: grid;
            grid-template:
                "h1" auto
                "h2" auto
                "h3" auto
                "s1" auto
                "c1" auto
                "s2" auto
                "c2" auto
                "b1" auto
                "." 1fr / 1fr;

            h1 {
                display        : flex;
                justify-content: center;
                align-items    : center;

                font-size: 40px;
            }

            h2 {
                display        : flex;
                justify-content: center;
                align-items    : center;

                font-size: 17px;
            }

            h3 {
                display        : flex;
                justify-content: center;
                align-items    : center;

                font-size  : 17px;
                font-weight: 400;
            }

            .separator {
                height: 45px;

                display        : flex;
                justify-content: flex-start;
                align-items    : center;

                color: white;
                font-size: 20px;
                font-weight: 600;

                border-radius: 50px;
                padding-left : 1rem;
                box-sizing   : border-box;

                background-color: #193498;
            }

            .team {
                display: grid;
                grid-template:
                    ". t ." 1fr / 50px 1fr 50px;

                padding-top   : 1rem;
                padding-bottom: 1rem;

                .table {
                    grid-area: t;

                    color: #848484;
                    font-family: 'Segoe UI';

                    tr {
                        td {
                            padding-top   : 0.50rem;
                            padding-bottom: 0.50rem;
                        }
                    }
                }
            }

            .form {
                display: grid;
                display        : flex;
                justify-content: flex-start;
                align-items    : flex-start;
                flex-direction : column;

                padding-top   : 1rem;
                padding-bottom: 1rem;
                padding-left : 50px;
                padding-right: 50px;

                font-family: 'Segoe UI';

                span {
                    font-size: 17px;
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

            .btn-group {
                display        : flex;
                justify-content: space-around;
                align-items    : flex-start;
                flex-direction : row;

                padding-top   : 2rem;

                button {
                    width : 180px;
                    height: 50px;
                    
                    display        : flex;
                    justify-content: center;
                    align-items    : center;

                    border-radius: 50px;
                    border : none;
                    outline: none;
                    
                    cursor: pointer;
                    
                    font-family: 'Poppins';
                    font-style : normal;
                    font-weight: 500;
                    font-size  : 18px;
                    line-height: 60px;
                }

                .primary {
                    color     : #ffffff;
                    background: #113CFC; // 6377EE

                    &:hover {
                        color: #113CFC; // 6377EE
                        background: #ffffff;
                    }
                }

                .secondary {
                    color     : #848484;
                    background: #ffffff;

                    &:hover {
                        color: #000000;
                        background: #c7c7c7;
                    }
                }
            }
        }
    }

    .footer {
        grid-area: f;

        width : 100%;
        height: 100px;

        display: grid;
        grid-template:
            "l c r" 1fr / 1fr;

        font-family: 'Segoe UI';

        .left {
            grid-area: l;

            display        : flex;
            justify-content: flex-start;
            align-items    : flex-end;

            padding-left  : 1rem;
            padding-bottom: 1rem;
        }
        
        .right {
            grid-area: r;

            display        : flex;
            justify-content: flex-end;
            align-items    : flex-end;

            padding-right : 1rem;
            padding-bottom: 1rem;

            .active {
                color: black;
                font-weight: 800;
            }
        }
    }
}

</style>
