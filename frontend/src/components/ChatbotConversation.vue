<template>
    <div class="chatbot-container">
      <div :class="['chatbot-box', { 'chatbot-box--active': hasMessages }]">
        <div class="mori-img" v-if="!hasMessages">
          <div class="moriImg"></div>
        </div>
        <div class="mori" id="moriChatBot" v-if="!hasMessages">Mori</div>
        <div class="chatbot-message" v-if="!hasMessages">
          How can I help you?
        </div>
        <div class="chatbot-messages" v-if="hasMessages">
            <div 
              v-for="(message, index) in messages" 
              :key="index" 
              :class="['message', message.isUser ? 'user' : 'bot']"
            >
              <div v-if="!message.isUser" class="bot-logo">
                <img src="../assets/mori.png" alt="Bot Logo" />
              </div>
              <span>{{ message.text }}</span>
              <div class="timestamp">{{ message.timestamp }}</div>
            </div>
          </div>
          
        <div :class="['chatbot-input', { 'chatbot-input--active': hasMessages }]">
          <input
            type="text"
            v-model="userInput"
            @keyup.enter="handleUserInput"
            placeholder="Type your message here..."
          />
          <button @click="handleUserInput">Send</button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: "ChatbotConversation",
    data() {
      return {
        userInput: "", // Stores the user's input
        messages: [], // Stores the messages
      };
    },
    computed: {
      hasMessages() {
        return this.messages.length > 0; // Check if there are any messages
      },
    },
    methods: {
      handleUserInput() {
        if (this.userInput.trim()) {
          // Add user message with timestamp
          this.messages.push({ 
            text: this.userInput, 
            isUser: true, 
            timestamp: this.formatTimestamp(new Date()) 
          });
          this.userInput = ""; // Clear input after sending
  
          // Simulate bot response
          setTimeout(() => {
            this.messages.push({ 
              text: "Thank you for your message!", 
              isUser: false, 
              timestamp: this.formatTimestamp(new Date()) 
            });
          }, 1000);
        }
      },
      formatTimestamp(date) {
        const options = { hour: '2-digit', minute: '2-digit', hour12: false };
        const time = date.toLocaleTimeString([], options);
        const day = date.toLocaleDateString();
        return `${day} ${time}`;
      },
    },
  };
  </script>
  
  <style scoped>

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
    margin-bottom: 130px;
    transition: all 0.5s ease;
  }
  
  .chatbot-box--active {
    justify-content: space-between;
    width: 80%;
    height: 85vh;
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
  .chatbot-input {
    display: flex;
    gap: 10px;
    align-items: center;
    position: absolute;
    top: 65%; /* Initially positioned below the "How can I help?" */
    left: 50%;
    transform: translate(-50%, -50%);
    width: calc(40% - 40px);
    border-radius: 10px;
    padding: 10px 20px;
    transition: top 0.7s ease, transform 0.7s ease, width 0.7s ease; /* Added width for smooth transition */
  }
  
  .chatbot-input--active {
    width: calc(50% - 40px); /* New width for active state */
    position: fixed;
    top: calc(97% - 80px); /* Slide to the bottom of the viewport */
    transform: translateX(-50%);
  }
  
  .chatbot-input input {
    flex: 1;
    padding: 10px;
    border: 1px solid var(--color-grey);
    border-radius: 10px;
    height: 50px;
    font-size: 16px;
  }
  
  .chatbot-input button {
    padding: 15px 20px;
    background-color: var(--purple-color);
    color: var(--color-white);
    border: none;
    border-radius: 10px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.3s;
  }
  
  .chatbot-input button:hover {
    background-color: var(--hover-background-color);
  }
  
  .chatbot-input input:focus {
    outline: none;
    border-color: var(--color-primary);
  }
  </style>
  