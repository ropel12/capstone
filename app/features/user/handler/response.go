package handler

import (
	"net/http"

	"github.com/education-hub/BE/errorr"
	"github.com/labstack/echo/v4"
)

const (
	URLFRONTEND       = "https://education-hub-fe-3q5c.vercel.app/login"
	URLFRONTENDUPDATE = "http://localhost:8000?status=success"
)

type (
	WebResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
	}
)

func CreateWebResponse(code int, message string, data any) any {
	return WebResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func CreateErrorResponse(err error, c echo.Context) error {
	if err, ok := err.(errorr.BadRequest); ok {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, err.Error(), nil))
	}
	return c.JSON(http.StatusInternalServerError, CreateWebResponse(http.StatusInternalServerError, err.Error(), nil))
}
