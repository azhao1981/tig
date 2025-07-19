package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/azhao1981/tig/internal/config"
	"github.com/azhao1981/tig/internal/git"
)

// StatusView represents the status view showing working directory state
type StatusView struct {
	*BaseView
	*Scrollable
	config   *config.Config
	client   git.Client
	status   *git.Status
	selected int
	repoPath string
	box      *DrawBox
	mode     StatusMode
}

// StatusMode represents the current status display mode
type StatusMode int

const (
	StatusModeFiles StatusMode = iota
	StatusModeStaged
	StatusModeModified
	StatusModeUntracked
	StatusModeConflict
)

// NewStatusView creates a new status view
func NewStatusView(config *config.Config, client git.Client) *StatusView {
	return &StatusView{
		BaseView:   NewBaseView(ViewTypeStatus),
		Scrollable: NewScrollable(),
		config:     config,
		client:     client,
		box:        NewDrawBox("Status", tcell.StyleDefault.Foreground(tcell.ColorWhite)),
		mode:       StatusModeFiles,
	}
}

// Render renders the status view
func (v *StatusView) Render(screen tcell.Screen, x, y, width, height int) error {
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

	// Render status content
	v.renderStatus(screen, contentX, contentY, contentWidth, contentHeight)

	return nil
}

// renderStatus renders the status content
func (v *StatusView) renderStatus(screen tcell.Screen, x, y, width, height int) {
	if v.status == nil {
		msg := "No repository status available"
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

	// Build content lines
	lines := v.buildStatusLines()

	// Calculate visible range
	maxVisible := len(lines)
	if maxVisible > height {
		maxVisible = height
	}

	v.SetMaxOffset(len(lines) - height)
	if v.GetOffset() > len(lines)-height {
		v.SetMaxOffset(len(lines) - height)
	}

	start := v.GetOffset()
	end := start + height
	if end > len(lines) {
		end = len(lines)
	}

	// Render each line
	for i := start; i < end; i++ {
		line := lines[i]
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

		// Format and render the line
		v.renderStatusLine(screen, x, lineY, width, line, style)
	}
}

// buildStatusLines builds the status content lines
func (v *StatusView) buildStatusLines() []string {
	lines := make([]string, 0)

	// Add branch information
	if v.status.Branch != "" {
		lines = append(lines, fmt.Sprintf("On branch %s", v.status.Branch))
		
		// Add ahead/behind information
		if v.status.Ahead > 0 || v.status.Behind > 0 {
			aheadBehind := fmt.Sprintf("Your branch is ahead of 'origin/%s' by %d commit(s)", v.status.Branch, v.status.Ahead)
			if v.status.Behind > 0 {
				aheadBehind = fmt.Sprintf("Your branch is behind 'origin/%s' by %d commit(s)", v.status.Branch, v.status.Behind)
			}
			if v.status.Ahead > 0 && v.status.Behind > 0 {
				aheadBehind = fmt.Sprintf("Your branch and 'origin/%s' have diverged by %d and %d commits respectively", v.status.Branch, v.status.Ahead, v.status.Behind)
			}
			lines = append(lines, aheadBehind)
		}
		lines = append(lines, "")
	}

	// Add staged files
	if len(v.status.Staged) > 0 {
		lines = append(lines, "Changes to be committed:")
		lines = append(lines, `  (use "git reset HEAD <file>..." to unstage)`)
		for _, file := range v.status.Staged {
			lines = append(lines, fmt.Sprintf("\t%s: %s", v.formatStatus(file.X), file.Path))
		}
		lines = append(lines, "")
	}

	// Add modified files
	if len(v.status.Modified) > 0 {
		lines = append(lines, "Changes not staged for commit:")
		lines = append(lines, "  (use \"git add <file>...\" to update what will be committed)")
		lines = append(lines, "  (use \"git checkout -- <file>...\" to discard changes in working directory)")
		for _, file := range v.status.Modified {
			lines = append(lines, fmt.Sprintf("\t%s: %s", v.formatStatus(file.Y), file.Path))
		}
		lines = append(lines, "")
	}

	// Add untracked files
	if len(v.status.Untracked) > 0 {
		lines = append(lines, "Untracked files:")
		lines = append(lines, `  (use "git add <file>..." to include in what will be committed)`)
		for _, file := range v.status.Untracked {
			lines = append(lines, fmt.Sprintf("\t%s", file.Path))
		}
		lines = append(lines, "")
	}

	// Add conflict files
	if len(v.status.Conflict) > 0 {
		lines = append(lines, "Unmerged paths:")
		lines = append(lines, `  (use "git add <file>..." to mark resolution)`)
		for _, file := range v.status.Conflict {
			lines = append(lines, fmt.Sprintf("\tboth modified: %s", file.Path))
		}
		lines = append(lines, "")
	}

	// Add summary
	if len(v.status.Staged) == 0 && len(v.status.Modified) == 0 && len(v.status.Untracked) == 0 && len(v.status.Conflict) == 0 {
		lines = append(lines, "nothing to commit, working tree clean")
	} else {
		var staged, modified, untracked, conflict int
		staged = len(v.status.Staged)
		modified = len(v.status.Modified)
		untracked = len(v.status.Untracked)
		conflict = len(v.status.Conflict)

		summary := "no changes added to commit"
		if staged > 0 || modified > 0 || untracked > 0 || conflict > 0 {
			parts := make([]string, 0)
			if staged > 0 {
				parts = append(parts, fmt.Sprintf("%d staged", staged))
			}
			if modified > 0 {
				parts = append(parts, fmt.Sprintf("%d modified", modified))
			}
			if untracked > 0 {
				parts = append(parts, fmt.Sprintf("%d untracked", untracked))
			}
			if conflict > 0 {
				parts = append(parts, fmt.Sprintf("%d conflicted", conflict))
			}
			summary = strings.Join(parts, ", ") + " changes"
		}
		lines = append(lines, summary)
	}
	
	// Add key bindings help
	lines = append(lines, "")
	lines = append(lines, "Key bindings:")
	lines = append(lines, "  a - stage/unstage selected file")
	lines = append(lines, "  u - unstage selected file")
	lines = append(lines, "  d - discard changes to selected file")
	lines = append(lines, "  A - stage all files")
	lines = append(lines, "  U - unstage all files")
	lines = append(lines, "  c - commit staged changes")
	lines = append(lines, "  s - switch display mode")
	lines = append(lines, "  q - quit")

	return lines
}

// formatStatus formats the git status character
func (v *StatusView) formatStatus(status string) string {
	switch status {
	case "M":
		return "modified"
	case "A":
		return "new file"
	case "D":
		return "deleted"
	case "R":
		return "renamed"
	case "C":
		return "copied"
	case "?":
		return "untracked"
	case "U":
		return "unmerged"
	default:
		return status
	}
}

// renderStatusLine renders a single status line
func (v *StatusView) renderStatusLine(screen tcell.Screen, x, y, width int, line string, style tcell.Style) {
	if width <= 0 {
		return
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
func (v *StatusView) HandleKey(key tcell.Key, ch rune, mod tcell.ModMask) bool {
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
		v.moveDown()
		return true
	case 'k':
		v.moveUp()
		return true
	case 'g':
		v.ScrollToTop()
		return true
	case 'G':
		v.ScrollToBottom()
		return true
	case 's':
		// Toggle between status modes
		v.toggleMode()
		return true
	case 'a':
		// Stage/unstage selected file
		v.stageSelectedFile()
		return true
	case 'u':
		// Unstage selected file
		if v.canUnstageSelectedFile() {
			v.unstageSelectedFile()
		}
		return true
	case 'd':
		// Discard changes to selected file
		v.discardSelectedFile()
		return true
	case 'A':
		// Stage all files
		v.stageAllFiles()
		return true
	case 'U':
		// Unstage all files
		v.unstageAllFiles()
		return true
	case 'c':
		// Commit staged changes
		v.commit()
		return true
	}

	return false
}

// moveUp moves selection up
func (v *StatusView) moveUp() {
	if v.selected > 0 {
		v.selected--
		if v.selected < v.GetOffset() {
			v.ScrollUp()
		}
	}
}

// moveDown moves selection down
func (v *StatusView) moveDown() {
	lines := v.buildStatusLines()
	if v.selected < len(lines)-1 {
		v.selected++
		visibleEnd := v.GetOffset() + v.getPageSize()
		if v.selected >= visibleEnd {
			v.ScrollDown()
		}
	}
}

// toggleMode toggles between different status display modes
func (v *StatusView) toggleMode() {
	v.mode = (v.mode + 1) % 5 // Cycle through 5 modes
	v.selected = 0
	v.ScrollToTop()
}

// Refresh refreshes the status content
func (v *StatusView) Refresh() error {
	if !v.client.IsRepository() {
		v.status = nil
		v.selected = 0
		return nil
	}

	repo, err := v.client.GetRepository()
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	// Get repository status
	status, err := repo.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	v.status = status
	v.selected = 0
	v.ScrollToTop()

	return nil
}

// SetRepoPath sets the repository path
func (v *StatusView) SetRepoPath(path string) {
	v.repoPath = path
}

// getPageSize returns the number of visible lines
func (v *StatusView) getPageSize() int {
	_, _, _, height := v.GetPosition()
	return height - 2 // Account for borders
}

// GetSelectedFile returns the currently selected file
func (v *StatusView) GetSelectedFile() *git.FileStatus {
	if v.status == nil {
		return nil
	}

	files := v.getAllFiles()
	if v.selected < 0 || v.selected >= len(files) {
		return nil
	}

	return &files[v.selected]
}

// GetStatus returns the current git status
func (v *StatusView) GetStatus() *git.Status {
	return v.status
}

// getAllFiles returns all files in the current status
func (v *StatusView) getAllFiles() []git.FileStatus {
	if v.status == nil {
		return []git.FileStatus{}
	}

	var files []git.FileStatus
	
	switch v.mode {
	case StatusModeStaged:
		files = v.status.Staged
	case StatusModeModified:
		files = v.status.Modified
	case StatusModeUntracked:
		files = v.status.Untracked
	case StatusModeConflict:
		files = v.status.Conflict
	default:
		// All files
		files = append(files, v.status.Staged...)
		files = append(files, v.status.Modified...)
		files = append(files, v.status.Untracked...)
		files = append(files, v.status.Conflict...)
	}
	
	return files
}

// stageSelectedFile stages the currently selected file
func (v *StatusView) stageSelectedFile() error {
	file := v.GetSelectedFile()
	if file == nil {
		return nil
	}

	if file.IsUntracked || file.IsModified {
		err := v.client.StageFile(file.Path)
		if err != nil {
			return fmt.Errorf("failed to stage %s: %w", file.Path, err)
		}
		
		// Refresh the status view
		return v.Refresh()
	}
	
	return nil
}

// unstageSelectedFile unstages the currently selected file
func (v *StatusView) unstageSelectedFile() error {
	file := v.GetSelectedFile()
	if file == nil {
		return nil
	}

	if v.canUnstageSelectedFile() {
		err := v.client.UnstageFile(file.Path)
		if err != nil {
			return fmt.Errorf("failed to unstage %s: %w", file.Path, err)
		}
		
		// Refresh the status view
		return v.Refresh()
	}
	
	return nil
}

// canUnstageSelectedFile checks if the selected file can be unstaged
func (v *StatusView) canUnstageSelectedFile() bool {
	file := v.GetSelectedFile()
	if file == nil {
		return false
	}
	
	// Check if file is staged
	for _, staged := range v.status.Staged {
		if staged.Path == file.Path {
			return true
		}
	}
	return false
}

// discardSelectedFile discards changes to the selected file
func (v *StatusView) discardSelectedFile() error {
	file := v.GetSelectedFile()
	if file == nil {
		return nil
	}

	if file.IsModified {
		err := v.client.DiscardChanges(file.Path)
		if err != nil {
			return fmt.Errorf("failed to discard changes to %s: %w", file.Path, err)
		}
		
		// Refresh the status view
		return v.Refresh()
	}
	
	return nil
}

// stageAllFiles stages all modified and untracked files
func (v *StatusView) stageAllFiles() error {
	err := v.client.StageAll()
	if err != nil {
		return fmt.Errorf("failed to stage all files: %w", err)
	}
	
	return v.Refresh()
}

// unstageAllFiles unstages all files
func (v *StatusView) unstageAllFiles() error {
	err := v.client.UnstageAll()
	if err != nil {
		return fmt.Errorf("failed to unstage all files: %w", err)
	}
	
	return v.Refresh()
}

// commit opens a commit interface
func (v *StatusView) commit() {
	// This would be implemented with a commit dialog
	// For now, we'll just log the action
	fmt.Printf("Commit functionality would be implemented here\n")
}