package devops

import (
	"encoding/json"
	"fmt"
	"os/exec"

	//"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/gmvc"
)

type DevOpsController struct {
	gmvc.Controller
}

func (r *DevOpsController) Get() {
	r.Response.WriteJson("GET")

}

func (r *DevOpsController) Post() {

	var (
		//data *gjson.Json
		body RequestBody
		err  error
	)

	data := r.Request.GetBody()
	if err = json.Unmarshal(data, &body); err != nil {
		r.Response.Status = 500
		r.Response.Write(err)
		return
	}
	// 判断 namespace 是否存在
	IsExistNamespace(body.Metadata.Namespace)
	if body.Kind == "Redis" || body.Kind == "Zookeeper" {
		if filename, err := SaveFile(r.Request.GetBody()); err != nil {
			r.Response.Status = 500
			r.Response.Write(err)
			return
		} else {
			cmd := "kubectl apply -f " + filename
			f, err := exec.Command("sh", "-c", cmd).Output()
			if err != nil {
				r.Response.Status = 500
				r.Response.Write(err)
				return
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
			r.Response.Status = 500
			r.Response.Write(err)
		} else {
			if err = SpbootController(meta); err != nil {
				r.Response.Status = 500
				r.Response.Write(err)
				return
			}

		}

	}
	r.Response.Write(err)
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
