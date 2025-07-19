package ui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/azhao1981/tig/internal/config"
	"github.com/azhao1981/tig/internal/git"
)

// TreeView represents the repository tree browser view
type TreeView struct {
	*BaseView
	*Scrollable
	config      *config.Config
	client      git.Client
	files       []*git.File
	selected    int
	currentPath string
	rootPath    string
	repoPath    string
}

// NewTreeView creates a new tree view
func NewTreeView(config *config.Config, client git.Client) *TreeView {
	return &TreeView{
		BaseView:    NewBaseView(ViewTypeTree),
		Scrollable:  NewScrollable(),
		config:      config,
		client:      client,
		files:       []*git.File{},
		currentPath: "",
		rootPath:    "",
	}
}

// Load loads the tree view content
func (v *TreeView) Load() error {
	if v.client == nil {
		return fmt.Errorf("no git client available")
	}

	// Get files from git repository
	files, err := v.client.GetFiles(v.currentPath)
	if err != nil {
		return fmt.Errorf("failed to get files: %w", err)
	}

	v.files = files
	v.sortFiles()
	return nil
}

// sortFiles sorts files by type (directories first) and name
func (v *TreeView) sortFiles() {
	sort.Slice(v.files, func(i, j int) bool {
		// Directories first
		if v.files[i].IsDir != v.files[j].IsDir {
			return v.files[i].IsDir
		}
		return v.files[i].Path < v.files[j].Path
	})
}

// Render renders the tree view
func (v *TreeView) Render(screen tcell.Screen, x, y, width, height int) error {
	v.SetPosition(x, y, width, height)
	v.SetHeight(height - 3) // Account for header and status bar
	
	if width == 0 || height == 0 {
		return fmt.Errorf("invalid screen dimensions")
	}

	// Draw header
	header := "Repository Tree"
	if v.currentPath != "" {
		header = fmt.Sprintf("Tree: %s", v.currentPath)
	}
	
	// Truncate header if too long
	if len(header) > width {
		header = header[:width-3] + "..."
	}
	
	style := tcell.StyleDefault.Bold(true)
	v.drawText(screen, 0, 0, style, header)
	
	// Draw separator
	for x := 0; x < width; x++ {
		screen.SetContent(x, 1, '-', nil, tcell.StyleDefault)
	}

	// Draw file list
	startY := 2
	maxRows := height - startY - 1
	
	if len(v.files) == 0 {
		msg := "No files found"
		if v.repoPath == "" {
			msg = "No repository opened"
		}
		msgX := (width - len(msg)) / 2
		msgY := height / 2
		v.drawText(screen, msgX, msgY, tcell.StyleDefault.Dim(true), msg)
	} else {
		// Calculate visible range
		visibleStart := 0
		visibleEnd := len(v.files)
		
		if len(v.files) > maxRows {
			visibleStart = v.GetOffset()
			visibleEnd = v.GetOffset() + maxRows
			if visibleEnd > len(v.files) {
				visibleEnd = len(v.files)
				if visibleEnd-maxRows >= 0 {
					visibleStart = visibleEnd - maxRows
				}
			}
		}

		// Draw files
		for i := visibleStart; i < visibleEnd; i++ {
			file := v.files[i]
			y := startY + (i - visibleStart)
			
			if i == v.selected {
				// Highlight selected item
				for xPos := 0; xPos < width; xPos++ {
					screen.SetContent(xPos, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue))
				}
			}

			// File icon and name
			icon := "📄"
			if file.IsDir {
				icon = "📁"
			} else if file.IsBinary {
				icon = "⚙️"
			}

			path := file.Path
			if v.currentPath != "" {
				path = strings.TrimPrefix(file.Path, v.currentPath)
				path = strings.TrimPrefix(path, "/")
			}

			line := fmt.Sprintf("%s %s", icon, path)
			
			// Truncate if too long
			maxLen := width - 4
			if len(line) > maxLen {
				line = line[:maxLen-3] + "..."
			}

			style := tcell.StyleDefault
			if file.IsDir {
				style = style.Bold(true)
			}
			
			v.drawText(screen, 2, y, style, line)
		}

		// Draw scrollbar if needed
		if len(v.files) > maxRows {
			v.drawScrollbar(screen, len(v.files), maxRows, v.GetOffset())
		}
	}

	// Draw status bar
	v.drawStatusBar(screen, width, height)
	
	// Position cursor
	if v.selected >= 0 && v.selected < len(v.files) {
		cursorY := startY + (v.selected - v.GetOffset())
		if cursorY >= startY && cursorY < height-1 {
			screen.ShowCursor(0, cursorY)
		}
	}

	return nil
}

