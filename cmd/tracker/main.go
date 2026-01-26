package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rangaroo/2gis-friends/internal/core"
	"github.com/rangaroo/2gis-friends/internal/ui"
)

func main() {
	// load config
	cfg, err := core.Load()
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	// initialize database
	db, err := core.NewDatabaseClient(cfg.DBpath)
	if err != nil {
		log.Fatal("Failed to init database:", err)
	}
	defer func() {
		//log.Println("Closing database...")
		db.Close()
	}()

	// start UI
	m := ui.NewModel(cfg, db)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Printf("Error starting UI: %v", err)
		os.Exit(1)
	}

	//log.Println("Exiting...")
}
