<template>
  <!-- <h3>Historique du Chatbot</h3> -->
  <ul id="list_chat_convo">
    <li v-for="(convo, index) in conversationsUpdate" :key="index">
      <!-- <div @click="loadConvo(convo[0].conversation_id)"></div> -->
      <div class="elmt_history">
        {{ convo[convo.length - 1].user_request }}
      </div>
      <button
        class="btn_delete_convo"
        @click="deleteConvo(convo[index].conversation_id, index)"
      >
        X
      </button>
    </li>
  </ul>
</template>

<script>
// import ChatbotConversation from "./ChatbotConversation.vue";

export default {
  data() {
    return {
      userInput: "",
      chatHistory: [],
      currentChat: { messages: [] },
    };
  },
  computed: {
    conversationsUpdate() {
      console.log(
        "allConversations.length: ",
        this.$store.getters.allConversations.length
      );
      if (this.$store.getters.allConversations.length === 1) {
        this.chatHistory.push(this.$store.getters.allConversations);
        this.chatHistory.reverse();
        this.$store.dispatch("clearChatHistory");
      }
      console.log("conversationsUpdate: ", this.chatHistory);
      return this.chatHistory;
    },
  },
  mounted() {
    this.getChatHistory();
  },
  methods: {
    async getMyUserID() {
      const response = await fetch("http://localhost:8081/currentUser", {
        credentials: "include",
        headers: new Headers({
          "Content-Type": "application/json",
        }),
        method: "POST",
      });
      if (!response.ok) {
        console.error(
          "Erreur lors de la récupération de l'ID de l'utilisateur :",
          response.statusText
        );
        return;
      } else {
        const resp = await response.json();
        console.log(resp);
        console.log(resp.users[0].id);
        return resp.users[0].id;
      }
    },

    async deleteConvo(conversation_id, index) {
      console.log("conversation_id FOR DELETING: ", conversation_id);
      const response = await fetch("http://localhost:8081/llmConvoDelete", {
        credentials: "include",
        headers: new Headers({
          "Content-Type": "application/json",
        }),
        method: "POST",
        body: JSON.stringify({
          user_id: await this.getMyUserID(),
          conversation_id: conversation_id,
        }),
      });
      if (!response.ok) {
        console.error(
          "Erreur lors de la suppression de la conversation de l'utilisateur :",
          response.statusText
        );
        return;
      } else {
        this.$store.dispatch("deleteConversation", index);
        this.$store.dispatch("clearMessages");
        this.chatHistory.splice(index, 1);
        console.log("Conversation supprimée");
      }
    },

    //Méthode pour obtenir la discussion selctionnée
    async loadConvo(convo_id) {
      const response = await fetch("http://localhost:8081/llmConvoSelected", {
        credentials: "include",
        headers: new Headers({
          "Content-Type": "application/json",
        }),

        method: "POST",
        body: JSON.stringify({
          user_id: await this.getMyUserID(),
          conversation_id: convo_id,
        }),
      });
      if (!response.ok) {
        console.error(
          "Erreur lors de la récupération de la conversation de l'utilisateur :",
          response.statusText
        );
        return;
      } else {
        const resp = await response.json();
        console.log("CONVERSATION: ", resp);
        this.currentChat = resp;
      }
    },

    //Méthode pour obtenir les conversations de l'utilisateur
    async getChatHistory() {
      const response = await fetch("http://localhost:8081/llmConvoGet", {
        credentials: "include",
        headers: new Headers({
          "Content-Type": "application/json",
        }),
        method: "POST",
        body: JSON.stringify({ user_id: await this.getMyUserID() }),
      });
      if (!response.ok) {
        console.error(
          "Erreur lors de la récupération des conversations de l'utilisateur :",
          response.statusText
        );
        return;
      }
      const resp = await response.json();

      this.$store.dispatch("clearChatHistory"); //clear the chat history variable
      let same_convo = [];
      // this.chatHistory = this.$store.getters.allConversations; //get the list variable of all conversations
      for (let i = 0; i < resp.length; i++) {
        if (
          i < resp.length - 1 &&
          resp[i].conversation_id === resp[i + 1].conversation_id
        ) {
          if (resp[i].llm_response !== "") {
            same_convo.push(resp[i]);
          }
        } else if (
          i < resp.length - 1 &&
          resp[i].conversation_id !== resp[i + 1].conversation_id
        ) {
          if (
            i > 0 &&
            resp[i].conversation_id === resp[i - 1].conversation_id &&
            resp[i].llm_response !== ""
          ) {
            same_convo.push(resp[i]);
          }
        }
      }
      this.chatHistory.push(same_convo);
      for (let i = 0; i < resp.length; i++) {
        if (!same_convo.includes(resp[i]) && resp[i] !== null) {
          this.chatHistory.push([resp[i]]);
        }
      }
      this.chatHistory = this.chatHistory.filter((convo) => convo.length > 0);
      this.chatHistory.reverse();
    },
  },
};
</script>

<style>
#app {
  display: flex;
  height: 100vh;
}

.sidebar {
  width: 25%;
  background: #f4f4f4;
  border-right: 1px solid #ddd;
  padding: 10px;
}

.sidebar ul {
  list-style: none;
  padding: 0;
}

.sidebar li {
  cursor: pointer;
  padding: 5px;
  border-bottom: 1px solid #ddd;
}

.sidebar li:hover {
  background: #eaeaea;
}

.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 10px;
}

.chat-content {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 10px;
}

.user-message {
  text-align: right;
  background: #d1ffd1;
  margin: 5px;
  padding: 5px 10px;
  border-radius: 5px;
}

.bot-message {
  text-align: left;
  background: #f0f0f0;
  margin: 5px;
  padding: 5px 10px;
  border-radius: 5px;
}

.elmt_history {
  padding: 10px;
  border: 1px solid #ddd;
  margin: 5px;
  border-radius: 5px;
  background-color: black;
  color: white;
}

.elmt_history:hover {
  background-color: none;
  border-color: none;
}

input[type="text"] {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  outline: none;
}
</style>
