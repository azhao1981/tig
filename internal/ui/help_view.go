package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/azhao1981/tig/internal/config"
	"github.com/azhao1981/tig/internal/git"
)

// HelpView represents the help view with key bindings and usage information
type HelpView struct {
	*BaseView
	*Scrollable
	config         *config.Config
	client         git.Client
	keyBindingMgr  *KeyBindingManager
	sections       []HelpSection
	currentSection int
	selected       int
	repoPath       string
	screen         tcell.Screen
}

// HelpSection represents a section in the help view
type HelpSection struct {
	Title   string
	Items   []HelpItem
}

// HelpItem represents a help item with key binding and description
type HelpItem struct {
	Key         string
	Description string
	Category    string
}

// NewHelpView creates a new help view
func NewHelpView(config *config.Config, client git.Client) *HelpView {
	return &HelpView{
		BaseView:       NewBaseView(ViewTypeHelp),
		Scrollable:     NewScrollable(),
		config:         config,
		client:         client,
		currentSection: 0,
		selected:       0,
	}
}

// Load loads the help view content
func (v *HelpView) Load() error {
	// Define help sections
	v.sections = []HelpSection{
		{
			Title: "Navigation",
			Items: []HelpItem{
				{Key: "j, ↓", Description: "Move selection down", Category: "navigation"},
				{Key: "k, ↑", Description: "Move selection up", Category: "navigation"},
				{Key: "g", Description: "Go to top", Category: "navigation"},
				{Key: "G", Description: "Go to bottom", Category: "navigation"},
				{Key: "PgUp", Description: "Page up", Category: "navigation"},
				{Key: "PgDn", Description: "Page down", Category: "navigation"},
			},
		},
		{
			Title: "Views",
			Items: []HelpItem{
				{Key: "l", Description: "Log view (main)", Category: "view"},
				{Key: "d", Description: "Diff view", Category: "view"},
				{Key: "s", Description: "Status view", Category: "view"},
				{Key: "t", Description: "Tree view", Category: "view"},
				{Key: "r", Description: "Refs view", Category: "view"},
				{Key: "h", Description: "Help view", Category: "view"},
			},
		},
		{
			Title: "Actions",
			Items: []HelpItem{
				{Key: "Enter", Description: "Select/open item", Category: "action"},
				{Key: "R", Description: "Refresh current view", Category: "action"},
				{Key: "q", Description: "Quit application", Category: "action"},
				{Key: "Ctrl+C", Description: "Quit application", Category: "action"},
			},
		},
		{
			Title: "Tree View",
			Items: []HelpItem{
				{Key: "Enter", Description: "Enter directory", Category: "tree"},
				{Key: "h, ←", Description: "Go up one directory", Category: "tree"},
				{Key: "l, →", Description: "Enter directory", Category: "tree"},
			},
		},
		{
			Title: "Refs View",
			Items: []HelpItem{
				{Key: "Tab", Description: "Cycle through sections", Category: "refs"},
				{Key: "1, b", Description: "Switch to branches", Category: "refs"},
				{Key: "2, t", Description: "Switch to tags", Category: "refs"},
				{Key: "3, r", Description: "Switch to remotes", Category: "refs"},
			},
		},
		{
			Title: "General",
			Items: []HelpItem{
				{Key: "Ctrl+L", Description: "Redraw screen", Category: "general"},
				{Key: "?", Description: "Show this help", Category: "general"},
			},
		},
	}

	return nil
}

