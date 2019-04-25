package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/google/logger"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/server"
)

func main() {
	file, log, conn := server_settings.SetUp()
	//noinspection GoUnhandledErrorResult
	defer file.Close()
	defer log.Close()
	//noinspection GoUnhandledErrorResult
	defer conn.Close()

	s := Server(*server_settings.ServerPort)
	logger.Info("Listening at ", s.Addr)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	logger.Info("Gracefully shutting down...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := s.Shutdown(ctx); err != nil {
		logger.Error(
			"Failed to shutdown gracefully: %v; shutting down forcefully...", err)
		if err := s.Close(); err != nil {
			logger.Fatalf(
				"Failed to shutdown forcefully: %v; ignoring errors and shutting down", err)
		}
	} else {
		logger.Info("Shutdown sequence complete")
	}
}
