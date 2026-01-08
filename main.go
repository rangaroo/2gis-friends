package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rangaroo/2gis-friends/internal/config"
	"github.com/rangaroo/2gis-friends/internal/database"

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
	db, err := database.NewClient(cfg.DBpath)
	if err != nil {
		log.Fatal("Failed to init database:", err)
	}
	defer db.Close()

	fmt.Println("Database connected successfully")


}
