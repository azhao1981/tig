package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Test loading with no config file (should use defaults)
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Test default values
	assert.Equal(t, 8, cfg.UI.TabSize)
	assert.Equal(t, "topo", cfg.UI.CommitOrder)
	assert.Equal(t, false, cfg.UI.IgnoreCase)
	assert.Equal(t, true, cfg.UI.ShowDate)
	assert.Equal(t, true, cfg.UI.ShowAuthor)
	assert.Equal(t, true, cfg.UI.ShowLineNumbers)

	assert.Equal(t, 20, cfg.Git.AuthorWidth)
	assert.Equal(t, "%Y-%m-%d", cfg.Git.DateFormat)
	assert.Equal(t, true, cfg.Git.ShowNotes)
	assert.Equal(t, true, cfg.Git.ShowDiffStat)
	assert.Equal(t, true, cfg.Git.ShowBranches)
	assert.Equal(t, true, cfg.Git.ShowRemotes)
	assert.Equal(t, true, cfg.Git.ShowTags)

	assert.Equal(t, true, cfg.Views.Main.ShowGraph)
	assert.Equal(t, true, cfg.Views.Main.ShowRefs)
	assert.Equal(t, false, cfg.Views.Main.ShowID)
	assert.Equal(t, true, cfg.Views.Main.ShowDate)
	assert.Equal(t, true, cfg.Views.Main.ShowAuthor)
	assert.Equal(t, true, cfg.Views.Main.ShowCommitTitle)

	assert.Equal(t, 3, cfg.Views.Diff.ContextLines)
	assert.Equal(t, true, cfg.Views.Diff.ShowStat)
	assert.Equal(t, false, cfg.Views.Diff.IgnoreSpace)

	assert.Equal(t, true, cfg.Views.Status.ShowUntracked)
	assert.Equal(t, false, cfg.Views.Status.ShowIgnored)

	assert.NotEmpty(t, cfg.General.Editor)
	assert.Equal(t, "less", cfg.General.Pager)
	assert.Equal(t, "topo", cfg.General.CommitOrder)
	assert.Equal(t, false, cfg.General.VerticalSplit)

	assert.NotEmpty(t, cfg.Keymaps.Bindings)
	assert.NotEmpty(t, cfg.Colors.Colors)
}

