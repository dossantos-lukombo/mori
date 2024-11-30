<template>
    <div class="profile-page">
      <div class="profile-header">
        <h1>Mon Profil</h1>
        <button @click="goToEditProfile">Modifier le profil</button>
      </div>
      <div class="profile-info">
        <div class="avatar-section">
          <img :src="avatarUrl" alt="Avatar" class="avatar" />
          <input type="file" ref="avatarInput" @change="changeAvatar" />
        </div>
        <div class="user-info">
          <div class="info-item">
            <label>Nom d'utilisateur:</label>
            <span>{{ userInfo.username }}</span>
          </div>
          <div class="info-item">
            <label>Email:</label>
            <span>{{ userInfo.email }}</span>
          </div>
          <button @click="changePassword">Changer le mot de passe</button>
        </div>
      </div>
      <div class="profile-actions">
        <button @click="deleteAccount" class="danger">Supprimer le compte</button>
        <button @click="deleteFavorites" class="danger">Supprimer les favoris</button>
        <button @click="deleteConversations" class="danger">Supprimer les conversations</button>
      </div>
    </div>
  </template>
  
  <script>
  import axios from "axios";
  axios.defaults.baseURL = 'http://localhost:8081/api';
  
  export default {
    name: "ProfilePage",
    data() {
      return {
        avatarUrl: "", 
        userInfo: {
          username: "",
          email: "",
        },
      };
    },
    created() {
      this.fetchUserProfile();
    },
    methods: {
      async fetchUserProfile() {
        try {
          const response = await axios.get("/api/profile");
          this.userInfo = response.data;
          this.avatarUrl = response.data.avatarUrl;
        } catch (error) {
          console.error("Erreur lors de la récupération du profil:", error);
        }
      },
      goToEditProfile() {
        this.$router.push({ name: "UpdateUserInfo" });
      },
      async changeAvatar() {
        try {
          const formData = new FormData();
          formData.append("avatar", this.$refs.avatarInput.files[0]);
          const response = await axios.post("/api/profile/avatar", formData);
          if (response.status === 200) {
            alert("Avatar mis à jour avec succès!");
            this.fetchUserProfile(); 
          } else {
            console.error("Erreur lors du changement de l'avatar");
          }
        } catch (error) {
          console.error("Erreur lors du changement de l'avatar:", error);
        }
      },
      async changePassword() {
        this.$router.push({ name: "ChangePassword" });
      },
      async deleteAccount() {
        try {
          const response = await axios.delete("/api/profile/delete-account");
          if (response.status === 200) {
            alert("Compte supprimé avec succès!");
            this.$router.push({ name: "LoginPage" }); // Rediriger après la suppression du compte
          } else {
            console.error("Erreur lors de la suppression du compte");
          }
        } catch (error) {
          console.error("Erreur lors de la suppression du compte:", error);
        }
      },
      async deleteFavorites() {
        try {
          const response = await axios.delete("/api/profile/delete-favorites");
          if (response.status === 200) {
            alert("Favoris supprimés avec succès!");
          } else {
            console.error("Erreur lors de la suppression des favoris");
          }
        } catch (error) {
          console.error("Erreur lors de la suppression des favoris:", error);
        }
      },
      async deleteConversations() {
        try {
          const response = await axios.delete("/api/profile/delete-conversations");
          if (response.status === 200) {
            alert("Conversations supprimées avec succès!");
          } else {
            console.error("Erreur lors de la suppression des conversations");
          }
        } catch (error) {
          console.error("Erreur lors de la suppression des conversations:", error);
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .profile-page {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    background-color: #2c2c2c;
    border-radius: 10px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
    color: #e4e4e4;
  }
  
  .profile-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }
  
  .profile-info {
    display: flex;
    gap: 20px;
    margin-bottom: 20px;
  }
  
  .avatar-section {
    text-align: center;
  }
  
  .avatar {
    width: 150px;
    height: 150px;
    border-radius: 50%;
    margin-bottom: 10px;
  }
  
  .user-info {
    flex: 1;
  }
  
  .info-item {
    margin-bottom: 10px;
  }
  
  .profile-actions {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  
  button {
    background-color: #9146bc;
    color: white;
    border: none;
    padding: 10px;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.3s;
  }
  
  button:hover {
    background-color: #7d3aa6;
  }
  
  button.danger {
    background-color: #ff4d4d;
  }
  
  button.danger:hover {
    background-color: #e03c3c;
  }
  </style>
  