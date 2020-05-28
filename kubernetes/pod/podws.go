package podws

import (
	"encoding/json"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	initConfig "k8s/config"
)

type PodWSController struct {
	gmvc.Controller
	ws *ghttp.WebSocket
}

var (
	clientset *kubernetes.Clientset
)

// ssh流式处理器
type streamHandler struct {
	wsConn      *ghttp.WebSocket
	resizeEvent chan remotecommand.TerminalSize
}

// web终端发来的包
type xtermMessage struct {
	MsgType string `json:"type"`  // 类型:resize客户端调整终端, input客户端输入
	Input   string `json:"input"` // msgtype=input情况下使用
	Rows    uint16 `json:"rows"`  // msgtype=resize情况下使用
	Cols    uint16 `json:"cols"`  // msgtype=resize情况下使用
}

// executor回调获取web是否resize
func (handler *streamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.resizeEvent
	size = &ret
	return
}

// executor回调读取web端的输入
func (handler *streamHandler) Read(p []byte) (size int, err error) {
	var (
		msg      []byte
		xtermMsg xtermMessage
	)

	// 读web发来的输入
	if _, msg, err = handler.wsConn.ReadMessage(); err != nil {
		return
	}

	// 解析客户端请求
	if err = json.Unmarshal(msg, &xtermMsg); err != nil {
		return
	}

	//web ssh调整了终端大小
	if xtermMsg.MsgType == "resize" {
		// 放到channel里，等remotecommand executor调用我们的Next取走
		handler.resizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	} else if xtermMsg.MsgType == "input" { // web ssh终端输入了字符
		// copy到p数组中
		size = len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return
}

// executor回调向web端输出
func (handler *streamHandler) Write(p []byte) (size int, err error) {
	var (
		copyData []byte
	)

	// 产生副本
	copyData = make([]byte, len(p))
	copy(copyData, p)
	size = len(p)
	err = handler.wsConn.WriteMessage(websocket.TextMessage, copyData)
	return
}

func (r *PodWSController) Websocket() {
	var (
		//wsConn *ws.WsConnection
		restConf      *rest.Config
		sshReq        *rest.Request
		podName       string
		podNs         string
		containerName string
		executor      remotecommand.Executor
		handler       *streamHandler
		err           error
	)

	// 解析GET参数
	podNs = r.Request.GetString("namespace")
	podName = r.Request.GetString("PodName")
	containerName = r.Request.GetString("ContainerName")

	// 得到websocket长连接
	if wsconn, err := r.Request.WebSocket(); err == nil {
		r.ws = wsconn
	} else {
		g.Log().Error(err)
		return
	}
	// 获取pods

	// 获取k8s rest client配置
	if restConf, err = initConfig.GetRestConf(); err != nil {
		goto END
	}
	clientset, _ = initConfig.InitClient()
	sshReq = clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(podNs).
		SubResource("exec")
	sshReq.VersionedParams(&v1.PodExecOptions{
		Container: containerName,
		Command:   []string{"bash"},
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	// 创建到容器的连接
	if executor, err = remotecommand.NewSPDYExecutor(restConf, "POST", sshReq.URL()); err != nil {
		goto END
	}

	// 配置与容器之间的数据流处理回调
	handler = &streamHandler{wsConn: r.ws, resizeEvent: make(chan remotecommand.TerminalSize)}
	if err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		goto END
	}
	return

END:
	initConfig.Logger(err)
	_ = r.ws.Close()
}
