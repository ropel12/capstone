package handler

import (
	"net/http"
	"strconv"

	entity "github.com/education-hub/BE/app/entities"
	"github.com/education-hub/BE/app/features/transaction/service"
	"github.com/education-hub/BE/config/dependency"
	"github.com/education-hub/BE/helper"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type Transaction struct {
	dig.In
	Service service.TransactionService
	Dep     dependency.Depend
}

func (u *Transaction) GetTransactionStudent(c echo.Context) error {
	uid := helper.GetUid(c.Get("user").(*jwt.Token))
	res, err := u.Service.GetAllTrasactionCart(c.Request().Context(), uid)
	if err != nil {
		c.Set("err", err)
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}

func (u *Transaction) CreateTransaction(c echo.Context) error {
	req := entity.ReqCheckout{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err)
		u.Dep.Log.Errorf("Error service: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	res, err := u.Service.CreateTransaction(c.Request().Context(), req, helper.GetUid(c.Get("user").(*jwt.Token)))
	if err != nil {
		c.Set("err", err)
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusCreated, CreateWebResponse(http.StatusCreated, "Status Created", map[string]any{"data": res}))
}

func (u *Transaction) GetDetailTransaction(c echo.Context) error {
	schid := c.Param("id")
	if schid == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "School Id is missing", nil))
	}
	newschid, err := strconv.Atoi(schid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Progress Id", nil))
	}
	res, err := u.Service.GetDetailTransaction(c.Request().Context(), newschid, helper.GetUid(c.Get("user").(*jwt.Token)))
	if err != nil {
		c.Set("err", err)
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusCreated, CreateWebResponse(http.StatusCreated, "Status Created", res))
}

func (u *Transaction) MidtransNotification(c echo.Context) error {
	midres := MidtransNotifResponse{}
	if err := c.Bind(&midres); err != nil {
		c.Set("err", err)
		u.Dep.Log.Errorf("[ERROR] When Binding Midtrans Reponse : %v", err)
	}
	switch midres.TransactionStatus {
	case "settlement":
		if err := u.Service.UpdateStatus(c.Request().Context(), "paid", midres.OrderID); err != nil {
			u.Dep.Log.Errorf("[ERROR]When update settlement status: %v", err)
		}
	case "expire":
		if err := u.Service.UpdateStatus(c.Request().Context(), "cancel", midres.OrderID); err != nil {
			u.Dep.Log.Errorf("[ERROR]When update Expire status: %v", err)
		}

	}
	return nil
}
