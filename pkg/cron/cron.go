package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var Cron *Client

type Client struct {
	client *cron.Cron
}

func New() *Client {
	return &Client{client: cron.New()}
}

func (c *Client) Start() {
	logrus.Info("start cron")
	c.client.Start()
}

func (c *Client) Add(spec string, cmd func()) {
	id, err := c.client.AddFunc(spec, cmd)
	if err != nil {
		logrus.Infof("add cron: %v, err: %v", id, err)
		return
	}
	logrus.Infof("add cron: %v", id)
}

func (c *Client) Stop() {
	logrus.Info("stop cron")
	c.client.Stop()
}

func init() {
	Cron = New()
}
