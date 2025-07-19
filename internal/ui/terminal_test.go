package ui

import (
	"testing"

	"github.com/azhao1981/tig/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTerminal(t *testing.T) {
	// Skip test in CI environments where terminal might not be available
	if testing.Short() {
		t.Skip("skipping terminal test in short mode")
	}

	terminal, err := NewTerminal()
	if err != nil {
		t.Skipf("terminal not available: %v", err)
	}

	require.NotNil(t, terminal)
	assert.Greater(t, terminal.width, 0)
	assert.Greater(t, terminal.height, 0)

	// Test size getter
	w, h := terminal.Size()
	assert.Equal(t, terminal.width, w)
	assert.Equal(t, terminal.height, h)

	// Test close
	err = terminal.Close()
	assert.NoError(t, err)
}

func TestTheme(t *testing.T) {
	// Create a mock config for testing
	cfg := &config.Config{}
	cfg.Colors.Colors = map[string]string{
		"default": "white",
		"header":  "blue",
		"commit":  "yellow",
	}

	theme := NewTheme(cfg)
	assert.NotNil(t, theme)

	// Check that the colors from the config are loaded correctly
	assert.Equal(t, tcell.ColorWhite, theme.GetColor("default"))
	assert.Equal(t, tcell.ColorBlue, theme.GetColor("header"))
	assert.Equal(t, tcell.ColorYellow, theme.GetColor("commit"))

	// Check a default color that wasn't in the config
	assert.Equal(t, tcell.ColorGreen, theme.GetColor("status"))
}

func TestTerminalMethods(t *testing.T) {
	// Test that methods exist and don't panic
	terminal := &Terminal{
		width:  80,
		height: 24,
		theme:  &Theme{},
	}

	w, h := terminal.Size()
	assert.Equal(t, 80, w)
	assert.Equal(t, 24, h)
}
