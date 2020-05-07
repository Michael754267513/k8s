package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	Pod "k8s/kubernetes/pod"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func init() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		pod := new(Pod.PodController)
		podws := new(Pod.PodWSController)
		group.Middleware(MiddlewareCORS) // 跨域处理
		group.ALL("/kubernetes/pod", pod)
		group.ALL("/kubernetes/pod", podws)
	})
}
