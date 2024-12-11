<template>
    <div class="sign-in__container">
        
        <SignIn />

    </div>
</template>

<script>

// @ is an alias to /src

import SignIn from '@/components/SignIn.vue';

export default {
    name: 'SignInView',
    components: {        
        SignIn
    }
}

</script>

<style>
.sign-in__container {
    display: grid;
    grid-template-rows: max-content minmax(min-content, 1fr);
    height: 100vh;

}

.sign-in__container .sign-in__wrapper {
    place-self: center;
}
</style>
