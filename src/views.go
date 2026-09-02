package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	t := currentTheme

	var renderedContent string

	switch m.state {
	case viewService:
		renderedContent = m.renderServiceView()
		return lipgloss.Place(
			m.winWidth,
			m.winHeight,
			lipgloss.Left,
			lipgloss.Top,
			renderedContent,
			lipgloss.WithWhitespaceBackground(t.Bg),
			lipgloss.WithWhitespaceChars(" "),
		)

	case viewTheme:
		modal := m.renderThemeView()
		return lipgloss.Place(
			m.winWidth,
			m.winHeight,
			lipgloss.Center,
			lipgloss.Center,
			modal,
			lipgloss.WithWhitespaceBackground(t.Bg),
			lipgloss.WithWhitespaceChars(" "),
		)

	case viewKillConfirm:
		modal := m.renderKillConfirmView()
		return lipgloss.Place(
			m.winWidth,
			m.winHeight,
			lipgloss.Center,
			lipgloss.Center,
			modal,
			lipgloss.WithWhitespaceBackground(t.Bg),
			lipgloss.WithWhitespaceChars(" "),
		)

	default:
		renderedMenu := m.renderMenuView()
		return lipgloss.Place(
			m.winWidth,
			m.winHeight,
			lipgloss.Center,
			lipgloss.Center,
			renderedMenu,
			lipgloss.WithWhitespaceBackground(t.Bg),
			lipgloss.WithWhitespaceChars(" "),
		)
	}
}

func (m model) renderMenuView() string {
	t := currentTheme

	menuWidth := m.winWidth - 6
	if menuWidth > maxMenuWidth {
		menuWidth = maxMenuWidth
	}
	if menuWidth < 40 {
		menuWidth = 40
	}

	badgeOn := lipgloss.NewStyle().Bold(true).Foreground(t.Bg).Background(t.Success).Padding(0, 1).SetString("● ON").Render()
	badgeOff := lipgloss.NewStyle().Bold(true).Foreground(t.Muted).Background(t.KeyBg).Padding(0, 1).SetString("○ OFF").Render()

	mySQLStatus := badgeOff
	if isMySQLRunning() {
		mySQLStatus = badgeOn
	}

	loginStatus := badgeOff
	if running, _ := isProcessRunning("login_server"); running {
		loginStatus = badgeOn
	}

	canaryStatus := badgeOff
	if running, _ := isProcessRunning("canary"); running {
		canaryStatus = badgeOn
	}

	lblSt := lipgloss.NewStyle().Foreground(t.Text).Background(t.Bg)
	spSt := lipgloss.NewStyle().Background(t.Bg)

	p1 := lblSt.Render("MySQL: ") + mySQLStatus
	p2 := spSt.Render("   ") + lblSt.Render("Login: ") + loginStatus
	p3 := spSt.Render("   ") + lblSt.Render("Canary: ") + canaryStatus
	quickBarText := p1 + p2 + p3

	innerBannerW := menuWidth - 4
	quickBarPadded := padLine(quickBarText, innerBannerW, t.Bg)

	topBanner := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		BorderBackground(t.Bg).
		Background(t.Bg).
		Padding(0, 1).
		Render(quickBarPadded)

	var menuRows []string

	for i, item := range m.menuItems {
		isSelected := (i == m.menuIndex)
		rowBg := t.Bg
		if isSelected {
			rowBg = t.KeyBg
		}
		rowSt := lipgloss.NewStyle().Background(rowBg)

		badgeMissing := lipgloss.NewStyle().Bold(true).Foreground(t.Bg).Background(t.Warning).Padding(0, 1).SetString("⚠ NÃO COMPILADO").Render()
		badgeAction := lipgloss.NewStyle().Bold(true).Foreground(t.Bg).Background(t.Secondary).Padding(0, 1).SetString("⚡ EXECUTAR").Render()
		badgeTheme := lipgloss.NewStyle().Bold(true).Foreground(t.Bg).Background(t.Primary).Padding(0, 1).SetString(t.ID).Render()

		var badge string
		if item.checkBin != "" && !fileExists(item.checkBin) {
			badge = badgeMissing
		} else if item.isMySQL {
			if isMySQLRunning() {
				badge = badgeOn
			} else {
				badge = badgeOff
			}
		} else if item.procName != "" {
			if running, _ := isProcessRunning(item.procName); running {
				badge = badgeOn
			} else {
				badge = badgeOff
			}
		} else if item.isTheme {
			badge = badgeTheme
		} else if item.isExit {
			badge = ""
		} else {
			badge = badgeAction
		}

		cursor := "  "
		titleSt := lipgloss.NewStyle().Bold(true).Foreground(t.Text).Background(rowBg)
		descSt := lipgloss.NewStyle().Foreground(t.Muted).Background(rowBg)

		if isSelected {
			cursor = lipgloss.NewStyle().Foreground(t.Primary).Background(rowBg).Bold(true).Render("▶ ")
			titleSt = titleSt.Foreground(t.Primary)
			descSt = descSt.Foreground(t.Text)
		} else {
			cursor = rowSt.Render("  ")
		}

		leftTitle := cursor + titleSt.Render(item.title)
		badgeWidth := lipgloss.Width(badge)
		leftWidth := lipgloss.Width(leftTitle)

		spacesCount := menuWidth - leftWidth - badgeWidth
		if spacesCount < 1 {
			spacesCount = 1
		}
		spaces := rowSt.Render(strings.Repeat(" ", spacesCount))

		line1Content := leftTitle + spaces + badge
		line2Content := rowSt.Render("   ") + descSt.Render(item.desc)

		line1 := padLine(line1Content, menuWidth, rowBg)
		line2 := padLine(line2Content, menuWidth, rowBg)

		menuRows = append(menuRows, line1, line2)

		if i < len(m.menuItems)-1 {
			menuRows = append(menuRows, padLine("", menuWidth, t.Bg))
		}
	}

	footerText := formatFooterBar("⚡ ATALHOS:", []string{
		formatAction("↑/↓", "Navegar", t),
		formatAction("ENTER", "Selecionar", t),
		formatAction("T", "Temas", t),
		formatAction("Q", "Sair", t),
	}, t)

	footerPadded := padLine(footerText, menuWidth, t.Bg)
	spacer := padLine("", menuWidth, t.Bg)

	allElements := []string{
		topBanner,
		spacer,
	}
	allElements = append(allElements, menuRows...)
	allElements = append(allElements, spacer, footerPadded)

	return lipgloss.JoinVertical(lipgloss.Left, allElements...)
}

