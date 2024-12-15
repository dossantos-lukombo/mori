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
    <div id="layout">
      <div class="chatbox-view-wrapper">
        <header class="chatbox-view-header">
          <div class="header-left-part">
            <div
              class="receiver-avatar"
              :style="{
                backgroundImage: user.avatar
                  ? `url(${user.avatar})`
                  : 'url(default-avatar.png)',
              }"
            ></div>
            <h1 class="receiver-name">{{ user.name || "Unnamed User" }}</h1>
          </div>
          <p class="follow-status">{{ user.following }}</p>
        </header>

        <div class="chatbox-view-content" ref="contentDiv">
          <div
            class="message"
            v-for="(message, index) in allMessages"
            :style="msgPosition(message)"
            :key="index"
          >
          <div class="receiver-avatar-name">
          <div class="receiver-avatar-chat" v-if="displayName(message, index)" :style="{
            backgroundImage: user.avatar
              ? `url(${user.avatar})`
              : 'url(default-avatar.png)',
          }"></div>
            <p class="message-author" v-if="displayName(message, index)">
              {{ message.sender.nickname }}
            </p>
          </div>
            <p :class="getClass(message)" class="message-content">
              {{ message.content }}

              <p class="message-timeStamp">{{ formatTime(message.createdAt) }}</p>
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
  </div>
</template>

<script>
import { mapState } from "vuex";
import NavBarOn from "@/components/NavBarOn.vue";
import Sidebar from "@/components/Sidebar.vue";
import Emojis from "../components/Chat/Emojis.vue";

export default {
  props: ["name", "receiverId", "type"],
  components: { NavBarOn, Sidebar, Emojis },
  data() {
    return {
      user: {
        name: "", // Default name
        avatar: "default-avatar.png", // Default avatar
      },
      previousMessages: [],
      isSidebarActive: false,
      contacts: [],
    };
  },

  computed: {
    allMessages() {
      // Ensure no duplicates between previousMessages and store messages
      const storeMessages = this.$store.getters.getMessages(
        this.receiverId,
        this.type
      );
      const uniqueMessages = storeMessages.filter(
        (msg) => !this.previousMessages.some((prevMsg) => prevMsg.id === msg.id)
      );

      return [...this.previousMessages, ...uniqueMessages];
    },
    ...mapState({
      myID: (state) => state.id,
    }),
  },

  watch: {
    allMessages() {
      this.$nextTick(() => {
        this.scrollToBottom();
      });
    },
    receiverId: {
      immediate: true,
      handler(newId) {
        this.fetchUserDetails(newId);
        this.getPreviousMessages();
      },
    },
  },
  methods: {
    formatTime(timestamp) {
      if (!timestamp) return "Invalid timestamp";

      const date = new Date(timestamp);

      // Extract date components
      const day = date.getDate().toString().padStart(2, "0"); // Two-digit day
      const month = (date.getMonth() + 1).toString().padStart(2, "0"); // Two-digit month
      const year = date.getFullYear();

      // Extract time components
      const hours = date.getHours().toString().padStart(2, "0"); // Two-digit hours
      const minutes = date.getMinutes().toString().padStart(2, "0"); // Two-digit minutes

      // Combine into desired format: DD/MM/YYYY HH:mm
      return `${day}/${month}/${year} ${hours}:${minutes}`;
    },
    removeZFromTimestamp(timestamp) {
      if (!timestamp || typeof timestamp !== "string") {
        return "Invalid timestamp";
      }
      return timestamp.replace("Z", "");
    },

    async fetchUserDetails(userId) {
      try {
        const response = await fetch("http://localhost:8081/allUsers", {
          credentials: "include",
        });
        const data = await response.json();
        const user = data.users.find((user) => user.id === userId);
        if (user) {
          if (user.follower == true) {
            user.following = "Follows you";
          } else {
            user.following = "Not Following you";
          }
          this.user = {
            name: user.nickname,
            following: user.following,
            avatar:
              `http://localhost:8081/${user.avatar}` || "default-avatar.png", // Include full URL for avatar
          };
        } else {
          this.user = { name: "Unknown User", avatar: "default-avatar.png" }; // Fallback for invalid user
        }
      } catch (error) {
        console.error("Error fetching user details:", error);
        this.user = { name: "Unknown User", avatar: "default-avatar.png" }; // Fallback for errors
      }
    },
    updateLayoutWidth() {
      const sidebar = document.querySelector(".sidebar");
      const layout = document.getElementById("layout");

      if (sidebar?.classList.contains("sidebar--active")) {
        layout.style.width = "70%";
      } else {
        layout.style.width = "100%";
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
        console.log(data);

        // Filter out messages already in Vuex
        const storeMessages = this.$store.getters.getMessages(
          this.receiverId,
          this.type
        );

        // Format the `createdAt` field
        this.previousMessages = (data.chatMessage || [])
          .filter(
            (msg) => !storeMessages.some((storeMsg) => storeMsg.id === msg.id)
          )
          .map((msg) => {
            return {
              ...msg,
              createdAt: this.removeZFromTimestamp(msg.createdAt), // Add formatted timestamp
            };
          });

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
        createdAt: new Date(),
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
          this.$store.dispatch("addNewChatMessage", {
            ...msgObj,
            senderId: this.myID,
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
    clearChatNewMessages() {
      if (this.type === "GROUP") {
        let msgs = this.$store.state.chat.newGroupChatMessages.filter(
          (msg) => msg.receiverId !== this.receiverId
        );
        this.$store.commit("updateNewGroupChatMessages", msgs);
      } else {
        let msgs = this.$store.state.chat.newChatMessages.filter(
          (msg) =>
            msg.receiverId !== this.receiverId &&
            msg.senderId !== this.receiverId
        );
        this.$store.commit("updateNewChatMessages", msgs);
      }
    },
    displayName(message, index) {
      if (message.senderId === this.myID) return false;
      if (index < 1) return true;
      return message.senderId !== this.allMessages[index - 1]?.senderId;
    },
    getClass(message) {
      return message.senderId === this.myID
        ? { "sent-message": true }
        : { "received-message": true };
    },
    msgPosition(message) {
      return {
        alignSelf: message.senderId === this.myID ? "flex-end" : "flex-start",
      };
    },
    scrollToBottom() {
      this.$nextTick(() => {
        if (this.$refs.contentDiv) {
          this.$refs.contentDiv.scrollTop = this.$refs.contentDiv.scrollHeight;
        }
      });
    },
  },
  async mounted() {
    this.fetchUserDetails(this.receiverId);

    const sidebar = document.querySelector(".sidebar");

    // Ensure layout width is set initially
    this.updateLayoutWidth();

    // Observe changes to the sidebar class
    this.sidebarObserver = new MutationObserver(() => {
      this.updateLayoutWidth();
    });

    this.sidebarObserver.observe(sidebar, {
      attributes: true, // Watch for changes to attributes (like `class`)
      attributeFilter: ["class"], // Only observe the `class` attribute
    });
  },
  beforeUnmount() {
    // Disconnect the observer to prevent memory leaks
    if (this.sidebarObserver) {
      this.sidebarObserver.disconnect();
    }
  },
  created() {
    console.log("Component created: chatbox");
    this.getPreviousMessages();
  },
  unmounted() {
    console.log("Component unmounted: chatbox");
    this.clearChatNewMessages();
  },
};
</script>

<style scoped>
#layout {
  display: flex;
  height: 95vh;
  width: 100%;
  position: fixed;
  bottom: 0px;
  align-items: center;
  justify-content: center;
  right: 0px;
  transition: all 0.3s ease;
}
.chatbox-view-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--page-bg);
  width: 60%;
  outline: 1px solid var(--bg-neutral);
}
.chatbox-view-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  background-color: var(--bg-neutral);
  padding: 15px;
  border-radius: 0 0 15px 15px;
  box-shadow: 0 2px 10px rgb(0, 0, 0);
  z-index: 1;
  font-size: 1.5em;
}

