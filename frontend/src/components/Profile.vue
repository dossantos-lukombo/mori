<template>
  <div v-if="user && $store.state.id !== ''" class="profile-container">
    <div class="profile-header">
      <div class="profile-cover-photo"></div>
      <div class="profile-avatar-container">
        <div 
          class="profile-avatar"
          :style="{ backgroundImage: `url(${getAvatarUrl(user.avatar)})` }"
        ></div>
      </div>
      
      <div class="profile-header-content">
        <div class="profile-user-info">
          <h1 class="profile-name">{{ user.firstName }} {{ user.lastName }}</h1>
          <h2 v-if="showNickname" class="profile-nickname">@{{ user.nickname }}</h2>
          <p class="profile-email" v-if="user.login">{{ user.login }}</p>
          <p class="profile-dateOfBirth" v-if="user.dateOfBirth">{{ user.dateOfBirth }}</p>
        </div>
        
        <div class="profile-actions">
          <PrivacyBtn v-if="isMyProfile" :status="user.status" class="profile-btn" />
          <component
            v-else
            :is="displayBtn"
            v-bind="{ user }"
            @follow="checkFollowRequest"
            @unfollow="unfollow"
            class="profile-btn"
          ></component>
        </div>
      </div>
    </div>

    <div class="profile-content" v-if="showProfileData">
      <div class="profile-sidebar">
        <div class="profile-card">
          <h3 class="card-title">Connections</h3>
          <div class="connections-stats">
            <div class="stat-item">
              <span class="stat-value">{{ followers ? followers.length : 0 }}</span>
              <span class="stat-label">Followers</span>
            </div>
            <div class="stat-divider"></div>
            <div class="stat-item">
              <span class="stat-value">{{ following ? following.length : 0 }}</span>
              <span class="stat-label">Following</span>
            </div>
          </div>
        </div>
        
        <div class="multiple-item-list">
          <Following :following="following" />
          <Followers :followers="followers" />
        </div>

        <Groups :groups="profileGroups" />
      </div>

      <div class="profile-main">
        <div v-if="user.about !== ''" class="profile-card about-card">
          <h3 class="card-title">About me</h3>
          <p class="about-text">{{ user.about }}</p>
        </div>

        <!-- Edit Profile Section (only for own profile) -->
        <div v-if="isMyProfile" class="profile-card edit-profile-card">
          <button class="toggle-edit-btn" @click="toggleEditBanner">
            {{ editBannerOpen ? "Close Edit Profile" : "Edit Profile" }}
          </button>
          
          <div v-if="editBannerOpen" class="edit-profile-form">
            <h3 class="card-title">Edit Profile</h3>

            <!-- Change Nickname -->
            <div class="edit-section">
              <label for="nickname">New Nickname</label>
              <input
                id="nickname"
                v-model="newNickname"
                type="text"
                placeholder="Enter your new nickname"
                class="edit-input"
              />
              <button @click="changeNickname" class="edit-btn edit-btn-primary">
                Update Nickname
              </button>
            </div>

            <!-- Change Avatar -->
            <div class="edit-section">
              <label for="avatar">New Avatar</label>
              <div class="file-upload-wrapper">
                <input
                  id="avatar"
                  type="file"
                  ref="avatarInput"
                  class="file-upload"
                  accept="image/*"
                />
                <span class="file-upload-text">Choose a file</span>
              </div>
              <button @click="changeAvatar" class="edit-btn edit-btn-primary">
                Update Avatar
              </button>
            </div>

            <!-- Delete Account -->
            <div class="edit-section danger-zone">
              <h4 class="danger-title">Danger Zone</h4>
              <button @click="openDeleteModal" class="edit-btn edit-btn-danger">
                Delete Account
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Private profile message -->
    <div v-if="!showProfileData" class="private-profile-message">
      <div class="lock-icon">
        <i class="fas fa-lock"></i>
      </div>
      <p>This profile is private</p>
    </div>

    <!-- Delete Account Modal -->
    <div v-if="showDeleteModal" class="modal-overlay">
      <div class="modal-content">
        <h3 class="modal-title">Are you sure you want to delete your account?</h3>
        <p class="modal-text">This action cannot be undone.</p>
        <div class="modal-buttons">
          <button @click="cancelDelete" class="modal-btn modal-btn-cancel">Cancel</button>
          <button @click="confirmDelete" class="modal-btn modal-btn-confirm">Delete Account</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Following from "./Following.vue";
import Followers from "./Followers.vue";
import FollowBtn from "./FollowBtn.vue";
import PrivacyBtn from "./PrivacyBtn.vue";
import UnfollowBtn from "./UnfollowBtn.vue";
import Groups from "./Groups.vue";
import { getAvatarUrl } from '../utils/imageHelper';

