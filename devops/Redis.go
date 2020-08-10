package devops

type RedisData struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type RedisSpec struct {
	Image     string     `json:"image"`
	Args      RedisArgs  `json:"args"`
	NodeList  []NodeList `json:"nodeList"`
	IsCluster bool       `json:"isCluster"`
}

type RedisArgs struct {
	Args []string `json:"args"`
}

type RedisMeta struct {
	ApiVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Metadata   Metadata  `json:"metadata"`
	Spec       RedisSpec `json:"spec"`
}