// HandleKey handles key events for the tree view
func (v *TreeView) HandleKey(key tcell.Key, ch rune, mod tcell.ModMask) bool {
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
	case key == tcell.KeyEnter || key == tcell.KeyRight:
		return v.enterDirectory()
	case key == tcell.KeyLeft || ch == 'h':
		return v.goUpDirectory()
	case ch == 'r':
		v.refresh()
		return true
	case ch == 'q':
		return false // Let view manager handle quit
	}
	return false
}

// moveUp moves the selection up
func (v *TreeView) moveUp() {
	if v.selected > 0 {
		v.selected--
		if v.selected < v.GetOffset() {
			v.SetOffset(v.selected)
		}
	}
}

// moveDown moves the selection down
func (v *TreeView) moveDown() {
	if v.selected < len(v.files)-1 {
		v.selected++
		_, _, _, height := v.GetPosition()
		maxRows := height - 3
		if v.selected >= v.GetOffset()+maxRows {
			v.SetOffset(v.selected - maxRows + 1)
		}
	}
}

// moveTop moves to the top of the list
func (v *TreeView) moveTop() {
	v.selected = 0
	v.SetOffset(0)
}

// moveBottom moves to the bottom of the list
func (v *TreeView) moveBottom() {
	v.selected = len(v.files) - 1
	if v.selected < 0 {
		v.selected = 0
	}
	v.adjustScroll()
}

// enterDirectory enters the selected directory
func (v *TreeView) enterDirectory() bool {
	if v.selected < 0 || v.selected >= len(v.files) {
		return false
	}

	file := v.files[v.selected]
	if !file.IsDir {
		return false
	}

	// Update current path
	if v.currentPath == "" {
		v.currentPath = file.Path
	} else {
		v.currentPath = filepath.Join(v.currentPath, file.Path)
	}

	// Reload the view
	v.Load()
	v.selected = 0
	v.SetOffset(0)
	return true
}

// goUpDirectory goes up one directory level
func (v *TreeView) goUpDirectory() bool {
	if v.currentPath == "" {
		return false
	}

	// Go up one level
	parent := filepath.Dir(v.currentPath)
	if parent == "." {
		parent = ""
	}

	// Find the directory we came from
	oldDir := filepath.Base(v.currentPath)

	v.currentPath = parent
	v.Load()

	// Try to select the old directory
	for i, file := range v.files {
		if file.Path == oldDir {
			v.selected = i
			break
		}
	}
	
	v.adjustScroll()
	return true
}

// refresh refreshes the tree view
func (v *TreeView) refresh() {
	v.Load()
}

// pageUp moves up by one page
func (v *TreeView) pageUp() {
	_, _, _, height := v.GetPosition()
	pageSize := height - 3
	
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
func (v *TreeView) pageDown() {
	_, _, _, height := v.GetPosition()
	pageSize := height - 3
	
	if v.selected+pageSize < len(v.files) {
		v.selected += pageSize
	} else {
		v.selected = len(v.files) - 1
	}
	
	v.adjustScroll()
}

// adjustScroll adjusts the scroll position based on selected index
func (v *TreeView) adjustScroll() {
	_, _, _, height := v.GetPosition()
	maxRows := height - 3
	
	if v.selected < v.GetOffset() {
		v.SetOffset(v.selected)
	} else if v.selected >= v.GetOffset()+maxRows {
		v.SetOffset(v.selected - maxRows + 1)
	}
}

// drawStatusBar draws the status bar
func (v *TreeView) drawStatusBar(screen tcell.Screen, width, height int) {
	if height < 2 {
		return
	}

	// Draw status bar background
	statusStyle := tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite)
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-1, ' ', nil, statusStyle)
	}

	// Status text
	status := "Tree View - Use ↑/↓ to navigate, Enter to enter dir, h/← to go up, r to refresh"
	if len(status) > width {
		status = status[:width-1]
	}
	v.drawText(screen, 0, height-1, statusStyle, status)
}

// GetType returns the view type
func (v *TreeView) GetType() ViewType {
	return ViewTypeTree
}

// drawText draws text at the specified position
func (v *TreeView) drawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	width, _ := screen.Size()
	for i, r := range text {
		if x+i >= width {
			break
		}
		screen.SetContent(x+i, y, r, nil, style)
	}
}

// drawScrollbar draws the scrollbar if needed
func (v *TreeView) drawScrollbar(screen tcell.Screen, totalItems, visibleItems, offset int) {
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

// Refresh refreshes the tree view
func (v *TreeView) Refresh() error {
	return v.Load()
}

// SetRepoPath sets the repository path
func (v *TreeView) SetRepoPath(path string) {
	v.repoPath = path
	v.rootPath = path
	v.currentPath = ""
	v.Load()
}