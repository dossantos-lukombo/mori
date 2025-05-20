<template>
  <div class="chatbot-container">
    <div :class="['chatbot-box', { 'chatbot-box--active': hasMessages }]">
      <div class="mori-img" v-if="!hasMessages">
        <div class="moriImg"></div>
      </div>
      <div class="mori" id="moriChatBot" v-if="!hasMessages">Mori</div>
      <div class="chatbot-message" v-if="!hasMessages">How can I help you?</div>
      <div class="chatbot-messages" v-if="hasMessages">
        <div
          v-for="(message, index) in allMessages"
          :key="index"
          :class="[
            'message',
            message.sender === 'Utilisateur' ? 'Utilisateur' : 'LLM',
          ]"
        >
          <div v-if="message.sender === 'LLM'" class="bot-logo">
            <img src="../assets/mori.png" alt="Bot Logo" />
          </div>
          <div
            v-if="message.sender === 'Utilisateur'"
            class="markdown-container-Utilisateur"
          >
            <Markdown :source="message.text" />
          </div>
          <div v-if="message.sender === 'LLM'" class="markdown-container-LLM">
            <Markdown class="markdownLLM" :source="message.text" />
          </div>
          <div class="timestamp">{{ message.timestamp }}</div>
        </div>
        <!-- Loading indicator -->
        <div v-if="isGenerating && isLoading" class="message LLM loading-message">
          <div class="loading-container">
            <div class="loading-text">Mori is thinking</div>
            <div class="loading-indicator">
              <div class="loading-dot"></div>
              <div class="loading-dot"></div>
              <div class="loading-dot"></div>
              <div class="loading-dot"></div>
              <div class="loading-dot"></div>
            </div>
          </div>
        </div>
      </div>

      <div
        :class="[
          'chatbot-input-container',
          { 'chatbot-input-container--active': hasMessages },
        ]"
      >
        <textarea
          ref="textarea"
          :rows="rows"
          class="chatbot-textarea"
          v-model="userInput"
          @keydown="handleKeydown"
          @click="handleKeydown"
          placeholder="Type your message here..."
          :disabled="isGenerating"
        ></textarea>
        <button 
          @click="isGenerating ? stopGeneration() : sendMessage()" 
          :class="['send-button', { 'stop-button': isGenerating }]"
          :disabled="!userInput.trim() && !isGenerating"
        >
          <template v-if="!isGenerating">
            <span>Send</span>
          </template>
          <template v-else>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="6" y="6" width="12" height="12" rx="2" fill="currentColor"/>
            </svg>
            <span>Stop</span>
          </template>
        </button>
        
        <!-- Kickstart suggestions (now inside the input container) -->
        <div v-if="!hasMessages" class="kickstart-container">
          <h3 class="kickstart-title">Try asking Mori about:</h3>
          <div class="kickstart-suggestions">
            <button 
              v-for="(suggestion, index) in kickstartSuggestions" 
              :key="index" 
              class="kickstart-button"
              @click="useKickstartSuggestion(suggestion)"
            >
              <div class="suggestion-icon">
                <svg v-if="index === 0" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z" fill="currentColor"/>
                </svg>
                <svg v-if="index === 1" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z" fill="currentColor"/>
                </svg>
                <svg v-if="index === 2" width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V5h14v14zm-7-2h2V7h-4v2h2z" fill="currentColor"/>
                </svg>
              </div>
              <span class="suggestion-text">{{ suggestion }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Markdown from "vue3-markdown-it";
const { v4: uuidv4 } = require("uuid");

export default {
  components: { Markdown },

  data() {
    return {
      userInput: "",
      messages: [],
      current_convID: "",
      rows: 10,
      sourceLLM: "",
      sourceUtilisateur: "",
      markdownText: "",
      conversation: {
        user_id: "",
        conversation_id: "",
        convo: [],
        new_conversation: false,
      },
      isLoading: false,
      abortController: null,
      isGenerating: false,
      kickstartSuggestions: [
        "Tell me about yourself and how you can help me",
        "What can you do for me as an AI assistant?",
        "Give me some ideas for my next project"
      ]
    };
  },
  computed: {
    hasMessages() {
      return this.$store.getters.allMessages.length > 0;
    },
    allMessages() {
      this.messages = this.$store.getters.allMessages;
      return this.messages;
    },
  },
  mounted() {
    // this.loadCurrentConvo();
  },
  methods: {
    //Méthode pour récupérer les messages de la conversation selectionné
    getCurrentMessages() {
      this.messages = this.$store.getters.allMessages;
    },
    async loadCurrentConvo() {
      const convo_id = localStorage.getItem("current_convo_id");
      console.log("convo_id in loadConvo: ", convo_id);
    
      if (convo_id === null) {
        console.log("No conversation selected");
        return;
      }
    
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
        this.$store.dispatch("clearMessages");
        this.convertMessages(resp);
      }
    },

    //Méthode pour récupérer l'ID de l'utilisateur
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
    appendMessage(sender, text) {
      let dict = {};
      dict = {
        sender,
        text,
        conversation_id: "",
      };

      this.$store.dispatch("addMessage", dict);

      this.$nextTick(() => {
        const chatBox = this.$el.querySelector(".chatbot-messages");
        chatBox.scrollTop = chatBox.scrollHeight;
      });
    },
    async sendData() {
      let accumulatedText = "";
      let hasStartedGenerating = false;
      
      try {
        const response = await fetch(`http://localhost:8081/llmConvo`, {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(this.conversation),
          signal: this.abortController.signal
        });

        if (!response.ok) {
          throw new Error(response.statusText);
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder("utf-8");
        this.appendMessage("LLM", "");
        let lastLLMMessage = this.$store.getters.allMessages[this.$store.getters.allMessages.length - 1];

        accumulatedText = "";
        let { done, value } = await reader.read();

        while (!done) {
          const chunk = decoder.decode(value, { stream: true });
          const lines = chunk.split("\n");
          
          for (const line of lines) {
            if (line.startsWith("data: ")) {
              const jsonData = line.replace("data: ", "").trim();
              try {
                const parsedData = JSON.parse(jsonData);
                if (parsedData.response) {
                  if (!hasStartedGenerating) {
                    hasStartedGenerating = true;
                    this.isLoading = false; // Clear loading state when we start receiving the response
                  }
                  accumulatedText += parsedData.response;
                  lastLLMMessage.text = accumulatedText;
                  // Force Vue to update the view
                  this.$forceUpdate();
                }
              } catch (error) {
                console.error("Erreur de parsing JSON :", error);
              }
            }
          }
          
          ({ done, value } = await reader.read());
        }

        // Update the final message
        if (accumulatedText) {
          lastLLMMessage.text = accumulatedText;
          this.$store.dispatch("updateLastMessage", {
            index: this.$store.getters.allMessages.length - 1,
            text: accumulatedText
          });
        }

        this.conversation.user_id = await this.getMyUserID();
        this.conversation.convo[this.conversation.convo.length - 1].llm_response = accumulatedText;

        if (this.messages.length <= 2) {
          this.conversation.conversation_id = uuidv4();
          this.messages.forEach((message) => {
            message.conversation_id = this.conversation.conversation_id;
          });
          this.conversation.new_conversation = true;
          this.$store.dispatch("addConversation", this.conversation);
        } else {
          this.conversation.conversation_id = this.messages[0].conversation_id;
          this.conversation.new_conversation = false;
          this.addMessageToExistingConversation(
            this.$store.getters.allConversations,
            this.conversation
          );
        }

        await this.sendConversation();
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error; // Re-throw abort errors to be handled by sendMessage
        }
        console.error("Error in sendData:", error);
      } finally {
        this.isLoading = false;
        this.isGenerating = false;
        this.abortController = null;
      }
    },
    async sendConversation() {
      const response = await fetch(`http://localhost:8081/llmConvoSave`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(this.conversation),
      });

      if (!response.ok) {
        console.error(
          "Erreur lors de l'envoi de la conversation :",
          response.statusText
        );
        return;
      }
      console.log("Conversation envoyée avec succès");
      this.conversation = {
        user_id: "",
        conversation_id: "",
        convo: [],
        new_conversation: false,
      };
    },
    //Méthode pour gérer les messages dans une conversation existante
    addMessageToExistingConversation(allConversation, currentConversation) {
      console.log(
        "Conversation in addMessageToExistingConversation: ",
        allConversation
      );
      console.log("Current conversation: ", currentConversation);
      console.log(
        "Current conversation ID: ",
        currentConversation.conversation_id
      );

      for (let t = 0; t < allConversation.length; t++) {
        const convo = allConversation[t];
        if (convo.conversation_id === currentConversation.conversation_id) {
          console.log("Conversation found");
          console.log("Current conversation: ", currentConversation.convo);
          currentConversation.convo.push(...convo.convo);
        }
      }
      console.log("Current conversation after: ", currentConversation.convo);
    },

    // Méthode pour gérer les événements de touche
    handleKeydown(event) {
      if (event.shiftKey && event.key === "Enter") {
        event.preventDefault();
        this.userInput += "\n";
        let textarea = this.$el.querySelector("textarea");
        textarea.style.height = `${textarea.scrollHeight + 10}px`;
      } else if (event.key === "Enter") {
        event.preventDefault();
        this.sendMessage();
      }
      // let textarea = this.$el.querySelector("textarea");
      let textarea = document.querySelector("textarea");
      const textLength = textarea.value.length;
      if (event.key === "Backspace" && textLength > 0) {
        const cursorPosition = textarea.selectionEnd; // Position actuelle du curseur

        // Vérifie si le caractère à supprimer est un retour chariot
        if (textarea.value[cursorPosition - 1] === "\n") {
          // Réduit la hauteur du textarea
          textarea.style.height = `${textarea.scrollHeight - 22}px`;
        }
      } else if (event.key === "Backspace" && textLength === 1) {
        textarea.style.height = `50px`;
      }
    },
    convertMessages(convo) {
      console.log("convo in convertMessages: ", convo);

      convo.convo.forEach((message) => {
        this.messages.push({
          sender: "Utilisateur",
          text: message.user_request,
          conversation_id: convo.conversation_id,
        });
        this.messages.push({
          sender: "LLM",
          text: message.llm_response,
          conversation_id: convo.conversation_id,
        });
      });
    },
    async sendMessage() {
      if (this.userInput.trim() === "" || this.isGenerating) return;
      
      this.isLoading = true;
      this.isGenerating = true;
      this.abortController = new AbortController();
      
      const userMessage = this.userInput;
      this.userInput = ""; // Clear input immediately
      
      this.appendMessage("Utilisateur", userMessage);
      this.conversation.convo.push({
        user_request: userMessage,
        llm_response: "",
      });

      try {
        await this.sendData();
      } catch (error) {
        if (error.name === 'AbortError') {
          console.log('Generation stopped by user');
          this.appendMessage("LLM", "[Generation stopped]");
        } else {
          console.error("Erreur lors de l'envoi du message :", error);
          this.appendMessage("LLM", "Sorry, an error occurred while generating the response.");
        }
      } finally {
        this.isLoading = false;
        this.isGenerating = false;
      }
    },
    stopGeneration() {
      if (this.abortController) {
        this.abortController.abort();
        this.isLoading = false;
        this.isGenerating = false;
        this.abortController = null;
      }
    },
    useKickstartSuggestion(suggestion) {
      this.userInput = suggestion;
      // Focus the textarea
      this.$nextTick(() => {
        this.$refs.textarea.focus();
      });
    },
  },
};
</script>

