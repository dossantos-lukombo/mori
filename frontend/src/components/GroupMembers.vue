<template>
    <div class="item-list__wrapper" id="groups">
        <h3>Members</h3>
        <ul class="item-list">
            <li v-for="member in this.groupMembers" :key="member.id">
                <div class="profile-image-container">
                    <img class="profile-image" src="../assets/icons/default-profile.svg" alt="Profile image">
                </div>
                <div class="item-text">
                    <router-link :to="{ path: `/profile/${member.id}`}">{{ member.nickname }}</router-link>
                </div>
            </li>
        </ul>
        <button v-if="this.isMember" class="btn form-submit" @click="toggleModal();getFollowers()">
            <i class="uil uil-user-plus"></i> Invite users
        </button>
        <button v-if="!this.isMember" class="btn form-submit join-btn" @click="this.joinGroup">
            <i class="uil uil-user-check"></i> Join group
        </button>
        <button v-if="this.isMember && !this.isAdmin" class="btn form-submit btn-danger" @click="openLeaveModal" title="Your messages in this group will remain with your name">
            <i class="uil uil-signout"></i> Leave group
        </button>

        <!-- Invite users modal -->
        <Modal v-if="this.isOpen" @closeModal="toggleModal">
            <template #title>Invite users</template>
            <template #body>
                <div class="invite-form">
                    <MultiselectDropdown v-model:checkedOptions="checkedNames" placeholder="Select followers"
                        :content="listForShowing" />
                    <div class="modal-actions">
                        <button class="btn btn-secondary" @click="toggleModal">Cancel</button>
                        <button class="btn form-submit" @click="toggleModal() ; inviteUsersToGroup()">
                            <i class="uil uil-user-plus"></i> Invite
                        </button>
                    </div>
                </div>
            </template>
        </Modal>

        <!-- Leave group confirmation modal -->
        <div v-if="showLeaveModal" class="modal-overlay" @click.self="cancelLeave">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 class="modal-title">Leave Group</h3>
                    <i class="uil uil-times close" @click="cancelLeave"></i>
                </div>
                <div class="modal-body">
                    <div class="modal-icon">
                        <i class="uil uil-exclamation-triangle"></i>
                    </div>
                    <p class="modal-message">Are you sure you want to leave this group? Your messages will remain visible to other members.</p>
                </div>
                <div class="modal-actions">
                    <button class="btn btn-secondary" @click="cancelLeave">Cancel</button>
                    <button class="btn btn-danger" @click="confirmLeave">
                        <i class="uil uil-signout"></i> Leave Group
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>



<script>
import Modal from './Modal.vue';
import MultiselectDropdown from './MultiselectDropdown.vue';
import { useToast } from 'vue-toast-notification';

