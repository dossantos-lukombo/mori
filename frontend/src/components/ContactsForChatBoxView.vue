<template>
  <div class="contacts-wrapper">
    <h2 class="titre">Contacts</h2>
    <h3 class="sous_titres">Amis</h3>
    <ul class="contacts-list">
      <li
        v-for="contact in chatUserList"
        :key="contact.id"
        @click="selectContact(contact, 'PERSON')"
        class="contact-item"
      >
        <div
          class="user-picture medium"
          :style="{ backgroundImage: `url(http://localhost:8081/${contact.avatar})` }"
        ></div>
        <div class="contact-name">{{ contact.nickname }}</div>
        <span
          v-if="totalUnreadMessagesCount(contact.id, 'PERSON') !== 0"
          class="unread-messages"
        >
          {{ totalUnreadMessagesCount(contact.id, 'PERSON') }}
        </span>
      </li>
    </ul>

    <div class="conversations-wrapper">
      <ul class="conversation-list">
        <li
          v-for="convMsg in friends"
          :key="convMsg.id"
          class="conversation-card"
          @click="selectContact(convMsg, 'PERSON')"
        >
          <!-- Avatar à gauche -->
          <div
            class="avatar"
            :style="{ backgroundImage: `url(http://localhost:8081/${convMsg.avatar})` }"
          ></div>

          <!-- Contenu principal -->
          <div class="content">
            <!-- Ligne du haut : nom + heure -->
            <div class="header">
              <span class="name">{{ convMsg.name }}</span>
              <span class="time">{{ formatTime(convMsg.lastMessageTime) }}</span>
            </div>
            <!-- Aperçu du dernier message -->
            <div class="message-preview">
              {{ convMsg.lastMessage.length > 40
                ? convMsg.lastMessage.slice(0, 40) + "…"
                : convMsg.lastMessage }}
            </div>
          </div>
        </li>
      </ul>
      <h3 class="sous_titres">Groupes</h3>
      <ul class="contacts-list">
        <li
          v-for="group in userGroups"
          :key="group.id"
          @click="selectContact(group, 'GROUP')"
          class="contact-item"
        >
          <div class="contact-name">{{ group.name }}</div>
          <span
            v-if="totalUnreadMessagesCount(group.id, 'GROUP') !== 0"
            class="unread-messages"
          >
            {{ totalUnreadMessagesCount(group.id, 'GROUP') }}
          </span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";

export default {
  name: "ContactsForChatBotView",
  data() {
    return {};
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
      console.log("[computed friends] conversationsMsg =", this.conversationsMsg);
      return this.conversationsMsg.filter((c) => c.type === "PERSON");
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
    totalUnreadMessagesCount(receiverId, type) {
      if (type === "PERSON") {
        return (
          this.getUnreadMessagesCount(receiverId) +
          this.getUnreadMsgsCountFromDB(receiverId)
        );
      } else {
        return (
          this.getUnreadGroupMessagesCount(receiverId) +
          this.getUnreadMsgsCountFromDB(receiverId)
        );
      }
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
  margin-bottom: 20px; /* Augmente l'espace sous le titre "Amis" */
}

/* Liste des contacts */
.contacts-list {
  display: flex;
  flex-wrap: wrap;
  gap: 20px; /* Espace entre les amis */
  list-style: none;
  padding: 0;
  margin: 0 0 30px 0; /* Espace entre les amis et la conversation */
}

.contact-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  cursor: pointer;
  width: 80px;
}

.contact-item:hover {
  background-color: var(--hover-color);
  border-radius: 10px;
}

.user-picture.medium {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-size: cover;
  background-position: center;
  margin-bottom: 5px;
}

.contact-name {
  font-size: 14px;
  color: var(--text-primary);
}

.unread-messages {
  margin-top: 5px;
  background-color: var(--purple-color);
  color: var(--color-white);
  padding: 5px;
  border-radius: 10px;
  font-size: 12px;
  text-align: center;
}

/* Liste des conversations */
.conversation-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.conversation-card {
  display: flex;
  align-items: center;
  background-color: #3a3a3a;
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 10px; /* Espace entre les conversations */
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
  margin-right: 10px;
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
/* Responsivité pour les écrans moyens */
@media (max-width: 768px) {
  .contacts-list {
    gap: 15px; /* Réduit l'espacement entre les amis */
    justify-content: center;
  }

  .contact-item {
    width: 70px; /* Réduit la largeur pour s'adapter */
  }

  .user-picture.medium {
    width: 50px;
    height: 50px; /* Taille plus petite pour l'avatar */
  }

  .conversation-card {
    flex-direction: column; /* Place les éléments verticalement */
    align-items: flex-start; /* Aligne le contenu à gauche */
    padding: 15px;
  }

  .avatar {
    margin-right: 0;
    margin-bottom: 10px;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
  }

  .name {
    font-size: 14px;
    margin-bottom: 5px;
  }

  .time {
    font-size: 12px;
  }

  .message-preview {
    font-size: 12px;
  }
}

/* Responsivité pour les petits écrans */
@media (max-width: 480px) {
  .contacts-list {
    flex-direction: column; /* Place les amis verticalement */
    align-items: center;
    gap: 10px; /* Réduit davantage l'espacement */
  }

  .contact-item {
    width: 100%; /* Prend toute la largeur disponible */
  }

  .conversation-card {
    flex-direction: column;
    align-items: flex-start;
    padding: 10px;
  }

  .avatar {
    width: 40px;
    height: 40px;
    margin-bottom: 10px;
  }

  .header {
    align-items: flex-start;
  }

  .name {
    font-size: 12px;
  }

  .time {
    font-size: 10px;
  }

  .message-preview {
    font-size: 11px;
  }
}
</style>
