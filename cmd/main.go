package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"story-time-server/database"
	"story-time-server/firebaseProvider"
	"story-time-server/server"
	"syscall"
	"time"
)

const shutDownTimeOut = 10 * time.Second

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// initiate firebase client
	firebaseProvider.InitFirebaseClient()

	// create server instance
	srv := server.SetupRoutes()
	if err := database.InitDB(); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Print("Migration Successful!")

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logrus.Panicf("Failed to run server with error: %+v", err)
		}
	}()
	logrus.Print("Server started")

	<-done

	logrus.Info("shutting down server")
	if err := database.ShutdownDatabase(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}
	if err := srv.Shutdown(shutDownTimeOut); err != nil {
		logrus.WithError(err).Panic("failed to gracefully shutdown server")
	}
}
