package devops

type HornetqMeta struct {
	ApiVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Metadata   Metadata    `json:"metadata"`
	Spec       HornetqSpec `json:"spec"`
}

type HornetqSpec struct {
	Image     string     `json:"image"`
	NodeList  []NodeList `json:"nodeList"`
	IsCluster bool       `json:"isCluster"`
}
