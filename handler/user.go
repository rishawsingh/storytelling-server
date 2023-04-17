package handler

import (
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"net/http"
	"time"

	"story-time-server/dbhelper"
	"story-time-server/firebaseProvider"
	"story-time-server/models"
	"story-time-server/utils"

	"github.com/go-playground/validator/v10"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	//authorizationToken := r.Header.Get("Authorization")
	//if authorizationToken != "" {
	//	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	//	if idToken == "" {
	//		utils.RespondError(w, http.StatusBadRequest, errors.New("empty authorization token"), "error authorizing")
	//		return
	//	}
	//
	//	userIdToken, err := firebaseProvider.FirebaseClient.VerifyIDToken(context.Background(), idToken)
	//	if err != nil {
	//		utils.RespondError(w, http.StatusBadRequest, errors.New("invalid token"), "error authorizing")
	//		return
	//	}
	//
	//	user, err := firebaseProvider.FirebaseClient.GetUser(context.Background(), userIdToken.UID)
	//	if err != nil {
	//		utils.RespondError(w, http.StatusInternalServerError, err, "unable to get details.")
	//		return
	//	}
	//
	//}
	var newUserReq models.User
	if parseErr := utils.ParseBody(r.Body, &newUserReq); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	// validate the request newUserReq
	validate := validator.New()
	err := validate.Struct(newUserReq)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "validation failed for signup")
		return
	}

	// check if user is an adult
	currentYear := time.Now().Year()
	if (currentYear - newUserReq.YearOfBirth) < utils.MinAge {
		utils.RespondError(w, http.StatusInternalServerError, errors.New("not eligible to signup"), "not eligible to signup as you are less than 18 year's old")
		return
	}

	// check if user already exist
	exists, existsErr := dbhelper.IsUserExists(newUserReq.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}

	// if user exists return error
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exists")
		return
	}

	// hash the given password
	hashedPassword, hasErr := utils.HashPassword(newUserReq.Password)
	if hasErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
		return
	}

	// create user in firebase
	user := new(auth.UserToCreate)
	user.Email(newUserReq.Email)
	user.Password(newUserReq.Password)
	userRecord, fbErr := firebaseProvider.FirebaseClient.CreateUser(context.Background(), user)
	if fbErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, fbErr, "Failed to create firebase user")
		return
	}

	// create user
	userID, createErr := dbhelper.CreateUser(newUserReq, hashedPassword, userRecord.UID)
	if createErr != nil {
		// delete user from firebase
		_ = firebaseProvider.FirebaseClient.DeleteUser(context.Background(), userRecord.UID)
		utils.RespondError(w, http.StatusInternalServerError, createErr, "error registering new user")
		return
	}

	// create firebase token
	claims := map[string]interface{}{
		"id": userID,
	}
	token, err := firebaseProvider.FirebaseClient.CustomTokenWithClaims(context.Background(), userRecord.UID, claims)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to create token.")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginReq models.UserLoginReq
	if parseErr := utils.ParseBody(r.Body, &loginReq); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "login error", "failed to parse request body")
		return
	}

	// validate request body
	validate := validator.New()
	err := validate.Struct(loginReq)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "validation failed for login")
		return
	}

	// get user details
	user, getErr := dbhelper.GetUserByEmail(loginReq)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "Error fetching user details")
		return
	}

	// create firebase token
	claims := map[string]interface{}{
		"id": user.ID,
	}
	token, err := firebaseProvider.FirebaseClient.CustomTokenWithClaims(context.Background(), user.FireBaseUID, claims)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to create token.")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
