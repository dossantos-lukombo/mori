package handlers

import (
	"fmt"
	"mori/pkg/models"
	"mori/pkg/utils"
	"net/http"
)

// handler for logout/ validate user by id and delete session
func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// delete session
	session := models.Session{UserID: userId}
	errSession := handler.repos.SessionRepo.Delete(session)
	if errSession != nil {
		fmt.Println("error on deleting session", errSession)
		return
	}
	// delete cookie
	utils.DeleteCookie(w)
	utils.RespondWithSuccess(w, "Logout successful", 200)
}
