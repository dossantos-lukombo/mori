<template>
    <div v-if="user && this.$store.state.id !== ''" class="user-profil">
        <div id="layout-profile">

            <div class="left-section ">
                <div class="user-profile__public">
                    <div class="user-picture" :style="{ backgroundImage: `url(http://localhost:8081/${user.avatar})` }">
                    </div>
                    <div class="user-profile__info">
                        <h2 class="username">{{user.firstName}} {{user.lastName}}</h2>
                        <h3 v-if="showNickname" class="username">{{ user.nickname }}</h3>
                        <p class="user-email" v-if="user.login">{{ user.login }}</p>
                        <p class="user-dateOfBirth" v-if="user.dateOfBirth">{{ user.dateOfBirth }}</p>
                    </div>

                    <div class="profile-btns">

                        <PrivacyBtn v-if="isMyProfile" :status="user.status" />

                        <!-- Follow/unfollow button -->
                        <component v-else :is="displayBtn" v-bind="{ user }" @follow="checkFollowRequest" @unfollow="unfollow"></component>

                    </div>

                </div>
                <div class="multiple-item-list" v-if="showProfileData">
                    <Following :following="following"/>
                    <Followers :followers="followers" />
                </div>

                <Groups :groups="profileGroups" v-if="showProfileData"/>


            </div>

            <div class="middle-section" v-if="showProfileData">

                <div class="about" v-if="user.about !== ''">
                    <h2 class="about-title">About me</h2>
                    <p class="about-text">{{ user.about }}</p>
                </div>

                <div v-if="isMyProfile" class="edit-profile">
                    <h2 class="edit-profile-title">Modifier le profil</h2>

                    <!-- Changer le pseudo -->
                    <div class="edit-section">
                        <h3 class="edit-section-title">Nouveau pseudo</h3>
                        <input
                        v-model="newNickname"
                        type="text"
                        placeholder="Entrez votre nouveau pseudo"
                        class="input-field"
                    />
                        <button @click="changeNickname" class="btn-primary">Modifier le pseudo</button>
                    </div>
                

                    <!-- Changer l'avatar -->
                    <div class="edit-section">
                        <h3 class="edit-section-title"> Nouvel avatar</h3>
                        <input type="file" ref="avatarInput" class="file-input" accept="image/*" />
                        <button @click="changeAvatar" class="btn-primary">Modifier l'avatar</button>
                    </div>
                </div>
            </div>

            <p v-else class="additional-info large"> This profile is private</p>

        </div>
    </div>

</template>

<script>

