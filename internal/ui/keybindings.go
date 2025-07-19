package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/azhao1981/tig/internal/config"
)

// KeyBinding represents a single key binding
type KeyBinding struct {
	Action string
	Key    tcell.Key
	Rune   rune
	Mods   tcell.ModMask
	Help   string
}

// KeyBindingManager manages key bindings for the application
type KeyBindingManager struct {
	bindings map[string]*KeyBinding
	config   *config.Config
}

// NewKeyBindingManager creates a new key binding manager
func NewKeyBindingManager(config *config.Config) *KeyBindingManager {
	manager := &KeyBindingManager{
		bindings: make(map[string]*KeyBinding),
		config:   config,
	}
	manager.loadDefaultBindings()
	return manager
}

// loadDefaultBindings loads the default key bindings
func (k *KeyBindingManager) loadDefaultBindings() {
	// Global bindings
	k.bindings["quit"] = &KeyBinding{
		Action: "quit",
		Key:    tcell.KeyRune,
		Rune:   'q',
		Help:   "Quit the application",
	}
	k.bindings["refresh"] = &KeyBinding{
		Action: "refresh",
		Key:    tcell.KeyRune,
		Rune:   'R',
		Help:   "Refresh all views",
	}
	k.bindings["help"] = &KeyBinding{
		Action: "help",
		Key:    tcell.KeyRune,
		Rune:   'h',
		Help:   "Show help",
	}

	// View switching
	k.bindings["status"] = &KeyBinding{
		Action: "status",
		Key:    tcell.KeyRune,
		Rune:   's',
		Help:   "Show status view",
	}
	k.bindings["diff"] = &KeyBinding{
		Action: "diff",
		Key:    tcell.KeyRune,
		Rune:   'd',
		Help:   "Show diff view",
	}
	k.bindings["log"] = &KeyBinding{
		Action: "log",
		Key:    tcell.KeyRune,
		Rune:   'l',
		Help:   "Show log view",
	}
	k.bindings["tree"] = &KeyBinding{
		Action: "tree",
		Key:    tcell.KeyRune,
		Rune:   't',
		Help:   "Show tree view",
	}
	k.bindings["refs"] = &KeyBinding{
		Action: "refs",
		Key:    tcell.KeyRune,
		Rune:   'r',
		Help:   "Show refs view",
	}

	// Navigation
	k.bindings["up"] = &KeyBinding{
		Action: "up",
		Key:    tcell.KeyUp,
		Help:   "Move selection up",
	}
	k.bindings["down"] = &KeyBinding{
		Action: "down",
		Key:    tcell.KeyDown,
		Help:   "Move selection down",
	}
	k.bindings["page-up"] = &KeyBinding{
		Action: "page-up",
		Key:    tcell.KeyPgUp,
		Help:   "Move selection up one page",
	}
	k.bindings["page-down"] = &KeyBinding{
		Action: "page-down",
		Key:    tcell.KeyPgDn,
		Help:   "Move selection down one page",
	}
	k.bindings["top"] = &KeyBinding{
		Action: "top",
		Key:    tcell.KeyRune,
		Rune:   'g',
		Help:   "Move to top",
	}
	k.bindings["bottom"] = &KeyBinding{
		Action: "bottom",
		Key:    tcell.KeyRune,
		Rune:   'G',
		Help:   "Move to bottom",
	}

	// Staging operations
	k.bindings["stage"] = &KeyBinding{
		Action: "stage",
		Key:    tcell.KeyRune,
		Rune:   'a',
		Help:   "Stage/unstage selected file",
	}
	k.bindings["unstage"] = &KeyBinding{
		Action: "unstage",
		Key:    tcell.KeyRune,
		Rune:   'u',
		Help:   "Unstage selected file",
	}
	k.bindings["stage-all"] = &KeyBinding{
		Action: "stage-all",
		Key:    tcell.KeyRune,
		Rune:   'A',
		Help:   "Stage all files",
	}
	k.bindings["unstage-all"] = &KeyBinding{
		Action: "unstage-all",
		Key:    tcell.KeyRune,
		Rune:   'U',
		Help:   "Unstage all files",
	}
	k.bindings["discard"] = &KeyBinding{
		Action: "discard",
		Key:    tcell.KeyRune,
		Rune:   'd',
		Help:   "Discard changes to selected file",
	}
	k.bindings["commit"] = &KeyBinding{
		Action: "commit",
		Key:    tcell.KeyRune,
		Rune:   'c',
		Help:   "Commit staged changes",
	}

	// Load custom bindings from config
	k.loadCustomBindings()
}

// loadCustomBindings loads custom bindings from configuration
func (k *KeyBindingManager) loadCustomBindings() {
	for action, binding := range k.config.Keymaps.Bindings {
		if binding == "" {
			continue
		}

		// Parse the binding string
		key, rune, mods := k.parseBinding(binding)
		
		if existing, ok := k.bindings[action]; ok {
			existing.Key = key
			existing.Rune = rune
			existing.Mods = mods
		}
	}
}

