package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/google/logger"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/formats/pb"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/auth"
)

func gracefulShutdown(s *grpc.Server, ctx context.Context) error {
	ch := make(chan struct{})
	go func() {
		s.GracefulStop()
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	file, log := auth_settings.SetUp()
	//noinspection GoUnhandledErrorResult
	defer file.Close()
	defer log.Close()

	logger.Info("Listening at :" + strconv.Itoa(*auth_settings.AuthPort))
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*auth_settings.AuthPort))
	if err != nil {
		logger.Fatalf(
			"Failed to listen at port %d : %v", *auth_settings.AuthPort, err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServer(s, &AuthService{})
	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	logger.Info("Gracefully shutting down...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := gracefulShutdown(s, ctx); err != nil {
		logger.Error(
			"Failed to shutdown gracefully: %v; shutting down forcefully...", err)
		s.Stop()
		logger.Info("Forceful shutdown sequence complete")
	} else {
		logger.Info("Shutdown sequence complete")
	}
}
