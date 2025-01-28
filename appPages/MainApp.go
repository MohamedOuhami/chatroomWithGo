package main

import (
	"github.com/rivo/tview"
	"net"
)

var ()
func DrawMainApp(conn net.Conn) *tview.Flex {
    // Message view setup
    MessageView.SetBorder(true).
        SetTitle("Chat Messages")

    // Users view setup
    ConnectedUsersView.SetBorder(true).
        SetTitle("Online Users")

    // Message form setup
    Messageform = tview.NewForm().
        AddInputField("Message", "", 0, nil, func(message string) {
            messageToSend = message
        }).
        AddButton("Send", func() {
            SendMessage(conn, MessageView, Messageform, messageToSend)
        }).
        AddButton("Quit", func() {
            App.Stop()
        })
    Messageform.SetBorder(true).
        SetTitle("Message Input")

    // Main content area (messages and users side by side)
    contentArea := tview.NewFlex().
        SetDirection(tview.FlexColumn).
        AddItem(MessageView, 0, 7, false).
        AddItem(ConnectedUsersView, 0, 3, false)

    // Main layout
    layout := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(welcomeText, 3, 0, false).
        AddItem(contentArea, 0, 1, false).
        AddItem(Messageform, 7, 0, true)

    return layout
}
