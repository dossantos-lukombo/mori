<template>
  <div class="contacts-wrapper">
    <h2 class="titre">Contacts</h2>

    <!-- Section Amis -->
    <h3 class="sous_titres">Amis</h3>
    <ul class="horizontal-list">
      <li
        v-for="contact in chatUserList"
        :key="contact.id"
        @click="selectContact(contact, 'PERSON')"
        class="contact-item-horizontal"
      >
        <div
          class="user-picture small"
          :style="{ backgroundImage: `url(http://localhost:8081/${contact.avatar})` }"
        ></div>
        <div class="contact-name">{{ contact.nickname }}</div>
      </li>
    </ul>

    <!-- Section Conversations Amis -->
    <h3 class="sous_titres">Conversations Amis</h3>
    <div class="conversation-card-wrapper">
      <div
        v-for="convMsg in friends"
        :key="convMsg.id"
        class="conversation-card"
        @click="selectContact(convMsg, 'PERSON')"
      >
        <div
          class="avatar"
          :style="{ backgroundImage: `url(http://localhost:8081/${convMsg.avatar})` }"
        ></div>
        <div class="content">
          <div class="header">
            <span class="name">{{ convMsg.name }}</span>
            <span class="time">{{ formatTime(convMsg.lastMessageTime) }}</span>
          </div>
          <div class="message-preview">
            {{ convMsg.lastMessage.length > 40
              ? convMsg.lastMessage.slice(0, 40) + "…"
              : convMsg.lastMessage }}
          </div>
        </div>
      </div>
    </div>

    <!-- Section Groupes -->
    <h3 class="sous_titres">Groupes</h3>
    <NewGroup />
    <ul class="horizontal-list">
      <li
        v-for="group in userGroups"
        :key="group.id"
        @click="selectContact(group, 'GROUP')"
        class="contact-item-horizontal"
      >
        <div
          class="user-picture small"
          :style="{
            backgroundImage: `url(http://localhost:8081/${group.avatar || 'defaultGroup.png'})`
          }"
        ></div>
        <div class="contact-name">{{ group.name }}</div>
      </li>
    </ul>

    <!-- Section Conversations Groupes -->
    <h3 class="sous_titres">Conversations Groupes</h3>
    <div class="conversation-card-wrapper">
      <div
        v-for="convMsg in groups"
        :key="convMsg.id"
        class="conversation-card"
        @click="selectContact(convMsg, 'GROUP')"
      >
        <div
          class="avatar"
          :style="{ backgroundImage: `url(http://localhost:8081/${convMsg.avatar || 'defaultGroup.png'})` }"
        ></div>
        <div class="content">
          <div class="header">
            <span class="name">{{ convMsg.name }}</span>
            <span class="time">{{ formatTime(convMsg.lastMessageTime) }}</span>
          </div>
          <div class="message-preview">
            {{ convMsg.lastMessage.length > 40
              ? convMsg.lastMessage.slice(0, 40) + "…"
              : convMsg.lastMessage }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import NewGroup from "@/components/NewGroup.vue";

export default {
  name: "ContactsForChatBotView",
  components: {
    NewGroup,
  },
  computed: {
    ...mapState({
      chatUserList: (state) => state.chat.chatUserList,
      userGroups: (state) => state.groups.userGroups,
      conversationsMsg: (state) => state.conversationsMsg,
    }),
    ...mapGetters([
      "getUnreadMessagesCount",
      "getUnreadGroupMessagesCount",
      "getUnreadMsgsCountFromDB",
    ]),
    friends() {
      return this.conversationsMsg.filter((c) => c.type === "PERSON");
    },
    groups() {
      return this.conversationsMsg.filter((c) => c.type === "GROUP");
    },
  },
  created() {
    this.$store.dispatch("fetchConversationsMsg");
  },
  methods: {
    selectContact(contact, type) {
      this.$emit("select-contact", {
        id: contact.id,
        name: contact.nickname || contact.name,
        type,
      });
    },
    formatTime(isoString) {
      const date = new Date(isoString);
      return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    },
  },
};
</script>

<style scoped>
.contacts-wrapper {
  padding: 20px;
  background-color: var(--bg-neutral);
}

.titre {
  color: white;
  font-size: 1.8em;
}

.sous_titres {
  color: white;
  font-size: 1.2em;
  margin-top: 2vh;
  margin-bottom: 10px;
}

/* Liste horizontale */
.horizontal-list {
  display: flex;
  flex-wrap: nowrap;
  overflow-x: auto;
  gap: 10px;
  list-style: none;
  padding: 0;
  margin: 10px 0;
}

.contact-item-horizontal {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  cursor: pointer;
  width: 60px;
}

.user-picture.small {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-size: cover;
  background-position: center;
  margin-bottom: 5px;
}

.contact-name {
  font-size: 12px;
  color: var(--text-primary);
}

/* Cards pour les conversations */
.conversation-card-wrapper {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.conversation-card {
  display: flex;
  align-items: center;
  background-color: #3a3a3a;
  border-radius: 10px;
  padding: 10px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.conversation-card:hover {
  background-color: #4a4a4a;
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-size: cover;
  background-position: center;
  margin-right: 15px;
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.name {
  font-size: 14px;
  font-weight: bold;
  color: white;
}

.time {
  font-size: 12px;
  color: #ccc;
}

.message-preview {
  font-size: 13px;
  color: #ccc;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
