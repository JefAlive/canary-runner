package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	maxMenuWidth = 84

	// URLs Oficiais Fixas (Padrão de fábrica HTTPS - sem senha)
	RepoCanaryURL     = "https://github.com/opentibiabr/canary.git"
	RepoLoginURL      = "https://github.com/opentibiabr/login-server.git"
	RepoClientEditURL = "https://github.com/opentibiabr/client-editor.git"
	ClientZipURL      = "https://github.com/dudantas/tibia-client/releases/download/15.25.0a00a0/tibia-client-15.25.0a00a0.zip"

	PathCanaryLog = "~/Code/canary/canary.log"
	PathLoginLog  = "~/Code/login-server/login_server.log"
)

type AppMode string

const (
	ModeWSL        AppMode = "wsl"
	ModeVPS        AppMode = "vps"
	ModeLinuxLocal AppMode = "local"
)

type AppConfig struct {
	Theme           string  `json:"theme"`
	Mode            AppMode `json:"mode"`
	PublicHost      string  `json:"public_host"`
	MySQLHost       string  `json:"mysql_host"`
	MySQLPort       string  `json:"mysql_port"`
	MySQLUser       string  `json:"mysql_user"`
	MySQLPass       string  `json:"mysql_pass"`
	MySQLDatabase   string  `json:"mysql_database"`
	CanaryDir       string  `json:"canary_dir"`
	LoginDir        string  `json:"login_dir"`
	ClientEditorDir string  `json:"client_editor_dir"`
	TibiaClientDir  string  `json:"tibia_client_dir"`
	TibiaClientExe  string  `json:"tibia_client_exe"`
}

var appConfig AppConfig

func getDefaultConfig() AppConfig {
	defaultClientExe := "~/Code/tibia-client/bin/client.exe"
	detectedMode := ModeLinuxLocal

	if _, err := os.Stat("/mnt/c/Users"); err == nil {
		detectedMode = ModeWSL
		entries, _ := os.ReadDir("/mnt/c/Users")
		for _, e := range entries {
			if e.IsDir() && e.Name() != "Public" && e.Name() != "Default" && e.Name() != "All Users" {
				defaultClientExe = fmt.Sprintf("/mnt/c/Users/%s/Code/tibia-client/bin/client.exe", e.Name())
				break
			}
		}
	}

	return AppConfig{
		Theme:           "dracula",
		Mode:            detectedMode,
		PublicHost:      "127.0.0.1",
		MySQLHost:       "127.0.0.1",
		MySQLPort:       "3306",
		MySQLUser:       "root",
		MySQLPass:       "",
		MySQLDatabase:   "canary",
		CanaryDir:       "~/Code/canary",
		LoginDir:        "~/Code/login-server",
		ClientEditorDir: "~/Code/client-editor",
		TibiaClientDir:  "~/Code/tibia-client",
		TibiaClientExe:  defaultClientExe,
	}
}

func getConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "canary-runner.json"
	}
	configDir := filepath.Join(home, ".config", "canary-runner")
	_ = os.MkdirAll(configDir, 0755)
	return filepath.Join(configDir, "config.json")
}

func loadConfig() {
	appConfig = getDefaultConfig()
	configFile := getConfigFilePath()

	data, err := os.ReadFile(configFile)
	if err == nil {
		_ = json.Unmarshal(data, &appConfig)
	} else {
		saveConfig()
	}

	for i, t := range themes {
		if t.ID == appConfig.Theme {
			currentTheme = themes[i]
			break
		}
	}
}

func saveConfig() {
	appConfig.Theme = currentTheme.ID
	configFile := getConfigFilePath()
	data, err := json.MarshalIndent(appConfig, "", "  ")
	if err == nil {
		_ = os.WriteFile(configFile, data, 0644)
	}
}

// ============================================================================
// ESTRUTURAS DE NAVEGAÇÃO DOS MENUS
// ============================================================================

type serviceTarget int

