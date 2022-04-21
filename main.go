package main

import (
	"github.com/ergoapi/util/ztime"
	"github.com/sirupsen/logrus"
)

func main() {
	holiday := ztime.HolidayGet(ztime.GetToday())
	if holiday.NeedWork {
		logrus.Infof("Today is %s, need work", holiday.Name)
	} else {
		logrus.Infof("Today is %s, no need work", holiday.Name)
	}
}
