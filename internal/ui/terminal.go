package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/azhao1981/tig/internal/config"
	"github.com/azhao1981/tig/internal/git"
)

type Terminal struct {
	screen          tcell.Screen
	width           int
	height          int
	running         bool
	eventCh         chan tcell.Event
	viewManager     *ViewManager
	lastUpdate      time.Time
	theme           *Theme
	keyBindingMgr   *KeyBindingManager
	commandMgr      *CommandManager
	commandMode     bool
}

func NewTerminal() (*Terminal, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("failed to create screen: %w", err)
	}

	if err := screen.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize screen: %w", err)
	}

	width, height := screen.Size()

	terminal := &Terminal{
		screen:     screen,
		width:      width,
		height:     height,
		running:    false,
		eventCh:    make(chan tcell.Event, 10),
		lastUpdate: time.Now(),
	}

	terminal.setupScreen()
	return terminal, nil
}

func (t *Terminal) setupScreen() {
	defaultStyle := tcell.StyleDefault
	t.screen.SetStyle(defaultStyle)
	t.screen.Clear()
	t.screen.HideCursor()
}

func (t *Terminal) Close() error {
	if t.screen != nil {
		t.screen.Fini()
	}
	return nil
}

func (t *Terminal) Run(cfg *config.Config, client git.Client, repoPath string) error {
	// Initialize theme
	t.theme = NewTheme(cfg)

	// Initialize key binding manager
	t.keyBindingMgr = NewKeyBindingManager(cfg)

	// Initialize command manager
	t.commandMgr = NewCommandManager()

	// Initialize view manager
	t.viewManager = NewViewManager(t.screen, cfg, client, t.keyBindingMgr)
	t.viewManager.SetSize(t.width, t.height)
	t.viewManager.SetRepoPath(repoPath)

	// Initial refresh of all views
	t.viewManager.RefreshAll()

	t.running = true
	defer func() { t.running = false }()

	// Initial draw
	t.draw()

	// Start event loop
	go t.pollEvents()

	// Start periodic refresh
	go t.periodicRefresh()

	for t.running {
		select {
		case ev := <-t.eventCh:
			if err := t.handleEvent(ev); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Terminal) pollEvents() {
	for t.running {
		ev := t.screen.PollEvent()
		if ev != nil {
			t.eventCh <- ev
		}
	}
}

func (t *Terminal) periodicRefresh() {
	refreshTicker := time.NewTicker(5 * time.Second)
	defer refreshTicker.Stop()

	for t.running {
		select {
		case <-refreshTicker.C:
			if t.viewManager != nil {
				t.viewManager.RefreshCurrent()
				t.draw()
			}
		}
	}
}

func (t *Terminal) handleEvent(ev tcell.Event) error {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		return t.handleKeyEvent(ev)
	case *tcell.EventResize:
		return t.handleResizeEvent(ev)
	case *tcell.EventMouse:
		return t.handleMouseEvent(ev)
	}
	return nil
}

func (t *Terminal) handleKeyEvent(ev *tcell.EventKey) error {
	// Handle command mode
	if t.commandMode {
		if handled := t.commandMgr.HandleKey(ev.Key(), ev.Rune(), ev.Modifiers()); handled {
			if ev.Key() == tcell.KeyEnter {
				// Execute command
				if err := t.executeCommand(); err != nil {
					// TODO: Show error message
				}
				t.commandMode = false
				t.draw()
				return nil
			} else if ev.Key() == tcell.KeyEsc {
				t.commandMode = false
				t.draw()
				return nil
			}
			t.draw()
			return nil
		}
		return nil
	}

	// Handle command mode activation
	if ev.Rune() == ':' {
		t.commandMode = true
		t.commandMgr.StartCommandMode()
		t.draw()
		return nil
	}

	// Handle global keys
	switch ev.Key() {
	case tcell.KeyEsc, tcell.KeyCtrlC:
		t.running = false
		return nil
	case tcell.KeyCtrlL:
		t.screen.Sync()
		t.viewManager.RefreshAll()
		t.draw()
		return nil
	}

	// Handle view-specific key events
	if t.viewManager != nil {
		if handled := t.viewManager.HandleKey(ev.Key(), ev.Rune(), ev.Modifiers()); handled {
			t.draw()
			return nil
		}
	}

	return nil
}

func (t *Terminal) handleResizeEvent(ev *tcell.EventResize) error {
	t.width, t.height = ev.Size()
	if t.viewManager != nil {
		t.viewManager.SetSize(t.width, t.height)
	}
	t.draw()
	return nil
}

func (t *Terminal) handleMouseEvent(ev *tcell.EventMouse) error {
	// For now, ignore mouse events
	return nil
}

func (t *Terminal) draw() {
	if t.viewManager == nil {
		t.drawWelcome()
		return
	}

	// Render current view
	t.viewManager.Render()
	t.lastUpdate = time.Now()

	// Render command line if active
	if t.commandMode {
		t.drawCommandLine()
	}
}

func (t *Terminal) drawWelcome() {
	t.screen.Clear()

	// Draw welcome message
	welcome := "Welcome to Go Tig"
	x := (t.width - len(welcome)) / 2
	y := t.height / 2
	t.drawText(x, y, tcell.StyleDefault.Bold(true), welcome)

	// Draw instructions
	instructions := []string{
		"",
		"l - Log view",
		"d - Diff view",
		"s - Status view",
		"t - Tree view",
		"r - Refs view",
		"h - Help view",
		"q - Quit",
		"",
		"Navigation:",
		"j/k or ↑/↓ - Move selection",
		"g/G - Top/Bottom",
		"PgUp/PgDn - Page navigation",
	}

	for i, line := range instructions {
		x = (t.width - len(line)) / 2
		t.drawText(x, y+1+i, tcell.StyleDefault, line)
	}

	// Draw status
	status := fmt.Sprintf("Repository: %s", "./")
	x = (t.width - len(status)) / 2
	t.drawText(x, t.height-2, tcell.StyleDefault.Dim(true), status)

	t.screen.Show()
}

func (t *Terminal) drawText(x, y int, style tcell.Style, text string) {
	for i, r := range text {
		if x+i >= t.width {
			break
		}
		t.screen.SetContent(x+i, y, r, nil, style)
	}
}

func (t *Terminal) Size() (int, int) {
	return t.width, t.height
}

func (t *Terminal) drawCommandLine() {
	if !t.commandMode {
		return
	}

	// Draw command line at the bottom
	cmdY := t.height - 1
	
	// Clear the command line area
	for x := 0; x < t.width; x++ {
		t.screen.SetContent(x, cmdY, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))
	}

	// Draw the command prompt
	prompt := ":"
	cmdBuffer := t.commandMgr.GetBuffer()
	cursorPos := t.commandMgr.GetCursor()

	fullText := prompt + cmdBuffer
	if len(fullText) > t.width {
		// Truncate if too long
		start := len(fullText) - t.width + 1
		if start < 0 {
			start = 0
		}
		fullText = fullText[start:]
		cursorPos = len(prompt) + t.commandMgr.GetCursor() - start
	}

	// Draw the text
	for i, r := range fullText {
		t.screen.SetContent(i, cmdY, r, nil, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))
	}

	// Position the cursor
	cursorX := cursorPos
	if cursorX >= t.width {
		cursorX = t.width - 1
	}
	t.screen.ShowCursor(cursorX, cmdY)
}

func (t *Terminal) executeCommand() error {
	if !t.commandMode {
		return nil
	}

	buffer := t.commandMgr.GetBuffer()
	if buffer == "" {
		return nil
	}

	// Execute the command
	return t.commandMgr.Execute()
}