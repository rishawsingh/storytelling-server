package dbhelper

import (
	"database/sql"
	"story-time-server/database"
	"story-time-server/models"
	"story-time-server/utils"
)

func IsUserExists(email string) (bool, error) {
	// language=SQL
	SQL := `
			SELECT 
			    id 
			FROM 
			    users 
			WHERE 
			    email = TRIM(LOWER($1)) 
			  AND 
			    archived_at IS NULL
			    `

	var id string
	err := database.DB.Get(&id, SQL, email)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func CreateUser(user models.User, hashPwd, firebaseUID string) (int, error) {
	// language=SQL
	SQL := `
			INSERT INTO 
			    users(name, year_of_birth, email, password, firebase_uid, avatar_id) 
			VALUES 
			    ($1, $2, TRIM(LOWER($3)), $4, $5, $6) 
			RETURNING id`
	var userID int
	if err := database.DB.QueryRowx(SQL, user.Name, user.YearOfBirth, user.Email, hashPwd, firebaseUID, user.AvatarID).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

func GetUserByEmail(loginReq models.UserLoginReq) (models.User, error) {
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
       		FROM
				users 
			WHERE
				archived_at IS NULL 
			  AND 
			    email = TRIM(LOWER($1))
				`

	var user models.User
	err := database.DB.Get(&user, SQL, loginReq.Email)
	if err != nil && err != sql.ErrNoRows {
		return models.User{}, err
	}
	if err == sql.ErrNoRows {
		return models.User{}, nil
	}

	// compare password
	if passwordErr := utils.CheckPassword(loginReq.Password, user.Password); passwordErr != nil {
		return models.User{}, passwordErr
	}
	return user, nil
}
