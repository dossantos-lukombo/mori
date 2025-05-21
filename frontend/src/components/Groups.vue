<template>
    <div class="item-list__wrapper" id="groups" :style="groupsSectionStyle">
        <h3 v-if="isProfilePage">Groups</h3>
        <h3 v-else>My groups</h3>
        <ul class="item-list">
            <li v-for="group in groups" v-bind:key="group.id" v-if="groups !== null">
                <div class="group-icon-container">
                    <img class="group-icon" src="../assets/icons/users-alt.svg" alt="Group icon">
                </div>
                <div class="item-text">
                    <router-link :to="{ path: `/group/${group.id}`}">{{ group.name }}</router-link>
                </div>
            </li>

            <p class="additional-info" v-else>{{noGroupsText}}</p>

        </ul>
        <NewGroup v-if="!isProfilePage"/>
    </div>

</template>


<script>
import NewGroup from '@/components/NewGroup.vue';

export default {
    props: ['groups'],
    name: 'Groups',
    components: { NewGroup },

    computed: {
        // what text to display if no groups in groups section
        noGroupsText() {
            if (this.isProfilePage) {
                return "User is not part of any group"
            } else {
                return "You are not part of any group"
                
            }
        },

        // groups section alignment is left on home page and middle on profile page
        groupsSectionStyle() {
            return {
                alignItems: this.isProfilePage ? 'center' : 'flex-start'
            }
        },

        isProfilePage() {
            return this.$route.name === "Profile"
        }
    }
}
</script>

<style scoped>
.group-icon-container {
    width: 40px;
    height: 40px;
    min-width: 40px;  /* Prevent shrinking */
    min-height: 40px; /* Prevent shrinking */
    overflow: hidden;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--bg-neutral);
    border: 2px solid var(--purple-color);
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
    position: relative; /* For absolute positioning of child */
}

.group-icon {
    width: 24px;
    height: 24px;
    object-fit: contain;
    display: block; /* Force block display for Chrome */
    -webkit-user-drag: none; /* Prevent drag in webkit */
    position: absolute; /* Force position */
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%); /* Center exactly */
}

.item-list li {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 0;
}
</style>
