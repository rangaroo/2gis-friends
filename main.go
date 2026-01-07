package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rangaroo/2gis-friends/internal/config"

	"github.com/gorilla/websocket"
)

// UserCache is a simple memory store to map IDs to Names
var userCache = make(map[string]string)

func main() {
	fmt.Println("Starting 2GIS Friend Tracker...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	// Initialize database connection
	db, err := NewClient(cfg.DBpath)
	if err != nil {
		log.Fatal("Failed to init database:", err)
	}
	defer db.Close()

	fmt.Println("Database connected successfully")

	url := cfg.WebSocketURL()

	headers := http.Header{}
	headers.Add("Origin", cfg.SiteDomain)
	headers.Add("User-Agent", cfg.UserAgent)

	// Connect to websocket
	log.Printf("Connecting to 2GIS...")
	c, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer c.Close()
	log.Println("Connected. Waiting for friends...")

	// Listen Loop
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		handleMessage(message, db)
	}
}
