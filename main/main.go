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

	log.Printf(m.build_dir_input.Value())
	// bd = m.getBuildDir()
	bd = m.build_dir_val
	log.Printf("main::m.val=%s", bd)
}

type buildDirModel struct {
	build_dir_input  textinput.Model
	build_dir_val string
	user_input textinput.Model
	user_val      string

	synsh_files []string

	// cursor int // Selection in table
	// selected map[int]struct{} // Which build job is selected

	conf_done bool
	state     string

	err error
}

// func (m buildDirModel) getBuildDir() string {
// 	// return m.build_dir_input.Value()
// 	return m.val
// }

func initialModel() buildDirModel {
	build_dir_ti := textinput.New()
	build_dir_ti.Placeholder = "/home/cbs/github/build_man/test"
	build_dir_ti.Focus()
	// build_dir_ti.CharLimit = 150
	// build_dir_ti.Width = 20

	user_ti := textinput.New()
	user_ti.Placeholder = "CBS"
	user_ti.Focus()

	return buildDirModel{
		build_dir_input:  build_dir_ti,
		build_dir_val:    build_dir_ti.Placeholder,
		user_input: user_ti,
		user_val:         user_ti.Placeholder,
		conf_done:        false,
		state:            "build_dir",
		// cursor:           int,
		// selcted:          make(map[int]struct{}),
		err:              nil,
	}
}

// Init() ///////////////////////////////////////////////////////////////////////
func (m buildDirModel) Init() tea.Cmd {
	// cmd := m.build_dir_input.Init()
	textInputBlinkCmd := textinput.Blink
	// return textinput.Blink
	// return textInputBlinkCmd
	return tea.Batch(textInputBlinkCmd)
}

// View() /////////////////////////////////////////////////////////////////////
func (m buildDirModel) View() string {
	s := &strings.Builder{}

	// log.Printf("View::conf_done=%b", m.conf_done)

	// s.WriteString("Build dir location: (ESC to quit)\n")
	// s.WriteString(m.build_dir_input.View())

	if m.conf_done == false {
		// Configure application
		s.WriteString(argsView(m))
	} else if m.conf_done == true {
		s.WriteString(monView(m))
	} else {
		s.WriteString("Unknown m.state: " + m.state)
		// return fmt.Sprintf("Text selection: %s", m.val)
		// s.WriteString(fmt.Sprintf("Text selection: %s", m.build_dir_val))
	}

	return s.String()
}

func monView(m buildDirModel) string {
	return "TODO: Monitor state" + ", " + m.build_dir_val + ", " + m.user_val
}

func argsView(m buildDirModel) string {
	s := &strings.Builder{}

	// log.Printf("argsView::m.state=%s", m.state)

	if m.state == "build_dir" {
		s.WriteString("Build dir location: (ESC to quit)\n")
		s.WriteString(m.build_dir_input.View())
	} else if m.state == "build_user" {
		// return fmt.Sprintf(
		// 	"Build dir location: (ESC to quit)\n%s\n",
		// 	m.build_dir_input.View(),
		// ) + "\n"
		s.WriteString("Build user: (ESC to quit)\n")
		s.WriteString(m.user_input.View())
		// s.WriteString(fmt.Sprintf("args_build_dir=%s", m.build_dir_val))
		// } else {
		// 	log.Fatal("FATAL: Unhandled state in argsView")
		// 	s.WriteString("TODO")
	}

	return s.String()
}

// Update() /////////////////////////////////////////////////////////////////////
// func (m buildDirModel) ConfUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
func (m *buildDirModel) ConfUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		// case tea.KeyCtrlC, tea.KeyEsc:
		// 	return m, tea.Quit
		case tea.KeyEnter:
			switch m.state {
			case "build_dir":
				// Store build_dir_input value in model
				m.build_dir_val = m.build_dir_input.Value()
				// Move to next input state
				m.state = "build_user"
				// m.build_dir_input, cmd = m.build_dir_input.Update(msg)
			case "build_user":
				m.user_val = m.user_input.Value()
				m.state = "monitor"
				m.conf_done = true
				// m.user_input, cmd = m.user_input.Update(msg)
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	// cmds = append(cmds, cmd)
	// cmds = append(cmds, cmd)

	tea.LogToFile("debug.log", "debug") // FIXME: Seems to log also send "log.Printf" to file with this command?

	log.Printf("ConfUpdate::m.state=%s", m.state)
	log.Printf("ConfUpdsate::m.build_dir_input.Value()=%s", m.build_dir_input.Value())
	log.Printf("ConfUpdsate::m.user_input.Value()=%s", m.user_input.Value())


	return m, cmd
}

func (m buildDirModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// switch msg.String() {
		// case "up", "k":
		// 	//
		// }
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	// log.Printf("Update::m.conf_done=%b", m.conf_done)
	// log.Printf("Update::m.state=%s", m.state)
	if (m.conf_done == false) {
		_, cmd = m.ConfUpdate(msg)
		cmds = append(cmds, cmd)
	} else {
		// m, cmd = MonitorUpdate(msg)
		// cmds = append(cmds, cmd)
	}





	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// 	switch msg.Type {
	// 	case tea.KeyCtrlC, tea.KeyEsc:
	// 		return m, tea.Quit
	// 	case tea.KeyEnter:
	// 		switch m.state {
	// 		case "build_dir":
	// 			// Store build_dir_input value in model
	// 			m.build_dir_val = m.build_dir_input.Value()
	// 			// log.Printf("val=%s", m.val)
	// 			// Move to next input state
	// 			// m.state = "args_user"
	// 			m.state = "build_user"
	// 		case "build_user":
	// 			m.user_val = m.build_dir_input.Value()
	// 			m.state = "monitor"
	// 			m.conf_done = true
	// 		}
	// 	}
	// case errMsg:
	// 	m.err = msg
	// 	return m, nil
	// }



	m.build_dir_input, cmd = m.build_dir_input.Update(msg)
	cmds = append(cmds, cmd)
	m.user_input, cmd = m.user_input.Update(msg)
	cmds = append(cmds, cmd)

	// return m, nil
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
