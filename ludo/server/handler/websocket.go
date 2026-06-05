package handler

import (
	"encoding/json"
	"log"
	"ludo-server/game"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

type WSClient struct {
	conn *websocket.Conn
	room *game.Room
	id   string
	name string
	send chan []byte
}

func (c *WSClient) Send(msg []byte) {
	c.send <- msg
}

func (c *WSClient) ID() string {
	return c.id
}

func (c *WSClient) Name() string {
	return c.name
}

func (c *WSClient) readPump() {
	defer func() {
		if c.room != nil {
			c.room.ActionChan() <- game.ActionMessage{Msg: game.Message{Type: game.MsgLeaveRoom}, Client: c}
		}
		c.conn.Close()
	}()

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg game.Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Println("Invalid msg:", err)
			continue
		}

		if msg.Type == game.MsgCreateRoom {
			var payload game.PayloadCreateRoom
			b, _ := json.Marshal(msg.Payload)
			json.Unmarshal(b, &payload)

			// Generate room ID
			roomID := "ROOM_" + time.Now().Format("150405")
			c.id = "p_" + time.Now().Format("150405.000")
			c.name = payload.PlayerName

			log.Println("Creating room:", roomID)
			room := game.GlobalHub.CreateRoom(roomID, payload.GameMode, payload.MaxPlayers, payload.AICount)
			c.room = room
			log.Println("Sending to RegisterChan...")
			room.RegisterChan() <- c
			log.Println("Sent to RegisterChan")

		} else if msg.Type == game.MsgJoinRoom {
			var payload game.PayloadJoinRoom
			b, _ := json.Marshal(msg.Payload)
			json.Unmarshal(b, &payload)

			room, err := game.GlobalHub.GetRoom(payload.RoomID)
			if err != nil {
				c.Send([]byte(`{"type":"error","payload":{"message":"Room not found"}}`))
				continue
			}

			c.id = "p_" + time.Now().Format("150405.000")
			c.name = payload.PlayerName
			c.room = room
			room.RegisterChan() <- c

		} else if c.room != nil {
			c.room.ActionChan() <- game.ActionMessage{
				Client: c,
				Msg:    msg,
			}
		}
	}
}

func (c *WSClient) writePump() {
	defer c.conn.Close()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &WSClient{
		conn: conn,
		send: make(chan []byte, 256),
	}
	go client.readPump()
	go client.writePump()
}
