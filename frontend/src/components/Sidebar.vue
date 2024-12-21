<template>
  <div :class="['sidebar', { 'sidebar--active': isActive }]">
    <div class="sidebar-content">
      <ul class="icon-container">
        <li @click="navigateToMessages" class="icon-wrapper">
          <div class="icon-circle">
            <img src="@/assets/icons/messages.svg" alt="Messagerie" />
          </div>
          <span>Messages</span>
        </li>
        <li @click="navigateToChatBot" class="icon-wrapper">
          <div class="icon-circle">
            <img src="@/assets/icons/chat.svg" alt="Chat" />
          </div>
          <span>ChatBot</span>
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

.icon-container {
  display: flex;
  justify-content: center; /* Center horizontally */
  gap: 20px; /* Space between items */
  margin-bottom: 20px;
}

.icon-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  text-align: center;
  color: var(--text-color);
  transition: transform 0.3s ease;
}

.icon-wrapper:hover {
  transform: scale(1.1);
}

.icon-circle {
  width: 70px;
  height: 70px;
  background-color: var(--purple-color);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
  transition: background-color 0.3s ease;
}

.icon-circle:hover {
  background-color: var(--hover-color);
}

.icon-circle img {
  width: 32px;
  height: 32px;
}

.icon-wrapper span {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}
</style>
