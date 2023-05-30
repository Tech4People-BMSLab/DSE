<script lang="ts">
    // ------------------------------------------------------------
    // Imports
    // ------------------------------------------------------------
    import { page }         from '../stores/store'
    import { lang }         from '../stores/store'

    // ------------------------------------------------------------
    // Components
    // ------------------------------------------------------------
    import Container from '../lib/Container.svelte'

    import SectionHeading from './SectionHeading.svelte'
    import SectionTeam    from './SectionTeam.svelte'
    import SectionConsent from './SectionConsent.svelte'

    import Button from '../lib/Button.svelte'

    // ------------------------------------------------------------
    // : Props
    // ------------------------------------------------------------
    const i18n = {
        'nl': {
            1: 'Akkoord',
            2: 'Annuleren',
        },
        'en': {
            1: 'Agree',
            2: 'Cancel',
        },
    }

    // ------------------------------------------------------------
    // : Functions
    // ------------------------------------------------------------
    function on_accept() {
        page.set('register')
    }
    
    function on_reject() {
        page.set('cancelled')
    }
</script>

<main>
    <Container>
        <div>
            <SectionHeading/>
            <SectionTeam/>
            <SectionConsent/>

            <div class='btn-group'>
                <Button primary={true} onclick={on_accept}>{i18n[$lang][1]}</Button>
                <Button secondary={true} onclick={on_reject}>{i18n[$lang][2]}</Button>
            </div>
        </div>
    </Container>
</main>

<style type="text/scss">
    main {
        width : 100%;
        height: 100%;

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
