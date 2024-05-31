package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup

	// Create users
	users := []*User{
		NewUser("User1"),
		NewUser("User2"),
		NewUser("User3"),
		NewUser("User4"),
		NewUser("User5"),
	}

	// Add initial balances
	for _, user := range users {
		user.Deposit(1000)
	}

	// Start web server
	srv := startServer()

	// Start transactions between users
	go func() {
		for {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sender := users[rand.Intn(len(users))]
				receiver := users[rand.Intn(len(users))]
				amount := rand.Intn(200) + 1

				if sender != receiver {
					success := sender.Send(receiver, amount)
					if success {
						transaction := Transaction{
							Sender:    sender.name,
							Receiver:  receiver.name,
							Amount:    amount,
							Timestamp: time.Now(),
						}
						broadcastTransaction(transaction)
					} else {
						log.Printf("Failed transaction: %s -> %s (%d)", sender.name, receiver.name, amount)
					}
				}
			}()
			time.Sleep(time.Millisecond * 500) // Adjust the interval as needed
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func startServer() *http.Server {
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server is running on http://localhost:8080")
	return srv
}
