package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/azhao1981/tig/internal/config"
	"github.com/azhao1981/tig/internal/git"
)

// RefItem represents a reference item (branch, tag, remote)
type RefItem struct {
	Type     string
	Name     string
	Hash     string
	Current  bool
	Remote   string
	Upstream string
}

// RefsView represents the references view (branches, tags, remotes)
type RefsView struct {
	*BaseView
	*Scrollable
	config         *config.Config
	client         git.Client
	branches       []*RefItem
	tags           []*RefItem
	remotes        []*RefItem
	sections       []string
	currentSection int
	selected       int
	repoPath       string
}

// NewRefsView creates a new references view
func NewRefsView(config *config.Config, client git.Client) *RefsView {
	return &RefsView{
		BaseView:       NewBaseView(ViewTypeRefs),
		Scrollable:     NewScrollable(),
		config:         config,
		client:         client,
		branches:       []*RefItem{},
		tags:           []*RefItem{},
		remotes:        []*RefItem{},
		sections:       []string{"Branches", "Tags", "Remotes"},
		currentSection: 0,
	}
}

// Load loads the references view content
func (v *RefsView) Load() error {
	if v.client == nil {
		return fmt.Errorf("no git client available")
	}

	// Load branches
	branches, err := v.client.GetBranches()
	if err != nil {
		return fmt.Errorf("failed to get branches: %w", err)
	}

	// Load tags
	tags, err := v.client.GetTags()
	if err != nil {
		return fmt.Errorf("failed to get tags: %w", err)
	}

	// Load remotes
	remotes, err := v.client.GetRemotes()
	if err != nil {
		return fmt.Errorf("failed to get remotes: %w", err)
	}

	// Convert to ref items
	v.branches = v.convertRefs(branches, "branch")
	v.tags = v.convertRefs(tags, "tag")
	v.remotes = v.convertRemotes(remotes)

	return nil
}

// convertRefs converts git refs to ref items
func (v *RefsView) convertRefs(refs []*git.Ref, refType string) []*RefItem {
	items := []*RefItem{}
	
	// Get current HEAD
	var currentRef string
	head, err := v.client.GetHead()
	if err == nil {
		currentRef = head.Name
	}

	for _, ref := range refs {
		name := ref.Name
		// Clean up branch names
		if refType == "branch" {
			name = strings.TrimPrefix(name, "refs/heads/")
		} else if refType == "tag" {
			name = strings.TrimPrefix(name, "refs/tags/")
		}

		item := &RefItem{
			Type:    refType,
			Name:    name,
			Hash:    ref.Hash,
			Current: (refType == "branch" && ref.Name == currentRef),
		}
		items = append(items, item)
	}

	return items
}

// convertRemotes converts git remotes to ref items
func (v *RefsView) convertRemotes(remotes []*git.Remote) []*RefItem {
	items := []*RefItem{}
	for _, remote := range remotes {
		item := &RefItem{
			Type:   "remote",
			Name:   remote.Name,
			Remote: remote.Name,
		}
		items = append(items, item)
	}
	return items
}