<style scoped>
  .Utilisateur {
    align-self: flex-end;
    background-color: var(--purple-color);
    color: var(--color-white);
  }
  
  .LLM {
    align-self: flex-start;
    text-align: left;
    background-color: var(--page-bg);
    color: var(--color-white);
  }
  
  .chatbot-container {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
    background-color: var(--page-bg);
    font-family: Arial, sans-serif;
  }
  
  .chatbot-box {
    display: flex;
    flex-direction: column;
    width: 100%;
    max-width: 800px;
    border-radius: 10px;
    padding: 30px;
    text-align: center;
    gap: 20px;
    margin-bottom: 205px;
    transition: all 0.5s ease;
  }
  
  .chatbot-box--active {
    justify-content: space-between;
    width: 90%;
    height: 85vh;
    padding: 2rem;
    margin-bottom: 100px;
  }

  .bot-logo {
    display: inline-block;
    vertical-align: top;
    margin-right: 10px;
    margin-top: -5px;
    margin-left: -5px;
  }
  
  .bot-logo img {
    width: 35px; /* Adjust size as needed */
    height: 35px; /* Adjust size as needed */
    border-radius: 50%; /* Optional: Make the image circular */
    background-color: var(--purple-color);
    object-fit: cover; /* Ensure the image scales properly */
  }
  
  
  .mori-img {
    display: flex;
    justify-content: center;
    align-items: center;
    transition: opacity 0.5s ease;
  }
  
  #moriChatBot {
    user-select: none;
    font-size: 50px;
    font-weight: bold;
    transition: opacity 0.5s ease;
  }
  
  .chatbot-message {
    user-select: none;
    margin-bottom: 20px;
    font-size: 20px;
    color: var(--color-white);
    transition: opacity 0.5s ease;
  }
  
  .chatbot-messages {
    flex: 1;
    overflow-y: auto;
    padding: 10px;
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  
  .message {
    max-width: 70%;
    padding: 10px;
    border-radius: 10px;
    font-size: 16px;
    position: relative;
  }
  
  .user {
    align-self: flex-end;
    background-color: var(--purple-color);
    color: var(--color-white);
  }
  
  .bot {
    align-self: flex-start;
    background-color: var(--bg-neutral);
    color: var(--color-white);
  }
  
  .timestamp {
    font-size: 12px;
    color: var(--color-grey);
    opacity: 0.5;
    text-align: right;
    margin-top: 5px;
  }
  
  /* Input field animation */
  /* 
  1. The container that slides down with an animation 
     (replaces .chatbot-input in your old code)
*/
.chatbot-input-container {
  display: flex;
  gap: 10px;
  align-items: center;
  position: absolute;
  top: 58%; /* Initially below the greeting message */
  left: 50%;
  transform: translate(-50%, -50%);
  width: calc(40% - 40px);
  border-radius: 10px;
  padding: 10px 20px;
  transition: top 0.7s ease, transform 0.7s ease, width 0.7s ease;
}

.chatbot-input-container--active {
  width: calc(50% - 40px); /* Widen the container */
  position: absolute;
  top: calc(97% - 80px);   /* Slide to bottom of viewport */
  transform: translateX(-50%);
}

/* 
  2. The textarea itself: 
     (new .chatbot-textarea class)
*/
.chatbot-textarea {
  flex: 1;
  border: 1px solid var(--color-grey);
  border-radius: 10px;
  font-size: 16px;
  min-height: 50px;   /* Ensure it matches your old input height */
  padding: 13px;
  resize: none;       /* Optional: remove manual resize handle */
  transition: all 0.3s ease;
}

/* 
  3. The Send button 
  (same rules as your old .chatbot-input button style)
*/
.chatbot-input-container button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 24px;
  background-color: var(--purple-color);
  color: var(--color-white);
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.3s ease;
  min-width: 100px;
}

