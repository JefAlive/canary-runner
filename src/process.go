package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

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

func getMySQLStatusOutput() string {
	out, _ := exec.Command("service", "mysql", "status").CombinedOutput()
	return string(out)
}

// Detecta o IP do WSL dinamicamente
func getWSLIP() string {
	cmd := exec.Command("bash", "-c", "hostname -I | awk '{print $1}'")
	out, err := cmd.Output()
	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		return "127.0.0.1"
	}
	return strings.TrimSpace(string(out))
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

	var filtered []string
	for _, line := range strings.Split(cleaned, "\n") {
		if strings.HasPrefix(line, "Script started on") || strings.HasPrefix(line, "Script done on") {
			continue
		}
		filtered = append(filtered, line)
	}

	return strings.Join(filtered, "\n")
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

// Dispara o client.exe diretamente no Windows através do WSL
func launchTibiaClientWindows() tea.Cmd {
	return func() tea.Msg {
		exePath := expandHome(appConfig.TibiaClientExe)
		if !fileExists(appConfig.TibiaClientExe) {
			return actionDoneMsg{
				target: targetLaunchClient,
				notice: "⚠️ Executável do client não encontrado no caminho configurado!",
			}
		}

		// Se estiver no WSL, converte o caminho para formato Windows e dispara via cmd.exe
		if strings.HasPrefix(exePath, "/mnt/") {
			driveLetter := strings.ToUpper(string(exePath[5]))
			windowsPath := driveLetter + ":" + strings.ReplaceAll(exePath[6:], "/", "\\")
			_ = exec.Command("cmd.exe", "/c", "start", "", windowsPath).Start()
		} else {
			_ = exec.Command("nohup", exePath).Start()
		}

		return actionDoneMsg{
			target: targetLaunchClient,
			notice: "🚀 Tibia Client disparado na tela do Windows!",
		}
	}
}

// Roda o patcher do client com o local.toml
func runClientPatchNative() tea.Cmd {
	return func() tea.Msg {
		expandedWorkDir := expandHome(appConfig.ClientEditorDir)
		expandedExe := expandHome(appConfig.TibiaClientExe)
		script := fmt.Sprintf("export CLICOLOR_FORCE=1 FORCE_COLOR=1; cd '%s' && ./client-editor-linux-x64 edit -t '%s' -c local.toml", expandedWorkDir, expandedExe)
		out, err := exec.Command("bash", "-c", script).CombinedOutput()
		outputStr := strings.ReplaceAll(string(out), "\r", "")
		if err != nil && len(strings.TrimSpace(outputStr)) == 0 {
			outputStr = "❌ Erro ao aplicar patch: " + err.Error()
		}
		time.Sleep(400 * time.Millisecond)
		return actionDoneMsg{
			target: targetClientPatch,
			output: outputStr,
			notice: "✔ Patch do client aplicado com sucesso!",
		}
	}
}

// Sincroniza IPs no local.toml do client-editor e no login-server
func runSyncConnectionsNative() tea.Cmd {
	return func() tea.Msg {
		var logs strings.Builder
		targetHost := "127.0.0.1"

		if appConfig.Mode == ModeWSL {
			targetHost = getWSLIP()
		} else if appConfig.Mode == ModeVPS && appConfig.PublicHost != "" {
			targetHost = appConfig.PublicHost
		}

		logs.WriteString(fmt.Sprintf("\033[1;36m=== SINCRONIZANDO CONEXÕES DE REDE ===\033[0m\n\n"))
		logs.WriteString(fmt.Sprintf("🌐 Endereço Alvo: \033[1;32m%s\033[0m (Modo: %s)\n\n", targetHost, appConfig.Mode))

		// 1. Atualiza local.toml no client-editor se existir
		localTomlPath := filepath.Join(expandHome(appConfig.ClientEditorDir), "local.toml")
		if fileExists(localTomlPath) {
			content, err := os.ReadFile(localTomlPath)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				for i, line := range lines {
					if strings.HasPrefix(strings.TrimSpace(line), "ip =") || strings.HasPrefix(strings.TrimSpace(line), "host =") {
						lines[i] = fmt.Sprintf("ip = \"%s\"", targetHost)
					}
				}
				_ = os.WriteFile(localTomlPath, []byte(strings.Join(lines, "\n")), 0644)
				logs.WriteString(fmt.Sprintf("✔ local.toml atualizado com o endereço %s\n", targetHost))
			}
		} else {
			logs.WriteString("ℹ️ Arquivo local.toml ainda não existe (será criado no primeiro patch).\n")
		}

		// 2. Atualiza config.toml no login-server se existir
		loginTomlPath := filepath.Join(expandHome(appConfig.LoginDir), "config.toml")
		if fileExists(loginTomlPath) {
			content, err := os.ReadFile(loginTomlPath)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				for i, line := range lines {
					if strings.HasPrefix(strings.TrimSpace(line), "ip =") {
						lines[i] = fmt.Sprintf("ip = \"%s\"", targetHost)
					}
				}
				_ = os.WriteFile(loginTomlPath, []byte(strings.Join(lines, "\n")), 0644)
				logs.WriteString(fmt.Sprintf("✔ login-server config.toml atualizado com o endereço %s\n", targetHost))
			}
		} else {
			logs.WriteString("ℹ️ config.toml do login-server não encontrado.\n")
		}

		logs.WriteString("\n\033[1;32m✔ Sincronização concluída com sucesso!\033[0m\n")

		return actionDoneMsg{
			target: targetSyncConnections,
			output: logs.String(),
			notice: "✔ Conexões de rede sincronizadas!",
		}
	}
}