const (
	targetCanary serviceTarget = iota
	targetLogin
	targetMySQL
	targetLaunchClient
	targetClientPatch
	targetFullSetup
	targetDownloadClient
	targetSyncConnections
	targetPlaceholder
)

type menuItem struct {
	id         string
	title      string
	desc       string
	target     serviceTarget
	procName   string
	isPlayers  bool
	isTools    bool
	isSetup    bool
	isTheme    bool
	isExit     bool
	isMySQL    bool
}

type subMenuItem struct {
	id         string
	title      string
	desc       string
	target     serviceTarget
	isAction   bool
	isModePick bool
	isEditPath bool
	isBack     bool
}

func getMainItems() []menuItem {
	return []menuItem{
		{
			id:       "canary",
			title:    "🚀 Canary Server (Servidor do Jogo)",
			desc:     "Liga o mundo do jogo, monstros, mapa, magias e sistemas",
			target:   targetCanary,
			procName: "canary",
		},
		{
			id:       "login",
			title:    "🔑 Login Server (Servidor de Contas)",
			desc:     "Liga o serviço que valida senhas e lista os personagens",
			target:   targetLogin,
			procName: "login_server",
		},
		{
			id:      "mysql",
			title:   "🗄️ MySQL Database (Banco de Dados)",
			desc:    "Liga onde ficam salvas as contas, personagens, itens e casas",
			target:  targetMySQL,
			isMySQL: true,
		},
		{
			id:     "launch_client",
			title:  "🎮 Abrir Tibia no Windows",
			desc:   "Dispara o jogo no Windows direto pelo terminal sem precisar caçar a pasta",
			target: targetLaunchClient,
		},
		{
			id:        "players",
			title:     "👥 Gerenciar Jogadores & Contas",
			desc:      "Painel para criar contas, criar personagens, GOD e adicionar Tibia Coins",
			isPlayers: true,
		},
		{
			id:      "tools",
			title:   "🛠️ Ferramentas & Calibração Global",
			desc:    "Criar monstros, bosses, alavancas, echo raids, itens, NPCs e calibrar magias",
			isTools: true,
		},
		{
			id:      "setup",
			title:   "⚙️ Setup & Conexões do Servidor",
			desc:    "Instalação 1-click, sincronizar IPs, baixar client, pastas e modo de uso",
			isSetup: true,
		},
		{
			id:      "theme",
			title:   "🎨 Escolher Tema",
			desc:    "Alterne as cores visuais do painel (Dracula, Tokyo, Catppuccin, Tron, etc.)",
			isTheme: true,
		},
		{
			id:     "exit",
			title:  "❌ Sair",
			desc:   "Fecha este painel (os servidores que estiverem ligados continuam rodando)",
			isExit: true,
		},
	}
}

