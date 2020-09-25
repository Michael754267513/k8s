package devops

import (
	"encoding/json"
	"fmt"
	"os/exec"

	//"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/gmvc"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

type DevOpsController struct {
	gmvc.Controller
}

func (r *DevOpsController) Get() {

	var (
		//data *gjson.Json
		body RequestBody
		err  error
	)

	data := r.Request.GetBody()
	if err = json.Unmarshal(data, &body); err != nil {
		goto ERROR
	}

	if body.Metadata.Namespace == "" {
		r.Response.Write("namespace is not null!")
		goto ERROR
	}

	if body.Kind == "Domain" {
		var (
			svc *v1.Service
		)
		if svc, err = PropGetController(body.Metadata.Namespace, body.Metadata.Name); err != nil {
			goto ERROR
		}
		r.Response.WriteJson(svc)
		return
	}

	if body.Kind == "Redis" {
	}

	if body.Kind == "Zookeeper" {
	}

	if body.Kind == "Spboot" {
		var (
			meta       SpbootMeta
			service    *v1.Service
			deployment *appsv1.Deployment
		)
		if err := json.Unmarshal(r.Request.GetBody(), &meta); err != nil {
			goto ERROR
		} else {
			service, deployment, err = SpbootGet(meta)
			r.Response.WriteJson(service)
			r.Response.WriteJson(deployment)
			if err != nil {
				goto ERROR
			}
			return
		}
	}

	if body.Kind == "Hservice" {
		var (
			meta    SVCMeta
			service *v1.Service
		)
		if err = json.Unmarshal(r.Request.GetBody(), &meta); err != nil {
			goto ERROR
		} else {
			if service, err = SVCget(meta); err != nil {
				goto ERROR
			}
			r.Response.WriteJson(service)
			return
		}
	}

	if body.Kind == "Env" {
		// 获取
	}

ERROR:
	fmt.Println(err)
	r.Response.Status = 500
	r.Response.Write(err)
	return
}

func (r *DevOpsController) Post() {

	var (
		//data *gjson.Json
		body RequestBody
		err  error
	)

	data := r.Request.GetBody()
	fmt.Println(string(data))
	if err = json.Unmarshal(data, &body); err != nil {
		goto ERROR
	}

	if body.Metadata.Namespace == "" {
		r.Response.Write("namespace is not null!")
		goto ERROR
	}
	// 判断 namespace 是否存在
	IsExistNamespace(body.Metadata.Namespace)
	if body.Kind == "Redis" || body.Kind == "Zookeeper" {
		if filename, err := SaveFile(r.Request.GetBody()); err != nil {
			goto ERROR
		} else {
			cmd := "kubectl apply -f " + filename
			f, err := exec.Command("sh", "-c", cmd).Output()
			if err != nil {
				goto ERROR
			} else {
				r.Response.Write(f)
				return
			}
		}
	}

	if body.Kind == "Spboot" {
		var (
			meta SpbootMeta
		)
		if err := json.Unmarshal(r.Request.GetBody(), &meta); err != nil {
			goto ERROR
		} else {
			if err = SpbootController(meta); err != nil {
				goto ERROR
			}

		}
		r.Response.WriteJson(err)
		return
	}

	if body.Kind == "Hservice" {
		var (
			meta SVCMeta
		)
		if err = json.Unmarshal(r.Request.GetBody(), &meta); err != nil {
			goto ERROR
		} else {
			if err = SVCcontroller(meta); err != nil {
				goto ERROR
			}
		}
		r.Response.WriteJson(err)
		return
	}

	if body.Kind == "Prop" {
		if err = PropController(body.Metadata.Namespace); err != nil {
			goto ERROR
		}
		r.Response.WriteJson(err)
		return
	}

	if body.Kind == "Hornetq" {
		var (
			meta HornetqMeta
		)
		if err = json.Unmarshal(r.Request.GetBody(), &meta); err != nil {
			goto ERROR
		} else {
			if err = HornetqController(meta); err != nil {
				goto ERROR
			}
		}
		r.Response.WriteJson(err)
		return

	}

	if body.Kind == "Env" {
		// 新建环境
	}

	r.Response.Status = 500
	r.Response.Write("resource not found")
	return
ERROR:
	fmt.Println(err)
	r.Response.Status = 500
	r.Response.Write(err)
	return

}

func (r *DevOpsController) Delete() {

	var (
		//data *gjson.Json
		body     RequestBody
		filename string
		err      error
	)

	data := r.Request.GetBody()
	if err = json.Unmarshal(data, &body); err != nil {
		goto ERROR
	}

	if body.Kind == "Redis" || body.Kind == "Zookeeper" {
		if filename, err = SaveFile(r.Request.GetBody()); err != nil {
			goto ERROR
		} else {
			cmd := "kubectl delete -f " + filename
			f, err := exec.Command("sh", "-c", cmd).Output()
			if err != nil {
				goto ERROR
			} else {
				r.Response.Write(f)
				return
			}

		}

	}

	if body.Kind == "Spboot" {
		var (
			meta SpbootMeta
		)
		if err := json.Unmarshal(r.Request.GetBody(), &meta); err != nil {
			goto ERROR
		} else {
			if err = SpbootDeleteController(meta); err != nil {
				goto ERROR
			}

		}

	}
	r.Response.Write(err)
	return
ERROR:
	fmt.Println(err)
	r.Response.Status = 500
	r.Response.Write(err)
	return

}

func (r *DevOpsController) Put() {

	r.Response.WriteJson("PUT")

}

func (r *DevOpsController) Patch() {

	r.Response.WriteJson("PATCH")

}
