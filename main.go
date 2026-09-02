package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ============================================================================
// CONFIGURAÇÕES E CAMINHOS
// ============================================================================

const (
	maxMenuWidth = 84

	PathCanaryDir = "~/Code/canary"
	PathCanaryBin = "~/Code/canary/canary"
	PathCanaryLog = "~/Code/canary/canary.log"

	PathLoginDir = "~/Code/login-server"
	PathLoginBin = "~/Code/login-server/login_server"
	PathLoginLog = "~/Code/login-server/login_server.log"

	PathClientDir    = "~/Code/client-editor"
	PathClientBin    = "~/Code/client-editor/client-editor-linux-x64"
	PathTibiaExe     = "/mnt/c/Users/T-GAMER/Code/tibia-client/bin/client.exe"
	ClientEditorArgs = "edit -t /mnt/c/Users/T-GAMER/Code/tibia-client/bin/client.exe -c local.toml"
)

// ============================================================================
// SISTEMA DE TEMAS DARK (PALETAS OFICIAIS)
// ============================================================================

type Theme struct {
	ID        string
	Name      string
	Bg        lipgloss.Color // Fundo oficial da aplicação
	Primary   lipgloss.Color // Cor primária / Destaques
	Secondary lipgloss.Color // Cor secundária
	Success   lipgloss.Color // Verde de status online
	Danger    lipgloss.Color // Vermelho de kill / Erro
	Warning   lipgloss.Color // Amarelo de aviso / Alerta
	Muted     lipgloss.Color // Texto atenuado / Comentários
	Text      lipgloss.Color // Texto principal
	Border    lipgloss.Color // Bordas de caixas
	KeyBg     lipgloss.Color // Fundo de keycaps e seleção
	KeyFg     lipgloss.Color // Texto de keycaps
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
		Border:    lipgloss.Color("#44475A"),
		KeyBg:     lipgloss.Color("#44475A"),
		KeyFg:     lipgloss.Color("#F8F8F2"),
	},
	{
		ID:        "tokyo-night",
		Name:      "🌌 Tokyo Night",
		Bg:        lipgloss.Color("#1A1B26"),
		Primary:   lipgloss.Color("#7AA2F7"),
		Secondary: lipgloss.Color("#BB9AF7"),
		Success:   lipgloss.Color("#9ECE6A"),
		Danger:    lipgloss.Color("#F7768E"),
		Warning:   lipgloss.Color("#E0AF68"),
		Muted:     lipgloss.Color("#565F89"),
		Text:      lipgloss.Color("#C0CAF5"),
		Border:    lipgloss.Color("#292E42"),
		KeyBg:     lipgloss.Color("#24283B"),
		KeyFg:     lipgloss.Color("#7AA2F7"),
	},
	{
		ID:        "catppuccin",
		Name:      "🌸 Catppuccin (Mocha)",
		Bg:        lipgloss.Color("#1E1E2E"),
		Primary:   lipgloss.Color("#CBA6F7"),
		Secondary: lipgloss.Color("#89B4FA"),
		Success:   lipgloss.Color("#A6E3A1"),
		Danger:    lipgloss.Color("#F38BA8"),
		Warning:   lipgloss.Color("#FAB387"),
		Muted:     lipgloss.Color("#6C7086"),
		Text:      lipgloss.Color("#CDD6F4"),
		Border:    lipgloss.Color("#313244"),
		KeyBg:     lipgloss.Color("#313244"),
		KeyFg:     lipgloss.Color("#CBA6F7"),
	},
	{
		ID:        "gruvbox",
		Name:      "🍂 Gruvbox (Dark)",
		Bg:        lipgloss.Color("#282828"),
		Primary:   lipgloss.Color("#FE8019"),
		Secondary: lipgloss.Color("#FABD2F"),
		Success:   lipgloss.Color("#B8BB26"),
		Danger:    lipgloss.Color("#FB4934"),
		Warning:   lipgloss.Color("#FABD2F"),
		Muted:     lipgloss.Color("#928374"),
		Text:      lipgloss.Color("#EBDBB2"),
		Border:    lipgloss.Color("#504945"),
		KeyBg:     lipgloss.Color("#3C3836"),
		KeyFg:     lipgloss.Color("#FE8019"),
	},
	{
		ID:        "orng",
		Name:      "🍊 Orng (Cyber Amber)",
		Bg:        lipgloss.Color("#1C120C"),
		Primary:   lipgloss.Color("#FF7700"),
		Secondary: lipgloss.Color("#FFAA00"),
		Success:   lipgloss.Color("#00E676"),
		Danger:    lipgloss.Color("#FF3D00"),
		Warning:   lipgloss.Color("#FFD600"),
		Muted:     lipgloss.Color("#8C6B53"),
		Text:      lipgloss.Color("#FFECD6"),
		Border:    lipgloss.Color("#4D2C19"),
		KeyBg:     lipgloss.Color("#2B1B12"),
		KeyFg:     lipgloss.Color("#FF7700"),
	},
	{
		ID:        "default",
		Name:      "⚡ Default (Cyberpunk Neon)",
		Bg:        lipgloss.Color("#0D1117"),
		Primary:   lipgloss.Color("#00E5FF"),
		Secondary: lipgloss.Color("#BD93F9"),
		Success:   lipgloss.Color("#00FF88"),
		Danger:    lipgloss.Color("#FF2A6D"),
		Warning:   lipgloss.Color("#FFEA00"),
		Muted:     lipgloss.Color("#6E7681"),
		Text:      lipgloss.Color("#F0F6FC"),
		Border:    lipgloss.Color("#30363D"),
		KeyBg:     lipgloss.Color("#161B22"),
		KeyFg:     lipgloss.Color("#00E5FF"),
	},
}

