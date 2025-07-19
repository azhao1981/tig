package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/jonas/tig/internal/config"
)

// Theme represents a color theme for the application
type Theme struct {
	colors map[string]tcell.Color
}

// NewTheme creates a new theme from configuration
func NewTheme(config *config.Config) *Theme {
	theme := &Theme{
		colors: make(map[string]tcell.Color),
	}
	theme.loadFromConfig(config)
	return theme
}

// loadFromConfig loads colors from the configuration
func (t *Theme) loadFromConfig(config *config.Config) {
	// Map color names to tcell colors
	colorMap := map[string]tcell.Color{
		"black":   tcell.ColorBlack,
		"red":     tcell.ColorRed,
		"green":   tcell.ColorGreen,
		"yellow":  tcell.ColorYellow,
		"blue":    tcell.ColorBlue,
		"magenta": tcell.ColorFuchsia,
		"cyan":    tcell.ColorAqua,
		"white":   tcell.ColorWhite,
		"gray":    tcell.ColorGray,
		"darkgray": tcell.ColorDarkGray,
		"lightgray": tcell.ColorLightGray,
	}

	// Load custom colors from config
	for key, colorName := range config.Colors.Colors {
		if color, ok := colorMap[strings.ToLower(colorName)]; ok {
			t.colors[key] = color
		} else {
			// Default to white if color not found
			t.colors[key] = tcell.ColorWhite
		}
	}

	// Set defaults for missing colors
	t.setDefaults()
}

// setDefaults sets default colors for missing theme entries
func (t *Theme) setDefaults() {
	defaults := map[string]tcell.Color{
		"default":       tcell.ColorWhite,
		"cursor":        tcell.ColorYellow,
		"status":        tcell.ColorGreen,
		"error":         tcell.ColorRed,
		"diff-header":   tcell.ColorFuchsia,
		"diff-add":      tcell.ColorGreen,
		"diff-del":      tcell.ColorRed,
		"branch":        tcell.ColorFuchsia,
		"tag":           tcell.ColorYellow,
		"author":        tcell.ColorAqua,
		"date":          tcell.ColorGreen,
		"id":            tcell.ColorBlue,
		"header":        tcell.ColorWhite,
		"line-number":   tcell.ColorDarkGray,
		"selected":      tcell.ColorYellow,
		"directory":     tcell.ColorBlue,
		"file":          tcell.ColorWhite,
		"binary":        tcell.ColorRed,
		"staged":        tcell.ColorGreen,
		"modified":      tcell.ColorYellow,
		"untracked":     tcell.ColorRed,
		"conflict":      tcell.ColorRed,
	}

	for key, defaultColor := range defaults {
		if _, exists := t.colors[key]; !exists {
			t.colors[key] = defaultColor
		}
	}
}

// GetColor returns a color by name
func (t *Theme) GetColor(name string) tcell.Color {
	if color, ok := t.colors[name]; ok {
		return color
	}
	return tcell.ColorWhite
}

// GetStyle returns a style with the specified color
func (t *Theme) GetStyle(colorName string) tcell.Style {
	color := t.GetColor(colorName)
	return tcell.StyleDefault.Foreground(color)
}

// GetStyleWithBackground returns a style with foreground and background colors
func (t *Theme) GetStyleWithBackground(foreground, background string) tcell.Style {
	fg := t.GetColor(foreground)
	bg := t.GetColor(background)
	return tcell.StyleDefault.Foreground(fg).Background(bg)
}

// GetSelectedStyle returns the style for selected items
func (t *Theme) GetSelectedStyle() tcell.Style {
	return t.GetStyleWithBackground("cursor", "blue")
}

// GetStatusBarStyle returns the style for the status bar
func (t *Theme) GetStatusBarStyle() tcell.Style {
	return t.GetStyleWithBackground("default", "darkgray")
}

// GetHeaderStyle returns the style for headers
func (t *Theme) GetHeaderStyle() tcell.Style {
	return t.GetStyle("header").Bold(true)
}

// GetDiffStyle returns styles for diff display
func (t *Theme) GetDiffStyle() map[string]tcell.Style {
	return map[string]tcell.Style{
		"header":  t.GetStyle("diff-header").Bold(true),
		"add":     t.GetStyle("diff-add"),
		"del":     t.GetStyle("diff-del"),
		"context": t.GetStyle("default"),
		"line":    t.GetStyle("line-number"),
	}
}

// GetFileStyle returns styles for file display
func (t *Theme) GetFileStyle(isDir bool) tcell.Style {
	if isDir {
		return t.GetStyle("directory").Bold(true)
	}
	return t.GetStyle("file")
}

// GetStatusStyle returns styles for status display
func (t *Theme) GetStatusStyle(status string) tcell.Style {
	switch status {
	case "M", "modified":
		return t.GetStyle("modified")
	case "A", "added":
		return t.GetStyle("staged")
	case "D", "deleted":
		return t.GetStyle("diff-del")
	case "R", "renamed":
		return t.GetStyle("directory")
	case "C", "copied":
		return t.GetStyle("staged")
	case "?", "untracked":
		return t.GetStyle("untracked")
	case "U", "conflict":
		return t.GetStyle("conflict")
	default:
		return t.GetStyle("default")
	}
}

// AvailableSchemes returns available color schemes
func AvailableSchemes() []string {
	return []string{
		"default",
		"dark",
		"light",
		"monochrome",
		"solarized-dark",
		"solarized-light",
	}
}

// LoadScheme loads a predefined color scheme
func LoadScheme(schemeName string) *Theme {
	config := &config.Config{Colors: config.ColorConfig{Scheme: schemeName}}
	
	// Define predefined schemes
	schemes := map[string]map[string]string{
		"dark": {
			"default":       "lightgray",
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
			"header":        "white",
			"line-number":   "darkgray",
			"selected":      "yellow",
			"directory":     "blue",
			"file":          "lightgray",
			"binary":        "red",
			"staged":        "green",
			"modified":      "yellow",
			"untracked":     "red",
			"conflict":      "red",
		},
		"light": {
			"default":       "black",
			"cursor":        "blue",
			"status":        "green",
			"error":         "red",
			"diff-header":   "magenta",
			"diff-add":      "green",
			"diff-del":      "red",
			"branch":        "magenta",
			"tag":           "brown",
			"author":        "blue",
			"date":          "green",
			"id":            "blue",
			"header":        "black",
			"line-number":   "gray",
			"selected":      "blue",
			"directory":     "blue",
			"file":          "black",
			"binary":        "red",
			"staged":        "green",
			"modified":      "brown",
			"untracked":     "red",
			"conflict":      "red",
		},
		"monochrome": {
			"default":       "white",
			"cursor":        "white",
			"status":        "white",
			"error":         "white",
			"diff-header":   "white",
			"diff-add":      "white",
			"diff-del":      "white",
			"branch":        "white",
			"tag":           "white",
			"author":        "white",
			"date":          "white",
			"id":            "white",
			"header":        "white",
			"line-number":   "white",
			"selected":      "white",
			"directory":     "white",
			"file":          "white",
			"binary":        "white",
			"staged":        "white",
			"modified":      "white",
			"untracked":     "white",
			"conflict":      "white",
		},
	}
	
	if scheme, ok := schemes[schemeName]; ok {
		config.Colors.Colors = scheme
	}
	
	return NewTheme(config)
}