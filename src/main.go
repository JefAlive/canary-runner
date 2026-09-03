package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func main() {
	// Força o perfil 24-bit TrueColor RGB
	lipgloss.SetColorProfile(termenv.TrueColor)

	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erro ao executar runner: %v\n", err)
		os.Exit(1)
	}
}