package service

import (
	"context"
	"encoding/json"
	"fmt"

	entity "github.com/education-hub/BE/app/entities"
	school "github.com/education-hub/BE/app/features/school/repository"
	"github.com/education-hub/BE/app/features/transaction/repository"
	user "github.com/education-hub/BE/app/features/user/repository"
	"github.com/education-hub/BE/config/dependency"
	"github.com/education-hub/BE/errorr"
	"github.com/education-hub/BE/helper"
	"github.com/go-playground/validator"
	"github.com/midtrans/midtrans-go"
)

type (
	transaction struct {
		repo       repository.TransactionRepo
		validator  *validator.Validate
		dep        dependency.Depend
		userrepo   user.UserRepo
		schoolrepo school.SchoolRepo
	}
	TransactionService interface {
		CreateTransaction(ctx context.Context, req entity.ReqCheckout, uid int) (*entity.ResTransaction, error)
		GetAllTrasactionCart(ctx context.Context, uid int) ([]entity.ResGetAllTrasaction, error)
		GetDetailTransaction(ctx context.Context, schid, uid int) (any, error)
		UpdateStatus(ctx context.Context, status, invoice string) error
	}
)

func NewTransactionService(repo repository.TransactionRepo, dep dependency.Depend, userrepo user.UserRepo, sch school.SchoolRepo) TransactionService {
	return &transaction{
		repo:       repo,
		dep:        dep,
		validator:  validator.New(),
		userrepo:   userrepo,
		schoolrepo: sch,
	}
}

func (t *transaction) CreateTransaction(ctx context.Context, req entity.ReqCheckout, uid int) (*entity.ResTransaction, error) {
	if err := t.validator.Struct(req); err != nil {
		t.dep.Log.Errorf("[ERROR]WHEN VALIDATE CHECKOUT REQ, err :%v", err)
		return nil, errorr.NewBad("Missing or Invalid Request Body")
	}
	var total int
	var invoice = helper.GenerateInvoice(req.SchoolID, uid)
	itemdetails := []midtrans.ItemDetails{}
	transactionitems := []entity.TransactionItems{}
	if req.Type != "herregistration" && req.Type != "registration" {
		return nil, errorr.NewBad("Invalid Req Body")
	}
	if req.Type == "registration" {
		total = 200000
		itemdetails = append(itemdetails, midtrans.ItemDetails{ID: "1", Name: "First Registration", Price: 200000, Qty: 1})
		transactionitems = append(transactionitems, entity.TransactionItems{ItemName: "First Registration", ItemPrice: 200000, TransactionInvoice: invoice})
	} else {
		data, err := t.repo.GetSchoolPayment(t.dep.Db.WithContext(ctx), req.SchoolID)
		if err != nil {
			return nil, err
		}
		for i, val := range data.Payments {
			itemdetail := midtrans.ItemDetails{
				ID:    fmt.Sprintf("%d", i+1),
				Name:  val.Description,
				Qty:   1,
				Price: int64(val.Price),
			}
			trxitem := entity.TransactionItems{
				TransactionInvoice: invoice,
				ItemName:           val.Description,
				ItemPrice:          val.Price,
			}
			transactionitems = append(transactionitems, trxitem)
			itemdetails = append(itemdetails, itemdetail)
			total += val.Price
		}
	}
	reqcharge := entity.ReqCharge{
		PaymentType:  req.PaymentMethod,
		Invoice:      invoice,
		Total:        total,
		ItemsDetails: &itemdetails,
	}
	res, err := t.dep.Mds.CreateCharge(reqcharge)
	if err != nil {
		return nil, errorr.NewBad(err.Error())
	}
	if res.Expire == "" {
		res.Expire = helper.GenerateExpiretime(res.TransactionTime, t.dep.Mds.ExpDuration)
	}
	trxdata := entity.Transaction{
		Invoice:          invoice,
		UserID:           uint(uid),
		SchoolID:         uint(req.SchoolID),
		Expire:           res.Expire,
		Total:            total,
		PaymentCode:      res.PaymentCode,
		PaymentMethod:    req.PaymentMethod,
		Status:           "pending",
		TransactionItems: transactionitems,
	}

	if err := t.repo.CreateTranscation(t.dep.Db.WithContext(ctx), trxdata, req.Type); err != nil {
		return nil, err
	}
	userdetail, _ := t.userrepo.GetById(t.dep.Db.WithContext(ctx), uid)
	encodeddata, _ := json.Marshal(map[string]any{"invoice": invoice, "total": total, "name": userdetail.FirstName + " " + userdetail.SureName, "email": userdetail.Email, "payment_code": res.PaymentCode, "payment_method": req.PaymentMethod, "expire": res.Expire})
	err = t.dep.Nsq.Publish("1", encodeddata)
	if err != nil {
		t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
	}
	return &entity.ResTransaction{Invoice: invoice, PaymentMethod: req.PaymentMethod, Total: total, PaymentCode: res.PaymentCode, ExpireDate: res.Expire}, nil
}

func (t *transaction) GetAllTrasactionCart(ctx context.Context, uid int) ([]entity.ResGetAllTrasaction, error) {
	data, err := t.repo.GetAllCartByuid(t.dep.Db.WithContext(ctx), uid)
	if err != nil {
		return nil, err
	}
	res := []entity.ResGetAllTrasaction{}
	for _, val := range data {
		trx := entity.ResGetAllTrasaction{
			SchoolName:  val.School.Name,
			SchoolImage: val.School.Image,
			SchoolId:    int(val.School.ID),
		}
		res = append(res, trx)
	}
	return res, nil
}

