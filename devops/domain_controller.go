package devops

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	initConfig "k8s/config"
)

func DomainGetController(namespace, name string) (svc *v1.Service, err error) {
	var (
		client *kubernetes.Clientset
	)
	// 固定svc
	name = "prop-server"
	if client, err = initConfig.InitClient(); err != nil {
		return
	}

	if svc, err = client.CoreV1().Services(namespace).Get(name, metav1.GetOptions{}); err != nil {
		return
	}

	return
}
