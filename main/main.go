package main

import (
	"fmt"
	"strings"
	// "net/http"
	// "os"
	// "time"
	// "os/exec"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	// bub "github.com/charmbracelet/bubbles"
)

// const url = "https://charm.sh/"

type errMsg error

func main() {
	var bd string
	// var m buildDirModel
	// m = initialModel()
	m := initialModel()
	// p := tea.NewProgram(initialModel())
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		log.Fatal(err)
	}
	log.Printf(m.textInput.Value())
	// bd = m.getBuildDir()
	bd = m.build_dir_val
	log.Printf("m.val=%s", bd)
}

type buildDirModel struct {
	textInput     textinput.Model

	state         string

	build_dir_val string
	user_val      string

	err           error
}

// func (m buildDirModel) getBuildDir() string {
// 	// return m.textInput.Value()
// 	return m.val
// }

func initialModel() buildDirModel {
	ti := textinput.New()
	ti.Placeholder = "/home/cbs/github/build_man/test"
	ti.Focus()
	// ti.CharLimit = 150
	// ti.Width = 20

	return buildDirModel{
		textInput:     ti,
		build_dir_val: ti.Placeholder,
		state:         "args_build_dir",
		err:           nil,
	}
}

// Init() ///////////////////////////////////////////////////////////////////////
func (m buildDirModel) Init() tea.Cmd {
	// cmd := m.textInput.Init()
	textInputBlinkCmd := textinput.Blink
	// return textinput.Blink
	// return textInputBlinkCmd
	return tea.Batch(textInputBlinkCmd)
}

// View() /////////////////////////////////////////////////////////////////////
func (m buildDirModel) View() string {
	s := &strings.Builder{}

	if m.build_dir_val == nil {
		// Configure application
		s.WriteString(argsView(m))
	} else {
		s.WriteString("Unknown m.state: " + m.state)
		// return fmt.Sprintf("Text selection: %s", m.val)
		// s.WriteString(fmt.Sprintf("Text selection: %s", m.build_dir_val))
	}

	return s.String()
}

func argsView(m buildDirModel) string {
	return "TODO"

	if m.build_dir_var == nil {
		s.WriteString("Build dir location: (ESC to quit)\n")
		s.WriteString(m.textInput.View())
	}
	if m.state == "args_build_dir" {
		// return fmt.Sprintf(
		// 	"Build dir location: (ESC to quit)\n%s\n",
		// 	m.textInput.View(),
		// ) + "\n"
		s.WriteString("Build dir location: (ESC to quit)\n")
		s.WriteString(m.textInput.View())
		// s.WriteString(fmt.Sprintf("args_build_dir=%s", m.build_dir_val))
	} else if m.state == "args_user" {
		s.WriteString("Build user: (ESC to quit)\n")
		s.WriteString(m.textInput.View())
		// b.WriteString(fmt.Sprintf("args_user=%s", m.user_val))
	} else if m.state == "monitor" {
		s.WriteString("TODO")
	}

}

// Update() /////////////////////////////////////////////////////////////////////
func (m buildDirModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			switch m.state {
			case "args_build_dir":
				// Store textInput value in model
				m.build_dir_val = m.textInput.Value()
				// log.Printf("val=%s", m.val)
				// Move to next input state
				// m.state = "args_user"
				m.state = "monitor" // TODO: Should be "args_user"
			case "args_user":
				m.user_val = m.textInput.Value()
				m.state = "monitor"
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	// return m, cmd
	return m, tea.Batch(cmds...)
}


// type model struct {
// 	status int
// 	err    error
// }

// func checkServer() tea.Msg {

// 	// Create an HTTP client and make a GET request.
// 	c := &http.Client{Timeout: 10 * time.Second}
// 	res, err := c.Get(url)

// 	if err != nil {
// 		// There was an error making our request. Wrap the error we received
// 		// in a message and return it.
// 		return errMsg{err}
// 	}
// 	// We received a response from the server. Return the HTTP status code
// 	// as a message.
// 	return statusMsg(res.StatusCode)
// }

// type zpoolMsg struct {err error}
// func checkZpool() tea.Msg {

// 	// Create an HTTP client and make a GET request.
// 	c := exec.Command("vim")
// 	return tea.ExecProcess(c, func(err error) tea.Msg {
// 		return zpoolMsg{err}
// 	})
// 	// res, err := c.Get(url)

// 	// if err != nil {
// 	// 	// There was an error making our request. Wrap the error we received
// 	// 	// in a message and return it.
// 	// 	return errMsg{err}
// 	// }
// 	// We received a response from the server. Return the HTTP status code
// 	// as a message.
// 	// return statusMsg(res.StatusCode)
// }

// type statusMsg int

// type errMsg struct{ err error }

// // For messages that contain errors it's often handy to also implement the
// // error interface on the message.
// func (e errMsg) Error() string { return e.err.Error() }

// func (m model) Init() (tea.Cmd) {
// 	// return checkServer
// 	return checkZpool
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {

// 	case statusMsg:
// 		// The server returned a status message. Save it to our model. Also
// 		// tell the Bubble Tea runtime we want to exit because we have nothing
// 		// else to do. We'll still be able to render a final view with our
// 		// status message.
// 		m.status = int(msg)
// 		return m, tea.Quit

// 	case errMsg:
// 		// There was an error. Note it in the model. And tell the runtime
// 		// we're done and want to quit.
// 		m.err = msg
// 		return m, tea.Quit

// 	case tea.KeyMsg:
// 		// Ctrl+c exits. Even with short running programs it's good to have
// 		// a quit key, just incase your logic is off. Users will be very
// 		// annoyed if they can't exit.
// 		if msg.Type == tea.KeyCtrlC {
// 			return m, tea.Quit
// 		}
// 	}

// 	// If we happen to get any other messages, don't do anything.
// 	return m, nil
// }

// func (m model) View() string {
// 	// If there's an error, print it out and don't do anything else.
// 	if m.err != nil {
// 		return fmt.Sprintf("We had some trouble: %v\n", m.err)
// 	}

// 	// Tell the user we're doing something.
// 	s := fmt.Sprintf("Checking %s ... ", url)

// 	// When the server responds with a status, add it to the current line.
// 	if m.status > 0 {
// 		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
// 	}

// 	// Send off whatever we came up with above for rendering.
// 	return s + "\n"
// }
