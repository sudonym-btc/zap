package tipView

import (
	"strconv"

	"github.com/muesli/reflow/indent"
	"github.com/sudonym-btc/zap/cmd/view"
	"github.com/sudonym-btc/zap/cmd/view/helper"
	"github.com/sudonym-btc/zap/cmd/view/helper/task"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type amountInputModel struct {
	task.Model
	maintainerModel     *maintainerModel
	detailConfirmParams *detailConfirmParams
	amountInput         textinput.Model
	textInput           textarea.Model
	isEmail             bool
	focused             int
}

func cvvValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func InitialAmountModel(maintainerModel *maintainerModel, detailConfirmParams *detailConfirmParams, isEmail bool) amountInputModel {

	input := textinput.New()
	input.Placeholder = "2100"
	input.Focus()
	input.CharLimit = 20
	input.Width = 30
	input.CharLimit = 8
	input.SetValue(strconv.Itoa(*detailConfirmParams.amount))

	input.Prompt = ""

	textInput := textarea.New()
	textInput.Placeholder = "Thanks for the great work!"
	textInput.SetValue(*detailConfirmParams.text)
	textInput.Prompt = ""
	if detailConfirmParams.textLimit != nil {
		textInput.CharLimit = *detailConfirmParams.textLimit
	}
	// input.Validate = ccnValidator
	return amountInputModel{
		maintainerModel:     maintainerModel,
		amountInput:         input,
		isEmail:             isEmail,
		textInput:           textInput,
		detailConfirmParams: detailConfirmParams,
		Model: task.New(task.Model{
			Id:       helper.Id{Id: 20},
			Progress: &task.TaskProgress{HideOnDone: true},
		}),
	}
}

func (m amountInputModel) Job() tea.Cmd {
	if !m.maintainerModel.tipModel.manual {
		return m.Save()
	}
	return func() tea.Msg {
		return nil
	}

}

func (m amountInputModel) Save() tea.Cmd {
	am, _ := strconv.Atoi(m.amountInput.Value())
	m.detailConfirmParams.amount = &am
	str := m.textInput.Value()
	m.detailConfirmParams.text = &str
	// validate
	return task.Completed(m.GetId())
}

func (m amountInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m amountInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 2)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
			// Only trigger save if current maintainer focused
			if m.Progress.InProgress {
				return m, m.Save()
			}
		case "enter":
			m.nextInput()
		case "shift+tab":
			m.prevInput()
		case "tab":
			m.nextInput()
		}
		m.amountInput.Blur()
		m.textInput.Blur()
		if m.focused == 0 {
			m.amountInput.Focus()
		} else {
			m.textInput.Focus()
		}

		// We handle errors just like any other message
		// case errMsg:
		// 	m.err = msg
		// 	return m, nil
	}

	m.amountInput, cmds[0] = m.amountInput.Update(msg)
	m.textInput, cmds[0] = m.textInput.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m amountInputModel) View() string {
	content := ""
	if m.isEmail {
		content = "## Gift URL attached here"
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		view.PadVertical.Render(""),
		view.GreenMessageStyle("Amount (sats)"),
		view.PadVertical.Render(m.amountInput.View()),
		view.GreenMessageStyle("Message"),
		view.PadVertical.Render(m.textInput.View()),
		indent.String(view.Faint.Render(content), 1),
		view.PadVertical.Render(
			indent.String(view.Faint.Render("(ctrl+n to accept)"), 1),
		),
	)
}

// nextInput focuses the next input field
func (m *amountInputModel) nextInput() {
	m.focused = (m.focused + 1) % 2
}

// prevInput focuses the previous input field
func (m *amountInputModel) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = 2 - 1
	}
}
