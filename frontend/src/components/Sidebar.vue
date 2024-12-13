<template>
  <div :class="['sidebar', { 'sidebar--active': isActive }]">
    <div class="sidebar-content">
      <ul>
        <li @click="navigateToMessages">
          <div class="box">Messages</div>
        </li>
        <li @click="navigateToChatBot">
          <div class="box">ChatBot</div>
        </li>
      </ul>

      <ContactsForChatBotView
        v-if="activeView === 'contacts'"
        @select-contact="handleContactSelection"
      />
    </div>
  </div>
</template>

<script>
import ContactsForChatBotView from "./ContactsForChatBoxView.vue";

export default {
  props: {
    isActive: {
      type: Boolean,
      required: true,
    },
    contactsList: {
      type: Array,
      required: true, // Ensure we receive the contacts list as a prop
    },
  },
  data() {
    return {
      activeView: null, // Manage the active view in the sidebar
    };
  },
  components: { ContactsForChatBotView },
  methods: {
    async navigateToMessages() {
      if (this.contactsList.length > 0) {
        const firstContact = this.contactsList[0];
        await this.$router.push({
          name: "messages",
          query: {
            name: firstContact.nickname,
            receiverId: firstContact.id,
            type: "PERSON",
          },
        });
      } else {
        this.activeView = "contacts";
      }
    },
    async navigateToChatBot() {
      await this.$router.push({ name: "mainpage" });
    },
    handleContactSelection(contact) {
      this.$router.push({
        name: "messages",
        query: {
          name: contact.nickname,
          receiverId: contact.id,
          type: "PERSON",
        },
      });
    },
  },
};
</script>


<style scoped>
/* Same styles */
.sidebar {
  position: fixed;
  top: 64.45px;
  left: -440px;
  width: 440px;
  height: calc(100% - 64.45px);
  background-color: var(--bg-neutral);
  transition: left 0.3s ease;
  z-index: 2;
  overflow-y: auto;
}

.sidebar--active {
  left: 0;
}

.sidebar-content {
  padding: 20px;
}

.box {
  background-color: var(--purple-color);
  padding: 15px;
  border-radius: 8px;
  text-align: center;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.box:hover {
  background-color: var(--hover-color);
}
</style>