// Executa a Instalação 1-Click
func runFullSetup1ClickNative() tea.Cmd {
	return func() tea.Msg {
		var logs strings.Builder
		logs.WriteString("\033[1;36m=== INICIANDO INSTALAÇÃO AUTOMÁTICA COMPLETA ===\033[0m\n\n")

		canaryDir := expandHome(appConfig.CanaryDir)
		loginDir := expandHome(appConfig.LoginDir)
		editorDir := expandHome(appConfig.ClientEditorDir)
		clientDir := expandHome(appConfig.TibiaClientDir)

		// 1. Canary
		if !fileExists(appConfig.CanaryDir) {
			logs.WriteString(fmt.Sprintf("📥 Clonando Canary em %s...\n", canaryDir))
			out, err := exec.Command("git", "clone", RepoCanaryURL, canaryDir).CombinedOutput()
			if err != nil {
				logs.WriteString(fmt.Sprintf("❌ Erro no Canary: %s\n", string(out)))
			} else {
				logs.WriteString("✔ Canary clonado com sucesso!\n")
			}
		} else {
			logs.WriteString("✔ Pasta do Canary já existe.\n")
		}

		// 2. Login Server
		if !fileExists(appConfig.LoginDir) {
			logs.WriteString(fmt.Sprintf("\n📥 Clonando Login Server em %s...\n", loginDir))
			out, err := exec.Command("git", "clone", RepoLoginURL, loginDir).CombinedOutput()
			if err != nil {
				logs.WriteString(fmt.Sprintf("❌ Erro no Login Server: %s\n", string(out)))
			} else {
				logs.WriteString("✔ Login Server clonado com sucesso!\n")
			}
		} else {
			logs.WriteString("✔ Pasta do Login Server já existe.\n")
		}

		// 3. Client Editor
		if !fileExists(appConfig.ClientEditorDir) {
			logs.WriteString(fmt.Sprintf("\n📥 Clonando Client Editor em %s...\n", editorDir))
			out, err := exec.Command("git", "clone", RepoClientEditURL, editorDir).CombinedOutput()
			if err != nil {
				logs.WriteString(fmt.Sprintf("❌ Erro no Client Editor: %s\n", string(out)))
			} else {
				logs.WriteString("✔ Client Editor clonado com sucesso!\n")
			}
		} else {
			logs.WriteString("✔ Pasta do Client Editor já existe.\n")
		}

		// 4. Download do Client 15.25
		logs.WriteString(fmt.Sprintf("\n📥 Baixando Tibia Client 15.25 em %s...\n", clientDir))
		_ = os.MkdirAll(clientDir, 0755)
		zipPath := filepath.Join(clientDir, "tibia-client-15.25.zip")

		dlCmd := exec.Command("curl", "-L", "-o", zipPath, ClientZipURL)
		out, err := dlCmd.CombinedOutput()
		if err != nil {
			logs.WriteString(fmt.Sprintf("❌ Erro ao baixar client: %s\n", string(out)))
		} else {
			logs.WriteString("📦 Extraindo client...\n")
			unzipCmd := exec.Command("unzip", "-o", zipPath, "-d", clientDir)
			unzipOut, unzipErr := unzipCmd.CombinedOutput()
			if unzipErr != nil {
				logs.WriteString(fmt.Sprintf("❌ Erro ao extrair: %s\n", string(unzipOut)))
			} else {
				_ = os.Remove(zipPath)
				logs.WriteString("✔ Tibia Client 15.25 extraído e pronto para uso!\n")
			}
		}

		logs.WriteString("\n\033[1;32m=== SETUP CONCLUÍDO COM SUCESSO! ===\033[0m\n")

		return actionDoneMsg{
			target: targetFullSetup,
			output: logs.String(),
			notice: "✔ Instalação completa finalizada!",
		}
	}
}

// Download e extração do client isolado
func runDownloadClientNative() tea.Cmd {
	return func() tea.Msg {
		var logs strings.Builder
		clientDir := expandHome(appConfig.TibiaClientDir)
		_ = os.MkdirAll(clientDir, 0755)
		zipPath := filepath.Join(clientDir, "tibia-client-15.25.zip")

		logs.WriteString(fmt.Sprintf("📥 Baixando Tibia Client 15.25 em: %s\nURL: %s\n\n", clientDir, ClientZipURL))

		dlCmd := exec.Command("curl", "-L", "-o", zipPath, ClientZipURL)
		out, err := dlCmd.CombinedOutput()
		if err != nil {
			logs.WriteString(fmt.Sprintf("❌ Falha no download: %s\n", string(out)))
			return actionDoneMsg{target: targetDownloadClient, output: logs.String(), notice: "❌ Erro no download do client!"}
		}

		logs.WriteString("📦 Extraindo arquivos do client...\n")
		unzipCmd := exec.Command("unzip", "-o", zipPath, "-d", clientDir)
		unzipOut, unzipErr := unzipCmd.CombinedOutput()
		if unzipErr != nil {
			logs.WriteString(fmt.Sprintf("❌ Falha na descompactação: %s\n", string(unzipOut)))
			return actionDoneMsg{target: targetDownloadClient, output: logs.String(), notice: "❌ Erro ao descompactar o client!"}
		}

		_ = os.Remove(zipPath)
		logs.WriteString("\n\033[1;32m✔ Tibia Client 15.25 pronto em: " + clientDir + "\033[0m\n")

		return actionDoneMsg{
			target: targetDownloadClient,
			output: logs.String(),
			notice: "✔ Download e extração do client finalizados!",
		}
	}
}