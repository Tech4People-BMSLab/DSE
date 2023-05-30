<script lang="ts">
    // ------------------------------------------------------------
    // Imports
    // ------------------------------------------------------------
    import { page }            from '../stores/store'
    import { error_message }   from '../stores/store';
    import { success_message } from '../stores/store';  

    import { storage }        from '../utils/util'
    import { report }         from '../utils/util'

    import { lang } from '../stores/store'

    // ------------------------------------------------------------
    // Components
    // ------------------------------------------------------------
    import InputEmail    from '../lib/InputEmail.svelte'
    import InputTextArea from '../lib/InputTextArea.svelte'
    import Button        from '../lib/Button.svelte'

    import { DateTime } from 'luxon' // To avoid users spam sending feedbacks

    // ------------------------------------------------------------
    // Props
    // ------------------------------------------------------------
    const i18n = {
        'nl': {
            1: 'E-mailadres (niet verplichts)',
            2: 'Waarom heb je het aanmeldproces geannuleerd?',

            3: 'Verstuur',
            4: 'Terug',

            5: 'Je kunt maximaal een feedback per 30 minuten versturen',
            6: 'Niet alle velden zijn correct ingevuld',
            7: 'Je bericht is te kort. Minimaal 25 tekens',
            8: 'Er is een fout opgetreden. Probeer het later opnieuw',

            9: 'Bedankt voor je feedback, je bericht is verzonden',
            
        },
        'en': {
            1: 'Email address (optional)',
            2: 'Why did you cancel the registration process?',

            3: 'Send',
            4: 'Back',

            5: 'You can only send a feedback once every 30 minutes',
            6: 'Not all fields are filled in correctly',
            7: 'Your message is too short. Minimum 25 characters',
            8: 'An error has occurred. Please try again later',

            9: 'Thank you for your feedback, your message has been sent',
        },
    }

    // Message of the user
    let email             = ''
    let message           = ''

    // Event variables
    let submit            = false // To notify the child components (inputs)

    // Shows the length of text area 
    let length      = 0
    let min_length  = 25
    let max_length  = 1000

    // Spam avoidance
    let timestamp

    // ------------------------------------------------------------
    // : Functions
    // ------------------------------------------------------------
    function goto_consent() {
        error_message.set('')
        page.set('consent')
    }

    /**
     * Sends a feedback email to info@digitalepolarisatie.nl
     */
    async function on_send() {
        // Check if the user is spamming (if less than 30 minutes)
        if (timestamp) {
            if (DateTime.local().diff(timestamp).as('minutes') < 30) {
                error_message.set(i18n[$lang][5])
                return
            }
        }

        // Trigger the submit event
        submit = true

        // Show error if email and message are not filled in
        if (email.length == 0 || message.length == 0) {
            error_message.set(i18n[$lang][6])
            return
        }

        // Show error if email is not in the correct format (using regex)
        if (!email.match(/^([a-z0-9_\.-]+)@([\da-z\.-]+)\.([a-z\.]{2,6})$/g)) {
            error_message.set(i18n[$lang][6])
            return
        }

        // Show error if message does not meet minimum length
        if (message.length < min_length) {
            error_message.set(i18n[$lang][7])
            return
        }
        
        // Prepare payload
        const payload = {
            email,
            message
        }
        
        const api_url  = `${await storage.get('api_url')}/feedback`
        
        try {
            const response = await fetch(api_url, {
                method: 'POST',
                headers: {
                    'content-type': 'application/json',
                },
                body: JSON.stringify(payload),
            })

            if (!response.ok) {
                error_message.set(i18n[$lang][8])
                throw new Error(`HTTP status code ${response.status}`)
            }

            // Show success message
            success_message.set(i18n[$lang][9])

            // Set the time that the feedback was sent
            timestamp = DateTime.local()
        } catch (error) {
            console.error(`Failed to send feedback: ${error}`)
            console.error(error)

            report(`Failed to send feedback (reason: ${error})`)
        }
    }

    async function on_back() {
        goto_consent()
    }

    // ------------------------------------------------------------
    // : Reactive
    // ------------------------------------------------------------
    $: {
        if (timestamp) {
            storage.set({'timestamp_feedback': timestamp})
        }
    }

</script>

<section class='section'>
    <div class='form'>

        <div class='reason-wrapper'>
            <div class='email'>
                <label class='label' for=''>{i18n[$lang][1]}</label>
                <InputEmail 
                    submit     = {submit}
                    bind:email = {email} 
                />
            </div>

            <div class='message'>
                <label class='label' for=''>{i18n[$lang][2]} (<span>{length}/{max_length}</span>)</label>
                <InputTextArea 
                    submit      = {submit} 
                    min_length  = {min_length}
                    max_length  = {max_length}
                    bind:text   = {message} 
                    bind:length = {length}
                />
            </div>
        </div>

        <div class='btn-group'>
            <Button primary={true} onclick={on_send}>{i18n[$lang][3]}</Button>
            <Button secondary={true} onclick={on_back}>{i18n[$lang][4]}</Button>
        </div>
    </div>
</section>


<style type="text/scss">
    .section {
        height: 100%;
        width : 100%;

        display        : flex;
        justify-content: flex-start;
        align-items    : flex-start;
        flex-direction : column;

        .form {
            width : 100%;
            height: 100%;

            .message {
                padding-top: 2rem;
            }
    
            .label {
                height: 30px;
                
                display        : flex;
                justify-content: flex-start;
                align-items    : center;
    
                padding-top   : 1rem;
    
                color: #848484;
                font-family: 'Roboto';
                font-style : normal;
                font-weight: 700;
                font-size  : 17px;
                line-height: 32px;
            }

            .btn-group {
                display        : flex;
                justify-content: center;
                align-items    : center;
                flex-direction : row;
                gap: 1rem;
    
                padding-top   : 2rem;
            }
        }
    }
</style>