# Kubernetes GoFrame 
## 介绍
kubernetes  GO 客户端与GoFrame提供接口对kubernetes相关资源的管理

## 完成功能
1. pod容器列表
2. pod websocket
3. deployment 列表
4. namespace  列表
5. service 列表

# 运行 
git clone https://github.com/Michael754267513/k8s.git

修改配置文件 admin.conf 为kubernetes的config文件，一般是 ~/.kube/config文件

运行 go build main.go # 自动下包建议设置一个代理
 export GO111MODULE=on GOPROXY=https://goproxy.cn;  


# 接口
  SERVER  | DOMAIN  | ADDRESS | METHOD |            ROUTE            |                        HANDLER                         |      MIDDLEWARE
|---------|---------|---------|--------|-----------------------------|--------------------------------------------------------|-----------------------|
  default | default | :80     | ALL    | /kubernetes/deployment/list | k8s/kubernetes/deployment.(*DeployMentController).List | router.MiddlewareCORS
|---------|---------|---------|--------|-----------------------------|--------------------------------------------------------|-----------------------|
  default | default | :80     | ALL    | /kubernetes/namespace/list  | k8s/kubernetes/namespace.(*NameSpaceController).List   | router.MiddlewareCORS
|---------|---------|---------|--------|-----------------------------|--------------------------------------------------------|-----------------------|
  default | default | :80     | ALL    | /kubernetes/pod/list        | k8s/kubernetes/pod.(*podws.PodController).List         | router.MiddlewareCORS
|---------|---------|---------|--------|-----------------------------|--------------------------------------------------------|-----------------------|
  default | default | :80     | ALL    | /kubernetes/pod/websocket   | k8s/kubernetes/pod.(*podws.PodWSController).Websocket  | router.MiddlewareCORS
|---------|---------|---------|--------|-----------------------------|--------------------------------------------------------|-----------------------|
  default | default | :80     | ALL    | /kubernetes/service/list    | k8s/kubernetes/service.(*ServiceController).List       | router.MiddlewareCORS
|---------|---------|---------|--------|-----------------------------|--------------------------------------------------------|-----------------------|
