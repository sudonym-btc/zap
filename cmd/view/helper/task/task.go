package task

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/slog"
	"github.com/sudonym-btc/zap/cmd/view/helper"
	"github.com/thoas/go-funk"
)

type ModelI interface {
	helper.IdModel
	Job() tea.Cmd
	GetProgress() *TaskProgress
	GetModel() *Model
	ViewThis() string
}

type Model struct {
	helper.Id
	focused     int
	Progress    *TaskProgress
	Children    []ModelI
	OnCompleted func() tea.Cmd
}

func New(m Model) Model {
	if m.Progress == nil {
		m.Progress = &TaskProgress{}
	}
	progress := &TaskProgress{
		ShouldFailOnFirstSubtaskFail:         m.Progress.ShouldFailOnFirstSubtaskFail,
		ShouldCompleteOnFirstSubtaskComplete: m.Progress.ShouldCompleteOnFirstSubtaskComplete,
		HideOnDone:                           m.Progress.HideOnDone,
		Spinner:                              spinner.New(),
	}
	// progress.Spinner.Spinner = spinner.MiniDot

	return Model{
		Id:          helper.NewId(),
		Progress:    progress,
		OnCompleted: m.OnCompleted,
		Children:    m.Children,
	}
}

type TaskProgress struct {
	Spinner                              spinner.Model
	ShouldReady                          bool
	ShouldFailOnFirstSubtaskFail         bool
	ShouldCompleteOnFirstSubtaskComplete bool
	HideOnDone                           bool
	ShouldBegin                          bool
	InProgress                           bool
	Ready                                bool
	Failed                               bool
	Completed                            bool
	Skipped                              bool
	Error                                *error
}

func (m Model) GetModel() *Model {
	return &m
}

func (m *Model) UpdateProgress(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if msg.ID == m.Progress.Spinner.ID() {
			var cmd tea.Cmd
			m.Progress.Spinner, cmd = m.Progress.Spinner.Update(msg)
			return cmd
		}
	case BeginMsg:
		if msg.Id == m.GetId() {
			m.Progress.InProgress = true
			return m.Progress.Spinner.Tick
		}
	case ReadyMsg:
		if msg.Id == m.GetId() {
			m.Progress.Ready = true
		}
	case CompletedMsg:
		if msg.Id == m.GetId() {
			m.Progress.Completed = true
			m.Progress.InProgress = false
		}
	case SkippedMsg:
		if msg.Id == m.GetId() {
			m.Progress.Skipped = true
			m.Progress.InProgress = false
		}
	case FailedMsg:
		if msg.Id == m.GetId() {
			m.Progress.Failed = true
			m.Progress.Error = msg.Error
			m.Progress.InProgress = false
		}
	}
	return nil
}

type BeginMsg struct {
	Id int
}

type ReadyMsg struct {
	Id int
}

type CompletedMsg struct {
	Id int
}

type FailedMsg struct {
	Id    int
	Error *error
}

type SkippedMsg struct {
	Id int
}

func Begin(id int) func() tea.Msg {
	return func() tea.Msg {
		slog.Debug("Begin task", id)
		return BeginMsg{Id: id}
	}
}
func Ready(id int) func() tea.Msg {
	return func() tea.Msg {
		slog.Debug("Ready task", id)
		return ReadyMsg{Id: id}
	}
}
func Skipped(id int) func() tea.Msg {
	return func() tea.Msg {
		slog.Debug("Skipped task", id)
		return SkippedMsg{Id: id}
	}
}
func Completed(id int) func() tea.Msg {
	return func() tea.Msg {
		slog.Debug("Completed task", id)
		return CompletedMsg{Id: id}
	}
}
func Failed(id int, err *error) func() tea.Msg {
	return func() tea.Msg {
		slog.Warn("Failed task", id, err)
		return FailedMsg{Id: id, Error: err}
	}
}

func (m Model) Job() tea.Cmd {
	slog.Error("Job task")
	return nil
}

func (m Model) ViewThis() string {
	return ""
}

func (m Model) GetProgress() *TaskProgress {
	return m.Progress
}

func (m Model) View() string {
	return ""
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// Here we need the modelInterface parameter
// as this will allow our override functions to be accessed from this base class

func UpdateTask(m *Model, mI ModelI, msg tea.Msg) (Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	for i, child := range m.Children {
		// Normal view update
		updated, cmd := child.Update(msg)
		m.Children[i] = updated.(ModelI)
		// Do task-model specific updates
		_, cmd2 := UpdateTask(m.Children[i].GetModel(), m.Children[i], msg)
		cmds = append(cmds, cmd, cmd2)
	}
	cmds = append(cmds, m.UpdateProgress(msg))

	switch msg := msg.(type) {
	case BeginMsg:

		if msg.Id == m.GetId() {
			cmd := mI.Job()
			// If the job is nil, we assume that the task is ready
			if cmd == nil {
				cmd = Ready(m.GetId())
			}
			cmds = append(cmds, cmd)
		}

	case ReadyMsg:
		if msg.Id == m.GetId() {
			cmds = append(cmds, m.Next())
		}

	case CompletedMsg:
		if msg.Id == m.GetId() && m.OnCompleted != nil {
			cmds = append(cmds, m.OnCompleted())
		} else {
			for _, child := range m.Children {
				if msg.Id == child.GetId() {
					if m.Progress.ShouldCompleteOnFirstSubtaskComplete {
						for _, child := range m.Children {
							if !child.GetProgress().Completed {
								child.GetModel().Progress.Skipped = true
							}
						}
						cmds = append(cmds, Completed(m.GetId()))
					} else {
						cmds = append(cmds, m.Next())
					}
				}
			}
		}

	case FailedMsg:

		for _, child := range m.Children {
			if msg.Id == child.GetId() {
				if m.Progress.ShouldFailOnFirstSubtaskFail {
					cmds = append(cmds, Failed(m.GetId(), nil))
				} else {
					cmds = append(cmds, m.Next())
				}
			}
		}
	}
	return *m, tea.Batch(cmds...)
}
func (m Model) Next() tea.Cmd {
	slog.Debug("Next child for ", m.GetId())

	var firstPending = funk.Find(m.Children, func(l ModelI) bool {
		return !l.GetProgress().Completed && !l.GetProgress().Failed && !l.GetProgress().Skipped
	})

	if firstPending != nil {
		return Begin(firstPending.(ModelI).GetId())
	}
	if m.Progress.Completed || m.Progress.Failed {
		return nil
	}

	return Completed(m.GetId())
}
