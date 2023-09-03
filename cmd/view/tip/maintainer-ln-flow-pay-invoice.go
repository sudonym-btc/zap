package tipView

import (
	"fmt"

	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	wallet "github.com/sudonym-btc/zap/service"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type maintainerLnFlowPayInvoiceModel struct {
	task.Model
	maintainerModel         *maintainerModel
	detailConfirmParams     *detailConfirmParams
	maintainerPayFlowParams *maintainerPayFlowParams
}

func initialMaintainerLnFlowPayInvoiceModel(i *maintainerModel, detailConfirmParams *detailConfirmParams, payFlowParams *maintainerPayFlowParams) *maintainerLnFlowPayInvoiceModel {
	return &maintainerLnFlowPayInvoiceModel{
		maintainerModel:         i,
		detailConfirmParams:     detailConfirmParams,
		maintainerPayFlowParams: payFlowParams,
		Model:                   task.New(task.Model{}),
	}
}

func (m maintainerLnFlowPayInvoiceModel) Job() tea.Cmd {

	return func() tea.Msg {
		wc, _ := wallet.Connect()
		content, err := wallet.Pay_invoice(wc, *m.maintainerPayFlowParams.invoice)
		if err != nil {
			return maintainerPaymentErrMsg{id: m.GetId(), content: err}
		}

		return maintainerPaymentMsg{id: m.GetId(), content: content}
	}
}

func (m maintainerLnFlowPayInvoiceModel) Init() tea.Cmd {
	return nil
}

func (m maintainerLnFlowPayInvoiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.(type) {
	case maintainerPaymentMsg:
		if msg.(maintainerPaymentMsg).id != m.GetId() {
			return m, nil
		}
		m.maintainerPayFlowParams.payment = msg.(maintainerPaymentMsg).content
		cmds = append(cmds, task.Ready(m.GetId()), func() tea.Msg {
			if m.maintainerPayFlowParams.kind == "lightning" {
				return TippedMsg{maintainer: *m.maintainerModel, method: "lightning", amount: *m.detailConfirmParams.amount}
			}
			return nil
		})

	case maintainerPaymentErrMsg:
		if msg.(maintainerPaymentErrMsg).id != m.GetId() {
			return m, nil
		}
		err := msg.(maintainerPaymentErrMsg).content
		cmds = append(cmds, task.Failed(m.GetId(), &err))

	}
	return m, tea.Batch(cmds...)
}

func (m maintainerLnFlowPayInvoiceModel) View() string {

	if m.Progress.InProgress && m.Progress.Ready == false {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Paying invoice...")
	} else if m.Progress.Completed {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.CheckMark.Render(), view.SuccessMessageStyle("Paid invoice"))
	} else if m.Progress.Failed {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.Cross.Render(), lipgloss.JoinVertical(lipgloss.Left, view.ErrorMessageStyle("Failed paying invoice"), fmt.Sprint((*m.Progress.Error))))
	}
	return ""
}

type maintainerPaymentMsg struct {
	id      int
	content *wallet.PayResponse
}
type maintainerPaymentErrMsg struct {
	id      int
	content error
}
