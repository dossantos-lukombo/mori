<template>
    <div id="app">
        <div class="sidebar">
            <h3>Historique du Chatbot</h3>
            <ul>
                <li v-for="(message, index) in chatHistory" :key="index" @click="loadMessage(index)">
                    {{ message.summary }}
                </li>
            </ul>
        </div>
        <div class="chat-window">
            <h3>Chatbot</h3>
            <div class="chat-content">
                <div v-for="(msg, idx) in currentChat.messages" :key="idx"
                    :class="{ 'user-message': msg.sender === 'user', 'bot-message': msg.sender === 'bot' }">
                    {{ msg.text }}
                </div>
            </div>
            <input type="text" v-model="userInput" @keyup.enter="sendMessage" placeholder="Écrivez un message..." />
        </div>
    </div>
</template>
  
<script>
export default {
    data() {
        return {
            userInput: "",
            chatHistory: [],
            currentChat: { messages: [] }
        };
    },
    methods: {
        sendMessage() {
            if (this.userInput.trim() !== "") {
                const userMessage = { sender: "user", text: this.userInput };
                const botResponse = { sender: "bot", text: "Réponse générée par le LLM" };

                this.currentChat.messages.push(userMessage, botResponse);

                if (!this.chatHistory.includes(this.currentChat)) {
                    this.chatHistory.push({
                        summary: this.userInput.slice(0, 20) + "...",
                        messages: [...this.currentChat.messages],
                    });
                }

                this.userInput = "";
            }
        },
        loadMessage(index) {
            this.currentChat = this.chatHistory[index];
        }
    }
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

input[type="text"] {
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 5px;
    outline: none;
}
</style>
  