package ui

import (
	"testing"

	"github.com/azhao1981/tig/internal/config"
	"github.com/azhao1981/tig/internal/git"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewViewManager(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)
	cfg := &config.Config{}
	client := git.NewClient()
	keyBindingMgr := NewKeyBindingManager(cfg)

	vm := NewViewManager(screen, cfg, client, keyBindingMgr)
	assert.NotNil(t, vm)
	assert.Equal(t, ViewTypeMain, vm.GetCurrentView())
	assert.NotNil(t, vm.GetView(ViewTypeMain))
}

func TestViewManagerSwitchView(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)
	cfg := &config.Config{}
	client := git.NewClient()
	keyBindingMgr := NewKeyBindingManager(cfg)

	vm := NewViewManager(screen, cfg, client, keyBindingMgr)
	err = vm.SwitchView(ViewTypeDiff)
	assert.NoError(t, err)
	assert.Equal(t, ViewTypeDiff, vm.GetCurrentView())

	err = vm.SwitchView(ViewType(99)) // Invalid view
	assert.Error(t, err)
}

func TestViewManagerRender(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)
	cfg := &config.Config{}
	client := git.NewClient()
	keyBindingMgr := NewKeyBindingManager(cfg)

	vm := NewViewManager(screen, cfg, client, keyBindingMgr)
	vm.SetSize(80, 24)

	err = vm.Render()
	assert.NoError(t, err)
}

func TestViewManagerHandleKey(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)
	cfg := &config.Config{}
	client := git.NewClient()
	keyBindingMgr := NewKeyBindingManager(cfg)

	vm := NewViewManager(screen, cfg, client, keyBindingMgr)
	vm.SetSize(80, 24)

	// Test quit
	handled := vm.HandleKey(tcell.KeyRune, 'q', 0)
	assert.False(t, handled) // Should signal to quit

	// Test view switching
	handled = vm.HandleKey(tcell.KeyRune, 's', 0)
	assert.True(t, handled)
	assert.Equal(t, ViewTypeStatus, vm.GetCurrentView())

	// Test navigation (should be handled by view)
	handled = vm.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
}

func TestViewManagerRefreshAll(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)
	cfg := &config.Config{}
	client := git.NewClient()
	keyBindingMgr := NewKeyBindingManager(cfg)

	vm := NewViewManager(screen, cfg, client, keyBindingMgr)
	vm.SetRepoPath(".") // Set a valid repo path for refreshing

	err = vm.refreshAll()
	assert.NoError(t, err)
}

func TestViewManagerGetSelectedCommit(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)
	cfg := &config.Config{}
	client := git.NewClient()
	keyBindingMgr := NewKeyBindingManager(cfg)

	vm := NewViewManager(screen, cfg, client, keyBindingMgr)

	// No commit selected initially
	commit := vm.GetSelectedCommit()
	assert.Nil(t, commit)

	// Switch to main view and select a commit
	mainView := vm.GetView(ViewTypeMain).(*MainView)
	mainView.commits = []*git.Commit{
		{Hash: "1", Message: "Commit 1"},
	}
	mainView.selected = 0

	commit = vm.GetSelectedCommit()
	assert.NotNil(t, commit)
	assert.Equal(t, "1", commit.Hash)
}

func TestViewManagerUpdateDiffView(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)
	cfg := &config.Config{}
	client := git.NewClient()
	keyBindingMgr := NewKeyBindingManager(cfg)

	vm := NewViewManager(screen, cfg, client, keyBindingMgr)
	vm.SetRepoPath(".")

	// Select a commit in main view
	mainView := vm.GetView(ViewTypeMain).(*MainView)
	mainView.commits = []*git.Commit{
		{Hash: "1", Message: "Commit 1"},
	}
	mainView.selected = 0

	err = vm.UpdateDiffView()
	assert.NoError(t, err)

	diffView := vm.GetView(ViewTypeDiff).(*DiffView)
	assert.Equal(t, "1", diffView.GetCommitHash())
}
