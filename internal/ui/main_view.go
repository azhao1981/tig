package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/jonas/tig/internal/config"
	"github.com/jonas/tig/internal/git"
)

// MainView represents the main commit log view
type MainView struct {
	*BaseView
	*Scrollable
	config   *config.Config
	client   git.Client
	commits  []*git.Commit
	selected int
	repoPath string
	box      *DrawBox
}

// NewMainView creates a new main view
func NewMainView(config *config.Config, client git.Client) *MainView {
	return &MainView{
		BaseView:  NewBaseView(ViewTypeMain),
		Scrollable: NewScrollable(),
		config:    config,
		client:    client,
		commits:   make([]*git.Commit, 0),
		box:       NewDrawBox("Log", tcell.StyleDefault.Foreground(tcell.ColorWhite)),
	}
}

// Render renders the main view
func (v *MainView) Render(screen tcell.Screen, x, y, width, height int) error {
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

	// Render commits
	v.renderCommits(screen, contentX, contentY, contentWidth, contentHeight)
	
	return nil
}

// renderCommits renders the commit list
func (v *MainView) renderCommits(screen tcell.Screen, x, y, width, height int) {
	if len(v.commits) == 0 {
		// Show loading or no commits message
		msg := "No commits found"
		if !v.client.IsRepository() {
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

	// Ensure selected index is valid
	if v.selected < 0 {
		v.selected = 0
	}
	if v.selected >= len(v.commits) {
		v.selected = len(v.commits) - 1
	}

	// Calculate visible range
	maxVisible := len(v.commits)
	if maxVisible > height {
		maxVisible = height
	}
	
	v.SetMaxOffset(max(0, len(v.commits) - height))
	
	start := v.GetOffset()
	end := start + height
	if end > len(v.commits) {
		end = len(v.commits)
	}
	
	// Ensure start is not negative
	if start < 0 {
		start = 0
	}
	if start >= len(v.commits) {
		start = len(v.commits) - 1
		if start < 0 {
			start = 0
		}
	}

	// Render each commit
	for i := start; i < end; i++ {
		if i < 0 || i >= len(v.commits) {
			continue
		}
		
		commit := v.commits[i]
		lineY := y + (i - start)
		
		if lineY >= y+height {
			break
		}
		
		// Determine style based on selection
		style := tcell.StyleDefault
		if i == v.selected && v.IsFocused() {
			style = style.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
		} else if i == v.selected {
			style = style.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)
		}
		
		// Format commit line
		v.renderCommitLine(screen, x, lineY, width, commit, style)
	}
}

// renderCommitLine renders a single commit line
func (v *MainView) renderCommitLine(screen tcell.Screen, x, y, width int, commit *git.Commit, style tcell.Style) {
	if width <= 0 {
		return
	}
	
	// Build the commit line
	var parts []string
	
	// Show graph if enabled
	if v.config.Views.Main.ShowGraph {
		// For now, use a simple asterisk for commits
		parts = append(parts, "*")
	} else {
		parts = append(parts, " ")
	}
	
	// Show refs if enabled
	if v.config.Views.Main.ShowRefs {
		refs := v.getCommitRefs(commit.Hash)
		if len(refs) > 0 {
			parts = append(parts, strings.Join(refs, ", ")+" ")
		}
	}
	
	// Show ID if enabled
	if v.config.Views.Main.ShowID {
		id := commit.Hash
		if len(id) > 7 {
			id = id[:7]
		}
		parts = append(parts, id+" ")
	}
	
	// Show date if enabled
	if v.config.Views.Main.ShowDate {
		date := commit.Author.Time.Format("2006-01-02")
		parts = append(parts, date+" ")
	}
	
	// Show author if enabled
	if v.config.Views.Main.ShowAuthor {
		author := commit.Author.Name
		if len(author) > 20 {
			author = author[:17] + "..."
		}
		parts = append(parts, fmt.Sprintf("%-20s ", author))
	}
	
	// Show commit title
	title := commit.Summary
	if title == "" {
		title = commit.Message
		if len(title) > 50 {
			title = title[:47] + "..."
		}
	}
	parts = append(parts, title)
	
	// Combine parts
	line := strings.Join(parts, "")
	if len(line) > width {
		line = line[:width]
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
		screen.SetContent(x+i, y, ' ', nil, style)
	}
}

// getCommitRefs returns refs (branches, tags) pointing to this commit
func (v *MainView) getCommitRefs(hash string) []string {
	// This is a placeholder - in real implementation, we'd query git for refs
	// For now, return empty slice
	return []string{}
}

// HandleKey handles keyboard input
func (v *MainView) HandleKey(key tcell.Key, ch rune, mod tcell.ModMask) bool {
	if !v.IsFocused() {
		return false
	}

	switch key {
	case tcell.KeyUp:
		v.moveUp()
		return true
	case tcell.KeyDown:
		v.moveDown()
		return true
	case tcell.KeyPgUp:
		v.ScrollPageUp()
		v.selected -= v.getPageSize()
		if v.selected < 0 {
			v.selected = 0
		}
		return true
	case tcell.KeyPgDn:
		v.ScrollPageDown()
		v.selected += v.getPageSize()
		if v.selected >= len(v.commits) {
			v.selected = len(v.commits) - 1
		}
		return true
	case tcell.KeyHome:
		v.ScrollToTop()
		v.selected = 0
		return true
	case tcell.KeyEnd:
		v.ScrollToBottom()
		v.selected = len(v.commits) - 1
		return true
	}

	switch ch {
	case 'j':
		v.moveDown()
		return true
	case 'k':
		v.moveUp()
		return true
	case 'g':
		v.ScrollToTop()
		v.selected = 0
		return true
	case 'G':
		v.ScrollToBottom()
		v.selected = len(v.commits) - 1
		return true
	}

	return false
}

// moveUp moves selection up
func (v *MainView) moveUp() {
	if v.selected > 0 {
		v.selected--
		if v.selected < v.GetOffset() {
			v.ScrollUp()
		}
	}
}

// moveDown moves selection down
func (v *MainView) moveDown() {
	if v.selected < len(v.commits)-1 {
		v.selected++
		// Check if we need to scroll
		visibleEnd := v.GetOffset() + v.getPageSize()
		if v.selected >= visibleEnd {
			v.ScrollDown()
		}
	}
}

// getPageSize returns the number of visible lines
func (v *MainView) getPageSize() int {
	_, _, _, height := v.GetPosition()
	return height - 2 // Account for borders
}

// Refresh refreshes the commit list
func (v *MainView) Refresh() error {
	if !v.client.IsRepository() {
		v.commits = make([]*git.Commit, 0)
		v.selected = 0
		return nil
	}

	repo, err := v.client.GetRepository()
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	// Get commits from HEAD
	commits, err := repo.GetCommits(&git.LogOptions{
		MaxCount: 100, // Limit to 100 commits for performance
		All:      true,
	})
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	v.commits = commits
	if v.selected >= len(v.commits) {
		v.selected = len(v.commits) - 1
	}
	if v.selected < 0 {
		v.selected = 0
	}

	return nil
}

// GetSelectedCommit returns the currently selected commit
func (v *MainView) GetSelectedCommit() *git.Commit {
	if v.selected < 0 || v.selected >= len(v.commits) {
		return nil
	}
	return v.commits[v.selected]
}

// SetRepoPath sets the repository path
func (v *MainView) SetRepoPath(path string) {
	v.repoPath = path
}