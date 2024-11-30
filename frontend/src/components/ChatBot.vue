<template>
  <div class="app-container">
    <aside class="sidebar" :class="{ hidden: isSidebarHidden }" ref="sidebar">
      <div ref="sidebarContent">
        <div class="sidebar-header">
          <div class="small_logo"></div>
          <div class="title">Conversations</div>
        </div>
        <div class="conversation-list" ref="conversationList">
          <!-- Les conversations précédentes seront listées ici -->
        </div>
      </div>
    </aside>
    <main class="chat-container" :class="{ 'full-width': isSidebarHidden }">
      <div class="chat-header">
        <button class="toggle-button" @click="toggleSidebar">☰</button>
        <div class="small-logo"></div>
        <div class="title">Mori Chatbot</div>
      </div>
      <div
        class="chat-window"
        ref="chatWindow"
        :class="{ 'full-width': isSidebarHidden }"
      >
        <!-- Les messages du chatbot et de l'utilisateur apparaîtront ici -->
        <div
          v-for="(message, index) in messages"
          :key="index"
          :class="['message', message.sender]"
        >
          <pre>{{ message.content }}</pre>
        </div>
      </div>
      <div class="chat-input" :class="{ 'full-width': isSidebarHidden }">
        <textarea
          style="resize: none"
          cols="71"
          rows="2"
          v-model="userInput"
          placeholder="Tapez votre message..."
          @keydown="handleKeydown"
        ></textarea>
        <button @click="sendMessage">Envoyer</button>
      </div>
    </main>
  </div>
</template>

<script>
export default {
  data() {
    return {
      userInput: "",
      messages: [
        {
          sender: "LLM",
          content: "Bienvenue chez Mori",
        },
      ],
      isSidebarHidden: false,
      conversation: {
        user_id: "",
        conversation_id: "",
        user_request: "",
        llm_response: "",
        session: "",
        new_conversation: false,
        history: [],
      },
    };
  },
  methods: {
    toggleSidebar() {
      this.isSidebarHidden = !this.isSidebarHidden;
      this.$nextTick(() => {
        this.scrollToBottom();
      });
    },
    handleKeydown(event) {
      if (event.key === "Enter") {
        if (event.shiftKey) {
          // Permettre le retour à la ligne avec Shift+Enter
        } else {
          event.preventDefault();
          this.sendMessage();
        }
      }
    },
    sendMessage() {
      const message = this.userInput.trim();
      if (message === "") return;

      // Ajouter le message de l'utilisateur
      this.messages.push({ sender: "Utilisateur", content: message });
      this.userInput = "";
      this.scrollToBottom();

      // Envoyer le message au serveur
      this.sendToServer(message);
    },
    async sendToServer(message) {
      this.conversation.user_id = "0001";
      this.conversation.user_request = message;
      this.conversation.conversation_id = window.location.pathname
        .split("/")
        .pop();
      this.conversation.session = window.location.pathname.split("/")[2];

      try {
        const response = await fetch(`${window.location.pathname}`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(this.conversation),
        });

        if (!response.ok || !response.body) {
          console.error(
            "Erreur lors de l'envoi des données :",
            response.statusText
          );
          return;
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder("utf-8");

        const messageElement = { sender: "LLM", content: "" };
        this.messages.push(messageElement);

        let done = false;
        while (!done) {
          const { value, done: readerDone } = await reader.read();
          done = readerDone;

          const chunk = decoder.decode(value, { stream: true });

          const lines = chunk.split("\n");
          for (const line of lines) {
            if (line.startsWith("data: ")) {
              let jsonData = line.replace("data: ", "").trim();
              try {
                let parsedData = JSON.parse(jsonData);
                const responseText = parsedData.response.message.content;
                messageElement.content += responseText;
              } catch (e) {
                console.error("Erreur de parsing JSON :", e);
              }
            }
          }

          // Forcer la mise à jour de Vue
          this.$forceUpdate();
          this.scrollToBottom();
        }

        this.conversation.llm_response = messageElement.content;
        this.sendConversation();
      } catch (error) {
        console.error("Erreur de lecture du flux", error);
      }
    },
    async sendConversation() {
      try {
        const response = await fetch(
          `http://localhost:8080/${window.location.pathname}`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(this.conversation),
          }
        );

        if (!response.ok) {
          console.error(
            "Erreur lors de l'envoi des données :",
            response.statusText
          );
        } else {
          console.log("Conversation envoyée avec succès !");
        }
      } catch (error) {
        console.error("Erreur lors de l'envoi de la conversation", error);
      }
    },
    scrollToBottom() {
      const chatWindow = this.$refs.chatWindow;
      chatWindow.scrollTop = chatWindow.scrollHeight;
    },
  },
  mounted() {
    this.scrollToBottom();
  },
};
</script>

<style>
@font-face {
  font-family: "TITLE";
  src: url(../assets/fonts/SinosansRegular-aYxZ5.otf) format("truetype");
  font-weight: normal;
  font-style: normal;
}

