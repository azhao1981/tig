package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/jonas/tig/internal/config"
	"github.com/jonas/tig/internal/git"
	"github.com/stretchr/testify/assert"
)

func TestNewStatusView(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)
	assert.NotNil(t, view)
	assert.Equal(t, ViewTypeStatus, view.GetType())
	assert.NotNil(t, view.Scrollable)
	assert.Equal(t, StatusModeFiles, view.mode)
}

func TestStatusViewRender(t *testing.T) {
	// Create a mock screen
	screen := &mockScreen{}

	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// Test rendering with no status
	err := view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)

	// Test rendering with status content
	status := &git.Status{
		Branch: "main",
		Ahead:  2,
		Behind: 1,
		Staged: []git.FileStatus{
			{Path: "file1.txt", X: "M", Y: " ", IsModified: true},
		},
		Modified: []git.FileStatus{
			{Path: "file2.txt", X: " ", Y: "M", IsModified: true},
		},
		Untracked: []git.FileStatus{
			{Path: "file3.txt", X: "?", Y: "?"},
		},
		Conflict: []git.FileStatus{
			{Path: "file4.txt", X: "U", Y: "U", IsConflict: true},
		},
	}
	view.status = status

	err = view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)
}

func TestStatusViewHandleKey(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)
	view.Focus()

	// Create test status content
	status := &git.Status{
		Branch: "main",
		Staged: []git.FileStatus{
			{Path: "file1.txt", X: "M", Y: " "},
		},
		Modified: []git.FileStatus{
			{Path: "file2.txt", X: " ", Y: "M"},
		},
	}
	view.status = status
	view.SetPosition(0, 0, 80, 24)

	// Test initial state
	assert.Equal(t, 0, view.selected)
	assert.Equal(t, StatusModeFiles, view.mode)

	// Test down navigation
	handled := view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	// Should move down

	// Test up navigation
	handled = view.HandleKey(tcell.KeyUp, 0, 0)
	assert.True(t, handled)
	// Should move up

	// Test page down
	handled = view.HandleKey(tcell.KeyPgDn, 0, 0)
	assert.True(t, handled)

	// Test page up
	handled = view.HandleKey(tcell.KeyPgUp, 0, 0)
	assert.True(t, handled)

	// Test home
	handled = view.HandleKey(tcell.KeyHome, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected)

	// Test end
	handled = view.HandleKey(tcell.KeyEnd, 0, 0)
	assert.True(t, handled)

	// Test vim-style navigation
	handled = view.HandleKey(tcell.KeyRune, 'j', 0)
	assert.True(t, handled)

	handled = view.HandleKey(tcell.KeyRune, 'k', 0)
	assert.True(t, handled)

	handled = view.HandleKey(tcell.KeyRune, 'g', 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected)

	handled = view.HandleKey(tcell.KeyRune, 'G', 0)
	assert.True(t, handled)

	// Test mode toggle
	originalMode := view.mode
	handled = view.HandleKey(tcell.KeyRune, 's', 0)
	assert.True(t, handled)
	assert.NotEqual(t, originalMode, view.mode)
	assert.Equal(t, 0, view.selected) // Should reset to top
}

func TestStatusViewBoundaryConditions(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)
	view.Focus()
	view.SetPosition(0, 0, 80, 24)

	// Test with no status
	view.status = nil

	// Test navigation with no content
	handled := view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected) // Should stay at 0

	handled = view.HandleKey(tcell.KeyUp, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected) // Should stay at 0

	// Test with empty status
	view.status = &git.Status{}

	// Test navigation with empty status
	handled = view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.selected) // Should stay at 0
}

func TestStatusViewFormatStatus(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	testCases := []struct {
		status   string
		expected string
	}{
		{"M", "modified"},
		{"A", "new file"},
		{"D", "deleted"},
		{"R", "renamed"},
		{"C", "copied"},
		{"?", "untracked"},
		{"U", "unmerged"},
		{"X", "X"},
		{"", ""},
	}

	for _, tc := range testCases {
		result := view.formatStatus(tc.status)
		assert.Equal(t, tc.expected, result)
	}
}

func TestStatusViewBuildStatusLines(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// Test with nil status
	view.status = nil
	lines := view.buildStatusLines()
	assert.Empty(t, lines)

	// Test with empty status
	view.status = &git.Status{}
	lines = view.buildStatusLines()
	assert.Contains(t, strings.Join(lines, "\n"), "nothing to commit")

	// Test with staged files
	view.status = &git.Status{
		Branch: "main",
		Staged: []git.FileStatus{
			{Path: "file1.txt", X: "M", Y: " "},
			{Path: "file2.txt", X: "A", Y: " "},
		},
	}
	lines = view.buildStatusLines()
	content := strings.Join(lines, "\n")
	assert.Contains(t, content, "On branch main")
	assert.Contains(t, content, "Changes to be committed")
	assert.Contains(t, content, "file1.txt")
	assert.Contains(t, content, "file2.txt")

	// Test with modified files
	view.status = &git.Status{
		Branch: "main",
		Modified: []git.FileStatus{
			{Path: "file1.txt", X: " ", Y: "M"},
		},
	}
	lines = view.buildStatusLines()
	content = strings.Join(lines, "\n")
	assert.Contains(t, content, "Changes not staged for commit")
	assert.Contains(t, content, "file1.txt")

	// Test with untracked files
	view.status = &git.Status{
		Branch: "main",
		Untracked: []git.FileStatus{
			{Path: "newfile.txt", X: "?", Y: "?"},
		},
	}
	lines = view.buildStatusLines()
	content = strings.Join(lines, "\n")
	assert.Contains(t, content, "Untracked files")
	assert.Contains(t, content, "newfile.txt")

	// Test with conflict files
	view.status = &git.Status{
		Branch: "main",
		Conflict: []git.FileStatus{
			{Path: "conflict.txt", X: "U", Y: "U"},
		},
	}
	lines = view.buildStatusLines()
	content = strings.Join(lines, "\n")
	assert.Contains(t, content, "Unmerged paths")
	assert.Contains(t, content, "conflict.txt")

	// Test with ahead/behind
	view.status = &git.Status{
		Branch: "main",
		Ahead:  2,
		Behind: 1,
	}
	lines = view.buildStatusLines()
	content = strings.Join(lines, "\n")
	assert.Contains(t, content, "diverged")
	assert.Contains(t, content, "2 and 1 commits")
}

