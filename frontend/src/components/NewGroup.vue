<template>
    <button class="btn create-group-btn" @click="toggleModal">
      <i class="uil uil-users-alt"></i>
      <span>New group</span>
      <i class="uil uil-plus"></i>
    </button>

    <Modal v-show="isOpen" @closeModal="toggleModal(); toggleClearInput();">
        <template #title>
            Create new group
        </template>

        <template #body>
            <form @submit.prevent="submitNewGroup" ref="theForm" class="create-group-form">
                <div class="form-input">
                    <label for="name">Group Name</label>
                    <input type="text" name="name" id="name" placeholder="Enter group name" required>
                </div>

                <div class="form-input">
                    <label for="description">Description</label>
                    <textarea id="description"
                              name="description"
                              rows="4"
                              cols="50"
                              required
                              placeholder="What is this group about?"></textarea>
                </div>
    
                <MultiselectDropdown
                    v-model:checkedOptions="checkedFollowers"
                    :content="getMyFollowersList"
                    label-name="Invite members"
                    placeholder="Select followers to invite" />

                <div class="form-actions">
                  <button type="button" class="btn btn-cancel" @click="toggleModal">Cancel</button>
                  <button class="btn form-submit" type="submit">
                    <i class="uil uil-users-alt"></i>
                    Create Group
                  </button>
                </div>
            </form>
        </template>
    </Modal>
</template>


<script>
import Modal from "@/components/Modal.vue";
import MultiselectDropdown from "./MultiselectDropdown.vue";
import { useToast } from 'vue-toast-notification';

export default {
  components: {
    Modal,
    MultiselectDropdown,
  },
  setup() {
    const toast = useToast();
    return { toast };
  },
  data() {
    return {
      checkedFollowers: [],
      isOpen: false,
      clearInput: false,
    };
  },
  created() {
    this.getMyFollowers();
  },
  computed: {
    getMyFollowersList() {
      // Return the list of followers from Vuex
      return this.$store.getters.followers;
    },
  },
  methods: {
    async getMyFollowers() {
      // Ensure this action properly fills the Vuex state
      await this.$store.dispatch("getMyFollowers");
    },
    toggleModal() {
      if (this.isOpen) {
        this.$refs.theForm.reset();
        this.checkedFollowers = [];
      }
      this.isOpen = !this.isOpen;
    },
    toggleClearInput() {
      this.clearInput = !this.clearInput;
    },
    async submitNewGroup(e) {
      const form = e.currentTarget;
      const formData = new FormData(form);
      const formDataObject = Object.fromEntries(formData.entries());
      formDataObject["invitations"] = this.getIds(this.checkedFollowers);

      try {
        const response = await fetch("http://localhost:8081/newGroup", {
          method: "POST",
          credentials: "include",
          body: JSON.stringify(formDataObject),
        });

        const json = await response.json();
        this.toast.open({
          message: "Group created successfully!",
          type: "success",
        });

        form.reset();
        this.toggleModal();
        this.toggleClearInput();
        this.$store.dispatch("getUserGroups");
      } catch (error) {
        this.toast.open({
          message: "Error creating group",
          type: "error",
        });
      }
    },
    getIds() {
      const arrOfIDS = [];
      for (const name of this.checkedFollowers) {
        for (const obj of this.$store.state.myFollowers) {
          if (obj.nickname === name.nickname) {
            arrOfIDS.push(obj.id);
          }
        }
      }
      return arrOfIDS;
    },
  },
};
</script>


<style scoped>
.create-group-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 15px;
  padding: 10px 15px;
  transition: all 0.3s ease;
}

.create-group-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.create-group-btn i {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2em;
  line-height: 1;
}

.create-group-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

.form-input label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  text-align: left;
}

.form-input input,
.form-input textarea {
  width: 100%;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background-color: rgba(255, 255, 255, 0.05);
  color: var(--color-white);
  font-size: 1em;
  transition: all 0.3s ease;
}

.form-input input:focus,
.form-input textarea:focus {
  outline: none;
  border-color: var(--purple-color);
  box-shadow: 0 0 0 2px rgba(145, 70, 188, 0.3);
}

.form-input textarea {
  min-height: 100px;
  resize: vertical;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 10px;
}

.btn-cancel {
  background-color: transparent;
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: var(--color-white);
}

.btn-cancel:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.form-submit {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 20px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.form-submit:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(145, 70, 188, 0.3);
}

.form-submit i {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2em;
  line-height: 1;
}
</style>


