package handler

import (
	"encoding/json"
	"log"

	"github.com/rangaroo/2gis-friends/internal/database"
	"github.com/rangaroo/2gis-friends/internal/models"
	"github.com/rangaroo/2gis-friends/internal/state"
)

type BaseMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Handler struct {
	db        *database.Client
	store     *state.GlobalStore
}

func New(db *database.Client, store *state.GlobalStore) *Handler {
	return &Handler{
		db:        db,
		store:     store,
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
		h.handleInitialState(base.Payload)
	case "friendState":
		h.handleFriendState(base.Payload)
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

	// update profile data of global state
	h.store.UpdateFromPayload(data)

	for _, s := range data.States {
		if err := h.db.SaveState(s); err != nil {
			log.Printf("Failed to save state to db: %v\n", err)
		}
	}
}

func (h *Handler) handleFriendState(payload json.RawMessage) {
	var state models.State
	if err := json.Unmarshal(payload, &state); err != nil {
		log.Println("Error parsing friendState:", err)
		return
	}

	if err := h.db.SaveState(state); err != nil {
		log.Printf("Failed to save state to db: %v\n", err)
	}
}
