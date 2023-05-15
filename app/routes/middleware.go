package routes

import (
	"net/http"

	"github.com/education-hub/BE/errorr"
	"github.com/education-hub/BE/helper"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func AdminMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := helper.GetRole(c.Get("user").(*jwt.Token))
		cookie, err := c.Request().Cookie("role")
		if err != nil {
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, map[string]any{"code": 401, "message": "Cookie Not Found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{"code": 500, "message": "Internal Server Error"})
		}
		if role != cookie.Value {
			return c.JSON(http.StatusUnauthorized, map[string]any{"code": 401, "message": "UnAuthorization"})
		}
		if role != "admin" {
			return c.JSON(http.StatusUnauthorized, map[string]any{"code": 401, "message": "UnAuthorization"})
		}
		return next(c)
	}
}
func StatusVerifiedMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isverified := helper.GetStatus(c.Get("user").(*jwt.Token))
		cookie, err := c.Request().Cookie("verified")
		if err != nil {
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, map[string]any{"code": 401, "message": "Cookie Not Found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{"code": 500, "message": "Internal Server Error"})
		}
		if cookie.Value != isverified {
			return c.JSON(http.StatusUnauthorized, map[string]any{"code": 401, "message": "UnAuthorization"})
		}
		if isverified != "true" {
			return c.JSON(http.StatusUnauthorized, map[string]any{"code": 400, "message": "Email Not Verified"})
		}
		return next(c)
	}
}

type CustomValidator struct {
	validator *validator.Validate
	log       *logrus.Logger
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		cv.log.Errorf("[ERROR]WHEN Validate CREATE GMEET REQ, Err: %v", err)
		return errorr.NewBad("Missing Request Body")
	}
	return nil
}