func (t *transaction) GetDetailTransaction(ctx context.Context, schid, uid int) (any, error) {

	trxdata, err := t.repo.GetTransaction(t.dep.Db.WithContext(ctx), schid, uid)
	if err == nil {
		return entity.ResDetailTransaction{
			Invoice:       trxdata.Invoice,
			PaymentMethod: trxdata.PaymentCode,
			Total:         trxdata.Total,
			PaymentCode:   trxdata.PaymentCode,
			Expire:        trxdata.Expire,
		}, nil
	}
	trxcart, err := t.repo.GetCart(t.dep.Db.WithContext(ctx), schid, uid)
	if err != nil {
		return nil, err
	}
	if trxcart.Type == "registration" {
		return entity.ResDetailRegisCart{ItemName: "First Registraion", ItemPrice: 200000, Type: "registration", Total: 200000}, nil
	}
	total := 0
	restrx := entity.ResDetailHerRegisCart{Type: "herregistration"}
	for _, val := range trxcart.School.Payments {
		detailpayment := entity.ResDetailPayment{Name: val.Description, Price: val.Price}
		if val.Type == "one" {
			restrx.OneTime = append(restrx.OneTime, detailpayment)
		} else {
			restrx.Interval = append(restrx.Interval, detailpayment)
		}
		total += val.Price
	}
	restrx.Total = total
	return restrx, nil
}

func (t *transaction) UpdateStatus(ctx context.Context, status, invoice string) error {
	trxdata, _ := t.repo.GetTransactionByInvoice(t.dep.Db.WithContext(ctx), invoice)
	encodeddata, _ := json.Marshal(map[string]any{"invoice": invoice, "email": trxdata.User.Email, "name": trxdata.User.FirstName + " " + trxdata.User.SureName})
	switch status {
	case "paid":
		err := t.repo.UpdateStatus(t.dep.Db.WithContext(ctx), invoice, "paid")
		if err != nil {
			t.dep.Log.Errorf("[ERROR]WHEN UPDATING TRASACTION STATUS,Err : %v", err)
		}
		cartdata, err := t.repo.GetCart(t.dep.Db.WithContext(ctx), int(trxdata.SchoolID), int(trxdata.UserID))
		if err != nil {
			t.dep.Log.Errorf("[ERROR]WHEN GETTING CART DATA,Err : %v", err)
		}
		if cartdata.Type == "registration" {
			err = t.schoolrepo.UpdateProgressByUid(t.dep.Db.WithContext(ctx), int(trxdata.UserID), int(trxdata.SchoolID), "Done Payment")
			if err != nil {
				t.dep.Log.Errorf("[ERROR]WHEN UPDATING PROGRESS STATUS,Err : %v", err)
			}
		} else {
			err = t.schoolrepo.UpdateProgressByUid(t.dep.Db.WithContext(ctx), int(trxdata.UserID), int(trxdata.SchoolID), "Already Paid Her-Registration")
			if err != nil {
				t.dep.Log.Errorf("[ERROR]WHEN UPDATING PROGRESS STATUS,Err : %v", err)
			}
		}
		go func() {
			encodeddata, _ := json.Marshal(map[string]any{"invoice": invoice, "email": trxdata.User.Email, "name": trxdata.User.FirstName + " " + trxdata.User.SureName})
			if err := t.dep.Nsq.Publish("2", encodeddata); err != nil {
				t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
			}
		}()
		err = t.repo.DeleteCart(t.dep.Db.WithContext(ctx), int(trxdata.SchoolID), int(trxdata.UserID))
		if err != nil {
			t.dep.Log.Errorf("[ERROR]WHEN DELETING CART,Err : %v", err)
		}
		if err := t.dep.Pusher.Publish(map[string]string{"username": trxdata.User.Username, "type": "payment", "school_name": trxdata.School.Name, "status": "success"}, 1); err != nil {
			t.dep.Log.Errorf("Failed to publish to PusherJs: %v", err)
		}
	case "cancel":
		err := t.repo.DeleteCart(t.dep.Db.WithContext(ctx), int(trxdata.SchoolID), int(trxdata.UserID))
		if err != nil {
			t.dep.Log.Errorf("[ERROR]WHEN DELETE CART,Err : %v", err)
		}
		err = t.repo.UpdateStatus(t.dep.Db.WithContext(ctx), invoice, "cancel")
		if err != nil {
			t.dep.Log.Errorf("[ERROR]WHEN UPDATING TRASACTION STATUS,Err : %v", err)
		}
		err = t.schoolrepo.UpdateProgressByUid(t.dep.Db.WithContext(ctx), int(trxdata.UserID), int(trxdata.UserID), "failed")
		if err != nil {
			t.dep.Log.Errorf("[ERROR]WHEN UPDATING PROGRESS STATUS,Err : %v", err)
		}
		go func() {
			if err := t.dep.Nsq.Publish("3", encodeddata); err != nil {
				t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
			}
		}()
		if err := t.dep.Pusher.Publish(map[string]string{"username": trxdata.User.Username, "type": "payment", "school_name": trxdata.School.Name, "status": "cancel"}, 1); err != nil {
			t.dep.Log.Errorf("Failed to publish to PusherJs: %v", err)
		}
	}
	return nil
}
