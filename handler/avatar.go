package handler

import (
	"net/http"
	"story-time-server/dbhelper"
	"story-time-server/utils"
)

func GetAvatars(w http.ResponseWriter, r *http.Request) {
	// get avatars
	avatars, err := dbhelper.GetAvatars()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get avatars")
		return
	}

	utils.RespondJSON(w, http.StatusOK, avatars)
}
