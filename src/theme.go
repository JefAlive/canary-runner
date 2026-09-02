package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	ID        string
	Name      string
	Bg        lipgloss.Color // Fundo oficial com subtom sutil
	Primary   lipgloss.Color // Destaques / Títulos / Cursor
	Secondary lipgloss.Color // Ações / Menus secundários
	Success   lipgloss.Color // Status ON / Concluído
	Danger    lipgloss.Color // Kill / Erro / Alerta Crítico
	Warning   lipgloss.Color // Avisos / Não compilado
	Muted     lipgloss.Color // Descrições / Comentários
	Text      lipgloss.Color // Texto principal
	Border    lipgloss.Color // Bordas de caixas
	KeyBg     lipgloss.Color // Fundo de keycaps e seleção
	KeyFg     lipgloss.Color // Texto das keycaps
}

var themes = []Theme{
	{
		ID:        "dracula",
		Name:      "🧛 Dracula",
		Bg:        lipgloss.Color("#282A36"), // Grafite arroxeado clássico
		Primary:   lipgloss.Color("#BD93F9"), // Roxo vibrante
		Secondary: lipgloss.Color("#8BE9FD"), // Ciano
		Success:   lipgloss.Color("#50FA7B"), // Verde
		Danger:    lipgloss.Color("#FF5555"), // Vermelho
		Warning:   lipgloss.Color("#FFB86C"), // Laranja
		Muted:     lipgloss.Color("#6272A4"), // Comentário cinza-azulado
		Text:      lipgloss.Color("#F8F8F2"), // Branco suave
		Border:    lipgloss.Color("#6272A4"),
		KeyBg:     lipgloss.Color("#44475A"),
		KeyFg:     lipgloss.Color("#BD93F9"),
	},
	{
		ID:        "tokyo-night",
		Name:      "🌌 Tokyo Night",
		Bg:        lipgloss.Color("#24283B"), // Azul-marinho noite (Storm)
		Primary:   lipgloss.Color("#7AA2F7"), // Azul elétrico
		Secondary: lipgloss.Color("#BB9AF7"), // Magenta suave
		Success:   lipgloss.Color("#9ECE6A"), // Verde oliva claro
		Danger:    lipgloss.Color("#F7768E"), // Coral avermelhado
		Warning:   lipgloss.Color("#E0AF68"), // Âmbar
		Muted:     lipgloss.Color("#565F89"), // Cinza escuro frio
		Text:      lipgloss.Color("#C0CAF5"), // Azul esbranquiçado
		Border:    lipgloss.Color("#414868"),
		KeyBg:     lipgloss.Color("#343B58"),
		KeyFg:     lipgloss.Color("#7AA2F7"),
	},
	{
		ID:        "catppuccin",
		Name:      "🌸 Catppuccin Mocha",
		Bg:        lipgloss.Color("#1E1E2E"), // Grafite aveludado mirtilo
		Primary:   lipgloss.Color("#CBA6F7"), // Mauve pastel
		Secondary: lipgloss.Color("#89B4FA"), // Sapphire
		Success:   lipgloss.Color("#A6E3A1"), // Menta
		Danger:    lipgloss.Color("#F38BA8"), // Flamingo
		Warning:   lipgloss.Color("#FAB387"), // Pêssego
		Muted:     lipgloss.Color("#6C7086"), // Overlay
		Text:      lipgloss.Color("#CDD6F4"), // Texto claro
		Border:    lipgloss.Color("#45475A"),
		KeyBg:     lipgloss.Color("#313244"),
		KeyFg:     lipgloss.Color("#CBA6F7"),
	},
	{
		ID:        "nord",
		Name:      "🌊 Nord Arctic",
		Bg:        lipgloss.Color("#2E3440"), // Cinza-ardósia glacial (Polar Night)
		Primary:   lipgloss.Color("#88C0D0"), // Frost ciano
		Secondary: lipgloss.Color("#81A1C1"), // Frost azul
		Success:   lipgloss.Color("#A3BE8C"), // Aurora verde
		Danger:    lipgloss.Color("#BF616A"), // Aurora vermelho
		Warning:   lipgloss.Color("#EBCB8B"), // Aurora amarelo
		Muted:     lipgloss.Color("#616E88"), // Ardósia médio
		Text:      lipgloss.Color("#D8DEE9"), // Neve suave
		Border:    lipgloss.Color("#4C566A"),
		KeyBg:     lipgloss.Color("#3B4252"),
		KeyFg:     lipgloss.Color("#88C0D0"),
	},
	{
		ID:        "gruvbox",
		Name:      "🍂 Gruvbox Dark",
		Bg:        lipgloss.Color("#282828"), // Carvão quente com toque terra
		Primary:   lipgloss.Color("#FE8019"), // Laranja queimado
		Secondary: lipgloss.Color("#FABD2F"), // Amarelo mostarda
		Success:   lipgloss.Color("#B8BB26"), // Verde oliva
		Danger:    lipgloss.Color("#FB4934"), // Vermelho terracota
		Warning:   lipgloss.Color("#D79921"), // Âmbar quente
		Muted:     lipgloss.Color("#928374"), // Cinza terra
		Text:      lipgloss.Color("#EBDBB2"), // Pergaminho claro
		Border:    lipgloss.Color("#504945"),
		KeyBg:     lipgloss.Color("#3C3836"),
		KeyFg:     lipgloss.Color("#FE8019"),
	},
	{
		ID:        "kanagawa",
		Name:      "🌊 Kanagawa Wave",
		Bg:        lipgloss.Color("#1F1F28"), // Nanquim japonês (Sumi Ink)
		Primary:   lipgloss.Color("#7E9CD8"), // Azul cristalino
		Secondary: lipgloss.Color("#FFA066"), // Laranja surimi
		Success:   lipgloss.Color("#98BB6C"), // Verde primavera
		Danger:    lipgloss.Color("#E46876"), // Vermelho outono
		Warning:   lipgloss.Color("#DCA561"), // Ouro antigo
		Muted:     lipgloss.Color("#727169"), // Cinza fuji
		Text:      lipgloss.Color("#DCD7BA"), // Branco fuji
		Border:    lipgloss.Color("#363646"),
		KeyBg:     lipgloss.Color("#2A2A37"),
		KeyFg:     lipgloss.Color("#7E9CD8"),
	},
	{
		ID:        "tron-legacy",
		Name:      "💠 Tron: Legacy",
		Bg:        lipgloss.Color("#11161B"), // Aço escuro da Grade
		Primary:   lipgloss.Color("#00E5FF"), // Azul/Ciano néon icônico
		Secondary: lipgloss.Color("#6DFCFF"), // Ciano claro
		Success:   lipgloss.Color("#00FF9D"), // Verde circuito
		Danger:    lipgloss.Color("#FF3344"), // Vermelho do Rinzler
		Warning:   lipgloss.Color("#FFA000"), // Laranja do CLU
		Muted:     lipgloss.Color("#4E6473"), // Aço desbotado
		Text:      lipgloss.Color("#DFE8ED"), // Branco holográfico
		Border:    lipgloss.Color("#1C2D37"),
		KeyBg:     lipgloss.Color("#1A242C"),
		KeyFg:     lipgloss.Color("#00E5FF"),
	},
	{
		ID:        "tron-ares",
		Name:      "🔴 Tron: Ares",
		Bg:        lipgloss.Color("#181214"), // Obsidiana com leve tom carmesim
		Primary:   lipgloss.Color("#FF1744"), // Vermelho néon do Ares
		Secondary: lipgloss.Color("#FF6D00"), // Âmbar Dillinger
		Success:   lipgloss.Color("#00E676"), // Verde néon
		Danger:    lipgloss.Color("#D50000"), // Carmesim profundo
		Warning:   lipgloss.Color("#FFB300"), // Amarelo alerta
		Muted:     lipgloss.Color("#6E5059"), // Cinza arroxeado escuro
		Text:      lipgloss.Color("#FBECEF"), // Branco rosado suave
		Border:    lipgloss.Color("#382026"),
		KeyBg:     lipgloss.Color("#281B20"),
		KeyFg:     lipgloss.Color("#FF1744"),
	},
	{
		ID:        "matrix",
		Name:      "📟 Matrix",
		Bg:        lipgloss.Color("#0D1410"), // Carbono com leve tom verde-digital
		Primary:   lipgloss.Color("#00FF66"), // Fósforo verde clássico
		Secondary: lipgloss.Color("#26A69A"), // Verde-azulado terminal
		Success:   lipgloss.Color("#00FF66"), // Verde código
		Danger:    lipgloss.Color("#FF5252"), // Vermelho agente
		Warning:   lipgloss.Color("#FFD54F"), // Amarelo sentinela
		Muted:     lipgloss.Color("#3E6B48"), // Verde musgo atenuado
		Text:      lipgloss.Color("#C8E6C9"), // Fósforo claro suave
		Border:    lipgloss.Color("#1E3626"),
		KeyBg:     lipgloss.Color("#16241C"),
		KeyFg:     lipgloss.Color("#00FF66"),
	},
}

var currentTheme = themes[0] // Dracula como padrão

// ============================================================================
// FORMATADORES DE LAYOUT (100% ANSI-SAFE)
// ============================================================================

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