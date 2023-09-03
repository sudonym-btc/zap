package tipView

import (
	"fmt"

	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	"github.com/sudonym-btc/zap/service/address"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type maintainerLnFlowGetInvoiceModel struct {
	task.Model
	maintainerModel *maintainerModel
	flowParams      *maintainerLnFlowParams
}

func initialMaintainerLnFlowGetInvoiceModel(i *maintainerModel, flowParams *maintainerLnFlowParams) *maintainerLnFlowGetInvoiceModel {
	return &maintainerLnFlowGetInvoiceModel{
		maintainerModel: i,
		flowParams:      flowParams,
		Model:           task.New(task.Model{}),
	}
}

func (m maintainerLnFlowGetInvoiceModel) Job() tea.Cmd {

	return func() tea.Msg {
		params, err := address.FetchFromParams(int64(*m.flowParams.amount)*1000, *m.flowParams.text, *m.flowParams.address)
		if err != nil {
			return maintainerInvoiceErrMsg{id: m.GetId(), content: err}
		}

		return maintainerInvoiceMsg{id: m.GetId(), content: *params}
	}
}

func (m maintainerLnFlowGetInvoiceModel) Init() tea.Cmd {
	return nil
}

func (m maintainerLnFlowGetInvoiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.(type) {
	case maintainerInvoiceMsg:
		if msg.(maintainerInvoiceMsg).id != m.GetId() {
			return m, nil
		}
		invoice := msg.(maintainerInvoiceMsg).content
		m.flowParams.invoice = &invoice
		cmds = append(cmds, task.Ready(m.GetId()))

	case maintainerInvoiceErrMsg:
		if msg.(maintainerInvoiceErrMsg).id != m.GetId() {
			return m, nil
		}
		err := msg.(maintainerInvoiceErrMsg).content
		cmds = append(cmds, task.Failed(m.GetId(), &err))
	}
	return m, tea.Batch(cmds...)
}

func (m maintainerLnFlowGetInvoiceModel) View() string {

	if m.Progress.InProgress && m.Progress.Ready == false {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Getting invoice...")
	} else if m.Progress.Completed {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.CheckMark.Render(), view.Faint.Render(view.SuccessMessageStyle("Got invoice")))
	} else if m.Progress.Failed {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.Cross.Render(), lipgloss.JoinVertical(lipgloss.Left, view.ErrorMessageStyle("Failed getting invoice"), fmt.Sprint((*m.Progress.Error))))
	}
	return ""
}

type maintainerInvoiceMsg struct {
	id      int
	content string
}
type maintainerInvoiceErrMsg struct {
	id      int
	content error
}
