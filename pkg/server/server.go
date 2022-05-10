package server

import (
	"context"
	"net/http"
	"os"

	"github.com/drone-stack/drone-kube-spot/pkg/cron"
	"github.com/drone-stack/drone-kube-spot/pkg/k8s"
	"github.com/ergoapi/exgin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func Serve(ctx context.Context) error {
	// pod重建时创建一次
	k8s.Pods()
	defer cron.Cron.Stop()
	cron.Cron.Start()
	g := exgin.Init(false)
	g.Use(exgin.ExCors())
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	cron.Cron.Add("20 8 * * *", func() {
		k8s.Pods()
	})
	cron.Cron.Add("0 18 * * *", func() {
		logrus.Info("try clean pods")
		k8s.Clean()
	})
	addr := "0.0.0.0:65001"
	srv := &http.Server{
		Addr:    addr,
		Handler: g,
	}
	logrus.Infof("http listen to %v, pid is %v", addr, os.Getpid())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Errorf("Failed to start http server, error: %s", err)
		return err
	}
	return nil
}