var currentTheme = themes[0] // Começa em Dracula

// ============================================================================
// FORMATADORES E PREENCHEDORES DE FUNDO (100% ANSI-SAFE)
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

// ============================================================================
// UTILITÁRIOS DE SISTEMA & STATUS
// ============================================================================

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func fileExists(path string) bool {
	_, err := os.Stat(expandHome(path))
	return err == nil
}

func isProcessRunning(processName string) (bool, int) {
	cmd := exec.Command("pgrep", "-x", processName)
	output, err := cmd.Output()
	if err != nil {
		return false, 0
	}
	pids := strings.Fields(strings.TrimSpace(string(output)))
	myPid := os.Getpid()
	for _, pStr := range pids {
		var pid int
		fmt.Sscanf(pStr, "%d", &pid)
		if pid > 0 && pid != myPid {
			return true, pid
		}
	}
	return false, 0
}

func isMySQLRunning() bool {
	if ok, _ := isProcessRunning("mysqld"); ok {
		return true
	}
	cmd := exec.Command("service", "mysql", "status")
	out, err := cmd.CombinedOutput()
	return err == nil && strings.Contains(strings.ToLower(string(out)), "is running")
}

func readLastLogs(logFilePath string, linesCount int) string {
	expanded := expandHome(logFilePath)
	if !fileExists(logFilePath) {
		return "Aguardando inicialização dos logs em: " + expanded
	}
	cmd := exec.Command("tail", fmt.Sprintf("-n%d", linesCount), expanded)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "Lendo logs... (" + err.Error() + ")"
	}
	raw := string(out)
	if len(strings.TrimSpace(raw)) == 0 {
		return "O arquivo de log está aguardando saída do processo..."
	}

	cleaned := strings.ReplaceAll(raw, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "")

	// Remove cabeçalhos gerados pelo terminal virtual
	var filtered []string
	for _, line := range strings.Split(cleaned, "\n") {
		if strings.HasPrefix(line, "Script started on") || strings.HasPrefix(line, "Script done on") {
			continue
		}
		filtered = append(filtered, line)
	}

	return strings.Join(filtered, "\n")
}

// ============================================================================
// COMANDOS TEA ASSÍNCRONOS
// ============================================================================

type tickMsg time.Time
type actionDoneMsg struct {
	target serviceTarget
	output string
	notice string
}

