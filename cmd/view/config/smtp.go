package configView

import (
	"net/url"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/slog"
	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/service/config"
	"github.com/sudonym-btc/zap/service/email"
)

type (
	errMsg error
)

const (
	host = iota
	port
	user
	pass
)

type EmailModel struct {
	inputs         []textinput.Model
	textarea       textarea.Model
	focused        int
	loading        bool
	loadingSpinner spinner.Model
	err            error
	timer          timer.Model
	quitting       bool
}

const timeout = time.Second * 5

func InitialEmailModel() EmailModel {
	var inputs []textinput.Model = make([]textinput.Model, 4)
	inputs[host] = textinput.New()
	inputs[host].Placeholder = "smtp.gmail.com"
	inputs[host].Focus()
	inputs[host].Prompt = ""

	inputs[port] = textinput.New()
	inputs[port].Placeholder = "587"
	inputs[port].CharLimit = 5
	inputs[port].Prompt = ""

	inputs[user] = textinput.New()
	inputs[user].Placeholder = "joe.smith@gmail.com"
	inputs[user].Prompt = ""

	inputs[pass] = textinput.New()
	inputs[pass].Placeholder = "pswd1234"
	inputs[pass].Prompt = ""

	textarea := textarea.New()
	config, _ := config.LoadConfig()
	// todo: why does setvalue persist here, but not in init function?
	textarea.SetValue(config.EmailTemplate)
	textarea.Placeholder = "Hi there. I'm tipping you a few sats for your work in open source development! Gift attached below."
	timer := timer.New(timeout)

	return EmailModel{
		loading:        false,
		loadingSpinner: spinner.New(),
		inputs:         inputs,
		textarea:       textarea,
		timer:          timer,
		focused:        0,
		err:            nil,
	}
}

func (m EmailModel) Init() tea.Cmd {
	config, _ := config.LoadConfig()
	if config.Smtp != "" {
		parsedUrl, _ := url.Parse(config.Smtp)
		password, _ := parsedUrl.User.Password()
		m.inputs[host].SetValue(parsedUrl.Hostname())
		m.inputs[port].SetValue(parsedUrl.Port())
		m.inputs[user].SetValue(parsedUrl.User.Username())
		m.inputs[pass].SetValue(password)
	}

	return tea.Batch(textinput.Blink, textarea.Blink, m.loadingSpinner.Tick, func() tea.Msg {
		m.textarea.SetValue(config.EmailTemplate)
		return nil
	})
}

func (m EmailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd = []tea.Cmd{}
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
	case testEmailMsg:
		return m, m.Save
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			m.loading = true
			m.err = nil
			return m, m.Test
		case tea.KeyEnter:
			if m.focused != len(m.inputs) {
				m.nextInput()
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.textarea.Blur()

		if m.focused == len(m.inputs) {
			m.textarea.Focus()
		} else {
			m.inputs[m.focused].Focus()
		}

	// We handle errors just like any other message
	case ErrMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	}

	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m EmailModel) View() string {
	var content = []string{lipgloss.JoinVertical(lipgloss.Left,
		view.GreenMessageStyle("Host"),
		m.inputs[host].View(),
		"",

		view.GreenMessageStyle("Port"),
		m.inputs[port].View(),
		"",

		view.GreenMessageStyle("User"),
		m.inputs[user].View(),
		"",

		view.GreenMessageStyle("Pass"),
		m.inputs[pass].View(),
		"",

		view.GreenMessageStyle("Default email message"),

		m.textarea.View(),
	)}
	if m.quitting {
		content = append(content, view.PadVertical.Render(view.Faint.Render("Successfully saved. Exiting in "+m.timer.View())))
	} else {
		if !m.loading {
			content = append(content, view.Faint.Render(view.PadVertical.Render("(ctrl+s) test and save email settings")))
		}
	}
	if m.err != nil {
		content = append(content, view.PadVertical.Render(view.ErrorMessageStyle(m.err.Error())))
	}
	if m.loading {
		content = append(content, view.PadVertical.Render(lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.loadingSpinner.View()), "Loading...")))
	}
	return view.AppStyle.Render(lipgloss.JoinVertical(lipgloss.Left, content...))
}

func (m EmailModel) Save() tea.Msg {
	c, _ := config.LoadConfig()
	c.Smtp = m.mergeUrl()
	c.EmailTemplate = m.textarea.Value()
	err := config.SetConfig(*c)
	slog.Error(c.EmailTemplate)
	if err != nil {
		return ErrMsg{err}
	}
	return saveMsg{}
}

func (m EmailModel) mergeUrl() string {
	url := "smtp://" + m.inputs[user].Value() + ":" + m.inputs[pass].Value() + "@" + m.inputs[host].Value() + ":" + m.inputs[port].Value()
	return url
}

func (m EmailModel) Test() tea.Msg {

	url := m.mergeUrl()

	_, err := email.Connect(url)
	if err != nil {
		return ErrMsg{err}
	}
	return testEmailMsg{str: url}
}

type testEmailMsg struct {
	str string
}

// nextInput focuses the next input field
func (m *EmailModel) nextInput() {
	m.focused = (m.focused + 1) % (len(m.inputs) + 1)
}

// prevInput focuses the previous input field
func (m *EmailModel) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs)
	}
}
