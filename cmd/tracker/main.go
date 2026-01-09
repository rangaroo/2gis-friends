package main

import (
	"fmt"
	"log"
	"context"
	"signal"

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
	defer func() {
		log.Println("Closing database...")
		db.Close()
	}
	fmt.Println("Database connected successfully")

	// create user cache to store 2GIS friend's profiles
	userCache := config.NewUserCache()

	// initialize handler
	handler := handler.New(db, userCache)

	// create context that cancels when Ctrl+C is pressed
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

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

func supervisorLoop(ctx context.Context, cfg *config.Config, db *database.Client) {
	timeout := 1 * time.Second

	for {
		// check if user pressed Ctrl+C
		if ctx.Err() != nil {
			return
		}

		err != runTracker(ctx, cfg, db)

		if ctx.Err() != nil {
			return
		}

		log.Printf("Connection lost: %v. Retrying in %v...", err, timeout)

		select {
		case <-ctx.Done():
			return
		case <-time.After(timeout):
			timeout *= 2
			if timeout > 10 * time.Second {
				timeout = 10 * time.Second
			}
		}
	}
}

func runTracker(ctx context.Context, cfg *config.Config, db *database.Client) error {
	// connect to client
	ws, err := client.Connect(cfg)
	if err != nil {
		return err
	}
	defer ws.Close()

	done := make(chan error, 1)

	go func() {
		err := ws.ReadMessages(handler.HandleMessage)
		done <- err
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-done:
		return err
	}
}
