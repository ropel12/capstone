package dependency

import (
	"github.com/education-hub/BE/config"
	"github.com/education-hub/BE/pkg"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Depend struct {
	dig.In
	Db         *gorm.DB
	Config     *config.Config
	Echo       *echo.Echo
	Log        *logrus.Logger
	Gcp        *pkg.StorageGCP
	Rds        *redis.Client
	Nsq        *pkg.NSQProducer
	Validation *pkg.Validation
}
