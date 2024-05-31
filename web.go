package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
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

func startServer() {
	r := gin.Default()

	// Set trusted proxies
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/ws", func(c *gin.Context) {
		handleConnections(c.Writer, c.Request)
	})

	go handleMessages()

	r.Run(":8080")
}

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
		formattedMessage := formatTransaction(msg)
		log.Println(formattedMessage)

		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(formattedMessage)
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
	return "Sender: " + transaction.Sender +
		", Receiver: " + transaction.Receiver +
		", Amount: " + fmt.Sprintf("%d", transaction.Amount) +
		", Date: " + timestamp[:10] +
		", Time: " + timestamp[11:]
}
