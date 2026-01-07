package main

import (
	"os"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/gorilla/websocket"
)

// UserCache is a simple memory store to map IDs to Names
var userCache = make(map[string]string)

type appConfig struct {
	accessToken string
	appVersion  string
	userAgent   string
	siteDomain  string
}

func main() {
	godotenv.Load(".env")

	accessToken := os.Getenv("ACCESS_TOKEN") //TODO: Figure out how to generate these tokens
	if accessToken == "" {
		log.Fatal("ACCESS_TOKEN must be set")
	}

	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		log.Fatal("APP_VERSION must be set")
	}

	userAgent := os.Getenv("USER_AGENT")
	if userAgent == "" {
		log.Fatal("USER_AGENT must be set")
	}

	siteDomain := os.Getenv("SITE_DOMAIN")
	if siteDomain == "" {
		log.Fatal("SITE_DOMAIN must be set")
	}

	cfg := appConfig{
		accessToken: accessToken,
		appVersion:  appVersion,
		userAgent:   userAgent,
		siteDomain:  siteDomain,
	}

	url := fmt.Sprintf(
		"wss://zond.api.2gis.ru/api/1.1/user/ws?appVersion=%s&channels=markers,sharing,routes&token=%s",
		cfg.appVersion,
		cfg.accessToken,
	)

	headers := http.Header{}
	headers.Add("Origin", cfg.siteDomain)
	headers.Add("User-Agent", cfg.userAgent)

	// Connect to websocket
	log.Printf("Connecting to 2GIS...")
	c, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer c.Close()
	log.Println("Connected. Waiting for friends...")

	// Listen Loop
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		handleMessage(message)
	}
}