import Following from './Following.vue'
import Followers from './Followers.vue'
import FollowBtn from './FollowBtn.vue'
import PrivacyBtn from './PrivacyBtn.vue'
import UnfollowBtn from './UnfollowBtn.vue'
import Groups from './Groups.vue'
export default {
    name: 'Profile',
    components: { Followers, Following, FollowBtn, PrivacyBtn, UnfollowBtn, Groups },
    data() {
        return {
            // flag: false,
            user: null,
            isMyProfile: false,
            followers: [],
            following: [],
            

            profileGroups: null,
            newNickname: "", // Nouveau pseudo
            selectedAvatar: null, // Fichier pour l'avatar 
        }
    },

    created() {
        this.updateProfileData()
    },
    computed: {
        showProfileData() {
            return (this.user.following || this.isMyProfile || this.user.status === "PUBLIC") ? true : false
        },
        showSendButton() {
            return !this.isMyProfile && this.user.status === "PUBLIC" && !this.user.following
        },

        displayBtn() {
            if (this.user.following) {
                return UnfollowBtn
            } else {
                return FollowBtn
            }
        },
        showNickname(){
            if (this.user.nickname == this.user.firstName + " "+ this.user.lastName){
                return false
            }
            return true
        }
    },
    methods: {
        updateProfileData(){
            this.getUserData()           
            this.getFollowers()
            this.getFollowing()
            this.checkProfile()
            this.getProfileGroups()            
        },
        async getUserData() {
            await fetch("http://localhost:8081/userData?userId=" + this.$route.params.id, {
                credentials: "include",
            })
                .then((r) => r.json())
                .then((json) => {
                    console.log("/getUserData", json)
                    console.log("id", this.$route.params.id)
                    this.user = json.users[0];
                });

        },

        async getProfileGroups() {
            const response = await fetch("http://localhost:8081/otherUserGroups?userId=" + this.$route.params.id, {
                credentials: 'include',
            });
            const data = await response.json();

            if (data.type == "Error") {
                console.log("/getProfileGroups error: ", data.message)
            } else {
                this.profileGroups = data.groups;
            }
        },

        async getMyUserID() {
            await this.$store.dispatch("getMyUserID")
        },

        async checkProfile() {
            await this.getMyUserID();
            const myID = this.$store.state.id;
            const profileID = this.$route.params.id;
            this.isMyProfile = (profileID === myID)
        },


        checkFollowRequest(action) {
            if (action === "followedUser") {
                this.$store.dispatch("fetchChatUserList");
                this.updateProfileData()
                this.toggleFollowingThisUser();

            }
        },

        toggleFollowingThisUser() {
            this.user.following = !this.user.following
        },
        unfollow(){
            this.updateProfileData();
            this.$store.dispatch("fetchChatUserList");
        },
        async getFollowers() {
            await fetch('http://localhost:8081/followers?userId=' + this.$route.params.id, {
                credentials: 'include'
            })
                .then((response => response.json()))
                .then((json => {
                    this.followers = json.users
                }))
        },
        async getFollowing() {
            await fetch('http://localhost:8081/following?userId=' + this.$route.params.id, {
                credentials: 'include'
            })
                .then((response => response.json()))      
                .then((json => {
                    this.following = json.users
                }))

        },
        async changeNickname() {
            await this.$store.dispatch("changeNickname", this.newNickname);
            this.updateProfileData();
        },

        async changeAvatar() {
            const input = this.$refs.avatarInput;
            if (!input || !input.files || input.files.length === 0) {
                alert("Veuillez sélectionner un fichier.");
                return;
            }
            this.selectedAvatar = input.files[0];
            await this.$store.dispatch("changeAvatar", this.selectedAvatar);
            this.updateProfileData();
        },

     
        
        addChat() {
            // check if user doesnt have a chat with that person already
            // ....

            let newChat = {
                name: this.user.nickname,
                receiverId: this.user.id,
                type: "PERSON"
            };

            this.$store.dispatch("addNewChat", newChat);
        }
    },
    watch: { //watching changes in route
        $route() {
            if (this.$route.name === "Profile") {
                this.updateProfileData()
            }

        }
    }
}
</script>

<style scoped>

.user-profile {
    overflow: scroll;
}

#layout-profile {
    display: grid;
    grid-template-columns: 1fr minmax(min-content, 550px) 1fr;
    column-gap: 50px;
    margin-top: 50px;
    /* justify-items: center; */ 

}

.middle-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 50px;
}


.left-section {
    display: flex;
    flex-direction: column;
    gap: 35px;
    max-width: 250px;
    justify-self: flex-end;

}


.user-profile__public,
.user-profile__private {
    display: flex;
    flex-direction: column;
    padding: var(--container-padding);
    background-color: var(--bg-neutral);
    box-shadow: 0 2px 10px rgb(0, 0, 0);
    border-radius: var(--container-border-radius);
    align-items: center;
    text-align: center;
    gap: 25px;


}

.user-profile__private p, h3, .user-profile__public p, h3{
    color: var(--color-white);
}

.user-profile__public h2, .user-profile__private h2 {
    color: var(--purple-color);
}

.user-profile__privacy{
    color: var(--color-white);
}

:is(.user-profile__public, .user-profile__private) .user-picture {
    height: 185px;
    width: 185px;
}

.user-profile__info {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.profile-btns {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
}

.additional-info {
    text-align: center;
}

.edit-profile {
  margin-top: 20px;
  padding: 20px;
  background-color: var(--bg-neutral);
  border-radius: var(--container-border-radius);
  box-shadow: 0 2px 10px rgb(0, 0, 0);
  width: 550px;
}

.edit-section {
  margin-bottom: 20px;
}

.input-field {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin-bottom: 10px;
}

.file-input {
  margin-bottom: 10px;
}

.btn-primary {
  padding: 10px 20px;
  background-color: var(--purple-color);
  color: #fff;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.btn-primary:hover {
  background-color: var(--purple-hover-color);
}

.edit-profile-title {
    text-align: center;
    color: var(--purple-color);
    margin-bottom: 20px;
}

.edit-section-title {
    margin-bottom: 15px;    
}
</style>