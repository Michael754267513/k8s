package devops

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	initConfig "k8s/config"
)

func IsExistNamespace(namespace string) {

	client, _ := initConfig.InitClient()
	if _, err := client.CoreV1().Namespaces().Get(namespace, v1.GetOptions{}); err != nil {
		if _, err := client.CoreV1().Namespaces().Create(&apiv1.Namespace{
			ObjectMeta: v1.ObjectMeta{
				Name: namespace,
			},
		}); err != nil {
			fmt.Println(err)
		}
	}

}
