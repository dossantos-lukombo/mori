<template>
    <router-view @hide-chat="toggleChatVisibility"></router-view>
    <Chat v-if="showSmallChat" />
  </template>
  
  <script>
  import Chat from "./components/Chat/Chat.vue";
  
  export default {
    name: "App",
    components: { Chat },
    data() {
      return {
        showSmallChat: true, // Controls visibility of the small chatbox
      };
    },
    methods: {
      toggleChatVisibility(hideChat) {
        this.showSmallChat = !hideChat;
      },
      createWebSocketConn() {
        if (this.$route.path === "/sign-in" || this.$route.path === "/register") {
          return;
        }
        this.$store.dispatch("createWebSocketConn");
      },
    },
  };
  </script>
  