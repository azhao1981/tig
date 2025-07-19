package ui

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/azhao1981/tig/internal/config"
	"github.com/azhao1981/tig/internal/git"
)

// ViewManager manages multiple views and handles view switching
type ViewManager struct {
	screen          tcell.Screen
	config          *config.Config
	client          git.Client
	views           map[ViewType]View
	currentView     ViewType
	repoPath        string
	mutex           sync.RWMutex
	width           int
	height          int
	keyBindingMgr   *KeyBindingManager
}

// NewViewManager creates a new view manager
func NewViewManager(screen tcell.Screen, config *config.Config, client git.Client, keyBindingMgr *KeyBindingManager) *ViewManager {
	vm := &ViewManager{
		screen:        screen,
		config:        config,
		client:        client,
		views:         make(map[ViewType]View),
		currentView:   ViewTypeMain,
		keyBindingMgr: keyBindingMgr,
	}

	// Initialize views
	vm.initializeViews()
	return vm
}

// initializeViews initializes all views
func (vm *ViewManager) initializeViews() {
	// Create main view
	mainView := NewMainView(vm.config, vm.client)
	vm.views[ViewTypeMain] = mainView

	// Create diff view
	diffView := NewDiffView(vm.config, vm.client)
	vm.views[ViewTypeDiff] = diffView

	// Create status view
	statusView := NewStatusView(vm.config, vm.client)
	vm.views[ViewTypeStatus] = statusView

	// Create tree view
	treeView := NewTreeView(vm.config, vm.client)
	vm.views[ViewTypeTree] = treeView

	// Create refs view
	refsView := NewRefsView(vm.config, vm.client)
	vm.views[ViewTypeRefs] = refsView

	// Create help view
	helpView := NewHelpView(vm.config, vm.client)
	vm.views[ViewTypeHelp] = helpView

	// Set initial focus
	vm.setFocus(vm.currentView)
}

// SetSize sets the screen dimensions
func (vm *ViewManager) SetSize(width, height int) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	vm.width = width
	vm.height = height

	// Update all view sizes
	for _, view := range vm.views {
		view.SetPosition(0, 0, width, height)
	}
}

// SetRepoPath sets the repository path for all views
func (vm *ViewManager) SetRepoPath(path string) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	vm.repoPath = path
	
	// Update repository path for all views
	for _, view := range vm.views {
		switch v := view.(type) {
		case *MainView:
			v.SetRepoPath(path)
		case *DiffView:
			v.SetRepoPath(path)
		case *StatusView:
			v.SetRepoPath(path)
		case *TreeView:
			v.SetRepoPath(path)
		case *RefsView:
			v.SetRepoPath(path)
		case *HelpView:
			v.SetRepoPath(path)
		}
	}

	// Refresh all views
	vm.refreshAll()
}

// SwitchView switches to a different view
func (vm *ViewManager) SwitchView(viewType ViewType) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	if _, exists := vm.views[viewType]; !exists {
		return fmt.Errorf("view type %d not found", viewType)
	}

	// Blur current view
	if current, exists := vm.views[vm.currentView]; exists {
		current.Blur()
	}

	// Switch to new view
	vm.currentView = viewType
	vm.setFocus(vm.currentView)

	return nil
}

// setFocus sets focus to the specified view
func (vm *ViewManager) setFocus(viewType ViewType) {
	if view, exists := vm.views[viewType]; exists {
		view.Focus()
	}
}

// Render renders the current view
func (vm *ViewManager) Render() error {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	if vm.width == 0 || vm.height == 0 {
		return fmt.Errorf("screen dimensions not set")
	}

	// Clear screen
	vm.screen.Clear()

	// Render current view
	view, exists := vm.views[vm.currentView]
	if !exists {
		return fmt.Errorf("current view %d not found", vm.currentView)
	}

	return view.Render(vm.screen, 0, 0, vm.width, vm.height)
}

