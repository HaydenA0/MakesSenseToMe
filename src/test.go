package main

//
// import (
// 	"fmt"
// 	"os"
//
// 	tea "github.com/charmbracelet/bubbletea"
// )
//
// type model struct {
// 	choices  []string
// 	cursor   int
// 	selected map[int]struct{}
// }
//
// func initModel() model {
// 	return model{choices: []string{"Something", "Choice2", "Choice3"}, selected: make(map[int]struct{})}
// }
//
// func (m model) Init() tea.Cmd {
// 	return nil
// }
//
// func (m model) View() string {
// 	s := "What choice ?\n"
// 	for i, choice := range m.choices {
// 		cursor := " "
// 		if m.cursor == i {
// 			cursor = ">"
// 		}
// 		checked := " "
// 		if _, ok := m.selected[i]; ok {
// 			checked = "o"
// 		}
// 		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
// 	}
// 	s += "\nPress q to quit.\n"
// 	return s
// }
//
// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
//
// 		case "ctrl+c", "q":
// 			return m, tea.Quit
//
// 		case "k":
// 			if m.cursor > 0 {
// 				m.cursor--
// 			}
// 		case "j":
// 			if m.cursor < len(m.choices)-1 {
// 				m.cursor++
// 			}
// 		case "enter":
// 			_, ok := m.selected[m.cursor]
// 			if ok {
// 				delete(m.selected, m.cursor)
// 			} else {
// 				m.selected[m.cursor] = struct{}{}
// 			}
// 		}
// 	}
// 	return m, nil
// }
//
// // func main() {
// // 	p := tea.NewProgram(initModel())
// // 	if _, err := p.Run(); err != nil {
// // 		fmt.Printf("Alas, there's been an error: %v", err)
// // 		os.Exit(1)
// // 	}
// // }
