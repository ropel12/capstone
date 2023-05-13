package user

import "gorm.io/gorm"

type (
	User struct {
		gorm.Model
		Username         string `gorm:"type:varchar(30);not null"`
		FirstName        string `gorm:"type:varchar(30);not null"`
		SureName         string `gorm:"type:varchar(30);not null"`
		Email            string `gorm:"type:varchar(255);not null"`
		Password         string `gorm:"type:varchar(80);not null"`
		Address          string `gorm:"type:varchar(255);not null"`
		Image            string `gorm:"type:varchar(255);not null;default:default.jpg"`
		Role             string `gorm:"not null"`
		IsVerified       bool   `gorm:"not null"`
		VerificationCode string `gorm:"not null"`
	}

	ForgotPass struct {
		Token     string
		Email     string
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}
	LoginReq struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	RegisterReq struct {
		Email     string `json:"email" validate:"required"`
		FirstName string `json:"firstname" validate:"required"`
		LastName  string `json:"lastname" validate:"required"`
		Username  string `json:"username" validate:"required"`
		Password  string `json:"password" validate:"required"`
		Address   string `json:"address" validate:"required"`
		Role      string `json:"role" validate:"required"`
	}
	UpdateReq struct {
		Id       int
		Email    string `form:"email"`
		Password string `form:"password"`
		Name     string `form:"name" `
		Address  string `form:"address"`
		Image    string `form:"image"`
	}
)
