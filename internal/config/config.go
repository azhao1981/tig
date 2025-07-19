package config

import (
	"os"
	"path/filepath"
)

// Config represents the main configuration structure
type Config struct {
	UI        UIConfig        `mapstructure:"ui"`
	Git       GitConfig       `mapstructure:"git"`
	Keymaps   KeymapConfig    `mapstructure:"keymaps"`
	Colors    ColorConfig     `mapstructure:"colors"`
	Views     ViewsConfig     `mapstructure:"views"`
	General   GeneralConfig   `mapstructure:"general"`
}

// UIConfig holds UI-related configuration
type UIConfig struct {
	TabSize       int    `mapstructure:"tab_size"`
	CommitOrder   string `mapstructure:"commit_order"`
	IgnoreSpace   string `mapstructure:"ignore_space"`
	IgnoreCase    bool   `mapstructure:"ignore_case"`
	ShowDate      bool   `mapstructure:"show_date"`
	ShowAuthor    bool   `mapstructure:"show_author"`
	ShowLineNumbers bool `mapstructure:"show_line_numbers"`
}

// GitConfig holds Git-related configuration
type GitConfig struct {
	AuthorWidth   int    `mapstructure:"author_width"`
	DateFormat    string `mapstructure:"date_format"`
	ShowNotes     bool   `mapstructure:"show_notes"`
	ShowDiffStat  bool   `mapstructure:"show_diff_stat"`
	ShowBranches  bool   `mapstructure:"show_branches"`
	ShowRemotes   bool   `mapstructure:"show_remotes"`
	ShowTags      bool   `mapstructure:"show_tags"`
}

// KeymapConfig holds key binding configuration
type KeymapConfig struct {
	Bindings map[string]string `mapstructure:"bindings"`
}

// ColorConfig holds color scheme configuration
type ColorConfig struct {
	Scheme string            `mapstructure:"scheme"`
	Colors map[string]string `mapstructure:"colors"`
}

// ViewsConfig holds view-specific configuration
type ViewsConfig struct {
	Main   MainViewConfig   `mapstructure:"main"`
	Diff   DiffViewConfig   `mapstructure:"diff"`
	Status StatusViewConfig `mapstructure:"status"`
}

// MainViewConfig holds main view configuration
type MainViewConfig struct {
	ShowGraph    bool `mapstructure:"show_graph"`
	ShowRefs     bool `mapstructure:"show_refs"`
	ShowID       bool `mapstructure:"show_id"`
	ShowDate     bool `mapstructure:"show_date"`
	ShowAuthor   bool `mapstructure:"show_author"`
	ShowCommitTitle bool `mapstructure:"show_commit_title"`
}

// DiffViewConfig holds diff view configuration
type DiffViewConfig struct {
	ContextLines int  `mapstructure:"context_lines"`
	ShowStat     bool `mapstructure:"show_stat"`
	IgnoreSpace  bool `mapstructure:"ignore_space"`
}

// StatusViewConfig holds status view configuration
type StatusViewConfig struct {
	ShowUntracked bool `mapstructure:"show_untracked"`
	ShowIgnored   bool `mapstructure:"show_ignored"`
}

// GeneralConfig holds general configuration
type GeneralConfig struct {
	Editor          string `mapstructure:"editor"`
	Pager           string `mapstructure:"pager"`
	Terminal        string `mapstructure:"terminal"`
	CommitOrder     string `mapstructure:"commit_order"`
	VerticalSplit   bool   `mapstructure:"vertical_split"`
}

// Load loads configuration from tigrc files and environment variables
func Load() (*Config, error) {
	config := &Config{}

	// Set default configuration
	setDefaults(config)

	// For now, just return the default configuration
	// TODO: Implement configuration file parsing
	return config, nil
}

// setDefaults sets default configuration values
func setDefaults(config *Config) {
	// UI defaults
	config.UI.TabSize = 8
	config.UI.CommitOrder = "topo"
	config.UI.IgnoreSpace = "no"
	config.UI.IgnoreCase = false
	config.UI.ShowDate = true
	config.UI.ShowAuthor = true
	config.UI.ShowLineNumbers = true

	// Git defaults
	config.Git.AuthorWidth = 20
	config.Git.DateFormat = "%Y-%m-%d"
	config.Git.ShowNotes = true
	config.Git.ShowDiffStat = true
	config.Git.ShowBranches = true
	config.Git.ShowRemotes = true
	config.Git.ShowTags = true

	// Views defaults
	config.Views.Main.ShowGraph = true
	config.Views.Main.ShowRefs = true
	config.Views.Main.ShowID = false
	config.Views.Main.ShowDate = true
	config.Views.Main.ShowAuthor = true
	config.Views.Main.ShowCommitTitle = true

	config.Views.Diff.ContextLines = 3
	config.Views.Diff.ShowStat = true
	config.Views.Diff.IgnoreSpace = false

	config.Views.Status.ShowUntracked = true
	config.Views.Status.ShowIgnored = false

	// General defaults
	config.General.Editor = getDefaultEditor()
	config.General.Pager = "less"
	config.General.CommitOrder = "topo"
	config.General.VerticalSplit = false

	// Keymaps defaults
	config.Keymaps.Bindings = map[string]string{
		"quit":            "q",
		"refresh":         "R",
		"status":          "s",
		"diff":            "d",
		"log":             "l",
		"tree":            "t",
		"refs":            "r",
		"help":            "h",
		"stage":           "u",
		"unstage":         "U",
		"commit":          "c",
	}

	// Colors defaults
	config.Colors.Scheme = "default"
	config.Colors.Colors = map[string]string{
		"default":       "white",
		"cursor":        "yellow",
		"status":        "green",
		"error":         "red",
		"diff-header":   "magenta",
		"diff-add":      "green",
		"diff-del":      "red",
		"branch":        "magenta",
		"tag":           "yellow",
		"author":        "cyan",
		"date":          "green",
		"id":            "blue",
	}
}

// getDefaultEditor returns the default editor from environment or system
func getDefaultEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	return "vi"
}

// GetConfigPaths returns the search paths for configuration files
func GetConfigPaths() []string {
	return []string{
		filepath.Join(".", "tigrc"),
		filepath.Join("$HOME", ".config", "tig", "tigrc"),
		filepath.Join("$HOME", ".tigrc"),
		filepath.Join("/etc", "tig", "tigrc"),
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here
	return nil
}

// Save saves the configuration to the default location
func (c *Config) Save() error {
	// Add save logic here
	return nil
}