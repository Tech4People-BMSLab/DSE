<script lang="ts">
    // ------------------------------------------------------------
    // : Imports
    // ------------------------------------------------------------
    import { onMount } from 'svelte'
    import { draw }    from 'svelte/transition'

    import Header from '@/components/header.svelte'
    import Footer from '@/components/footer.svelte'

    import storage from '@/services/storage'
    // ------------------------------------------------------------
    // : Init
    // ------------------------------------------------------------
    const base_path = import.meta.env.BASE_URL || ''
    
    let raw : any
    let form: any
    let data: any
    let animate = false


    function convert(questions: { [key: string]: any }): any {
        const result: any = {
            age          : '',
            sex          : '',
            income       : '',
            social       : {},
            browser      : {},
            language     : {},
            postcode     : { value: '' },
            resident     : '',
            education    : '',
            political    : '',
            employment   : '',
            search_engine: {}
        }

        Object.keys(questions).forEach((key) => {
            const question = questions[key]
            if (!question || !question.name) return
            
            const { name, answer = {}, options = [] } = question
            const selected = answer.selected || {}
            
            switch (name) {
                case 'age':
                case 'sex':
                case 'income':
                case 'education':
                case 'political':
                case 'employment':
                case 'resident':
                    result[name] = selected.value || ''
                    break
                case 'postcode':
                    result.postcode.value = answer.postcode || ''
                    break
                case 'social':
                case 'browser':
                case 'language':
                case 'search_engine':
                    result[name] = options.reduce((acc: any, option: any) => {
                        acc[option.value] = selected.some((sel: any) => sel.value === option.value) || false
                        return acc
                    }, {})
                    break
            }
        })

        return result
    }

    onMount(async () => {
        animate = true

        raw  = await storage.get('form')
        form = convert(raw.questions)

        data = {form, raw}
        data = JSON.stringify(data, null, 0)
    })

</script>


<div class="page-complete">

    <div class="header">
        <Header/>
    </div>

    <div class="content">
        <div class="r1"><h1>Thank for your cooperation/help</h1></div>
        <div class="r2"><p>Thank you for installing and registerin the browser extension. If you experience problems or have questions, you can always contact us. You can now use your browser normally.</p></div>
        <div class="r3">
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

        <div class="r4">
            <div id="export">{data}</div>

            <!-- <textarea disabled>
                {data}
            </textarea> -->
        </div>
        
        <div class="r5">
            <p>You can proceed by closing this tab/window.</p>
        </div>
    </div>

    <div class="footer">
        <Footer/>
    </div>

</div>


<style lang="scss" scoped>

    .page-complete {
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

        justify-items: center;

        font-family: 'Poppins';

        .header {
            grid-area: h;
        }

        .content {
            grid-area: c;

            display: grid;
            grid-template:
                "r1" min-content        // Title
                "r2" min-content        // Description
                "r3" min-content        // Sign check
                "r4" min-content        // Data
                "r5" min-content / 1fr; // Closing

            .r1 {
                grid-area: r1;

                display        : flex;
                justify-content: center;
                align-items    : center;
            }

            .r2 {
                grid-area: r2;

                max-width: 800px;

                display        : flex;
                justify-content: center;
                align-items    : center;

                text-align: center;
            }

            .r3 {
                grid-area: r3;


                display        : flex;
                justify-content: center;
                align-items    : center;
            }

            .r4 {
                grid-area: r4;
                
                display        : flex;
                justify-content: center;
                align-items    : center;

                #export {
                    width : 700px;
                    height: 100px;
    
                    overflow-y: hidden;
                    overflow-x: hidden;
                    display: none;
                }
            }

            .r5 {
                grid-area: r5;

                display        : flex;
                justify-content: center;
                align-items    : center;
            }

        }

        .footer {
            grid-area: f;
        }

    }

</style>