.chatbot-input-container button:hover:not(:disabled) {
  background-color: var(--hover-background-color);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(var(--purple-color-rgb), 0.2);
}

.chatbot-input-container button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.chatbot-input-container button.stop-button {
  background-color: #ff4444;
  animation: pulse 2s infinite;
}

.chatbot-input-container button.stop-button:hover:not(:disabled) {
  background-color: #ff3333;
  box-shadow: 0 4px 12px rgba(255, 68, 68, 0.2);
}

.chatbot-input-container button svg {
  width: 16px;
  height: 16px;
  transition: transform 0.2s ease;
}

.chatbot-input-container button:hover:not(:disabled) svg {
  transform: scale(1.1);
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(255, 68, 68, 0.4);
  }
  70% {
    box-shadow: 0 0 0 10px rgba(255, 68, 68, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(255, 68, 68, 0);
  }
}

.loading-message {
  opacity: 0.9;
  animation: fadeIn 0.3s ease-in-out;
  background: linear-gradient(135deg, rgba(var(--purple-color-rgb), 0.1), rgba(var(--purple-color-rgb), 0.05));
  border: 1px solid rgba(var(--purple-color-rgb), 0.1);
  backdrop-filter: blur(8px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  margin-left: 45px; /* Add margin to align with other messages */
}

.loading-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px 16px;
  min-width: 120px;
  position: relative;
  min-height: 60px; /* Reduced height */
}

