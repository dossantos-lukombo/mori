<template>
  <!-- <h3>Historique du Chatbot</h3> -->
  <div class="chat-history-container">
    <div class="chat-history-header">
      <div class="search-container">
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="Search conversations..." 
          class="search-input"
        />
        <div class="search-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M21 21L15 15M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
      </div>
  <button class="btn_new_convo" @click.stop="newConversation">
        <span class="plus-icon">+</span>
    New Conversation
  </button>
    </div>

    <div v-if="filteredConversations.length === 0" class="empty-state">
      <div class="empty-state-icon">💭</div>
      <p class="empty-state-text">{{ searchQuery ? 'No conversations found' : 'No conversations yet' }}</p>
      <p class="empty-state-subtext">{{ searchQuery ? 'Try different search terms' : 'Start a new conversation to begin' }}</p>
    </div>

    <div v-else class="conversations-list">
      <div
        v-for="(convo, index) in filteredConversations"
      :key="convo.conversation_id"
        :class="['conversation-card', { 
          'selected': selectedConvoId === convo.conversation_id,
          'new': isNewConversation(convo)
        }]"
        @click.stop="loadConvo(convo.conversation_id)"
      >
        <div class="conversation-content">
          <div class="conversation-header">
            <span class="conversation-type" :class="getConversationType(convo)">
              {{ getConversationTypeLabel(convo) }}
            </span>
            <span class="conversation-time">{{ formatTime(convo.convo[convo.convo.length - 1].created_at) }}</span>
          </div>
          
          <div class="conversation-preview">
            <div class="preview-message user">
              <span class="message-icon">👤</span>
              {{ truncateText(convo.convo[convo.convo.length - 1].user_request, 60) }}
            </div>
            <div class="preview-message ai">
              <span class="message-icon">🤖</span>
              {{ truncateText(convo.convo[convo.convo.length - 1].llm_response, 60) }}
            </div>
          </div>

          <div class="conversation-meta">
            <span class="message-count" :class="{ updating: isUpdatingMessageCount(convo) }">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="message-icon">
                <path d="M20 2H4C2.9 2 2 2.9 2 4V22L6 18H20C21.1 18 22 17.1 22 16V4C22 2.9 21.1 2 20 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              {{ convo.messageCount || getMessageCount(convo) }} {{ (convo.messageCount || getMessageCount(convo)) === 1 ? 'message' : 'messages' }}
            </span>
            <span class="conversation-status" :class="getConversationStatus(convo)">
              {{ getConversationStatusLabel(convo) }}
            </span>
          </div>
        </div>

        <div class="delete-button-wrapper" @click.stop>
      <button
        class="btn_delete_convo"
        @click.stop="deleteConvo(convo.conversation_id)"
            title="Delete conversation"
      >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M19 7L18.1327 19.1425C18.0579 20.1891 17.187 21 16.1378 21H7.86224C6.81296 21 5.94208 20.1891 5.86732 19.1425L5 7M10 11V17M14 11V17M15 7V4C15 3.44772 14.5523 3 14 3H10C9.44772 3 9 3.44772 9 4V7M4 7H20" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
      </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal-overlay" @click.self="cancelDelete">
      <div class="modal-content">
        <div class="modal-icon">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 9V14M12 19C7.58172 19 4 15.4183 4 11C4 6.58172 7.58172 3 12 3C16.4183 3 20 6.58172 20 11C20 15.4183 16.4183 19 12 19ZM12.0498 16V16.1L11.9502 16.1V16H12.0498Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <h3 class="modal-title">Delete Conversation</h3>
        <p class="modal-message">Are you sure you want to delete this conversation? This action cannot be undone.</p>
        <div class="modal-actions">
          <button class="modal-btn cancel" @click="cancelDelete">Cancel</button>
          <button class="modal-btn delete" @click="confirmDelete">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      userInput: "",
      selectedConvoId: null,
      searchQuery: "",
      newConversations: new Set(), // Track new conversations for animation
      showDeleteModal: false,
      conversationToDelete: null,
    };
  },
  computed: {
    conversationsUpdate() {
      // Always return a new reversed array, don't mutate chatHistory directly
      console.log("computed conversationsUpdate: ", this.$store.getters.allConversations);
      // return [...this.$store.getters.allConversations].reverse();
      return [...this.$store.getters.allConversations].reverse();
    },
    filteredConversations() {
      if (!this.searchQuery) return this.conversationsUpdate;
      
      const query = this.searchQuery.toLowerCase();
      return this.conversationsUpdate.filter(convo => {
        const lastMessage = convo.convo[convo.convo.length - 1];
        return lastMessage.user_request.toLowerCase().includes(query) ||
               lastMessage.llm_response.toLowerCase().includes(query);
      });
    },
    currentMessages() {
      return this.$store.getters.allMessages;
    }
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
      this.selectedConvoId = null;
      this.$store.dispatch("clearMessages");
      
      // Navigate to the chatbot page if not already there
      if (this.$route.name !== "mainpage") {
        this.$router.push({ name: "mainpage" });
      }
    },

    async deleteConvo(conversation_id) {
      this.conversationToDelete = conversation_id;
      this.showDeleteModal = true;
    },

    cancelDelete() {
      this.showDeleteModal = false;
      this.conversationToDelete = null;
    },

    async confirmDelete() {
      if (!this.conversationToDelete) return;

      try {
      const response = await fetch("http://localhost:8081/llmConvoDelete", {
        credentials: "include",
          headers: {
          "Content-Type": "application/json",
          },
        method: "POST",
        body: JSON.stringify({
          user_id: await this.getMyUserID(),
            conversation_id: this.conversationToDelete,
        }),
      });

      if (!response.ok) {
          throw new Error(response.statusText);
        }

        // Remove from store
        this.$store.dispatch("deleteConversationById", this.conversationToDelete);
        this.$store.dispatch("clearMessages");
        
        // Clear selection if this was the selected conversation
        if (this.selectedConvoId === this.conversationToDelete) {
          this.selectedConvoId = null;
        }

        // Close modal and reset state
        this.showDeleteModal = false;
        this.conversationToDelete = null;
      } catch (error) {
        console.error("Error deleting conversation:", error);
        // You could show an error toast here
      }
    },

    //Méthode pour obtenir la discussion selctionnée
    async loadConvo(convo_id) {
      console.log("convo_id in loadConvo: ", convo_id);
      
      // Navigate to the chatbot page if not already there
      if (this.$route.name !== "mainpage") {
        await this.$router.push({ name: "mainpage" });
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
        console.log("RESPONSE loadConvo: ", resp);
        console.log("Current convo: ", resp.convo);
        console.log("Current convo_id loadConvo: ", resp.conversation_id);
        this.$store.dispatch("clearMessages");
        this.convertMessages(resp);
        this.selectedConvoId = convo_id;
        this.newConversations.delete(convo_id);
        
        // Update message count after loading
        this.updateMessageCount(convo_id);
      }
    },

    convertMessages(convo) {
      console.log("convo in convertMessages: ", convo);

      if (convo === null) {
        console.log("convo is null ");
        return;
      }
      
      this.$store.dispatch("clearMessages");
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
      
      // Update message count after converting messages
      this.updateMessageCount(convo.conversation_id);
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
      resp.forEach((convo) => {
        this.$store.dispatch("addConversation", convo);
        this.newConversations.add(convo.conversation_id);
      });
      // Ne touche plus à this.chatHistory ici !
      console.log("Conversations ajoutées au store.");
      // Clear new status after 5 seconds
      setTimeout(() => {
        this.newConversations.clear();
      }, 5000);
    },

    formatTime(timestamp) {
      if (!timestamp) return '';
      const date = new Date(timestamp);
      if (isNaN(date.getTime())) return '';
      
      const now = new Date();
      const diff = now - date;
      const days = Math.floor(diff / (1000 * 60 * 60 * 24));
      
      if (days === 0) {
        return 'Today';
      } else if (days === 1) {
        return 'Yesterday';
      } else if (days < 7) {
        return `${days} days ago`;
      } else {
        return date.toLocaleDateString();
      }
    },

    truncateText(text, length) {
      if (!text) return '';
      return text.length > length ? text.substring(0, length) + '...' : text;
    },

    isNewConversation(convo) {
      return this.newConversations.has(convo.conversation_id);
    },

    getConversationType(convo) {
      // You can implement logic to determine conversation type
      // For now, returning a default type
      return 'general';
    },

    getConversationTypeLabel(convo) {
      const type = this.getConversationType(convo);
      const labels = {
        general: 'General',
        code: 'Code',
        creative: 'Creative',
        analysis: 'Analysis'
      };
      return labels[type] || 'General';
    },

    getConversationStatus(convo) {
      // Implement logic to determine conversation status
      return 'active';
    },

    getConversationStatusLabel(convo) {
      const status = this.getConversationStatus(convo);
      const labels = {
        active: 'Active',
        completed: 'Completed',
        pending: 'Pending'
      };
      return labels[status] || 'Active';
    },

    getMessageCount(convo) {
      if (!convo || !convo.convo) return 0;
      
      // If this is the current conversation, use the live message count
      if (convo.conversation_id === this.selectedConvoId) {
        const count = this.currentMessages.length;
        this.$store.dispatch('updateConversationMessageCount', {
          conversationId: convo.conversation_id,
          count
        });
        return count;
      }
      
      // For other conversations, count from the stored messages
      const uniqueMessages = new Set();
      convo.convo.forEach(message => {
        if (message.user_request) uniqueMessages.add(message.user_request);
        if (message.llm_response) uniqueMessages.add(message.llm_response);
      });
      
      const count = uniqueMessages.size;
      this.$store.dispatch('updateConversationMessageCount', {
        conversationId: convo.conversation_id,
        count
      });
      
      return count;
    },

    // Update the store when messages change
    updateMessageCount(conversationId) {
      const conversation = this.$store.getters.allConversations.find(
        conv => conv.conversation_id === conversationId
      );
      
      if (conversation) {
        const count = this.getMessageCount(conversation);
        this.$store.dispatch('updateConversationMessageCount', {
          conversationId,
          count
        });
      }
    },

    isUpdatingMessageCount(convo) {
      // Implement logic to determine if the message count is updating
      return false; // Placeholder, actual implementation needed
    }
  },
  
  watch: {
    // Watch for changes in the store's conversations
    '$store.getters.allConversations': {
      handler(newConversations) {
        newConversations.forEach(convo => {
          if (!convo.messageCount) {
            this.updateMessageCount(convo.conversation_id);
          }
        });
      },
      deep: true
    },

    // Watch for changes in messages to update count in real-time
    currentMessages: {
      handler(newMessages) {
        if (this.selectedConvoId) {
          // Update the count for the current conversation
          this.$store.dispatch('updateConversationMessageCount', {
            conversationId: this.selectedConvoId,
            count: newMessages.length
          });
        }
      },
      deep: true
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

.btn_elmt_history {
  padding: 10px;
  border: 1px solid #ddd;
  margin: 5px;
  border-radius: 5px;
  background-color: black;
  color: white;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.btn_elmt_history:hover {
  background-color: none;
  border-color: none;
}

input[type="text"] {
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 5px;
  outline: none;
  padding-left: 40px;
}
.selected .elmt_history {
  background-color: #333;
  color: #fff;
}
</style>

<style>
.chat-history-container {
  padding: 1.5rem;
  background: linear-gradient(135deg, var(--bg-neutral), var(--page-bg));
  height: 100%;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  border-right: none;
}

.chat-history-header {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  position: sticky;
  top: 0;
  padding-bottom: 1rem;
  z-index: 10;
}

.search-container {
  position: relative;
  width: 100%;
  margin-bottom: 0.5rem;
}

.search-input {
  width: 100%;
  padding: 0.875rem 1rem 0.875rem 3rem;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-color);
  font-size: 0.95rem;
  transition: all 0.3s ease;
  backdrop-filter: blur(8px);
  box-shadow: none;
}

.search-input::placeholder {
  color: var(--text-color-secondary);
  opacity: 0.7;
  padding-left: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.search-input:focus {
  outline: none;
  border-color: var(--purple-color);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 2px rgba(var(--purple-color-rgb), 0.2);
}

.search-icon {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-color-secondary);
  opacity: 0.7;
  pointer-events: none;
  transition: opacity 0.3s ease;
  z-index: 1;
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-input:focus + .search-icon {
  opacity: 1;
  color: var(--purple-color);
}

.btn_new_convo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.875rem 1rem;
  border-radius: 12px;
  background: var(--purple-color);
  color: white;
  border: none;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.btn_new_convo:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  background: var(--hover-color);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 1rem;
  text-align: center;
  color: var(--text-color-secondary);
}

.empty-state-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state-text {
  font-size: 1.25rem;
  font-weight: 500;
  margin-bottom: 0.5rem;
  color: var(--text-color);
}

.empty-state-subtext {
  font-size: 0.95rem;
  opacity: 0.7;
}

.conversations-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.conversation-card {
  display: flex;
  align-items: stretch;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
  gap: 0;
  cursor: pointer;
}

.conversation-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(45deg, transparent, rgba(255, 255, 255, 0.03), transparent);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.conversation-card:hover::before {
  opacity: 1;
}

.conversation-card:hover {
  transform: translateX(4px) translateY(-2px);
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.conversation-card.selected {
  background: rgba(var(--purple-color-rgb), 0.15);
  border-color: var(--purple-color);
}

.conversation-card.new {
  animation: newConversation 0.5s ease-out;
}

@keyframes newConversation {
  0% {
    transform: translateX(-100%);
    opacity: 0;
  }
  100% {
    transform: translateX(0);
    opacity: 1;
  }
}

.conversation-content {
  flex: 1;
  padding: 1.25rem;
  min-width: 0;
}

.conversation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.conversation-type {
  font-size: 0.8rem;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  background: rgba(var(--purple-color-rgb), 0.1);
  color: var(--purple-color);
}

.conversation-type.code { background: rgba(52, 152, 219, 0.1); color: #3498db; }
.conversation-type.creative { background: rgba(155, 89, 182, 0.1); color: #9b59b6; }
.conversation-type.analysis { background: rgba(46, 204, 113, 0.1); color: #2ecc71; }

.conversation-preview {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.preview-message {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  font-size: 0.9rem;
  line-height: 1.4;
  color: var(--text-color);
  opacity: 0.9;
}

.preview-message.user { color: var(--text-color); }
.preview-message.ai { color: var(--text-color-secondary); }

.message-icon {
  font-size: 1rem;
  opacity: 0.7;
}

.conversation-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  color: var(--text-color-secondary);
}

.message-count {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 0.8rem;
  color: var(--text-color-secondary);
  opacity: 0.9;
  transition: all 0.2s ease;
}

.message-count.updating {
  animation: pulse 0.5s ease-in-out;
}

@keyframes pulse {
  0% {
    opacity: 0.9;
  }
  50% {
    opacity: 0.5;
  }
  100% {
    opacity: 0.9;
  }
}

.conversation-status {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
}

.conversation-status.active { background: rgba(46, 204, 113, 0.1); color: #2ecc71; }
.conversation-status.completed { background: rgba(52, 152, 219, 0.1); color: #3498db; }
.conversation-status.pending { background: rgba(241, 196, 15, 0.1); color: #f1c40f; }

.conversation-actions {
  display: flex;
  align-items: center;
  padding: 0 0.75rem;
  border-left: 1px solid rgba(255, 255, 255, 0.1);
  opacity: 0;
  transition: opacity 0.2s ease;
}

.conversation-card:hover .conversation-actions {
  opacity: 1;
}

.delete-button-wrapper {
  display: flex;
  align-items: center;
  padding: 0 0.75rem;
  border-left: 1px solid rgba(255, 255, 255, 0.1);
  opacity: 0;
  transition: opacity 0.2s ease;
  pointer-events: auto;
}

.conversation-card:hover .delete-button-wrapper {
  opacity: 1;
}

.btn_delete_convo {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: none;
  color: var(--text-color-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  border-radius: 8px;
  padding: 0;
  opacity: 0.7;
}

.btn_delete_convo:hover {
  color: #ff4444;
  background: rgba(255, 68, 68, 0.1);
  transform: scale(1.1);
  opacity: 1;
}

.btn_delete_convo svg {
  width: 16px;
  height: 16px;
  transition: transform 0.2s ease;
}

.btn_delete_convo:hover svg {
  transform: scale(1.1);
}

/* Scrollbar styling */
.chat-history-container::-webkit-scrollbar {
  width: 6px;
}

.chat-history-container::-webkit-scrollbar-track {
  background: transparent;
}

.chat-history-container::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

.chat-history-container::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .chat-history-container {
    padding: 1rem;
  }
  
  .conversation-card {
    border-radius: 12px;
  }
  
  .conversation-content {
    padding: 1rem;
  }
  
  .preview-message {
    font-size: 0.85rem;
  }
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease-out;
}

.modal-content {
  background: var(--bg-neutral);
  border-radius: 20px;
  padding: 2rem;
  width: 90%;
  max-width: 400px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.1);
  animation: slideUp 0.3s ease-out;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 1.5rem;
}

.modal-icon {
  color: #ff4444;
  background: rgba(255, 68, 68, 0.1);
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.5rem;
}

.modal-title {
  color: var(--text-color);
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

.modal-message {
  color: var(--text-color-secondary);
  font-size: 1rem;
  line-height: 1.5;
  margin: 0;
  opacity: 0.9;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  width: 100%;
  margin-top: 0.5rem;
}

.modal-btn {
  flex: 1;
  padding: 0.875rem 1.5rem;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.modal-btn.cancel {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-color);
}

.modal-btn.cancel:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-1px);
}

.modal-btn.delete {
  background: #ff4444;
  color: white;
}

.modal-btn.delete:hover {
  background: #ff3333;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(255, 68, 68, 0.2);
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
