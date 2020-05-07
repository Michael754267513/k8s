package main

import (
	"github.com/gogf/gf/frame/g"
	_ "k8s/router"
)

func main() {
	g.Server().Run()
}
