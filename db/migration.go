package db

import (
	user "github.com/education-hub/BE/app/entities/user"
	"github.com/education-hub/BE/config"
)

func Migrate(c *config.Config) {
	db, err := config.GetConnection(c)
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(user.User{}, user.ForgotPass{}); err != nil {
		panic(err)
	}
}
