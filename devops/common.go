package devops

type RequestBody struct {
	Kind     string   `json:"kind"`
	Metadata Metadata `json:"metadata"`
}

type NodeList struct {
	Port int    `json:"port"`
	Name string `json:"host"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
