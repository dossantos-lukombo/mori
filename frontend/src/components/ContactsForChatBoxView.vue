<template>
    <div class="contacts-wrapper">
      <h2 class="titre">Contacts</h2>
      <ul class="contacts-list">
        <h3 class="sous_titres">Amis</h3>
        <li v-for="contact in chatUserList" :key="contact.id" @click="selectContact(contact, 'PERSON')"> 
          <div class="user-picture medium"
            :style="{ backgroundImage: `url(http://localhost:8081/${contact.avatar})` }"></div>        
          <div class="contact-item">
            <div class="contact-name">{{ contact.nickname }}</div>
            <span
              v-if="totalUnreadMessagesCount(contact.id, 'PERSON') !== 0"
              class="unread-messages"
            >
              {{ totalUnreadMessagesCount(contact.id, 'PERSON') }}
            </span>
          </div>
        </li>
        <h3 class="sous_titres">Groupes</h3>
        <li v-for="group in userGroups" :key="group.id" @click="selectContact(group, 'GROUP')">          
          <div class="contact-item">
            <div class="contact-name">{{ group.name }}</div>
            <span
              v-if="totalUnreadMessagesCount(group.id, 'GROUP') !== 0"
              class="unread-messages"
            >
              {{ totalUnreadMessagesCount(group.id, 'GROUP') }}
            </span>
          </div>
        </li>
      </ul>
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
      }),
      ...mapGetters(["getUnreadMessagesCount", "getUnreadGroupMessagesCount", "getUnreadMsgsCountFromDB"]),
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
          return this.getUnreadMessagesCount(receiverId) + this.getUnreadMsgsCountFromDB(receiverId);
        } else {
          return this.getUnreadGroupMessagesCount(receiverId) + this.getUnreadMsgsCountFromDB(receiverId);
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .contacts-wrapper {
    padding: 20px;
    background-color: var(--bg-neutral);
  }
  
  .contacts-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  
  .contact-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px;
    border-bottom: 1px solid var(--color-grey-light);
    cursor: pointer;
  }
  
  .contact-item:hover {
    background-color: var(--hover-color);
  }
  
  .unread-messages {
    background-color: var(--purple-color);
    color: var(--color-white);
    padding: 5px;
    border-radius: 10px;
    font-size: 12px;
  }.titre {
    color: white;
    font-size: 1.8em;
  }

  .sous_titres {
    color: white;
    font-size: 1.2em;
    margin-top: 2vh;
  }
  </style>
  