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

var (
    // Map to store client connections with usernames as keys
    clients = make(map[string]net.Conn)
    // Map to store username for each connection
    usernames = make(map[net.Conn]string)
    // Mutex for concurrency
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

func broadcastMessage(message string, username string) {
    mutex.RLock()
    defer mutex.RUnlock()
    
    messageWithTimestamp := message + " ---- " + time.Now().Format(time.RFC1123) + "\n"
    
    for _, conn := range clients {
        _, err := conn.Write([]byte(message + "\n"))
        if err != nil {
            println("There was an error broadcasting the message:", err)
        }
    }
    
    appendToLog(messageWithTimestamp)
}

func broadcastConnectedUsers() {
    mutex.RLock()
    defer mutex.RUnlock()
    
    var userList strings.Builder
    userList.WriteString("UserList:\n")
    userList.WriteString("-----------------\n")
    userList.WriteString(fmt.Sprintf("Total Users: %d\n", len(clients)))
    userList.WriteString("-----------------\n")
    
    for username := range clients {
        userList.WriteString(fmt.Sprintf("%s\n", username))
    }
    
    message := userList.String()
    
    for _, conn := range clients {
        _, err := conn.Write([]byte(message))
        if err != nil {
            println("Error sending user list:", err)
        }
    }
}

func handleClient(conn net.Conn) {
    // First message from client should be the username
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Error reading username:", err)
        conn.Close()
        return
    }
    
    username := strings.TrimSpace(string(buffer[:n]))
    
    // Lock for client registration
    mutex.Lock()
    clients[username] = conn
    usernames[conn] = username
    mutex.Unlock()
    
    // Announce new connection
    connectionMsg := fmt.Sprintf("%s has joined the chat", username)
    broadcastMessage(connectionMsg, username)
    
    // Broadcast updated user list after short delay
    time.Sleep(100 * time.Millisecond)
    broadcastConnectedUsers()
    
    // Cleanup on disconnect
    defer func() {
        mutex.Lock()
        delete(clients, username)
        delete(usernames, conn)
        mutex.Unlock()
        
        disconnectMsg := fmt.Sprintf("%s has left the chat", username)
        broadcastMessage(disconnectMsg, username)
        broadcastConnectedUsers()
        conn.Close()
    }()
    
    messageBuffer := make([]byte, 1024)
    for {
        n, err := conn.Read(messageBuffer)
        if err != nil {
            fmt.Printf("%s disconnected: %v\n", username, err)
            return
        }
        
        message := string(messageBuffer[:n])
        message = strings.TrimSpace(message)
        formattedMsg := fmt.Sprintf("%s: %s", username, message)
        broadcastMessage(formattedMsg, username)
    }
}

