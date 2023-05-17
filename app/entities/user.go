package entities

import (
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model       `json:"-"`
		Username         string `gorm:"type:varchar(30);not null" json:"username,omitempty"`
		FirstName        string `gorm:"type:varchar(30);not null" json:"fname,omitempty"`
		SureName         string `gorm:"type:varchar(30);not null" json:"sname,omitempty"`
		Email            string `gorm:"type:varchar(255);not null" json:"email,omitempty"`
		Password         string `gorm:"type:varchar(80);not null" json:"password,omitempty"`
		Address          string `gorm:"type:varchar(255);not null" json:"address,omitempty"`
		Image            string `gorm:"type:varchar(255);not null;default:default.jpg" json:"image,omitempty"`
		Role             string `gorm:"not null" json:"-"`
		IsVerified       bool   `gorm:"not null" json:"-"`
		VerificationCode string `gorm:"not null" json:"-"`
		School           School
		Progresses       []Progress
		Submission       []Submission
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
		Id        int
		Email     string `form:"email"`
		Password  string `form:"password"`
		Username  string `form:"username"`
		FirstName string `form:"fname"`
		SureName  string `form:"sname"`
		Address   string `form:"address"`
		Image     string `form:"image"`
	}
)