export default {
  name: "Profile",
  components: { Following, Followers, FollowBtn, PrivacyBtn, UnfollowBtn, Groups },
  data() {
    return {
      user: null,
      isMyProfile: false,
      followers: [],
      following: [],
      profileGroups: null,
      newNickname: "", // New nickname
      selectedAvatar: null, // File for avatar
      editBannerOpen: false, // Controls visibility of the edit-profile banner
      showDeleteModal: false, // Controls visibility of the delete account modal
    };
  },
  created() {
    this.updateProfileData();
  },
  computed: {
    showProfileData() {
      return this.user && (this.user.following || this.isMyProfile || this.user.status === "PUBLIC");
    },
    showNickname() {
      return this.user && this.user.nickname !== `${this.user.firstName} ${this.user.lastName}`;
    },
    displayBtn() {
      return this.user && this.user.following ? UnfollowBtn : FollowBtn;
    },
  },
  methods: {
    getAvatarUrl,
    updateProfileData() {
      this.getUserData();
      this.getFollowers();
      this.getFollowing();
      this.checkProfile();
      this.getProfileGroups();
    },
    async getUserData() {
      await fetch("http://localhost:8081/userData?userId=" + this.$route.params.id, {
        credentials: "include",
      })
        .then((r) => r.json())
        .then((json) => {
          console.log("/getUserData", json);
          console.log("id", this.$route.params.id);
          this.user = json.users[0];
        });
    },
    async getProfileGroups() {
      const response = await fetch("http://localhost:8081/otherUserGroups?userId=" + this.$route.params.id, {
        credentials: "include",
      });
      const data = await response.json();
      if (data.type == "Error") {
        console.log("/getProfileGroups error: ", data.message);
      } else {
        this.profileGroups = data.groups;
      }
    },
    async getMyUserID() {
      await this.$store.dispatch("getMyUserID");
    },
    async checkProfile() {
      await this.getMyUserID();
      const myID = this.$store.state.id;
      const profileID = this.$route.params.id;
      this.isMyProfile = (profileID === myID);
    },
    checkFollowRequest(action) {
      if (action === "followedUser") {
        this.$store.dispatch("fetchChatUserList");
        this.updateProfileData();
        this.toggleFollowingThisUser();
      }
    },
    toggleFollowingThisUser() {
      this.user.following = !this.user.following;
    },
    unfollow() {
      this.updateProfileData();
      this.$store.dispatch("fetchChatUserList");
    },
    async getFollowers() {
      await fetch("http://localhost:8081/followers?userId=" + this.$route.params.id, {
        credentials: "include",
      })
        .then((response) => response.json())
        .then((json) => {
          this.followers = json.users || [];
        });
    },
    async getFollowing() {
      await fetch("http://localhost:8081/following?userId=" + this.$route.params.id, {
        credentials: "include",
      })
        .then((response) => response.json())
        .then((json) => {
          this.following = json.users || [];
        });
    },
    async changeNickname() {
      await this.$store.dispatch("changeNickname", this.newNickname);
      this.updateProfileData();
    },
    async changeAvatar() {
      const input = this.$refs.avatarInput;
      if (!input || !input.files || input.files.length === 0) {
        this.$toast.open({
          message: "Please select a file.",
          type: "warning",
        });
        return;
      }
      this.selectedAvatar = input.files[0];
      await this.$store.dispatch("changeAvatar", this.selectedAvatar);
      this.updateProfileData();
    },
    addChat() {
      let newChat = {
        name: this.user.nickname,
        receiverId: this.user.id,
        type: "PERSON",
      };
      this.$store.dispatch("addNewChat", newChat);
    },
    toggleEditBanner() {
      this.editBannerOpen = !this.editBannerOpen;
    },
    // Delete Account logic
    openDeleteModal() {
      this.showDeleteModal = true;
    },
    cancelDelete() {
      this.showDeleteModal = false;
    },
    async confirmDelete() {
      try {
        await this.$store.dispatch("deleteAccount");
        this.$toast.open({
          message: "Account deleted successfully.",
          type: "success",
        });
      } catch (err) {
        console.error("Error deleting account:", err);
        this.$toast.open({
          message: err.message,
          type: "error",
        });
      }
    },
  },
  watch: {
    $route() {
      if (this.$route.name === "Profile") {
        this.updateProfileData();
      }
    },
  },
};
</script>

<style scoped>
.profile-container {
  margin: 35px auto;
  width: 100%;
  max-width: 1200px;
  padding: 0 20px;
  color: var(--color-white);
  animation: fadeIn-bf1681ae 0.5s ease;
}