// parseBinding parses a binding string into key components
func (k *KeyBindingManager) parseBinding(binding string) (tcell.Key, rune, tcell.ModMask) {
	binding = strings.ToLower(binding)
	
	var mods tcell.ModMask
	
	// Handle modifier keys
	if strings.HasPrefix(binding, "ctrl-") {
		mods |= tcell.ModCtrl
		binding = strings.TrimPrefix(binding, "ctrl-")
	}
	if strings.HasPrefix(binding, "alt-") {
		mods |= tcell.ModAlt
		binding = strings.TrimPrefix(binding, "alt-")
	}
	if strings.HasPrefix(binding, "shift-") {
		mods |= tcell.ModShift
		binding = strings.TrimPrefix(binding, "shift-")
	}
	
	// Handle special keys
	switch binding {
	case "up":
		return tcell.KeyUp, 0, mods
	case "down":
		return tcell.KeyDown, 0, mods
	case "left":
		return tcell.KeyLeft, 0, mods
	case "right":
		return tcell.KeyRight, 0, mods
	case "pgup", "pageup":
		return tcell.KeyPgUp, 0, mods
	case "pgdn", "pagedown":
		return tcell.KeyPgDn, 0, mods
	case "home":
		return tcell.KeyHome, 0, mods
	case "end":
		return tcell.KeyEnd, 0, mods
	case "enter":
		return tcell.KeyEnter, 0, mods
	case "esc", "escape":
		return tcell.KeyEsc, 0, mods
	case "tab":
		return tcell.KeyTab, 0, mods
	case "backspace":
		return tcell.KeyBackspace, 0, mods
	case "del", "delete":
		return tcell.KeyDelete, 0, mods
	case "space":
		return tcell.KeyRune, ' ', mods
	default:
		if len(binding) == 1 {
			return tcell.KeyRune, rune(binding[0]), mods
		}
		return tcell.KeyRune, 0, mods
	}
}

// GetBinding returns a key binding by action name
func (k *KeyBindingManager) GetBinding(action string) (*KeyBinding, bool) {
	binding, ok := k.bindings[action]
	return binding, ok
}

// GetAllBindings returns all key bindings
func (k *KeyBindingManager) GetAllBindings() map[string]*KeyBinding {
	return k.bindings
}

// MatchEvent matches a keyboard event to a key binding
func (k *KeyBindingManager) MatchEvent(key tcell.Key, ch rune, mod tcell.ModMask) (string, bool) {
	for action, binding := range k.bindings {
		if k.matches(binding, key, ch, mod) {
			return action, true
		}
	}
	return "", false
}

// matches checks if an event matches a key binding
func (k *KeyBindingManager) matches(binding *KeyBinding, key tcell.Key, ch rune, mod tcell.ModMask) bool {
	if binding.Key != key {
		return false
	}
	
	if key == tcell.KeyRune {
		if binding.Rune != ch {
			return false
		}
	}
	
	return binding.Mods == mod
}

// GetHelpText returns help text for all bindings
func (k *KeyBindingManager) GetHelpText() []string {
	var help []string
	
	// Group bindings by category
	categories := map[string][]string{
		"Global":    {"quit", "refresh", "help"},
		"Views":     {"status", "diff", "log", "tree", "refs"},
		"Navigation":{"up", "down", "page-up", "page-down", "top", "bottom"},
		"Staging":   {"stage", "unstage", "stage-all", "unstage-all", "discard", "commit"},
	}
	
	for category, actions := range categories {
		help = append(help, fmt.Sprintf("\n%s:", category))
		for _, action := range actions {
			if binding, ok := k.bindings[action]; ok {
				keyStr := k.bindingToString(binding)
				help = append(help, fmt.Sprintf("  %-12s %s", keyStr, binding.Help))
			}
		}
	}
	
	return help
}

// bindingToString converts a key binding to a display string
func (k *KeyBindingManager) bindingToString(binding *KeyBinding) string {
	var parts []string
	
	if binding.Mods&tcell.ModCtrl != 0 {
		parts = append(parts, "Ctrl+")
	}
	if binding.Mods&tcell.ModAlt != 0 {
		parts = append(parts, "Alt+")
	}
	if binding.Mods&tcell.ModShift != 0 {
		parts = append(parts, "Shift+")
	}
	
	keyName := k.keyToString(binding.Key, binding.Rune)
	parts = append(parts, keyName)
	
	return strings.Join(parts, "")
}

// keyToString converts a key to a display string
func (k *KeyBindingManager) keyToString(key tcell.Key, ch rune) string {
	switch key {
	case tcell.KeyUp:
		return "↑"
	case tcell.KeyDown:
		return "↓"
	case tcell.KeyLeft:
		return "←"
	case tcell.KeyRight:
		return "→"
	case tcell.KeyPgUp:
		return "PgUp"
	case tcell.KeyPgDn:
		return "PgDn"
	case tcell.KeyHome:
		return "Home"
	case tcell.KeyEnd:
		return "End"
	case tcell.KeyEnter:
		return "Enter"
	case tcell.KeyEsc:
		return "Esc"
	case tcell.KeyTab:
		return "Tab"
	case tcell.KeyBackspace:
		return "Backspace"
	case tcell.KeyDelete:
		return "Delete"
	case tcell.KeyRune:
		return string(ch)
	default:
		return "Unknown"
	}
}