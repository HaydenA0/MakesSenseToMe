package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	textinput textinput.Model
	rn        int
	g         int
}

func initModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 20
	randomNumber := generateRandomInteger()
	return model{textinput: ti, rn: randomNumber}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(mesg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch mesg := mesg.(type) {
	case tea.KeyMsg:
		switch mesg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			v_string := m.textinput.Value()
			guess, err := strconv.Atoi(v_string)
			if err != nil {
				// later
			}
			m.g = guess
		}
	}
	m.textinput, cmd = m.textinput.Update(mesg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Enter your guess : \n%s\nThe randomNumber is %d.\nYour guess is %d.\n",
		m.textinput.View(),
		m.rn,
		m.g,
	) + "\n"
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
