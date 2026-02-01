package core

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketConn struct {
	conn *websocket.Conn
}

func ConnectToWebSocket(cfg Config) (*WebSocketConn, error) {
	headers := http.Header{}
	headers.Add("Origin", cfg.SiteDomain)
	headers.Add("User-Agent", cfg.UserAgent)

	// Connect to websocket
	conn, _, err := websocket.DefaultDialer.Dial(cfg.WebSocketURL(), headers)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	//log.Println("Connected to 2GIS\n")

	return &WebSocketConn{conn: conn}, nil
}

func (c *WebSocketConn) ReadMessages() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	return message, err
}

func (c *WebSocketConn) Close() error {
	if c.conn == nil {
		return nil
	}

	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return c.conn.Close()
}
