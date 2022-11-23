package main

import (
	"fmt"
	"strings"
	"os"
	"syscall"
	// "time"
	"strconv"
	"time"
	"log"
	// "os/exec"
	// "io/ioutil"
	"path/filepath"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	// bub "github.com/charmbracelet/bubbles"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)


type errMsg error

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Help  key.Binding
	Quit  key.Binding
}
// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up},   // first column
		{k.Down}, // second column
		{k.Help}, // third column
		{k.Quit}, // fourth column
	}
}
var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}


func main() {
	// Log to file if DEBUG
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	log.Printf("\n\nvvvvvvvvvvv starting vvvvvvvvvvvvvvv\n\n")

	// var m mainModel
	// m = initialModel()
	m := initialModel()

	// p := tea.NewProgram(initialModel())
	p := tea.NewProgram(m)
	// if model,err := p.Run(); err != nil {
	if _,err := p.Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		log.Fatal(err)
	}

	// log.Printf(m.build_dir_input.Value())
	// // bd = m.getBuildDir()
	// bd = m.build_dir_val
	// log.Printf("main::m.val=%s", bd)

}

type SynshFileStruct struct {
	file string
	// startTime int
	startTime time.Time
	runTime time.Duration
}

type mainModel struct {
	build_dir_input  textinput.Model
	build_dir_val string
	user_input textinput.Model
	user_val      string

	synsh_files []SynshFileStruct
	build_job_dirs []string

	// cursor int // Selection in table
	// selected map[int]struct{} // Which build job is selected

	conf_done bool
	state     string

	tbl table.Model

	keys keyMap

	// helpStyle lipgloss.Style
	help help.Model

	err error
}


func initialModel() mainModel {
	build_dir_ti := textinput.New()
	build_dir_ti.Placeholder = "/home/cbs/github/build-mon/test/FPGA"
	build_dir_ti.Focus()
	// build_dir_ti.CharLimit = 150
	// build_dir_ti.Width = 20

	user_ti := textinput.New()
	user_ti.Placeholder = "CBS"
	user_ti.Focus()

	columns := []table.Column{
		{Title: "Job", Width: 90},
		{Title: "Start time", Width: 20},
		{Title: "Run time", Width: 25},
		// {Title: "Country", Width: 10},
		// {Title: "Population", Width: 10},
	}
	rows := []table.Row{
		{"No jobs", ""},
		// {"1", "Tokyo", "Japan", "37,274,000"},
		// {"2", "Delhi", "India", "32,065,760"},
		// {"3", "Shanghai", "China", "28,516,904"},
		// {"4", "Dhaka", "Bangladesh", "22,478,116"},
		// {"5", "São Paulo", "Brazil", "22,429,800"},
		// {"6", "Mexico City", "Mexico", "22,085,140"},
		// {"7", "Cairo", "Egypt", "21,750,020"},
		// {"8", "Beijing", "China", "21,333,332"},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return mainModel{
		build_dir_input:  build_dir_ti,
		build_dir_val:    build_dir_ti.Placeholder,
		user_input:       user_ti,
		user_val:         user_ti.Placeholder,
		conf_done:        true, // Skipping conf until later version
		state:            "build_dir",

		// synsh_files:        make(chan SynshFileStruct[string, int]),
		// cursor:           int,
		// selected:         make(map[string]struct{}),

		tbl: t,

		keys: keys,

		// helpStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#008000")),
		help:      help.New(),

		err:              nil,
	}
}

func (m *mainModel) getJobs() tea.Msg {
	var files []SynshFileStruct

	// Find .synsh files and note creation time
	err := filepath.Walk(m.build_dir_val + "/New_file", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".synsh" {
			tmp := SynshFileStruct{path, time.Unix(0,0), time.Duration(0)}
			tmp.file = path

			// Get file creation time
			fi, err := os.Stat(path)
			if err != nil {
				return err
			}
			stat := fi.Sys().(*syscall.Stat_t)
			ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
			tmp.startTime = ctime

			// Get job runtime
			tmp.runTime = time.Now().Sub(ctime)

			files = append(files, tmp)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)

	}
	m.synsh_files = files
	log.Printf("[getJobs::len(m.synsh_files)]: " + strconv.Itoa(len(m.synsh_files)))

	// Create list of build jobs dirs for each .synsh file
	var dirs []string
	for _,strct := range files {
		// dir := filepath.Join(filepath.Dir(strct.file.(string)), "/../syntese/", strings.TrimSuffix(filepath.Base(strct.file.(string)), filepath.Ext(strct.file.(string))))
		dir := filepath.Join(filepath.Dir(strct.file), "/../syntese/", strings.TrimSuffix(filepath.Base(strct.file), filepath.Ext(strct.file)))
		dirs = append(dirs, dir)
	}

	m.build_job_dirs = dirs
	log.Printf("[getJobs::len(m.build_job_dirs)]: " + strconv.Itoa(len(m.build_job_dirs)))

	// return 0
	return m
}


