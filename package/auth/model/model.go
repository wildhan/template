package model

type UserAuth struct {
	Id           string `json:"id"       gorm:"column:id"`
	Username     string `json:"username" gorm:"column:username"      validate:"required"`
	Email        string `json:"email"    gorm:"column:email"         validate:"required,email"`
	Password     string `json:"password" gorm:"column:password"      validate:"required,strong_password"`
	HashPassword string `json:"-"        gorm:"column:hash_password"`
}

type LoginParameter struct {
	Email    string `json:"email"    gorm:"column:email"         validate:"required,email"`
	Password string `json:"password" gorm:"column:password"      validate:"required"`
}
