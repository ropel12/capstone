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
		AddAchievement(db *gorm.DB, achv entity.Achievement) (int, error)
		DeleteAchievement(db *gorm.DB, id int) error
		UpdateAchievement(db *gorm.DB, achv entity.Achievement) (*entity.Achievement, error)
		AddExtracurricular(db *gorm.DB, achv entity.Extracurricular) (int, error)
		DeleteExtracurricular(db *gorm.DB, id int) error
		UpdateExtracurricular(db *gorm.DB, achv entity.Extracurricular) (*entity.Extracurricular, error)
		GetByUid(db *gorm.DB, uid int) (*entity.School, error)
		GetById(db *gorm.DB, id int) (*entity.School, error)
		AddFaq(db *gorm.DB, faq entity.Faq) (int, error)
		DeleteFaq(db *gorm.DB, id int) error
		UpdateFaq(db *gorm.DB, extrac entity.Faq) (*entity.Faq, error)
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
func (u *school) GetByUid(db *gorm.DB, uid int) (*entity.School, error) {
	res := entity.School{}
	if err := db.Preload("Achievements", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,school_id,description,image,title")
	}).Preload("Extracurriculars", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,description,image,title")
	}).Where("user_id=?", uid).Find(&res).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN GETTING The School Data BY UID, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	if res.Name == "" {
		return nil, errorr.NewBad("Data Not Found")
	}
	return &res, nil
}
func (u *school) GetById(db *gorm.DB, id int) (*entity.School, error) {
	res := entity.School{}
	if err := db.Preload("Achievements", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,description,image,title")
	}).Preload("Extracurriculars", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,description,image,title")
	}).Where("id=?", id).Find(&res).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN GETTING The School Data BY SchoolID, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	if res.Name == "" {
		return nil, errorr.NewBad("Data Not Found")
	}
	return &res, nil
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

func (u *school) AddAchievement(db *gorm.DB, achv entity.Achievement) (int, error) {
	if err := db.Save(&achv).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN ADDING ACHIEVEMENT, Err: %v", err)
		return 0, errorr.NewInternal("internal Server Error")
	}
	return int(achv.SchoolID), nil
}

func (u *school) DeleteAchievement(db *gorm.DB, id int) error {

	if err := db.Where("id=?", id).First(&entity.Achievement{}).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			u.log.Errorf("[ERROR]WHEN GETTING The Achievement Data, Err: %v", err)
			return errorr.NewInternal("Internal Server Error")
		}
		return errorr.NewBad("Id Not Found")
	}
	if err := db.Where("id=?", id).Delete(&entity.Achievement{}).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN DELETING Achievement, Err: %v", err)
		return errorr.NewInternal("Internal Server Error")
	}
	return nil
}

func (s *school) UpdateAchievement(db *gorm.DB, achv entity.Achievement) (*entity.Achievement, error) {
	newdata := entity.Achievement{}
	if err := db.First(&newdata, achv.ID).Error; err == gorm.ErrRecordNotFound {
		s.log.Errorf("ERROR]WHEN UPDATE ACHIEVEMENT,Error: %v ", err)
		return nil, errorr.NewBad("Id Not Found")
	}
	v := reflect.ValueOf(achv)
	n := reflect.ValueOf(&newdata).Elem()

	for i := 0; i < v.NumField(); i++ {
		if val, ok := v.Field(i).Interface().(string); ok {
			if val != "" {
				n.Field(i).SetString(val)
			}
		}
	}
	if err := db.Save(&newdata).Error; err != nil {
		s.log.Errorf("[ERROR]WHEN UPDATING ACHIEVEMENT, Err: %v", err)
		return nil, errorr.NewInternal("Internal server error")
	}
	return &newdata, nil
}

func (u *school) AddExtracurricular(db *gorm.DB, extrac entity.Extracurricular) (int, error) {
	if err := db.Save(&extrac).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN ADDING Extracurricular, Err: %v", err)
		return 0, errorr.NewInternal("internal Server Error")
	}
	return int(extrac.SchoolID), nil
}

func (u *school) DeleteExtracurricular(db *gorm.DB, id int) error {

	if err := db.Where("id=?", id).First(&entity.Extracurricular{}).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			u.log.Errorf("[ERROR]WHEN GETTING The Extracurricular Data, Err: %v", err)
			return errorr.NewInternal("Internal Server Error")
		}
		return errorr.NewBad("Id Not Found")
	}
	if err := db.Where("id=?", id).Delete(&entity.Extracurricular{}).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN DELETING Extracurricular, Err: %v", err)
		return errorr.NewInternal("Internal Server Error")
	}
	return nil
}

func (s *school) UpdateExtracurricular(db *gorm.DB, extrac entity.Extracurricular) (*entity.Extracurricular, error) {
	newdata := entity.Extracurricular{}
	if err := db.First(&newdata, extrac.ID).Error; err == gorm.ErrRecordNotFound {
		s.log.Errorf("ERROR]WHEN UPDATE Extracurricular,Error: %v ", err)
		return nil, errorr.NewBad("Id Not Found")
	}
	v := reflect.ValueOf(extrac)
	n := reflect.ValueOf(&newdata).Elem()

	for i := 0; i < v.NumField(); i++ {
		if val, ok := v.Field(i).Interface().(string); ok {
			if val != "" {
				n.Field(i).SetString(val)
			}
		}
	}
	if err := db.Save(&newdata).Error; err != nil {
		s.log.Errorf("[ERROR]WHEN UPDATING Extracurricular, Err: %v", err)
		return nil, errorr.NewInternal("Internal server error")
	}
	return &newdata, nil
}

func (u *school) AddFaq(db *gorm.DB, faq entity.Faq) (int, error) {
	if err := db.Save(&faq).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN ADDING Faq, Err: %v", err)
		return 0, errorr.NewInternal("internal Server Error")
	}
	return int(faq.SchoolID), nil
}

func (u *school) DeleteFaq(db *gorm.DB, id int) error {

	if err := db.Where("id=?", id).First(&entity.Faq{}).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			u.log.Errorf("[ERROR]WHEN GETTING The Faq Data, Err: %v", err)
			return errorr.NewInternal("Internal Server Error")
		}
		return errorr.NewBad("Id Not Found")
	}
	if err := db.Where("id=?", id).Delete(&entity.Faq{}).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN DELETING Faq, Err: %v", err)
		return errorr.NewInternal("Internal Server Error")
	}
	return nil
}
func (s *school) UpdateFaq(db *gorm.DB, extrac entity.Faq) (*entity.Faq, error) {
	newdata := entity.Faq{}
	if err := db.First(&newdata, extrac.ID).Error; err == gorm.ErrRecordNotFound {
		s.log.Errorf("ERROR]WHEN UPDATE Faq,Error: %v ", err)
		return nil, errorr.NewBad("Id Not Found")
	}
	v := reflect.ValueOf(extrac)
	n := reflect.ValueOf(&newdata).Elem()

	for i := 0; i < v.NumField(); i++ {
		if val, ok := v.Field(i).Interface().(string); ok {
			if val != "" {
				n.Field(i).SetString(val)
			}
		}
	}
	if err := db.Save(&newdata).Error; err != nil {
		s.log.Errorf("[ERROR]WHEN UPDATING Faq, Err: %v", err)
		return nil, errorr.NewInternal("Internal server error")
	}
	return &newdata, nil
}
