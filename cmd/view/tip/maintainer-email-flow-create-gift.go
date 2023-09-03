package tipView

import (
	"fmt"

	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	lightninggifts "github.com/sudonym-btc/zap/service/lightning-gifts"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type maintainerEmailFlowCreateGiftModel struct {
	task.Model
	maintainerModel *maintainerModel
	flowParams      *maintainerEmailFlowParams
}

func initialMaintainerEmailFlowCreateGiftModel(i *maintainerModel, flowParams *maintainerEmailFlowParams) *maintainerEmailFlowCreateGiftModel {
	return &maintainerEmailFlowCreateGiftModel{
		maintainerModel: i,
		flowParams:      flowParams,
		Model:           task.New(task.Model{}),
	}
}

func (m maintainerEmailFlowCreateGiftModel) Job() tea.Cmd {
	return func() tea.Msg {
		content, err := lightninggifts.CreateGift(*m.flowParams.amount, "anonymous", *m.flowParams.text)
		if err != nil {
			return maintainerGiftErrMsg{id: m.GetId(), content: err}
		}

		return maintainerGiftMsg{id: m.GetId(), content: content}
	}
}

func (m maintainerEmailFlowCreateGiftModel) Init() tea.Cmd {
	return nil
}

func (m maintainerEmailFlowCreateGiftModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.(type) {
	case maintainerGiftMsg:
		if msg.(maintainerGiftMsg).id != m.GetId() {
			return m, nil
		}
		m.flowParams.gift = msg.(maintainerGiftMsg).content
		m.flowParams.maintainerPayFlowParams.invoice = &m.flowParams.gift.LightningInvoice.Payreq
		cmds = append(cmds, task.Ready(m.GetId()))

	case maintainerGiftErrMsg:
		if msg.(maintainerGiftErrMsg).id != m.GetId() {
			return m, nil
		}
		err := msg.(maintainerGiftErrMsg).content
		cmds = append(cmds, task.Failed(m.GetId(), &err))
	}
	return m, tea.Batch(cmds...)
}

func (m maintainerEmailFlowCreateGiftModel) View() string {

	if m.Progress.InProgress && m.Progress.Ready == false {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Creating lightning gift...")
	} else if m.Progress.Completed {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.CheckMark.Render(), view.SuccessMessageStyle("Created lightning gift"))
	} else if m.Progress.Failed {
		return lipgloss.JoinHorizontal(lipgloss.Left, view.Cross.Render(), view.ErrorMessageStyle("Failed creating lightning gift: "+fmt.Sprint(*m.Progress.Error)))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, view.PadRight.Render(m.Progress.Spinner.View()), "Waiting to create lightning gift...")

}

type maintainerGiftMsg struct {
	id      int
	content *lightninggifts.LightningGift
}
type maintainerGiftErrMsg struct {
	id      int
	content error
}
