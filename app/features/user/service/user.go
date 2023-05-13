package service

import (
	"context"
	"encoding/base32"

	entity "github.com/education-hub/BE/app/entities/user"
	"github.com/education-hub/BE/app/features/user/repository"
	dependcy "github.com/education-hub/BE/config/dependency"
	"github.com/education-hub/BE/errorr"
	"github.com/education-hub/BE/helper"
	"github.com/go-playground/validator"
)

type (
	user struct {
		repo      repository.UserRepo
		validator *validator.Validate
		dep       dependcy.Depend
	}
	UserService interface {
		Login(ctx context.Context, req entity.LoginReq) (int, string, error)
		Register(ctx context.Context, req entity.RegisterReq) error
		VerifyEmail(ctx context.Context, verificationcode string) error
		ForgetPass(ctx context.Context, email string) error
		ResetPass(ctx context.Context, token string, newpass string) error
	}
)

func NewUserService(repo repository.UserRepo, dep dependcy.Depend) UserService {
	return &user{repo: repo, dep: dep, validator: validator.New()}
}

func (u *user) Login(ctx context.Context, req entity.LoginReq) (int, string, error) {
	if err := u.validator.Struct(req); err != nil {
		u.dep.Log.Errorf("[ERROR] WHEN VALIDATE LOGIN REQ, Error: %v", err)
		return 0, "", errorr.NewBad("Missing or Invalid Request Body")
	}
	user, err := u.repo.FindByUsername(u.dep.Db.WithContext(ctx), req.Username)
	if err != nil {
		return 0, "", err
	}
	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		u.dep.Log.Errorf("Error Service : %v", err)
		return 0, "", errorr.NewBad("Wrong password")
	}
	if user.IsVerified == false {
		return 0, "", errorr.NewBad("Email Not Verified")
	}
	return int(user.ID), user.Role, nil
}

func (u *user) Register(ctx context.Context, req entity.RegisterReq) error {
	if err := u.validator.Struct(req); err != nil {
		u.dep.Log.Errorf("[ERROR] WHEN VALIDATE Regis REQ, Error: %v", err)
		return errorr.NewBad("Request body not valid")
	}
	_, err := u.repo.FindByUsername(u.dep.Db.WithContext(ctx), req.Username)
	if err == nil {
		return errorr.NewBad("Username already registered")
	}
	user, err := u.repo.FindByEmail(u.dep.Db.WithContext(ctx), req.Email)
	if err == nil {
		if user.IsVerified == true {
			return errorr.NewBad("Email already registered")
		}
	}
	passhash, err := helper.HashPassword(req.Password)
	if err != nil {
		u.dep.Log.Errorf("Erorr service: %v", err)
		return errorr.NewBad("Register failed")
	}
	hashedEmailString := base32.StdEncoding.EncodeToString([]byte(req.Email))
	data := entity.User{
		Username:         req.Username,
		Email:            req.Email,
		Address:          req.Address,
		Password:         passhash,
		FirstName:        req.FirstName,
		SureName:         req.LastName,
		Role:             req.Role,
		IsVerified:       false,
		VerificationCode: hashedEmailString,
	}
	go func() {
		err := u.dep.Nsq.Publish("5", []byte(hashedEmailString))
		if err != nil {
			u.dep.Log.Errorf("[FAILED] to publish to NSQ: %v", err)
			return
		}
	}()
	err = u.repo.Create(u.dep.Db.WithContext(ctx), data)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) VerifyEmail(ctx context.Context, verificationcode string) error {
	if err := u.repo.VerifyEmail(u.dep.Db.WithContext(ctx), verificationcode); err != nil {
		return err
	}
	return nil
}
func (u *user) ForgetPass(ctx context.Context, email string) error {
	user, err := u.repo.FindByEmail(u.dep.Db.WithContext(ctx), email)
	if err != nil {
		return errorr.NewBad("Email not registered")
	}
	if user.IsVerified == false {
		return errorr.NewBad("Email not verified")
	}
	hashedEmailString := base32.StdEncoding.EncodeToString([]byte(user.Email))
	if err := u.repo.InsertForgotPassToken(u.dep.Db.WithContext(ctx), entity.ForgotPass{Token: hashedEmailString, Email: user.Email}); err != nil {
		return err
	}
	go func() {
		err := u.dep.Nsq.Publish("6", []byte(hashedEmailString))
		if err != nil {
			u.dep.Log.Errorf("[FAILED] to publish to NSQ: %v", err)
			return
		}
	}()
	return nil

}

func (u *user) ResetPass(ctx context.Context, token string, newpass string) error {
	if err := u.repo.ResetPass(u.dep.Db.WithContext(ctx), newpass, token); err != nil {
		return err
	}
	return nil
}
