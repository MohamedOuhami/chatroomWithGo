package main

import (
	"fmt"
	"net"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	App                = tview.NewApplication()
	Pages              = tview.NewPages()
	welcomeText        = tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("(q) to quit")
	Messageform        = tview.NewForm()
	messageToSend      string
	MessageView        *tview.TextView
	ConnectedUsersView *tview.TextView
	conn               net.Conn
)

func main() {
	var err error
	conn, err = net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("There was an error", err)
		return
	}
	defer conn.Close()

	// Initialize the views
	MessageView = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() { App.Draw() }).
		SetScrollable(true)

	ConnectedUsersView = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() { App.Draw() })

	// Now we can pass the pointers directly
	go ReceiveMessages(conn, App, ConnectedUsersView, MessageView)
	SetupUI(conn)
}
