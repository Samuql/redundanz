package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// describes the application state
type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		choices: getSelection(), // currently string

		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("redundanz")
}

// handle updates and modify model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			getFileNames(m.choices[m.cursor], true)
			// TODO handle multiple params
			// ctn := []string{getWd()}
			// ctn = append(ctn, getContent(m.choices[m.cursor])...)
			// fmt.Printf("\n %v\n", ctn)
			return m, tea.Quit
		}
	}

	return m, nil
}

// render the application
func (m model) View() string {
	s := "Which folder do you want to eliminate redundancy in?\n[space] to select, [enter] to run.\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