.loading-text {
  font-size: 0.9rem;
  color: var(--text-color);
  opacity: 0.8;
  font-weight: 500;
  letter-spacing: 0.3px;
}

.loading-indicator {
  display: flex;
  gap: 6px;
  justify-content: center;
  align-items: center;
  padding: 4px 0;
}

.loading-dot {
  width: 6px;
  height: 6px;
  background: var(--purple-color);
  border-radius: 50%;
  animation: wave 1.5s infinite ease-in-out;
  transform-origin: center;
  opacity: 0.6;
  box-shadow: 0 0 8px rgba(var(--purple-color-rgb), 0.3);
}

.loading-dot:nth-child(1) { animation-delay: -0.4s; }
.loading-dot:nth-child(2) { animation-delay: -0.3s; }
.loading-dot:nth-child(3) { animation-delay: -0.2s; }
.loading-dot:nth-child(4) { animation-delay: -0.1s; }
.loading-dot:nth-child(5) { animation-delay: 0s; }

@keyframes wave {
  0%, 100% {
    transform: translateY(0) scale(1);
    opacity: 0.6;
  }
  50% {
    transform: translateY(-8px) scale(1.2);
    opacity: 1;
  }
}

@keyframes fadeIn {
  from { 
    opacity: 0;
    transform: translateY(10px);
  }
  to { 
    opacity: 0.9;
    transform: translateY(0);
  }
}

