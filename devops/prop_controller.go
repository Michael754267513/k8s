package devops

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	initConfig "k8s/config"
)

func PropController(namespace string) (err error) {
	var (
		client *kubernetes.Clientset
	)
	if client, err = initConfig.InitClient(); err != nil {
		return
	}
	deployment := PropDeployment(namespace)
	service := PropService(namespace)
	if _, err1 := client.CoreV1().Services(namespace).Get("prop-server", metav1.GetOptions{}); err1 != nil {
		if errors.IsNotFound(err1) {
			if _, err = client.CoreV1().Services(namespace).Create(service); err != nil {
				return
			}
		}
	}
	if _, err1 := client.AppsV1().Deployments(namespace).Get("prop-server", metav1.GetOptions{}); err1 != nil {
		if errors.IsNotFound(err1) {
			if _, err = client.AppsV1().Deployments(namespace).Create(deployment); err != nil {
				return
			}
		}
	}
	return
}

func PropDeployment(namespace string) (deployment *appsv1.Deployment) {

	var (
		replicas int32 = 1
	)

	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prop-server",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"prop-server": "prop-server",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"prop-server": "prop-server",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "prop-server",
							Image: "nginx",
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	return

}

func PropService(namespace string) (service *v1.Service) {
	var (
		Ports []v1.ServicePort
	)

	Ports = append(Ports, v1.ServicePort{
		Name:     "prop-server",
		Protocol: "TCP",
		Port:     80,
		TargetPort: intstr.IntOrString{
			IntVal: 80,
		},
	})
	service = &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prop-server",
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"prop-server": "prop-server",
			},
			Type:  v1.ServiceTypeClusterIP,
			Ports: Ports,
			//SessionAffinity: "ClientIP",
		},
	}
	return
}

func PropGetController(namespace, name string) (svc *v1.Service, err error) {
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
