package tipView

import (
	"github.com/fiatjaf/go-lnurl"
	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	"github.com/sudonym-btc/zap/service/address"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type maintainerLnFlowCheckAddressModel struct {
	task.Model
	maintainerModel *maintainerModel
	flowParams      *maintainerLnFlowParams
}

func initialMaintainerLnFlowCheckAddressModel(i *maintainerModel, flowParams *maintainerLnFlowParams) *maintainerLnFlowCheckAddressModel {
	return &maintainerLnFlowCheckAddressModel{
		maintainerModel: i,
		flowParams:      flowParams,
		Model:           task.New(task.Model{}),
	}
}

func (m maintainerLnFlowCheckAddressModel) Job() tea.Cmd {
	return func() tea.Msg {
		params, err := address.Fetch(m.maintainerModel.address)
		if err != nil {
			return maintainerAddressErrMsg{id: m.GetId(), content: err}
		}

		return maintainerAddressMsg{id: m.GetId(), content: params}
	}
}

func (m maintainerLnFlowCheckAddressModel) Init() tea.Cmd {
	return nil
}

func (m maintainerLnFlowCheckAddressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.(type) {
	case maintainerAddressMsg:
		if msg.(maintainerAddressMsg).id != m.GetId() {
			return m, nil
		}
		m.flowParams.address = msg.(maintainerAddressMsg).content
		cmds = append(cmds, task.Ready(m.GetId()))
	case maintainerAddressErrMsg:
		if msg.(maintainerAddressErrMsg).id != m.GetId() {
			return m, nil
		}
		err := msg.(maintainerAddressErrMsg).content
		cmds = append(cmds, task.Failed(m.GetId(), &err))
	}
	return m, tea.Batch(cmds...)
}

func (m maintainerLnFlowCheckAddressModel) View() string {

	if m.Progress.InProgress && m.Progress.Ready == false {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Loading payment address...")
	} else if m.Progress.Completed {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.CheckMark.Render(), view.Faint.Render(view.SuccessMessageStyle("Got address information")))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, view.Cross.Render(), view.ErrorMessageStyle("Could not get lightning information"))
}

type maintainerAddressMsg struct {
	id      int
	content *lnurl.LNURLPayParams
}
type maintainerAddressErrMsg struct {
	id      int
	content error
}
