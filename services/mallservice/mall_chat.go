package mallservice

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"newbee/global"
	"newbee/models/jsontime"
	"newbee/models/mall"
	"newbee/models/mall/response"
	"newbee/pkg/dates"
	"strconv"
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

var clientMap = make(map[int]*Node, 0)
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

func UpDateCache(key string, field string, val string, isRecv bool) {
	data, err := global.Redis.HGet(context.Background(), key, field).Result()
	if err != nil {
		return
	}
	contact := new(response.ContactResponse)
	_ = jsoniter.Unmarshal([]byte(data), contact)
	if isRecv && val != "" {
		contact.Count += 1
		contact.MessageContent = val
		contact.MessageTime = jsontime.JSONTime{
			Time: time.Now(),
		}.Format("2006-01-02 15:04:05")
	}
	if isRecv && val == "" {
		contact.Count = 0
	}
	if !isRecv {
		contact.MessageContent = val
		contact.MessageTime = jsontime.JSONTime{
			Time: time.Now(),
		}.Format("2006-01-02 15:04:05")
	}
	resp, _ := jsoniter.Marshal(contact)
	fmt.Println(string(resp))
	global.Redis.HSet(context.Background(), key, field, string(resp))
}

func sendProc(node *Node) {
	defer func() {
		if _, ok := clientMap[node.UserId]; ok {
			delete(clientMap, node.UserId)
		}
	}()
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
	defer func() {
		if _, ok := clientMap[node.UserId]; ok {
			delete(clientMap, node.UserId)
		}
	}()
	cacheUserContact := fmt.Sprintf("%s%v", global.CacheUserContactPrefix, node.UserId)
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			return
		}
		msg := mall.MallMessage{
			SendId:        node.UserId,
			MessageStatus: 1,
			CreateTime:    jsontime.JSONTime{Time: time.Now()},
		}
		err = jsoniter.Unmarshal(data, &msg)
		if msg.RecvId == 0 || msg.Content == "" {
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		if msg.Type == -1 {
			currentTime := dates.NowTimestamp()
			node.Heartbeat(currentTime)
		} else {
			global.DB.Create(&msg)
			data, _ = jsoniter.Marshal(msg)
			node.DataQueue <- data
			UpDateCache(cacheUserContact, strconv.Itoa(msg.RecvId), msg.Content, false)
			cacheUserContactRecv := fmt.Sprintf("%s%v", global.CacheUserContactPrefix, msg.RecvId)
			UpDateCache(cacheUserContactRecv, strconv.Itoa(node.UserId), msg.Content, true)
			sendMsg(msg.RecvId, data)
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
	global.DB.Table("tb_newbee_mall_message").Where("send_id = ? AND recv_id = ?", id, user.UserId).Update("message_status", 0)
	cacheUserContactRecv := fmt.Sprintf("%s%v", global.CacheUserContactPrefix, user.UserId)
	UpDateCache(cacheUserContactRecv, strconv.Itoa(id), "", true)
	if err != nil {
		return nil, "", err
	}
	return messages, user.NickName, nil
}
