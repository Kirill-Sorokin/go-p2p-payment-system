# Payment System Monitor

This project simulates a payment system and monitors transactions in real-time using Go, Gin, and WebSocket. It allows multiple users to perform transactions between each other and displays these transactions in real-time on a web interface.

## Features

- **Simulates transactions** between multiple users
- **Displays transactions in real-time** using WebSocket
- **User-friendly interface** to view transaction details
- **Concurrent transactions** using Go routines and mutexes

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or later
- [Git](https://git-scm.com/)

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/Kirill-Sorokin/payment-system-monitor.git
   cd payment-system-monitor

2. Initialise a new Go module (if not already initialised):
   ```sh
   go mod init payment-system-monitor

3. Install dependencies:
   ```sh
   go mod tidy


### Running the Application
1. Run the application:
   ```sh
   go run *.go

2. Open your browser and navigate to http://localhost:8080.

### Project Structure
- **main.go**: Entry point of the application. Initializes users and starts the web server.
- **user.go**: Contains the User struct and methods for deposit, withdrawal, and sending money.
- **transaction.go**: Defines the Transaction struct.
- **web.go**: Sets up the Gin web server and WebSocket connections.
- **templates/index.html**: HTML template for the web interface.
- **static/style.css**: CSS for the web interface.
