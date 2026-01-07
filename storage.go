package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	db *sql.DB
}

func NewClient(filepath string) (*Client, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	c := &Client{db: db}
	if err := c.createTables(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) createTables() error {
	// 1. Table for static info (Name, Avatar)
	queryProfiles := `
	CREATE TABLE IF NOT EXISTS profiles (
		user_id TEXT PRIMARY KEY,
		name TEXT,
		avatar_url TEXT,
		last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 2. Table for the history log
	queryHistory := `
	CREATE TABLE IF NOT EXISTS locations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT,
		lat REAL,
		lon REAL,
		battery_level INTEGER,
		is_charging BOOLEAN,
		status TEXT,
		recorded_at DATETIME,
		FOREIGN KEY(user_id) REFERENCES profiles(user_id)
	);`

	if _, err := c.db.Exec(queryProfiles); err != nil {
		return err
	}
	if _, err := c.db.Exec(queryHistory); err != nil {
		return err
	}

	return nil
}

func (c *Client) SaveProfile(id, name, avatar string) error {
	query := `
	INSERT OR REPLACE INTO profiles (user_id, name, avatar_url, last_updated) 
	VALUES (?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := c.db.Exec(query, id, name, avatar)
	return err
}

func (c *Client) SaveState(state State) error {
	query := `
	INSERT INTO locations (user_id, lat, lon, battery_level, is_charging, status, recorded_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	// Convert 2GIS timestamp (milliseconds) to Go Time
	t := time.Unix(state.LastSeen/1000, 0)

	batteryPct := int(state.Battery.Level * 100)

	// NOTE: Leave it empty for now
	status := "unknown"

	_, err := c.db.Exec(query,
		state.ID,
		state.Location.Lat,
		state.Location.Lon,
		batteryPct,
		state.Battery.IsCharging,
		status,
		t,
	)
	return err
}

func (c *Client) Close() {
	c.db.Close()
}
