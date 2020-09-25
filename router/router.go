package router

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	DevOps "k8s/devops"
	Deploy "k8s/kubernetes/deployment"
	Namespace "k8s/kubernetes/namespace"
	Pod "k8s/kubernetes/pod"
	Service "k8s/kubernetes/service"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func init() {
	s := g.Server()
	s.SetConfigWithMap(g.Map{
		"AccessLogEnabled": true,
		"ErrorLogEnabled":  true,
	})
	s.Group("/kubernetes", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS) // 跨域处理
		// pod 处理段
		group.ALL("/pod/", new(Pod.PodController))
		group.ALL("/pod/", new(Pod.PodWSController))
		// deployment 处理段
		group.ALL("/deployment", new(Deploy.DeployMentController))
		// service 处理段
		group.ALL("/service", new(Service.ServiceController))
		// namespace 处理段
		group.REST("/namespace", new(Namespace.NameSpaceController))
	})
	s.Group("/handpay", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS) // 跨域处理
		group.REST("/devops/kubernetes", new(DevOps.DevOpsController))
		group.REST("/devops/build", new(DevOps.BuildController))
	})
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:9999", nil))
	}()
}
