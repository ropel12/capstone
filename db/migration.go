package db

import (
	entity "github.com/education-hub/BE/app/entities"
	"github.com/education-hub/BE/config"
)

func Migrate(c *config.Config) {
	db, err := config.GetConnection(c)
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(entity.User{}, entity.ForgotPass{}, entity.School{}, entity.Achievement{}, entity.Extracurricular{}, entity.Faq{}, entity.Payment{}); err != nil {
		panic(err)
	}
}
