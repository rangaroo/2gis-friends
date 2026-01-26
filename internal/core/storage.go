package core

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DatabaseClient struct {
	db *sql.DB
}

func NewDatabaseClient(filepath string) (*DatabaseClient, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	c := &DatabaseClient{db: db}
	if err := c.initSchema(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *DatabaseClient) initSchema() error {
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
	_, err := c.db.Exec(schema)
	return err
}

func (c *DatabaseClient) SaveState(state State) error {
	query := `INSERT INTO locations (user_id, lat, lon, accuracy, speed, battery_level, is_charging, timestamp)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := c.db.Exec(query,
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

func (c *DatabaseClient) Reset() error {
	if _, err := c.db.Exec("DELETE FROM locations"); err != nil {
		return fmt.Errorf("failed to reset table locations: %w", err)
	}

	return nil
}

func (c *DatabaseClient) Close() {
	c.db.Close()
}