// Init() ///////////////////////////////////////////////////////////////////////
type tickMsg time.Time
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m mainModel) Init() tea.Cmd {

	// // cmd := m.build_dir_input.Init()
	// textInputBlinkCmd := textinput.Blink
	// // return textinput.Blink
	// // return textInputBlinkCmd
	// return tea.Batch(textInputBlinkCmd)

	// return m.getJobs
	return tickCmd()
	// return nil
}

// View() /////////////////////////////////////////////////////////////////////
func (m mainModel) View() string {
	s := &strings.Builder{}

	log.Printf("[View::conf_done]: %b", m.conf_done)

	// Print help text
	helpView := m.help.View(m.keys)
	height := 3 - strings.Count(s.String(), "\n") - strings.Count(helpView, "\n")
	s.WriteString(strings.Repeat("\n", height))
	s.WriteString(helpView)
	s.WriteString("\n")

	// s.WriteString(m.helpStyle.Render("Hello"))
	// s.WriteString("\n")

	if m.conf_done == false {
		// Configure application
		s.WriteString(argsView(m))
	} else if m.conf_done == true {
		s.WriteString(monView(m))
	} else {
		s.WriteString("Unknown m.state: " + m.state)
	}

	return s.String()
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func monView(m mainModel) string {
	s := &strings.Builder{}

	log.Printf("--------monView-----------\n")
	log.Printf("[monView::m.build_dir_val]  " + m.build_dir_val + "\n")

	log.Printf("[monView::m.build_job_dirs] len=" + strconv.Itoa(len(m.build_job_dirs)))
	// log.Printf("[monView::m.synsh_files]    len=" + strconv.Itoa(len(m.synsh_files)))
	// for _, val := range m.synsh_files {
	// 	log.Printf("[monView]:" + val)
	// }

	s.WriteString(baseStyle.Render(m.tbl.View()) + "\n")

	log.Printf("--------------------------\n")
	return s.String()
}

func argsView(m mainModel) string {
	s := &strings.Builder{}

	log.Printf("argsView::m.state=%s", m.state)

	if m.state == "build_dir" {
		s.WriteString("Build dir location: (ESC to quit)\n")
		s.WriteString(m.build_dir_input.View())
	} else if m.state == "build_user" {
		s.WriteString("Build user: (ESC to quit)\n")
		s.WriteString(m.user_input.View())
	}

	return s.String()
}

// Update() /////////////////////////////////////////////////////////////////////

func (m *mainModel) ConfUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	return m, cmd
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
	case tea.KeyMsg:
		// switch msg.Type {
		// case tea.KeyCtrlC, tea.KeyEsc:
		// 	return m, tea.Quit
		// }
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.tbl.Focused() {
				m.tbl.Blur()
			} else {
				m.tbl.Focus()
			}
		case "enter":
			return m, tea.Batch(
				tea.Printf("Selected %s", m.tbl.SelectedRow()[0]),
			)
		// case key.Matches(msg, m.keys.Help):
		case "?":
			m.help.ShowAll = !m.help.ShowAll}
	case tickMsg:
		log.Printf("[Update::tickMsg]")
		// Execute new tickCmd
		cmds = append(cmds, tickCmd())
	case errMsg:
		m.err = msg
		return m, nil
	}

	// Look for new .synsh files and update table
	m.getJobs()
	log.Printf("[Update::m.synsh_files]   len=" + strconv.Itoa(len(m.synsh_files)))
	s := &strings.Builder{}
	var eol string
	for i,f := range m.synsh_files {
		// Terminate each row with newline except for last row
		if (i==len(m.synsh_files)-1) {eol = ""} else {eol = "\n"}
		// Add string for each row. Separate columns with tabs. // FIXME: Choose other separator than tab since filename might have tabs(?)
		s.WriteString(
			f.file + "\t" +
			f.startTime.Format("2006-01-02 15:04:09") + "\t" +
			f.runTime.String() +
			eol)
			// "\n")

	}
	m.tbl.FromValues(s.String(), "\t")
	m.tbl, cmd = m.tbl.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
