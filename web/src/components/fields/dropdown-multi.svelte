<script lang="ts">
    // ------------------------------------------------------------
    // : Imports
    // ------------------------------------------------------------
    import { onMount, onDestroy } from 'svelte'

    import { computePosition } from '@floating-ui/dom'
    import { autoUpdate }      from '@floating-ui/dom'
    import { flip }            from '@floating-ui/dom'
    import { offset }          from '@floating-ui/dom'

    import { map }        from '@/utils/utils'
    import { filter }     from '@/utils/utils'
    import { find_index } from '@/utils/utils'
    import { each }       from '@/utils/utils'
    import { is_array }   from '@/utils/utils'

    // ------------------------------------------------------------
    // : Initialization
    // ------------------------------------------------------------
    export let ref: any

    let opened           : boolean = false
    let element_component: HTMLElement
    let element_arrow    : HTMLElement
    let element_anchor   : HTMLElement
    let element_float    : HTMLElement

    let cleanup: Function

    onMount(async () => {
        ref.answer.selected = map(ref.options, (option: any) => {
            return option
        })

        cleanup = autoUpdate(element_anchor, element_float, () => {
            computePosition(element_anchor, element_float, {
                placement: 'bottom-start',
                middleware: [
                    offset(5),
                    flip(),
                ]
            }).then(({x,y}) => {
                element_float.style.left = `${x}px`
                element_float.style.top  = `${y}px`
            })
        })
    })

    onDestroy(async () => {
        cleanup()
    })

    function onOutsideClick(node: HTMLElement) {
        const handler = (event: any) => {
            if (node && !node.contains(event.target) && !event.defaultPrevented) {
                opened = false
            }
        }

        document.addEventListener('click', handler, true)

        return {
            destroy() {
                document.removeEventListener('click', handler, true)
            }
        }
    }

    function toggle() {
        opened = !opened
    }

    function select(selected: any) {
        if (selected === 'unselected') {
            each(ref.answer.selected, (option: any) => {
                option.checked = false
            })
        }

        const index = find_index(ref.answer.selected, (option: any) => {
            return option.value == selected.value
        })

        ref.answer.selected[index].checked = !ref.answer.selected[index].checked

        let selections = []
        selections = filter(ref.answer.selected, (option: any) => {
            return option.checked
        })

        selections = map(selections, (option: any) => {
            return option.label
        })

        ref.answer.label = selections.join(', ')
        ref.answer.valid = selections.length > 0
    }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="component" bind:this={element_component} use:onOutsideClick>
    
    <div class="title">{ref.question}</div>

    <div class="fields">
        <div class="field-selected" bind:this={element_anchor} on:click={toggle}>
            <span>{ref.answer.label}</span>
        </div>

        <div class="field-arrow" on:click={toggle}>
            <img 
                src          = './images/arrow.svg'
                alt          = "arrow"
                class:closed = {opened == false} 
                class:opened = {opened == true}
            />
        </div>

        <div class="field-options visible" bind:this={element_float} class:opened={opened}>
            {#if is_array(ref.answer.selected)}
                {#each ref.answer.selected as option}
                
                    <div class="field-option" on:click={() => select(option)}>
                        <div><input type="checkbox" bind:checked={option.checked} on:click={() => select(option)}/></div>
                        <div>{option.label}</div>
                    </div>
                
                {/each}
            {/if}
        </div>
    </div>
</div>


<style lang="scss" scoped>

    ::-webkit-scrollbar {
        width: 2px;
    }
    
    ::-webkit-scrollbar-track {
        background: #f1f1f1;
    }
    
    ::-webkit-scrollbar-thumb {
        background: #888;
    }
    
    ::-webkit-scrollbar-thumb:hover {
        background: #555;
    }


    .component {
        min-width: 100px;
        max-width: max-content;
        width: 100%;

        display: grid;
        grid-template:
            "t" 1fr
            "f" 1fr / 1fr;

        padding-left  : 1rem;
        padding-right : 1rem;
        padding-top   : 1rem;
        padding-bottom: 1rem;

        box-sizing: border-box;

        .title {
            grid-area: t;

            display        : flex;
            justify-content: flex-start;
            align-items    : center;
        }

        .fields {
            grid-area: f;

            width : 338px;
            height: 40px;

            display: grid;
            grid-template:
                "l r" 1fr / auto auto;

            border       : 1px solid #e9e9e9;
            border-radius: 8px;
            box-shadow   : 0px 4px 15px rgba(0,0,0, 0.1);

            overflow-x: hidden;
            overflow-y: hidden;

            background: white;

            .field-selected {
                grid-area: l;

                width : 300px;
                height: 40px;

                display        : flex;
                justify-content: flex-start;
                align-items    : center;

                padding-left : 1rem;
                box-sizing: border-box;

                overflow-x: hidden;
                overflow-y: hidden;

                cursor: pointer;
            }

            .field-arrow {
                grid-area: r;

                width : 20px;
                height: 40px;

                display        : flex;
                justify-content: center;
                align-items    : center;

                padding-right: 0.5rem;
                padding-left : 0.5rem;

                img {
                    width : 20px;
                    height: 20px;
                }

                .opened {
                    animation-name: rotate-arrow;
                    animation-duration: 100ms;
                    animation-direction: reverse;
                    animation-fill-mode: forwards;
                }

                .closed {
                    animation-name: rotate-arrow-reverse;
                    animation-duration: 100ms;
                    animation-direction: normal;
                    animation-fill-mode: forwards;
                }

                @keyframes rotate-arrow {
                    0%   { 
                        opacity: 1;
                        transform: rotate(180deg); 
                    }
                    100% { 
                        opacity: 1;
                        transform: rotate(0deg); 
                    }
                }
                
                @keyframes rotate-arrow-reverse {
                    0%   { 
                        opacity: 1;
                        transform: rotate(180deg); 
                    }
                    100% { 
                        opacity: 1;
                        transform: rotate(0deg); 
                    }
                }
            }

            .field-options {
                visibility: hidden;

                z-index: 100;

                width: calc(300px + 20px + 1rem);
                max-height: 300px;

                position: absolute;
                top : 0;
                left: 0;

                overflow-x: hidden;
                overflow-y: scroll;

                border       : 1px solid #e9e9e9;
                border-radius: 8px;
                box-shadow   : 0px 4px 15px rgba(0,0,0, 0.1);

                background: white;

                .field-option {
                    height: 40px;

                    display        : flex;
                    justify-content: flex-start;
                    align-items    : center;
                    flex-direction : row;

                    gap: 1rem;

                    overflow-x: hidden;
                    overflow-y: hidden;
                    
                    padding-left : 1rem;
                    box-sizing: border-box;

                    &:hover {
                        color: white;
                        background: #113cfc;
                    }
                }
            }

            .opened {
                visibility: visible;
            }
        }
        
        
    }

</style>
