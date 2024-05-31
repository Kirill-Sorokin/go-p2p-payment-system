package main

import (
	"log"
	"math/rand"
	"sync"
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
	go startServer()

	// Start transactions between users
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
}
