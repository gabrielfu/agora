package dialogs

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/tui/messages"
	"github.com/gabrielfu/tipi/tui/styles"
	"github.com/gabrielfu/tipi/tui/views"
)

type DoubleTextInputCmdFunc func(string, string) tea.Cmd

type DoubleTextInputDialog struct {
	width          int
	maxWidth       int
	upperTitle     []string
	upperFooter    []string
	lowerTitle     []string
	lowerFooter    []string
	submitCmdFunc  DoubleTextInputCmdFunc // func to generate a Cmd that submits the input values
	exitView       views.View
	upperTextInput textinput.Model
	lowerTextInput textinput.Model
	focusUpper     bool
}

func NewDoubleTextInputDialog(maxWidth int, upperTitle, upperFooter, lowerTitle, lowerFooter []string, submitCmdFunc DoubleTextInputCmdFunc, exitView views.View) DoubleTextInputDialog {
	upper := textinput.New()
	upper.Prompt = ""
	lower := textinput.New()
	lower.Prompt = ""
	return DoubleTextInputDialog{
		width:          maxWidth,
		maxWidth:       maxWidth,
		upperTitle:     upperTitle,
		upperFooter:    upperFooter,
		lowerTitle:     lowerTitle,
		lowerFooter:    lowerFooter,
		submitCmdFunc:  submitCmdFunc,
		exitView:       exitView,
		upperTextInput: upper,
		lowerTextInput: lower,
		focusUpper:     true,
	}
}

func (m *DoubleTextInputDialog) SetUpperPrompt(prompt string) {
	m.upperTextInput.Prompt = prompt
}

func (m *DoubleTextInputDialog) SetLowerPrompt(prompt string) {
	m.lowerTextInput.Prompt = prompt
}

func (m *DoubleTextInputDialog) SetUpperValue(value string) {
	m.upperTextInput.SetValue(value)
}

func (m *DoubleTextInputDialog) SetLowerValue(value string) {
	m.lowerTextInput.SetValue(value)
}

func (m *DoubleTextInputDialog) FocusUpper() {
	m.lowerTextInput.Blur()
	m.upperTextInput.Focus()
	m.focusUpper = true
}

func (m *DoubleTextInputDialog) FocusLower() {
	m.upperTextInput.Blur()
	m.lowerTextInput.Focus()
	m.focusUpper = false
}

func (m DoubleTextInputDialog) generateStyle(upper bool) lipgloss.Style {
	var title, footer []string
	var color string
	if upper {
		title = m.upperTitle
		footer = m.upperFooter
		color = styles.FocusBorderColor
	} else {
		title = m.lowerTitle
		footer = m.lowerFooter
		color = styles.DefaultBorderColor
	}
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: title, Footer: footer},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(color)).
		Width(m.width).
		Padding(0, 1)
}

func (m DoubleTextInputDialog) exit() tea.Cmd {
	return messages.ExitDialogCmd(m.exitView)
}

func (m *DoubleTextInputDialog) SetWidth(windowWidth int) {
	m.width = min(m.maxWidth, windowWidth-4)
}

func (m *DoubleTextInputDialog) SetCmdFunc(cmdFunc DoubleTextInputCmdFunc) {
	m.submitCmdFunc = cmdFunc
}

func (m *DoubleTextInputDialog) updateUpper(msg tea.Msg) (any, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.FocusLower()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, m.exit()
		}
	}
	m.upperTextInput, cmd = m.upperTextInput.Update(msg)
	return m, tea.Batch(textinput.Blink, cmd)
}

func (m *DoubleTextInputDialog) updateLower(msg tea.Msg) (any, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			upper := m.upperTextInput.Value()
			lower := m.lowerTextInput.Value()
			return m, tea.Batch(m.exit(), m.submitCmdFunc(upper, lower))
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, m.exit()
		}
	}
	m.lowerTextInput, cmd = m.lowerTextInput.Update(msg)
	return m, tea.Batch(textinput.Blink, cmd)
}

func (m *DoubleTextInputDialog) Update(msg tea.Msg) (any, tea.Cmd) {
	if m.focusUpper {
		return m.updateUpper(msg)
	} else {
		return m.updateLower(msg)
	}
}

func (m *DoubleTextInputDialog) View() string {
	return m.generateStyle(true).Render(m.upperTextInput.View()) +
		"\n" +
		m.generateStyle(false).Render(m.lowerTextInput.View())
}
