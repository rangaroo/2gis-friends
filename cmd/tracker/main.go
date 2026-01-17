package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

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

	// initialize UI state
	store := state.NewStore()

	// initialize handler
	h := handler.New(db, store)


	// start UI
	m := ui.NewModel(store, cfg, h)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Printf("Error starting UI: %v", err)
		os.Exit(1)
	}

	//log.Println("Exiting...")
}
