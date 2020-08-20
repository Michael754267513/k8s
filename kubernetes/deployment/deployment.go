package deployment

import (
	"github.com/gogf/gf/frame/gmvc"
	apps_v1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	initConfig "k8s/config"
)

//type DeployMent interface {
//
//}

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

// 定义更新deployment返回参数
type UpdateDeployStatus struct {
	Name      string
	NameSpace string
	Status    bool
	Image     string
	OldImage  string
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
	name := r.Request.GetString("name")
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

	if deploylist, err = clientset.AppsV1().Deployments(namespace).List(meta_v1.ListOptions{
		LabelSelector: name,
	}); err != nil {
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
	if err = r.Response.WriteJson(resDeployList); err != nil {
		goto ERROR
	}
ERROR:
	initConfig.Logger(err)
}

func (r *DeployMentController) Update() {
	namespace := r.Request.GetString("namespace")
	deployment := r.Request.GetString("deployment")
	image := r.Request.GetString("image")
	var (
		clientset          *kubernetes.Clientset
		updateDeployStatus UpdateDeployStatus
		deployMeta         *apps_v1.Deployment
		err                error
	)
	updateDeployStatus.Name = deployment
	updateDeployStatus.NameSpace = namespace
	if clientset, err = initConfig.InitClient(); err != nil {
		updateDeployStatus.Message = err
		goto ERROR
	}
	// 获取deployment相关信息
	if deployMeta, err = clientset.AppsV1().Deployments(namespace).Get(deployment, meta_v1.GetOptions{}); err != nil {
		updateDeployStatus.Message = err
		goto ERROR
	}
	// 历史镜像
	updateDeployStatus.OldImage = deployMeta.Spec.Template.Spec.Containers[0].Image
	updateDeployStatus.Image = image
	// 更新新增镜像配置
	deployMeta.Spec.Template.Spec.Containers[0].Image = image
	if deployMeta, err = clientset.AppsV1().Deployments(namespace).Update(deployMeta); err != nil {
		updateDeployStatus.Message = err
	} else {
		updateDeployStatus.Status = true
	}
	r.Response.Write(updateDeployStatus)
ERROR:
	initConfig.Logger(err)
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
		deployStatus.Message = err
		goto ERROR
	}
	if _, err = clientset.AppsV1().Deployments(namespace).Get(deployment, meta_v1.GetOptions{}); err != nil {
		deployStatus.Status = true
		goto ERROR
	}
	if err = clientset.AppsV1().Deployments(namespace).Delete(deployment, &meta_v1.DeleteOptions{}); err != nil {
		deployStatus.Message = err
	} else {
		deployStatus.Status = true
	}
	// 数据格式化json失败会抛出异常
	if err = r.Response.WriteJson(deployStatus); err != nil {
		goto ERROR
	}
ERROR:
	initConfig.Logger(err)
}

func (r *DeployMentController) Create() {
	/*
		根据实际情况定义参数
		https://github.com/kubernetes/client-go/blob/master/examples/dynamic-create-update-delete-deployment/main.go
	*/

}
