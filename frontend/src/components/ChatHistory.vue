<template>
    <!-- <div id="app"> -->
        <!-- <div class="sidebar"> -->
            <!-- <h3>Historique du Chatbot</h3> -->
            <ul>
                <li v-for="(convo,index) in chatHistory" :key="index">
                    <div v-if="convo.length > 0" @click="loadConvo(convo[0].conversation_id)">
                        <button type="button">{{convo[convo.length-1].user_request}}</button>
                    </div>
                </li>
            </ul>
        <!-- </div> -->
    <!-- </div> -->
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
            })
            if (!response.ok) {
                console.error("Erreur lors de la récupération de l'ID de l'utilisateur :", response.statusText);
                return;
            } else {
                const resp = await response.json();
                console.log(resp);
                console.log(resp.users[0].id);
                return resp.users[0].id;
            }
        },

        //Méthode pour obtenir la discussion selctionnée
        async loadConvo(convo_id) {
            const response = await fetch("http://localhost:8081/llmConvoGet", {
                credentials: "include",
                headers: new Headers({
                    "Content-Type": "application/json",
                }),

                method: "POST",
                body: JSON.stringify({
                    user_id: await this.getMyUserID(),
                    conversation_id: convo_id
                }),
            })
            if (!response.ok) {
                console.error("Erreur lors de la récupération de la conversation de l'utilisateur :", response.statusText);
                return;
            } else {
                const resp = await response.json();
                console.log("CONVERSATION: ", resp);
                this.currentChat = resp;
            }
        },

        //Méthode pour obtenir les conversations de l'utilisateur
        async getChatHistory() {
            let same_conversations_id = {};

            const response = await fetch("http://localhost:8081/llmConvoGet", {
                credentials: "include",
                headers: new Headers({
                    "Content-Type": "application/json",
                }),

                method: "POST",
                body: JSON.stringify({ user_id: await this.getMyUserID() }),
            })
            if (!response.ok) {
                console.error("Erreur lors de la récupération des conversations de l'utilisateur :", response.statusText);
                return;
            } else {
                const resp = await response.json();

                console.log("LIST OF CONVERSATION: ", resp);
                let same_convo = [];
                let diff_convo = [];
                for (let i = 0; i < resp.length; i++) {
                    if (i<resp.length-1 && resp[i].conversation_id === resp[i+1].conversation_id) {
                        same_convo.push(resp[i]);
                        
                    }else if (i<resp.length-1 && resp[i].conversation_id !== resp[i+1].conversation_id) {
                        if (i>0 && resp[i].conversation_id === resp[i-1].conversation_id) {
                            same_convo.push(resp[i]);
                        }
                    }
                }
                this.chatHistory.push(same_convo);
                for (let i = 0; i < resp.length; i++) {
                    if (!same_convo.includes(resp[i])) {
                        this.chatHistory.push([resp[i]]);
                    }
                }

                
                
                
                // return same_conversations_id;
                // this.chatHistory = same_conversations_id;
                console.log("SAME CONVERSATION ID: ", this.chatHistory);
            }
            // this.chatHistory = same_conversations_id;
        },
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
  