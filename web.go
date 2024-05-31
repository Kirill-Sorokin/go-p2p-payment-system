package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Transaction)
var mu sync.Mutex

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	for {
		var msg Transaction
		err := ws.ReadJSON(&msg)
		if err != nil {
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		log.Println("Broadcasting transaction: ", msg)

		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("WebSocket error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func broadcastTransaction(transaction Transaction) {
	broadcast <- transaction
}

func formatTransaction(transaction Transaction) string {
	timestamp := transaction.Timestamp.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("Sender: %s, Receiver: %s, Amount: %d, Date: %s, Time: %s",
		transaction.Sender, transaction.Receiver, transaction.Amount, timestamp[:10], timestamp[11:])
}
