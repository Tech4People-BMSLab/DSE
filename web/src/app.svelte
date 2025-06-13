<script lang="ts">
    // ------------------------------------------------------------
    // : Imports
    // ------------------------------------------------------------
    import { onMount }                 from 'svelte'
    import { Router, Route, navigate } from 'svelte-routing'

    import PageConsent  from '@/pages/consent/consent.svelte'
    import PageForm     from '@/pages/form/form.svelte'
    import PageComplete from '@/pages/complete/complete.svelte'

    import { is_debug } from '@/utils/utils'
    // ------------------------------------------------------------
    // : Helpers
    // ------------------------------------------------------------
    const base_path = import.meta.env.BASE_URL || ''

    // ------------------------------------------------------------
    // : Init
    // ------------------------------------------------------------
    onMount(() => {
        let url = new URL(window.location.href)

        switch (true) {
            case is_debug():
                navigate(`${base_path}/consent`)
                break

            case url.pathname === `${base_path}/`:
                navigate(`${base_path}/consent`)
                break
        }
    })
</script>

<main class="page-index">
    <div class="left"></div>
    <div class="content">
        <Router>
            <Route path={`${base_path}/`}         component={PageConsent} />
            <Route path={`${base_path}/consent`}  component={PageConsent} />
            <Route path={`${base_path}/form`}     component={PageForm} />
            <Route path={`${base_path}/complete`} component={PageComplete} />
        </Router>
    </div>
    <div class="right"></div>
</main>

<style lang="scss" scoped>
@import url('https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap');
@import url('https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap');
@import url('https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100;0,300;0,400;0,500;0,700;0,900;1,100;1,300;1,400;1,500;1,700;1,900&display=swap');

.page-index {
    position: absolute;
    top : 0;
    left: 0;

    width : 100vw;
    height: 100vh;

    min-width : 1024px;
    min-height: 768px;

    display: grid;
    grid-template:
        "l c r" 1fr / minmax(200px, 1fr) auto minmax(200px, 1fr);

    .content {
        grid-area: c;
    }
}   
</style>
