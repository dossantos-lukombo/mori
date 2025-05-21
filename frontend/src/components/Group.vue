<template>
    <div class="content" v-if="groupData">

        <div class="left-section">
            <GroupMembers v-bind:isMember="isMemberOfGroup" />
            
                 
        </div>

        <div class="middle-section">
            <div class="about">
                <h2 class="about-title">{{ this.groupData.name }}</h2>
                <p class="about-text">{{ this.groupData.description }}</p>
                <div v-if="isAdmin" class="admin-actions">
                    <button class="btn btn-danger" @click="openDeleteModal">
                        <i class="uil uil-trash-alt"></i> Delete Group
                    </button>
                </div>
            </div>            
            <p class="additional-info large" v-if="!this.isMemberOfGroup">Only group members can see additional
                information.
            </p>            
        </div>
      
        <GroupJoinRequests v-bind:isAdmin="this.isAdmin" class="right-section"/>  
     
        <!-- Delete Group Confirmation Modal -->
        <div v-if="showDeleteModal" class="modal-overlay" @click.self="cancelDelete">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 class="modal-title">Delete Group</h3>
                    <i class="uil uil-times close" @click="cancelDelete"></i>
                </div>
                <div class="modal-body">
                    <div class="modal-icon warning">
                        <i class="uil uil-exclamation-triangle"></i>
                    </div>
                    <p class="modal-message">Are you sure you want to delete this group? This action will remove all group messages and cannot be undone.</p>
                </div>
                <div class="modal-actions">
                    <button class="btn btn-secondary" @click="cancelDelete">Cancel</button>
                    <button class="btn btn-danger" @click="confirmDelete">
                        <i class="uil uil-trash-alt"></i> Delete Permanently
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import Groups from './Groups.vue';
import Notifications from './Notifications.vue';
import GroupMembers from './GroupMembers.vue';
import GroupJoinRequests from './GroupJoinRequests.vue';
import { useToast } from 'vue-toast-notification';

export default {
    name: "Group",
    created() {
        this.getGroupInfo();
    },
    setup() {
        const toast = useToast();
        return { toast };
    },
    watch: {
        $route() {
            if (this.$route.path.includes("group")){
                this.isMemberOfGroup=false;
                this.getGroupInfo(); 
            }
            // this.getGroupInfo()
        }
    },
    data() {
        return {
            groupData: null,
            isMemberOfGroup: false,
            isAdmin:false,
            showDeleteModal: false
        };
    },
    methods: {
        async getGroupInfo() {
            await fetch("http://localhost:8081/groupInfo?groupId=" + this.$route.params.id, {
                credentials: "include"
            })
                .then((r => r.json()))
                .then((json => {
                    // console.log("/groupInfo response", json);
                    this.groupData = json.groups[0];
                    if (json.groups[0].admin === true || json.groups[0].member === true) {
                        this.isMemberOfGroup = true
                    }
                    if(json.groups[0].admin === true){
                        this.isAdmin = true
                    }
                }));
        },
        
        openDeleteModal() {
            this.showDeleteModal = true;
        },
        
        cancelDelete() {
            this.showDeleteModal = false;
        },
        
        async confirmDelete() {
            try {
                const response = await fetch("http://localhost:8081/deleteGroup?groupId=" + this.$route.params.id, {
                    credentials: "include"
                });
                const json = await response.json();
                
                this.toast.open({
                    message: json.message,
                    type: "success",
                });
                
                // Close modal and redirect
                this.showDeleteModal = false;
                this.$router.push('/');
            } catch (error) {
                this.toast.open({
                    message: "Failed to delete group",
                    type: "error",
                });
                this.showDeleteModal = false;
            }
        }
    },
    components: { Groups, Notifications, GroupMembers, GroupJoinRequests }
}
</script>


<style scoped>
.content {
    margin-top: 50px;
    padding: 0 30px;
    display: grid;
    grid-template-columns: minmax(0, 1fr) minmax(0, 550px) minmax(0, 1fr);
    column-gap: 50px;
}



.middle-section {
    justify-self: center;
    display: flex;
    flex-direction: column;
    gap: 35px;


}


.left-section {
    justify-self: flex-end;
    display: flex;
    flex-direction: column;
    gap: 35px;
}

.right-section {
    justify-self: flex-start;
    min-width: 250px;
  
}

.admin-actions {
    margin-top: 15px;
}

.btn-danger {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    background-color: #dc3545;
    color: white;
    transition: all 0.3s ease;
}

.btn-danger:hover {
    background-color: #c82333;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(220, 53, 69, 0.3);
}

.btn-secondary {
    background-color: transparent;
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: var(--color-white);
}

.btn-secondary:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.btn-danger i, 
.btn-secondary i {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.2em;
    line-height: 1;
}

.modal-actions .btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
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
    margin-bottom: 15px;
}

.modal-icon.warning {
    color: #ffc107;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 15px;
}

.modal-icon.warning i {
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
    padding: 15px 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
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

.box {
    height: 300px;
    width: 550px;
    border: 2px solid blue;
}

@media only screen and (max-width: 1250px) {
    .content {
        grid-template-columns: minmax(min-content, max-content) minmax(min-content, 550px);
        grid-template-rows: repeat(2, minmax(auto, max-content));
        row-gap: 35px;
        justify-content: center;
        grid-template-areas:
            "left-section middle-section"
            "right-section middle-section"
            "... middle-section";

    }


    .middle-section {
        grid-area: middle-section;
    }

    .left-section {
        grid-area: left-section;
    }

    .right-section {
        grid-area: right-section;
    }

}
</style>