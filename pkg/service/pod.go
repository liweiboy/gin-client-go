package service

import (
	"context"
	"encoding/json"
	"errors"
	"gin-client-go/gin-client-go/pkg/client"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog"
	"net/http"
	"sync"
)

func GetPods(namespace string) ([]v1.Pod, error) {
	clientset, err := client.GetClientset()
	if err != nil {
		return nil, err
	}
	list, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func DeletePods(namespace string, names []string) error {
	clientset, err := client.GetClientset()
	if err != nil {
		return err
	}
	for _, name := range names {
		err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// ----------------webssh部分 需要用到websocket协议-------------------------

type WsMessage struct {
	MessageType int
	Data        []byte
}

type WsConnection struct {
	wsSocket  *websocket.Conn
	inChan    chan *WsMessage
	outChan   chan *WsMessage
	mutex     sync.Mutex
	isClosed  bool
	closeChan chan byte
}

// WsClose close websocket
func (wsConn *WsConnection) WsClose() {
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	err := wsConn.wsSocket.Close()
	if err != nil {
		klog.Errorln(err)
		return
	}
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}

// wsReadLoop read msg into inChan from websocket
func (wsConn *WsConnection) wsReadLoop() {
	var (
		msgType int
		data    []byte
		msg     *WsMessage
		err     error
	)
	for {
		if msgType, data, err = wsConn.wsSocket.ReadMessage(); err != nil {
			goto ERROR
		}
		msg = &WsMessage{
			MessageType: msgType,
			Data:        data,
		}
		select {
		case wsConn.inChan <- msg:
		case <-wsConn.closeChan:
			goto CLOSED
		}
	}
ERROR:
	wsConn.WsClose()
CLOSED:
}

// wsWriteLoop write msg into websocket from outChan
func (wsConn *WsConnection) wsWriteLoop() {
	var (
		msg *WsMessage
		err error
	)
	for {
		select {
		case msg = <-wsConn.outChan:
			err = wsConn.wsSocket.WriteMessage(msg.MessageType, msg.Data)
			if err != nil {
				goto ERROR
			}
		case <-wsConn.closeChan:
			goto CLOSED
		}

	}
ERROR:
	wsConn.WsClose()
CLOSED:
}

func (wsConn *WsConnection) WsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &WsMessage{MessageType: messageType, Data: data}:
		return nil
	case <-wsConn.closeChan:
		return errors.New("websocket closed")
	}
}

func (wsConn *WsConnection) WsRead() (*WsMessage, error) {
	select {
	case msg := <-wsConn.inChan:
		return msg, nil
	case <-wsConn.closeChan:
		return nil, errors.New("websocket closed")
	}
}

// http 升级到 ws 的升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func InitWebsocket(w http.ResponseWriter, r *http.Request) (*WsConnection, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	wsConn := &WsConnection{
		wsSocket:  conn,
		inChan:    make(chan *WsMessage),
		outChan:   make(chan *WsMessage),
		mutex:     sync.Mutex{},
		isClosed:  false,
		closeChan: make(chan byte),
	}
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()
	return wsConn, nil
}

type streamHandler struct {
	wsConn      *WsConnection
	resizeEvent chan remotecommand.TerminalSize
}

type xtermMessage struct {
	MsgType string `json:"type"`
	Input   string `json:"input"`
	Rows    uint16 `json:"rows"`
	Cols    uint16 `json:"cols"`
}

// Write write msg into outChan
func (handler *streamHandler) Write(p []byte) (int, error) {
	size := len(p)
	copyData := make([]byte, size)
	copy(copyData, p)
	err := handler.wsConn.WsWrite(websocket.TextMessage, copyData)
	if err != nil {
		return size, err
	}
	return size, nil
}

// Read read msg from inChan
func (handler *streamHandler) Read(p []byte) (size int, err error) {
	var (
		xtermMsg xtermMessage
		msg      *WsMessage
	)
	if msg, err = handler.wsConn.WsRead(); err != nil {
		klog.Errorln(err)
		return
	}
	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		return
	}
	if xtermMsg.MsgType == "resize" {
		handler.resizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	} else if xtermMsg.MsgType == "input" {
		size = len(xtermMsg.Input)
		copy(p, xtermMsg.MsgType)
	}
	return
}

// Next get elem from resizeEvent channel
func (handler *streamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.resizeEvent
	return &ret
}

func WebSSH(namespaceName, podName, containerName, method string, resp http.ResponseWriter, req *http.Request) error {
	kubeconfig, err := client.GetKubeconfig()
	if err != nil {
		return err
	}
	clientset, err := client.GetClientset()
	if err != nil {
		return err
	}
	reqSSH := clientset.CoreV1().RESTClient().Post().Resource("pods").Namespace(namespaceName).Name(podName).SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   []string{method},
			Stderr:    true,
			Stdin:     true,
			Stdout:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	executor, err := remotecommand.NewSPDYExecutor(kubeconfig, "POST", reqSSH.URL())
	if err != nil {
		return err
	}
	wsConn, err := InitWebsocket(resp, req)
	if err != nil {
		return err
	}

	handler := &streamHandler{
		wsConn:      wsConn,
		resizeEvent: make(chan remotecommand.TerminalSize),
	}

	if err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		klog.Errorln(err)
		wsConn.WsClose()
		return err
	}
	return err
}
