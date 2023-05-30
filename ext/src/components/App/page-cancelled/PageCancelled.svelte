<script lang="ts">
    // ------------------------------------------------------------
    // Components
    // ------------------------------------------------------------
    import Container      from '../lib/Container.svelte'
    import SectionHeading from './SectionHeading.svelte'
    import SectionForm    from './SectionForm.svelte'
</script>


<main>
    <Container>
        <div>
            <SectionHeading/>
            <SectionForm/>
        </div>
    </Container>
</main>


<style type="text/scss">
    main {
        height: 100vh; // TODO: Maybe change this
        width : 100%;
    }
</style>