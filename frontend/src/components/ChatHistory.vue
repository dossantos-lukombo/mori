<template>
  <!-- <h3>Historique du Chatbot</h3> -->
  <button class="btn_new_convo" @click="newConversation">
    New Conversation
  </button>
  <ul id="list_chat_convo">
    <li v-for="(convo, index) in conversationsUpdate">
      <div @click="loadConvo(convo.conversation_id)">
        <div class="elmt_history">
          {{ convo.convo[convo.convo.length - 1].user_request }}
        </div>
        <button
          class="btn_delete_convo"
          @click="deleteConvo(convo.conversation_id, convo.convo.length - 1)"
        >
          X
        </button>
      </div>
    </li>
  </ul>
</template>

<script>
export default {
  data() {
    return {
      userInput: "",
      chatHistory: [],
    };
  },
  computed: {
    conversationsUpdate() {
      console.log(
        "store allConversations: ",
        this.$store.getters.allConversations
      );
      this.chatHistory = [];
      this.$store.getters.allConversations.forEach((convo) => {
        console.log("c: ", convo);
        this.chatHistory.push(convo);
      });

      return this.chatHistory.reverse();
    },
  },
  mounted() {
    this.getHistory();
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

    newConversation() {
      this.$store.dispatch("clearMessages");
    },

    async deleteConvo(conversation_id, index) {
      console.log("conversation_id FOR DELETING: ", conversation_id);
      console.log("index FOR DELETING: ", index);

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
        // this.chatHistory.splice(index, 1);
        const newArray = [
          ...this.chatHistory.slice(0, index),
          ...this.chatHistory.slice(index + 1)
        ];
        this.chatHistory = newArray.reverse();
        console.log("chatHistory after deleting: ", this.chatHistory);
        console.log("History into store: ", this.$store.getters.allConversations);
        console.log("Conversation supprimée");
      }
    },

    //Méthode pour obtenir la discussion selctionnée
    async loadConvo(convo_id) {
      console.log("convo_id in loadConvo: ", convo_id);
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
        console.log("Current convo: ", resp.convo);
        console.log("Current convo_id loadConvo: ", resp.conversation_id);
        this.$store.dispatch("clearMessages");
        this.convertMessages(resp);
      }
    },

    convertMessages(convo) {
      console.log("convo in convertMessages: ", convo);

      if (convo !== null) {
        console.log("convo is null ");
        return;
      }
      convo.convo.reverse().forEach((message) => {
        this.$store.dispatch("addMessage", {
          sender: "Utilisateur",
          text: message.user_request,
          conversation_id: convo.conversation_id,
        });
        this.$store.dispatch("addMessage", {
          sender: "LLM",
          text: message.llm_response,
          conversation_id: convo.conversation_id,
        });
      });
    },
    convertLastMessages(convo, conversation_id) {
      console.log("convo in convertLastMessages: ", convo);

      convo.forEach((c) => {
        if (c.conversation_id === conversation_id) {
          c.convo.forEach((message) => {
            this.$store.dispatch("addMessage", {
              sender: "Utilisateur",
              text: message.user_request,
              conversation_id: message.conversation_id,
            });
            this.$store.dispatch("addMessage", {
              sender: "LLM",
              text: message.llm_response,
              conversation_id: message.conversation_id,
            });
          });
        }
      });
    },

    async getHistory() {
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
      console.log("RESPONSE getHistory: ", resp);
      this.$store.dispatch("clearMessages");

      this.$store.dispatch("clearChatHistory");

      // this.convertLastMessages(resp, localStorage.getItem("current_convo_id"));
      this.chatHistory = [];
      resp.forEach((convo) => {
        this.$store.dispatch("addConversation", convo);
        this.chatHistory.push(convo);
      });
      this.chatHistory.reverse();
      console.log("chatHistory: ", this.chatHistory);
    },

    //Méthode pour obtenir les conversations de l'utilisateur
    // async getChatHistory() {
    //   const response = await fetch("http://localhost:8081/llmConvoGet", {
    //     credentials: "include",
    //     headers: new Headers({
    //       "Content-Type": "application/json",
    //     }),
    //     method: "POST",
    //     body: JSON.stringify({ user_id: await this.getMyUserID() }),
    //   });
    //   if (!response.ok) {
    //     console.error(
    //       "Erreur lors de la récupération des conversations de l'utilisateur :",
    //       response.statusText
    //     );
    //     return;
    //   }
    //   const resp = await response.json();

    //   // this.$store.dispatch("clearChatHistory"); //clear the chat history variable
    //   this.$store.dispatch("clearMessages"); //clear the messages variable
    //   // this.chatHistory = [];
    //   console.log("chatHistory before Reload: ", this.chatHistory);
    //   console.log("RESPONSE getChatHistory: ", resp);
    //   let same_convo = [];
    //   let conv_ids = [];
    //   console.log("same_convo: ", same_convo);

    //   for (let i = 0; i < resp.length; i++) {
    //     for (let j = i; j < resp.length; j++) {
    //       if (resp[i].conversation_id === resp[j].conversation_id) {
    //         if (!same_convo.includes(resp[j])) {
    //           same_convo.push(resp[j]);
    //         } else {
    //           this.chatHistory.push(same_convo);
    //           if (!conv_ids.includes(resp[j].conversation_id)) {
    //             conv_ids.push(resp[j].conversation_id);
    //           }
    //           same_convo = [];
    //           break;
    //         }
    //       }
    //     }
    //   }

    //   // this.chatHistory = this.chatHistory.filter((convo) => convo.length > 1);
    //   console.log("conversation ids: ", conv_ids);
    //   for (let j = 0; j < resp.length; j++) {
    //     if (!conv_ids.includes(resp[j].conversation_id)) {
    //       this.chatHistory.push([resp[j]]);
    //     }
    //   }

    //   // this.chatHistory = this.chatHistory.filter((convo) => convo.length > 0);
    //   this.chatHistory.reverse();
    //   console.log("chatHistory after Reload: ", this.chatHistory);
    // },
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
