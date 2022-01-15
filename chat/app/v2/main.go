package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type model struct {
	todos []string
	cursor int
	selected map[int]struct{}
}

func (m model) Init() tea.Cmd {
	return nil
}

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
			if m.cursor < len(m.todos) -1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			}else {
				m.selected[m.cursor] = struct{}{}
			}

		}
	}
	return m, nil
}

func (m model) View() string{
	s := "todo list:\n\n"
	for i, choice := range m.todos {
		curor := " "
		if m.cursor == i {
			curor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", curor, checked, choice)
	}
	s += "\nPress q to quit.\n"


	return s
}

var initModl = model{
	todos: []string{
		"cleaning",
		"wash clothes",
		"write a blog",
	},
	selected: make(map[int]struct{}),
}

func main() {
	cmd := tea.NewProgram(initModl)
	if err := cmd.Start(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}
}