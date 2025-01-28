package main

import (
	"ChatroomWithGo/models"
	"ChatroomWithGo/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var Authform = tview.NewForm()
var newUser models.UserModel
var authFlex = tview.NewFlex()
var RegisterTitle = tview.NewTextView()
var errorText = tview.NewTextView()
func DrawRegisterApp() *tview.Flex {
	db, err := utils.ConnectToDb()
	if err != nil {
		println("There was an error connecting ", err)
	}

	RegisterTitle.SetText("======= Register in the chatroom =======").SetTextAlign(tview.AlignCenter)
	errorText.SetTextColor(tcell.ColorRed)

	// Set up the message form
	Authform.AddInputField("Username", "", 0, nil, func(username string) {
		newUser.Username = username
	})

	passwordInput := tview.NewInputField().
		SetLabel("Password").
		SetMaskCharacter('*'). // Mask the password
		SetChangedFunc(func(password string) {
			newUser.Password = password
		})

	Authform.AddFormItem(passwordInput)

	Authform.AddButton("Register", func() {
		registererr := utils.InsertNewUser(db, newUser)

		if registererr != nil {
			errorText.SetText("The user already exists")
		} else {
			Pages.SwitchToPage("Login Page")
		}
	})

	Authform.AddButton("Login", func() {
		Pages.SwitchToPage("Login Page")
	})

	Authform.AddButton("Quit", func() {
		App.Stop()
	})

	// Create a Flex container for the form and error message
	formAndError := tview.NewFlex().
		SetDirection(tview.FlexRow). // Stack items vertically
		AddItem(Authform, 8, 1, true). // Give the form a fixed height (6 lines)
		AddItem(errorText, 1, 0, false) // Add the error text below the form (1 line height)

	// Create the main Flex container with title and the form + error message
	authFlex.SetDirection(tview.FlexRow).
		AddItem(RegisterTitle, 3, 0, false). // Title takes 3 lines
		AddItem(formAndError, 0, 1, true)   // Form and error take the remaining space

	return authFlex
}
