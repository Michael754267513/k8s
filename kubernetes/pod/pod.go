package podws

import (
	"github.com/gogf/gf/frame/gmvc"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	initConfig "k8s/config"
)

type PodController struct {
	gmvc.Controller
}

// 定义pods数据返回结构体
type Pods struct {
	ContainerName string
	Message       string
	NameSpace     string
	PodName       string
	PodIP         string
	HostIP        string
	RestartCount  int32
	Status        core_v1.PodPhase
}

func (r *PodController) List() {
	namespace := r.Request.GetString("namespace")
	var (
		clientset  *kubernetes.Clientset
		podsList   *core_v1.PodList
		resPodList []Pods
		pod        Pods
		err        error
	)
	if clientset, err = initConfig.InitClient(); err != nil {
		goto ERROR
	}
	if podsList, err = clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{}); err != nil {
		goto ERROR
	}
	// 遍历获取到的所有pod添加到pods结构体内，然后append到pods数组
	for _, v := range podsList.Items {
		pod.ContainerName = v.Spec.Containers[0].Name
		pod.NameSpace = v.Namespace
		pod.PodName = v.Name
		pod.PodIP = v.Status.PodIP
		pod.HostIP = v.Status.HostIP
		pod.Status = v.Status.Phase
		pod.Message = v.Status.Message
		pod.RestartCount = v.Status.ContainerStatuses[0].RestartCount
		resPodList = append(resPodList, pod)
	}
	r.Response.WriteJson(resPodList)
ERROR:
	r.Response.Write(err)
}
