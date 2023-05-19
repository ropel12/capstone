package routes

import (
	"net/http"
	"strconv"

	"github.com/education-hub/BE/errorr"
	"github.com/education-hub/BE/helper"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func AdminMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := helper.GetRole(c.Get("user").(*jwt.Token))
		if role != "administrator" {
			return c.JSON(http.StatusUnauthorized, map[string]any{"code": 401, "message": "UnAuthorization"})
		}
		return next(c)
	}
}
func StudentMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := helper.GetRole(c.Get("user").(*jwt.Token))
		if role != "student" {
			return c.JSON(http.StatusUnauthorized, map[string]any{"code": 401, "message": "UnAuthorization"})
		}
		return next(c)
	}
}
func StatusVerifiedMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isverified := helper.GetStatus(c.Get("user").(*jwt.Token))
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

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path", "method"})

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status",
			Help: "Status of HTTP response",
		},
		[]string{"path", "method", "status"})

	httpLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_time",
			Help: "Duration of HTTP requests.",
		}, []string{"path", "method"})

	httpError = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Number of errors.",
		}, []string{"path", "method", "exception"})
)

func init() {
	prometheus.Register(httpLatency)
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(httpError)
	prometheus.MustRegister(responseStatus)
}
func MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		method := c.Request().Method

		timer := prometheus.NewTimer(httpLatency.WithLabelValues(path, method))

		defer timer.ObserveDuration()
		totalRequests.WithLabelValues(path, method).Inc()
		next(c)
		statusCode := c.Response().Status
		responseStatus.WithLabelValues(path, method, strconv.Itoa(statusCode)).Inc()
		if statusCode >= 400 {
			exception := c.Get("err").(string)
			httpError.WithLabelValues(path, method, exception).Inc()
		}
		return nil
	}
}
