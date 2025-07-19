package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonas/tig/internal/ui"
	"github.com/jonas/tig/internal/config"
	"github.com/jonas/tig/internal/git"
)

var (
	Version = "dev"
	Build   = "unknown"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get current working directory
	repoPath, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Initialize git client
	client := git.NewClient()
	if err := client.Open(repoPath); err != nil {
		// Continue without git repository - we'll show appropriate messages
	}

	terminal, err := ui.NewTerminal()
	if err != nil {
		return fmt.Errorf("failed to initialize terminal: %w", err)
	}
	defer terminal.Close()

	return terminal.Run(cfg, client, repoPath)
}