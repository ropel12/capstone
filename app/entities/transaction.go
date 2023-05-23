package entities

import (
	"github.com/midtrans/midtrans-go"
	"gorm.io/gorm"
)

type (
	ReqCharge struct {
		PaymentType     string
		Invoice         string
		Total           int
		ItemsDetails    *[]midtrans.ItemDetails
		CustomerDetails *midtrans.CustomerDetails
	}
	Carts struct {
		UserID    uint           `gorm:"not null" `
		SchoolID  uint           `gorm:"not null"`
		Type      string         `gorm:"not null"`
		DeletedAt gorm.DeletedAt `gorm:"index"`
		School    School
		User      User
	}
	Transaction struct {
		Invoice          string `gorm:"primaryKey;not null;type:varchar(20)" json:"invoice,omitempty"`
		UserID           uint   `gorm:"not null"`
		SchoolID         uint   `gorm:"not null"`
		Expire           string `gorm:"not null"`
		Total            int    `gorm:"not null"`
		PaymentCode      string `gorm:"not null"`
		PaymentMethod    string `gorm:"not null"`
		Status           string `gorm:"not null"`
		User             User
		School           School
		TransactionItems []TransactionItems
	}
	TransactionItems struct {
		TransactionInvoice string
		ItemName           string
		ItemPrice          int
	}

	ReqCheckout struct {
		SchoolID      int    `json:"school_id" validate:"required"`
		Type          string `json:"type" validate:"required"`
		PaymentMethod string `json:"payment_method" validate:"required"`
	}
	ResTransaction struct {
		Invoice       string `json:"invoice"`
		PaymentMethod string `json:"payment_method"`
		Total         int    `json:"total"`
		PaymentCode   string `json:"payment_code"`
		ExpireDate    string `json:"expire_date"`
	}
	ResGetAllTrasaction struct {
		SchoolName  string `json:"school_name"`
		SchoolImage string `json:"school_image"`
		SchoolId    int    `json:"school_id"`
	}
	ResDetailRegisCart struct {
		ItemName  string `json:"item_name"`
		ItemPrice int    `json:"item_price"`
		Type      string `json:"type"`
		Total     int    `json:"total"`
	}

	ResDetailPayment struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}

	ResDetailHerRegisCart struct {
		OneTime  []ResDetailPayment `json:"one_time"`
		Interval []ResDetailPayment `json:"interval"`
		Type     string             `json:"type"`
		Total    int                `json:"total"`
	}
	ResDetailTransaction struct {
		Invoice       string
		PaymentMethod string
		Total         int
		PaymentCode   string
		Expire        string
	}
	PaymentSchedule struct {
		
	}
)
