package models

type User struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" validate:"gt=0"`
	Email       string `json:"email" db:"email" validate:"email,required"`
	Password    string `json:"password" db:"password" validate:"gt=4"`
	YearOfBirth int    `json:"yearOfBirth" db:"year_of_birth" validate:"required"`
	FireBaseUID string `json:"fireBaseId" db:"firebase_uid"`
	AvatarID    int    `json:"avatarId" db:"avatar_id"`
}

type UserLoginReq struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"  validate:"gt=4"`
}
