package devops

type SVCMeta struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       SVCSpec  `json:"spec"`
}

type SVCSpec struct {
	Port    []int32 `json:"port"`
	Type    string  `json:"type"`
	Address string  `json:"address"`
}
