package dbhelper

import (
	"story-time-server/database"
	"story-time-server/models"
)

func GetAvatars() ([]models.Avatars, error) {
	// language=SQL
	SQL := `
			SELECT 
			    id,
			    avatar_url
			FROM
			    avatars
			    `

	avatars := make([]models.Avatars, 0)
	err := database.DB.Select(&avatars, SQL)

	return avatars, err
}