// Submenu contextual de Setup
func getSetupItems() []subMenuItem {
	switch appConfig.Mode {
	case ModeVPS:
		return []subMenuItem{
			{
				id:     "full_setup",
				title:  "⚡ Instalação Automática Completa (1-Click)",
				desc:   "Baixa e instala todos os arquivos do servidor nesta máquina",
				target: targetFullSetup,
			},
			{
				id:     "sync_conn",
				title:  "🔄 Sincronizar Conexões do Servidor",
				desc:   "Aplica o endereço da internet e o banco de dados em todos os serviços",
				target: targetSyncConnections,
			},
			{
				id:     "patch_client",
				title:  "🎮 Preparar o Jogo (Client) para os Jogadores",
				desc:   "Configura o jogo com o endereço da internet para enviar aos seus amigos",
				target: targetClientPatch,
			},
			{
				id:     "download_client",
				title:  "📥 Baixar e Extrair Tibia Client 15.25",
				desc:   "Baixa a release oficial do jogo se você ainda não tiver",
				target: targetDownloadClient,
			},
			{
				id:         "edit_paths",
				title:      "📁 Alterar Pastas dos Arquivos",
				desc:       "Apenas se os arquivos do servidor estiverem em pastas personalizadas",
				isEditPath: true,
			},
			{
				id:         "change_mode",
				title:      "🌐 Trocar Modo de Uso (Ativo: 🌐 Servidor na Nuvem)",
				desc:       "Mude para o modo Windows/WSL ou Linux Local",
				isModePick: true,
			},
			{
				id:     "back",
				title:  "❮ Voltar ao Menu Principal",
				desc:   "Retorna à tela inicial do runner",
				isBack: true,
			},
		}

	case ModeLinuxLocal:
		return []subMenuItem{
			{
				id:     "full_setup",
				title:  "⚡ Instalação Automática Completa (1-Click)",
				desc:   "Baixa e prepara todo o ambiente no seu computador automaticamente",
				target: targetFullSetup,
			},
			{
				id:     "sync_conn",
				title:  "🔄 Sincronizar Conexão Local",
				desc:   "Garante que todos os serviços conversem em 127.0.0.1 no seu computador",
				target: targetSyncConnections,
			},
			{
				id:     "patch_client",
				title:  "🎮 Configurar o Tibia (Client)",
				desc:   "Prepara o jogo para conectar no servidor deste computador",
				target: targetClientPatch,
			},
			{
				id:     "download_client",
				title:  "📥 Baixar e Extrair Tibia Client 15.25",
				desc:   "Baixa os arquivos do jogo prontos para jogar",
				target: targetDownloadClient,
			},
			{
				id:         "edit_paths",
				title:      "📁 Alterar Pastas dos Arquivos",
				desc:       "Apenas se você já tiver os arquivos em pastas personalizadas",
				isEditPath: true,
			},
			{
				id:         "change_mode",
				title:      "🌐 Trocar Modo de Uso (Ativo: 💻 Linux Local)",
				desc:       "Mude se você estiver em Windows/WSL ou em um servidor na nuvem",
				isModePick: true,
			},
			{
				id:     "back",
				title:  "❮ Voltar ao Menu Principal",
				desc:   "Retorna à tela inicial do runner",
				isBack: true,
			},
		}

	default: // ModeWSL
		return []subMenuItem{
			{
				id:     "full_setup",
				title:  "⚡ Instalação Automática Completa (1-Click)",
				desc:   "Baixa todos os arquivos do servidor, baixa o jogo e deixa tudo pronto",
				target: targetFullSetup,
			},
			{
				id:     "sync_conn",
				title:  "🔄 Sincronizar Conexão do Jogo com o Windows",
				desc:   "Ajusta o servidor e o jogo com o endereço atual para você conseguir entrar",
				target: targetSyncConnections,
			},
			{
				id:     "patch_client",
				title:  "🎮 Configurar o Tibia (Client.exe)",
				desc:   "Prepara o jogo para abrir e conectar direto no seu servidor",
				target: targetClientPatch,
			},
			{
				id:     "download_client",
				title:  "📥 Apenas Baixar o Tibia 15.25",
				desc:   "Baixa os arquivos do jogo prontos para jogar se você ainda não tiver",
				target: targetDownloadClient,
			},
			{
				id:         "edit_paths",
				title:      "📁 Alterar Pastas dos Arquivos",
				desc:       "Apenas para quem já baixou ou moveu as pastas do servidor manualmente",
				isEditPath: true,
			},
			{
				id:         "change_mode",
				title:      "🌐 Trocar Modo de Uso (Ativo: 🖥️ Windows/WSL)",
				desc:       "Mude se você estiver rodando em uma máquina na nuvem ou Linux puro",
				isModePick: true,
			},
			{
				id:     "back",
				title:  "❮ Voltar ao Menu Principal",
				desc:   "Retorna à tela inicial do runner",
				isBack: true,
			},
		}
	}
}