// Render renders the refs view
func (v *RefsView) Render(screen tcell.Screen, x, y, width, height int) error {
	v.SetPosition(x, y, width, height)
	v.SetHeight(height - 5) // Account for header, tabs, and status bar
	
	if width == 0 || height == 0 {
		return fmt.Errorf("invalid screen dimensions")
	}

	screen.Clear()

	// Draw header
	header := "References"
	style := tcell.StyleDefault.Bold(true)
	v.drawText(screen, 0, 0, style, header)
	
	// Draw separator
	for xPos := 0; xPos < width; xPos++ {
		screen.SetContent(xPos, 1, '-', nil, tcell.StyleDefault)
	}

	// Draw section tabs
	v.drawSectionTabs(screen, width)

	// Draw content based on current section
	contentStartY := 3
	maxRows := height - contentStartY - 1

	var items []*RefItem
	var title string
	
	switch v.currentSection {
	case 0: // Branches
		items = v.branches
		title = fmt.Sprintf("Branches (%d)", len(v.branches))
	case 1: // Tags
		items = v.tags
		title = fmt.Sprintf("Tags (%d)", len(v.tags))
	case 2: // Remotes
		items = v.remotes
		title = fmt.Sprintf("Remotes (%d)", len(v.remotes))
	}

	// Draw section title
	v.drawText(screen, 0, contentStartY, tcell.StyleDefault.Bold(true), title)
	contentStartY++

	// Draw separator
	for xPos := 0; xPos < width; xPos++ {
		screen.SetContent(xPos, contentStartY, '-', nil, tcell.StyleDefault)
	}
	contentStartY++

	if len(items) == 0 {
		msg := "No items found"
		msgX := (width - len(msg)) / 2
		msgY := height / 2
		v.drawText(screen, msgX, msgY, tcell.StyleDefault.Dim(true), msg)
	} else {
		// Calculate visible range
		visibleStart := 0
		visibleEnd := len(items)
		
		if len(items) > maxRows {
			visibleStart = v.GetOffset()
			visibleEnd = v.GetOffset() + maxRows
			if visibleEnd > len(items) {
				visibleEnd = len(items)
				if visibleEnd-maxRows >= 0 {
					visibleStart = visibleEnd - maxRows
				}
			}
		}

		// Draw items
		for i := visibleStart; i < visibleEnd; i++ {
			item := items[i]
			y := contentStartY + (i - visibleStart)
			
			if i == v.selected {
				// Highlight selected item
				for xPos := 0; xPos < width; xPos++ {
					screen.SetContent(xPos, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue))
				}
			}

			// Item icon and name
			var icon, prefix string
			var itemStyle tcell.Style

			switch item.Type {
			case "branch":
				icon = "🌿"
				if item.Current {
					icon = "🌿*"
					itemStyle = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorGreen)
				}
			case "tag":
				icon = "🏷️"
				itemStyle = tcell.StyleDefault.Foreground(tcell.ColorYellow)
			case "remote":
				icon = "🌐"
				itemStyle = tcell.StyleDefault.Foreground(tcell.ColorBlue)
			}

			if item.Current {
				prefix = "* "
			} else {
				prefix = "  "
			}

			line := fmt.Sprintf("%s%s %s", prefix, icon, item.Name)
			
			// Truncate if too long
			maxLen := width - 4
			if len(line) > maxLen {
				line = line[:maxLen-3] + "..."
			}

			v.drawText(screen, 2, y, itemStyle, line)

			// Show hash for branches and tags
			if (item.Type == "branch" || item.Type == "tag") && item.Hash != "" {
				hash := item.Hash[:8]
				if len(hash)+len(line)+3 < width {
					hashLine := fmt.Sprintf(" %s", hash)
					v.drawText(screen, width-len(hashLine)-2, y, tcell.StyleDefault.Dim(true), hashLine)
				}
			}
		}

		// Draw scrollbar if needed
		if len(items) > maxRows {
			v.drawScrollbar(screen, len(items), maxRows, v.GetOffset())
		}
	}

	// Draw status bar
	v.drawStatusBar(screen, width, height)
	
	// Position cursor
	cursorY := contentStartY + v.selected - v.GetOffset()
	if cursorY >= contentStartY && cursorY < height-1 {
		screen.ShowCursor(0, cursorY)
	}

	return nil
}

// drawSectionTabs draws the section tabs
func (v *RefsView) drawSectionTabs(screen tcell.Screen, width int) {
	startX := 0
	for i, section := range v.sections {
		style := tcell.StyleDefault
		if i == v.currentSection {
			style = style.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
		}
		
		label := fmt.Sprintf(" %s ", section)
		v.drawText(screen, startX, 2, style, label)
		startX += len(label) + 1
	}
}

// drawText draws text at the specified position
func (v *RefsView) drawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

// drawStatusBar draws the status bar
func (v *RefsView) drawStatusBar(screen tcell.Screen, width, height int) {
	if height < 2 {
		return
	}

	// Draw status bar background
	statusStyle := tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite)
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-1, ' ', nil, statusStyle)
	}

	// Status text
	status := "Refs View - Use ↑/↓ to navigate, 1/b for branches, 2/t for tags, 3/r for remotes, Tab to cycle, R to refresh"
	if len(status) > width {
		status = status[:width-1]
	}
	v.drawText(screen, 0, height-1, statusStyle, status)
}

// drawScrollbar draws the scrollbar if needed
func (v *RefsView) drawScrollbar(screen tcell.Screen, totalItems, visibleItems, offset int) {
	if totalItems <= visibleItems {
		return
	}

	width, height := screen.Size()
	if height < 5 {
		return
	}

	scrollbarHeight := height - 4
	scrollRatio := float64(visibleItems) / float64(totalItems)
	thumbHeight := int(float64(scrollbarHeight) * scrollRatio)
	if thumbHeight < 1 {
		thumbHeight = 1
	}

	scrollPos := int(float64(offset) / float64(totalItems-visibleItems) * float64(scrollbarHeight-thumbHeight))

	// Draw scrollbar background
	scrollbarX := width - 1
	for y := 2; y < height-2; y++ {
		screen.SetContent(scrollbarX, y, '│', nil, tcell.StyleDefault.Dim(true))
	}

	// Draw scrollbar thumb
	for y := 2 + scrollPos; y < 2+scrollPos+thumbHeight && y < height-2; y++ {
		screen.SetContent(scrollbarX, y, '█', nil, tcell.StyleDefault)
	}
}

