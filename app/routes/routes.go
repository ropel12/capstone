package routes

import (
	"net/http"

	schoolhand "github.com/education-hub/BE/app/features/school/handler"
	userhand "github.com/education-hub/BE/app/features/user/handler"
	"github.com/education-hub/BE/config/dependency"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
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
	ro.Validator = &CustomValidator{validator: validator.New(), log: r.Depend.Log}
	ro.Use(middleware.RemoveTrailingSlash())
	ro.Use(middleware.Logger())
	ro.Use(middleware.Recover())
	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{"https://education-hub-fe-3q5c.vercel.app"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}
	ro.Use(middleware.CORSWithConfig(corsConfig))
	//No Auth
	ro.POST("/login", r.User.Login)
	ro.POST("/register", r.User.Register)
	ro.GET("/verify/:verifcode", r.User.Verify)
	ro.GET("/updateverif/:verifcode", r.User.UpdateVerif)
	ro.POST("/forgot", r.User.Forgotpass)
	ro.POST("/reset/:token", r.User.ResetPass)
	ro.GET("/getcaptcha", r.User.GetCaptcha)
	ro.POST("/verifycaptcha", r.User.VerifyCaptcha)

	ro.GET("/school/search", r.School.Search)
	ro.GET("/gmeet", r.School.CreateGmeet)
	// AUTH
	rauth := ro.Group("", middleware.JWT([]byte(r.Depend.Config.JwtSecret)))
	//user
	rauth.PUT("/users", r.User.Update)
	rauth.DELETE("/users", r.User.Delete)
	rauth.GET("/users", r.User.GetProfile)
	//school
	rverif := rauth.Group("", StatusVerifiedMiddleWare)
	//ADMIN AREA
	radm := rverif.Group("", AdminMiddleWare)
	radm.POST("/school", r.School.Create)
	radm.PUT("/school", r.School.Update)
	radm.POST("/achievements", r.School.AddAchievement)
	radm.PUT("/achievements", r.School.UpdateAchievement)
	radm.DELETE("/achievements/:id", r.School.DeleteAchievement)
	radm.POST("/faqs", r.School.AddFaq)
	radm.PUT("/faqs", r.School.UpdateFaq)
	radm.DELETE("/faqs/:id", r.School.DeleteFaq)
	radm.POST("/extracurriculars", r.School.AddExtracurricular)
	radm.PUT("/extracurriculars", r.School.UpdateExtracurricular)
	radm.DELETE("/extracurriculars/:id", r.School.DeleteExtracurricular)
	radm.POST("/gmeet", r.School.GenerateUrl)
	radm.POST("/payments", r.School.AddPayment)
	radm.PUT("/payments", r.School.UpdatePayment)
	radm.DELETE("/payments/:id", r.School.DeletePayment)
}
