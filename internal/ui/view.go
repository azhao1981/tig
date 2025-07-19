package ui

import (
	"github.com/gdamore/tcell/v2"
)

// ViewType represents different view types in tig
type ViewType int

const (
	ViewTypeMain ViewType = iota
	ViewTypeDiff
	ViewTypeStatus
	ViewTypeTree
	ViewTypeRefs
	ViewTypeHelp
)

// View represents a generic interface for all views
type View interface {
	// Render renders the view content to the screen
	Render(screen tcell.Screen, x, y, width, height int) error
	
	// HandleKey handles keyboard input for this view
	HandleKey(key tcell.Key, ch rune, mod tcell.ModMask) bool
	
	// Focus sets focus to this view
	Focus()
	
	// Blur removes focus from this view
	Blur()
	
	// IsFocused returns whether this view has focus
	IsFocused() bool
	
	// GetType returns the view type
	GetType() ViewType
	
	// Refresh refreshes the view content
	Refresh() error
	
	// SetPosition sets the view position and size
	SetPosition(x, y, width, height int)
	
	// GetPosition returns the view position and size
	GetPosition() (x, y, width, height int)
}

// BaseView provides common functionality for all views
type BaseView struct {
	x      int
	y      int
	width  int
	height int
	focused bool
	viewType ViewType
}

// NewBaseView creates a new base view
func NewBaseView(viewType ViewType) *BaseView {
	return &BaseView{
		viewType: viewType,
	}
}

// Focus sets focus to this view
func (v *BaseView) Focus() {
	v.focused = true
}

// Blur removes focus from this view
func (v *BaseView) Blur() {
	v.focused = false
}

// IsFocused returns whether this view has focus
func (v *BaseView) IsFocused() bool {
	return v.focused
}

// GetType returns the view type
func (v *BaseView) GetType() ViewType {
	return v.viewType
}

// SetPosition sets the view position and size
func (v *BaseView) SetPosition(x, y, width, height int) {
	v.x = x
	v.y = y
	v.width = width
	v.height = height
}

// GetPosition returns the view position and size
func (v *BaseView) GetPosition() (int, int, int, int) {
	return v.x, v.y, v.width, v.height
}

// DrawBox draws a box around the view area
type DrawBox struct {
	Title string
	Style tcell.Style
}

// NewDrawBox creates a new draw box
func NewDrawBox(title string, style tcell.Style) *DrawBox {
	return &DrawBox{
		Title: title,
		Style: style,
	}
}

// Draw draws the box
func (db *DrawBox) Draw(screen tcell.Screen, x, y, width, height int) {
	if width <= 0 || height <= 0 {
		return
	}

	// Draw borders
	for i := 0; i < width; i++ {
		screen.SetContent(x+i, y, tcell.RuneHLine, nil, db.Style)
		if height > 1 {
			screen.SetContent(x+i, y+height-1, tcell.RuneHLine, nil, db.Style)
		}
	}

	for i := 0; i < height; i++ {
		screen.SetContent(x, y+i, tcell.RuneVLine, nil, db.Style)
		if width > 1 {
			screen.SetContent(x+width-1, y+i, tcell.RuneVLine, nil, db.Style)
		}
	}

	// Draw corners
	screen.SetContent(x, y, tcell.RuneULCorner, nil, db.Style)
	screen.SetContent(x+width-1, y, tcell.RuneURCorner, nil, db.Style)
	if height > 1 {
		screen.SetContent(x, y+height-1, tcell.RuneLLCorner, nil, db.Style)
		screen.SetContent(x+width-1, y+height-1, tcell.RuneLRCorner, nil, db.Style)
	}

	// Draw title
	if db.Title != "" && width > 2 {
		title := db.Title
		if len(title) > width-2 {
			title = title[:width-2]
		}
		titleX := x + (width-len(title))/2
		for i, char := range title {
			screen.SetContent(titleX+i, y, char, nil, db.Style)
		}
	}
}

// Scrollable provides scrolling functionality
type Scrollable struct {
	offset    int
	maxOffset int
	height    int
}

// NewScrollable creates a new scrollable
func NewScrollable() *Scrollable {
	return &Scrollable{}
}

// SetMaxOffset sets the maximum scroll offset
func (s *Scrollable) SetMaxOffset(max int) {
	s.maxOffset = max
	if s.offset > max {
		s.offset = max
	}
}

// SetHeight sets the visible height
func (s *Scrollable) SetHeight(height int) {
	s.height = height
}

// GetOffset returns the current scroll offset
func (s *Scrollable) GetOffset() int {
	return s.offset
}

// SetOffset sets the current scroll offset
func (s *Scrollable) SetOffset(offset int) {
	s.offset = offset
	if s.offset < 0 {
		s.offset = 0
	}
	if s.offset > s.maxOffset {
		s.offset = s.maxOffset
	}
}

// ScrollUp scrolls up by one line
func (s *Scrollable) ScrollUp() {
	if s.offset > 0 {
		s.offset--
	}
}

// ScrollDown scrolls down by one line
func (s *Scrollable) ScrollDown() {
	if s.offset < s.maxOffset {
		s.offset++
	}
}

// ScrollPageUp scrolls up by one page
func (s *Scrollable) ScrollPageUp() {
	s.offset -= s.height
	if s.offset < 0 {
		s.offset = 0
	}
}

// ScrollPageDown scrolls down by one page
func (s *Scrollable) ScrollPageDown() {
	s.offset += s.height
	if s.offset > s.maxOffset {
		s.offset = s.maxOffset
	}
}

// ScrollToTop scrolls to the top
func (s *Scrollable) ScrollToTop() {
	s.offset = 0
}

// ScrollToBottom scrolls to the bottom
func (s *Scrollable) ScrollToBottom() {
	s.offset = s.maxOffset
}

// IsAtTop returns whether the view is at the top
func (s *Scrollable) IsAtTop() bool {
	return s.offset == 0
}

// IsAtBottom returns whether the view is at the bottom
func (s *Scrollable) IsAtBottom() bool {
	return s.offset == s.maxOffset
}