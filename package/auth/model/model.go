package model

type UserAuth struct {
	Id           string `json:"id"       gorm:"column:id"`
	Username     string `json:"username" gorm:"column:username"      validate:"required"`
	Email        string `json:"email"    gorm:"column:email"         validate:"required,email"`
	Password     string `json:"password" gorm:"column:password"      validate:"required,strong_password"`
	HashPassword string `json:"-"        gorm:"column:hash_password"`
}

type UserProfile struct {
	Id            string  `json:"id"              gorm:"column:id"`
	FirstName     *string `json:"first_name"      gorm:"column:first_name"`
	LastName      *string `json:"last_name"       gorm:"column:last_name"`
	Brithday      *string `json:"birthday"        gorm:"column:birthday"`
	Address       *string `json:"address"         gorm:"column:address"`
	Phone         *string `json:"phone"           gorm:"column:phone"`
	CreatedAt     *string `json:"created_at"      gorm:"column:created_at"`
	CreatedByUser *string `json:"created_by_user" gorm:"column:created_by_user"`
	CreatedByRole *string `json:"created_by_role" gorm:"column:created_by_role"`
	UpdatedAt     *string `json:"updated_at"      gorm:"column:updated_at"`
	UpdatedByUser *string `json:"updated_by_user" gorm:"column:updated_by_user"`
	UpdatedByRole *string `json:"updated_by_role" gorm:"column:updated_by_role"`
	DeletedAt     *string `json:"deleted_at"      gorm:"column:deleted_at"`
	DeletedByUser *string `json:"deleted_by_user" gorm:"column:deleted_by_user"`
	DeletedByRole *string `json:"deleted_by_role" gorm:"column:deleted_by_role"`
}

type LoginParameter struct {
	Email    string `json:"email"    gorm:"column:email"         validate:"required,email"`
	Password string `json:"password" gorm:"column:password"      validate:"required"`
}

type LoginResponse struct {
	Username     string `json:"username"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