func (m model) renderServiceView() string {
	t := currentTheme
	fullWidth := m.winWidth - 4
	if fullWidth < 40 {
		fullWidth = 40
	}

	var name string
	var statusStr string
	var pidStr string
	var isRunning bool

	badgeOn := lipgloss.NewStyle().Bold(true).Foreground(t.Bg).Background(t.Success).Padding(0, 1).SetString("● ON").Render()
	badgeOff := lipgloss.NewStyle().Bold(true).Foreground(t.Muted).Background(t.KeyBg).Padding(0, 1).SetString("○ OFF").Render()
	badgeAction := lipgloss.NewStyle().Bold(true).Foreground(t.Bg).Background(t.Secondary).Padding(0, 1).SetString("⚡ EXECUTAR").Render()

	switch m.activeTarget {
	case targetCanary:
		name = "CANARY OT SERVER"
		running, pid := isProcessRunning("canary")
		isRunning = running
		if running {
			statusStr = badgeOn
			pidStr = lipgloss.NewStyle().Foreground(t.Success).Background(t.Bg).Bold(true).Render(fmt.Sprintf("PID: %d", pid))
		} else {
			statusStr = badgeOff
			pidStr = lipgloss.NewStyle().Foreground(t.Muted).Background(t.Bg).Render("Inativo")
		}

	case targetLogin:
		name = "LOGIN SERVER"
		running, pid := isProcessRunning("login_server")
		isRunning = running
		if running {
			statusStr = badgeOn
			pidStr = lipgloss.NewStyle().Foreground(t.Success).Background(t.Bg).Bold(true).Render(fmt.Sprintf("PID: %d", pid))
		} else {
			statusStr = badgeOff
			pidStr = lipgloss.NewStyle().Foreground(t.Muted).Background(t.Bg).Render("Inativo")
		}

	case targetMySQL:
		name = "MYSQL DATABASE"
		isRunning = isMySQLRunning()
		if isRunning {
			statusStr = badgeOn
			pidStr = lipgloss.NewStyle().Foreground(t.Success).Background(t.Bg).Bold(true).Render("Ativo")
		} else {
			statusStr = badgeOff
			pidStr = lipgloss.NewStyle().Foreground(t.Muted).Background(t.Bg).Render("Inativo")
		}

	case targetClient:
		name = "CLIENT PATCHER 15.25"
		statusStr = badgeAction
		pidStr = lipgloss.NewStyle().Foreground(t.Secondary).Background(t.Bg).Render("Target: client.exe")
		isRunning = false
	}

	lblSt := lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Background(t.Bg)
	txtSt := lipgloss.NewStyle().Foreground(t.Text).Background(t.Bg)
	spSt := lipgloss.NewStyle().Background(t.Bg)

	h1 := lblSt.Render("🖥️  PAINEL DO SERVIÇO: ") + txtSt.Render(name)
	h2 := spSt.Render("   ") + lblSt.Render("Status: ") + statusStr
	h3 := spSt.Render("   (") + pidStr + txtSt.Render(")")
	headerText := h1 + h2 + h3
	headerRendered := padLine(headerText, fullWidth, t.Bg)

	logsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		BorderBackground(t.Bg).
		Background(lipgloss.Color("#000000")).
		Padding(0, 1).
		Render(m.logsViewport.View())

	var notice string
	if m.actionNotice != "" {
		notice = padLine(lipgloss.NewStyle().Foreground(t.Warning).Background(t.Bg).Bold(true).Render(m.actionNotice), fullWidth, t.Bg)
	}

	var actionsBar string
	if m.activeTarget == targetClient {
		actionsBar = formatFooterBar("⚡ AÇÕES:", []string{
			formatAction("ENTER", "Aplicar Novamente", t),
			formatAction("Esc", "❮ Voltar ao Menu", t),
		}, t)
	} else if isRunning {
		actionsBar = formatFooterBar("⚡ AÇÕES:", []string{
			formatAction("S", "🛑 Safe Stop (Salvar)", t),
			formatAction("K", "💀 Force Kill (-9)", t),
			formatAction("R", "🔄 Atualizar Logs", t),
			formatAction("Esc", "❮ Voltar ao Menu", t),
		}, t)
	} else {
		actionsBar = formatFooterBar("⚡ AÇÕES:", []string{
			formatAction("ENTER / S", "🚀 Iniciar Novamente", t),
			formatAction("R", "🔄 Atualizar Logs", t),
			formatAction("Esc", "❮ Voltar ao Menu", t),
		}, t)
	}

	actionsRendered := padLine(actionsBar, fullWidth, t.Bg)

	elements := []string{
		headerRendered,
		logsBox,
	}

	if notice != "" {
		elements = append(elements, notice)
	}
	elements = append(elements, actionsRendered)

	return lipgloss.JoinVertical(lipgloss.Left, elements...)
}

