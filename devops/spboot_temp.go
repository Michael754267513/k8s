package devops

import (
	"fmt"
	"io/ioutil"
	"os"
	//"os/exec"

	"github.com/gogf/gf/util/guuid"
)

func BulidImage(meta BuildMeta) (image string, err error) {
	var (
		path string
	)
	image = "harbor.devops.hpay/java/" + meta.Metadata.Name + ":" + guuid.New().String()
	if path, err = GetDockerFile(meta); err != nil {
		return
	}
	if err = GetImage(); err != nil {
		return
	}
	fmt.Println(path)
	// 进行打包镜像
	//cmd := "cd "+ path + " && docker build -t " + image + "."
	//_,err = exec.Command("sh","-c",cmd).Output()
	return
}

func GetDockerFile(meta BuildMeta) (path string, err error) {
	var (
		data []byte
	)
	path = "tmp/" + guuid.New().String()
	// docker file 字段
	baseimage := "FROM harbor.devops.hpay/base/hjdk:" + meta.Spec.JDK + "\r\n"
	copyimage := "COPY " + meta.Spec.PKG + " /opt/spboot/" + meta.Spec.PKG + "\r\n"
	startboot := "ENTRYPOINT exec java $JAVA_OPTS /opt/spboot/" + meta.Spec.PKG
	dockerfile := baseimage + copyimage + startboot
	data = []byte(dockerfile)
	_ = os.Mkdir(path, os.ModePerm)
	err = ioutil.WriteFile(path+"/Dockerfile", data, os.ModeAppend)
	if err != nil {
		return path, err
	}

	return path, err
}

func GetImage() (err error) {
	return
}
