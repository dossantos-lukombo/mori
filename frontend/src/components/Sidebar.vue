<template>
  <div :class="['sidebar', { 'sidebar--active': isActive }]">
    <div class="sidebar-content">
      <!-- Icônes du haut (accès direct) -->
      <ul class="icon-container">
        <li @click="navigateToMessages" class="icon-wrapper">
          <div class="icon-circle">
            <img src="@/assets/icons/messages.png" alt="Messagerie" />
          </div>
          <span>Messages</span>
        </li>

        <li @click="navigateToChatBot" class="icon-wrapper">
          <div class="icon-circle">
            <img src="@/assets/icons/chat.png" alt="Chat" />
          </div>
          <span>Mori Chatbot</span>
        </li>
      </ul>

      <!-- Zone d'affichage des contacts (amis + groupes) -->
      <ContactsForChatBotView
        v-if="activeView === 'contacts'"
        @select-contact="handleContactSelection"
      />
    </div>
  </div>
</template>

<script>
import ContactsForChatBotView from "./ContactsForChatBoxView.vue";

export default {
  name: "Sidebar",
  props: {
    isActive: {
      type: Boolean,
      required: true,
    },
    // Si tu n'utilises plus ce tableau "contactsList" directement,
    // tu peux le laisser ou le retirer selon ta logique
    contactsList: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      // Permet de savoir si on affiche la liste des contacts ou non
      activeView: null,
    };
  },
  components: { ContactsForChatBotView },
  methods: {
    async navigateToMessages() {
      // Si tu n'as aucun contact, on ouvre la vue "contacts" pour en ajouter ou voir
      if (this.contactsList.length === 0) {
        this.activeView = "contacts";
      } else {
        // EXEMPLE : si tu veux ouvrir directement le premier contact => DM
        // ou si tu préfères forcer l'utilisateur à cliquer => tu ouvres juste "contacts"
        // Ici, on choisit de juste ouvrir la liste :
        this.activeView = "contacts";

        // -- OU si tu veux ouvrir le 1er contact en DM, fais par ex. :
        
        const firstContact = this.contactsList[0];
        await this.$router.push({
          name: "messages",
          query: {
            name: firstContact.nickname,
            receiverId: firstContact.id,
            type: "PERSON", // <= si on sait que c'est un ami
          },
        });
      
      }
    },

    async navigateToChatBot() {
      await this.$router.push({ name: "mainpage" });
    },

    // >>> Correction principale <<<
    // On récupère l'objet { id, name, type } (émit par ContactsForChatBotView)
    // puis on navigue vers la route "messages" en passant 'type' tel quel.
    handleContactSelection({ id, name, type }) {
      // type peut valoir "PERSON" ou "GROUP" selon l'élément cliqué
      this.$router.push({
        name: "messages",
        query: {
          name,
          receiverId: id,
          type,
        },
      });
    },
  },
};
</script>

<style scoped>
.sidebar {
  position: fixed;
  top: 64.45px;
  left: -440px;
  width: 440px;
  height: calc(100% - 64.45px);
  background-color: var(--bg-neutral);
  transition: left 0.3s ease;
  z-index: 2;
  overflow-y: auto;
}

.sidebar--active {
  left: 0;
}

.sidebar-content {
  padding: 20px;
}

.icon-container {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
}

.icon-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  text-align: center;
  color: var(--text-color);
  transition: transform 0.3s ease;
}

.icon-wrapper:hover {
  transform: scale(1.1);
}

.icon-circle {
  width: 100px;
  height: 70px;
  background-color: var(--purple-color);
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
  transition: background-color 0.3s ease;
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.3);
}

.icon-circle:hover {
  background-color: var(--hover-color);
}

.icon-circle img {
  width: 32px;
  height: 32px;
}

.icon-wrapper span {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}
</style>
