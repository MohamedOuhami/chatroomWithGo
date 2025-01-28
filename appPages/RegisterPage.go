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

func DrawAuthApp() *tview.Flex {

	db, err := utils.ConnectToDb()

	if err != nil {
		println("There was en error connecting ", err)

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
		err = utils.InsertNewUser(db, newUser)

		if err != nil {

			errorText.SetText("The username already exists")
		} else {

			Pages.SwitchToPage("Main App")

		}
	})

	Authform.AddButton("Login", func() {

		Pages.SwitchToPage("Login Page")
	})
	Authform.AddButton("Quit", func() {
		App.Stop()
	})

	authFlex.SetDirection(tview.FlexRow).AddItem(RegisterTitle, 3, 1, false).AddItem(Authform, 0, 1, false).AddItem(errorText, 0, 1, false)

	return authFlex
}
