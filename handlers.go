package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type BaseMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LocationData struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Accuracy float64 `json:"accuracy"`
	Speed    float64 `json:"speed"`
}

type State struct {
	ID        string       `json:"id"`
	LastSeen  int64        `json:"lastSeen"`
	Location  LocationData `json:"location"`
	Battery   struct {
		Level      float64 `json:"level"`
		IsCharging bool    `json:"isCharging"`
	} `json:"battery"`
}

type InitialStatePayload struct {
	Profiles []Profile `json:"profiles"`
	States   []State   `json:"states"`
}

func handleMessage(data []byte) {
	var base BaseMessage
	if err := json.Unmarshal(data, &base); err != nil {
		log.Println("Failed to parse base message:", err)
		return
	}

	switch base.Type {
	case "initialState":
		var payload InitialStatePayload
		if err := json.Unmarshal(base.Payload, &payload); err != nil {
			log.Println("Error parsing initialState:", err)
			return
		}
		
		// Fill the phonebook
		fmt.Println("--- Loading Friends List ---")
		for _, p := range payload.Profiles {
			userCache[p.ID] = p.Name
			fmt.Printf("Found friend: %s (%s)\n", p.Name, p.ID)
		}
		
		// Process initial locations
		for _, s := range payload.States {
			logState(s)
		}

	case "friendState":
		var state State
		// note: friendState payload is a single State object, not a list
		if err := json.Unmarshal(base.Payload, &state); err != nil {
			log.Println("Error parsing friendState:", err)
			return
		}
		logState(state)

	default:
		// Ignore heartbeats or other messages
		// fmt.Printf("Unknown message type: %s\n", base.Type)
	}
}

func logState(s State) {
	name, exists := userCache[s.ID]
	if !exists {
		name = "Unknown User"
	}

	// Convert timestamp from ms to Time
	t := time.Unix(s.LastSeen/1000, 0)
    
	fmt.Printf("[UPDATE] %s is at [%f, %f] (Battery: %.0f%%) - Time: %s\n", 
        name, s.Location.Lat, s.Location.Lon, s.Battery.Level*100, t.Format("15:04:05"))
}
