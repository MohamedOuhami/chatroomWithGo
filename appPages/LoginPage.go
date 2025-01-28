package main

import (
	"ChatroomWithGo/models"
	"ChatroomWithGo/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var Loginform = tview.NewForm()
var connectingUser models.UserModel
var LoginTitle = tview.NewTextView()
var LoginFlex = tview.NewFlex()

var errorLoginText = tview.NewTextView()

func DrawLoginApp() *tview.Flex {

	db, err := utils.ConnectToDb()

	if err != nil {
		println("There was en error connecting ", err)

	}

	LoginTitle.SetText("======= Login to the Chatroom =======").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	errorLoginText.SetTextColor(tcell.ColorRed)
	// Set up the message form
	Loginform.AddInputField("Username", "", 0, nil, func(username string) {
		connectingUser.Username = username
	})

	passwordInput := tview.NewInputField().
		SetLabel("Password").
		SetMaskCharacter('*'). // Mask the password
		SetChangedFunc(func(password string) {
			connectingUser.Password = password
		})

	Loginform.AddFormItem(passwordInput)
	Loginform.AddButton("Login", func() {

		err = utils.ConnectWithUsername(db, connectingUser.Username, connectingUser.Password)

		if err != nil {
			errorLoginText.SetText("Error : " + err.Error())

		} else {

			Pages.SwitchToPage("Main App")
		}

	})
	Loginform.AddButton("Register", func() {

		Pages.SwitchToPage("Register Page")
	})

	Loginform.AddButton("Quit", func() {
		App.Stop()
	})

	LoginFlex.SetDirection(tview.FlexRow).
		AddItem(LoginTitle, 3, 1, false).    // Add title as *tview.TextView
		AddItem(Loginform, 0, 1, true).      // Add form
		AddItem(errorLoginText, 0, 1, false) // Add the error if there is any
	return LoginFlex
}