@keyframes fadeIn-bf1681ae {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* Profile Header */
.profile-header {
  position: relative;
  margin-bottom: 30px;
  background-color: var(--bg-neutral);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.profile-cover-photo {
  height: 200px;
  background: linear-gradient(135deg, var(--clear-purple) 0%, var(--purple-color) 100%);
  position: relative;
}

.profile-cover-photo::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewBox="0 0 100 100"><rect fill="none" width="100%" height="100%"/><path d="M0 0L100 100M100 0L0 100" stroke="rgba(255,255,255,0.05)" stroke-width="1"/></svg>');
  opacity: 0.8;
}

.profile-avatar-container {
  position: absolute;
  left: 50px;
  bottom: 5px;
  z-index: 2;
}

.profile-avatar {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  border: 5px solid var(--bg-neutral);
  background-color: var(--purple-color);
  background-size: cover;
  background-position: center;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

.profile-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  padding: 20px 50px 20px 220px;
  min-height: 100px;
}

.profile-user-info {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.profile-name {
  font-size: 28px;
  font-weight: 600;
  margin: 0;
  color: var(--color-white);
}

.profile-nickname {
  font-size: 18px;
  color: var(--purple-color);
  margin: 0;
}

.profile-email, .profile-dateOfBirth {
  font-size: 14px;
  color: #bbbbbb;
  margin: 0;
}

.profile-actions {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.profile-btn {
  transform: scale(1.1);
}

/* Profile Content */
.profile-content {
  display: grid;
  grid-template-columns: minmax(250px, 300px) 1fr;
  gap: 30px;
}

.profile-sidebar, .profile-main {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.profile-card {
  background-color: var(--bg-neutral);
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
}

.card-title {
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 20px;
  color: var(--purple-color);
  position: relative;
}

.card-title::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 0;
  width: 40px;
  height: 3px;
  background-color: var(--purple-color);
  border-radius: 3px;
}

.connections-stats {
  display: flex;
  justify-content: space-around;
  align-items: center;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-white);
}

.stat-label {
  font-size: 14px;
  color: #bbbbbb;
}

.stat-divider {
  width: 1px;
  height: 40px;
  background-color: rgba(255, 255, 255, 0.1);
}

.about-text {
  line-height: 1.6;
  color: #e0e0e0;
}

/* Edit Profile Section */
.toggle-edit-btn {
  width: 100%;
  padding: 12px;
  background-color: var(--purple-color);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.2s ease;
}

.toggle-edit-btn:hover {
  background-color: var(--hover-background-color);
}

.edit-profile-form {
  margin-top: 20px;
}

.edit-section {
  margin-bottom: 25px;
}

.edit-section label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  color: #bbbbbb;
}

.edit-input {
  width: 100%;
  padding: 12px;
  background-color: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: white;
  margin-bottom: 12px;
}

.edit-input:focus {
  border-color: var(--purple-color);
  outline: none;
}

.file-upload-wrapper {
  position: relative;
  margin-bottom: 12px;
}

.file-upload {
  position: absolute;
  top: 0;
  left: 0;
  opacity: 0;
  width: 100%;
  height: 100%;
  cursor: pointer;
}

.file-upload-text {
  display: block;
  padding: 12px;
  background-color: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #bbbbbb;
  text-align: center;
}

.edit-btn {
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
}

.edit-btn-primary {
  background-color: var(--purple-color);
  color: white;
}

.edit-btn-primary:hover {
  background-color: var(--hover-background-color);
}

.danger-zone {
  margin-top: 30px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding-top: 20px;
}

.danger-title {
  color: #ff5c5c;
  margin-bottom: 15px;
  font-size: 16px;
}

.edit-btn-danger {
  background-color: rgba(255, 87, 87, 0.1);
  color: #ff5c5c;
  border: 1px solid #ff5c5c;
}

.edit-btn-danger:hover {
  background-color: #ff5c5c;
  color: white;
}

/* Private Profile Message */
.private-profile-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 50px;
  background-color: var(--bg-neutral);
  border-radius: 12px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
  text-align: center;
  margin-top: 30px;
}

.lock-icon {
  font-size: 50px;
  margin-bottom: 20px;
  color: var(--purple-color);
}

.private-profile-message p {
  font-size: 18px;
  color: var(--color-white);
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background-color: #fff;
  width: 90%;
  max-width: 400px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.3);
  animation: modalFadeIn 0.3s ease;
}

@keyframes modalFadeIn {
  from { transform: scale(0.9); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

.modal-title {
  padding: 20px;
  margin: 0;
  color: #333;
  font-size: 18px;
  font-weight: 500;
  text-align: center;
  border-bottom: 1px solid #eee;
}

.modal-text {
  padding: 20px;
  color: #666;
  text-align: center;
}

.modal-buttons {
  display: flex;
  border-top: 1px solid #eee;
}

.modal-btn {
  flex: 1;
  padding: 15px;
  border: none;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.2s ease;
}

.modal-btn-cancel {
  background-color: #f8f8f8;
  color: #333;
}

.modal-btn-cancel:hover {
  background-color: #eee;
}

.modal-btn-confirm {
  background-color: #ff5c5c;
  color: white;
}

.modal-btn-confirm:hover {
  background-color: #e74c3c;
}

/* Responsive Design */
@media (max-width: 768px) {
  .profile-content {
    grid-template-columns: 1fr;
  }
  
  .profile-header-content {
    flex-direction: column;
    align-items: flex-start;
    padding: 20px 20px 20px 20px;
  }
  
  .profile-avatar-container {
    position: relative;
    left: 20px;
    bottom: -75px;
  }
  
  .profile-avatar {
    width: 100px;
    height: 100px;
  }
  
  .profile-actions {
    margin-top: 15px;
    margin-bottom: 0;
  }
}

.small {
  height: 1.4em;
  width: 1.4em;
  color: var(--color-white) !important;
  margin-bottom: 5px;
}
</style>