export default {
    name: "GroupMembers",
    props: {
        isMember: false
    },
    setup() {
        const toast = useToast();
        return { toast };
    },
    data() {
        return {
            groupMembers: null,
            isOpen: false,
            followers: [],
            listForShowing: [],
            allUsers: [],
            checkedNames: [],
            clearInput: false,
            isAdmin: false,
            showLeaveModal: false
        };
    },
    created() {
        this.getGroupMembers();
        this.checkIfAdmin();
    },

    computed: {
        allFollowersNames() {
            return this.listForShowing
        }
    },

    watch: {
        $route() {
            this.getGroupMembers();
            this.checkIfAdmin();
        }
    },
  
    methods: {
        async getFollowers() {
            this.$store.dispatch("getMyFollowers");
            this.createFollowersListForShowing(this.$store.state.myFollowers, this.groupMembers)
        },
        async getGroupMembers() {
            await fetch("http://localhost:8081/groupMembers?groupId=" + this.$route.params.id, {
                credentials: "include"
            })
                .then((response => response.json()))
                .then((json => {
                    // console.log("GroupMembers:", json);
                    this.groupMembers = json.users;
                }));
        },
        createFollowersListForShowing(followers, members) {
            this.listForShowing = [];
            let isUserInGroup = false
            for (let i = 0; i < Object.keys(followers).length; i++) {
                for (let j = 0; j < Object.keys(members).length; j++) {
                    if (followers[i].nickname === members[j].nickname) {
                        isUserInGroup = true
                    }
                }
                if (!isUserInGroup) {
                    this.listForShowing.push(followers[i])
                }
                isUserInGroup = false
            }
        },
        toggleModal() {
            this.isOpen = !this.isOpen;
        },

        getIds() {
            let arrOfIDS = [];
            for (let name of this.checkedNames) {
                for (let obj of this.listForShowing) {
                    if (obj.nickname === name.nickname) {
                        arrOfIDS.push(obj.id)
                    }
                }
            }
            return arrOfIDS
        },
        
        async joinGroup(){
            await fetch("http://localhost:8081/newGroupRequest?groupId=" + this.$route.params.id, {
                credentials: 'include',
            })
            .then(response=>response.json())
            .then(json=>{
                this.toast.open({
                    message: json.message,
                    type: "success",
                });
            })
        },

        openLeaveModal() {
            this.showLeaveModal = true;
        },

        cancelLeave() {
            this.showLeaveModal = false;
        },

        async confirmLeave() {
            try {
                const response = await fetch("http://localhost:8081/leaveGroup?groupId=" + this.$route.params.id, {
                    credentials: 'include',
                });
                const json = await response.json();
                
                this.toast.open({
                    message: json.message,
                    type: "success",
                });
                
                // Close modal and redirect
                this.showLeaveModal = false;
                this.$router.push('/');
                
            } catch (error) {
                this.toast.open({
                    message: "Failed to leave group",
                    type: "error",
                });
                this.showLeaveModal = false;
            }
        },

        async checkIfAdmin() {
            await fetch("http://localhost:8081/groupInfo?groupId=" + this.$route.params.id, {
                credentials: "include"
            })
            .then(response => response.json())
            .then(json => {
                if (json.groups && json.groups[0]) {
                    this.isAdmin = json.groups[0].admin === true;
                }
            });
        },

        async inviteUsersToGroup() {
            await fetch("http://localhost:8081/newGroupInvite", {
                method: 'POST',
                credentials: 'include',
                body: JSON.stringify({ invitations: this.getIds(), id: this.$route.params.id })
            })
                .then((response => response.json()))
                .then((json => {
                    this.toast.open({
                        message: "Invitations sent successfully",
                        type: "success",
                    });
                    this.clearInput = true;
                }));
        },
    },
    components: { Modal, MultiselectDropdown }
}
</script>


<style>
.btn-danger {
    background-color: #dc3545;
    color: white;
    margin-top: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.btn-danger:hover {
    background-color: #c82333;
    transform: translateY(-2px);
}

.btn-secondary {
    background-color: transparent;
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: var(--color-white);
}

.btn-secondary:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.form-submit {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    transition: all 0.3s ease;
}

.form-submit i {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.2em;
    line-height: 1;
}

.form-submit:hover {
    transform: translateY(-2px);
}

.join-btn {
    background-color: #28a745;
}

.join-btn:hover {
    background-color: #218838;
}

.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    animation: fadeIn 0.3s ease-out;
}

.modal-content {
    background-color: var(--bg-neutral);
    border-radius: 10px;
    width: 90%;
    max-width: 450px;
    overflow: hidden;
    box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
    animation: slideUp 0.3s ease-out;
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-title {
    margin: 0;
    font-size: 1.3rem;
    color: var(--color-white);
}

.modal-body {
    padding: 20px;
    text-align: center;
}

.modal-icon {
    font-size: 3rem;
    color: #ffc107;
    margin-bottom: 15px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 15px;
}

.modal-icon i {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 3rem;
    height: 3rem;
}

.modal-message {
    color: var(--color-white);
    margin-bottom: 20px;
    line-height: 1.5;
}

.modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 20px;
}

.modal-actions .btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.modal-actions .btn i {
    display: flex;
    align-items: center;
    justify-content: center;
    line-height: 1;
}

.invite-form {
    display: flex;
    flex-direction: column;
    width: 100%;
}

.close {
    cursor: pointer;
    font-size: 1.5rem;
    color: var(--color-white);
    opacity: 0.7;
    transition: opacity 0.2s ease;
}

.close:hover {
    opacity: 1;
}

@keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
}

@keyframes slideUp {
    from { transform: translateY(30px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
}

.profile-image-container {
    width: 40px;
    height: 40px;
    min-width: 40px;  /* Prevent shrinking */
    min-height: 40px; /* Prevent shrinking */
    overflow: hidden;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--bg-neutral);
    border: 2px solid var(--purple-color);
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
    position: relative; /* For absolute positioning of child */
}

.profile-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block; /* Force block display for Chrome */
    -webkit-user-drag: none; /* Prevent drag in webkit */
    position: absolute; /* Force position */
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%); /* Center exactly */
}

/* Make sure item-list rows are aligned properly */
.item-list li {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 0;
}
</style>