package namespace

import (
	"github.com/gogf/gf/frame/gmvc"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	initConfig "k8s/config"
)

type NameSpaceController struct {
	gmvc.Controller
}

type NameSpace struct {
	Name   string
	Status core_v1.NamespacePhase
}

func (r *NameSpaceController) List() {
	//namespace := r.Request.GetString("namespace")
	var (
		clientset        *kubernetes.Clientset
		namespaceList    *core_v1.NamespaceList
		resNameSpaceList []NameSpace
		namespace        NameSpace
		err              error
	)
	if clientset, err = initConfig.InitClient(); err != nil {
		goto ERROR
	}
	if namespaceList, err = clientset.CoreV1().Namespaces().List(meta_v1.ListOptions{}); err != nil {
		goto ERROR
	}
	for _, v := range namespaceList.Items {
		namespace.Name = v.Name
		namespace.Status = v.Status.Phase
		resNameSpaceList = append(resNameSpaceList, namespace)
	}
	r.Response.WriteJson(resNameSpaceList)
ERROR:
	initConfig.Logger(err)
}

func (r *NameSpaceController) Delete() {
	namespace := r.Request.GetString("namespace")

	var (
		clientset *kubernetes.Clientset
		err       error
	)
	if clientset, err = initConfig.InitClient(); err != nil {
		goto ERROR
	}

	err = clientset.CoreV1().Namespaces().Delete(namespace, &meta_v1.DeleteOptions{})
	if err != nil {
		goto ERROR
	}

	r.Response.WriteJson("")
	return
ERROR:
	initConfig.Logger(err)
	r.Response.Status = 500
	r.Response.Write(err)
}

func (r *NameSpaceController) GET() {
	namespace := r.Request.GetString("namespace")

	var (
		clientset     *kubernetes.Clientset
		namespaceList *core_v1.Namespace
		err           error
	)
	if clientset, err = initConfig.InitClient(); err != nil {
		goto ERROR
	}

	namespaceList, err = clientset.CoreV1().Namespaces().Get(namespace, meta_v1.GetOptions{})
	if err != nil {
		goto ERROR
	}

	r.Response.WriteJson(namespaceList)
	return
ERROR:
	initConfig.Logger(err)
	r.Response.Status = 500
	r.Response.Write(err)
}

func (r *NameSpaceController) POST() {
	namespace := r.Request.GetString("namespace")

	var (
		clientset     *kubernetes.Clientset
		namespaceList *core_v1.Namespace
		err           error
	)
	if clientset, err = initConfig.InitClient(); err != nil {
		goto ERROR
	}

	namespaceList, err = clientset.CoreV1().Namespaces().Create(&core_v1.Namespace{

		ObjectMeta: meta_v1.ObjectMeta{
			Name:      namespace,
			Namespace: namespace,
		},
	})
	if err != nil {
		goto ERROR
	}

	r.Response.WriteJson(namespaceList)
	return
ERROR:
	initConfig.Logger(err)
	r.Response.Status = 500
	r.Response.Write(err)

}