// HandleKey handles key events for the refs view
func (v *RefsView) HandleKey(key tcell.Key, ch rune, mod tcell.ModMask) bool {
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
	case ch == '1' || ch == 'b':
		v.switchSection(0)
		return true
	case ch == '2' || ch == 't':
		v.switchSection(1)
		return true
	case ch == '3' || ch == 'r':
		v.switchSection(2)
		return true
	case key == tcell.KeyTab:
		v.nextSection()
		return true
	case ch == 'R':
		v.refresh()
		return true
	case ch == 'q':
		return false // Let view manager handle quit
	}
	return false
}

// moveUp moves the selection up
func (v *RefsView) moveUp() {
	items := v.getCurrentItems()
	if v.selected > 0 {
		v.selected--
		if v.selected < v.GetOffset() {
			v.SetOffset(v.selected)
		}
		if v.selected >= len(items) {
			v.selected = len(items) - 1
		}
	}
}

// moveDown moves the selection down
func (v *RefsView) moveDown() {
	items := v.getCurrentItems()
	if v.selected < len(items)-1 {
		v.selected++
		_, _, _, height := v.GetPosition()
		maxRows := height - 5 // Account for header, tabs, etc.
		if v.selected >= v.GetOffset()+maxRows {
			v.SetOffset(v.selected - maxRows + 1)
		}
	}
}

// moveTop moves to the top of the list
func (v *RefsView) moveTop() {
	v.selected = 0
	v.SetOffset(0)
}

// moveBottom moves to the bottom of the list
func (v *RefsView) moveBottom() {
	items := v.getCurrentItems()
	v.selected = len(items) - 1
	if v.selected < 0 {
		v.selected = 0
	}
	v.adjustScroll()
}

// switchSection switches to a specific section
func (v *RefsView) switchSection(section int) {
	if section >= 0 && section < len(v.sections) {
		v.currentSection = section
		v.selected = 0
		v.SetOffset(0)
	}
}

// nextSection cycles to the next section
func (v *RefsView) nextSection() {
	v.currentSection = (v.currentSection + 1) % len(v.sections)
	v.selected = 0
	v.SetOffset(0)
}

// getCurrentItems returns the items for the current section
func (v *RefsView) getCurrentItems() []*RefItem {
	switch v.currentSection {
	case 0:
		return v.branches
	case 1:
		return v.tags
	case 2:
		return v.remotes
	}
	return []*RefItem{}
}

// refresh refreshes the refs view
func (v *RefsView) refresh() {
	v.Load()
}

// pageUp moves up by one page
func (v *RefsView) pageUp() {
	items := v.getCurrentItems()
	_, _, _, height := v.GetPosition()
	maxRows := height - 5 // Account for header, tabs, etc.
	
	if v.selected >= maxRows {
		v.selected -= maxRows
	} else {
		v.selected = 0
	}
	
	if v.selected < v.GetOffset() {
		v.SetOffset(v.selected)
	}
	
	if len(items) > 0 && v.selected >= len(items) {
		v.selected = len(items) - 1
	}
}

// pageDown moves down by one page
func (v *RefsView) pageDown() {
	items := v.getCurrentItems()
	_, _, _, height := v.GetPosition()
	maxRows := height - 5 // Account for header, tabs, etc.
	
	if v.selected+maxRows < len(items) {
		v.selected += maxRows
	} else {
		v.selected = len(items) - 1
	}
	
	v.adjustScroll()
}

// adjustScroll adjusts the scroll position based on selected index
func (v *RefsView) adjustScroll() {
	items := v.getCurrentItems()
	_, _, _, height := v.GetPosition()
	maxRows := height - 5 // Account for header, tabs, etc.
	
	if v.selected < v.GetOffset() {
		v.SetOffset(v.selected)
	} else if v.selected >= v.GetOffset()+maxRows {
		v.SetOffset(v.selected - maxRows + 1)
	}
	
	if v.GetOffset() < 0 {
		v.SetOffset(0)
	}
	if len(items) > 0 && v.selected >= len(items) {
		v.selected = len(items) - 1
	}
}

// GetType returns the view type
func (v *RefsView) GetType() ViewType {
	return ViewTypeRefs
}

// Refresh refreshes the refs view
func (v *RefsView) Refresh() error {
	return v.Load()
}

// SetRepoPath sets the repository path
func (v *RefsView) SetRepoPath(path string) {
	v.repoPath = path
	v.Load()
}