package devops

import (
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	initConfig "k8s/config"

	"k8s.io/client-go/kubernetes"
)

func GetHornetqLables(meta HornetqMeta, node NodeList) map[string]string {
	lables := map[string]string{}
	lables[meta.Metadata.Name] = meta.Metadata.Name
	lables["serviceName"] = node.Name
	return lables
}

func HornetqController(meta HornetqMeta) (err error) {
	var (
		deployment *appsv1.Deployment
		service    *apiv1.Service
		client     *kubernetes.Clientset
	)
	if client, err = initConfig.InitClient(); err != nil {
		return
	}
	for _, node := range meta.Spec.NodeList {
		if deployment, err = HornetqDeployment(meta, node); err != nil {
			return
		}
		if service, err = HornetqService(meta, node); err != nil {
			return
		}
		// 判断deployment 是否存在，不存在则新建,存在则更新
		if _, err1 := client.AppsV1().Deployments(meta.Metadata.Namespace).Get(node.Name, metav1.GetOptions{}); err1 != nil && errors.IsNotFound(err1) {
			if _, err = client.AppsV1().Deployments(meta.Metadata.Namespace).Create(deployment); err != nil {
				return
			}
		} else {
			if _, err = client.AppsV1().Deployments(meta.Metadata.Namespace).Update(deployment); err != nil {
				return
			}
		}
		// 判断service 是否存在，不存在则新建,存在则更新
		if svc, err1 := client.CoreV1().Services(meta.Metadata.Namespace).Get(node.Name, metav1.GetOptions{}); err1 != nil && errors.IsNotFound(err1) {
			if _, err = client.CoreV1().Services(meta.Metadata.Namespace).Create(service); err != nil {
				return
			}
		} else {
			// 判断存在的和新建的是否存在不一样
			if !reflect.DeepEqual(svc, service) {
				svc.Spec.Ports = service.Spec.Ports
				svc.Spec.Selector = service.Spec.Selector
				if _, err = client.CoreV1().Services(meta.Metadata.Namespace).Update(svc); err != nil {
					return
				}
			}
		}

	}
	return

}

func HornetqDeployment(meta HornetqMeta, node NodeList) (deployment *appsv1.Deployment, err error) {
	//测试环境副本数固定是1
	var (
		replicas int32 = 1
		env      []apiv1.EnvVar
	)
	//var volume []apiv1.Volume
	//var volumeMount []apiv1.VolumeMount
	// 自定义lable标签
	lables := GetHornetqLables(meta, node)
	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      node.Name,
			Namespace: meta.Metadata.Namespace,
			Labels:    GetHornetqLables(meta, node),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: lables,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: lables,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            meta.Metadata.Name,
							ImagePullPolicy: apiv1.PullAlways,
							Image:           meta.Spec.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: int32(node.Port),
								},
							},
						},
					},
				},
			},
		},
	}
	// 添加容器环境变量
	//env = meta.PodEnv
	env = append(env, AddEnvString("LANG", "en_US.UTF-8"))
	env = append(env, AddEnvString("TZ", "Asia/Shanghai"))
	//env = append(env, AddPodNameEnv("POD_NAME", "name"))
	//env = append(env, AddPodNameEnv("POD_NAMESPACE", "namespace"))
	deployment.Spec.Template.Spec.Containers[0].Env = env
	// 处理容器日志持久化到node节点，默认日志路径 /logs node节点存放日志 ${meta.NodeLogDir} / namespace /deployment/ podname
	//volume = append(volume, AddHostVolume(meta.NodeLogDir, "hostpath"))
	//volumeMount = append(volumeMount, AddVolumeMount(meta.ServiceName, meta.LogDir, "hostpath"))
	//deployment.Spec.Template.Spec.Volumes = volume
	//deployment.Spec.Template.Spec.Containers[0].VolumeMounts = volumeMount
	deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &apiv1.Probe{
		Handler: apiv1.Handler{
			TCPSocket: &apiv1.TCPSocketAction{
				Port: intstr.FromInt(node.Port),
			},
		},
		InitialDelaySeconds: 10, // 容器启动多长时间开始使用Liveness探针
		TimeoutSeconds:      3,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
	deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &apiv1.Probe{
		Handler: apiv1.Handler{
			TCPSocket: &apiv1.TCPSocketAction{
				Port: intstr.FromInt(node.Port),
			},
		},
		InitialDelaySeconds: 10, // 容器启动多长时间开始使用Readiness探针
		TimeoutSeconds:      3,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
	return
}

func HornetqService(meta HornetqMeta, node NodeList) (service *apiv1.Service, err error) {
	var (
		Ports []apiv1.ServicePort
	)

	Ports = append(Ports, apiv1.ServicePort{
		Name:     meta.Metadata.Name,
		Protocol: "TCP",
		Port:     int32(node.Port),
		TargetPort: intstr.IntOrString{
			IntVal: int32(node.Port),
		},
	})
	service = &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.Metadata.Name,
			Namespace: meta.Metadata.Namespace,
			Labels:    GetHornetqLables(meta, node),
		},
		Spec: apiv1.ServiceSpec{
			Selector: GetHornetqLables(meta, node),
			Type:     apiv1.ServiceTypeClusterIP,
			//Type:                  apiv1.ServiceTypeNodePort,
			//ExternalTrafficPolicy: apiv1.ServiceExternalTrafficPolicyTypeLocal,
			Ports: Ports,
			//SessionAffinity: "ClientIP",
		},
	}
	return
}
