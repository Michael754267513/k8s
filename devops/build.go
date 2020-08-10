package devops

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/gogf/gf/frame/gmvc"
)

type MetaData struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type SpecMeta struct {
	Args      []string `json:"args"`
	DB        []string `json:"db"`
	Directory []string `json:"directory"`
	Image     string   `json:"image"`
	IsCluster bool     `json:"isCluster"`
	JDK       string   `json:"jdk"`
	Layer     string   `json:"layer"`
	Listen    []int    `json:"listen"`
	Version   string   `json:"version"`
	PKG       string   `json:"pkg"`
}

type BuildMeta struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   MetaData `json:"metadata"`
	Spec       SpecMeta `json:"spec"`
}

type BuildController struct {
	gmvc.Controller
}

func (r *BuildController) Get() {
	r.Response.WriteJson("GET")
}

func (r *BuildController) Post() {
	var (
		meta BuildMeta
	)
	data := r.Request.GetBody()
	if err := json.Unmarshal(data, &meta); err != nil {
		r.Response.Status = 500
		r.Response.Write(err)
	} else {
		if meta.Kind == "spboot" {
			if image, err := BulidImage(meta); err != nil {
				r.Response.Status = 500
				r.Response.Write(err)
			} else {
				r.Response.WriteJson(image)
			}
		}
	}

}

func (r *BuildController) Delete() {
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

func (r *BuildController) Put() {

	r.Response.WriteJson("PUT")

}

func (r *BuildController) Patch() {

	r.Response.WriteJson("PATCH")

}
