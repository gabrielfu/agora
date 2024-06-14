package dialogs

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/tui/messages"
	"github.com/gabrielfu/tipi/tui/styles"
	"github.com/gabrielfu/tipi/tui/views"
)

type TextInputCmdFunc func(string) tea.Cmd

type TextInputDialog struct {
	width         int
	maxWidth      int
	title         []string
	footer        []string
	submitCmdFunc TextInputCmdFunc // func to generate a Cmd that submits the input value
	exitView      views.View
	textInput     textinput.Model
}

func NewTextInputDialog(maxWidth int, title, footer []string, submitCmdFunc TextInputCmdFunc, exitView views.View) TextInputDialog {
	t := textinput.New()
	t.Prompt = ""
	return TextInputDialog{
		width:         maxWidth,
		maxWidth:      maxWidth,
		title:         title,
		footer:        footer,
		submitCmdFunc: submitCmdFunc,
		exitView:      exitView,
		textInput:     t,
	}
}

func (m *TextInputDialog) SetPrompt(prompt string) {
	m.textInput.Prompt = prompt
}

func (m *TextInputDialog) SetValue(value string) {
	m.textInput.SetValue(value)
}

func (m *TextInputDialog) Focus() {
	m.textInput.Focus()
}

func (m TextInputDialog) generateStyle() lipgloss.Style {
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: m.title, Footer: m.footer},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(styles.FocusBorderColor)).
		Width(m.width).
		Padding(0, 1)
}

func (m TextInputDialog) exit() tea.Cmd {
	return messages.ExitDialogCmd(m.exitView)
}

func (m *TextInputDialog) SetWidth(windowWidth int) {
	m.width = min(m.maxWidth, windowWidth-4)
}

func (m *TextInputDialog) SetCmdFunc(cmdFunc TextInputCmdFunc) {
	m.submitCmdFunc = cmdFunc
}

func (m *TextInputDialog) Update(msg tea.Msg) (any, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Batch(m.exit(), m.submitCmdFunc(m.textInput.Value()))
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, m.exit()
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, tea.Batch(textinput.Blink, cmd)
}

func (m *TextInputDialog) View() string {
	return m.generateStyle().Render(m.textInput.View())
}
