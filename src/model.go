package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	viewMenu viewState = iota
	viewSetupMenu
	viewEditPaths
	viewModeSelector
	viewPlayersMenu
	viewToolsMenu
	viewService
	viewKillConfirm
	viewTheme
)

type tickMsg time.Time

type actionDoneMsg struct {
	target serviceTarget
	output string
	notice string
}

type model struct {
	state          viewState
	activeTarget   serviceTarget
	menuIndex      int
	menuItems      []menuItem
	setupIndex     int
	setupItems     []subMenuItem
	playersIndex   int
	playersItems   []subMenuItem
	toolsIndex     int
	toolsItems     []subMenuItem
	modeIndex      int
	pathInputs     []textinput.Model
	pathInputIndex int
	logsViewport   viewport.Model
	killInput      textinput.Model
	killFeedback   string
	actionNotice   string
	genericLogs    string
	mysqlLogs      string
	themeIndex     int
	winWidth       int
	winHeight      int
}

func initialModel() model {
	loadConfig()

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")).
		Foreground(lipgloss.Color("#FFFFFF"))

	ti := textinput.New()
	ti.Placeholder = "Digite KILL em maiúsculas"
	ti.CharLimit = 10
	ti.Width = 30

	labels := []string{
		appConfig.CanaryDir,
		appConfig.LoginDir,
		appConfig.ClientEditorDir,
		appConfig.TibiaClientDir,
		appConfig.TibiaClientExe,
	}

	inputs := make([]textinput.Model, 5)
	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].SetValue(labels[i])
		inputs[i].Width = 60
	}
	inputs[0].Focus()

	activeThemeIdx := 0
	for i, t := range themes {
		if t.ID == currentTheme.ID {
			activeThemeIdx = i
			break
		}
	}

	modeIdx := 0
	switch appConfig.Mode {
	case ModeVPS:
		modeIdx = 1
	case ModeLinuxLocal:
		modeIdx = 2
	}

	return model{
		state:        viewMenu,
		menuIndex:    0,
		menuItems:    getMainItems(),
		setupIndex:   0,
		setupItems:   getSetupItems(),
		playersIndex: 0,
		playersItems: getPlayersItems(),
		toolsIndex:   0,
		toolsItems:   getToolsItems(),
		modeIndex:    modeIdx,
		pathInputs:   inputs,
		logsViewport: vp,
		killInput:    ti,
		themeIndex:   activeThemeIdx,
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
		if msg.target == targetMySQL {
			m.mysqlLogs = msg.output
		} else {
			m.genericLogs = msg.output
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
				if selected.isPlayers {
					m.state = viewPlayersMenu
					m.playersIndex = 0
					return m, nil
				}
				if selected.isTools {
					m.state = viewToolsMenu
					m.toolsIndex = 0
					return m, nil
				}
				if selected.isSetup {
					m.state = viewSetupMenu
					m.setupIndex = 0
					m.setupItems = getSetupItems()
					return m, nil
				}

				if selected.target == targetLaunchClient {
					return m, launchTibiaClientWindows()
				}

				m.activeTarget = selected.target
				m.actionNotice = ""
				m.state = viewService

				switch selected.target {
				case targetCanary:
					canaryBin := filepath.Join(expandHome(appConfig.CanaryDir), "canary")
					if !fileExists(canaryBin) {
						m.actionNotice = "⚠️ O binário canary não existe. Execute o Setup ou compile antes de rodar!"
						return m, nil
					}
					running, _ := isProcessRunning("canary")
					if !running {
						m.actionNotice = "🚀 Iniciando Canary Server em background..."
						return m, startBackgroundProcess(targetCanary, appConfig.CanaryDir, "canary", PathCanaryLog)
					}
					m.refreshServiceLogs()
					m.logsViewport.GotoBottom()
					return m, nil

				case targetLogin:
					loginBin := filepath.Join(expandHome(appConfig.LoginDir), "login_server")
					if !fileExists(loginBin) {
						m.actionNotice = "⚠️ O binário login_server não existe. Execute o Setup ou compile antes de rodar!"
						return m, nil
					}
					running, _ := isProcessRunning("login_server")
					if !running {
						m.actionNotice = "🚀 Iniciando Login Server em background..."
						return m, startBackgroundProcess(targetLogin, appConfig.LoginDir, "login_server", PathLoginLog)
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
				}
			}

		// -------------------------------------------------------------
		// ESTADO 2: SUBMENU DE SETUP & FERRAMENTAS
		// -------------------------------------------------------------
		case viewSetupMenu:
			switch msg.String() {
			case "esc", "b", "B":
				m.state = viewMenu
				return m, nil
			case "up", "k":
				m.setupIndex--
				if m.setupIndex < 0 {
					m.setupIndex = len(m.setupItems) - 1
				}
				return m, nil
			case "down", "j":
				m.setupIndex++
				if m.setupIndex >= len(m.setupItems) {
					m.setupIndex = 0
				}
				return m, nil
			case "enter":
				selected := m.setupItems[m.setupIndex]
				if selected.isBack {
					m.state = viewMenu
					return m, nil
				}
				if selected.isModePick {
					m.state = viewModeSelector
					return m, nil
				}
				if selected.isEditPath {
					m.state = viewEditPaths
					m.pathInputIndex = 0
					m.pathInputs[0].SetValue(appConfig.CanaryDir)
					m.pathInputs[1].SetValue(appConfig.LoginDir)
					m.pathInputs[2].SetValue(appConfig.ClientEditorDir)
					m.pathInputs[3].SetValue(appConfig.TibiaClientDir)
					m.pathInputs[4].SetValue(appConfig.TibiaClientExe)
					for i := range m.pathInputs {
						m.pathInputs[i].Blur()
					}
					m.pathInputs[0].Focus()
					return m, textinput.Blink
				}

				m.activeTarget = selected.target
				m.actionNotice = ""
				m.state = viewService

				switch selected.target {
				case targetFullSetup:
					m.actionNotice = "⚡ Executando instalação automática completa..."
					return m, runFullSetup1ClickNative()

				case targetSyncConnections:
					m.actionNotice = "🔄 Sincronizando endereços de conexão..."
					return m, runSyncConnectionsNative()

				case targetClientPatch:
					clientEditBin := filepath.Join(expandHome(appConfig.ClientEditorDir), "client-editor-linux-x64")
					if !fileExists(clientEditBin) {
						m.actionNotice = "⚠️ O executável client-editor-linux-x64 não foi encontrado!"
						return m, nil
					}
					m.actionNotice = "⚡ Aplicando patch no Tibia Client 15.25..."
					return m, runClientPatchNative()

				case targetDownloadClient:
					m.actionNotice = "📥 Baixando e extraindo Tibia Client 15.25..."
					return m, runDownloadClientNative()
				}
			}

		// -------------------------------------------------------------
		// ESTADO 3: SUBMENU DE JOGADORES & CONTAS (PREVIEW)
		// -------------------------------------------------------------
		case viewPlayersMenu:
			switch msg.String() {
			case "esc", "b", "B":
				m.state = viewMenu
				return m, nil
			case "up", "k":
				m.playersIndex--
				if m.playersIndex < 0 {
					m.playersIndex = len(m.playersItems) - 1
				}
				return m, nil
			case "down", "j":
				m.playersIndex++
				if m.playersIndex >= len(m.playersItems) {
					m.playersIndex = 0
				}
				return m, nil
			case "enter":
				selected := m.playersItems[m.playersIndex]
				if selected.isBack {
					m.state = viewMenu
					return m, nil
				}
				// Entra na tela de log/serviço exibindo preview informativo
				m.activeTarget = targetPlaceholder
				m.state = viewService
				m.genericLogs = fmt.Sprintf(
					"\033[1;36m=== PAINEL DE ADMINISTRAÇÃO: %s ===\033[0m\n\n"+
						"Ação selecionada: \033[1;32m%s\033[0m\n"+
						"Descrição: %s\n\n"+
						"\033[0;33m[ℹ️  Módulo mapeado e pronto para integração de banco de dados SQL]\033[0m\n",
					selected.title, selected.title, selected.desc,
				)
				m.actionNotice = "ℹ️ Formulário de administração selecionado."
				m.refreshServiceLogs()
				return m, nil
			}

		// -------------------------------------------------------------
		// ESTADO 4: SUBMENU DE FERRAMENTAS & CALIBRAÇÃO (PREVIEW)
		// -------------------------------------------------------------
		case viewToolsMenu:
			switch msg.String() {
			case "esc", "b", "B":
				m.state = viewMenu
				return m, nil
			case "up", "k":
				m.toolsIndex--
				if m.toolsIndex < 0 {
					m.toolsIndex = len(m.toolsItems) - 1
				}
				return m, nil
			case "down", "j":
				m.toolsIndex++
				if m.toolsIndex >= len(m.toolsItems) {
					m.toolsIndex = 0
				}
				return m, nil
			case "enter":
				selected := m.toolsItems[m.toolsIndex]
				if selected.isBack {
					m.state = viewMenu
					return m, nil
				}
				m.activeTarget = targetPlaceholder
				m.state = viewService
				m.genericLogs = fmt.Sprintf(
					"\033[1;36m=== FERRAMENTA DE CALIBRAÇÃO & CRIAÇÃO ===\033[0m\n\n"+
						"Ferramenta: \033[1;32m%s\033[0m\n"+
						"Finalidade: %s\n\n"+
						"\033[0;33m[ℹ️  Gerador de scripts Canary pronto para entrada de parâmetros]\033[0m\n",
					selected.title, selected.desc,
				)
				m.actionNotice = "ℹ️ Ferramenta de criação selecionada."
				m.refreshServiceLogs()
				return m, nil
			}

		// -------------------------------------------------------------
		// ESTADO 5: SELETOR DE MODO DE USO (WSL / VPS / LOCAL)
		// -------------------------------------------------------------
		case viewModeSelector:
			switch msg.String() {
			case "esc":
				m.state = viewSetupMenu
				return m, nil
			case "up", "k":
				m.modeIndex--
				if m.modeIndex < 0 {
					m.modeIndex = 2
				}
				return m, nil
			case "down", "j":
				m.modeIndex++
				if m.modeIndex > 2 {
					m.modeIndex = 0
				}
				return m, nil
			case "enter":
				switch m.modeIndex {
				case 0:
					appConfig.Mode = ModeWSL
				case 1:
					appConfig.Mode = ModeVPS
				case 2:
					appConfig.Mode = ModeLinuxLocal
				}
				saveConfig()
				m.setupItems = getSetupItems() // Atualiza o menu de setup dinamicamente
				m.state = viewSetupMenu
				return m, nil
			}

		// -------------------------------------------------------------
		// ESTADO 6: FORMULÁRIO DE EDIÇÃO DE PASTAS
		// -------------------------------------------------------------
		case viewEditPaths:
			switch msg.String() {
			case "esc":
				m.state = viewSetupMenu
				return m, nil
			case "tab", "down":
				m.pathInputs[m.pathInputIndex].Blur()
				m.pathInputIndex = (m.pathInputIndex + 1) % len(m.pathInputs)
				m.pathInputs[m.pathInputIndex].Focus()
				return m, textinput.Blink
			case "shift+tab", "up":
				m.pathInputs[m.pathInputIndex].Blur()
				m.pathInputIndex = (m.pathInputIndex - 1 + len(m.pathInputs)) % len(m.pathInputs)
				m.pathInputs[m.pathInputIndex].Focus()
				return m, textinput.Blink
			case "enter":
				appConfig.CanaryDir = m.pathInputs[0].Value()
				appConfig.LoginDir = m.pathInputs[1].Value()
				appConfig.ClientEditorDir = m.pathInputs[2].Value()
				appConfig.TibiaClientDir = m.pathInputs[3].Value()
				appConfig.TibiaClientExe = m.pathInputs[4].Value()
				saveConfig()
				m.state = viewSetupMenu
				return m, nil
			}

			var cmd tea.Cmd
			m.pathInputs[m.pathInputIndex], cmd = m.pathInputs[m.pathInputIndex].Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		// -------------------------------------------------------------
		// ESTADO 7: VISUALIZADOR DE SERVIÇO & LOGS
		// -------------------------------------------------------------
		case viewService:
			switch msg.String() {
			case "esc", "b", "B":
				if m.activeTarget == targetFullSetup || m.activeTarget == targetClientPatch || m.activeTarget == targetDownloadClient || m.activeTarget == targetSyncConnections {
					m.state = viewSetupMenu
				} else if m.activeTarget == targetPlaceholder {
					m.state = viewMenu
				} else {
					m.state = viewMenu
				}
				return m, nil

			case "r", "R":
				m.refreshServiceLogs()
				return m, nil

			case "enter":
				switch m.activeTarget {
				case targetCanary:
					if running, _ := isProcessRunning("canary"); !running {
						m.actionNotice = "🚀 Reiniciando Canary Server..."
						return m, startBackgroundProcess(targetCanary, appConfig.CanaryDir, "canary", PathCanaryLog)
					}
				case targetLogin:
					if running, _ := isProcessRunning("login_server"); !running {
						m.actionNotice = "🚀 Reiniciando Login Server..."
						return m, startBackgroundProcess(targetLogin, appConfig.LoginDir, "login_server", PathLoginLog)
					}
				case targetMySQL:
					if !isMySQLRunning() {
						m.actionNotice = "🚀 Iniciando MySQL..."
						return m, manageMySQLNative("start")
					}
				case targetClientPatch:
					m.actionNotice = "⚡ Reaplicando patch no Tibia Client..."
					return m, runClientPatchNative()
				case targetFullSetup:
					m.actionNotice = "⚡ Reexecutando instalação completa..."
					return m, runFullSetup1ClickNative()
				case targetDownloadClient:
					m.actionNotice = "📥 Baixando client novamente..."
					return m, runDownloadClientNative()
				case targetSyncConnections:
					m.actionNotice = "🔄 Sincronizando conexões novamente..."
					return m, runSyncConnectionsNative()
				}

			case "s", "S":
				switch m.activeTarget {
				case targetCanary:
					if running, _ := isProcessRunning("canary"); running {
						m.actionNotice = "🛑 Enviando sinal SIGINT... Salvando dados com segurança..."
						return m, stopProcessGraceful(targetCanary, "canary")
					}
					m.actionNotice = "🚀 Iniciando Canary Server..."
					return m, startBackgroundProcess(targetCanary, appConfig.CanaryDir, "canary", PathCanaryLog)

				case targetLogin:
					if running, _ := isProcessRunning("login_server"); running {
						m.actionNotice = "🛑 Finalizando Login Server com SIGINT..."
						return m, stopProcessGraceful(targetLogin, "login_server")
					}
					m.actionNotice = "🚀 Iniciando Login Server..."
					return m, startBackgroundProcess(targetLogin, appConfig.LoginDir, "login_server", PathLoginLog)

				case targetMySQL:
					if isMySQLRunning() {
						m.actionNotice = "🛑 Desligando serviço MySQL..."
						return m, manageMySQLNative("stop")
					}
					m.actionNotice = "🚀 Iniciando MySQL..."
					return m, manageMySQLNative("start")
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
		// ESTADO 8: SELETOR DE TEMAS
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
				saveConfig()
				m.state = viewMenu
				return m, nil
			}

		// -------------------------------------------------------------
		// ESTADO 9: CONFIRMAÇÃO SCARY DE FORCE KILL
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
			out := getMySQLStatusOutput()
			m.logsViewport.SetContent(out)
		}
	default:
		if m.genericLogs != "" {
			m.logsViewport.SetContent(m.genericLogs)
		} else {
			m.logsViewport.SetContent("Aguardando execução da ação...")
		}
	}
}