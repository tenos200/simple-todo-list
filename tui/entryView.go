package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"todoList/dbhandler"
)

type mode int

const (
	nav mode = iota
	edit
	create
)

type model struct {
	mode   mode
	choice string
	tasks  []dbhandler.DbRow
	list   list.Model
	input  textinput.Model
}

// Item struct that is used for list items.
type item struct {
	name, status, date string
}

// These function needs to be implemeted as a part of bubbles list component.
func (i item) Title() string { return i.name }
func (i item) Description() string {
	return fmt.Sprintf("Status: %s Date: %s", i.status, i.date)
}
func (i item) Status() string {
	return i.status
}

func (i item) Date() string {
	return i.date
}
func (i item) FilterValue() string { return i.name }

func InitialModel(fetchedTasks []dbhandler.DbRow) model {

	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "Project name"
	input.CharLimit = 250
	input.Width = 50

	var items []list.Item
	for _, task := range fetchedTasks {
		items = append(items, item{name: task.Name,
			status: task.Status, date: task.Date})
	}

	return model{
		mode:  nav,
		tasks: fetchedTasks,
		list:  list.New(items, list.NewDefaultDelegate(), 40, 40),
		input: input,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.mode {
	case nav:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter", " ":
				m.input.Focus()
				m.mode = edit
				m.input.Placeholder = m.list.SelectedItem().FilterValue()
				return m, nil
			}
		}
		m.list, cmd = m.list.Update(msg)

	case edit:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.mode = nav
				// Ensure that we clear text input after exit.
				m.input.SetValue("")
				m.input.Blur()
				return m, nil
			case "enter":
				// TODO: We need multiple models to be here.
				inputValue := m.input.Value()
				var newItem item
				newItem.name = inputValue
				newItem.date = "test date"
				newItem.status = "test"

				// Get the correct index for current item.
				index := m.list.Index()
				m.list.SetItem(index, newItem)
				m.input.Blur()

				// Ensure to clear the input after we have entered.
				m.input.SetValue("")
				m.mode = nav
			}
		}
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.Update(msg)
		return m, inputCmd
	}
	return m, cmd
}

func (m model) View() string {
	// TODO: We need more models to load several stages edit, create, delete
	if m.mode == edit {
		return lipgloss.JoinHorizontal(lipgloss.Center, m.input.View())
	} else {
		// Apply the style for the tasks list
		return lipgloss.JoinHorizontal(lipgloss.Center, m.list.View())
	}
}
