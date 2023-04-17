package models

type Avatars struct {
	ID        int    `json:"id" db:"id"`
	AvatarURL string `json:"avatarUrl" db:"avatar_url"`
}
