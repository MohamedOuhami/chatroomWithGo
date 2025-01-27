package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Starting the id from 1
var id int = 1

func main() {

	// Listening for the incoming connections
	listener, err := net.Listen("tcp", "localhost:8000")

	// Checking if there was an error
	if err != nil {
		fmt.Println("There was an error ", err)
		return
	}

	// Create the serverLog
	createLogFile()

	// Make It the last thing to execute = defer
	defer listener.Close()

	println("Server is listening on port 8000")

	for {
		// Accepting the incoming connections
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("There was an error", err)
			continue
		}

		// Handle the client connection
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {

	connectionMessage := "Client " + strconv.Itoa(id) + " is now connected ----- " + time.Now().Format(time.RFC1123) + "\n"

	// Logging the entry in the serverLog
	file, Openerr := os.OpenFile("../logs/serverLog.txt", os.O_APPEND|os.O_WRONLY, 0666)

	if Openerr != nil {
		log.Fatal("There was an error", Openerr)
	}

	_, fileErr := file.Write([]byte(connectionMessage))

	if fileErr != nil {
		log.Fatal("There was an error", fileErr)
	}

	defer file.Close()
	userId := id
	id++

	defer conn.Close()
	// Create a buffer to put the received data in It ( It is a slice that we will change its content)
	buffer := make([]byte, 1024) // it is of length 1024

	for {
		// Read the data coming from the client
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Print the received data

		formatedMessage := "Client " + strconv.Itoa(userId) + " " + string(buffer[:n])
		fmt.Print(formatedMessage)

		// Logging the entry in the serverLog
		file, Openerr := os.OpenFile("../logs/serverLog.txt", os.O_APPEND|os.O_WRONLY, 0666)

		if Openerr != nil {
			log.Fatal("There was an error", Openerr)
		}

		formatedMessage = strings.TrimSpace(formatedMessage)
		_, fileErr := file.Write([]byte(formatedMessage + " ---- " + time.Now().Format(time.RFC1123) + "\n"))

		if fileErr != nil {
			log.Fatal("There was an error", fileErr)
		}

		defer file.Close()
	}

}

func createLogFile() {
	_, err := os.Create("../logs/serverLog.txt")
	if err != nil {
		log.Fatal("There was an error", err)
	}

}
