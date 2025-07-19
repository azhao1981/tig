package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/jonas/tig/internal/config"
	"github.com/jonas/tig/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewViewManager(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)
	assert.NotNil(t, vm)
	assert.Equal(t, ViewTypeMain, vm.GetCurrentView())
	assert.Len(t, vm.views, 3)
	assert.Contains(t, vm.views, ViewTypeMain)
	assert.Contains(t, vm.views, ViewTypeDiff)
	assert.Contains(t, vm.views, ViewTypeStatus)
}

func TestViewManagerSetSize(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test setting size
	vm.SetSize(80, 24)
	assert.Equal(t, 80, vm.width)
	assert.Equal(t, 24, vm.height)

	// Test with zero dimensions
	vm.SetSize(0, 0)
	assert.Equal(t, 0, vm.width)
	assert.Equal(t, 0, vm.height)
}

func TestViewManagerSwitchView(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test switching to valid views
	err := vm.SwitchView(ViewTypeMain)
	assert.NoError(t, err)
	assert.Equal(t, ViewTypeMain, vm.GetCurrentView())

	err = vm.SwitchView(ViewTypeDiff)
	assert.NoError(t, err)
	assert.Equal(t, ViewTypeDiff, vm.GetCurrentView())

	err = vm.SwitchView(ViewTypeStatus)
	assert.NoError(t, err)
	assert.Equal(t, ViewTypeStatus, vm.GetCurrentView())

	// Test switching to non-existent view
	err = vm.SwitchView(ViewType(99))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestViewManagerSetRepoPath(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test setting repository path
	vm.SetRepoPath("/path/to/repo")
	assert.Equal(t, "/path/to/repo", vm.repoPath)

	// Test with empty path
	vm.SetRepoPath("")
	assert.Equal(t, "", vm.repoPath)
}

func TestViewManagerRender(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test render without size set
	err := vm.Render()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dimensions not set")

	// Test render with size set
	vm.SetSize(80, 24)
	err = vm.Render()
	assert.NoError(t, err)

	// Test render with invalid current view
	vm.currentView = ViewType(99)
	err = vm.Render()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestViewManagerHandleKey(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)
	vm.SetSize(80, 24)

	// Test global key bindings
	handled := vm.HandleKey(tcell.KeyRune, 'q', 0)
	assert.False(t, handled) // q should not be handled (triggers quit)

	handled = vm.HandleKey(tcell.KeyRune, 'l', 0)
	assert.True(t, handled)
	assert.Equal(t, ViewTypeMain, vm.GetCurrentView())

	handled = vm.HandleKey(tcell.KeyRune, 'd', 0)
	assert.True(t, handled)
	assert.Equal(t, ViewTypeDiff, vm.GetCurrentView())

	handled = vm.HandleKey(tcell.KeyRune, 's', 0)
	assert.True(t, handled)
	assert.Equal(t, ViewTypeStatus, vm.GetCurrentView())

	// Test view-specific key bindings
	vm.SwitchView(ViewTypeMain)
	handled = vm.HandleKey(tcell.KeyUp, 0, 0)
	assert.True(t, handled) // Should be handled by main view
}

func TestViewManagerGetView(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test getting existing views
	mainView := vm.GetView(ViewTypeMain)
	assert.NotNil(t, mainView)

	diffView := vm.GetView(ViewTypeDiff)
	assert.NotNil(t, diffView)

	statusView := vm.GetView(ViewTypeStatus)
	assert.NotNil(t, statusView)

	// Test getting non-existent view
	view := vm.GetView(ViewType(99))
	assert.Nil(t, view)
}

func TestViewManagerRefreshAll(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test refreshing all views
	err := vm.RefreshAll()
	assert.NoError(t, err)
}

func TestViewManagerRefreshCurrent(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test refreshing current view
	err := vm.RefreshCurrent()
	assert.NoError(t, err)

	// Test with invalid current view
	vm.currentView = ViewType(99)
	err = vm.RefreshCurrent()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestViewManagerSetDiffCommit(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test setting diff commit
	err := vm.SetDiffCommit("abc123")
	assert.NoError(t, err)

	// Test with empty hash
	err = vm.SetDiffCommit("")
	assert.NoError(t, err)

	// Test with invalid view type (this should not happen in normal usage)
	// We'll simulate by creating a mock view manager without proper initialization
	// This is more of a sanity check
}

func TestViewManagerGetSelectedCommit(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test getting selected commit
	commit := vm.GetSelectedCommit()
	assert.Nil(t, commit) // Should be nil when no repository is open
}

func TestViewManagerUpdateDiffView(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test updating diff view
	err := vm.UpdateDiffView()
	assert.NoError(t, err) // Should handle nil commit gracefully
}

func TestViewManagerShowHelp(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test showing help
	err := vm.ShowHelp()
	assert.NoError(t, err) // Placeholder for now
}

func TestViewManagerExit(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)

	// Test exit (should not panic)
	vm.Exit()
}

func TestViewManagerFocusManagement(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)
	vm.SetSize(80, 24)

	// Test focus management during view switching
	mainView := vm.GetView(ViewTypeMain).(*MainView)
	diffView := vm.GetView(ViewTypeDiff).(*DiffView)

	// Initially, main view should be focused
	assert.True(t, mainView.IsFocused())
	assert.False(t, diffView.IsFocused())

	// Switch to diff view
	vm.SwitchView(ViewTypeDiff)
	assert.False(t, mainView.IsFocused())
	assert.True(t, diffView.IsFocused())

	// Switch back to main view
	vm.SwitchView(ViewTypeMain)
	assert.True(t, mainView.IsFocused())
	assert.False(t, diffView.IsFocused())
}