func startBackgroundProcess(target serviceTarget, workDir, binName, logPath string) tea.Cmd {
	return func() tea.Msg {
		expandedWorkDir := expandHome(workDir)
		expandedLogPath := expandHome(logPath)
		script := fmt.Sprintf(
			"export CLICOLOR_FORCE=1 FORCE_COLOR=1 COLORTERM=truecolor; cd '%s' && (script -qef -c './%s' '%s' >/dev/null 2>&1 &)",
			expandedWorkDir, binName, expandedLogPath,
		)
		_ = exec.Command("bash", "-c", script).Run()
		time.Sleep(300 * time.Millisecond)
		return actionDoneMsg{
			target: target,
			notice: "✔ Processo iniciado com cores ativas!",
		}
	}
}

func stopProcessGraceful(target serviceTarget, procName string) tea.Cmd {
	return func() tea.Msg {
		_ = exec.Command("pkill", "-INT", "-x", procName).Run()
		for i := 0; i < 20; i++ {
			time.Sleep(200 * time.Millisecond)
			if running, _ := isProcessRunning(procName); !running {
				break
			}
		}
		time.Sleep(300 * time.Millisecond)
		return actionDoneMsg{
			target: target,
			notice: "✔ Servidor finalizado e salvo com segurança.",
		}
	}
}

func killProcessHard(target serviceTarget, procName string) tea.Cmd {
	return func() tea.Msg {
		_ = exec.Command("pkill", "-9", "-x", procName).Run()
		time.Sleep(400 * time.Millisecond)
		return actionDoneMsg{
			target: target,
			notice: "💀 Processo aniquilado com SIGKILL (-9)!",
		}
	}
}

func manageMySQLNative(action string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("sudo", "service", "mysql", action)
		out, err := cmd.CombinedOutput()
		if err != nil {
			cmd = exec.Command("service", "mysql", action)
			out, _ = cmd.CombinedOutput()
		}
		time.Sleep(400 * time.Millisecond)

		statusOut, _ := exec.Command("service", "mysql", "status").CombinedOutput()
		statusText := string(statusOut)
		if len(strings.TrimSpace(statusText)) == 0 {
			statusText = string(out)
		}

		notice := "✔ MySQL iniciado com sucesso!"
		if action == "stop" {
			notice = "✔ MySQL desligado com sucesso."
		}

		return actionDoneMsg{
			target: targetMySQL,
			output: statusText,
			notice: notice,
		}
	}
}

func runClientPatchNative() tea.Cmd {
	return func() tea.Msg {
		expandedWorkDir := expandHome(PathClientDir)
		script := fmt.Sprintf("export CLICOLOR_FORCE=1 FORCE_COLOR=1; cd '%s' && ./client-editor-linux-x64 %s", expandedWorkDir, ClientEditorArgs)
		out, err := exec.Command("bash", "-c", script).CombinedOutput()
		outputStr := strings.ReplaceAll(string(out), "\r", "")
		if err != nil && len(strings.TrimSpace(outputStr)) == 0 {
			outputStr = "❌ Erro ao aplicar patch: " + err.Error()
		}
		time.Sleep(400 * time.Millisecond)
		return actionDoneMsg{
			target: targetClient,
			output: outputStr,
			notice: "✔ Patch do client aplicado com sucesso!",
		}
	}
}

// ============================================================================
// MODELO DA APLICAÇÃO
// ============================================================================

type viewState int

const (
	viewMenu viewState = iota
	viewService
	viewKillConfirm
	viewTheme
)

type serviceTarget int

const (
	targetCanary serviceTarget = iota
	targetLogin
	targetMySQL
	targetClient
)

type menuItem struct {
	id       string
	title    string
	desc     string
	target   serviceTarget
	checkBin string
	procName string
	isTheme  bool
	isExit   bool
	isMySQL  bool
}

type model struct {
	state        viewState
	activeTarget serviceTarget
	menuIndex    int
	menuItems    []menuItem
	logsViewport viewport.Model
	killInput    textinput.Model
	killFeedback string
	actionNotice string
	clientLogs   string
	mysqlLogs    string
	themeIndex   int
	winWidth     int
	winHeight    int
}

