package repository

import (
	entity "github.com/education-hub/BE/app/entities"
	"github.com/education-hub/BE/errorr"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	transaction struct {
		log *logrus.Logger
	}

	TransactionRepo interface {
		CreateTranscation(db *gorm.DB, data entity.Transaction, typee string) error
		GetTransaction(db *gorm.DB, schoolid int, userid int) (*entity.Transaction, error)
		GetAllCartByuid(db *gorm.DB, uid int) ([]entity.Carts, error)
		GetCart(db *gorm.DB, schid int, userid int) (*entity.Carts, error)
		DeleteCart(db *gorm.DB, schid int, uid int) error
		UpdateStatus(db *gorm.DB, invoice string, status string) error
		GetSchoolPayment(db *gorm.DB, schid int) (*entity.School, error)
		GetTransactionByInvoice(db *gorm.DB, invoice string) (*entity.Transaction, error)
	}
)

func NewTransactionRepo(log *logrus.Logger) TransactionRepo {
	return &transaction{log: log}
}

func (t *transaction) CreateTranscation(db *gorm.DB, data entity.Transaction, typee string) error {
	return db.Transaction(func(db *gorm.DB) error {
		if err := db.Create(&data).Error; err != nil {
			t.log.Errorf("[ERROR]WHEN GETTING Transaction Data, Err: %v", err)
			return errorr.NewInternal("Internal Server Error")
		}
		status := ""
		if typee == "registration" {
			status = "Sending Detail Cost Registration"
		} else {
			status = "Sending Detail Cost Her-Registration"
		}
		if err := db.Model(&entity.Progress{}).Where("user_id=? AND school_id=? AND status != 'Finish' AND status != 'failed'", data.UserID, data.SchoolID).Update("status", status).Error; err != nil {
			return errorr.NewInternal("Internal Server Error")
		}
		return nil
	})
}

func (t *transaction) GetTransaction(db *gorm.DB, schoolid int, userid int) (*entity.Transaction, error) {
	res := entity.Transaction{}
	if err := db.Preload("TransactionItems").Where("school_id = ?  AND user_id = ? AND (status != 'paid' OR status != 'cancel')", schoolid, userid).Find(&res).Error; err != nil {
		t.log.Errorf("[ERROR]WHEN GETTING TRANSACTION DATA, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	if res.Invoice == "" {
		return nil, errorr.NewBad("Data Not Found")
	}
	return &res, nil
}

func (t *transaction) GetTransactionByInvoice(db *gorm.DB, invoice string) (*entity.Transaction, error) {
	res := entity.Transaction{}
	if err := db.Preload("TransactionItems").Preload("User").Where("invoice =?", invoice).Find(&res).Error; err != nil {
		t.log.Errorf("[ERROR]WHEN GETTING TRANSACTION DATA, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	if res.Invoice == "" {
		return nil, errorr.NewBad("Data Not Found")
	}
	return &res, nil
}

func (t *transaction) GetAllCartByuid(db *gorm.DB, uid int) ([]entity.Carts, error) {
	res := []entity.Carts{}
	if err := db.Preload("School").Where("user_id=? AND deleted_at IS NULL", uid).Find(&res).Error; err != nil {
		t.log.Errorf("[ERROR]WHEN GETTING TRANSACTION DATA, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	if len(res) == 0 {
		return nil, errorr.NewBad("Data Not Found")
	}
	return res, nil
}
func (t *transaction) GetCart(db *gorm.DB, schid int, userid int) (*entity.Carts, error) {
	res := entity.Carts{}
	if err := db.Preload("School", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Payments")
	}).Where("school_id=? AND user_id=?", schid, userid).Find(&res).Error; err != nil {
		t.log.Errorf("[ERROR]WHEN GETTING CART DATA, ERR: %v", err)
		return nil, errorr.NewBad("Internal Server Error")
	}
	if int(res.SchoolID) == 0 {
		return nil, errorr.NewBad("Data Not Found")
	}
	return &res, nil
}
func (t *transaction) DeleteCart(db *gorm.DB, schid int, uid int) error {
	if err := db.Where("user_id=? AND school_id=? AND deleted_at is null", uid, schid).Delete(&entity.Carts{}).Error; err != nil {
		t.log.Errorf("[ERROR]WHEN DELETEING CART, Err :%v", err)
		return errorr.NewInternal("Internal Server Error")
	}
	return nil
}

func (t *transaction) UpdateStatus(db *gorm.DB, invoice string, status string) error {
	data := entity.Transaction{}
	if err := db.Where("invoice=?", invoice).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errorr.NewBad("Data Not Found")
		}
		t.log.Errorf("[ERROR]WHEN GETTING TRANSACTION Data, Err : %v", err)
		return errorr.NewInternal("Internal Server Error")
	}

	data.Status = status
	if err := db.Save(&data).Error; err != nil {
		t.log.Errorf("[ERORR]WHEN UPDATING TRANSACTION STATUS, Err : %v", err)
		return errorr.NewInternal("Internal Server Erorr")
	}
	return nil
}
func (t *transaction) GetSchoolPayment(db *gorm.DB, schid int) (*entity.School, error) {
	res := entity.School{}
	if err := db.Preload("Payments").Where("id=?", schid).Find(&res).Error; err != nil {
		t.log.Errorf("[ERROR]WHEN GETTING SCHOOL PAYMENT, Error: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	if res.Name == "" {
		return nil, errorr.NewBad("Data Not Found")
	}
	return &res, nil
}
