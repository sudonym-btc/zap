package tipView

import (
	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	taskView "github.com/sudonym-btc/zap/cmd/view/task"
	wallet "github.com/sudonym-btc/zap/service"
	"github.com/sudonym-btc/zap/service/config"
	lightninggifts "github.com/sudonym-btc/zap/service/lightning-gifts"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type maintainerEmailFlowModel struct {
	task.Model
	maintainerModel *maintainerModel
}

type maintainerEmailFlowParams struct {
	*detailConfirmParams
	*maintainerPayFlowParams
	gift     *lightninggifts.LightningGift
	sendInfo string
}

func initialMaintainerEmailFlowModel(i *maintainerModel) *maintainerEmailFlowModel {
	str := ""
	invoice := &str
	payment := &wallet.PayResponse{}
	amount := i.tipModel.each
	text := i.tipModel.comment
	if text == "" {
		conf, _ := config.LoadConfig()
		if conf != nil {
			text = conf.EmailTemplate
		}
	}
	payFlowParams := &maintainerPayFlowParams{invoice: invoice, payment: payment}
	flowParams := &maintainerEmailFlowParams{maintainerPayFlowParams: payFlowParams, detailConfirmParams: &detailConfirmParams{amount: &amount, text: &text}}
	return &maintainerEmailFlowModel{
		maintainerModel: i,
		Model: task.New(task.Model{
			Progress: &task.TaskProgress{ShouldFailOnFirstSubtaskFail: true},
			Children: []task.ModelI{
				InitialAmountModel(i, flowParams.detailConfirmParams, true),
				initialMaintainerEmailFlowCreateGiftModel(i, flowParams),
				initialMaintainerLnFlowPayInvoiceModel(i, flowParams.detailConfirmParams, flowParams.maintainerPayFlowParams),
				initialMaintainerEmailFlowSendEmailModel(i, flowParams),
			},
		}),
	}
}

func (m maintainerEmailFlowModel) Job() tea.Cmd {
	return nil
}

func (m maintainerEmailFlowModel) Init() tea.Cmd {
	return nil
}

func (m maintainerEmailFlowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	return m, tea.Batch(cmds...)
}

func (m maintainerEmailFlowModel) View() string {

	content := []string{}

	if !m.Progress.Completed && m.Progress.InProgress {
		content = append(content, lipgloss.JoinHorizontal(lipgloss.Left,
			view.PadRight.Render(m.Progress.Spinner.View()),
			"Gifting via email..."))
	}

	content = append(content, taskView.DisplayOnlyDoneOrInProgress(m.Children))

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}
