package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	ID        string
	Name      string
	Bg        lipgloss.Color
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Danger    lipgloss.Color
	Warning   lipgloss.Color
	Muted     lipgloss.Color
	Text      lipgloss.Color
	Border    lipgloss.Color
	KeyBg     lipgloss.Color
	KeyFg     lipgloss.Color
}

var themes = []Theme{
	{
		ID:        "dracula",
		Name:      "🧛 Dracula",
		Bg:        lipgloss.Color("#282A36"),
		Primary:   lipgloss.Color("#BD93F9"),
		Secondary: lipgloss.Color("#8BE9FD"),
		Success:   lipgloss.Color("#50FA7B"),
		Danger:    lipgloss.Color("#FF5555"),
		Warning:   lipgloss.Color("#FFB86C"),
		Muted:     lipgloss.Color("#6272A4"),
		Text:      lipgloss.Color("#F8F8F2"),
		Border:    lipgloss.Color("#6272A4"),
		KeyBg:     lipgloss.Color("#44475A"),
		KeyFg:     lipgloss.Color("#BD93F9"),
	},
	{
		ID:        "tokyo-night",
		Name:      "🌌 Tokyo Night",
		Bg:        lipgloss.Color("#24283B"),
		Primary:   lipgloss.Color("#7AA2F7"),
		Secondary: lipgloss.Color("#BB9AF7"),
		Success:   lipgloss.Color("#9ECE6A"),
		Danger:    lipgloss.Color("#F7768E"),
		Warning:   lipgloss.Color("#E0AF68"),
		Muted:     lipgloss.Color("#565F89"),
		Text:      lipgloss.Color("#C0CAF5"),
		Border:    lipgloss.Color("#414868"),
		KeyBg:     lipgloss.Color("#343B58"),
		KeyFg:     lipgloss.Color("#7AA2F7"),
	},
	{
		ID:        "catppuccin",
		Name:      "🌸 Catppuccin Mocha",
		Bg:        lipgloss.Color("#1E1E2E"),
		Primary:   lipgloss.Color("#CBA6F7"),
		Secondary: lipgloss.Color("#89B4FA"),
		Success:   lipgloss.Color("#A6E3A1"),
		Danger:    lipgloss.Color("#F38BA8"),
		Warning:   lipgloss.Color("#FAB387"),
		Muted:     lipgloss.Color("#6C7086"),
		Text:      lipgloss.Color("#CDD6F4"),
		Border:    lipgloss.Color("#45475A"),
		KeyBg:     lipgloss.Color("#313244"),
		KeyFg:     lipgloss.Color("#CBA6F7"),
	},
	{
		ID:        "nord",
		Name:      "🌊 Nord Arctic",
		Bg:        lipgloss.Color("#2E3440"),
		Primary:   lipgloss.Color("#88C0D0"),
		Secondary: lipgloss.Color("#81A1C1"),
		Success:   lipgloss.Color("#A3BE8C"),
		Danger:    lipgloss.Color("#BF616A"),
		Warning:   lipgloss.Color("#EBCB8B"),
		Muted:     lipgloss.Color("#616E88"),
		Text:      lipgloss.Color("#D8DEE9"),
		Border:    lipgloss.Color("#4C566A"),
		KeyBg:     lipgloss.Color("#3B4252"),
		KeyFg:     lipgloss.Color("#88C0D0"),
	},
	{
		ID:        "gruvbox",
		Name:      "🍂 Gruvbox Dark",
		Bg:        lipgloss.Color("#282828"),
		Primary:   lipgloss.Color("#FE8019"),
		Secondary: lipgloss.Color("#FABD2F"),
		Success:   lipgloss.Color("#B8BB26"),
		Danger:    lipgloss.Color("#FB4934"),
		Warning:   lipgloss.Color("#D79921"),
		Muted:     lipgloss.Color("#928374"),
		Text:      lipgloss.Color("#EBDBB2"),
		Border:    lipgloss.Color("#504945"),
		KeyBg:     lipgloss.Color("#3C3836"),
		KeyFg:     lipgloss.Color("#FE8019"),
	},
	{
		ID:        "kanagawa",
		Name:      "🌊 Kanagawa Wave",
		Bg:        lipgloss.Color("#1F1F28"),
		Primary:   lipgloss.Color("#7E9CD8"),
		Secondary: lipgloss.Color("#FFA066"),
		Success:   lipgloss.Color("#98BB6C"),
		Danger:    lipgloss.Color("#E46876"),
		Warning:   lipgloss.Color("#DCA561"),
		Muted:     lipgloss.Color("#727169"),
		Text:      lipgloss.Color("#DCD7BA"),
		Border:    lipgloss.Color("#363646"),
		KeyBg:     lipgloss.Color("#2A2A37"),
		KeyFg:     lipgloss.Color("#7E9CD8"),
	},
	{
		ID:        "tron-legacy",
		Name:      "💠 Tron: Legacy",
		Bg:        lipgloss.Color("#11161B"),
		Primary:   lipgloss.Color("#00E5FF"),
		Secondary: lipgloss.Color("#6DFCFF"),
		Success:   lipgloss.Color("#00FF9D"),
		Danger:    lipgloss.Color("#FF3344"),
		Warning:   lipgloss.Color("#FFA000"),
		Muted:     lipgloss.Color("#4E6473"),
		Text:      lipgloss.Color("#DFE8ED"),
		Border:    lipgloss.Color("#1C2D37"),
		KeyBg:     lipgloss.Color("#1A242C"),
		KeyFg:     lipgloss.Color("#00E5FF"),
	},
	{
		ID:        "tron-ares",
		Name:      "🔴 Tron: Ares",
		Bg:        lipgloss.Color("#181214"),
		Primary:   lipgloss.Color("#FF1744"),
		Secondary: lipgloss.Color("#FF6D00"),
		Success:   lipgloss.Color("#00E676"),
		Danger:    lipgloss.Color("#D50000"),
		Warning:   lipgloss.Color("#FFB300"),
		Muted:     lipgloss.Color("#6E5059"),
		Text:      lipgloss.Color("#FBECEF"),
		Border:    lipgloss.Color("#382026"),
		KeyBg:     lipgloss.Color("#281B20"),
		KeyFg:     lipgloss.Color("#FF1744"),
	},
	{
		ID:        "matrix",
		Name:      "📟 Matrix",
		Bg:        lipgloss.Color("#0D1410"),
		Primary:   lipgloss.Color("#00FF66"),
		Secondary: lipgloss.Color("#26A69A"),
		Success:   lipgloss.Color("#00FF66"),
		Danger:    lipgloss.Color("#FF5252"),
		Warning:   lipgloss.Color("#FFD54F"),
		Muted:     lipgloss.Color("#3E6B48"),
		Text:      lipgloss.Color("#C8E6C9"),
		Border:    lipgloss.Color("#1E3626"),
		KeyBg:     lipgloss.Color("#16241C"),
		KeyFg:     lipgloss.Color("#00FF66"),
	},
}

var currentTheme = themes[0]

func padLine(content string, targetWidth int, bg lipgloss.Color) string {
	w := lipgloss.Width(content)
	if w >= targetWidth {
		return content
	}
	spaces := lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", targetWidth-w))
	return content + spaces
}

func formatKeycap(key string, t Theme) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.KeyFg).
		Background(t.KeyBg).
		Padding(0, 1).
		Render(key)
}

func formatAction(key, label string, t Theme) string {
	k := formatKeycap(key, t)
	l := lipgloss.NewStyle().Foreground(t.Text).Background(t.Bg).Render(" " + label)
	return k + l
}

func formatFooterBar(title string, actions []string, t Theme) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary).
		Background(t.Bg).
		Render(title + " ")

	sep := lipgloss.NewStyle().
		Foreground(t.Muted).
		Background(t.Bg).
		Render("  •  ")

	return titleStyle + strings.Join(actions, sep)
}