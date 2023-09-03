package configView

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sudonym-btc/zap/cmd/view"
	wallet "github.com/sudonym-btc/zap/service"
	"github.com/sudonym-btc/zap/service/config"
)

type EmailTemplate struct {
	loading        bool
	loadingSpinner spinner.Model
	inputs         []textinput.Model
	err            error
}

func InitialEmailTemplateModel() EmailTemplate {

	input := textinput.New()
	input.Placeholder = "Enter your wallet connect URI"
	input.Focus()
	return EmailTemplate{
		loading:        false,
		loadingSpinner: spinner.New(),
		inputs:         []textinput.Model{input},
	}
}

func (m EmailTemplate) Init() tea.Cmd {

	c, _ := config.LoadConfig()

	if c != nil && c.WalletConnect != "" {
		m.inputs[0].SetValue(c.WalletConnect)
	}

	return tea.Batch(m.loadingSpinner.Tick)
}

func (m EmailTemplate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
		cmds = append(cmds, cmd)
	case saveMsg:
		m.loading = false
		return m, tea.Quit
	case testMsg:
		return m, m.Save
	case ErrMsg:
		m.loading = false
		m.err = msg.err
		return m, m.Save
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" {
				m.loading = true
				m.err = nil
				return m, m.Test
			}
			return m, nil
		}

	}

	var cmd tea.Cmd
	m.inputs[0], cmd = m.inputs[0].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m EmailTemplate) View() string {
	var arg = []string{view.PadVertical.Render(m.inputs[0].View())}
	if m.loading {
		arg = append(arg, lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.loadingSpinner.View()), "Loading.."))
	}
	if m.err != nil {
		arg = append(arg, lipgloss.JoinHorizontal(lipgloss.Left, view.Cross.Render(), view.ErrorMessageStyle(m.err.Error())))
	}
	return view.AppStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		arg...,
	))

}

func (m EmailTemplate) Save() tea.Msg {
	c, _ := config.LoadConfig()
	c.WalletConnect = m.inputs[0].Value()
	err := config.SetConfig(*c)
	if err != nil {
		return ErrMsg{err}
	}
	return saveMsg{}
}

func (m EmailTemplate) Test() tea.Msg {

	wc, err := wallet.Parse_and_connect(m.inputs[0].Value())
	if err != nil {
		return ErrMsg{err}
	}
	return testMsg{wc}
}

func (m *EmailTemplate) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
