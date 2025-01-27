package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
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

func broadcastMessage(message string, senderId int){
  mutex.RLock()
  defer mutex.RUnlock()

  for id,conn := range clients {
    if id != senderId {
      _,err := conn.Write([]byte(message))
      if err != nil {
        println("There was an error broadcasting the message :",err)
      }
    }
    
  }
}

func broadcastConnectedUsers(){
  mutex.RLock()
  defer mutex.RUnlock()

  var allUserIds string
// Iterate over the clients map and collect all user IDs
	for id := range clients {
		allUserIds += "User " + strconv.Itoa(id) + "\n" // Concatenate each ID with a newline to separate them
	}

  for _,conn := range clients {

    _,err := conn.Write([]byte(allUserIds))

    if err != nil {
      println("There was an error sending the connected users ",err)
    }

  }
}

func handleClient(conn net.Conn) {

	mutex.Lock()
	userId := id
	clients[userId] = conn
	id++
	mutex.Unlock()

	connectionMessage := "Client " + strconv.Itoa(userId) + " is now connected" 
  connectionMessageTimeStamped := connectionMessage + " ---- " + time.Now().Format(time.RFC1123) + "\n"
	println(connectionMessageTimeStamped)
	appendToLog(connectionMessageTimeStamped)

  broadcastMessage(connectionMessage,userId)
  broadcastConnectedUsers()

	defer func() {

		mutex.Lock()
		delete(clients, userId)
		mutex.Unlock()
		conn.Close()

	}()

	// Create a buffer to put the received data in It ( It is a slice that we will change its content)
	buffer := make([]byte, 1024) // it is of length 1024

	for {
		// Read the data coming from the client
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("There was a connection error : ", err)
			return
		}

		// Print the received data

    formatedMessage := "Client " + strconv.Itoa(userId) + " : " + string(buffer[:n])
		formatedMessage = strings.TrimSpace(formatedMessage)

		appendToLog(formatedMessage + " ---- " + time.Now().Format(time.RFC1123) + "\n")

    broadcastMessage(formatedMessage,userId)


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
