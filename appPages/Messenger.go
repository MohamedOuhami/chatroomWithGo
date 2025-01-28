package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/rivo/tview"
)

// Send messages to the users

func SendMessage(conn net.Conn, messageView *tview.TextView, messageForm *tview.Form, message string) {
	if conn == nil {
		println("There is no connection")
		return
	}
	if message == "" {
		println("There is no message")
		return
	}

	// Send the raw message without "You: " prefix to the server
	data := []byte(message)
	_, err := conn.Write(data)
	if err != nil {
		println("There was an error in sending the data ", err)
		return
	}

	// Clear the input field
	messageForm.GetFormItemByLabel("Message").(*tview.InputField).SetText("")
	// Don't update the messageView here, let the server broadcast handle it
}

// Add a function to receive messages
func ReceiveMessages(conn net.Conn, app *tview.Application, connectedUsersView *tview.TextView, messageView *tview.TextView) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}

		message := string(buffer[:n])

		app.QueueUpdateDraw(func() {
			if strings.Contains(message, "UserList:") {
				connectedUsersView.SetText(message)
			} else {
				currentText := messageView.GetText(false)
				messageView.SetText(currentText + message)
				messageView.ScrollToEnd()
			}
		})
	}
}
