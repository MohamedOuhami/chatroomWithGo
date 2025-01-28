package main 

import (
	"net"

	"github.com/gdamore/tcell/v2"
)

var (
)

func SetupUI(conn net.Conn) {

	App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			App.Stop()
		}
		return event
	})

	Mainlayout := DrawMainApp(conn)
	AuthLayout := DrawAuthApp()
	LoginLayout := DrawLoginApp()


	Pages.AddPage("Main App", Mainlayout, true, false)
	Pages.AddPage("Register Page", AuthLayout, true, true)
	Pages.AddPage("Login Page", LoginLayout, true, false)

	if err := App.SetRoot(Pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
