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