package configView

import (
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sudonym-btc/zap/cmd/view"
	wallet "github.com/sudonym-btc/zap/service"
	"github.com/sudonym-btc/zap/service/config"
)

type WalletConnectModel struct {
	loading        bool
	quitting       bool
	timer          timer.Model
	loadingSpinner spinner.Model
	inputs         []textinput.Model
	err            error
	focused        int
}

func InitialConnectModel() WalletConnectModel {

	input := textinput.New()
	input.Placeholder = "nostr+walletconnect://554534h543543i53h5o34h543oi53534?relay=wss://relay.getalby.com/v1&secret=supersecretsecret&lud16=frostysun783@getalby.com"
	input.Focus()

	amountInput := textinput.New()
	amountInput.Placeholder = "100"

	timer := timer.New(timeout)
	return WalletConnectModel{
		loading:        false,
		timer:          timer,
		loadingSpinner: spinner.New(),
		inputs:         []textinput.Model{input, amountInput},
		focused:        0,
	}
}

func (m WalletConnectModel) Init() tea.Cmd {

	c, _ := config.LoadConfig()

	if c != nil && c.WalletConnect != "" {
		m.inputs[0].SetValue(c.WalletConnect)
	}
	if c != nil && c.DefaultEach != nil {
		m.inputs[1].SetValue(strconv.Itoa(*c.DefaultEach))
	}

	return tea.Batch(m.loadingSpinner.Tick)
}

func (m WalletConnectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
		cmds = append(cmds, cmd)
	case saveMsg:
		m.loading = false
		m.quitting = true
		return m, m.timer.Init()
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		return m, tea.Quit
	case testMsg:
		return m, m.Save
	case ErrMsg:
		m.loading = false
		m.err = msg.err
		return m, m.Save
	case tea.KeyMsg:
		switch msg.String() {
		case "shift+tab":
			m.prevInput()
		case "tab", "enter":
			m.nextInput()
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "ctrl+s":
			m.loading = true
			m.err = nil
			return m, m.Test

		}

	}

	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		if m.focused == i {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m WalletConnectModel) View() string {
	var content = []string{view.GreenMessageStyle("Wallet Connect URL"), view.PadVertical.Render(m.inputs[0].View()),
		view.GreenMessageStyle("Default tip amount"), view.PadVertical.Render(m.inputs[1].View())}
	if m.quitting {
		content = append(content, view.PadVertical.Render(view.Faint.Render("Successfully saved. Exiting in "+m.timer.View())))
	} else {
		if !m.loading {
			content = append(content, view.Faint.Render(view.PadVertical.Render("(ctrl+s) test and save wallet connection")))
		}
	}
	if m.err != nil {
		content = append(content, view.PadVertical.Render(view.ErrorMessageStyle(m.err.Error())))
	}
	if m.loading {
		content = append(content, view.PadVertical.Render(lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.loadingSpinner.View()), "Loading...")))
	}
	return view.AppStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		content...,
	))

}

func (m WalletConnectModel) Save() tea.Msg {
	c, _ := config.LoadConfig()
	c.WalletConnect = m.inputs[0].Value()
	val, err := strconv.Atoi(m.inputs[1].Value())
	c.DefaultEach = &val
	_, errConnect := wallet.Parse_and_connect(c.WalletConnect)
	if errConnect != nil {
		return ErrMsg{err: errConnect}
	}
	err = config.SetConfig(*c)
	if err != nil {
		return ErrMsg{err}
	}
	return saveMsg{}
}

func (m WalletConnectModel) Test() tea.Msg {

	wc, err := wallet.Parse_and_connect(m.inputs[0].Value())
	if err != nil {
		return ErrMsg{err}
	}
	return testMsg{wc}
}

type saveMsg struct {
}
type testMsg struct {
	wc *wallet.WalletConnect
}
type ErrMsg struct {
	err error
}

func (m *WalletConnectModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// nextInput focuses the next input field
func (m *WalletConnectModel) nextInput() {
	m.focused = (m.focused + 1) % (len(m.inputs) + 1)
}

// prevInput focuses the previous input field
func (m *WalletConnectModel) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs)
	}
}
