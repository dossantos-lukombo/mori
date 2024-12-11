package main

import (
	"fmt"
	"net/http"

	sqlite "mori/pkg/db/PostgreSql"
	"mori/pkg/handlers"
	"mori/pkg/utils"
	ws "mori/pkg/wsServer"
)

func main() {
	// initialize database
	db := sqlite.InitDB()
	defer db.Close()
	// temp
	// initialize repositories
	repos := sqlite.InitRepositories(db)
	// initialize handlers with connection to repositories
	handler := handlers.InitHandlers(repos)
	// initialize wsServer
	wsServer := ws.StartServer(repos)

	// set up server address and routes
	server := &http.Server{
		Addr:    ":8081",
		Handler: setRoutes(handler, wsServer),
	}

	fmt.Printf("Server started at http://localhost" + server.Addr + "\n")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error", err)
	}
}

// Set up all routes
func setRoutes(handler *handlers.Handler, wsServer *ws.Server) http.Handler {
	mux := http.NewServeMux()
	/* ------------------------------ image server ------------------------------ */
	fs := http.FileServer(http.Dir("./imageUpload"))
	mux.Handle("/imageUpload/", http.StripPrefix("/imageUpload/", utils.ConfigFSHeader(fs)))
	/* ------------------------------- auth route ------------------------------- */
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/signin", handler.Signin)
	mux.HandleFunc("/logout", handler.Auth(handler.Logout))
	mux.HandleFunc("/sessionActive", handler.SessionActive)

	/* ---------------------------------- users --------------------------------- */
	mux.HandleFunc("/allUsers", handler.Auth(handler.AllUsers))       // all users + info except current
	mux.HandleFunc("/followers", handler.Auth(handler.GetFollowers))  // follower list
	mux.HandleFunc("/following", handler.Auth(handler.GetFollowing))  // following list
	mux.HandleFunc("/currentUser", handler.Auth(handler.CurrentUser)) // current user data
	mux.HandleFunc("/userData", handler.Auth(handler.UserData))       // userd data based on following status
	mux.HandleFunc("/changeStatus", handler.Auth(handler.UserStatus)) // change status

	mux.HandleFunc("/follow", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.Follow(wsServer, w, r)
	})) // follow user
	mux.HandleFunc("/cancelFollowRequest", handler.Auth(handler.CancelFollowRequest))
	mux.HandleFunc("/unfollow", handler.Auth(handler.Unfollow))
	mux.HandleFunc("/responseFollowRequest", handler.Auth(handler.ResponseFollowRequest))	

	/* --------------------------------- groups --------------------------------- */
	mux.HandleFunc("/allGroups", handler.Auth(handler.AllGroups))             // group list
	mux.HandleFunc("/userGroups", handler.Auth(handler.UserGroups))           // group list of user groups
	mux.HandleFunc("/otherUserGroups", handler.Auth(handler.OtherUserGroups)) // group list for specific user

	mux.HandleFunc("/groupInfo", handler.Auth(handler.GroupInfo))                     // get group info
	mux.HandleFunc("/groupMembers", handler.Auth(handler.GroupMembers))               // get group members
	
	mux.HandleFunc("/groupRequests", handler.Auth(handler.GroupRequests))             // get group member requests
	mux.HandleFunc("/cancelGroupRequests", handler.Auth(handler.CancelGroupRequests)) // cancel request or joing group

	mux.HandleFunc("/newGroup", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewGroup(wsServer, w, r)
	})) // create new group
	
	mux.HandleFunc("/newGroupInvite", handler.Auth(func(w http.ResponseWriter, r *http.Request) { // invite new users to group
		handler.NewGroupInvite(wsServer, w, r)
	}))
	mux.HandleFunc("/newGroupRequest", handler.Auth(func(w http.ResponseWriter, r *http.Request) { // invite new users to group
		handler.NewGroupRequest(wsServer, w, r)
	}))
	mux.HandleFunc("/responseGroupRequest", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.ResponseGroupRequest(wsServer, w, r)
	})) // response to join request
	mux.HandleFunc("/responseInviteRequest", handler.Auth(handler.ResponseInviteRequest)) // response to invite request
	

	/* ------------------------------ notifications ----------------------------- */
	mux.HandleFunc("/notifications", handler.Auth(handler.Notifications)) // get all notifs from db on login

	/* ------------------------------ chat messages ----------------------------- */
	mux.HandleFunc("/messages", handler.Auth(handler.Messages))             // get all chat messages for specific chat
	mux.HandleFunc("/unreadMessages", handler.Auth(handler.UnreadMessages)) // get list of messages that isn't read
	mux.HandleFunc("/messageRead", handler.Auth(handler.MessageRead))       // mark message as read
	mux.HandleFunc("/newMessage", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.NewMessage(wsServer, w, r)
	})) // new chat message
	mux.HandleFunc("/chatList", handler.Auth(handler.ChatList))                       // get list of users to display in chatbox
	mux.HandleFunc("/responseChatRequest", handler.Auth(handler.ResponseChatRequest)) // response to chat request

	/* ---------------------------- websocket server ---------------------------- */
	mux.HandleFunc("/ws", handler.Auth(func(w http.ResponseWriter, r *http.Request) {
		handler.SocketHandler(wsServer, w, r)
	}))

	return mux
}