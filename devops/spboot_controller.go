package devops

import (
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	initConfig "k8s/config"
)

func AddEnvString(name string, value string) apiv1.EnvVar {
	var env apiv1.EnvVar
	env.Name = name
	env.Value = value
	return env
}

func GetLables(meta SpbootMeta) map[string]string {
	lables := map[string]string{}
	lables["jdkVersion"] = meta.Spec.Jdk
	lables[meta.Metadata.Name] = meta.Metadata.Name
	lables["serviceName"] = meta.Metadata.Name
	return lables
}

func SpbootController(meta SpbootMeta) (err error) {
	var (
		deployment *appsv1.Deployment
		service    *apiv1.Service
		client     *kubernetes.Clientset
	)
	if deployment, err = SpbootDeployment(meta); err != nil {
		return
	}

	if service, err = SpbootService(meta); err != nil {
		return
	}

	if client, err = initConfig.InitClient(); err != nil {
		return
	}

	// 判断deployment 是否存在，不存在则新建,存在则更新
	if _, err1 := client.AppsV1().Deployments(meta.Metadata.Namespace).Get(meta.Metadata.Name, metav1.GetOptions{}); err1 != nil && errors.IsNotFound(err1) {
		if _, err = client.AppsV1().Deployments(meta.Metadata.Namespace).Create(deployment); err != nil {
			return
		}
	} else {
		if _, err = client.AppsV1().Deployments(meta.Metadata.Namespace).Update(deployment); err != nil {
			return
		}
	}
	// 判断service 是否存在，不存在则新建,存在则更新
	if svc, err1 := client.CoreV1().Services(meta.Metadata.Namespace).Get(meta.Metadata.Name, metav1.GetOptions{}); err1 != nil && errors.IsNotFound(err1) {
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

	return

}

func SpbootDeleteController(meta SpbootMeta) (err error) {
	var (
		client *kubernetes.Clientset
	)
	client, err = initConfig.InitClient()
	if err = client.AppsV1().Deployments(meta.Metadata.Namespace).Delete(meta.Metadata.Name, &metav1.DeleteOptions{}); err != nil {
		return
	}
	if err = client.CoreV1().Services(meta.Metadata.Namespace).Delete(meta.Metadata.Name, &metav1.DeleteOptions{}); err != nil {
		return
	}

	return
}

func SpbootDeployment(meta SpbootMeta) (deployment *appsv1.Deployment, err error) {
	//测试环境副本数固定是1
	var (
		replicas int32 = 1
		env      []apiv1.EnvVar
	)
	//var volume []apiv1.Volume
	//var volumeMount []apiv1.VolumeMount
	// 判断值是否存在
	if meta.Spec.Replicas != 0 {
		replicas = meta.Spec.Replicas
	}
	// 自定义lable标签
	lables := GetLables(meta)
	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.Metadata.Name,
			Namespace: meta.Metadata.Namespace,
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
							Name:  meta.Metadata.Name,
							Image: meta.Spec.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 8000,
								},
							},
						},
					},
				},
			},
		},
	}
	// 添加hosts解析
	deployment.Spec.Template.Spec.HostAliases = meta.Spec.Hosts
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
	return
}

func SpbootService(meta SpbootMeta) (service *apiv1.Service, err error) {
	var (
		Ports []apiv1.ServicePort
	)

	Ports = append(Ports, apiv1.ServicePort{
		Name:     meta.Metadata.Name,
		Protocol: "TCP",
		Port:     8000,
		TargetPort: intstr.IntOrString{
			IntVal: 8000,
		},
	})
	service = &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.Metadata.Name,
			Namespace: meta.Metadata.Namespace,
		},
		Spec: apiv1.ServiceSpec{
			Selector: GetLables(meta),
			Type:     apiv1.ServiceTypeClusterIP,
			Ports:    Ports,
			//SessionAffinity: "ClientIP",
		},
	}
	return
}
