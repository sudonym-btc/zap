package tipView

import (
	"github.com/fiatjaf/go-lnurl"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"
	taskView "github.com/sudonym-btc/zap/cmd/view/task"
	wallet "github.com/sudonym-btc/zap/service"

	tea "github.com/charmbracelet/bubbletea"
)

type maintainerLnFlowModel struct {
	task.Model
	maintainerModel *maintainerModel
}
type maintainerPayFlowParams struct {
	invoice *string
	kind    string
	payment *wallet.PayResponse
}
type maintainerLnFlowParams struct {
	*maintainerPayFlowParams
	*detailConfirmParams
	address *lnurl.LNURLPayParams
}

func initialMaintainerLnFlowModel(i *maintainerModel) *maintainerLnFlowModel {
	str := ""
	amount := i.tipModel.each
	text := i.tipModel.comment
	invoice := &str
	payment := &wallet.PayResponse{}
	payFlowParams := &maintainerPayFlowParams{invoice: invoice, kind: "lightning", payment: payment}
	flowParams := &maintainerLnFlowParams{maintainerPayFlowParams: payFlowParams, detailConfirmParams: &detailConfirmParams{amount: &amount, text: &text}}
	return &maintainerLnFlowModel{
		maintainerModel: i,
		Model: task.New(task.Model{
			Progress: &task.TaskProgress{ShouldFailOnFirstSubtaskFail: true},
			Children: []task.ModelI{
				initialMaintainerLnFlowCheckAddressModel(i, flowParams),
				InitialAmountModel(i, flowParams.detailConfirmParams, false),
				initialMaintainerLnFlowGetInvoiceModel(i, flowParams),
				initialMaintainerLnFlowPayInvoiceModel(i, flowParams.detailConfirmParams, flowParams.maintainerPayFlowParams),
			},
		}),
	}
}

func (m maintainerLnFlowModel) Job() tea.Cmd {
	return nil
}

func (m maintainerLnFlowModel) Init() tea.Cmd {
	return nil
}

func (m maintainerLnFlowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	return m, tea.Batch(cmds...)
}

func (m maintainerLnFlowModel) View() string {
	return taskView.DisplayOnlyDoneOrInProgress(m.Children)
}
