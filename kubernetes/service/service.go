package service

import (
	"github.com/gogf/gf/frame/gmvc"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	initConfig "k8s/config"
)

type ServiceController struct {
	gmvc.Controller
}

type Service struct {
	Name                  string
	NameSpace             string
	Ports                 []core_v1.ServicePort
	ClusterIP             string
	Type                  core_v1.ServiceType
	ExternalIPs           []string
	SessionAffinity       core_v1.ServiceAffinity
	LoadBalanceIP         string
	ExternalName          string
	ExternalTrafficPolicy core_v1.ServiceExternalTrafficPolicyType
	SessionAffinityConfig *core_v1.SessionAffinityConfig
	Ingress               []core_v1.LoadBalancerIngress
}

func (r *ServiceController) List() {
	namespace := r.Request.GetString("namespace")
	name := r.Request.GetString("name")
	var (
		clientset      *kubernetes.Clientset
		servicelist    *core_v1.ServiceList
		resServiceList []Service
		service        Service
		err            error
	)
	if clientset, err = initConfig.InitClient(); err != nil {
		goto ERROR
	}
	if servicelist, err = clientset.CoreV1().Services(namespace).List(meta_v1.ListOptions{LabelSelector: name}); err != nil {
		goto ERROR
	}
	for _, v := range servicelist.Items {
		service.Name = v.Name
		service.NameSpace = v.Namespace
		service.SessionAffinityConfig = v.Spec.SessionAffinityConfig
		service.SessionAffinity = v.Spec.SessionAffinity
		service.ClusterIP = v.Spec.ClusterIP
		service.ExternalIPs = v.Spec.ExternalIPs
		service.ExternalName = v.Spec.ExternalName
		service.ExternalTrafficPolicy = v.Spec.ExternalTrafficPolicy
		service.Ingress = v.Status.LoadBalancer.Ingress
		service.LoadBalanceIP = v.Spec.LoadBalancerIP
		service.Type = v.Spec.Type
		service.Ports = v.Spec.Ports
		resServiceList = append(resServiceList, service)
	}
	r.Response.WriteJson(resServiceList)
ERROR:
	initConfig.Logger(err)
}

func (r *ServiceController) Delete() {
	namespace := r.Request.GetString("namespace")
	serviceName := r.Request.GetString("serviceName")
	var (
		clientset *kubernetes.Clientset
		err       error
	)
	if clientset, err = initConfig.InitClient(); err != nil {
		goto ERROR
	}
	if err = clientset.CoreV1().Services(namespace).Delete(serviceName, &meta_v1.DeleteOptions{}); err != nil {
		goto ERROR
	}
	r.Response.Write(err)
ERROR:
	initConfig.Logger(err)
}

func (r *ServiceController) Create() {}

func (r *ServiceController) Update() {}