func getMenuItems() []menuItem {
	return []menuItem{
		{
			id:       "canary",
			title:    "🚀 Canary Server (OT Server)",
			desc:     "Iniciar, monitorar logs ao vivo ou encerrar com segurança",
			target:   targetCanary,
			checkBin: PathCanaryBin,
			procName: "canary",
		},
		{
			id:       "login",
			title:    "🔑 Login Server",
			desc:     "Iniciar, monitorar logs ao vivo ou encerrar",
			target:   targetLogin,
			checkBin: PathLoginBin,
			procName: "login_server",
		},
		{
			id:      "mysql",
			title:   "🗄️ MySQL Database",
			desc:    "Ligar, desligar ou verificar estado do banco",
			target:  targetMySQL,
			isMySQL: true,
		},
		{
			id:       "client",
			title:    "🎮 Patchear Tibia Client 15.25",
			desc:     "Executa client-editor com o local.toml",
			target:   targetClient,
			checkBin: PathClientBin,
		},
		{
			id:      "theme",
			title:   "🎨 Escolher Tema",
			desc:    "Alternar esquema de cores (Dracula, Tokyo, Catppuccin, Gruvbox, etc)",
			isTheme: true,
		},
		{
			id:     "exit",
			title:  "❌ Sair",
			desc:   "Fecha o painel (os servidores continuarão rodando)",
			isExit: true,
		},
	}
}

func initialModel() model {
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")).
		Foreground(lipgloss.Color("#FFFFFF"))

	ti := textinput.New()
	ti.Placeholder = "Digite KILL em maiúsculas"
	ti.CharLimit = 10
	ti.Width = 30

	return model{
		state:        viewMenu,
		menuIndex:    0,
		menuItems:    getMenuItems(),
		logsViewport: vp,
		killInput:    ti,
		themeIndex:   0,
	}
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return doTick()
}