func (m model) renderThemeView() string {
	t := currentTheme
	const modalWidth = 66
	const innerWidth = modalWidth - 4
	const targetNameWidth = 32

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary).
		Background(t.Bg).
		Render("🎨 SELECIONAR TEMA VISUAL (DARK)")

	titlePadded := padLine(title, innerWidth, t.Bg)

	var items []string

	for i, theme := range themes {
		cursor := "  "
		itemStyle := lipgloss.NewStyle().Foreground(t.Text).Background(t.Bg)

		if i == m.themeIndex {
			cursor = lipgloss.NewStyle().Foreground(theme.Primary).Background(t.Bg).Bold(true).Render("▶ ")
			itemStyle = itemStyle.Bold(true).Foreground(theme.Primary)
		} else {
			cursor = lipgloss.NewStyle().Background(t.Bg).Render("  ")
		}

		visualLen := lipgloss.Width(theme.Name)
		pad := targetNameWidth - visualLen
		if pad < 1 {
			pad = 1
		}
		spaces := lipgloss.NewStyle().Background(t.Bg).Render(strings.Repeat(" ", pad))

		chip1 := lipgloss.NewStyle().Foreground(theme.Primary).Background(t.Bg).Render("██")
		chip2 := lipgloss.NewStyle().Foreground(theme.Secondary).Background(t.Bg).Render("██")
		chip3 := lipgloss.NewStyle().Foreground(theme.Success).Background(t.Bg).Render("██")
		chip4 := lipgloss.NewStyle().Foreground(theme.Warning).Background(t.Bg).Render("██")
		chip5 := lipgloss.NewStyle().Foreground(theme.Danger).Background(t.Bg).Render("██")
		sp := lipgloss.NewStyle().Background(t.Bg).Render(" ")

		chips := chip1 + sp + chip2 + sp + chip3 + sp + chip4 + sp + chip5
		nameRendered := itemStyle.Render(theme.Name)

		lineContent := cursor + nameRendered + spaces + chips
		items = append(items, padLine(lineContent, innerWidth, t.Bg))
	}

	footer := lipgloss.NewStyle().
		Foreground(t.Muted).
		Background(t.Bg).
		Render("Navegue com [↑/↓]  •  [ENTER] Confirmar  •  [Esc] Cancelar")

	footerPadded := padLine(footer, innerWidth, t.Bg)
	spacer := padLine("", innerWidth, t.Bg)

	var allLines []string
	allLines = append(allLines, titlePadded, spacer)
	allLines = append(allLines, items...)
	allLines = append(allLines, spacer, footerPadded)

	boxContent := lipgloss.JoinVertical(lipgloss.Left, allLines...)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		BorderBackground(t.Bg).
		Background(t.Bg).
		Padding(1, 2).
		Width(modalWidth).
		Render(boxContent)
}

