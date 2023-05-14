package repository

import (
	"reflect"

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
		Update(db *gorm.DB, school entity.School) (*entity.School, error)
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
	if err := db.Where("npsn=?", npsn).Find(&data).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN GETTING SCHOOL DATA,Error: %v ", err)
		return errorr.NewInternal("Internal Server Error")
	}
	if data.Name == "" {
		return errorr.NewBad("Data Not Found")
	}
	return nil
}

func (s *school) Update(db *gorm.DB, school entity.School) (*entity.School, error) {
	newdata := entity.School{}
	if err := db.First(&newdata, school.ID).Error; err == gorm.ErrRecordNotFound {
		s.log.Errorf("ERROR]WHEN UPDATE School,Error: %v ", err)
		return nil, errorr.NewBad("Id Not Found")
	}
	v := reflect.ValueOf(school)
	n := reflect.ValueOf(&newdata).Elem()
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Interface().(type) {
		case string:
			val := v.Field(i).Interface().(string)
			if val != "" {
				n.Field(i).SetString(val)
			}
		case int:
			val := v.Field(i).Interface().(int)
			if val != 0 {
				n.Field(i).SetInt(int64(val))
			}
		}
	}

	if err := db.Save(&newdata).Error; err != nil {
		s.log.Errorf("[ERROR]WHEN UPDATING SCHOOL, Err: %v", err)
		return nil, errorr.NewInternal("Internal server error")
	}
	return &newdata, nil
}