func TestStatusViewRefresh(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// Test refresh without repository
	err := view.Refresh()
	assert.NoError(t, err)
	assert.Nil(t, view.status)
	assert.Equal(t, 0, view.selected)

	// Test refresh with selected index adjustment
	view.status = &git.Status{}
	view.selected = 5
	view.SetMaxOffset(10)
	view.ScrollToBottom()

	err = view.Refresh()
	assert.NoError(t, err)
	assert.Equal(t, 0, view.selected) // Should reset to top
	assert.Equal(t, 0, view.GetOffset())
}

func TestStatusViewToggleMode(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// Test mode cycling
	modes := []StatusMode{StatusModeFiles, StatusModeStaged, StatusModeModified, StatusModeUntracked, StatusModeConflict}

	for i := 0; i < len(modes)*2; i++ {
		expected := modes[i%len(modes)]
		assert.Equal(t, expected, view.mode)
		view.toggleMode()
	}
}

func TestStatusViewGetSelectedFile(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// Test with nil status
	file := view.GetSelectedFile()
	assert.Nil(t, file)

	// Test with empty status
	view.status = &git.Status{}
	file = view.GetSelectedFile()
	assert.Nil(t, file)

	// Test with files
	view.status = &git.Status{
		Staged: []git.FileStatus{
			{Path: "file1.txt", X: "M", Y: " "},
		},
		Modified: []git.FileStatus{
			{Path: "file2.txt", X: " ", Y: "M"},
		},
	}

	// Test with valid selection
	view.selected = 0
	file = view.GetSelectedFile()
	assert.NotNil(t, file)
	assert.Equal(t, "file1.txt", file.Path)

	// Test with invalid selection
	view.selected = 10
	file = view.GetSelectedFile()
	assert.Nil(t, file)

	// Test with negative selection
	view.selected = -1
	file = view.GetSelectedFile()
	assert.Nil(t, file)
}

func TestStatusViewGetStatus(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// Test empty status
	status := view.GetStatus()
	assert.Nil(t, status)

	// Test with status
	expected := &git.Status{Branch: "main"}
	view.status = expected
	status = view.GetStatus()
	assert.Equal(t, expected, status)
}

func TestStatusViewSetRepoPath(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// Test setting repository path
	view.SetRepoPath("/path/to/repo")
	assert.Equal(t, "/path/to/repo", view.repoPath)
}

func TestStatusViewRenderWithZeroDimensions(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)
	screen := &mockScreen{}

	// Test rendering with zero width
	err := view.Render(screen, 0, 0, 0, 24)
	assert.NoError(t, err)

	// Test rendering with zero height
	err = view.Render(screen, 0, 0, 80, 0)
	assert.NoError(t, err)

	// Test rendering with minimal dimensions
	err = view.Render(screen, 0, 0, 1, 1)
	assert.NoError(t, err)
}

func TestStatusViewRenderStatusLine(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)
	screen := &mockScreen{}

	// Test normal rendering
	view.renderStatusLine(screen, 0, 0, 80, "Test status line", tcell.StyleDefault)

	// Test with truncation
	view.renderStatusLine(screen, 0, 1, 10, "This is a very long status line that should be truncated", tcell.StyleDefault)

	// Test with empty line
	view.renderStatusLine(screen, 0, 2, 80, "", tcell.StyleDefault)

	// Test with zero width
	view.renderStatusLine(screen, 0, 3, 0, "Test", tcell.StyleDefault)
}

func TestStatusViewScrollableIntegration(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// Create status with many files
	status := &git.Status{
		Staged: make([]git.FileStatus, 50),
		Modified: make([]git.FileStatus, 50),
		Untracked: make([]git.FileStatus, 50),
	}
	for i := 0; i < 50; i++ {
		status.Staged[i] = git.FileStatus{Path: fmt.Sprintf("staged%d.txt", i+1)}
		status.Modified[i] = git.FileStatus{Path: fmt.Sprintf("modified%d.txt", i+1)}
		status.Untracked[i] = git.FileStatus{Path: fmt.Sprintf("untracked%d.txt", i+1)}
	}
	view.status = status
	view.SetPosition(0, 0, 80, 24)

	// Test scrollable integration
	lines := view.buildStatusLines()
	assert.Greater(t, len(lines), 100)

	// Test scrollable bounds
	view.SetMaxOffset(len(lines) - view.getPageSize())
	assert.Equal(t, len(lines)-22, view.getMaxOffset())

	// Test scrolling
	view.ScrollToBottom()
	assert.Equal(t, len(lines)-22, view.GetOffset())

	view.ScrollToTop()
	assert.Equal(t, 0, view.GetOffset())
}