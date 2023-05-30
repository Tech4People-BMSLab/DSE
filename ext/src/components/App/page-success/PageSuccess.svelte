<script lang="ts">
    // ------------------------------------------------------------
    // Imports
    // ------------------------------------------------------------
    import { onMount } from 'svelte'
	import { draw    } from 'svelte/transition'

    import { lang } from '../stores/store'

    // ------------------------------------------------------------
    // Components
    // ------------------------------------------------------------
    import Container      from '../lib/Container.svelte'
    import SectionHeading from './SectionHeading.svelte'
    import Button         from '../lib/Button.svelte'

    // ------------------------------------------------------------
    // : Props
    // ------------------------------------------------------------
    const i18n = {
        'nl': {
            1: 'Sluiten'
        },
        'en': {
            1: 'Close'
        },
    }

    let animate = false

    // ------------------------------------------------------------
    // : Functions
    // ------------------------------------------------------------
    function on_close() {
        window.close()
    }

    onMount(() => {
        animate = true  
    })
</script>

<main>
    <Container>
        <div>
            <SectionHeading/>
            <div class="wrapper">
                {#if animate}
                    <svg
                        xmlns:dc="http://purl.org/dc/elements/1.1/"
                        xmlns:cc="http://creativecommons.org/ns#"
                        xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
                        xmlns:svg="http://www.w3.org/2000/svg"
                        xmlns="http://www.w3.org/2000/svg"
                        id="svg5046"
                        version="1.1"
                        viewBox="0 0 50 50"
                        height="50mm"
                        width="50mm">
                        <defs id="defs5040" />
                        <g
                            style="display:inline"
                            transform="translate(0,-247)"
                            id="draw-line">
                            <path in:draw="{{duration: 800}}"
                                d="m 18.963249,271.54142 5.11928,5.88254 16.621216,-13.3733 c -4.624996,-7.33644 -9.696966,-10.5123 -18.398843,-9.418 -9.407725,1.18306 -15.8732057,10.07183 -14.7734199,19.63013 1.0997857,9.55831 9.3987219,16.23758 19.2969299,15.20827 11.395441,-1.18501 19.680898,-12.82477 13.875333,-25.4204"
                                style="opacity:1;fill:none;stroke:#848484;stroke-width:4.84156466;stroke-linecap:round;stroke-linejoin:round;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers"
                                id="path5046" 
                            />
                        </g>
                    </svg>
                {/if}
            </div>
            <div class='btn-group'>
                <Button primary={true} onclick={on_close}>{i18n[$lang][1]}</Button>
            </div>
        </div>
    </Container>
</main>

<style type="text/scss">
    main {
        width : 100%;
        height: 100vh;

        .wrapper {
            display        : flex;
            justify-content: center;
            align-items    : center;

            padding-top   : 5rem;
            padding-bottom: 5rem;
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
