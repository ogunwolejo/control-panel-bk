package internal

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	wsClient = "websocket_client"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// User connection tracking (alternative to Redis for debugging)
var connections = make(map[string]*websocket.Conn)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	// Validate AWS TOKEN
	token := r.URL.Query().Get("token")

	if token == "" {
		log.Println("No token was sent")
		return
	}

	// TODO Validating Token against AWS Cognito
	header := w.Header()
	conn, e := upgrade.Upgrade(w, r, header)
	if e != nil {
		return
	}

	defer conn.Close()

	c := r.URL.Query().Get("role")
	userID := c

	// Store in local memory for debugging
	connections[userID] = conn

	// Keep the connection alive
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			removeWsConnection(userID)
			delete(connections, userID)
			break
		}

		log.Printf("Received from %s: %s\n", userID, msg)
	}
}

func storeWsConnection(userId, role string) error {
	userId = fmt.Sprintf("%s:%s", wsClient, userId)
	if err := RedisClient.HSet(context.TODO(), userId, map[string]interface{}{
		"role": role,
	}).Err(); err != nil {
		return err
	}

	return nil
}

func removeWsConnection(userId string) error {
	userId = fmt.Sprintf("%s:%s", wsClient, userId)
	if err := RedisClient.HDel(context.TODO(), userId).Err(); err != nil {
		return err
	}

	return nil
}
