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

// 定义deploy删除返回
type DeployStatus struct {
	Name      string
	NameSpace string
	Status    bool
	Message   error
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

func (r *DeployMentController) Update() {

}

func (r *DeployMentController) Delete() {
	namespace := r.Request.GetString("namespace")
	deployment := r.Request.GetString("deployment")
	var (
		clientset    *kubernetes.Clientset
		deployStatus DeployStatus
		err          error
	)
	deployStatus.Name = deployment
	deployStatus.NameSpace = namespace
	if clientset, err = initConfig.InitClient(); err != nil {
		goto ERROR
	}
	if _, err = clientset.AppsV1().Deployments(namespace).Get(deployment, meta_v1.GetOptions{}); err != nil {
		deployStatus.Status = true
		goto ERROR
	}
	if err = clientset.AppsV1().Deployments(namespace).Delete(deployment, &meta_v1.DeleteOptions{}); err != nil {
		deployStatus.Message = err
		r.Response.WriteJson(deployStatus)
	} else {
		deployStatus.Status = true
		r.Response.WriteJson(deployStatus)
	}
ERROR:
}

func (r *DeployMentController) Create() {

}
