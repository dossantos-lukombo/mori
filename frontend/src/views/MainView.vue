<template>
    <NavBarOn @toggle-sidebar="toggleSidebar" />
    <Sidebar
      :isActive="isSidebarActive"
      :contactsList="contacts"
      @navigate-to="navigateTo"
    />
    <div id="layout">
      <div class="main-content">
        <ChatbotConversation />
      </div>
    </div>
  </template>
  
  <script>
  import NavBarOn from "@/components/NavBarOn.vue";
  import Sidebar from "@/components/Sidebar.vue";
  import ChatbotConversation from "@/components/ChatbotConversation.vue";
  
  export default {
    components: {
      NavBarOn,
      Sidebar,
      ChatbotConversation,
    },
    data() {
      return {
        contacts: [],
        isSidebarActive: false,
      };
    },
    methods: {
      async fetchContacts() {
        try {
          const response = await fetch("http://localhost:8081/contacts", { credentials: "include" });
          const data = await response.json();
          this.contacts = data.contacts || [];
        } catch (error) {
          console.error("Error fetching contacts:", error);
        }
      },
      toggleSidebar() {
        this.isSidebarActive = !this.isSidebarActive;
      },
      navigateTo(target) {
        if (target === "chatbot") {
          this.$router.push({ name: "mainpage" });
        } else if (target === "messages") {
          this.$router.push({ name: "messages" });
        }
      },
    },
    created() {
      this.fetchContacts();
    },
  };
  </script>
  

<style>
html,
body {
  overflow-y: hidden;
}

#layout {
  display: flex;
  height: 95vh;
  width: 100%;
  position: fixed;
  bottom: 0px;
}

.main-content {
  height: 100%;
  width: fit-content;
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: var(--page-bg);
}
</style>