body {
  font-family: "Arial", sans-serif;
  margin: 0;
  background-color: #1a1a1a;
  color: #e4e4e4;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.app-container {
  display: flex;
  width: 100%;
  height: 100vh;
}

.sidebar {
  width: 30%;
  background: #2c2c2c;
  display: flex;
  flex-direction: column;
  padding: 10px;
  overflow-y: auto;
  animation: appearFromLeft 0.5s ease;
  transition: transform 0.3s ease;
  animation: grow 0.5s ease;
}

.sidebar.hidden {
  animation: none;
  display: none;
  transform: translateX(-100%);
  transition: transform 0.3s ease;
}

#sidebar_content {
  display: flex;
  width: 100%;
  height: 100%;
  flex-direction: column;
  opacity: 0;
  animation: opacity 0.5s ease 0.5s forwards;
}

.sidebar-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.small_logo {
  width: 50px;
  height: 50px;
  background-image: url("../assets/images/mori.png");
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
  border-radius: 50%;
  margin-right: 10px;
}

.small-logo.full-width {
  display: block;
  width: 50px;
  height: 50px;
  background-image: url("../assets/images/mori.png");
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
  border-radius: 50%;
  margin-right: 10px;
  transform: translateX(0);
  transition: transform 0.5s ease;
}

.title {
  font-size: 24px;
  font-family: "arial";
  font-weight: bold;
  color: #e2e2e2;
}

.conversation-list {
  flex-grow: 1;
}

.conversation-item {
  padding: 10px;
  background: #3a3a3a;
  margin-bottom: 10px;
  border-radius: 8px;
  cursor: pointer;
}

.conversation-item:hover {
  background: #4a4a4a;
}

.chat-container {
  width: 70%;
  background: #1a1a1a;
  display: flex;
  flex-direction: column;
  margin-left: 0px;
  animation: shrink 0.5s ease;
}

.chat-container.full-width {
  width: 100%;
  transform: translateX(0);
  transition: transform 0.5s ease;
  animation: none;
}

.chat-header {
  display: flex;
  align-items: center;
  padding: 10px;
  background: #9146bc;
  height: 55px;
}

.toggle-button {
  background: none;
  border: none;
  color: #e4e4e4;
  font-size: 24px;
  cursor: pointer;
  margin-right: 10px;
}

.chat-window {
  display: flex;
  flex-direction: column;
  gap: 5px;
  overflow-y: auto;
  padding: 10px;
  height: 400px;
  border: none;
}

.chat-window.full-width {
  display: flex;
  flex-direction: column;
  padding: 10px;
  height: 400px;
  width: 70%;
  margin: 0 auto;
  justify-content: center;
  align-items: center;
  transform: translateX(0);
  transition: width 0.3s ease;
}

.chat-message {
  padding: 8px;
  margin: 5px 0;
  border-radius: 5px;
  max-width: 80%;
  color: #ededed;
}

.message {
  padding: 8px;
  margin: 5px 0;
  border-radius: 5px;
  max-width: 80%;
  color: #ededed;
}

.message.Utilisateur {
  background: #9146bc;
  color: #fff;
  align-self: flex-end;
  text-align: left;
}

.message.LLM {
  background: #444;
  color: #e4e4e4;
  align-self: flex-start;
  text-align: left;
}

.chat-input {
  display: flex;
  padding: 10px;
  width: 55%;
  margin: 0 auto;
}

#welcome_message {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  font-size: 24px;
  font-family: Arial;
}

#welcome_message.hidden {
  display: none;
}

textarea {
  flex-grow: 1;
  padding: 10px;
  border: none;
  border-radius: 8px;
  background: #444;
  color: #e4e4e4;
  font-size: 14px;
  font-family: Arial, Helvetica, sans-serif;
}

textarea:focus {
  outline: 2px solid #9146bc;
}

.chat-input.full-width {
  width: 55%;
  margin: 0 auto;
  transform: translateX(0);
  transition: width 0.3s ease;
}

#send-btn {
  padding: 10px 20px;
  margin-left: 10px;
  border: none;
  border-radius: 8px;
  background: #9146bc;
  color: #fff;
  font-size: 16px;
  cursor: pointer;
}

#send-btn:hover {
  background: #7d3aa6;
}

pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}

@media screen and (max-width: 768px) {
  .app-container {
    flex-direction: column;
  }
  .sidebar {
    width: 100%;
    height: 30%;
    transform: translateY(-100%);
  }
  .sidebar.hidden {
    transform: translateY(0);
  }
  .chat-container {
    width: 100%;
    height: 70%;
  }
}

@keyframes appearFromLeft {
  from {
    transform: translateX(-100%);
  }
  to {
    transform: translateX(0);
  }
}

@keyframes shrink {
  from {
    width: 100%;
  }
  to {
    width: 70%;
  }
}

@keyframes grow {
  from {
    width: 0%;
  }
  to {
    width: 30%;
  }
}

@keyframes opacity {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>
