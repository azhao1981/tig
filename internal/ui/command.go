package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Command represents a command that can be executed
type Command struct {
	Name        string
	Description string
	Handler     func(args []string) error
	Usage       string
}

// CommandManager manages the command system
type CommandManager struct {
	commands map[string]*Command
	buffer   string
	active   bool
	cursor   int
	history  []string
	historyIndex int
}

// NewCommandManager creates a new command manager
func NewCommandManager() *CommandManager {
	cm := &CommandManager{
		commands: make(map[string]*Command),
		buffer:   "",
		active:   false,
		cursor:   0,
		history:  make([]string, 0),
	}
	cm.registerCommands()
	return cm
}

// registerCommands registers all available commands
func (cm *CommandManager) registerCommands() {
	// Navigation commands
	cm.Register(&Command{
		Name:        "log",
		Description: "Show log/commit view",
		Handler:     cm.handleViewCommand,
		Usage:       "log",
	})

	cm.Register(&Command{
		Name:        "status",
		Description: "Show status view",
		Handler:     cm.handleViewCommand,
		Usage:       "status",
	})

	cm.Register(&Command{
		Name:        "diff",
		Description: "Show diff view",
		Handler:     cm.handleViewCommand,
		Usage:       "diff",
	})

	cm.Register(&Command{
		Name:        "tree",
		Description: "Show tree view",
		Handler:     cm.handleViewCommand,
		Usage:       "tree",
	})

	cm.Register(&Command{
		Name:        "refs",
		Description: "Show refs view",
		Handler:     cm.handleViewCommand,
		Usage:       "refs",
	})

	cm.Register(&Command{
		Name:        "help",
		Description: "Show help view",
		Handler:     cm.handleViewCommand,
		Usage:       "help",
	})

	// Git commands
	cm.Register(&Command{
		Name:        "commit",
		Description: "Commit changes",
		Handler:     cm.handleCommitCommand,
		Usage:       "commit [message]",
	})

	cm.Register(&Command{
		Name:        "add",
		Description: "Add files to staging area",
		Handler:     cm.handleAddCommand,
		Usage:       "add [files...]",
	})

	cm.Register(&Command{
		Name:        "reset",
		Description: "Reset files from staging area",
		Handler:     cm.handleResetCommand,
		Usage:       "reset [files...]",
	})

	cm.Register(&Command{
		Name:        "status",
		Description: "Show working tree status",
		Handler:     cm.handleStatusCommand,
		Usage:       "status",
	})

	// Search commands
	cm.Register(&Command{
		Name:        "search",
		Description: "Search commits by message",
		Handler:     cm.handleSearchCommand,
		Usage:       "search [pattern]",
	})

	cm.Register(&Command{
		Name:        "grep",
		Description: "Search files by content",
		Handler:     cm.handleGrepCommand,
		Usage:       "grep [pattern]",
	})

	// System commands
	cm.Register(&Command{
		Name:        "refresh",
		Description: "Refresh all views",
		Handler:     cm.handleRefreshCommand,
		Usage:       "refresh",
	})

	cm.Register(&Command{
		Name:        "quit",
		Description: "Quit the application",
		Handler:     cm.handleQuitCommand,
		Usage:       "quit",
	})

	cm.Register(&Command{
		Name:        "exit",
		Description: "Quit the application",
		Handler:     cm.handleQuitCommand,
		Usage:       "exit",
	})
}

// Register registers a new command
func (cm *CommandManager) Register(cmd *Command) {
	cm.commands[cmd.Name] = cmd
}

// Get returns a command by name
func (cm *CommandManager) Get(name string) (*Command, bool) {
	cmd, ok := cm.commands[name]
	return cmd, ok
}

// GetCommands returns all available commands
func (cm *CommandManager) GetCommands() map[string]*Command {
	return cm.commands
}

// StartCommandMode starts command mode
func (cm *CommandManager) StartCommandMode() {
	cm.active = true
	cm.buffer = ""
	cm.cursor = 0
	cm.historyIndex = -1
}

// StopCommandMode stops command mode
func (cm *CommandManager) StopCommandMode() {
	cm.active = false
	cm.buffer = ""
	cm.cursor = 0
	cm.historyIndex = -1
}

// IsActive returns whether command mode is active
func (cm *CommandManager) IsActive() bool {
	return cm.active
}

// GetBuffer returns the current command buffer
func (cm *CommandManager) GetBuffer() string {
	return cm.buffer
}

// GetCursor returns the cursor position
func (cm *CommandManager) GetCursor() int {
	return cm.cursor
}

// InsertChar inserts a character at the cursor position
func (cm *CommandManager) InsertChar(ch rune) {
	if cm.cursor == len(cm.buffer) {
		cm.buffer += string(ch)
	} else {
		cm.buffer = cm.buffer[:cm.cursor] + string(ch) + cm.buffer[cm.cursor:]
	}
	cm.cursor++
}

// Backspace removes the character before the cursor
func (cm *CommandManager) Backspace() {
	if cm.cursor > 0 {
		cm.buffer = cm.buffer[:cm.cursor-1] + cm.buffer[cm.cursor:]
		cm.cursor--
	}
}