func TestViewManagerConcurrentAccess(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)
	vm.SetSize(80, 24)

	// Test concurrent access (basic sanity check)
	done := make(chan bool)
	
	go func() {
		for i := 0; i < 10; i++ {
			vm.SwitchView(ViewTypeMain)
			vm.GetCurrentView()
			vm.Render()
		}
		done <- true
	}()
	
	go func() {
		for i := 0; i < 10; i++ {
			vm.SwitchView(ViewTypeStatus)
			vm.GetCurrentView()
			vm.RefreshAll()
		}
		done <- true
	}()

	// Wait for goroutines to finish
	<-done
	<-done
}

func TestViewManagerSetRepoPathPropagation(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)
	vm.SetSize(80, 24)

	// Test that repo path is propagated to all views
	repoPath := "/test/repo"
	vm.SetRepoPath(repoPath)

	// Check that views have the repo path
	mainView := vm.GetView(ViewTypeMain).(*MainView)
	assert.Equal(t, repoPath, mainView.repoPath)

	diffView := vm.GetView(ViewTypeDiff).(*DiffView)
	assert.Equal(t, repoPath, diffView.repoPath)

	statusView := vm.GetView(ViewTypeStatus).(*StatusView)
	assert.Equal(t, repoPath, statusView.repoPath)
}

func TestViewManagerIntegration(t *testing.T) {
	screen := &mockScreen{}
	cfg := &config.Config{}
	client := git.NewClient()

	vm := NewViewManager(screen, cfg, client)
	vm.SetSize(80, 24)
	vm.SetRepoPath(".")

	// Test full integration flow
	err := vm.RefreshAll()
	require.NoError(t, err)

	err = vm.Render()
	require.NoError(t, err)

	// Test view switching
	views := []ViewType{ViewTypeMain, ViewTypeDiff, ViewTypeStatus}
	for _, viewType := range views {
		err := vm.SwitchView(viewType)
		assert.NoError(t, err)
		assert.Equal(t, viewType, vm.GetCurrentView())

		err = vm.Render()
		assert.NoError(t, err)
	}

	// Test key handling across views
	vm.SwitchView(ViewTypeMain)
	handled := vm.HandleKey(tcell.KeyRune, 'd', 0)
	assert.True(t, handled)
	assert.Equal(t, ViewTypeDiff, vm.GetCurrentView())
}