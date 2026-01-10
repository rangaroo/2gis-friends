package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rangaroo/2gis-friends/internal/client"
	"github.com/rangaroo/2gis-friends/internal/config"
	"github.com/rangaroo/2gis-friends/internal/database"
	"github.com/rangaroo/2gis-friends/internal/handler"
)

func main() {
	log.Println("Starting 2GIS Friend Tracker...")

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
	}()
	log.Println("Database connected successfully")

	// create user cache to store 2GIS friend's profiles
	userCache := config.NewUserCache()

	// initialize handler
	h := handler.New(db, userCache)

	// create context that cancels when Ctrl+C is pressed
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// start supervisor loop that cancels on Ctrl+C and reconnects to websocket if connection is lost
	supervisorLoop(ctx, cfg, db, h)

	log.Println("Exiting...")
}

func supervisorLoop(ctx context.Context, cfg *config.Config, db *database.Client, h *handler.Handler) {
	timeout := 1 * time.Second

	for {
		// check if user pressed Ctrl+C
		if ctx.Err() != nil {
			return
		}

		err := runTracker(ctx, cfg, db, h.HandleMessage)

		if ctx.Err() != nil {
			return
		}

		log.Printf("Connection lost: %v. Retrying in %v...", err, timeout)

		select {
		case <-ctx.Done():
			return
		case <-time.After(timeout):
			timeout *= 2
			if timeout > 10*time.Second {
				timeout = 10 * time.Second
			}
		}
	}
}

func runTracker(ctx context.Context, cfg *config.Config, db *database.Client, handleMessage func([]byte)) error {
	// connect to client
	ws, err := client.Connect(cfg)
	if err != nil {
		return err
	}
	defer ws.Close()

	done := make(chan error, 1)

	go func() {
		err := ws.ReadMessages(handleMessage)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-done:
		return err
	}
}
