package ui

import (
	"testing"

	"github.com/azhao1981/tig/internal/config"
	"github.com/azhao1981/tig/internal/git"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewStatusView(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)
	assert.NotNil(t, view)
	assert.Equal(t, ViewTypeStatus, view.GetType())
	assert.NotNil(t, view.Scrollable)
}

func TestStatusViewRender(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)

	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// Test rendering with no status
	err = view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)

	// Test rendering with status
	view.status = &git.Status{
		Staged: []git.FileStatus{
			{Path: "staged.txt", X: "M"},
		},
		Modified: []git.FileStatus{
			{Path: "modified.txt", Y: "M"},
		},
		Untracked: []git.FileStatus{
			{Path: "untracked.txt", X: "?", Y: "?"},
		},
	}

	err = view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)
}

func TestStatusViewHandleKey(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)
	view.Focus()
	view.SetPosition(0, 0, 80, 24)

	// Manually set status for testing key handling
	view.status = &git.Status{
		Staged:   []git.FileStatus{{Path: "file1.txt"}},
		Modified: []git.FileStatus{{Path: "file2.txt"}},
	}

	assert.Equal(t, 0, view.GetOffset())

	handled := view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 1, view.GetOffset())

	handled = view.HandleKey(tcell.KeyUp, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.GetOffset())
}

func TestStatusViewRefresh(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// Test with no repo path
	err := view.Refresh()
	assert.NoError(t, err) // Should not error

	// Test with repo path
	view.SetRepoPath(".")
	err = view.Refresh()
	assert.NoError(t, err)
}

func TestStatusViewGetSelectedFile(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)

	// No file selected
	assert.Nil(t, view.GetSelectedFile())

	// Select a file
	view.status = &git.Status{
		Staged: []git.FileStatus{{Path: "file1.txt"}},
	}
	view.buildStatusLines()
	view.SetOffset(1) // Select the file under the "Staged files" header

	assert.NotNil(t, view.GetSelectedFile())
	assert.Equal(t, "file1.txt", view.GetSelectedFile().Path)
}

func TestStatusViewBoundaryConditions(t *testing.T) {
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)
	view.Focus()
	view.SetPosition(0, 0, 80, 24)

	// Test with no status, should not panic
	view.status = &git.Status{}

	// Test navigation with no content
	handled := view.HandleKey(tcell.KeyDown, 0, 0)
	assert.True(t, handled)
	assert.Equal(t, 0, view.GetOffset())
}

func TestStatusViewRenderEmpty(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)
	view.status = &git.Status{} // Empty status

	err = view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)
}

func TestStatusViewRenderWithNilStatus(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)
	cfg := &config.Config{}
	client := git.NewClient()

	view := NewStatusView(cfg, client)
	view.status = nil

	err = view.Render(screen, 0, 0, 80, 24)
	assert.NoError(t, err)
}
