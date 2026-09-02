package main

import (
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
			m.logsViewport.SetContent(getMySQLStatusOutput())
		}
	case targetClient:
		if m.clientLogs != "" {
			m.logsViewport.SetContent(m.clientLogs)
		} else {
			m.logsViewport.SetContent("Pronto para aplicar patch no client 15.25.\nPressione [ENTER] para executar.")
		}
	}
}