package main

import (
	"fmt"
	"log"

	"github.com/rangaroo/2gis-friends/internal/client"
	"github.com/rangaroo/2gis-friends/internal/config"
	"github.com/rangaroo/2gis-friends/internal/database"
	"github.com/rangaroo/2gis-friends/internal/handler"
)

func main() {
	fmt.Println("Starting 2GIS Friend Tracker...")

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	// initialize database
	db, err := database.NewClient(cfg.DBpath)
	if err != nil {
		log.Fatal("Failed to init database:", err)
	}
	defer db.Close()
	fmt.Println("Database connected successfully")

	// create user cache to store 2GIS friend's profiles
	userCache := config.NewUserCache()

	// initialize handler (or whatever the correct verb that comes with handlers)
	handler := handler.New(db, userCache)

	// connect to client
	ws, err := client.Connect(cfg)
	if err != nil {
		log.Fatal("WebSocket connection failed:", err)
	}
	defer ws.Close()

	if err := ws.ReadMessages(handler.HandleMessage); err != nil {
		log.Println("WebSocket error:", err)
	}

}
