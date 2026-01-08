package client

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rangaroo/2gis-friends/internal/config"
)

type Client struct {
	conn *websocket.Conn
}

func Connect(cfg *config.Config) (*Client, error) {
	headers := http.Header{}
	headers.Add("Origin", cfg.SiteDomain)
	headers.Add("User-Agent", cfg.UserAgent)

	// Connect to websocket
	log.Printf("Connecting to 2GIS...")
	conn, _, err := websocket.DefaultDialer.Dial(cfg.WebSocketURL(), headers)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}

	log.Println("Connected. Waiting for friends...")
	return &Client{conn: conn}, nil
}

func (c *Client) ReadMessages(handler func([]byte)) error {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}
		handler(message)
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}
