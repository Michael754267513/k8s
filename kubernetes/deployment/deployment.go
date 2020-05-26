package deployment

import (
	"github.com/gogf/gf/frame/gmvc"
	apps_v1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	initConfig "k8s/config"
)

type DeployMentController struct {
	gmvc.Controller
}

type Deploys struct {
	Name                string
	NameSpace           string
	Replicas            int32
	Strategy            apps_v1.DeploymentStrategy
	ObservedGeneration  int64
	UpdatedReplices     int32
	ReadyReplicas       int32
	AvailableReplicas   int32
	UnavailableReplicas int32
	Conditions          []apps_v1.DeploymentCondition
}

func (r *DeployMentController) List() {
	namespace := r.Request.GetString("namespace")
	var (
		clientset     *kubernetes.Clientset
		deploylist    *apps_v1.DeploymentList
		resDeployList []Deploys
		deploy        Deploys
		err           error
	)
	if clientset, err = initConfig.InitClient(); err != nil {
		goto ERROR
	}
	if deploylist, err = clientset.AppsV1().Deployments(namespace).List(meta_v1.ListOptions{}); err != nil {
		goto ERROR
	}
	// 获取deployment 相关数据
	for _, v := range deploylist.Items {
		deploy.Name = v.Name
		deploy.AvailableReplicas = v.Status.AvailableReplicas
		deploy.Strategy = v.Spec.Strategy
		deploy.Conditions = v.Status.Conditions
		deploy.ObservedGeneration = v.Status.ObservedGeneration
		deploy.ReadyReplicas = v.Status.ReadyReplicas
		deploy.Replicas = v.Status.Replicas
		deploy.NameSpace = v.Namespace
		deploy.UnavailableReplicas = v.Status.UnavailableReplicas
		deploy.UpdatedReplices = v.Status.UnavailableReplicas
		resDeployList = append(resDeployList, deploy)

	}
	r.Response.WriteJson(resDeployList)
ERROR:
	r.Response.Write(err)
}

