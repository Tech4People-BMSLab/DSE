<script lang="ts">

</script>

<header class="header">
    <div class="left">
        <img src='./images/logo/dp_logo.png' alt='DP logo'>
    </div>

    <div class="right">
        <img src='./images/logo/ut_logo.png' alt='UT logo'>
    </div>
</header>


<style lang="scss" scoped>

    .header {
        width : 100vw;
        height: 100px;

        display: grid;
        grid-template:
            "l . r" 1fr / 1fr 1fr 1fr;

        .left {
            grid-area: l;

            display        : flex;
            justify-content: flex-start;
            align-items    : center;

            padding-left: 1rem;

            img {
                height: 70px;
            }
        }

        .right {
            grid-area: r;

            display        : flex;
            justify-content: flex-end;
            align-items    : center;

            padding-right: 1rem;

            img {
                height: 70px;
            }
        }
    }

</style>
