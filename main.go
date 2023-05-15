package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/education-hub/BE/app/routes"
	"github.com/education-hub/BE/config/dependency"
	"github.com/education-hub/BE/config/dependency/container"
	"github.com/education-hub/BE/db"
)

func main() {
	container.RunAll()
	err := container.Container.Invoke(func(depend dependency.Depend, ro routes.Routes) {
		db.Migrate(depend.Config)
		var sig = make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		ro.RegisterRoutes()
		go func() {
			depend.Log.Infof("Starting server on port %s", depend.Config.Server.Port)
			if err := depend.Echo.Start(fmt.Sprintf(":%s", depend.Config.Server.Port)); err != nil {
				depend.Log.Errorf("Failed to start server: %v", err)
				sig <- syscall.SIGTERM
			}
		}()
		<-sig
		depend.Nsq.Stop()
		depend.Log.Info("Shutting down server")
	})
	if err != nil {
		log.Print(err)
	}
}
