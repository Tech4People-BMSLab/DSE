<script lang="ts">
    // ------------------------------------------------------------
    // Components
    // ------------------------------------------------------------
    import Container      from '../lib/Container.svelte'
    import SectionHeading from './SectionHeading.svelte'
</script>


<main>
    <Container>
        <div>
            <SectionHeading/>
            <div class="wrapper">
                <img src='./images/process.svg' alt='Process'/>
            </div>
        </div>
    </Container>
</main>


<style type="text/scss">
    main {
        height: 100vh; // TODO: Maybe change this
        width : 100%;

        .wrapper {
           display        : flex;
           justify-content: center;
           align-items    : center;

           padding-top   : 3rem;
        }
    }
</style>