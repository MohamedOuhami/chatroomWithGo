package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

  // Initializing the reader
  reader := bufio.NewReader(os.Stdin)

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8000")

	if err != nil {
		fmt.Println("There was an error", err)
		return
	}

	for {

		fmt.Print("Enter your message : ")
    message,_ := reader.ReadString('\n')

		data := []byte(message)
		_, err = conn.Write(data)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

	}
	defer conn.Close()
}
