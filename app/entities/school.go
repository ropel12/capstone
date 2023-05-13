package entities

import "gorm.io/gorm"

type (
	School struct {
		gorm.Model
		UserID          uint
		Npsn            string `gorm:"type:varchar(12);not null"`
		Name            string `gorm:"type:varchar(150);not null"`
		Description     string `gorm:"type:varchar(255);not null"`
		Image           string `gorm:"type:varchar(150);not null"`
		Video           string `gorm:"type:varchar(150);not null"`
		Pdf             string `gorm:"type:varchar(150);not null"`
		Web             string `gorm:"type:varchar(150);not null"`
		Province        string `gorm:"type:varchar(150);not null"`
		City            string `gorm:"type:varchar(150);not null"`
		District        string `gorm:"type:varchar(150);not null"`
		Village         string `gorm:"type:varchar(150);not null"`
		Detail          string `gorm:"type:varchar(150);not null"`
		ZipCode         string `gorm:"type:varchar(150);not null"`
		Students        int    `gorm:"not null"`
		Teachers        int    `gorm:"not null"`
		Staff           int    `gorm:"not null"`
		Accreditation   string `gorm:"type:varchar(3);not null"`
		Gmeet           string `gorm:"type:varchar(35);default: "`
		QuizLinkPub     string `gorm:"type:varchar(70);default:"`
		QuizLinkPreview string `gorm:"type:varchar(70);default:"`
	}
	ReqCreateSchool struct {
		UserId        int
		Npsn          string `form:"npsn" validate:"required"`
		Name          string `form:"name" validate:"required"`
		Description   string `form:"description" validate:"required"`
		Image         string `form:"image" validate:"required"`
		Video         string `form:"video" validate:"required"`
		Pdf           string `form:"pdf" validate:"required"`
		Web           string `form:"web" validate:"required"`
		Province      string `form:"province" validate:"required"`
		City          string `form:"city" validate:"required"`
		District      string `form:"district" validate:"required"`
		Village       string `form:"village" validate:"required"`
		Detail        string `form:"detail" validate:"required"`
		ZipCode       string `form:"zipcode" validate:"required"`
		Students      int    `form:"students" validate:"required"`
		Teachers      int    `form:"teachers" validate:"required"`
		Staff         int    `form:"staff" validate:"required"`
		Accreditation string `form:"accreditation" validate:"required"`
	}
)