/* Update the send button loading spinner */
.loading-spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;
  box-shadow: 0 0 8px rgba(255, 255, 255, 0.2);
}

@keyframes spin {
  to { 
    transform: rotate(360deg);
  }
}

/* Update disabled states */
button:disabled {
  opacity: 0.8;
  cursor: not-allowed;
  background: linear-gradient(135deg, var(--purple-color), var(--hover-color));
  box-shadow: 0 2px 8px rgba(var(--purple-color-rgb), 0.2);
}

.chatbot-textarea:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(var(--purple-color-rgb), 0.2);
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.05);
  pointer-events: none;
}

.stop-generation-btn,
.stop-generation-btn:hover,
.stop-generation-btn:disabled,
.stop-generation-btn svg,
.stop-generation-btn:hover svg {
  display: none;
}

/* Kickstart suggestions styles */
.kickstart-container {
  position: absolute;
  top: calc(100% + 20px);
  left: 0;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.2rem;
  animation: fadeUpIn 0.8s cubic-bezier(0.22, 1, 0.36, 1);
  padding: 1.5rem;
  background: rgba(30, 30, 40, 0.7);
  border-radius: 24px;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  z-index: 5;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.25), 
              inset 0 1px 1px rgba(255, 255, 255, 0.1),
              0 0 0 1px rgba(255, 255, 255, 0.05);
}

.kickstart-title {
  font-size: 1rem;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 600;
  margin: 0;
  background: linear-gradient(135deg, var(--purple-color), var(--hover-color));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 0.02em;
}

.kickstart-suggestions {
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
  width: 100%;
}

.kickstart-button {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.2rem 1.5rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  cursor: pointer;
  text-align: left;
  transition: all 0.3s cubic-bezier(0.22, 1, 0.36, 1);
  backdrop-filter: blur(8px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  color: var(--text-color);
  position: relative;
  overflow: hidden;
}

.kickstart-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, 
    rgba(255, 255, 255, 0),
    rgba(255, 255, 255, 0.2),
    rgba(255, 255, 255, 0));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.kickstart-button::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, 
    rgba(255, 255, 255, 0),
    rgba(255, 255, 255, 0.1),
    rgba(255, 255, 255, 0));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.kickstart-button:hover {
  background: rgba(var(--purple-color-rgb), 0.1);
  border-color: rgba(var(--purple-color-rgb), 0.3);
  transform: translateY(-3px) scale(1.01);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3),
              0 0 0 1px rgba(var(--purple-color-rgb), 0.2);
}

.kickstart-button:hover::before,
.kickstart-button:hover::after {
  opacity: 1;
}

.suggestion-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  background: linear-gradient(135deg, var(--purple-color), var(--hover-color));
  border-radius: 12px;
  flex-shrink: 0;
  color: white;
  box-shadow: 0 4px 15px rgba(var(--purple-color-rgb), 0.3),
              inset 0 1px 1px rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
  transform: rotate(0deg);
}

.kickstart-button:hover .suggestion-icon {
  transform: rotate(10deg) scale(1.1);
  box-shadow: 0 6px 20px rgba(var(--purple-color-rgb), 0.4),
              inset 0 1px 1px rgba(255, 255, 255, 0.3);
}

.suggestion-text {
  font-size: 1rem;
  line-height: 1.5;
  font-weight: 500;
  transition: color 0.3s ease;
}

.kickstart-button:hover .suggestion-text {
  color: rgba(var(--purple-color-rgb), 1);
}

@keyframes fadeUpIn {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .kickstart-container {
    top: calc(100% + 15px);
    padding: 1.2rem;
  }
  
  .kickstart-button {
    padding: 1rem 1.2rem;
  }
  
  .suggestion-icon {
    width: 34px;
    height: 34px;
  }
}
</style>
