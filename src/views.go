package main

import (
	"fmt"
	"path/filepath"
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

	case viewSetupMenu:
		renderedContent = m.renderSetupMenuView()
		return lipgloss.Place(
			m.winWidth,
			m.winHeight,
			lipgloss.Center,
			lipgloss.Center,
			renderedContent,
			lipgloss.WithWhitespaceBackground(t.Bg),
			lipgloss.WithWhitespaceChars(" "),
		)

	case viewPlayersMenu:
		renderedContent = m.renderPlayersMenuView()
		return lipgloss.Place(
			m.winWidth,
			m.winHeight,
			lipgloss.Center,
			lipgloss.Center,
			renderedContent,
			lipgloss.WithWhitespaceBackground(t.Bg),
			lipgloss.WithWhitespaceChars(" "),
		)

	case viewToolsMenu:
		renderedContent = m.renderToolsMenuView()
		return lipgloss.Place(
			m.winWidth,
			m.winHeight,
			lipgloss.Center,
			lipgloss.Center,
			renderedContent,
			lipgloss.WithWhitespaceBackground(t.Bg),
			lipgloss.WithWhitespaceChars(" "),
		)

	case viewModeSelector:
		modal := m.renderModeSelectorView()
		return lipgloss.Place(
			m.winWidth,
			m.winHeight,
			lipgloss.Center,
			lipgloss.Center,
			modal,
			lipgloss.WithWhitespaceBackground(t.Bg),
			lipgloss.WithWhitespaceChars(" "),
		)

	case viewEditPaths:
		renderedContent = m.renderEditPathsView()
		return lipgloss.Place(
			m.winWidth,
			m.winHeight,
			lipgloss.Center,
			lipgloss.Center,
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

		badgeMissing := lipgloss.NewStyle().Bold(true).Foreground(t.Bg).Background(t.Warning).Padding(0, 1).SetString("⚠ NÃO INSTALADO").Render()
		badgeSubmenu := lipgloss.NewStyle().Bold(true).Foreground(t.Bg).Background(t.Secondary).Padding(0, 1).SetString("⚡ SUBMENU").Render()
		badgeAction := lipgloss.NewStyle().Bold(true).Foreground(t.Bg).Background(t.Secondary).Padding(0, 1).SetString("⚡ EXECUTAR").Render()
		badgeTheme := lipgloss.NewStyle().Bold(true).Foreground(t.Bg).Background(t.Primary).Padding(0, 1).SetString(t.ID).Render()

		var badge string
		if item.id == "canary" && !fileExists(filepath.Join(expandHome(appConfig.CanaryDir), "canary")) {
			badge = badgeMissing
		} else if item.id == "login" && !fileExists(filepath.Join(expandHome(appConfig.LoginDir), "login_server")) {
			badge = badgeMissing
		} else if item.isMySQL {
			if isMySQLRunning() {
				badge = badgeOn
			} else {
				badge = badgeOff
			}
		} else if item.isPlayers || item.isTools || item.isSetup {
			badge = badgeSubmenu
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

func (m model) renderSetupMenuView() string {
	t := currentTheme
	menuWidth := m.winWidth - 6
	if menuWidth > maxMenuWidth {
		menuWidth = maxMenuWidth
	}

	modeStr := "🖥️ Jogando pelo Windows (WSL)"
	if appConfig.Mode == ModeVPS {
		modeStr = "🌐 Servidor na Nuvem (VPS)"
	} else if appConfig.Mode == ModeLinuxLocal {
		modeStr = "💻 Linux Local"
	}

	titleText := lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Background(t.Bg).Render(fmt.Sprintf("⚙️  SETUP & CONEXÕES   [ Modo: %s ]", modeStr))
	innerBannerW := menuWidth - 4
	titlePadded := padLine(titleText, innerBannerW, t.Bg)

	topBanner := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		BorderBackground(t.Bg).
		Background(t.Bg).
		Padding(0, 1).
		Render(titlePadded)

	var menuRows []string

	for i, item := range m.setupItems {
		isSelected := (i == m.setupIndex)
		rowBg := t.Bg
		if isSelected {
			rowBg = t.KeyBg
		}
		rowSt := lipgloss.NewStyle().Background(rowBg)

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
		line1 := padLine(leftTitle, menuWidth, rowBg)
		line2 := padLine(rowSt.Render("   ")+descSt.Render(item.desc), menuWidth, rowBg)

		menuRows = append(menuRows, line1, line2)

		if i < len(m.setupItems)-1 {
			menuRows = append(menuRows, padLine("", menuWidth, t.Bg))
		}
	}

	footerText := formatFooterBar("⚡ ATALHOS:", []string{
		formatAction("↑/↓", "Navegar", t),
		formatAction("ENTER", "Executar", t),
		formatAction("Esc", "❮ Voltar", t),
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

func (m model) renderPlayersMenuView() string {
	t := currentTheme
	menuWidth := m.winWidth - 6
	if menuWidth > maxMenuWidth {
		menuWidth = maxMenuWidth
	}

	titleText := lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Background(t.Bg).Render("👥  GERENCIAR JOGADORES & CONTAS")
	innerBannerW := menuWidth - 4
	titlePadded := padLine(titleText, innerBannerW, t.Bg)

	topBanner := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		BorderBackground(t.Bg).
		Background(t.Bg).
		Padding(0, 1).
		Render(titlePadded)

	var menuRows []string

	for i, item := range m.playersItems {
		isSelected := (i == m.playersIndex)
		rowBg := t.Bg
		if isSelected {
			rowBg = t.KeyBg
		}
		rowSt := lipgloss.NewStyle().Background(rowBg)

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
		line1 := padLine(leftTitle, menuWidth, rowBg)
		line2 := padLine(rowSt.Render("   ")+descSt.Render(item.desc), menuWidth, rowBg)

		menuRows = append(menuRows, line1, line2)

		if i < len(m.playersItems)-1 {
			menuRows = append(menuRows, padLine("", menuWidth, t.Bg))
		}
	}

	footerText := formatFooterBar("⚡ ATALHOS:", []string{
		formatAction("↑/↓", "Navegar", t),
		formatAction("ENTER", "Selecionar", t),
		formatAction("Esc", "❮ Voltar", t),
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

func (m model) renderToolsMenuView() string {
	t := currentTheme
	menuWidth := m.winWidth - 6
	if menuWidth > maxMenuWidth {
		menuWidth = maxMenuWidth
	}

	titleText := lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Background(t.Bg).Render("🛠️  FERRAMENTAS & CALIBRAÇÃO GLOBAL")
	innerBannerW := menuWidth - 4
	titlePadded := padLine(titleText, innerBannerW, t.Bg)

	topBanner := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		BorderBackground(t.Bg).
		Background(t.Bg).
		Padding(0, 1).
		Render(titlePadded)

	var menuRows []string

	for i, item := range m.toolsItems {
		isSelected := (i == m.toolsIndex)
		rowBg := t.Bg
		if isSelected {
			rowBg = t.KeyBg
		}
		rowSt := lipgloss.NewStyle().Background(rowBg)

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
		line1 := padLine(leftTitle, menuWidth, rowBg)
		line2 := padLine(rowSt.Render("   ")+descSt.Render(item.desc), menuWidth, rowBg)

		menuRows = append(menuRows, line1, line2)

		if i < len(m.toolsItems)-1 {
			menuRows = append(menuRows, padLine("", menuWidth, t.Bg))
		}
	}

	footerText := formatFooterBar("⚡ ATALHOS:", []string{
		formatAction("↑/↓", "Navegar", t),
		formatAction("ENTER", "Selecionar", t),
		formatAction("Esc", "❮ Voltar", t),
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

func (m model) renderModeSelectorView() string {
	t := currentTheme
	const modalWidth = 74
	const innerWidth = modalWidth - 6

	title := padLine(lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Background(t.Bg).Render("🌐 ESCOLHA COMO VOCÊ ESTÁ USANDO ESTE SERVIDOR"), innerWidth, t.Bg)
	spacer := padLine("", innerWidth, t.Bg)

	options := []struct {
		title string
		desc  string
	}{
		{
			title: "🖥️  Jogando pelo Windows (WSL)",
			desc:  "Você joga pelo Windows e o servidor roda no WSL. O endereço é atualizado sozinho.",
		},
		{
			title: "🌐  Servidor na Nuvem / Para Amigos (VPS)",
			desc:  "O servidor fica em uma máquina na internet para outras pessoas jogarem.",
		},
		{
			title: "💻  Linux Completo (Tudo no mesmo PC)",
			desc:  "Você usa apenas Linux para tudo (o servidor e o jogo rodam no mesmo sistema).",
		},
	}

	var lines []string
	lines = append(lines, title, spacer)

	for i, opt := range options {
		isSelected := (i == m.modeIndex)
		rowBg := t.Bg
		if isSelected {
			rowBg = t.KeyBg
		}
		rowSt := lipgloss.NewStyle().Background(rowBg)

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

		line1 := padLine(cursor+titleSt.Render(opt.title), innerWidth, rowBg)
		line2 := padLine(rowSt.Render("   ")+descSt.Render(opt.desc), innerWidth, rowBg)

		lines = append(lines, line1, line2)
		if i < len(options)-1 {
			lines = append(lines, padLine("", innerWidth, t.Bg))
		}
	}

	footer := formatFooterBar("⚡ ATALHOS:", []string{
		formatAction("↑/↓", "Navegar", t),
		formatAction("ENTER", "Confirmar Modo", t),
		formatAction("Esc", "Cancelar", t),
	}, t)

	lines = append(lines, spacer, padLine(footer, innerWidth, t.Bg))

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		BorderBackground(t.Bg).
		Background(t.Bg).
		Padding(1, 2).
		Width(modalWidth).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func (m model) renderEditPathsView() string {
	t := currentTheme
	const modalWidth = 76
	const innerWidth = modalWidth - 6

	labels := []string{
		"1. Pasta do Canary Server:",
		"2. Pasta do Login Server:",
		"3. Pasta do Client Editor:",
		"4. Pasta de Download do Client:",
		"5. Caminho do Executável (client.exe):",
	}

	lines := []string{
		padLine(lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Background(t.Bg).Render("📁 CONFIGURAR CAMINHOS DAS PASTAS"), innerWidth, t.Bg),
		padLine(lipgloss.NewStyle().Foreground(t.Muted).Background(t.Bg).Render("Altere apenas se você já baixou os projetos em pastas manuais:"), innerWidth, t.Bg),
		padLine("", innerWidth, t.Bg),
	}

	for i := range m.pathInputs {
		lblSt := lipgloss.NewStyle().Foreground(t.Text).Background(t.Bg)
		if i == m.pathInputIndex {
			lblSt = lblSt.Bold(true).Foreground(t.Primary)
		}
		lines = append(lines, padLine(lblSt.Render(labels[i]), innerWidth, t.Bg))
		lines = append(lines, padLine(m.pathInputs[i].View(), innerWidth, t.Bg))
		lines = append(lines, padLine("", innerWidth, t.Bg))
	}

	footer := formatFooterBar("⚡ ATALHOS:", []string{
		formatAction("TAB / ↑↓", "Alternar Campo", t),
		formatAction("ENTER", "Salvar Configurações", t),
		formatAction("Esc", "Cancelar", t),
	}, t)

	lines = append(lines, padLine(footer, innerWidth, t.Bg))

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		BorderBackground(t.Bg).
		Background(t.Bg).
		Padding(1, 2).
		Width(modalWidth).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
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

	case targetClientPatch:
		name = "CONFIGURAR TIBIA CLIENT 15.25"
		statusStr = badgeAction
		pidStr = lipgloss.NewStyle().Foreground(t.Secondary).Background(t.Bg).Render("Target: client.exe")
		isRunning = false

	case targetFullSetup:
		name = "INSTALAÇÃO AUTOMÁTICA COMPLETA"
		statusStr = badgeAction
		pidStr = lipgloss.NewStyle().Foreground(t.Secondary).Background(t.Bg).Render("Status: Setup")
		isRunning = false

	case targetDownloadClient:
		name = "DOWNLOAD TIBIA CLIENT 15.25"
		statusStr = badgeAction
		pidStr = lipgloss.NewStyle().Foreground(t.Secondary).Background(t.Bg).Render("Status: Download")
		isRunning = false

	case targetSyncConnections:
		name = "SINCRONIZAR CONEXÕES DE REDE"
		statusStr = badgeAction
		pidStr = lipgloss.NewStyle().Foreground(t.Secondary).Background(t.Bg).Render("Status: Sync")
		isRunning = false

	default:
		name = "PAINEL DE FERRAMENTA"
		statusStr = badgeAction
		pidStr = lipgloss.NewStyle().Foreground(t.Secondary).Background(t.Bg).Render("Status: Preview")
		isRunning = false
	}

	lblSt := lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Background(t.Bg)
	txtSt := lipgloss.NewStyle().Foreground(t.Text).Background(t.Bg)
	spSt := lipgloss.NewStyle().Background(t.Bg)

	h1 := lblSt.Render("🖥️  PAINEL: ") + txtSt.Render(name)
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
	if m.activeTarget == targetClientPatch || m.activeTarget == targetFullSetup || m.activeTarget == targetDownloadClient || m.activeTarget == targetSyncConnections || m.activeTarget == targetPlaceholder {
		actionsBar = formatFooterBar("⚡ AÇÕES:", []string{
			formatAction("ENTER", "Executar Novamente", t),
			formatAction("Esc", "❮ Voltar", t),
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