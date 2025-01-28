
package main

import (
	"ChatroomWithGo/models"
	"ChatroomWithGo/utils"
	"fmt"
	"net"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var Loginform = tview.NewForm()
var connectingUser models.UserModel
var LoginTitle = tview.NewTextView()
var LoginFlex = tview.NewFlex()
var errorLoginText = tview.NewTextView()

func DrawLoginApp() *tview.Flex {
	// Database connection
	db, err := utils.ConnectToDb()
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
	}

	// Title configuration
	LoginTitle.SetText("======= Login to the Chatroom =======").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetTextColor(tcell.ColorGreen)

	// Error message configuration
	errorLoginText.SetText("").
		SetTextColor(tcell.ColorRed).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	// Form inputs
	Loginform.AddInputField("Username", "", 0, nil, func(username string) {
		connectingUser.Username = username
	})

	passwordInput := tview.NewInputField().
		SetLabel("Password").
		SetMaskCharacter('*').
		SetChangedFunc(func(password string) {
			connectingUser.Password = password
		})

	Loginform.AddFormItem(passwordInput)

	// Form buttons
	Loginform.AddButton("Login", func() {
		err = utils.ConnectWithUsername(db, connectingUser.Username, connectingUser.Password)
		if err != nil {
			errorLoginText.SetText("Wrong credentials. Please verify your username and password.")
		} else {
			// Initialize connection
			if err := initializeConnection(connectingUser.Username); err != nil {
				errorLoginText.SetText("Connection error: " + err.Error())
				return
			}

			// Navigate to the main app page
			mainApp := DrawMainApp(conn)
			Pages.AddPage("Main App", mainApp, true, false)
			Pages.SwitchToPage("Main App")
		}
	})

	Loginform.AddButton("Register", func() {
		Pages.SwitchToPage("Register Page")
	})

	Loginform.AddButton("Quit", func() {
		App.Stop()
	})

	// Layout: Title, Form, and Error Message
	LoginFlex.SetDirection(tview.FlexRow).
		AddItem(LoginTitle, 3, 1, false).          // Fixed height for title
		AddItem(Loginform, 10, 1, true).          // Fixed height for form
		AddItem(errorLoginText, 1, 0, false).     // Error message directly below the form
		AddItem(nil, 0, 1, false)                 // Spacer to prevent stretching

	return LoginFlex
}

// Initialize connection with error handling
func initializeConnection(username string) error {
	var err error
	conn, err = net.Dial("tcp", "localhost:8000")
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	// Send username immediately
	_, err = conn.Write([]byte(username))
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to send username: %v", err)
	}

	// Initialize views
	MessageView = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() { App.Draw() }).
		SetScrollable(true)

	ConnectedUsersView = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() { App.Draw() })

	welcomeText.SetText(fmt.Sprintf("Welcome %s! (q to quit)", username))

	// Start receiving messages
	go ReceiveMessages(conn, App, ConnectedUsersView, MessageView)

	return nil
}
