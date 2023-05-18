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
		Delete(db *gorm.DB, id int, uid int) error
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
		AddPayment(db *gorm.DB, paym entity.Payment) (int, error)
		DeletePayment(db *gorm.DB, id int) error
		GetAll(db *gorm.DB, limit, offset int, search string) ([]entity.School, int, error)
		UpdatePayment(db *gorm.DB, paym entity.Payment) (*entity.Payment, error)
		CreateSubmission(db *gorm.DB, subm entity.Submission) (int, error)
		UpdateProgress(db *gorm.DB, id int, status string) (*entity.Progress, error)
		GetAllProgressByuid(db *gorm.DB, uid int) ([]entity.Progress, error)
		GetProgressByid(db *gorm.DB, id int) (*entity.Progress, error)
		GetAllProgressAndSubmissionByuid(db *gorm.DB, uid int) (*entity.School, error)
		GetSubmissionByid(db *gorm.DB, id int) (*entity.Submission, error)
		AddReview(db *gorm.DB, data entity.Reviews) (int, error)
		UpdateProgressByUid(db *gorm.DB, uid int, schid int, status string) error
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
		return db.Select("id,school_id,description,image,title")
	}).Preload("Faqs", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,school_id,question,answer")
	}).Preload("Payments").Where("user_id=?", uid).Find(&res).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN GETTING The School Data BY UID, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	if res.Name == "" {
		return nil, errorr.NewBad("Data Not Found")
	}
	return &res, nil
}
func (u *school) GetAll(db *gorm.DB, limit, offset int, search string) ([]entity.School, int, error) {
	res := []entity.School{}
	search = "%" + search + "%"
	var total int64
	db.Model(&entity.School{}).Where("deleted_at IS NULL AND (name like ? or province like ? or district like ? or village like ? or detail like ?)", search, search, search, search, search).Count(&total)
	if err := db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,username")
	}).Where("deleted_at IS NULL AND (name like ? or province like ? or district like ? or village like ? or detail like ?)", search, search, search, search, search).Find(&res).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN GETTING SCHOOL DATA, Err : %v", err)
		return nil, 0, errorr.NewInternal("Internal Server Error")
	}
	return res, int(total), nil
}

func (u *school) Delete(db *gorm.DB, id int, uid int) error {

	if err := db.Where("id=? AND user_id=?", id, uid).First(&entity.School{}).Error; err != nil {
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
func (u *school) GetById(db *gorm.DB, id int) (*entity.School, error) {
	res := entity.School{}
	if err := db.Preload("Achievements", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, school_id, description, image, title")
	}).Preload("Extracurriculars", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, school_id, description, image, title")
	}).Preload("Faqs", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, school_id, question, answer")
	}).Preload("Reviews", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User").Select("user_id,school_id,review")
	}).Preload("Payments").Where("id=?", id).Find(&res).Error; err != nil {
		u.log.Errorf("[ERROR] WHEN GETTING The School Data BY SchoolID, Err: %v", err)
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

func (u *school) AddPayment(db *gorm.DB, paym entity.Payment) (int, error) {
	if err := db.Save(&paym).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN ADDING Payment, Err: %v", err)
		return 0, errorr.NewInternal("internal Server Error")
	}
	return int(paym.SchoolID), nil
}

func (u *school) DeletePayment(db *gorm.DB, id int) error {

	if err := db.Where("id=?", id).First(&entity.Payment{}).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			u.log.Errorf("[ERROR]WHEN GETTING The Payment Data, Err: %v", err)
			return errorr.NewInternal("Internal Server Error")
		}
		return errorr.NewBad("Id Not Found")
	}
	if err := db.Where("id=?", id).Delete(&entity.Payment{}).Error; err != nil {
		u.log.Errorf("[ERROR]WHEN DELETING Payment, Err: %v", err)
		return errorr.NewInternal("Internal Server Error")
	}
	return nil
}

func (s *school) UpdatePayment(db *gorm.DB, paym entity.Payment) (*entity.Payment, error) {
	newdata := entity.Payment{}
	if err := db.First(&newdata, paym.ID).Error; err == gorm.ErrRecordNotFound {
		s.log.Errorf("ERROR]WHEN UPDATE Payment,Error: %v ", err)
		return nil, errorr.NewBad("Id Not Found")
	}
	v := reflect.ValueOf(paym)
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
			if val != -1 {
				n.Field(i).SetInt(int64(val))
			}
		}
	}
	if err := db.Save(&newdata).Error; err != nil {
		s.log.Errorf("[ERROR]WHEN UPDATING Payment, Err: %v", err)
		return nil, errorr.NewInternal("Internal server error")
	}
	return &newdata, nil
}

