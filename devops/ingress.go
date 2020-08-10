package devops

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/gogf/gf/frame/gmvc"
)

type IngressController struct {
	gmvc.Controller
}

func (r *IngressController) Get() {
	r.Response.WriteJson("GET")

}

func (r *IngressController) Post() {
	var (
		meta BuildMeta
	)
	data := r.Request.GetBody()
	if err := json.Unmarshal(data, &meta); err != nil {
		r.Response.Status = 500
		r.Response.Write(err)
	} else {
		if meta.Kind == "domain" {

		}
	}

}

func (r *IngressController) Delete() {
	if data, err := r.Request.GetJson(); err != nil {
		r.Response.Status = 500
		r.Response.Write(err)
	} else {
		fmt.Println(string(data.MustToYaml()))
	}
	if filename, err := SaveFile(r.Request.GetBody()); err != nil {
		fmt.Println(filename)
		r.Response.Status = 500
		r.Response.Write(err)
	} else {
		cmd := "kubectl apply -f " + filename
		f, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			r.Response.Status = 500
			r.Response.Write(err)
		} else {
			r.Response.Write(f)
		}
	}

}

func (r *IngressController) Put() {

	r.Response.WriteJson("PUT")

}

func (r *IngressController) Patch() {

	r.Response.WriteJson("PATCH")

}
