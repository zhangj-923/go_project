package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"flying-server/game"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client WebSocket客户端
type Client struct {
	ID   string
	Conn *websocket.Conn
	Room *game.Room
	Send chan []byte
}

// Hub WebSocket中心
type Hub struct {
	mu         sync.RWMutex
	Clients    map[string]*Client
	Rooms      map[string]*game.Room
	Register   chan *Client
	Unregister chan *Client
}

// NewHub 创建Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[string]*game.Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("Client connected: %s", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)

				// 从房间移除
				if client.Room != nil {
					client.Room.RemovePlayer(client.ID)
				}
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: %s", client.ID)
		}
	}
}

// HandleWebSocket 处理WebSocket连接
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	// 生成唯一客户端ID
	clientID := fmt.Sprintf("player_%d", time.Now().UnixNano())

	client := &Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	h.Register <- client

	// 发送客户端ID
	h.sendToClient(client, game.GameMessage{
		Type: "connected",
		Data: map[string]string{"id": clientID},
	})

	go h.writePump(client)
	go h.readPump(client)
}

// readPump 读取消息
func (h *Hub) readPump(client *Client) {
	defer func() {
		h.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var msg game.GameMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Unmarshal error: %v", err)
			continue
		}

		h.handleMessage(client, msg)
	}
}

// writePump 写入消息
func (h *Hub) writePump(client *Client) {
	defer client.Conn.Close()

	for message := range client.Send {
		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

// handleMessage 处理消息
func (h *Hub) handleMessage(client *Client, msg game.GameMessage) {
	switch msg.Type {
	case "create_room":
		h.handleCreateRoom(client)
	case "join_room":
		h.handleJoinRoom(client, msg)
	case "start_game":
		h.handleStartGame(client)
	case "bid":
		h.handleBid(client, msg)
	case "play":
		h.handlePlay(client, msg)
	case "pass":
		h.handlePass(client)
	case "get_hint":
		h.handleGetHint(client)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// handleCreateRoom 创建房间
func (h *Hub) handleCreateRoom(client *Client) {
	roomID := fmt.Sprintf("room_%d", time.Now().UnixNano())
	room := game.NewRoom(roomID)
	room.OnMessage = func(playerID string, msg game.GameMessage) {
		h.mu.RLock()
		c, ok := h.Clients[playerID]
		h.mu.RUnlock()
		if ok {
			h.sendToClient(c, msg)
		}
	}

	// 创建者加入房间
	player := game.NewHumanPlayer(client.ID, "玩家")
	room.AddPlayer(player)

	h.mu.Lock()
	h.Rooms[roomID] = room
	client.Room = room
	h.mu.Unlock()

	h.sendToClient(client, game.GameMessage{
		Type: "room_created",
		Data: map[string]string{"room_id": roomID},
	})

	// 发送房间信息
	h.sendToClient(client, game.GameMessage{
		Type: "room_info",
		Data: room.GetRoomInfo(),
	})
}

// handleJoinRoom 加入房间
func (h *Hub) handleJoinRoom(client *Client, msg game.GameMessage) {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return
	}

	roomID, ok := data["room_id"].(string)
	if !ok {
		return
	}

	h.mu.RLock()
	room, exists := h.Rooms[roomID]
	h.mu.RUnlock()

	if !exists {
		h.sendToClient(client, game.GameMessage{
			Type: "error",
			Data: "房间不存在",
		})
		return
	}

	if room.IsFull() {
		h.sendToClient(client, game.GameMessage{
			Type: "error",
			Data: "房间已满",
		})
		return
	}

	player := game.NewHumanPlayer(client.ID, "玩家")
	room.AddPlayer(player)

	h.mu.Lock()
	client.Room = room
	h.mu.Unlock()

	h.sendToClient(client, game.GameMessage{
		Type: "room_joined",
		Data: map[string]string{"room_id": roomID},
	})

	// 发送房间信息
	h.sendToClient(client, game.GameMessage{
		Type: "room_info",
		Data: room.GetRoomInfo(),
	})

	// 如果房间满了，自动开始
	if room.IsFull() {
		h.autoFillAndStart(room)
	}
}

// autoFillAndStart 自动填充AI并开始
func (h *Hub) autoFillAndStart(room *game.Room) {
	// 开始游戏（会自动添加AI玩家）
	room.StartGame()
}

// handleStartGame 开始游戏
func (h *Hub) handleStartGame(client *Client) {
	if client.Room == nil {
		return
	}

	// 自动填充AI
	h.autoFillAndStart(client.Room)
}

// handleBid 处理叫地主
func (h *Hub) handleBid(client *Client, msg game.GameMessage) {
	if client.Room == nil {
		return
	}

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return
	}

	score, ok := data["score"].(float64)
	if !ok {
		return
	}

	err := client.Room.HandleBid(client.ID, int(score))
	if err != nil {
		h.sendToClient(client, game.GameMessage{
			Type: "error",
			Data: err.Error(),
		})
	}
}

// handlePlay 处理出牌
func (h *Hub) handlePlay(client *Client, msg game.GameMessage) {
	if client.Room == nil {
		return
	}

	// 解析出牌数据
	data, err := json.Marshal(msg.Data)
	if err != nil {
		return
	}

	var cards game.Cards
	if err := json.Unmarshal(data, &cards); err != nil {
		// 尝试解析为数组
		var cardArray []struct {
			Suit int `json:"suit"`
			Rank int `json:"rank"`
		}
		if err := json.Unmarshal(data, &cardArray); err != nil {
			return
		}
		for _, c := range cardArray {
			cards = append(cards, game.NewCard(game.Suit(c.Suit), game.Rank(c.Rank)))
		}
	}

	err = client.Room.HandlePlay(client.ID, cards)
	if err != nil {
		h.sendToClient(client, game.GameMessage{
			Type: "error",
			Data: err.Error(),
		})
	}
}

// handlePass 处理不出
func (h *Hub) handlePass(client *Client) {
	if client.Room == nil {
		return
	}

	err := client.Room.HandlePlay(client.ID, nil)
	if err != nil {
		h.sendToClient(client, game.GameMessage{
			Type: "error",
			Data: err.Error(),
		})
	}
}

// handleGetHint 处理提示请求
func (h *Hub) handleGetHint(client *Client) {
	if client.Room == nil {
		return
	}

	player := client.Room.GetPlayer(client.ID)
	if player == nil {
		return
	}

	// 获取提示
	hint := game.SuggestPlay(player.GetCards(), client.Room.LastCards, client.Room.LastResult)

	h.sendToClient(client, game.GameMessage{
		Type: "hint",
		Data: hint,
	})
}

// sendToClient 发送消息给客户端
func (h *Hub) sendToClient(client *Client, msg game.GameMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return
	}

	select {
	case client.Send <- data:
	default:
		close(client.Send)
		delete(h.Clients, client.ID)
	}
}
