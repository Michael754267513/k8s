package devops

import "k8s.io/api/core/v1"

type SpbootSpec struct {
	Layer string `json:"layer"`
	Jdk   string `json:"jdk"`
	Image string `json:"image"`
	//Args   SpbootArgs 	`json:"args"`
	NodeList  []NodeList     `json:"nodeList"`
	IsCluster bool           `json:"isCluster"`
	DB        []DB           `json:"db"`
	Hosts     []v1.HostAlias `json:"hosts"`
	Directory []string       `json:"directory"`
	Listen    []int          `json:"listen"`
	Replicas  int32          `json:"replicas"`
}

type SpbootArgs struct {
	//Args []string		`json:"args"`

}

type DB struct {
	Jndi        string    `json:"jndi"`
	Address     []Address `json:"address"`
	ServiceName string    `json:"service_name"`
	Name        string    `json:"name"`
	Schema      string    `json:"schema"`
}

type Address struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type SpbootMeta struct {
	ApiVersion string     `json:"apiVersion"`
	Kind       string     `json:"kind"`
	Metadata   Metadata   `json:"metadata"`
	Spec       SpbootSpec `json:"spec"`
}
