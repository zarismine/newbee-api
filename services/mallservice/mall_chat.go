package mallservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"newbee/global"
	"newbee/models/jsontime"
	"newbee/models/mall"
	"newbee/pkg/dates"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// "newbee/global"
// "newbee/models/mall/response"
// "newbee/models/manage"

var ChatService = newChatService()

func newChatService() *chatService {
	return &chatService{}
}

type chatService struct {
}

type Node struct {
	UserId        int
	Conn          *websocket.Conn //连接
	Addr          string          //客户端地址
	HeartbeatTime int64           //心跳时间
	LoginTime     int64           //登录时间
	DataQueue     chan []byte     //消息
}

func (node *Node) Heartbeat(heartbeat int64) *Node {
	node.HeartbeatTime = heartbeat
	return node
}

var clientMap map[int]*Node = make(map[int]*Node, 0)
var rwLocker sync.RWMutex

func (c *chatService) Chat(writer http.ResponseWriter, request *http.Request) {
	protocols := websocket.Subprotocols(request)
	if len(protocols) == 0 || protocols[0] == "" {
		http.Error(writer, "Missing token", http.StatusBadRequest)
		return
	}
	token := protocols[0]
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		fmt.Println(err)
	}
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	upgrader.Subprotocols = protocols
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
	}
	node := &Node{
		UserId:        userToken.UserId,
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(),
		HeartbeatTime: dates.NowTimestamp(),
		LoginTime:     dates.NowTimestamp(),
		DataQueue:     make(chan []byte, 500),
	}
	rwLocker.Lock()
	clientMap[userToken.UserId] = node
	rwLocker.Unlock()
	go sendProc(node)
	go recvProc(node)
}

func sendProc(node *Node) {
	for {
		data := <-node.DataQueue
		fmt.Println("[ws]sendProc >>>> msg :", string(data))
		err := node.Conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := mall.MallMessage{}
		msg.SendId = node.UserId
		msg.CreateTime = jsontime.JSONTime{Time: time.Now()}
		err = json.Unmarshal(data, &msg)
		if msg.RecvId == 0 || msg.Content == "" {
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		data, err = json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
		}
		if msg.Type == -1 {
			currentTime := dates.NowTimestamp()
			node.Heartbeat(currentTime)
		} else {
			node.DataQueue <- data
			sendMsg(msg.RecvId, data)
			global.DB.Save(&msg)
		}

	}
}

func sendMsg(recvId int, data []byte) {
	rwLocker.RLock()
	nodeRecv, ok := clientMap[recvId]
	rwLocker.RUnlock()
	if ok {
		nodeRecv.DataQueue <- data
	}
}

func (c *chatService) GetRecord(id int, token string) ([]mall.MallMessage, string, error) {
	user, err := MallUserService.GetUserByToken(token)
	if err != nil {
		return nil, "", err
	}
	var messages []mall.MallMessage
	err = global.DB.Table("tb_newbee_mall_message").Where("(send_id = ? AND recv_id = ?) OR (send_id = ? AND recv_id = ?)",
		user.UserId, id, id, user.UserId).Find(&messages).Error
	if err != nil {
		return nil, "", err
	}
	return messages, user.NickName, nil
}
