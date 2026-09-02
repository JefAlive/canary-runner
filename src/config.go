package main

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