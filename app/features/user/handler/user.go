package handler

import (
	"net/http"

	entity "github.com/education-hub/BE/app/entities/user"
	"github.com/education-hub/BE/app/features/user/service"
	"github.com/education-hub/BE/config/dependency"
	"github.com/education-hub/BE/helper"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type User struct {
	dig.In
	Service service.UserService
	Dep     dependency.Depend
}

func (u *User) Login(c echo.Context) error {
	var req entity.LoginReq
	var token string
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING LOGIN, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	uid, role, err := u.Service.Login(c.Request().Context(), req)
	token = helper.GenerateJWT(uid, role, u.Dep)

	if err != nil {
		return CreateErrorResponse(err, c)
	}
	c.SetCookie(&http.Cookie{Name: "role", Value: role})
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"token": token}))
}

func (u *User) Register(c echo.Context) error {
	var req entity.RegisterReq
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REGISTER, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	if err := u.Service.Register(c.Request().Context(), req); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusCreated, CreateWebResponse(http.StatusCreated, "Status Created", nil))
}

func (u *User) Verify(c echo.Context) error {
	verifcode := c.Param("verifcode")
	if verifcode == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Missing verifcode in path", nil))
	}
	if err := u.Service.VerifyEmail(c.Request().Context(), verifcode); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}

func (u *User) Forgotpass(c echo.Context) error {
	req := struct {
		Email string `json:"email"`
	}{}
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING FORGOTPASS, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Missing Email Request", nil))
	}
	if err := u.Service.ForgetPass(c.Request().Context(), req.Email); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}

func (u *User) ResetPass(c echo.Context) error {
	req := struct {
		Password string `json:"password"`
	}{}
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING RESETPASS, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	token := c.Param("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Missing token", nil))
	}
	hashpass, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}
	req.Password = hashpass
	if err := u.Service.ResetPass(c.Request().Context(), token, req.Password); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}

func (u *User) GetCaptcha(c echo.Context) error {

	id, captcha, err := helper.GenerateCaptcha()
	if err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN GENERATE CAPTCHA, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Missing token", nil))
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"captchaid": id, "image": captcha}))
}

func (u *User) VerifyCaptcha(c echo.Context) error {
	req := struct {
		CaptchaID string `json:"captcha_id"`
		Value     string `json:"value"`
	}{}
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING CAPTCHA, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	if err := helper.VerifyCaptcha(req.CaptchaID, req.Value); err == true {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Wrong Answer", nil))
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}