.header-left-part {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: left;
  color: var(--color-white);
  gap: 15px;
  font-size: 1.5em;
}

.receiver-name {
  align-items: center;
  justify-content: center;
  text-align: center;
  margin-top: 15px;
  text-align: center;
  font-size: 23px;
}
.receiver-avatar {
  margin-top: 15px;
  margin-left: 10px;
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-position: center;
  background-repeat: no-repeat;
  background-size: contain;
  background-color: wheat;
  border: 2px solid var(--purple-color);
}

.receiver-avatar-name{
  display: flex;
  align-items: end;
  gap: 6px;
}
.receiver-avatar-chat{
  box-shadow: 0 2px 10px rgb(0, 0, 0);
  background-color: wheat;
  border: 2px solid var(--purple-color);
  border-radius: 50%;
  background-position: center;
  background-repeat: no-repeat;
  background-size: contain;
  width: 40px;
  height: 40px;
  margin-bottom: 15px;
}

.follow-status {
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  margin-top: 18px;
  margin-right: 30px;
  color: var(--purple-color);
  font-size: 16px;
}

.chatbox-view-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
border-radius: 15px;
  gap: 10px;
}

.message-author {
  font-size: 0.9em;
  color: var(--purple-color);
  margin-bottom: 15px;
}

.message-content {
  padding: 10px;
  border-radius: 10px;
  word-break: break-word;
}

.message-timeStamp{
  font-size: 0.7em;
  color: var(--color-grey);
  opacity: 0.5;
  text-align: right;
  margin-top: 5px;
}

.message{
  max-width: 80%;
}
.sent-message {
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.3);
  background-color: var(--purple-color);
  color: var(--color-white);
}

.received-message {
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.3);
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
  box-shadow: 0 2px 10px rgb(0, 0, 0);
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