// Delete removes the character at the cursor
func (cm *CommandManager) Delete() {
	if cm.cursor < len(cm.buffer) {
		cm.buffer = cm.buffer[:cm.cursor] + cm.buffer[cm.cursor+1:]
	}
}

// MoveCursor moves the cursor position
func (cm *CommandManager) MoveCursor(delta int) {
	cm.cursor += delta
	if cm.cursor < 0 {
		cm.cursor = 0
	}
	if cm.cursor > len(cm.buffer) {
		cm.cursor = len(cm.buffer)
	}
}

// MoveCursorToStart moves the cursor to the start
func (cm *CommandManager) MoveCursorToStart() {
	cm.cursor = 0
}

// MoveCursorToEnd moves the cursor to the end
func (cm *CommandManager) MoveCursorToEnd() {
	cm.cursor = len(cm.buffer)
}

// ClearBuffer clears the command buffer
func (cm *CommandManager) ClearBuffer() {
	cm.buffer = ""
	cm.cursor = 0
}

// Execute executes the current command
func (cm *CommandManager) Execute() error {
	if !cm.active || cm.buffer == "" {
		return nil
	}

	// Add to history
	cm.history = append(cm.history, cm.buffer)
	if len(cm.history) > 100 {
		cm.history = cm.history[1:]
	}

	// Parse command
	parts := strings.Fields(cm.buffer)
	if len(parts) == 0 {
		return nil
	}

	cmdName := parts[0]
	args := parts[1:]

	// Find and execute command
	if cmd, ok := cm.commands[cmdName]; ok {
		return cmd.Handler(args)
	}

	return fmt.Errorf("unknown command: %s", cmdName)
}

// AutoComplete returns auto-completion suggestions
func (cm *CommandManager) AutoComplete(prefix string) []string {
	var matches []string
	for name := range cm.commands {
		if strings.HasPrefix(name, prefix) {
			matches = append(matches, name)
		}
	}
	return matches
}

// HandleKey handles key events in command mode
func (cm *CommandManager) HandleKey(key tcell.Key, ch rune, mod tcell.ModMask) bool {
	if !cm.active {
		return false
	}

	switch key {
	case tcell.KeyEnter:
		// Execute command
		return true
	case tcell.KeyEsc:
		// Cancel command mode
		cm.StopCommandMode()
		return true
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		cm.Backspace()
		return true
	case tcell.KeyDelete:
		cm.Delete()
		return true
	case tcell.KeyLeft:
		cm.MoveCursor(-1)
		return true
	case tcell.KeyRight:
		cm.MoveCursor(1)
		return true
	case tcell.KeyHome:
		cm.MoveCursorToStart()
		return true
	case tcell.KeyEnd:
		cm.MoveCursorToEnd()
		return true
	case tcell.KeyUp:
		// Navigate history
		if len(cm.history) > 0 {
			if cm.historyIndex == -1 {
				cm.historyIndex = len(cm.history) - 1
			} else if cm.historyIndex > 0 {
				cm.historyIndex--
			}
			if cm.historyIndex >= 0 && cm.historyIndex < len(cm.history) {
				cm.buffer = cm.history[cm.historyIndex]
				cm.cursor = len(cm.buffer)
			}
		}
		return true
	case tcell.KeyDown:
		// Navigate history
		if cm.historyIndex < len(cm.history)-1 {
			cm.historyIndex++
			cm.buffer = cm.history[cm.historyIndex]
			cm.cursor = len(cm.buffer)
		} else {
			cm.historyIndex = -1
			cm.buffer = ""
			cm.cursor = 0
		}
		return true
	case tcell.KeyTab:
		// Auto-complete
		if cm.buffer != "" {
			parts := strings.Fields(cm.buffer)
			if len(parts) == 1 {
				suggestions := cm.AutoComplete(parts[0])
				if len(suggestions) == 1 {
					cm.buffer = suggestions[0]
					cm.cursor = len(cm.buffer)
				}
			}
		}
		return true
	default:
		if key == tcell.KeyRune && ch != 0 {
			cm.InsertChar(ch)
			return true
		}
	}
	return false
}

// Command handlers
func (cm *CommandManager) handleViewCommand(args []string) error {
	_ = args
	// This would be implemented by the view manager
	return nil
}

func (cm *CommandManager) handleCommitCommand(args []string) error {
	_ = args
	// This would be implemented by the git client
	return nil
}

func (cm *CommandManager) handleAddCommand(args []string) error {
	_ = args
	// This would be implemented by the git client
	return nil
}

func (cm *CommandManager) handleResetCommand(args []string) error {
	_ = args
	// This would be implemented by the git client
	return nil
}

func (cm *CommandManager) handleStatusCommand(args []string) error {
	_ = args
	// This would be implemented by the git client
	return nil
}

func (cm *CommandManager) handleSearchCommand(args []string) error {
	_ = args
	// This would be implemented by the git client
	return nil
}

func (cm *CommandManager) handleGrepCommand(args []string) error {
	_ = args
	// This would be implemented by the git client
	return nil
}

func (cm *CommandManager) handleRefreshCommand(args []string) error {
	_ = args
	// This would be implemented by the view manager
	return nil
}

func (cm *CommandManager) handleQuitCommand(args []string) error {
	_ = args
	// This would be implemented by the application
	return nil
}