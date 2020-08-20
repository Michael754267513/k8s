package devops

import (
	"github.com/gogf/gf/util/gconv"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	initConfig "k8s/config"
)

func SVCcontroller(meta SVCMeta) (err error) {
	var (
		client *kubernetes.Clientset
	)
	if client, err = initConfig.InitClient(); err != nil {
		return
	}

	if meta.Spec.Type == "local" {
		if _, err1 := client.CoreV1().Services(meta.Metadata.Namespace).Get(meta.Metadata.Name, metav1.GetOptions{}); err1 != nil {
			if errors.IsNotFound(err1) {
				if _, err = client.CoreV1().Services(meta.Metadata.Namespace).Create(SVCservice(meta)); err != nil {
					return
				}
			}
		}
	}

	if meta.Spec.Type == "external" {
		if _, err1 := client.CoreV1().Services(meta.Metadata.Namespace).Get(meta.Metadata.Name, metav1.GetOptions{}); err1 != nil {
			if errors.IsNotFound(err1) {
				if _, err = client.CoreV1().Services(meta.Metadata.Namespace).Create(SVCservice(meta)); err != nil {
					return
				}
			}
		}
		if _, err1 := client.CoreV1().Endpoints(meta.Metadata.Namespace).Get(meta.Metadata.Name, metav1.GetOptions{}); err1 != nil {
			if errors.IsNotFound(err1) {
				if _, err = client.CoreV1().Endpoints(meta.Metadata.Namespace).Create(SVCendpoints(meta)); err != nil {
					return
				}
			}

		} else {
			if _, err = client.CoreV1().Endpoints(meta.Metadata.Namespace).Update(SVCendpoints(meta)); err != nil {
				return
			}
		}
	}

	return

}

func SVCendpoints(meta SVCMeta) (ep *apiv1.Endpoints) {

	var (
		subset  []apiv1.EndpointSubset
		address []apiv1.EndpointAddress
		ports   []apiv1.EndpointPort
	)
	address = append(address, apiv1.EndpointAddress{
		IP: meta.Spec.Address,
	})
	for _, v := range meta.Spec.Port {
		if v == 0 {
			continue
		}
		ports = append(ports, apiv1.EndpointPort{
			Name: meta.Metadata.Name + gconv.String(v),
			Port: v,
		})
	}

	subset = append(subset, apiv1.EndpointSubset{
		Addresses: address,
		Ports:     ports,
	})

	ep = &apiv1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: meta.Metadata.Namespace,
			Name:      meta.Metadata.Name,
			Labels: map[string]string{
				"serviceName":      meta.Metadata.Name,
				meta.Metadata.Name: meta.Metadata.Name,
			},
		},
		Subsets: subset,
	}

	return
}

func SVCservice(meta SVCMeta) (service *apiv1.Service) {
	var (
		Ports []apiv1.ServicePort
	)

	for _, v := range meta.Spec.Port {
		if v == 0 {
			continue
		}
		Ports = append(Ports, apiv1.ServicePort{
			Name:     meta.Metadata.Name + gconv.String(v),
			Protocol: "TCP",
			Port:     v,
			TargetPort: intstr.IntOrString{
				IntVal: v,
			},
		})
	}

	if meta.Spec.Type == "local" {
		service = &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      meta.Metadata.Name,
				Namespace: meta.Metadata.Namespace,
				Labels: map[string]string{
					"serviceName":      meta.Metadata.Name,
					meta.Metadata.Name: meta.Metadata.Name,
				},
			},
			Spec: apiv1.ServiceSpec{
				Selector: map[string]string{
					"serviceName":      meta.Metadata.Name,
					meta.Metadata.Name: meta.Metadata.Name,
				},
				Type:  apiv1.ServiceTypeClusterIP,
				Ports: Ports,
				//SessionAffinity: "ClientIP",
			},
		}
	}

	if meta.Spec.Type == "external" {
		service = &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      meta.Metadata.Name,
				Namespace: meta.Metadata.Namespace,
				Labels: map[string]string{
					"serviceName":      meta.Metadata.Name,
					meta.Metadata.Name: meta.Metadata.Name,
				},
			},
			Spec: apiv1.ServiceSpec{
				Type:  apiv1.ServiceTypeClusterIP,
				Ports: Ports,
			},
		}

	}

	return
}

func SVCget(meta SVCMeta) (service *apiv1.Service, err error) {
	var (
		client *kubernetes.Clientset
	)
	if client, err = initConfig.InitClient(); err != nil {
		return
	}

	if service, err = client.CoreV1().Services(meta.Metadata.Namespace).Get(meta.Metadata.Name, metav1.GetOptions{}); err != nil {
		return
	}
	return
}
