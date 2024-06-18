package dialogs

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/tui/messages"
	"github.com/gabrielfu/tipi/tui/styles"
	"github.com/gabrielfu/tipi/tui/views"
)

type TextAreaCmdFunc func(string) tea.Cmd

type TextAreaDialog struct {
	width         int
	maxWidth      int
	height        int
	maxHeight     int
	title         []string
	footer        []string
	submitCmdFunc TextAreaCmdFunc // func to generate a Cmd that submits the input value
	exitView      views.View
	textArea      textarea.Model
}

func NewTextAreaDialog(maxWidth, maxHeight int, title, footer []string, submitCmdFunc TextAreaCmdFunc, exitView views.View) TextAreaDialog {
	t := textarea.New()
	t.Prompt = ""
	t.FocusedStyle.LineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	t.BlurredStyle.LineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	return TextAreaDialog{
		width:         maxWidth,
		maxWidth:      maxWidth,
		height:        maxHeight,
		maxHeight:     maxHeight,
		title:         title,
		footer:        footer,
		submitCmdFunc: submitCmdFunc,
		exitView:      exitView,
		textArea:      t,
	}
}

func (m *TextAreaDialog) SetValue(value string) {
	m.textArea.SetValue(value)
}

func (m *TextAreaDialog) Focus() {
	m.textArea.Focus()
}

func (m TextAreaDialog) generateStyle() lipgloss.Style {
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: m.title, Footer: m.footer},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(styles.FocusBorderColor)).
		Width(m.width).
		Height(m.height).
		Padding(0, 1)
}

func (m TextAreaDialog) exit() tea.Cmd {
	return messages.ExitDialogCmd(m.exitView)
}

func (m *TextAreaDialog) SetWidth(windowWidth int) {
	m.width = min(m.maxWidth, windowWidth-4)
}

func (m *TextAreaDialog) SetHeight(windowHeight int) {
	m.height = min(m.maxHeight, windowHeight-7)
}

func (m *TextAreaDialog) SetCmdFunc(cmdFunc TextAreaCmdFunc) {
	m.submitCmdFunc = cmdFunc
}

func (m *TextAreaDialog) Update(msg tea.Msg) (any, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+w":
			return m, tea.Batch(m.exit(), m.submitCmdFunc(m.textArea.Value()))
		case "tab":
			m.textArea.InsertString("  ")
		case "ctrl+c", "esc":
			return m, m.exit()
		}
	}
	m.textArea, cmd = m.textArea.Update(msg)
	return m, tea.Batch(textarea.Blink, cmd)
}

func (m *TextAreaDialog) View() string {
	return m.generateStyle().Render(m.textArea.View())
}
