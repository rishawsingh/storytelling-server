package middlewares

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"story-time-server/database"
	"story-time-server/firebaseProvider"
	"story-time-server/models"
	"story-time-server/utils"
	"strings"
)

const (
	UserCtx = "user_context"
)

func FirebaseAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationToken := r.Header.Get("Authorization")
		idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
		if idToken == "" {
			utils.RespondError(w, http.StatusBadRequest, errors.New("empty authorization token"), "error authorizing")
			return
		}
		token, err := firebaseProvider.FirebaseClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, errors.New("invalid token"), "error authorizing")
			return
		}

		userID := token.Claims["id"].(int)
		var user models.User

		// language=SQL
		SQL := `
				SELECT 
				    id, 
				    name, 
				    year_of_birth, 
				    email, 
				    password, 
				    firebase_uid, 
				    avatar_id 
				FROM users
				WHERE
				    id = $1 
				  AND
				    firebase_uid = $2
				    `

		err = database.DB.Get(&user, SQL, userID, token.UID)
		if err != nil {
			logrus.Errorln("user not found ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context) *models.User {
	user, ok := ctx.Value(UserCtx).(*models.User)
	if !ok {
		return nil
	}
	return user
}