func (m model) renderKillConfirmView() string {
	t := currentTheme
	const modalWidth = 76
	const innerWidth = modalWidth - 6

	procName := "CANARY SERVER"
	targetProc := "canary"

	if m.activeTarget == targetLogin {
		procName = "LOGIN SERVER"
		targetProc = "login_server"
	} else if m.activeTarget == targetMySQL {
		procName = "MYSQL DATABASE"
		targetProc = "mysqld"
	}

	_, pid := isProcessRunning(targetProc)

	lines := []string{
		padLine("☠️  ⚠️  ALERTA CRÍTICO: FORCE KILL (SIGKILL -9)  ⚠️  ☠️", innerWidth, t.Bg),
		padLine("", innerWidth, t.Bg),
		padLine(fmt.Sprintf("Você está prestes a EXTERMINAR o processo %s (PID: %d)!", procName, pid), innerWidth, t.Bg),
		padLine("", innerWidth, t.Bg),
		padLine("❌ NÃO HAVERÁ SALVAMENTO DE MAPA OU CASAS!", innerWidth, t.Bg),
		padLine("❌ JOGADORES ONLINE SOFRERÃO ROLLBACK IMEDIATO!", innerWidth, t.Bg),
		padLine("❌ RISCO DE TRANSAÇÕES PENDENTES NO BANCO DE DADOS!", innerWidth, t.Bg),
		padLine("", innerWidth, t.Bg),
		padLine("Para confirmar a aniquilação forçada, digite KILL abaixo e tecle ENTER:", innerWidth, t.Bg),
		padLine("", innerWidth, t.Bg),
		padLine(m.killInput.View(), innerWidth, t.Bg),
		padLine(lipgloss.NewStyle().Bold(true).Foreground(t.Danger).Background(t.Bg).Render(m.killFeedback), innerWidth, t.Bg),
		padLine("", innerWidth, t.Bg),
		padLine("[ Pressione ESC para cancelar e voltar à segurança ]", innerWidth, t.Bg),
	}

	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(t.Danger).
		BorderBackground(t.Bg).
		Background(t.Bg).
		Padding(1, 2).
		Width(modalWidth).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}