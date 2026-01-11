package client

import (
	//"log"
	"net/http"
	"fmt"

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
	//log.Println("Connecting to 2GIS...")
	conn, _, err := websocket.DefaultDialer.Dial(cfg.WebSocketURL(), headers)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	//log.Println("Connected to 2GIS\n")

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
	if c.conn == nil {
        return nil
    }

	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return c.conn.Close()
}
