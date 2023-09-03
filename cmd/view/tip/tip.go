package tipView

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	wallet "github.com/sudonym-btc/zap/service"
	"github.com/sudonym-btc/zap/service/config"
	"github.com/thoas/go-funk"
)

type TippedMsg struct {
	id         int
	maintainer maintainerModel
	method     string
	amount     int
}
type TipModel struct {
	task.Model
	wc               *wallet.WalletConnect
	timer            timer.Model
	quitting         bool
	cwd              string
	sendEmails       bool
	name             string
	version          string
	total            int
	each             int
	comment          string
	manual           bool
	commandResolvers []CommandResolver
	tips             []TippedMsg
}

type TipModelParams struct {
	CommandResolver *CommandResolver
	Name            *string
	Version         *string
	Cwd             *string
	Amount          *int
	Comment         *string
	Manual          *bool
	SendEmails      *bool
}

const timeout = time.Second * 20

func InitialTipModel(params TipModelParams) TipModel {
	commandResolvers := []CommandResolver{}
	if params.CommandResolver != nil {
		commandResolvers = append(commandResolvers, *params.CommandResolver)
	} else {
		commandResolvers = funk.Filter(PackageCommandResolvers, func(cr CommandResolver) bool {
			return cr.Detect(*params.Cwd)
		}).([]CommandResolver)
	}

	timer := timer.New(timeout)
	each := params.Amount
	if each == nil || *each == 0 {
		c, _ := config.LoadConfig()
		if c != nil && c.DefaultEach != nil {
			each = c.DefaultEach
		}
	}
	if params.Comment == nil || *params.Comment == "" {
		comment := ""
		params.Comment = &comment
	}
	if params.SendEmails == nil {
		sendEmails := false
		conf, _ := config.LoadConfig()

		if conf != nil && conf.Smtp != "" {
			sendEmails = true
		}
		params.SendEmails = &sendEmails
	}
	if params.Manual == nil {
		manual := false
		params.Manual = &manual
	}
	model := TipModel{
		tips:             []TippedMsg{},
		timer:            timer,
		comment:          *params.Comment,
		each:             *each,
		cwd:              *params.Cwd,
		manual:           *params.Manual,
		name:             *params.Name,
		sendEmails:       *params.SendEmails,
		commandResolvers: commandResolvers,
		Model:            task.New(task.Model{}),
	}

	return model
}

func (m TipModel) Job() tea.Cmd {
	return func() tea.Msg { return commandResolverMsg(m.commandResolvers) }
}

func (m TipModel) Init() tea.Cmd {
	return task.Begin(m.GetId())
}

func (m TipModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	_, cmd = task.UpdateTask(&m.Model, &m, msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {

	case TippedMsg:
		m.tips = append(m.tips, msg)
	case commandResolverMsg:
		for _, cr := range m.commandResolvers {
			m.Children = append(m.Children, cr.Model(m))
			cmds = append(cmds, m.Children[len(m.Children)-1].Init())
		}
		cmds = append(cmds, task.Ready(m.GetId()))

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

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}

	case task.CompletedMsg:
		if msg.Id == m.GetId() {
			m.quitting = true
			return m, m.timer.Init()
		}
	}

	return m, tea.Batch(cmds...)
}
func (m TipModel) View() string {

	return view.AppStyle.Render(func() string {

		return lipgloss.JoinVertical(lipgloss.Left, func() string {
			if m.Progress.InProgress && m.Progress.Ready == false {
				return lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Checking for package managers...")
			} else if len(m.Children) > 0 {
				return lipgloss.JoinVertical(lipgloss.Left, funk.Map(m.Children, func(l tea.Model) string {
					return l.View()
				}).([]string)...)
			} else {
				return "No supported package manager detected in the working directory."
			}
		}(), m.GetTotals())

	}())

}

func (m TipModel) GetTotals() string {
	if m.quitting {
		return view.PadVertical.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				"Tipped "+view.GreenMessageStyle(fmt.Sprint(funk.Reduce(m.tips, func(acc int, i TippedMsg) int {
					if i.method == "lightning" {
						return acc + 1
					}
					return acc
				}, 0)))+" maintainers "+view.GreenMessageStyle(fmt.Sprint(funk.Reduce(m.tips, func(acc int, i TippedMsg) int {
					if i.method == "lightning" {
						return acc + i.amount
					}
					return acc
				}, 0))+" sats")+" direct to lightning addresses",
				"Tipped "+view.GreenMessageStyle(fmt.Sprint(funk.Reduce(m.tips, func(acc int, i TippedMsg) int {
					if i.method == "email" {
						return acc + 1
					}
					return acc
				}, 0)))+" maintainers "+view.GreenMessageStyle(fmt.Sprint(funk.Reduce(m.tips, func(acc int, i TippedMsg) int {
					if i.method == "email" {
						return acc + i.amount
					}
					return acc
				}, 0))+" sats")+" via email gifts",
				view.PadVertical.Render(view.Faint.Render("Exiting in "+m.timer.View()))))
	}
	return ""
}

type commandResolverMsg []CommandResolver

type QuitMsg struct{}

func Quit() tea.Msg {
	return QuitMsg{}
}
