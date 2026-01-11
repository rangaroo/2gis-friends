package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rangaroo/2gis-friends/internal/client"
	"github.com/rangaroo/2gis-friends/internal/config"
	"github.com/rangaroo/2gis-friends/internal/database"
	"github.com/rangaroo/2gis-friends/internal/handler"
	"github.com/rangaroo/2gis-friends/internal/ui"
	"github.com/rangaroo/2gis-friends/internal/state"
)

func main() {
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

	// initialize UI state
	store := state.NewStore()

	// initialize handler
	h := handler.New(db, userCache, store)

	// create context that cancels when Ctrl+C is pressed
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// start supervisor loop that cancels on Ctrl+C and reconnects to websocket if connection is lost
	go supervisorLoop(ctx, cfg, h, store)

	// start UI
	p := tea.NewProgram(ui.NewModel(store))
	if _, err := p.Run(); err != nil {
		log.Printf("Error starting UI: %v", err)
		os.Exit(1)
	}

	//log.Println("Exiting...")
}

func supervisorLoop(
	ctx context.Context,
	cfg *config.Config,
	h *handler.Handler,
	store *state.GlobalStore,
) {
	timeout := 1 * time.Second

	for {
		// check if user pressed Ctrl+C
		if ctx.Err() != nil {
			return
		}

		_ = runTracker(ctx, cfg, h.HandleMessage, store)

		if ctx.Err() != nil {
			return
		}

		// log.Printf("Connection lost: %v. Retrying in %v...", err, timeout)

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

func runTracker(
	ctx context.Context,
	cfg *config.Config,
	handleMessage func([]byte),
	store *state.GlobalStore,
) error {
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
