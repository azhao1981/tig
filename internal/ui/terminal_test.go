package ui

import (
	"testing"

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
	theme := &Theme{
		Default: 0,
		Header:  1,
		Footer:  2,
		Commit:  3,
		Author:  4,
		Date:    5,
		Branch:  6,
	}

	assert.NotNil(t, theme)
	assert.Equal(t, 0, theme.Default)
	assert.Equal(t, 1, theme.Header)
	assert.Equal(t, 2, theme.Footer)
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