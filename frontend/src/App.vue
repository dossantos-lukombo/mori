<template>
  <router-view @hide-chat="toggleChatVisibility"></router-view>
  <Chat v-if="showSmallChat" />
</template>

<script>
import Chat from "./components/Chat/Chat.vue";
import { mapState } from "vuex";

export default {
  name: "App",
  components: { Chat },
  data() {
    return {
      showSmallChat: true, // Controls visibility of the small chatbox
    };
  },
  computed: {
    ...mapState({
      currentRoute: (state) => state.route, // Map current route from Vuex or Vue Router
    }),
    shouldHideSmallChat() {
      // Hide small chatbox on specific routes
      const hiddenRoutes = ["/sign-in", "/register", "/messages"];
      return hiddenRoutes.includes(this.$route.path);
    },
  },
  watch: {
    // Watch route changes to toggle small chat visibility
    $route() {
      this.toggleChatVisibility();
    },
  },
  methods: {
    toggleChatVisibility() {
      this.showSmallChat = !this.shouldHideSmallChat;
    },
    createWebSocketConn() {
      const excludedPaths = ["/sign-in", "/register"];
      if (excludedPaths.includes(this.$route.path)) {
        return;
      }
      this.$store.dispatch("createWebSocketConn");
    },
  },
  created() {
    // Set the initial state of the small chatbox based on the current route
    this.toggleChatVisibility();
  },
};
</script>
