package main

import (
	"ChatroomWithGo/utils"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Starting the id from 1
var (
	// First id
	id int = 1

	// Create a map of client connection
	clients = make(map[int]net.Conn)

	// Create a Mutex for concurrency
	mutex sync.RWMutex
)

func main() {

	utils.ConnectToDb()

	// Start the server
	listener := startServer()

	// Create the serverLog
	createLogFile()

	if listener == nil {
		return
	}

	// Make It the last thing to execute = defer
	defer listener.Close()

	println("Server is listening on port 8000")

	for {
		// Accepting the incoming connections
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("There was a connection error", err)
			continue
		}

		// Handle the client connection
		go handleClient(conn)
	}
}

func startServer() net.Listener {
    // Listening for the incoming connections
    listener, err := net.Listen("tcp", "localhost:8000")
    // Checking if there was an error
    if err != nil {
        fmt.Println("There was a server error ", err)
        return nil
    }
    return listener
}
func appendToLog(message string) {
    // Logging the entry in the serverLog
    file, Openerr := os.OpenFile("../logs/serverLog.txt", os.O_APPEND|os.O_WRONLY, 0666)
    if Openerr != nil {
        log.Fatal("There was an error opening the logFile", Openerr)
    }
    _, fileErr := file.Write([]byte(message))
    if fileErr != nil {
        log.Fatal("There was an error writing to Log File", fileErr)
    }
    defer file.Close()
}
func createLogFile() {
    _, err := os.Create("../logs/serverLog.txt")
    if err != nil {
        log.Fatal("There was a file creation error", err)
    }
}

func broadcastMessage(message string, senderId int) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Format the message to include timestamp
	messageWithTimestamp := message + " ---- " + time.Now().Format(time.RFC1123) + "\n"

	// Broadcast to all clients, including sender
	for _, conn := range clients {
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			println("There was an error broadcasting the message:", err)
		}
	}

	// Log the message
	appendToLog(messageWithTimestamp)
}

func broadcastConnectedUsers() {
	mutex.RLock()
	defer mutex.RUnlock()

	// Format the connected users list
	var userList strings.Builder
	userList.WriteString("UserList:\n")
	userList.WriteString("-----------------\n")
	userList.WriteString(fmt.Sprintf("Total Users: %d\n", len(clients)))
	userList.WriteString("-----------------\n")

	for clientId := range clients {
		userList.WriteString(fmt.Sprintf("User %d\n", clientId))
	}

	message := userList.String()

	// Send to all connected clients
	for _, conn := range clients {
		_, err := conn.Write([]byte(message))
		if err != nil {
			println("Error sending user list:", err)
		}
	}
}

func handleClient(conn net.Conn) {
	// Lock for client registration
	mutex.Lock()
	userId := id
	clients[userId] = conn
	id++
	mutex.Unlock()

	// Announce new connection
	connectionMsg := fmt.Sprintf("Client %d has joined the chat", userId)
	broadcastMessage(connectionMsg, userId)

	// Broadcast updated user list after short delay
	time.Sleep(100 * time.Millisecond)
	broadcastConnectedUsers()

	// Cleanup on disconnect
	defer func() {
		mutex.Lock()
		delete(clients, userId)
		mutex.Unlock()

		// Announce disconnection
		disconnectMsg := fmt.Sprintf("Client %d has left the chat", userId)
		broadcastMessage(disconnectMsg, userId)

		// Update user list
		broadcastConnectedUsers()

		conn.Close()
	}()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Client %d disconnected: %v\n", userId, err)
			return
		}

		message := string(buffer[:n])
		message = strings.TrimSpace(message)

		// Format and broadcast the message
		formattedMsg := fmt.Sprintf("Client %d: %s", userId, message)
		broadcastMessage(formattedMsg, userId)
	}
}