// HandleKey handles keyboard input
func (vm *ViewManager) HandleKey(key tcell.Key, ch rune, mod tcell.ModMask) bool {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	// Check for key bindings using the key binding manager
	if action, ok := vm.keyBindingMgr.MatchEvent(key, ch, mod); ok {
		switch action {
		case "quit":
			return false
		case "refresh":
			vm.RefreshAll()
			return true
		case "status":
			_ = vm.SwitchView(ViewTypeStatus)
			return true
		case "diff":
			_ = vm.SwitchView(ViewTypeDiff)
			return true
		case "log":
			_ = vm.SwitchView(ViewTypeMain)
			return true
		case "tree":
			_ = vm.SwitchView(ViewTypeTree)
			return true
		case "refs":
			_ = vm.SwitchView(ViewTypeRefs)
			return true
		case "help":
			_ = vm.SwitchView(ViewTypeHelp)
			return true
		case "up":
			// Let views handle navigation
			if view, exists := vm.views[vm.currentView]; exists {
				return view.HandleKey(tcell.KeyUp, 0, 0)
			}
		case "down":
			// Let views handle navigation
			if view, exists := vm.views[vm.currentView]; exists {
				return view.HandleKey(tcell.KeyDown, 0, 0)
			}
		case "page-up":
			// Let views handle navigation
			if view, exists := vm.views[vm.currentView]; exists {
				return view.HandleKey(tcell.KeyPgUp, 0, 0)
			}
		case "page-down":
			// Let views handle navigation
			if view, exists := vm.views[vm.currentView]; exists {
				return view.HandleKey(tcell.KeyPgDn, 0, 0)
			}
		case "top":
			// Let views handle navigation
			if view, exists := vm.views[vm.currentView]; exists {
				return view.HandleKey(tcell.KeyHome, 0, 0)
			}
		case "bottom":
			// Let views handle navigation
			if view, exists := vm.views[vm.currentView]; exists {
				return view.HandleKey(tcell.KeyEnd, 0, 0)
			}
		}
	}

	// Handle view-specific key bindings
	if view, exists := vm.views[vm.currentView]; exists {
		return view.HandleKey(key, ch, mod)
	}

	return false
}

// GetCurrentView returns the current view type
func (vm *ViewManager) GetCurrentView() ViewType {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()
	return vm.currentView
}

// GetView returns a specific view
func (vm *ViewManager) GetView(viewType ViewType) View {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()
	return vm.views[viewType]
}

// RefreshAll refreshes all views
func (vm *ViewManager) RefreshAll() error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	return vm.refreshAll()
}

// refreshAll refreshes all views (internal, without lock)
func (vm *ViewManager) refreshAll() error {
	var lastErr error
	
	for _, view := range vm.views {
		if err := view.Refresh(); err != nil {
			lastErr = err
		}
	}
	
	return lastErr
}

// RefreshCurrent refreshes the current view
func (vm *ViewManager) RefreshCurrent() error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	if view, exists := vm.views[vm.currentView]; exists {
		return view.Refresh()
	}
	return fmt.Errorf("current view %d not found", vm.currentView)
}

// SetDiffCommit sets the commit hash for the diff view
func (vm *ViewManager) SetDiffCommit(hash string) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	if diffView, ok := vm.views[ViewTypeDiff].(*DiffView); ok {
		diffView.SetCommitHash(hash)
		return nil
	}
	return fmt.Errorf("diff view not found")
}

// GetSelectedCommit returns the selected commit from the main view
func (vm *ViewManager) GetSelectedCommit() *git.Commit {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	if mainView, ok := vm.views[ViewTypeMain].(*MainView); ok {
		return mainView.GetSelectedCommit()
	}
	return nil
}

// UpdateDiffView updates the diff view with the selected commit
func (vm *ViewManager) UpdateDiffView() error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	commit := vm.GetSelectedCommit()
	if commit != nil {
		if diffView, ok := vm.views[ViewTypeDiff].(*DiffView); ok {
			diffView.SetCommitHash(commit.Hash)
			return nil
		}
	}
	return fmt.Errorf("no commit selected")
}

// ShowHelp shows the help view
func (vm *ViewManager) ShowHelp() error {
	// TODO: Implement help view
	return nil
}

// Exit exits the application
func (vm *ViewManager) Exit() {
	vm.screen.Fini()
}