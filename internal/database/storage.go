package database

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
	"github.com/rangaroo/2gis-friends/internal/models"
)

type Client struct {
	db *sql.DB
}

func NewClient(filepath string) (*Client, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	c := &Client{db: db}
	if err := c.createTables(); err != nil {
		return nil, err
	}

	return c, nil
}

func (db *Client) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS locations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		lat REAL NOT NULL,
		lon REAL NOT NULL,
		accuracy REAL,
		speed REAL,
		battery_level REAL,
		is_charging BOOLEAN,
		timestamp INTEGER NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_user_timestamp ON locations(user_id, timestamp DESC);
	`
	_, err := db.conn.Exec(schema)
	return err
}

func (db *Client) SaveState(state models.State) error {
	query := `INSERT INTO locations (user_id, lat, lon, accuracy, speed, battery_level, is_charging, timestamp)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.conn.Exec(query,
		state.ID,
		state.Location.Lat,
		state.Location.Lon,
		state.Location.Accuracy,
		state.Location.Speed,
		state.Battery.Level,
		state.Battery.IsCharging,
		state.LastSeen,
	)
	return err
}

func (c *Client) Reset() error {
	if _, err := c.db.Exec("DELETE FROM locations"); err != nil {
		return fmt.Errorf("failed to reset table locations: %w", err)
	}
}

func (c *Client) Close() {
	c.db.Close()
}
