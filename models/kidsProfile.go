package models

type KidProfile struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name" validate:"gt=0"`
	Age      int    `json:"age" db:"age" validate:"required"`
	AvatarID int    `json:"avatarId" db:"avatar_id"`
	ParentID int    `json:"parentId"  db:"user_id" validate:"required"`
}
