package main
import (
    "net"
    "github.com/rivo/tview"
)

var (
)

func DrawMainApp(conn net.Conn) *tview.Flex {

    // Set up the message form
    Messageform.AddInputField("Message", "", 0, nil, func(message string) {
        messageToSend = message
    })

    Messageform.AddButton("Send", func() {
        SendMessage(conn, MessageView, Messageform, messageToSend)
    })

    Messageform.AddButton("Quit", func() {
        App.Stop()
    })

    // Add borders and titles to the views
    messageBox := tview.NewBox().
        SetBorder(true).
        SetTitle("Chat Messages")
    
    usersBox := tview.NewBox().
        SetBorder(true).
        SetTitle("Connected Users")

    // Create a flex for the message view with border
    messageWithBorder := tview.NewFlex().
        AddItem(messageBox, 0, 1, false)
    messageWithBorder.Box = messageBox
    messageWithBorder.AddItem(MessageView, 0, 1, false)

    // Create a flex for the users view with border
    usersWithBorder := tview.NewFlex().
        AddItem(usersBox, 0, 1, false)
    usersWithBorder.Box = usersBox
    usersWithBorder.AddItem(ConnectedUsersView, 0, 1, false)

    middleLayout := tview.NewFlex().
        SetDirection(tview.FlexColumn).
        AddItem(messageWithBorder, 0, 2, false).
        AddItem(usersWithBorder, 0, 1, false)

    layout := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(welcomeText, 3, 1, false).
        AddItem(middleLayout, 0, 2, false).
        AddItem(Messageform, 3, 1, true)

    return layout
}


