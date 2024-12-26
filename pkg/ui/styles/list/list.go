package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Item struct {
	Name, Desc string
}

func (i Item) Title() string       { return i.Name }
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.Name }

type List struct {
	list list.Model
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m List) View() string {
	return docStyle.Render(m.list.View())
}

func NewList(title string, items []Item) List {
	i := []list.Item{}
	for _, item := range items {
		i = append(i, item)
	}
	l := List{list: list.New(i, list.NewDefaultDelegate(), 0, 0)}

	l.list.Title = title

	return l
}
