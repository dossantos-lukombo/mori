<template>

    <div class="register__container">       
        <Register />
    </div>

</template>

<script>
import Register from "@/components/Register.vue";

export default {
    name: 'RegisterView',
    components: { Register }
}
</script>


<style>
.register__container {
    display: grid;
    grid-template-rows: max-content 1fr;
    height: 100vh;
}
</style>