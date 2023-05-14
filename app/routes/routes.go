package routes

import (
	schoolhand "github.com/education-hub/BE/app/features/school/handler"
	userhand "github.com/education-hub/BE/app/features/user/handler"
	"github.com/education-hub/BE/config/dependency"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

type Routes struct {
	dig.In
	Depend dependency.Depend
	User   userhand.User
	School schoolhand.School
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
	ro.GET("/updateverif/:verifcode", r.User.UpdateVerif)
	ro.POST("/forgot", r.User.Forgotpass)
	ro.POST("/reset/:token", r.User.ResetPass)
	ro.GET("/getcaptcha", r.User.GetCaptcha)
	ro.POST("/verifycaptcha", r.User.VerifyCaptcha)

	// AUTH
	rauth := ro.Group("", middleware.JWT([]byte(r.Depend.Config.JwtSecret)))
	//user
	rauth.PUT("/users", r.User.Update)
	rauth.DELETE("/users", r.User.Delete)
	rauth.GET("/users", r.User.GetProfile)
	//school
	rauth.GET("/school/search", r.School.Search)
	rverif := rauth.Group("", StatusVerifiedMiddleWare)
	//ADMIN AREA
	radm := rverif.Group("", AdminMiddleWare)
	radm.POST("/school", r.School.Create)
	radm.PUT("/school", r.School.Update)
}
