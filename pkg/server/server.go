package server

import (
	"context"
	"net/http"
	"os"

	"github.com/drone-stack/drone-kube-spot/pkg/cron"
	"github.com/drone-stack/drone-kube-spot/pkg/k8s"
	"github.com/ergoapi/exgin"
	"github.com/ergoapi/util/ztime"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func Serve(ctx context.Context) error {
	defer cron.Cron.Stop()
	cron.Cron.Start()
	g := exgin.Init(false)
	g.Use(exgin.ExCors())
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	cron.Cron.Add("20 8 * * *", func() {
		holiday := ztime.HolidayGet(ztime.GetToday())
		if holiday.NeedWork {
			logrus.Infof("Today is %s, need work", holiday.Name)
			k8s.Pods()
		} else {
			logrus.Infof("Today is %s, no need work", holiday.Name)
		}
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
