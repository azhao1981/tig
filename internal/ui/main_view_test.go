package ui

import (
	"testing"
	"time"

	"github.com/azhao1981/tig/internal/config"
	"github.com/azhao1981/tig/internal/git"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewMainView(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewMainView(cfg, client)
	assert.NotNil(t, view)
	assert.Equal(t, ViewTypeMain, view.GetType())
	assert.NotNil(t, view.Scrollable)
}

func TestMainViewRender(t *testing.T) {
	// Create a mock screen
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)

	cfg := &config.Config{}
	client := git.NewClient()

	view := NewMainView(cfg, client)

	// Test rendering with no commits
	err = view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)

	// Test rendering with commits
	commits := []*git.Commit{
		{
			Hash:    "abc123def456",
			Message: "Initial commit",
			Summary: "Initial commit",
			Author: git.Signature{
				Name:  "Test User",
				Email: "test@example.com",
				Time:  time.Now(),
			},
		},
	}
	view.commits = commits
	view.selected = 0

	err = view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)
}

func TestMainViewHandleKey(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewMainView(cfg, client)
	view.Focus()

	// Create test commits
	commits := []*git.Commit{
		{Hash: "1", Message: "Commit 1"},
		{Hash: "2", Message: "Commit 2"},
		{Hash: "3", Message: "Commit 3"},
		{Hash: "4", Message: "Commit 4"},
		{Hash: "5", Message: "Commit 5"},
	}
	view.commits = commits
	view.SetPosition(0, 0, 80, 24)

	// Test initial state
	assert.Equal(t, 0, view.selected)

	// Test down navigation
	handled := view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 1, view.selected)

	// Test up navigation
	handled = view.HandleKey(tcell.KeyUp, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected)

	// Test page down
	handled = view.HandleKey(tcell.KeyPgDn, 0, 0)
	assert.True(t, handled)
	// Should move down by page size

	// Test page up
	handled = view.HandleKey(tcell.KeyPgUp, 0, 0)
	assert.True(t, handled)

	// Test home
	handled = view.HandleKey(tcell.KeyHome, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected)

	// Test end
	handled = view.HandleKey(tcell.KeyEnd, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 4, view.selected)

	// Test vim-style navigation
	view.selected = 0
	handled = view.HandleKey(tcell.KeyRune, 'j', 0)
	assert.True(t, handled)
	assert.Equal(t, 1, view.selected)

	handled = view.HandleKey(tcell.KeyRune, 'k', 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected)

	handled = view.HandleKey(tcell.KeyRune, 'G', 0)
	assert.True(t, handled)
	assert.Equal(t, 4, view.selected)

	handled = view.HandleKey(tcell.KeyRune, 'g', 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected)
}

func TestMainViewBoundaryConditions(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewMainView(cfg, client)
	view.Focus()
	view.SetPosition(0, 0, 80, 24)

	// Test with no commits
	view.commits = []*git.Commit{}
	view.selected = 0

	// Test navigation with no commits
	handled := view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected) // Should stay at 0

	handled = view.HandleKey(tcell.KeyUp, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected) // Should stay at 0

	// Test with single commit
	view.commits = []*git.Commit{{Hash: "1", Message: "Commit 1"}}
	view.selected = 0

	handled = view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected) // Should stay at 0

	handled = view.HandleKey(tcell.KeyUp, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected) // Should stay at 0
}

func TestMainViewGetSelectedCommit(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewMainView(cfg, client)

	// Test with no commits
	commit := view.GetSelectedCommit()
	assert.Nil(t, commit)

	// Test with commits
	commits := []*git.Commit{
		{Hash: "1", Message: "Commit 1"},
		{Hash: "2", Message: "Commit 2"},
	}
	view.commits = commits
	view.selected = 1

	commit = view.GetSelectedCommit()
	assert.NotNil(t, commit)
	assert.Equal(t, "2", commit.Hash)

	// Test out of bounds
	view.selected = -1
	commit = view.GetSelectedCommit()
	assert.Nil(t, commit)

	view.selected = 10
	commit = view.GetSelectedCommit()
	assert.Nil(t, commit)
}

func TestMainViewRenderCommitLine(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewMainView(cfg, client)

	// Create mock screen
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)

	commit := &git.Commit{
		Hash:    "abc123def4567890",
		Message: "Add new feature\n\nThis adds a new feature to the application.",
		Summary: "Add new feature",
		Author: git.Signature{
			Name:  "John Doe",
			Email: "john@example.com",
			Time:  time.Now(),
		},
	}

	style := tcell.StyleDefault

	// Test normal rendering
	view.renderCommitLine(screen, 0, 0, 80, commit, style)

	// Test with different config options
	cfg.Views.Main.ShowGraph = true
	cfg.Views.Main.ShowRefs = true
	cfg.Views.Main.ShowID = true
	cfg.Views.Main.ShowDate = true
	cfg.Views.Main.ShowAuthor = true

	view.renderCommitLine(screen, 0, 1, 80, commit, style)

	// Test with limited width
	view.renderCommitLine(screen, 0, 2, 20, commit, style)
}

func TestMainViewRefresh(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewMainView(cfg, client)

	// Test refresh without repository
	err := view.Refresh()
	assert.NoError(t, err)
	assert.Empty(t, view.commits)

	// Test refresh with selected index adjustment
	view.commits = []*git.Commit{{Hash: "1", Message: "Test"}}
	view.selected = 5
	err = view.Refresh()
	assert.NoError(t, err)
	assert.Equal(t, 0, view.selected) // Should adjust to valid range
}

func TestMainViewConfigIntegration(t *testing.T) {
	cfg := &config.Config{}
	cfg.Views.Main.ShowGraph = true

	client := git.NewClient()
	view := NewMainView(cfg, client)

	// Test that config values are used in rendering
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)

	commits := []*git.Commit{
		{
			Hash:    "abc123def456",
			Message: "Test commit",
			Author:  git.Signature{Name: "Test", Email: "test@test.com", Time: time.Now()},
		},
	}
	view.commits = commits
	view.selected = 0

	err = view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)
}
