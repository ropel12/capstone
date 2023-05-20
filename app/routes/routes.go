package routes

import (
	"text/template"

	schoolhand "github.com/education-hub/BE/app/features/school/handler"
	trxhand "github.com/education-hub/BE/app/features/transaction/handler"
	userhand "github.com/education-hub/BE/app/features/user/handler"
	"github.com/education-hub/BE/config/dependency"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/dig"
)

type Routes struct {
	dig.In
	Depend dependency.Depend
	User   userhand.User
	School schoolhand.School
	Trx    trxhand.Transaction
}

func (r *Routes) RegisterRoutes() {
	ro := r.Depend.Echo
	ro.Validator = &CustomValidator{validator: validator.New(), log: r.Depend.Log}
	ro.Use(MetricsMiddleware)
	ro.Use(middleware.RemoveTrailingSlash())
	ro.Use(middleware.Logger())
	ro.Use(middleware.Recover())
	ro.Use(middleware.CORS())
	ro.GET("/prometheus", echo.WrapHandler(promhttp.Handler()))
	//static
	ro.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("./template/*.html")),
	}
	//No Auth
	ro.GET("/quiz/:url", r.School.PreviewQuiz)
	ro.POST("/login", r.User.Login)
	ro.POST("/register", r.User.Register)
	ro.GET("/verify/:verifcode", r.User.Verify)
	ro.GET("/updateverif/:verifcode", r.User.UpdateVerif)
	ro.POST("/forgot", r.User.Forgotpass)
	ro.POST("/reset/:token", r.User.ResetPass)
	ro.GET("/getcaptcha", r.User.GetCaptcha)
	ro.POST("/verifycaptcha", r.User.VerifyCaptcha)
	//school
	ro.GET("/schools", r.School.GetAll)
	ro.GET("/schools/search", r.School.Search)
	ro.GET("/schools/:id", r.School.GetById)
	ro.GET("/gmeet", r.School.CreateGmeet)
	///Third-Party Payment Notification
	ro.POST("/notif", r.Trx.MidtransNotification)
	// AUTH
	rauth := ro.Group("", middleware.JWT([]byte(r.Depend.Config.JwtSecret)))
	rauth.GET("/quiz/set/:token", r.School.SetNewToken, SuperAdmin)
	//User
	rauth.PUT("/users", r.User.Update)
	rauth.DELETE("/users", r.User.Delete)
	rauth.GET("/users", r.User.GetProfile)
	rauth.GET("/progresses/:id", r.School.GetProgressById)
	rverif := rauth.Group("", StatusVerifiedMiddleWare)

	rstdnt := rverif.Group("", StudentMiddleWare)
	///student Area
	rstdnt.GET("/transactions", r.Trx.GetTransactionStudent)
	rstdnt.GET("/transactions/:id", r.Trx.GetDetailTransaction)
	rstdnt.POST("/transactions/checkout", r.Trx.CreateTransaction)
	rstdnt.POST("/school/register", r.School.CreateSubbmision)
	rstdnt.GET("/users/progress", r.School.GetAllProgressByUid)

	//ADMIN AREA
	radm := rverif.Group("", AdminMiddleWare)
	radm.POST("/school", r.School.Create)
	radm.GET("/admin/school", r.School.GetByUid)
	radm.GET("/admin/admission", r.School.GetAllAdmission)
	radm.GET("/admin/admission/:id", r.School.GetSubmissionByid)
	radm.PUT("/progresses/:id", r.School.UpdateProgressByid)
	radm.DELETE("/school/:id", r.School.Delete)
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
	radm.POST("/quiz", r.School.CreateQuiz)
	radm.GET("/quiz", r.School.GetTestResult)
}
