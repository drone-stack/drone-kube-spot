package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/drone-stack/drone-kube-spot/pkg/server"
	"github.com/sirupsen/logrus"
)

func init() {
	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	DisableColors: true,
	// 	FullTimestamp: true,
	// })
	// logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ctx.Done()
		stop()
	}()

	if err := server.Serve(ctx); err != nil {
		logrus.Fatalf("run serve: %v", err)
	}
}
