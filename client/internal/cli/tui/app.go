package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func NewModel() (*model, error) {
	// We need to initialize a new text input model.
	ti := textinput.New()
	ti.CharLimit = 30
	ti.Placeholder = "Type in your event"
	// Nest the text input in our application state.
	return &model{input: ti}, nil
}

type model struct {
	nameInput string
	listInput string
	event     string
	// Add the text input to our main application state.  It's a subcomponent
	// which has its own state, etc.
	input textinput.Model
}

func (m model) Init() tea.Cmd {
	// Call Init() on our submodel.  If we had > 1 submodel and command, we would
	// create a slice of commands to batch:
	//
	// return tea.Batch(cmds...)
	//cmd := m.input.Init()
	//return cmd
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		_, _ = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlBackslash:
			return m, tea.Quit
		}
	}
	// We call Bubbletea using our model as the top-level application.  Bubbletea
	// will call Update() in our model only.  It's up to us to call Update() on
	// our text input to update its state.  Without this, typing won't fill out
	// the text box.
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	// store the text inputs value in our top-level state.
	m.nameInput = m.input.Value()
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.event != "" {
		return fmt.Sprintf("You've selected: %s", m.event)
	}

	b := &strings.Builder{}
	b.WriteString("Enter your event:\n")
	// render the text input.  All we need to do to show the full
	// input is call View() and return the string.
	b.WriteString(m.input.View())
	return b.String()
}
