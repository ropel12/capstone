package routes

import (
	userhand "github.com/education-hub/BE/app/features/user/handler"
	"github.com/education-hub/BE/config/dependency"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

type Routes struct {
	dig.In
	Depend dependency.Depend
	User   userhand.User
}

func (r *Routes) RegisterRoutes() {
	ro := r.Depend.Echo
	ro.Use(middleware.RemoveTrailingSlash())
	ro.Use(middleware.Logger())
	ro.Use(middleware.Recover())
	ro.Use(middleware.CORS())
	//No Auth
	ro.POST("/login", r.User.Login)
	ro.POST("/register", r.User.Register)
	ro.GET("/verify/:verifcode", r.User.Verify)
	ro.POST("/forgot", r.User.Forgotpass)
	ro.POST("/reset/:token", r.User.ResetPass)
	ro.GET("/getcaptcha", r.User.GetCaptcha)
	ro.POST("/verifycaptcha", r.User.VerifyCaptcha)
}
