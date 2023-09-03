package tipView

import (
	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	"github.com/sudonym-btc/zap/service/email"
	lightninggifts "github.com/sudonym-btc/zap/service/lightning-gifts"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type maintainerEmailFlowSendEmailModel struct {
	task.Model
	maintainerModel *maintainerModel
	flowParams      *maintainerEmailFlowParams
}

func initialMaintainerEmailFlowSendEmailModel(i *maintainerModel, flowParams *maintainerEmailFlowParams) *maintainerEmailFlowSendEmailModel {
	return &maintainerEmailFlowSendEmailModel{
		maintainerModel: i,
		flowParams:      flowParams,
		Model:           task.New(task.Model{}),
	}
}

func (m maintainerEmailFlowSendEmailModel) Job() tea.Cmd {

	return func() tea.Msg {
		err := email.SendMail(m.maintainerModel.address, "Donation", *m.flowParams.text+"\n\n"+lightninggifts.RedeemUrl(*&m.flowParams.gift.OrderId))
		if err != nil {
			return maintainerEmailErrMsg{id: m.GetId(), content: err}
		}

		return maintainerEmailMsg{id: m.GetId(), content: m.flowParams.text}
	}
}

func (m maintainerEmailFlowSendEmailModel) Init() tea.Cmd {
	return nil
}

func (m maintainerEmailFlowSendEmailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.(type) {
	case maintainerEmailMsg:
		if msg.(maintainerEmailMsg).id != m.GetId() {
			return m, nil
		}
		m.flowParams.sendInfo = *msg.(maintainerEmailMsg).content
		cmds = append(cmds, task.Ready(m.GetId()), func() tea.Msg {
			return TippedMsg{maintainer: *m.maintainerModel, method: "email", amount: *m.flowParams.amount}
		})

	case maintainerEmailErrMsg:
		if msg.(maintainerEmailErrMsg).id != m.GetId() {
			return m, nil
		}
		err := msg.(maintainerEmailErrMsg).content
		cmds = append(cmds, task.Failed(m.GetId(), &err))
	}
	return m, tea.Batch(cmds...)
}

func (m maintainerEmailFlowSendEmailModel) View() string {

	if m.Progress.InProgress && m.Progress.Ready == false {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Sending email...")
	} else if m.Progress.Completed {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.CheckMark.Render(), view.SuccessMessageStyle("Sent email"))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, view.Cross.Render(), view.ErrorMessageStyle("Failed sending email"))
}

type maintainerEmailMsg struct {
	id      int
	content *string
}
type maintainerEmailErrMsg struct {
	id      int
	content error
}