// Render renders the help view
func (v *HelpView) Render(screen tcell.Screen, x, y, width, height int) error {
	if width == 0 || height == 0 {
		return fmt.Errorf("invalid screen dimensions")
	}

	screen.Clear()

	// Draw header
	header := "Go Tig Help"
	style := tcell.StyleDefault.Bold(true)
	v.drawText(screen, 0, 0, style, header)
	
	// Draw version info
	version := "Version: dev (Go implementation)"
	if len(version) < width {
		v.drawText(screen, width-len(version), 0, style, version)
	}

	// Draw separator
	for x := 0; x < width; x++ {
		screen.SetContent(x, 1, '-', nil, tcell.StyleDefault)
	}

	// Draw section tabs
	v.drawSectionTabs(screen, width)

	// Draw content
	contentStartY := 3
	maxRows := height - contentStartY - 1

	if len(v.sections) == 0 {
		msg := "No help content available"
		msgX := (width - len(msg)) / 2
		msgY := height / 2
		v.drawText(screen, msgX, msgY, tcell.StyleDefault.Dim(true), msg)
	} else {
		section := v.sections[v.currentSection]
		
		// Draw section title
		v.drawText(screen, 0, contentStartY, tcell.StyleDefault.Bold(true), section.Title)
		contentStartY += 2

		// Calculate visible items
		items := section.Items
		visibleStart := 0
		visibleEnd := len(items)
		
		if len(items) > maxRows-2 {
			visibleStart = v.GetOffset()
			visibleEnd = v.GetOffset() + maxRows - 2
			if visibleEnd > len(items) {
				visibleEnd = len(items)
				if visibleEnd-(maxRows-2) >= 0 {
					visibleStart = visibleEnd - (maxRows - 2)
				}
			}
		}

		// Draw items
		for i := visibleStart; i < visibleEnd; i++ {
			item := items[i]
			itemY := contentStartY + (i - visibleStart)
			
			if i == v.selected {
				// Highlight selected item
				for xPos := 0; xPos < width; xPos++ {
					screen.SetContent(xPos, itemY, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue))
				}
			}

			// Format key binding
			keyStyle := tcell.StyleDefault.Bold(true).Foreground(tcell.ColorYellow)
			descStyle := tcell.StyleDefault

			keyText := fmt.Sprintf("%-12s", item.Key)
			descText := item.Description

			// Ensure we don't overflow
			maxDescLen := width - 15
			if len(descText) > maxDescLen {
				descText = descText[:maxDescLen-3] + "..."
			}

			v.drawText(screen, 2, itemY, keyStyle, keyText)
			v.drawText(screen, 15, itemY, descStyle, descText)
		}
	}

	// Draw status bar
	v.drawStatusBar(screen, width, height)
	
	// Position cursor (hidden in help view)
	screen.HideCursor()

	return nil
}

// drawSectionTabs draws the section tabs
func (v *HelpView) drawSectionTabs(screen tcell.Screen, width int) {
	startX := 0
	for i, section := range v.sections {
		style := tcell.StyleDefault
		if i == v.currentSection {
			style = style.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
		}
		
		label := fmt.Sprintf(" %s ", section.Title)
		v.drawText(screen, startX, 2, style, label)
		startX += len(label) + 1
	}
}

// drawText draws text at the specified position
func (v *HelpView) drawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

// drawStatusBar draws the status bar
func (v *HelpView) drawStatusBar(screen tcell.Screen, width, height int) {
	if height < 2 {
		return
	}

	// Draw status bar background
	statusStyle := tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite)
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-1, ' ', nil, statusStyle)
	}

	// Status text
	status := "Help View - Use ↑/↓ to navigate, Tab or 1-6 to switch sections, q/Esc to close"
	if len(status) > width {
		status = "Help: ↑/↓ navigate, Tab/1-6 sections, q/Esc close"
	}
	v.drawText(screen, 0, height-1, statusStyle, status)
}

// HandleKey handles key events for the help view
func (v *HelpView) HandleKey(key tcell.Key, ch rune, mod tcell.ModMask) bool {
	switch {
	case key == tcell.KeyUp || ch == 'k':
		v.moveUp()
		return true
	case key == tcell.KeyDown || ch == 'j':
		v.moveDown()
		return true
	case key == tcell.KeyHome || ch == 'g':
		v.moveTop()
		return true
	case key == tcell.KeyEnd || ch == 'G':
		v.moveBottom()
		return true
	case key == tcell.KeyPgUp:
		v.pageUp()
		return true
	case key == tcell.KeyPgDn:
		v.pageDown()
		return true
	case key == tcell.KeyTab:
		v.nextSection()
		return true
	case ch == '1':
		v.switchSection(0)
		return true
	case ch == '2':
		v.switchSection(1)
		return true
	case ch == '3':
		v.switchSection(2)
		return true
	case ch == '4':
		v.switchSection(3)
		return true
	case ch == '5':
		v.switchSection(4)
		return true
	case ch == '6':
		v.switchSection(5)
		return true
	case ch == 'q' || key == tcell.KeyEsc:
		return false // Let view manager handle quit
	}
	return false
}