func TestConfigValidation(t *testing.T) {
	cfg := &Config{
		UI: UIConfig{
			TabSize: 4,
		},
		Git: GitConfig{
			AuthorWidth: 15,
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestGetDefaultEditor(t *testing.T) {
	// Test that we get a default editor
	editor := getDefaultEditor()
	assert.NotEmpty(t, editor)
}

func TestGetConfigPaths(t *testing.T) {
	paths := GetConfigPaths()
	assert.Len(t, paths, 4)
	assert.Contains(t, paths[0], "tigrc")
	assert.Contains(t, paths[1], ".config/tig/tigrc")
	assert.Contains(t, paths[2], ".tigrc")
	assert.Contains(t, paths[3], "/etc/tig/tigrc")
}

func TestStructIntegrity(t *testing.T) {
	cfg := &Config{}

	// Test that all fields can be set
	cfg.UI.TabSize = 4
	cfg.UI.CommitOrder = "date"
	cfg.UI.IgnoreSpace = "all"
	cfg.UI.IgnoreCase = true
	cfg.UI.ShowDate = false
	cfg.UI.ShowAuthor = false
	cfg.UI.ShowLineNumbers = false

	cfg.Git.AuthorWidth = 25
	cfg.Git.DateFormat = "%Y-%m-%d %H:%M:%S"
	cfg.Git.ShowNotes = false
	cfg.Git.ShowDiffStat = false
	cfg.Git.ShowBranches = false
	cfg.Git.ShowRemotes = false
	cfg.Git.ShowTags = false

	cfg.Views.Main.ShowGraph = false
	cfg.Views.Main.ShowRefs = false
	cfg.Views.Main.ShowID = true
	cfg.Views.Main.ShowDate = false
	cfg.Views.Main.ShowAuthor = false
	cfg.Views.Main.ShowCommitTitle = false

	cfg.Views.Diff.ContextLines = 5
	cfg.Views.Diff.ShowStat = false
	cfg.Views.Diff.IgnoreSpace = true

	cfg.Views.Status.ShowUntracked = false
	cfg.Views.Status.ShowIgnored = true

	cfg.General.Editor = "vim"
	cfg.General.Pager = "more"
	cfg.General.CommitOrder = "date"
	cfg.General.VerticalSplit = true

	cfg.Keymaps.Bindings = map[string]string{
		"test": "T",
	}

	cfg.Colors.Scheme = "custom"
	cfg.Colors.Colors = map[string]string{
		"test": "blue",
	}

	// Verify all values were set correctly
	assert.Equal(t, 4, cfg.UI.TabSize)
	assert.Equal(t, "date", cfg.UI.CommitOrder)
	assert.Equal(t, "all", cfg.UI.IgnoreSpace)
	assert.Equal(t, true, cfg.UI.IgnoreCase)
	assert.Equal(t, false, cfg.UI.ShowDate)
	assert.Equal(t, false, cfg.UI.ShowAuthor)
	assert.Equal(t, false, cfg.UI.ShowLineNumbers)

	assert.Equal(t, 25, cfg.Git.AuthorWidth)
	assert.Equal(t, "%Y-%m-%d %H:%M:%S", cfg.Git.DateFormat)
	assert.Equal(t, false, cfg.Git.ShowNotes)
	assert.Equal(t, false, cfg.Git.ShowDiffStat)
	assert.Equal(t, false, cfg.Git.ShowBranches)
	assert.Equal(t, false, cfg.Git.ShowRemotes)
	assert.Equal(t, false, cfg.Git.ShowTags)

	assert.Equal(t, false, cfg.Views.Main.ShowGraph)
	assert.Equal(t, false, cfg.Views.Main.ShowRefs)
	assert.Equal(t, true, cfg.Views.Main.ShowID)
	assert.Equal(t, false, cfg.Views.Main.ShowDate)
	assert.Equal(t, false, cfg.Views.Main.ShowAuthor)
	assert.Equal(t, false, cfg.Views.Main.ShowCommitTitle)

	assert.Equal(t, 5, cfg.Views.Diff.ContextLines)
	assert.Equal(t, false, cfg.Views.Diff.ShowStat)
	assert.Equal(t, true, cfg.Views.Diff.IgnoreSpace)

	assert.Equal(t, false, cfg.Views.Status.ShowUntracked)
	assert.Equal(t, true, cfg.Views.Status.ShowIgnored)

	assert.Equal(t, "vim", cfg.General.Editor)
	assert.Equal(t, "more", cfg.General.Pager)
	assert.Equal(t, "date", cfg.General.CommitOrder)
	assert.Equal(t, true, cfg.General.VerticalSplit)

	assert.Equal(t, "T", cfg.Keymaps.Bindings["test"])
	assert.Equal(t, "custom", cfg.Colors.Scheme)
	assert.Equal(t, "blue", cfg.Colors.Colors["test"])
}

func TestEmptyConfig(t *testing.T) {
	cfg := &Config{}
	
	// Test that we can access all nested fields without panicking
	_ = cfg.UI.TabSize
	_ = cfg.UI.CommitOrder
	_ = cfg.UI.IgnoreSpace
	_ = cfg.UI.IgnoreCase
	_ = cfg.UI.ShowDate
	_ = cfg.UI.ShowAuthor
	_ = cfg.UI.ShowLineNumbers

	_ = cfg.Git.AuthorWidth
	_ = cfg.Git.DateFormat
	_ = cfg.Git.ShowNotes
	_ = cfg.Git.ShowDiffStat
	_ = cfg.Git.ShowBranches
	_ = cfg.Git.ShowRemotes
	_ = cfg.Git.ShowTags

	_ = cfg.Views.Main.ShowGraph
	_ = cfg.Views.Main.ShowRefs
	_ = cfg.Views.Main.ShowID
	_ = cfg.Views.Main.ShowDate
	_ = cfg.Views.Main.ShowAuthor
	_ = cfg.Views.Main.ShowCommitTitle

	_ = cfg.Views.Diff.ContextLines
	_ = cfg.Views.Diff.ShowStat
	_ = cfg.Views.Diff.IgnoreSpace

	_ = cfg.Views.Status.ShowUntracked
	_ = cfg.Views.Status.ShowIgnored

	_ = cfg.General.Editor
	_ = cfg.General.Pager
	_ = cfg.General.CommitOrder
	_ = cfg.General.VerticalSplit

	_ = cfg.Keymaps.Bindings
	_ = cfg.Colors.Scheme
	_ = cfg.Colors.Colors
}

func TestDefaultValues(t *testing.T) {
	cfg := &Config{}
	setDefaults(cfg)

	// Test that defaults are actually set
	assert.NotZero(t, cfg.UI.TabSize)
	assert.NotEmpty(t, cfg.UI.CommitOrder)
	assert.NotEmpty(t, cfg.Git.DateFormat)
	assert.NotEmpty(t, cfg.General.Editor)
	assert.NotEmpty(t, cfg.Keymaps.Bindings)
	assert.NotEmpty(t, cfg.Colors.Colors)
}

func TestKeymapBindings(t *testing.T) {
	cfg := &Config{}
	setDefaults(cfg)

	// Test that keymap bindings contain expected keys
	assert.Contains(t, cfg.Keymaps.Bindings, "quit")
	assert.Contains(t, cfg.Keymaps.Bindings, "refresh")
	assert.Contains(t, cfg.Keymaps.Bindings, "status")
	assert.Contains(t, cfg.Keymaps.Bindings, "diff")
	assert.Contains(t, cfg.Keymaps.Bindings, "log")
	assert.Contains(t, cfg.Keymaps.Bindings, "tree")
	assert.Contains(t, cfg.Keymaps.Bindings, "refs")
	assert.Contains(t, cfg.Keymaps.Bindings, "help")
	assert.Contains(t, cfg.Keymaps.Bindings, "stage")
	assert.Contains(t, cfg.Keymaps.Bindings, "unstage")
	assert.Contains(t, cfg.Keymaps.Bindings, "commit")
}

func TestColorSchemes(t *testing.T) {
	cfg := &Config{}
	setDefaults(cfg)

	// Test that color scheme contains expected keys
	assert.Contains(t, cfg.Colors.Colors, "default")
	assert.Contains(t, cfg.Colors.Colors, "cursor")
	assert.Contains(t, cfg.Colors.Colors, "status")
	assert.Contains(t, cfg.Colors.Colors, "error")
	assert.Contains(t, cfg.Colors.Colors, "diff-header")
	assert.Contains(t, cfg.Colors.Colors, "diff-add")
	assert.Contains(t, cfg.Colors.Colors, "diff-del")
	assert.Contains(t, cfg.Colors.Colors, "branch")
	assert.Contains(t, cfg.Colors.Colors, "tag")
	assert.Contains(t, cfg.Colors.Colors, "author")
	assert.Contains(t, cfg.Colors.Colors, "date")
	assert.Contains(t, cfg.Colors.Colors, "id")
}