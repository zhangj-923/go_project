package game

import (
	"errors"
	"sync"
)

type Hub struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

var GlobalHub = &Hub{
	rooms: make(map[string]*Room),
}

func (h *Hub) CreateRoom(id, mode string, maxP, aiCount int) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	room := NewRoom(id, mode, maxP, aiCount)
	h.rooms[id] = room
	go room.Run()
	return room
}

func (h *Hub) GetRoom(id string) (*Room, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	room, ok := h.rooms[id]
	if !ok {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (h *Hub) RemoveRoom(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms, id)
}
