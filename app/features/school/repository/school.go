package repository

import (
	entity "github.com/education-hub/BE/app/entities"
	"github.com/education-hub/BE/errorr"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	school struct {
		log *logrus.Logger
	}
	SchoolRepo interface {
		Create(db *gorm.DB, school entity.School) (int, error)
		FindByNPSN(db *gorm.DB, npsn string) error
	}
)

func NewSchoolRepo(log *logrus.Logger) SchoolRepo {
	return &school{log}
}

func (u *school) Create(db *gorm.DB, school entity.School) (int, error) {
	if err := db.Create(&school).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN CREATE USER,Error: %v ", err)
		return 0, errorr.NewInternal("Internal Server Error")
	}
	return int(school.ID), nil
}

func (u *school) FindByNPSN(db *gorm.DB, npsn string) error {
	data := entity.School{}
	if err := db.Find(&data).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN GETTING SCHOOL DATA,Error: %v ", err)
		return errorr.NewInternal("Internal Server Error")
	}
	if data.Name == "" {
		return errorr.NewBad("Data Not Found")
	}
	return nil
}
