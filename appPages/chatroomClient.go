// In chatroomClient.go

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
    welcomeText        = tview.NewTextView().SetTextColor(tcell.ColorGreen)
    Messageform        = tview.NewForm()
    messageToSend      string
    MessageView        *tview.TextView
    ConnectedUsersView *tview.TextView
    conn               net.Conn
    username           string // Add username variable
)

func main() {
    // Initialize the basic UI components
    App.SetRoot(Pages, true)
    
    // Only add login and register pages initially
    Pages.AddPage("Login Page", DrawLoginApp(), true, true)
    Pages.AddPage("Register Page", DrawRegisterApp(), true, false)
    
    if err := App.Run(); err != nil {
        fmt.Printf("Error running application: %s\n", err)
    }
}