func (s *school) CreateSubmission(db *gorm.DB, subm entity.Submission) (int, error) {
	progress := entity.Progress{UserID: subm.UserID, SchoolID: subm.SchoolID, Status: "Checking File"}
	err := db.Transaction(func(db *gorm.DB) error {

		if err := db.Create(&subm).Error; err != nil {
			s.log.Errorf("[ERROR]WHEN CREATING Submission, Err: %v", err)
			return errorr.NewInternal("Internal server error")
		}
		if err := db.Create(&progress).Error; err != nil {
			return errorr.NewInternal("Internal server error")
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return int(progress.ID), nil
}

func (s *school) UpdateProgress(db *gorm.DB, id int, status string) (*entity.Progress, error) {
	prog := entity.Progress{}
	if err := db.Where("status != 'Finish' AND status != 'Failed'").First(&prog, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorr.NewBad("Data not found")
		}
		s.log.Errorf("[ERORR]WHEN GETTING Progress DATA, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	if err := db.Model(&entity.Progress{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		s.log.Errorf("[ERROR]WHEN UPDATING Submission, Err : %v", err)
		return nil, errorr.NewInternal("Internal Server Erorr")
	}
	if status == "Send Detail Costs Registration" {
		db.Create(&entity.Carts{UserID: prog.UserID, SchoolID: prog.SchoolID, Type: "registration"})
	} else if status == "Send Detail Costs Her-Registration" {
		db.Create(&entity.Carts{UserID: prog.UserID, SchoolID: prog.SchoolID, Type: "herregistration"})
	}
	return &prog, nil
}
func (s *school) UpdateProgressByUid(db *gorm.DB, uid int, schid int, status string) error {
	if err := db.Where("status != 'Finish' AND status != 'Failed' AND user_id=? AND school_id=?", uid, schid).First(&entity.Progress{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errorr.NewBad("Data not found")
		}
		s.log.Errorf("[ERORR]WHEN GETTING Progress DATA, Err: %v", err)
		return errorr.NewInternal("Internal Server Error")
	}
	if err := db.Model(&entity.Progress{}).Where("user_id = ? AND school_id=?", uid, schid).Update("status", status).Error; err != nil {
		s.log.Errorf("[ERROR]WHEN UPDATING PROGRESS, Err : %v", err)
		return errorr.NewInternal("Internal Server Erorr")
	}
	return nil
}
func (s *school) GetAllProgressByuid(db *gorm.DB, uid int) ([]entity.Progress, error) {
	res := []entity.Progress{}
	if err := db.Preload("School", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,image,name,web")
	}).Where("user_id=? AND status != 'Failed' AND status !='Finish'", uid).Find(&res).Error; err != nil {
		s.log.Errorf("[ERROR]WHEN GETTING Student Progress Data, Err : %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	if len(res) == 0 {
		return nil, errorr.NewBad("Data Not Found")
	}
	return res, nil
}

func (s *school) GetProgressByid(db *gorm.DB, id int) (*entity.Progress, error) {
	res := entity.Progress{}
	if err := db.First(&res).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorr.NewBad("Data Not Found")
		}
		s.log.Errorf("[ERROR]WHEN GETTING Student Progress Data, Err : %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	return &res, nil
}

func (s *school) GetAllProgressAndSubmissionByuid(db *gorm.DB, uid int) (*entity.School, error) {
	res := entity.School{}
	if err := db.Preload("Progresses", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,first_name,sure_name,image")
		}).Select("school_id,id,user_id")
	}).Preload("Submissions", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,school_id")
	}).Joins("join progresses p on p.school_id=schools.id").Where("schools.user_id=? AND p.status != 'Failed'", uid).Find(&res).Error; err != nil {
		s.log.Errorf("[ERROR]WHEN GETTING PRORGRESS AND SUBMISSION DATA, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Erorr")
	}
	if res.Name == "" {
		return nil, errorr.NewBad("Data Not Found")
	}
	return &res, nil
}

func (s *school) GetSubmissionByid(db *gorm.DB, id int) (*entity.Submission, error) {
	res := entity.Submission{}
	if err := db.First(&res, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorr.NewBad("Data Not Found")
		}
		s.log.Errorf("[ERROR]WHEN GETTING SUBMISSION DATA, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	return &res, nil
}

func (s *school) AddReview(db *gorm.DB, data entity.Reviews) (int, error) {
	if err := db.Create(&data).Error; err != nil {
		s.log.Errorf("[ERROR]WHEN CREATING review, err :%v", err)
		return 0, errorr.NewInternal("Internal Server Error")
	}
	return int(data.SchoolID), nil
}
