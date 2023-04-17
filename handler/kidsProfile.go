package handler

import (
	"github.com/go-playground/validator/v10"
	"net/http"
	"story-time-server/dbhelper"
	"story-time-server/models"
	"story-time-server/utils"
)

func AddKidsProfile(w http.ResponseWriter, r *http.Request) {
	var addKidReq models.KidProfile
	if parseErr := utils.ParseBody(r.Body, &addKidReq); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "login error", "failed to parse request body")
		return
	}

	// validate request body
	validate := validator.New()
	err := validate.Struct(addKidReq)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "validation failed for request")
		return
	}

	// add kid's profile
	_, err = dbhelper.AddKidsProfile(addKidReq)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "unable to add kid's profile")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "Successfully Added Kid's Profile",
	})
}