// Submenu de Gerenciamento de Jogadores & Contas
func getPlayersItems() []subMenuItem {
	return []subMenuItem{
		{
			id:       "create_acc",
			title:    "👤 Criar Nova Conta",
			desc:     "Cria uma nova conta de acesso (Login e Senha) no banco de dados",
			isAction: true,
		},
		{
			id:       "create_char",
			title:    "🧙‍♂️ Criar Novo Personagem",
			desc:     "Cria um personagem escolhendo Nome, Vocação, Sexo e Nível inicial",
			isAction: true,
		},
		{
			id:       "promote_god",
			title:    "⭐ Promover Personagem para GOD (Administrador)",
			desc:     "Dá poderes de administrador a um personagem para usar comandos no jogo",
			isAction: true,
		},
		{
			id:       "add_coins",
			title:    "💰 Adicionar Tibia Coins na Conta",
			desc:     "Adiciona moedas na Store do jogo para a conta de um jogador",
			isAction: true,
		},
		{
			id:       "list_players",
			title:    "📋 Listar Jogadores & Contas",
			desc:     "Exibe uma lista com todas as contas, personagens criados, níveis e cargos",
			isAction: true,
		},
		{
			id:       "delete_char",
			title:    "🗑️ Deletar Personagem ou Conta",
			desc:     "Remove um personagem ou conta do banco de dados com segurança",
			isAction: true,
		},
		{
			id:     "back",
			title:  "❮ Voltar ao Menu Principal",
			desc:   "Retorna à tela inicial do runner",
			isBack: true,
		},
	}
}

// Submenu de Ferramentas & Calibração Global (Itens 1 a 8 + 14)
func getToolsItems() []subMenuItem {
	return []subMenuItem{
		{
			id:       "tool_monster",
			title:    "🐉 Criar / Atualizar Monstro ou Boss Oficial",
			desc:     "Gera arquivo do monstro com Vida, XP, ataques, imunidades e tabela de Loot",
			isAction: true,
		},
		{
			id:       "tool_lever",
			title:    "🕹️ Criar Alavanca de Sala de Boss",
			desc:     "Gera a mecânica oficial de sala: até 5 jogadores, teleporte, tempo limite e spawn",
			isAction: true,
		},
		{
			id:       "tool_echo",
			title:    "🌀 Criar / Editar Echo Raid",
			desc:     "Configura o sistema oficial de zonas de Eco com ondas de criaturas sombrias",
			isAction: true,
		},
		{
			id:       "tool_raid",
			title:    "🚨 Criar Invasão / Raid Padrão",
			desc:     "Configura invasões com horários e avisos automáticos no chat global",
			isAction: true,
		},
		{
			id:       "tool_item",
			title:    "⚔️ Cadastrar Equipamento ou Item Oficial",
			desc:     "Gera armas, armaduras e itens com atributos, imbuements e resistências",
			isAction: true,
		},
		{
			id:       "tool_npc",
			title:    "🗣️ Criar NPC de Loja ou Rota de Viagem",
			desc:     "Cria NPCs de comércio com lista de compra/venda ou barqueiros de viagem",
			isAction: true,
		},
		{
			id:       "tool_quest",
			title:    "📜 Criar Baú de Quest / Recompensa Oficial",
			desc:     "Gera baús com level mínimo, chaves, storages e itens de recompensa",
			isAction: true,
		},
		{
			id:       "tool_mount",
			title:    "🐎 Cadastrar Montaria ou Outfit Oficial",
			desc:     "Adiciona montarias com item de domar (taming) e novos visuais de personagens",
			isAction: true,
		},
		{
			id:       "tool_spells",
			title:    "⚖️ Ajustar Fórmulas e Cooldowns de Magias",
			desc:     "Calibra dano, cura, custo de mana e tempo de recarga das magias oficiais",
			isAction: true,
		},
		{
			id:     "back",
			title:  "❮ Voltar ao Menu Principal",
			desc:   "Retorna à tela inicial do runner",
			isBack: true,
		},
	}
}