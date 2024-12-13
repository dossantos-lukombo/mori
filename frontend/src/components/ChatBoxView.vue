<template>
  <div>
    <!-- Navbar -->
    <NavBarOn @toggle-sidebar="toggleSidebar" />
    <!-- Sidebar -->
    <Sidebar
      :isActive="isSidebarActive"
      :contactsList="contacts"
      @navigate-to="navigateTo"
    />
    <!-- Chatbox Content -->
    <div class="chatbox-view-wrapper">
      <header class="chatbox-view-header">
        <h1>{{ name }}</h1>
      </header>

      <div class="chatbox-view-content" ref="contentDiv">
        <div
          class="message"
          v-for="(message, index) in allMessages"
          :style="msgPosition(message)"
          :key="index"
        >
          <p class="message-author" v-if="displayName(message, index)">
            {{ message.sender.nickname }}
          </p>
          <p :class="getClass(message)" class="message-content">
            {{ message.content }}
          </p>
        </div>
      </div>

      <form
        @submit.prevent="sendMessage"
        class="chatbox-view-input"
        autocomplete="off"
        @keyup.enter="sendMessage"
      >
        <input
          type="text"
          placeholder="Type your message here..."
          ref="sendMessageInput"
        />
        <button type="submit"><i class="uil uil-message"></i></button>
        <Emojis
          :input="this.$refs.sendMessageInput"
          :messagebox="this.$refs.contentDiv"
        />
      </form>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import NavBarOn from "@/components/NavBarOn.vue";
import Sidebar from "@/components/Sidebar.vue";
import Emojis from "./Chat/Emojis.vue";

export default {
  props: ["name", "receiverId", "type"],
  components: { NavBarOn, Sidebar, Emojis },
  data() {
    return {
      previousMessages: [],
      isSidebarActive: false,
      contacts: [],
    };
  },
  computed: {
    allMessages() {
      return [...this.previousMessages];
    },
    ...mapState({
      myID: (state) => state.id,
    }),
  },
  watch: {
    receiverId: {
      immediate: true,
      handler() {
        this.getPreviousMessages();
      },
    },
  },
  methods: {
    async fetchContacts() {
      try {
        const response = await fetch("http://localhost:8081/contacts", {
          credentials: "include",
        });
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
    async getPreviousMessages() {
      try {
        const response = await fetch("http://localhost:8081/messages", {
          credentials: "include",
          method: "POST",
          body: JSON.stringify({
            type: this.type,
            receiverId: this.receiverId,
          }),
        });
        const data = await response.json();
        this.previousMessages = data.chatMessage || [];
        this.scrollToBottom();
      } catch (error) {
        console.error("Error fetching messages:", error);
      }
    },
    async sendMessage() {
      const sendMessageInput = this.$refs.sendMessageInput;
      if (sendMessageInput.value === "") return;

      const msgObj = {
        receiverId: this.receiverId,
        content: sendMessageInput.value,
        type: this.type,
      };

      try {
        const response = await fetch("http://localhost:8081/newMessage", {
          body: JSON.stringify(msgObj),
          method: "POST",
          credentials: "include",
        });
        const data = await response.json();

        if (data.type === "Success") {
          this.previousMessages.push({
            ...msgObj,
            senderId: this.myID,
            sender: { nickname: "You" },
          });
          this.scrollToBottom();
        } else {
          this.$toast.open({
            message: data.message,
            type: "warning",
          });
        }

        sendMessageInput.value = "";
      } catch (error) {
        console.error("Error sending message:", error);
      }
    },
    displayName(message, index) {
      if (message.senderId === this.myID) return false;
      if (index < 1) return true;
      return message.senderId !== this.allMessages[index - 1]?.senderId;
    },
    getClass(message) {
      return message.senderId === this.myID ? "sent-message" : "received-message";
    },
    msgPosition(message) {
      return {
        alignSelf: message.senderId === this.myID ? "flex-end" : "flex-start",
      };
    },
    scrollToBottom() {
      this.$nextTick(() => {
        this.$refs.contentDiv.scrollTop = this.$refs.contentDiv.scrollHeight;
      });
    },
  },
  created() {
    this.getPreviousMessages();
    this.fetchContacts();
  },
};
</script>


<style scoped>
.chatbox-view-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--page-bg);
  width: 60%;
}

.chatbox-view-header {
  color: var(--color-white);
  padding: 20px;
  text-align: center;
  font-size: 1.5em;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.chatbox-view-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.message-author {
  font-size: 0.9em;
  color: var(--purple-color);
  margin-bottom: 5px;
}

.message-content {
  padding: 10px;
  border-radius: 10px;
  word-break: break-word;
}

.sent-message {
  background-color: var(--purple-color);
  color: var(--color-white);
}

.received-message {
  background-color: var(--bg-neutral);
  color: var(--color-white);
}

.chatbox-view-input {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border-radius: 15px 15px 0 0;
  background-color: var(--bg-neutral);
}

.chatbox-view-input input {
  flex: 1;
  padding: 20px;
  border-radius: 15px;
  height: 45px;
  border: 1px solid var(--color-grey-light);
}

.chatbox-view-input button {
  background-color: var(--purple-color);
  color: var(--color-white);
  border: none;
  padding: 10px 15px;
  border-radius: 15px;
  font-size: 1.2em;
  cursor: pointer;
  transition: all 0.3s ease;
}

.chatbox-view-input button:hover {
  background-color: var(--hover-color);
}
</style>
