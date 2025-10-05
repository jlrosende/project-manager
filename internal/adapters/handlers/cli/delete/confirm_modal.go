package delete

import (
	"context"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/components"
	stylespkg "github.com/jlrosende/project-manager/internal/adapters/handlers/tui/styles"
)

const (
	confirmButtonLabel = "Delete"
	cancelButtonLabel  = "Cancel"
)

func confirmWithModal(
	ctx context.Context,
	req ConfirmationRequest,
	stdin, stdout *os.File,
	palette map[string]string,
) (bool, error) {
	colors := palette
	if colors == nil {
		colors = tui.DefaultPalette()
	} else {
		colors = clonePalette(colors)
	}

	tokens := stylespkg.NewTokensFromHex(colors)
	theme := stylespkg.BuildTheme(tokens)
	cs := stylespkg.BuildComponentStyles(theme)

	model := newConfirmModel(req, cs, colors)

	program := tea.NewProgram(model, tea.WithContext(ctx), tea.WithInput(stdin), tea.WithOutput(stdout))

	finalModel, err := program.Run()
	if err != nil {
		return false, err
	}

	m, ok := finalModel.(*confirmModel)
	if !ok {
		return false, fmt.Errorf("unexpected confirmation model type %T", finalModel)
	}

	if m.err != nil {
		return false, m.err
	}

	return m.confirmed, nil
}

type confirmModel struct {
	req               ConfirmationRequest
	modal             components.Modal
	buttons           components.ButtonGroup
	labelStyle        lipgloss.Style
	projectValueStyle lipgloss.Style
	scopeValueStyle   lipgloss.Style
	valueStyle        lipgloss.Style
	helpStyle         lipgloss.Style
	confirmed         bool
	err               error
}

func newConfirmModel(req ConfirmationRequest, cs stylespkg.ComponentStyles, palette map[string]string) *confirmModel {
	modal := components.NewModal(
		"\nConfirm deletion",
		"",
		"",
		cs.ModalWrap,
		cs.ModalTitle,
		cs.ModalBody,
		cs.ModalFooter,
	)
	modal.Visible = true

	cancelBtn := components.NewButton(
		cancelButtonLabel,
		components.Info,
		cs.ButtonInfo,
		cs.ButtonPrimary,
		cs.ButtonDisabled,
	)
	confirmBtn := components.NewButton(
		confirmButtonLabel,
		components.Danger,
		cs.ButtonDanger,
		cs.ButtonPrimary,
		cs.ButtonDisabled,
	)

	buttons := components.NewButtonGroup([]components.Button{cancelBtn, confirmBtn})
	buttons.Active = true

	labelColor := paletteValue(palette, "section", "#81A1C1")
	projectColor := paletteValue(palette, "title", "#88C0D0")
	scopeColor := paletteValue(palette, "buttonSelBg", "#A3BE8C")
	valueColor := paletteValue(palette, "text", "#D8DEE9")

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(labelColor)).Bold(true)
	projectValueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(projectColor)).Bold(true)
	scopeValueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(scopeColor)).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(valueColor))

	return &confirmModel{
		req:               req,
		modal:             modal,
		buttons:           buttons,
		labelStyle:        labelStyle,
		projectValueStyle: projectValueStyle,
		scopeValueStyle:   scopeValueStyle,
		valueStyle:        valueStyle,
		helpStyle:         cs.ModalFooter,
	}
}

func (m *confirmModel) Init() tea.Cmd { return nil }

func (m *confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		switch v.String() {
		case "esc":
			m.confirmed = false
			return m, tea.Quit
		case "ctrl+c":
			m.err = context.Canceled
			return m, tea.Quit
		}
	case components.ButtonChosenMsg:
		if v.Label == confirmButtonLabel {
			m.confirmed = true
		} else {
			m.confirmed = false
		}

		return m, tea.Quit
	}

	var cmd tea.Cmd

	m.buttons, cmd = m.buttons.Update(msg)

	return m, cmd
}

func (m *confirmModel) View() string {
	modal := m.modal
	modal.Body = m.bodyText()
	help := m.helpStyle.Render("Use ←/→ to choose, Enter to confirm, Esc to cancel")

	return fmt.Sprintf("%s\n%s\n%s", modal.View(), m.buttons.View(), help)
}

func (m *confirmModel) bodyText() string {
	var lines []string

	target := targetLabel(m.req.Target)
	projectLine := fmt.Sprintf("%s %s",
		m.labelStyle.Render("Project:"),
		m.projectValueStyle.Render(target),
	)
	lines = append(lines, projectLine)

	scopeLine := fmt.Sprintf("%s %s",
		m.labelStyle.Render("Scope:"),
		m.scopeValueStyle.Render(m.req.Scope.String()),
	)
	lines = append(lines, scopeLine)
	lines = append(lines, "")

	if m.req.BackupDestination != "" {
		lines = append(lines, fmt.Sprintf("%s %s",
			m.labelStyle.Render("Backup archive:"),
			m.valueStyle.Render(m.req.BackupDestination),
		))
	}

	if m.req.Scope.IncludesWorkspace() {
		lines = append(lines, "Workspace files will be deleted.")
	} else {
		lines = append(lines, "Workspace files will be kept.")
	}

	lines = append(lines, "This action cannot be undone.")

	return strings.Join(lines, "\n")
}

func paletteValue(palette map[string]string, key, fallback string) string {
	if palette == nil {
		return fallback
	}

	if v := strings.TrimSpace(palette[key]); v != "" {
		return v
	}

	return fallback
}
