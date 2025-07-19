package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/jonas/tig/internal/config"
	"github.com/jonas/tig/internal/git"
)

// DiffView represents the diff view
type DiffView struct {
	*BaseView
	*Scrollable
	config     *config.Config
	client     git.Client
	commitHash string
	diff       string
	lines      []string
	repoPath   string
	box        *DrawBox
}

// NewDiffView creates a new diff view
func NewDiffView(config *config.Config, client git.Client) *DiffView {
	return &DiffView{
		BaseView:   NewBaseView(ViewTypeDiff),
		Scrollable: NewScrollable(),
		config:     config,
		client:     client,
		lines:      make([]string, 0),
		box:        NewDrawBox("Diff", tcell.StyleDefault.Foreground(tcell.ColorWhite)),
	}
}

// Render renders the diff view
func (v *DiffView) Render(screen tcell.Screen, x, y, width, height int) error {
	v.SetPosition(x, y, width, height)
	v.SetHeight(height - 2) // Account for borders
	
	// Draw box
	v.box.Draw(screen, x, y, width, height)
	
	// Draw content area
	contentX := x + 1
	contentY := y + 1
	contentWidth := width - 2
	contentHeight := height - 2
	
	if contentWidth <= 0 || contentHeight <= 0 {
		return nil
	}

	// Render diff content
	v.renderDiff(screen, contentX, contentY, contentWidth, contentHeight)
	
	return nil
}

// renderDiff renders the diff content
func (v *DiffView) renderDiff(screen tcell.Screen, x, y, width, height int) {
	if len(v.lines) == 0 {
		msg := "No diff to display"
		if v.commitHash == "" {
			msg = "No commit selected"
		} else if !v.client.IsRepository() {
			msg = "Not in a git repository"
		}
		
		msgX := x + (width-len(msg))/2
		msgY := y + height/2
		if msgX >= x && msgY >= y {
			for i, char := range msg {
				screen.SetContent(msgX+i, msgY, char, nil, tcell.StyleDefault)
			}
		}
		return
	}

	// Calculate visible range
	maxVisible := len(v.lines)
	if maxVisible > height {
		maxVisible = height
	}
	
	v.SetMaxOffset(len(v.lines) - height)
	if v.GetOffset() > len(v.lines)-height {
		v.SetMaxOffset(len(v.lines) - height)
	}
	
	start := v.GetOffset()
	end := start + height
	if end > len(v.lines) {
		end = len(v.lines)
	}

	// Render each line
	for i := start; i < end; i++ {
		line := v.lines[i]
		lineY := y + (i - start)
		
		if lineY >= y+height {
			break
		}
		
		// Format and render the line
		v.renderDiffLine(screen, x, lineY, width, line)
	}
}

// renderDiffLine renders a single diff line with syntax highlighting
func (v *DiffView) renderDiffLine(screen tcell.Screen, x, y, width int, line string) {
	if width <= 0 {
		return
	}

	// Determine line style based on diff content
	style := tcell.StyleDefault
	
	// Apply syntax highlighting for diff
	if strings.HasPrefix(line, "+") {
		// Added lines
		if strings.HasPrefix(line, "+++ ") {
			style = style.Foreground(tcell.ColorAqua)
		} else {
			style = style.Foreground(tcell.ColorGreen)
		}
	} else if strings.HasPrefix(line, "-") {
		// Removed lines
		if strings.HasPrefix(line, "--- ") {
			style = style.Foreground(tcell.ColorAqua)
		} else {
			style = style.Foreground(tcell.ColorRed)
		}
	} else if strings.HasPrefix(line, "@@ ") {
		// Hunk headers
		style = style.Foreground(tcell.ColorPurple).Bold(true)
	} else if strings.HasPrefix(line, "index ") {
		// Index lines
		style = style.Foreground(tcell.ColorYellow)
	} else if strings.HasPrefix(line, "diff ") {
		// Diff headers
		style = style.Foreground(tcell.ColorBlue).Bold(true)
	} else if strings.HasPrefix(line, "new file mode ") || strings.HasPrefix(line, "deleted file mode ") {
		// File mode changes
		style = style.Foreground(tcell.ColorYellow)
	}

	// Handle line truncation if needed
	if len(line) > width {
		line = line[:width-3] + "..."
	}

	// Draw the line
	for i, char := range line {
		if x+i >= x+width {
			break
		}
		screen.SetContent(x+i, y, char, nil, style)
	}

	// Fill remaining space with background
	for i := len(line); i < width; i++ {
		screen.SetContent(x+i, y, ' ', nil, tcell.StyleDefault)
	}
}

// HandleKey handles keyboard input
func (v *DiffView) HandleKey(key tcell.Key, ch rune, mod tcell.ModMask) bool {
	if !v.IsFocused() {
		return false
	}

	switch key {
	case tcell.KeyUp:
		v.ScrollUp()
		return true
	case tcell.KeyDown:
		v.ScrollDown()
		return true
	case tcell.KeyPgUp:
		v.ScrollPageUp()
		return true
	case tcell.KeyPgDn:
		v.ScrollPageDown()
		return true
	case tcell.KeyHome:
		v.ScrollToTop()
		return true
	case tcell.KeyEnd:
		v.ScrollToBottom()
		return true
	}

	switch ch {
	case 'j':
		v.ScrollDown()
		return true
	case 'k':
		v.ScrollUp()
		return true
	case 'g':
		v.ScrollToTop()
		return true
	case 'G':
		v.ScrollToBottom()
		return true
	}

	return false
}

// Refresh refreshes the diff content
func (v *DiffView) Refresh() error {
	if !v.client.IsRepository() {
		v.diff = ""
		v.lines = []string{}
		return nil
	}

	if v.commitHash == "" {
		v.diff = ""
		v.lines = []string{}
		return nil
	}

	repo, err := v.client.GetRepository()
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	// Get diff for the commit
	diff, err := repo.GetCommitDiff(v.commitHash)
	if err != nil {
		return fmt.Errorf("failed to get commit diff: %w", err)
	}

	v.diff = diff
	v.lines = strings.Split(diff, "\n")
	
	// Reset scroll position
	v.SetMaxOffset(len(v.lines) - v.getPageSize())
	if v.GetOffset() > len(v.lines)-v.getPageSize() {
		v.SetMaxOffset(len(v.lines) - v.getPageSize())
	}

	return nil
}

// SetCommitHash sets the commit hash to display diff for
func (v *DiffView) SetCommitHash(hash string) {
	v.commitHash = hash
	v.Refresh()
}

// GetCommitHash returns the current commit hash
func (v *DiffView) GetCommitHash() string {
	return v.commitHash
}

// SetRepoPath sets the repository path
func (v *DiffView) SetRepoPath(path string) {
	v.repoPath = path
}

// getPageSize returns the number of visible lines
func (v *DiffView) getPageSize() int {
	_, _, _, height := v.GetPosition()
	return height - 2 // Account for borders
}

// Clear clears the diff content
func (v *DiffView) Clear() {
	v.commitHash = ""
	v.diff = ""
	v.lines = []string{}
	v.SetMaxOffset(0)
	v.ScrollToTop()
}

// GetDiffContent returns the raw diff content
func (v *DiffView) GetDiffContent() string {
	return v.diff
}