// moveUp moves the selection up
func (v *HelpView) moveUp() {
	if v.selected > 0 {
		v.selected--
		if v.selected < v.GetOffset() {
			v.SetOffset(v.selected)
		}
	}
}

// moveDown moves the selection down
func (v *HelpView) moveDown() {
	items := v.getCurrentItems()
	if v.selected < len(items)-1 {
		v.selected++
		_, _, _, height := v.GetPosition()
		maxRows := height - 6 // Account for header, tabs, etc.
		if v.selected >= v.GetOffset()+maxRows {
			v.SetOffset(v.selected - maxRows + 1)
		}
	}
}

// moveTop moves to the top of the list
func (v *HelpView) moveTop() {
	v.selected = 0
	v.SetOffset(0)
}

// moveBottom moves to the bottom of the list
func (v *HelpView) moveBottom() {
	items := v.getCurrentItems()
	v.selected = len(items) - 1
	if v.selected < 0 {
		v.selected = 0
	}
	v.adjustScroll()
}

// switchSection switches to a specific section
func (v *HelpView) switchSection(section int) {
	if section >= 0 && section < len(v.sections) {
		v.currentSection = section
		v.selected = 0
		v.SetOffset(0)
	}
}

// nextSection cycles to the next section
func (v *HelpView) nextSection() {
	v.currentSection = (v.currentSection + 1) % len(v.sections)
	v.selected = 0
	v.SetOffset(0)
}

// getCurrentItems returns the items for the current section
func (v *HelpView) getCurrentItems() []HelpItem {
	if v.currentSection >= 0 && v.currentSection < len(v.sections) {
		return v.sections[v.currentSection].Items
	}
	return []HelpItem{}
}

// adjustScroll adjusts the scroll position based on selected index
func (v *HelpView) adjustScroll() {
	items := v.getCurrentItems()
	_, _, _, height := v.GetPosition()
	maxRows := height - 6 // Account for header, tabs, etc.
	
	if v.selected < v.GetOffset() {
		v.SetOffset(v.selected)
	} else if v.selected >= v.GetOffset()+maxRows {
		v.SetOffset(v.selected - maxRows + 1)
	}
	
	if v.GetOffset() < 0 {
		v.SetOffset(0)
	}
	if v.GetOffset() > len(items)-maxRows && len(items) > maxRows {
		v.SetOffset(len(items) - maxRows)
	}
}

// pageUp moves up by one page
func (v *HelpView) pageUp() {
	_, _, _, height := v.GetPosition()
	pageSize := height - 6
	
	if v.selected >= pageSize {
		v.selected -= pageSize
	} else {
		v.selected = 0
	}
	
	if v.selected < v.GetOffset() {
		v.SetOffset(v.selected)
	}
}

// pageDown moves down by one page
func (v *HelpView) pageDown() {
	items := v.getCurrentItems()
	_, _, _, height := v.GetPosition()
	pageSize := height - 6
	
	if v.selected+pageSize < len(items) {
		v.selected += pageSize
	} else {
		v.selected = len(items) - 1
	}
	
	v.adjustScroll()
}

// Refresh refreshes the help view
func (v *HelpView) Refresh() error {
	return v.Load()
}

// SetPosition sets the view position and size
func (v *HelpView) SetPosition(x, y, width, height int) {
	v.BaseView.SetPosition(x, y, width, height)
	v.SetHeight(height - 6) // Account for header and status bar
}

// GetType returns the view type
func (v *HelpView) GetType() ViewType {
	return ViewTypeHelp
}

// SetRepoPath sets the repository path
func (v *HelpView) SetRepoPath(path string) {
	v.repoPath = path
}

// Focus sets focus to this view
func (v *HelpView) Focus() {
	v.BaseView.Focus()
}

// Blur removes focus from this view
func (v *HelpView) Blur() {
	v.BaseView.Blur()
}