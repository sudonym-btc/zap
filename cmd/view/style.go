package view

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	DefaultIndent uint = 0
	Orange             = lipgloss.Color("202")
	Red                = lipgloss.Color("9")
	AppStyle           = lipgloss.NewStyle().Padding(1, 2)
	PadRight           = lipgloss.NewStyle().Padding(0, 1)
	PadVertical        = lipgloss.NewStyle().Padding(1, 0)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1).
			Margin(1, 0)
	BoldStyle = lipgloss.NewStyle().
			Bold(true).
			PaddingRight(1)

	GreenMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")).
				Render
	ErrorMessageStyle = lipgloss.NewStyle().
				Foreground(Red).
				Render
	SuccessMessageStyle = lipgloss.NewStyle().
				Foreground(Orange).
				Render
	Faint     = lipgloss.NewStyle().Faint(true)
	CheckMark = lipgloss.NewStyle().Foreground(Orange).MarginRight(1).SetString("✓")
	Cross     = lipgloss.NewStyle().Foreground(Red).MarginRight(1).SetString("✕")
)

func ItemStyle(title, subtitle string) string {
	return lipgloss.NewStyle().MarginTop(1).MarginBottom(1).Render(BoldStyle.Render(title) + " " + subtitle)
}
