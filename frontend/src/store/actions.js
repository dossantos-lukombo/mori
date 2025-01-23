export default {
   

    async getMyUserID({ commit }) {

        await fetch("http://localhost:8081/currentUser", {
            credentials: "include",
        })
            .then((r) => r.json())
            .then((json) => {
                //console.log("JSON response", json)
                commit("updateMyUserID", json.users[0].id)
            });
    },

    async getMyProfileInfo(context) {
        await context.dispatch("getMyUserID");
        const userID = context.state.id;
        await fetch("http://localhost:8081/userData?userId=" + userID, {
            credentials: "include",
        })
            .then((r) => r.json())
            .then((json) => {
                let userInfo = json.users[0];
                // console.log(userInfo);
                this.commit("updateProfileInfo", userInfo);
                // console.log("userinfo -", json);
            });
    },

    async getAllUsers() {
        await fetch("http://localhost:8081/allUsers", {
            credentials: "include",
        })
            .then((r) => r.json())
            .then((json) => {
                let users = json.users;
                this.commit("updateAllUsers", users);
                // console.log("allUsers:", json.users);
            });
    },
    async getAllGroups() {
        await fetch("http://localhost:8081/allGroups", {
            credentials: "include",
        })
            .then((r) => r.json())
            .then((json) => {
                let groups = json.groups;
                this.commit("updateAllGroups", groups);
                //console.log("Allgroups:", json.groups);
            });
    },

    async getUserGroups(context) {
        const response = await fetch(`http://localhost:8081/userGroups`, {
            credentials: 'include'
        });

        const data = await response.json();
        // console.log("/getUserGroups data", data)
        // context.state.groups.userGroups.loaded = true;

        context.commit("updateUserGroups", data.groups)
        context.commit("updateDataLoaded", "userGroups")

    },

    addUserGroup({ state, commit }, userGroup) {
        let userGroups = state.groups.userGroups;
        console.log("userGroups state", userGroups)
        if (userGroups === null) { userGroups = [] };
        userGroups.push(userGroup);

        console.log("userGroup", userGroup)
        commit("updateUserGroups", userGroups)
    },

    async getMyFollowers(context) {
        await context.dispatch("getMyProfileInfo");
        const myID = context.state.profileInfo.id;

        const response = await fetch(`http://localhost:8081/followers?userId=${myID}`, {
            credentials: 'include'
        });

        const data = await response.json();

        context.commit("updateMyFollowers", data.users)

    },
    

    

    async isLoggedIn() {
        const response = await fetch('http://localhost:8081/sessionActive', {
            credentials: 'include'
        });

        const data = await response.json();

        if (data.message === "Session active") {
            // console.log("ah yes")
            return true
        } else {
            // console.log("ah no")
            return false
        }

    },

    async fetchConversationsMsg({ commit }) {
        const resp = await fetch("http://localhost:8081/conversationsMsg", {
          credentials: "include"
        });
        const data = await resp.json();
      
        // data devrait ressembler à :
        // { conversationsMsg: [ { ... }, { ... } ] }
        commit("setConversationsMsg", data.conversationsMsg);
    },

    async changeNickname({ dispatch }, newNickname) {
        if (!newNickname.trim()) {
            alert("Veuillez entrer un pseudo valide.");
            return;
        }
    
        try {
            const response = await fetch("http://localhost:8081/updateNickname", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify({ nickname: newNickname }),
            });
    
            const result = await response.json();
            if (!response.ok) {
                throw new Error(result.message || "Erreur lors de la modification du pseudo");
            }
    
            alert("Pseudo modifié avec succès !");
            await dispatch("getMyProfileInfo");
        } catch (err) {
            console.error("Erreur lors de la modification du pseudo :", err.message);
            alert(err.message);
        }
    },
    
      /**
       * Action pour changer l'avatar de l'utilisateur
       */
      async changeAvatar({ dispatch }, avatarFile) {
        if (!avatarFile) {
            alert("Veuillez sélectionner un avatar.");
            return;
        }
        const formData = new FormData();
        formData.append("avatar", avatarFile);
        try {
            const response = await fetch("http://localhost:8081/updateAvatar", {
                method: "POST",
                credentials: "include",
                body: formData,
            });
            if (!response.ok) {
                throw new Error("Erreur lors de la modification de l'avatar.");
            }
            alert("Avatar modifié avec succès !");
            await dispatch("getMyProfileInfo");
        } catch (err) {
            console.error(err);
            alert(err.message);
        }
    },
    

    
      

    createWebSocketConn({ commit, dispatch, state }) {
        const ws = new WebSocket("ws://localhost:8081/ws");
      
        ws.addEventListener("message", (e) => {
            const data = JSON.parse(e.data);
            if (data.action == "chat") {
                // only broadcast messages when participants(sender and reciever) chat is open

                const isParticipantsChatOpen = state.chat.openChats.some((chat) => {
                    // Chat is open with the person who sent the message
                    if (data.chatMessage.type === "PERSON" && data.chatMessage.senderId  === chat.receiverId) {
                        return true
                    }

                    if (data.chatMessage.type === "GROUP" && data.chatMessage.receiverId === chat.receiverId) {
                        return true
                    }


                })
                if (isParticipantsChatOpen) {
                    dispatch("addNewChatMessage", data.chatMessage)
                    dispatch("markMessageRead", data.chatMessage)
                } else {
                    if (data.message === "NEW") {
                        dispatch("fetchChatUserList");
                    }
    
                    dispatch("addUnreadChatMessage", data.chatMessage)
                }
            } else if (data.action == "notification") {
                dispatch("addNewNotification", data.notification);

            } else if(data.action == "groupAccept"){
                dispatch("getUserGroups");
            }

        })

        commit("updateWebSocketConn", ws)  
    }

}