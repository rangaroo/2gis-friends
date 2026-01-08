package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rangaroo/2gis-friends/internal/models"
	"github.com/rangaroo/2gis-friends/internal/database"
)

type BaseMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Handler {
	db database.Client
}

func New(db *database.Client) {
	return &Handler{
		db: db,
	}
}

func (h *Handler) HandleMessage(data []byte) {
	var base BaseMessage
	if err := json.Unmarshal(data, &base); err != nil {
		log.Println("Failed to parse base message:", err)
		return
	}

	switch base.Type {
	case "initialState":

	case "friendState":

	default:
		// Ignore heartbeats or other messages
	}
}

func (h *Handler) handleInitialState(payload json.RawMessage) {
	var data models.InitialStatePayload
	if err := json.Unmarshal(payload, &data); err != nil {
		log.Println("Error parsing initialState:", err)
		return
	}

	fmt.Println("--- Loading Friends List ---")
	for _, p := range data.Profiles {
		userCache[p.ID] = p.Name
		fmt.Printf("Found friend: %s (%s)\n", p.Name, p.ID)
	}

	for _, s := range data.States {
		h.logState(s)
	}
}

func (h *Handler) handleFriendState(payload json.RawMessage) {
	var state models.State
	if err := json.Unmarshal(payload, &state); err != nil {
		log.Println("Error parsing friendState:", err)
		return
	}

	h.logState(state)
}

func (h *Handler) logState(s models.State) {
	name, exists := userCache[s.ID]
	if !exists {
		name = "Unknown User"
	}

	t := time.Unix(s.LastSeen/1000, 0)

	fmt.Printf("[UPDATE] %s is at [%f, %f] (Battery: %.0f%%) - Time: %s\n",
		name, s.Location.Lat, s.Location.Lon, s.Battery.Level*100, t.Format("15:04:05"))

	if err != h.db.SaveState(s); err != nil {
		log.Printf("Failed to save state to db: %v\n", err)
	}
}
