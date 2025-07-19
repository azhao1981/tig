package ui

import (
	"fmt"
	"testing"

	"github.com/azhao1981/tig/internal/config"
	"github.com/azhao1981/tig/internal/git"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewDiffView(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)
	assert.NotNil(t, view)
	assert.Equal(t, ViewTypeDiff, view.GetType())
	assert.NotNil(t, view.Scrollable)
	assert.Equal(t, "", view.GetCommitHash())
}

func TestDiffViewRender(t *testing.T) {
	// Create a mock screen
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)

	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)

	// Test rendering with no diff
	err = view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)

	// Test rendering with diff content
	diff := `diff --git a/file.txt b/file.txt
index 1234567..abcdefg 100644
--- a/file.txt
+++ b/file.txt
@@ -1,3 +1,4 @@
 line 1
 line 2
-line 3
+line 3 modified
+line 4 added
`
	view.diff = diff
	view.lines = []string{
		"diff --git a/file.txt b/file.txt",
		"index 1234567..abcdefg 100644",
		"--- a/file.txt",
		"+++ b/file.txt",
		"@@ -1,3 +1,4 @@",
		" line 1",
		" line 2",
		"-line 3",
		"+line 3 modified",
		"+line 4 added",
	}

	err = view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)
}

func TestDiffViewHandleKey(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)
	view.Focus()

	// Create test diff content
	lines := make([]string, 100)
	for i := 0; i < 100; i++ {
		lines[i] = fmt.Sprintf("Line %d", i+1)
	}
	view.lines = lines
	view.SetPosition(0, 0, 80, 24)

	// Test initial state
	assert.Equal(t, 0, view.GetOffset())

	// Test down navigation
	handled := view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 1, view.GetOffset())

	// Test up navigation
	handled = view.HandleKey(tcell.KeyUp, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.GetOffset())

	// Test page down
	handled = view.HandleKey(tcell.KeyPgDn, 0, 0)
	assert.True(t, handled)
	// Should move down by page size

	// Test page up
	handled = view.HandleKey(tcell.KeyPgUp, 0, 0)
	assert.True(t, handled)

	// Test home
	handled = view.HandleKey(tcell.KeyHome, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.GetOffset())

	// Test end
	handled = view.HandleKey(tcell.KeyEnd, 0, 0)
	assert.True(t, handled)
	// Should move to last page

	// Test vim-style navigation
	view.SetMaxOffset(0) // Reset to top
	handled = view.HandleKey(tcell.KeyRune, 'j', 0)
	assert.True(t, handled)
	assert.Equal(t, 1, view.GetOffset())

	handled = view.HandleKey(tcell.KeyRune, 'k', 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.GetOffset())

	handled = view.HandleKey(tcell.KeyRune, 'G', 0)
	assert.True(t, handled)
	// Should move to bottom

	handled = view.HandleKey(tcell.KeyRune, 'g', 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.GetOffset())
}

func TestDiffViewBoundaryConditions(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)
	view.Focus()
	view.SetPosition(0, 0, 80, 24)

	// Test with no lines
	view.lines = []string{}

	// Test navigation with no content
	handled := view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.GetOffset()) // Should stay at 0

	handled = view.HandleKey(tcell.KeyUp, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.GetOffset()) // Should stay at 0

	// Test with single line
	view.lines = []string{"Single line"}

	handled = view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.GetOffset()) // Should stay at 0

	handled = view.HandleKey(tcell.KeyUp, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.GetOffset()) // Should stay at 0
}

func TestDiffViewSetCommitHash(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)

	// Test setting commit hash
	view.SetCommitHash("abc123")
	assert.Equal(t, "abc123", view.GetCommitHash())

	// Test clearing commit hash
	view.SetCommitHash("")
	assert.Equal(t, "", view.GetCommitHash())
}

func TestDiffViewClear(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)

	// Set some content
	view.SetCommitHash("abc123")
	view.diff = "some diff content"
	view.lines = []string{"line1", "line2"}
	view.SetMaxOffset(10)
	view.ScrollToBottom()

	// Clear the view
	view.Clear()

	// Verify everything is cleared
	assert.Equal(t, "", view.GetCommitHash())
	assert.Equal(t, "", view.GetDiffContent())
	assert.Empty(t, view.lines)
	assert.Equal(t, 0, view.GetOffset())
	assert.Equal(t, 0, view.getMaxOffset())
}

func TestDiffViewGetDiffContent(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)

	// Test empty diff
	assert.Equal(t, "", view.GetDiffContent())

	// Test with content
	diff := "diff --git a/file.txt b/file.txt"
	view.diff = diff
	assert.Equal(t, diff, view.GetDiffContent())
}

func TestDiffViewRenderDiffLine(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)

	// Create mock screen
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)

	testCases := []struct {
		name string
		line string
	}{
		{"diff header", "diff --git a/file.txt b/file.txt"},
		{"index line", "index 1234567..abcdefg 100644"},
		{"old file", "--- a/file.txt"},
		{"new file", "+++ b/file.txt"},
		{"hunk header", "@@ -1,3 +1,4 @@"},
		{"context line", " line 1"},
		{"added line", "+new line"},
		{"removed line", "-old line"},
		{"file mode", "new file mode 100644"},
		{"deleted file", "deleted file mode 100644"},
		{"empty line", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			view.renderDiffLine(screen, 0, 0, 80, tc.line)
			// Just verify it doesn't panic
		})
	}
}

func TestDiffViewRefresh(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)

	// Test refresh with no repo path
	err := view.Refresh()
	assert.NoError(t, err) // Should not error

	// Test with valid commit hash
	view.SetCommitHash("abc123")
	err = view.Refresh()
	assert.NoError(t, err)
}

func TestDiffViewSetRepoPath(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)
	view.SetRepoPath("/path/to/repo")
	assert.Equal(t, "/path/to/repo", view.repoPath)
}

func TestDiffViewRenderWithZeroDimensions(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)

	// Test rendering with zero width
	err = view.Render(screen, 0, 0, 0, 24)
	assert.NoError(t, err)

	// Test rendering with zero height
	err = view.Render(screen, 0, 0, 80, 0)
	assert.NoError(t, err)
}

func TestDiffViewScrollableIntegration(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewDiffView(cfg, client)

	// Set up content and dimensions
	lines := make([]string, 100)
	for i := 0; i < 100; i++ {
		lines[i] = fmt.Sprintf("Line %d", i)
	}
	view.lines = lines
	view.SetPosition(0, 0, 80, 24)

	// Test scrollable integration
	assert.Equal(t, 0, view.GetOffset())
	assert.Equal(t, 22, view.getPageSize()) // 24 - 2 for borders

	// Test max offset calculation
	expectedMax := len(lines) - view.getPageSize()
	assert.Equal(t, expectedMax, view.getMaxOffset())

	// Test scrolling to bottom
	view.ScrollToBottom()
	assert.Equal(t, expectedMax, view.GetOffset())

	// Test scrolling back to top
	view.ScrollToTop()
	assert.Equal(t, 0, view.GetOffset())
}

// getMaxOffset is a helper for testing
func (v *DiffView) getMaxOffset() int {
	return v.Scrollable.maxOffset
}
