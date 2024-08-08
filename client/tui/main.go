package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
	ChatView   tview.TextView
	InputField tview.InputField
	Flex       tview.Flex
	app        tview.Application
	Incoming   chan string
}

func NewApp() *TUI {
	t := TUI{
		app:        *tview.NewApplication(),
		ChatView:   *tview.NewTextView(),
		InputField: *tview.NewInputField(),
		Flex:       *tview.NewFlex(),
	}

	t.ChatView.
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			t.app.Draw()
		})

	t.InputField.
		SetLabel("Message: ").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				msg := t.InputField.GetText()
				if msg != "" {
					fmt.Fprintf(&t.ChatView, "[green::b]%s\n", msg)
					t.InputField.SetText("")
				}
			}
		}).
		SetFieldBackgroundColor(tcell.ColorBlack)

	t.Flex.
		SetDirection(tview.FlexRow).
		AddItem(&t.ChatView, 0, 1, false).
		AddItem(&t.InputField, 1, 0, true)

	return &t
}

func (t *TUI) Run() error {
	if err := t.app.SetRoot(&t.Flex, true).Run(); err != nil {
		return err
	}

	return nil
}

func (t *TUI) InputDoneFunc(key tcell.Key) {
	if key == tcell.KeyEnter {
		msg := t.InputField.GetText()
		if msg != "" {
			// need to handle outgoing messages
			fmt.Fprintf(&t.ChatView, "[green::b]%s\n", msg)
			t.InputField.SetText("")
		}
	}
}

func (t *TUI) displayMessage(msg string) {
	fmt.Println("displaying message!")
	fmt.Fprintf(&t.ChatView, "[blue::b]%s\n", msg)
}

func (t *TUI) Error(msg string) {
	fmt.Fprintf(&t.ChatView, "[red::b]Error: %s\n", msg)
}

func (t *TUI) Listen() {
	for {
		select {
		case msg := <-t.Incoming:
			t.displayMessage(msg)
		}
	}
}
