package devops

import (
	"io/ioutil"
	"os"

	"github.com/gogf/gf/util/guuid"
)

func SaveFile(data []byte) (string, error) {
	filename := "/tmp/" + "k8s." + guuid.New().String() + ".yaml"
	err := ioutil.WriteFile(filename, data, os.ModeAppend)
	if err != nil {
		return filename, err
	}
	return filename, nil
}