// ============================================================================
// UPDATE
// ============================================================================

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.winHeight = msg.Height

		logsWidth := msg.Width - 8
		if logsWidth < 40 {
			logsWidth = 40
		}
		m.logsViewport.Width = logsWidth

		logsHeight := msg.Height - 8
		if logsHeight < 8 {
			logsHeight = 8
		}
		m.logsViewport.Height = logsHeight

	case tickMsg:
		if m.state == viewService {
			m.refreshServiceLogs()
		}
		cmds = append(cmds, doTick())

	case actionDoneMsg:
		m.actionNotice = msg.notice
		if msg.target == targetClient {
			m.clientLogs = msg.output
		} else if msg.target == targetMySQL {
			m.mysqlLogs = msg.output
		}
		m.refreshServiceLogs()
		m.logsViewport.GotoBottom()
		return m, nil

	case tea.KeyMsg:
		switch m.state {

		// -------------------------------------------------------------
		// ESTADO 1: MENU PRINCIPAL
		// -------------------------------------------------------------
		case viewMenu:
			switch msg.String() {
			case "ctrl+c", "q", "Q":
				return m, tea.Quit
			case "up", "k":
				m.menuIndex--
				if m.menuIndex < 0 {
					m.menuIndex = len(m.menuItems) - 1
				}
				return m, nil
			case "down", "j":
				m.menuIndex++
				if m.menuIndex >= len(m.menuItems) {
					m.menuIndex = 0
				}
				return m, nil
			case "t", "T":
				m.state = viewTheme
				return m, nil
			case "enter":
				selected := m.menuItems[m.menuIndex]
				if selected.isExit {
					return m, tea.Quit
				}
				if selected.isTheme {
					m.state = viewTheme
					return m, nil
				}

				m.activeTarget = selected.target
				m.actionNotice = ""
				m.state = viewService

				switch selected.target {
				case targetCanary:
					if !fileExists(PathCanaryBin) {
						m.actionNotice = "⚠️ O binário canary não existe. Compile antes de rodar!"
						return m, nil
					}
					running, _ := isProcessRunning("canary")
					if !running {
						m.actionNotice = "🚀 Iniciando Canary Server em background..."
						return m, startBackgroundProcess(targetCanary, PathCanaryDir, "canary", PathCanaryLog)
					}
					m.refreshServiceLogs()
					m.logsViewport.GotoBottom()
					return m, nil

				case targetLogin:
					if !fileExists(PathLoginBin) {
						m.actionNotice = "⚠️ O binário login_server não existe. Compile antes de rodar!"
						return m, nil
					}
					running, _ := isProcessRunning("login_server")
					if !running {
						m.actionNotice = "🚀 Iniciando Login Server em background..."
						return m, startBackgroundProcess(targetLogin, PathLoginDir, "login_server", PathLoginLog)
					}
					m.refreshServiceLogs()
					m.logsViewport.GotoBottom()
					return m, nil

				case targetMySQL:
					if !isMySQLRunning() {
						m.actionNotice = "🚀 Iniciando serviço MySQL..."
						return m, manageMySQLNative("start")
					}
					m.refreshServiceLogs()
					return m, nil

				case targetClient:
					if !fileExists(PathClientBin) {
						m.actionNotice = "⚠️ O executável client-editor-linux-x64 não foi encontrado!"
						return m, nil
					}
					m.actionNotice = "⚡ Aplicando patch no Tibia Client 15.25..."
					return m, runClientPatchNative()
				}
			}

		// -------------------------------------------------------------
		// ESTADO 2: VISUALIZADOR DE SERVIÇO & LOGS
		// -------------------------------------------------------------
		case viewService:
			switch msg.String() {
			case "esc", "b", "B":
				m.state = viewMenu
				return m, nil

			case "r", "R":
				m.refreshServiceLogs()
				return m, nil

			case "enter":
				switch m.activeTarget {
				case targetCanary:
					if running, _ := isProcessRunning("canary"); !running {
						m.actionNotice = "🚀 Reiniciando Canary Server..."
						return m, startBackgroundProcess(targetCanary, PathCanaryDir, "canary", PathCanaryLog)
					}
				case targetLogin:
					if running, _ := isProcessRunning("login_server"); !running {
						m.actionNotice = "🚀 Reiniciando Login Server..."
						return m, startBackgroundProcess(targetLogin, PathLoginDir, "login_server", PathLoginLog)
					}
				case targetMySQL:
					if !isMySQLRunning() {
						m.actionNotice = "🚀 Iniciando MySQL..."
						return m, manageMySQLNative("start")
					}
				case targetClient:
					m.actionNotice = "⚡ Reaplicando patch no Tibia Client..."
					return m, runClientPatchNative()
				}

			case "s", "S":
				switch m.activeTarget {
				case targetCanary:
					if running, _ := isProcessRunning("canary"); running {
						m.actionNotice = "🛑 Enviando sinal SIGINT... Salvando dados com segurança..."
						return m, stopProcessGraceful(targetCanary, "canary")
					}
					m.actionNotice = "🚀 Iniciando Canary Server..."
					return m, startBackgroundProcess(targetCanary, PathCanaryDir, "canary", PathCanaryLog)

				case targetLogin:
					if running, _ := isProcessRunning("login_server"); running {
						m.actionNotice = "🛑 Finalizando Login Server com SIGINT..."
						return m, stopProcessGraceful(targetLogin, "login_server")
					}
					m.actionNotice = "🚀 Iniciando Login Server..."
					return m, startBackgroundProcess(targetLogin, PathLoginDir, "login_server", PathLoginLog)

				case targetMySQL:
					if isMySQLRunning() {
						m.actionNotice = "🛑 Desligando serviço MySQL..."
						return m, manageMySQLNative("stop")
					}
					m.actionNotice = "🚀 Iniciando MySQL..."
					return m, manageMySQLNative("start")

				case targetClient:
					m.actionNotice = "⚡ Aplicando patch..."
					return m, runClientPatchNative()
				}

			case "k", "K":
				if m.activeTarget == targetCanary || m.activeTarget == targetLogin || m.activeTarget == targetMySQL {
					m.state = viewKillConfirm
					m.killInput.Reset()
					m.killInput.Focus()
					m.killFeedback = ""
					return m, textinput.Blink
				}
			}

			var cmd tea.Cmd
			m.logsViewport, cmd = m.logsViewport.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		// -------------------------------------------------------------
		// ESTADO 3: SELETOR DE TEMAS
		// -------------------------------------------------------------
		case viewTheme:
			switch msg.String() {
			case "esc", "q", "Q":
				m.state = viewMenu
				return m, nil
			case "up", "k":
				m.themeIndex--
				if m.themeIndex < 0 {
					m.themeIndex = len(themes) - 1
				}
				currentTheme = themes[m.themeIndex]
				return m, nil
			case "down", "j":
				m.themeIndex++
				if m.themeIndex >= len(themes) {
					m.themeIndex = 0
				}
				currentTheme = themes[m.themeIndex]
				return m, nil
			case "enter":
				currentTheme = themes[m.themeIndex]
				m.state = viewMenu
				return m, nil
			}

		// -------------------------------------------------------------
		// ESTADO 4: CONFIRMAÇÃO SCARY DE FORCE KILL
		// -------------------------------------------------------------
		case viewKillConfirm:
			switch msg.String() {
			case "esc":
				m.state = viewService
				return m, nil

			case "enter":
				if strings.TrimSpace(m.killInput.Value()) == "KILL" {
					m.state = viewService
					m.actionNotice = "💀 Processo exterminado com SIGKILL (-9)!"

					switch m.activeTarget {
					case targetCanary:
						return m, killProcessHard(targetCanary, "canary")
					case targetLogin:
						return m, killProcessHard(targetLogin, "login_server")
					case targetMySQL:
						return m, killProcessHard(targetMySQL, "mysqld")
					}
				}
				m.killFeedback = "⚠️ DIGITE EXATAMENTE 'KILL' PARA CONFIRMAR OU [ESC] PARA CANCELAR!"
				return m, nil
			}

			var cmd tea.Cmd
			m.killInput, cmd = m.killInput.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) refreshServiceLogs() {
	switch m.activeTarget {
	case targetCanary:
		m.logsViewport.SetContent(readLastLogs(PathCanaryLog, 200))
	case targetLogin:
		m.logsViewport.SetContent(readLastLogs(PathLoginLog, 200))
	case targetMySQL:
		if m.mysqlLogs != "" {
			m.logsViewport.SetContent(m.mysqlLogs)
		} else {
			out, _ := exec.Command("service", "mysql", "status").CombinedOutput()
			m.logsViewport.SetContent(string(out))
		}
	case targetClient:
		if m.clientLogs != "" {
			m.logsViewport.SetContent(m.clientLogs)
		} else {
			m.logsViewport.SetContent("Pronto para aplicar patch no client 15.25.\nPressione [ENTER] para executar.")
		}
	}
}

// ============================================================================
// VIEWS (100% PREENCHIMENTO SEM NENHUM GAP PRETO)
// ============================================================================

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

	// 1. Status Bar do Topo com todos os fragmentos estilizados em t.Bg
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

	// O interior da caixa ocupa (menuWidth - 4) colunas para fechar perfeitamente com a moldura
	innerBannerW := menuWidth - 4
	quickBarPadded := padLine(quickBarText, innerBannerW, t.Bg)

	topBanner := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		BorderBackground(t.Bg).
		Background(t.Bg).
		Padding(0, 1).
		Render(quickBarPadded)

	// 2. Renderização de cada Item do Menu
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

		// Espaçamento entre opções preenchido 100% com t.Bg
		if i < len(m.menuItems)-1 {
			menuRows = append(menuRows, padLine("", menuWidth, t.Bg))
		}
	}

	// 3. Rodapé
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

	// Fundo preto puro exclusivo para a caixa interna dos logs
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
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erro ao executar runner: %v\n", err)
		os.Exit(1)
	}
}