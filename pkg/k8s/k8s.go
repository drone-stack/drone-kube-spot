package k8s

import (
	"context"
	"strings"

	"github.com/drone-stack/drone-kube-spot/internal/kube"
	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/ztime"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var K8SClient kubernetes.Interface

func init() {
	var err error
	kubecfg := &kube.ClientConfig{}
	K8SClient, err = kube.New(kubecfg)
	if err != nil {
		panic(err)
	}
}

func Pods() {
	holiday := ztime.HolidayGet(ztime.GetToday())
	if holiday.NeedWork {
		logrus.Infof("Today is %s, need work", holiday.Name)
	} else {
		logrus.Infof("Today is %s, no need work", holiday.Name)
		return
	}
	label := environ.GetEnv("label", "pool=spot")
	labels := strings.Split(label, "=")
	if len(labels) != 2 {
		return
	}
	key := labels[0]
	value := labels[1]
	nodes, err := K8SClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		panic(err)
	}
	if len(nodes.Items) > 0 {
		logrus.Info("exist node in cluster")
		return
	}
	logrus.Warn("not found in cluster, will create one")
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "spot-node",
			Labels: map[string]string{
				"kube-bot": "spot",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "spot-node",
					Image:           "nginx",
					ImagePullPolicy: corev1.PullIfNotPresent,
				},
			},
			NodeSelector: map[string]string{
				key: value,
			},
			Tolerations: []corev1.Toleration{
				{
					Operator: corev1.TolerationOpExists,
				},
			},
		},
	}
	if _, err := K8SClient.CoreV1().Pods("default").Get(context.Background(), pod.Name, metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			if _, err := K8SClient.CoreV1().Pods("default").Create(context.Background(), &pod, metav1.CreateOptions{}); err != nil {
				panic(err)
			}
			logrus.Info("create spot-node success")
			return
		} else {
			panic(err)
		}
	}
	logrus.Info("exist spot-node")
}

func Clean() {
	err := K8SClient.CoreV1().Pods("default").Delete(context.Background(), "spot-node", metav1.DeleteOptions{})
	if err != nil && errors.IsNotFound(err) {
		logrus.Errorf("delete spot-node error: %s", err)
		return
	}
	logrus.Info("delete spot-node success")
}
