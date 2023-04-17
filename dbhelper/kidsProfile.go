package dbhelper

import (
	"story-time-server/database"
	"story-time-server/models"
)

func AddKidsProfile(kidsDetails models.KidProfile) (int, error) {
	// language=SQL
	SQL := `
			INSERT INTO 
			    kids_profile
			    (
			     name, 
			     age, 
			     avatar_id, 
			     user_id
			     ) 
			VALUES 
			(
			 $1, $2, $3, $4
			)
			RETURNING id
			`

	var id int
	err := database.DB.QueryRowx(SQL, kidsDetails.Name, kidsDetails.Age, kidsDetails.AvatarID, kidsDetails.ParentID).Scan(&id)

	return id, err
}
