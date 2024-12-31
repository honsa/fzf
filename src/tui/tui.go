package tui

import (
	"strconv"
	"time"

	"github.com/junegunn/fzf/src/util"
	"github.com/rivo/uniseg"
)

// Types of user action
//
//go:generate stringer -type=EventType
type EventType int

const (
	Rune EventType = iota

	CtrlA
	CtrlB
	CtrlC
	CtrlD
	CtrlE
	CtrlF
	CtrlG
	CtrlH
	Tab
	CtrlJ
	CtrlK
	CtrlL
	CtrlM
	CtrlN
	CtrlO
	CtrlP
	CtrlQ
	CtrlR
	CtrlS
	CtrlT
	CtrlU
	CtrlV
	CtrlW
	CtrlX
	CtrlY
	CtrlZ
	Esc
	CtrlSpace
	CtrlDelete

	// https://apple.stackexchange.com/questions/24261/how-do-i-send-c-that-is-control-slash-to-the-terminal
	CtrlBackSlash
	CtrlRightBracket
	CtrlCaret
	CtrlSlash

	ShiftTab
	Backspace

	Delete
	PageUp
	PageDown

	Up
	Down
	Left
	Right
	Home
	End
	Insert

	ShiftUp
	ShiftDown
	ShiftLeft
	ShiftRight
	ShiftDelete

	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12

	AltBackspace

	AltUp
	AltDown
	AltLeft
	AltRight

	AltShiftUp
	AltShiftDown
	AltShiftLeft
	AltShiftRight

	Alt
	CtrlAlt

	Invalid
	Fatal

	Mouse
	DoubleClick
	LeftClick
	RightClick
	SLeftClick
	SRightClick
	ScrollUp
	ScrollDown
	SScrollUp
	SScrollDown
	PreviewScrollUp
	PreviewScrollDown

	// Events
	Resize
	Change
	BackwardEOF
	Start
	Load
	Focus
	One
	Zero
	Result
	Jump
	JumpCancel
	ClickHeader
)

func (t EventType) AsEvent() Event {
	return Event{t, 0, nil}
}

func (t EventType) Int() int {
	return int(t)
}

func (t EventType) Byte() byte {
	return byte(t)
}

func (e Event) Comparable() Event {
	// Ignore MouseEvent pointer
	return Event{e.Type, e.Char, nil}
}

func (e Event) KeyName() string {
	if e.Type >= Invalid {
		return ""
	}

	switch e.Type {
	case Rune:
		return string(e.Char)
	case Alt:
		return "alt-" + string(e.Char)
	case CtrlAlt:
		return "ctrl-alt-" + string(e.Char)
	case CtrlBackSlash:
		return "ctrl-\\"
	case CtrlRightBracket:
		return "ctrl-]"
	case CtrlCaret:
		return "ctrl-^"
	case CtrlSlash:
		return "ctrl-/"
	}

	return util.ToKebabCase(e.Type.String())
}

func Key(r rune) Event {
	return Event{Rune, r, nil}
}

func AltKey(r rune) Event {
	return Event{Alt, r, nil}
}

func CtrlAltKey(r rune) Event {
	return Event{CtrlAlt, r, nil}
}

const (
	doubleClickDuration = 500 * time.Millisecond
)

type Color int32

func (c Color) IsDefault() bool {
	return c == colDefault
}

func (c Color) is24() bool {
	return c > 0 && (c&(1<<24)) > 0
}

type ColorAttr struct {
	Color Color
	Attr  Attr
}

func NewColorAttr() ColorAttr {
	return ColorAttr{Color: colUndefined, Attr: AttrUndefined}
}

const (
	colUndefined Color = -2
	colDefault   Color = -1
)

const (
	colBlack Color = iota
	colRed
	colGreen
	colYellow
	colBlue
	colMagenta
	colCyan
	colWhite
)

type FillReturn int

const (
	FillContinue FillReturn = iota
	FillNextLine
	FillSuspend
)

type ColorPair struct {
	fg   Color
	bg   Color
	attr Attr
}

func HexToColor(rrggbb string) Color {
	r, _ := strconv.ParseInt(rrggbb[1:3], 16, 0)
	g, _ := strconv.ParseInt(rrggbb[3:5], 16, 0)
	b, _ := strconv.ParseInt(rrggbb[5:7], 16, 0)
	return Color((1 << 24) + (r << 16) + (g << 8) + b)
}

func NewColorPair(fg Color, bg Color, attr Attr) ColorPair {
	return ColorPair{fg, bg, attr}
}

func (p ColorPair) Fg() Color {
	return p.fg
}

func (p ColorPair) Bg() Color {
	return p.bg
}

func (p ColorPair) Attr() Attr {
	return p.attr
}

func (p ColorPair) HasBg() bool {
	return p.attr&Reverse == 0 && p.bg != colDefault ||
		p.attr&Reverse > 0 && p.fg != colDefault
}

func (p ColorPair) merge(other ColorPair, except Color) ColorPair {
	dup := p
	dup.attr = dup.attr.Merge(other.attr)
	if other.fg != except {
		dup.fg = other.fg
	}
	if other.bg != except {
		dup.bg = other.bg
	}
	return dup
}

func (p ColorPair) WithAttr(attr Attr) ColorPair {
	dup := p
	dup.attr = dup.attr.Merge(attr)
	return dup
}

func (p ColorPair) MergeAttr(other ColorPair) ColorPair {
	return p.WithAttr(other.attr)
}

func (p ColorPair) Merge(other ColorPair) ColorPair {
	return p.merge(other, colUndefined)
}

func (p ColorPair) MergeNonDefault(other ColorPair) ColorPair {
	return p.merge(other, colDefault)
}

type ColorTheme struct {
	Colored          bool
	Input            ColorAttr
	Disabled         ColorAttr
	Fg               ColorAttr
	Bg               ColorAttr
	ListFg           ColorAttr
	ListBg           ColorAttr
	SelectedFg       ColorAttr
	SelectedBg       ColorAttr
	SelectedMatch    ColorAttr
	PreviewFg        ColorAttr
	PreviewBg        ColorAttr
	DarkBg           ColorAttr
	Gutter           ColorAttr
	Prompt           ColorAttr
	Match            ColorAttr
	Current          ColorAttr
	CurrentMatch     ColorAttr
	Spinner          ColorAttr
	Info             ColorAttr
	Cursor           ColorAttr
	Marker           ColorAttr
	Header           ColorAttr
	Separator        ColorAttr
	Scrollbar        ColorAttr
	Border           ColorAttr
	PreviewBorder    ColorAttr
	PreviewLabel     ColorAttr
	PreviewScrollbar ColorAttr
	BorderLabel      ColorAttr
	ListLabel        ColorAttr
	ListBorder       ColorAttr
}

type Event struct {
	Type       EventType
	Char       rune
	MouseEvent *MouseEvent
}

type MouseEvent struct {
	Y      int
	X      int
	S      int
	Left   bool
	Down   bool
	Double bool
	Mod    bool
}

type BorderShape int

const (
	BorderUndefined BorderShape = iota
	BorderNone
	BorderRounded
	BorderSharp
	BorderBold
	BorderBlock
	BorderThinBlock
	BorderDouble
	BorderHorizontal
	BorderVertical
	BorderTop
	BorderBottom
	BorderLeft
	BorderRight
)

func (s BorderShape) HasLeft() bool {
	switch s {
	case BorderNone, BorderRight, BorderTop, BorderBottom, BorderHorizontal: // No Left
		return false
	}
	return true
}

func (s BorderShape) HasRight() bool {
	switch s {
	case BorderNone, BorderLeft, BorderTop, BorderBottom, BorderHorizontal: // No right
		return false
	}
	return true
}

func (s BorderShape) HasTop() bool {
	switch s {
	case BorderNone, BorderLeft, BorderRight, BorderBottom, BorderVertical: // No top
		return false
	}
	return true
}

func (s BorderShape) HasBottom() bool {
	switch s {
	case BorderNone, BorderLeft, BorderRight, BorderTop, BorderVertical: // No bottom
		return false
	}
	return true
}

func (s BorderShape) Visible() bool {
	return s != BorderNone
}

type BorderStyle struct {
	shape       BorderShape
	top         rune
	bottom      rune
	left        rune
	right       rune
	topLeft     rune
	topRight    rune
	bottomLeft  rune
	bottomRight rune
}

type BorderCharacter int

func MakeBorderStyle(shape BorderShape, unicode bool) BorderStyle {
	if shape == BorderNone {
		return BorderStyle{
			shape:       shape,
			top:         ' ',
			bottom:      ' ',
			left:        ' ',
			right:       ' ',
			topLeft:     ' ',
			topRight:    ' ',
			bottomLeft:  ' ',
			bottomRight: ' '}
	}
	if !unicode {
		return BorderStyle{
			shape:       shape,
			top:         '-',
			bottom:      '-',
			left:        '|',
			right:       '|',
			topLeft:     '+',
			topRight:    '+',
			bottomLeft:  '+',
			bottomRight: '+',
		}
	}
	switch shape {
	case BorderSharp:
		return BorderStyle{
			shape:       shape,
			top:         '─',
			bottom:      '─',
			left:        '│',
			right:       '│',
			topLeft:     '┌',
			topRight:    '┐',
			bottomLeft:  '└',
			bottomRight: '┘',
		}
	case BorderBold:
		return BorderStyle{
			shape:       shape,
			top:         '━',
			bottom:      '━',
			left:        '┃',
			right:       '┃',
			topLeft:     '┏',
			topRight:    '┓',
			bottomLeft:  '┗',
			bottomRight: '┛',
		}
	case BorderBlock:
		// ▛▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▜
		// ▌                  ▐
		// ▙▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▟
		return BorderStyle{
			shape:       shape,
			top:         '▀',
			bottom:      '▄',
			left:        '▌',
			right:       '▐',
			topLeft:     '▛',
			topRight:    '▜',
			bottomLeft:  '▙',
			bottomRight: '▟',
		}

	case BorderThinBlock:
		// 🭽▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔🭾
		// ▏                  ▕
		// 🭼▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁🭿
		return BorderStyle{
			shape:       shape,
			top:         '▔',
			bottom:      '▁',
			left:        '▏',
			right:       '▕',
			topLeft:     '🭽',
			topRight:    '🭾',
			bottomLeft:  '🭼',
			bottomRight: '🭿',
		}

	case BorderDouble:
		return BorderStyle{
			shape:       shape,
			top:         '═',
			bottom:      '═',
			left:        '║',
			right:       '║',
			topLeft:     '╔',
			topRight:    '╗',
			bottomLeft:  '╚',
			bottomRight: '╝',
		}
	}
	return BorderStyle{
		shape:       shape,
		top:         '─',
		bottom:      '─',
		left:        '│',
		right:       '│',
		topLeft:     '╭',
		topRight:    '╮',
		bottomLeft:  '╰',
		bottomRight: '╯',
	}
}

type TermSize struct {
	Lines    int
	Columns  int
	PxWidth  int
	PxHeight int
}

type WindowType int

const (
	WindowBase WindowType = iota
	WindowList
	WindowPreview
)

type Renderer interface {
	DefaultTheme() *ColorTheme
	Init() error
	Resize(maxHeightFunc func(int) int)
	Pause(clear bool)
	Resume(clear bool, sigcont bool)
	Clear()
	RefreshWindows(windows []Window)
	Refresh()
	Close()
	PassThrough(string)
	NeedScrollbarRedraw() bool
	ShouldEmitResizeEvent() bool

	GetChar() Event

	Top() int
	MaxX() int
	MaxY() int

	Size() TermSize

	NewWindow(top int, left int, width int, height int, windowType WindowType, borderStyle BorderStyle, erase bool) Window
}

type Window interface {
	Top() int
	Left() int
	Width() int
	Height() int

	DrawBorder()
	DrawHBorder()
	Refresh()
	FinishFill()

	X() int
	Y() int
	Enclose(y int, x int) bool

	Move(y int, x int)
	MoveAndClear(y int, x int)
	Print(text string)
	CPrint(color ColorPair, text string)
	Fill(text string) FillReturn
	CFill(fg Color, bg Color, attr Attr, text string) FillReturn
	LinkBegin(uri string, params string)
	LinkEnd()
	Erase()
	EraseMaybe() bool
}

type FullscreenRenderer struct {
	theme        *ColorTheme
	mouse        bool
	forceBlack   bool
	prevDownTime time.Time
	clicks       [][2]int
}

func NewFullscreenRenderer(theme *ColorTheme, forceBlack bool, mouse bool) Renderer {
	r := &FullscreenRenderer{
		theme:        theme,
		mouse:        mouse,
		forceBlack:   forceBlack,
		prevDownTime: time.Unix(0, 0),
		clicks:       [][2]int{}}
	return r
}

var (
	Default16 *ColorTheme
	Dark256   *ColorTheme
	Light256  *ColorTheme

	ColPrompt               ColorPair
	ColNormal               ColorPair
	ColInput                ColorPair
	ColDisabled             ColorPair
	ColMatch                ColorPair
	ColCursor               ColorPair
	ColCursorEmpty          ColorPair
	ColMarker               ColorPair
	ColSelected             ColorPair
	ColSelectedMatch        ColorPair
	ColCurrent              ColorPair
	ColCurrentMatch         ColorPair
	ColCurrentCursor        ColorPair
	ColCurrentCursorEmpty   ColorPair
	ColCurrentMarker        ColorPair
	ColCurrentSelectedEmpty ColorPair
	ColSpinner              ColorPair
	ColInfo                 ColorPair
	ColHeader               ColorPair
	ColSeparator            ColorPair
	ColScrollbar            ColorPair
	ColBorder               ColorPair
	ColPreview              ColorPair
	ColPreviewBorder        ColorPair
	ColBorderLabel          ColorPair
	ColPreviewLabel         ColorPair
	ColPreviewScrollbar     ColorPair
	ColPreviewSpinner       ColorPair
	ColListBorder           ColorPair
	ColListLabel            ColorPair
)

func EmptyTheme() *ColorTheme {
	return &ColorTheme{
		Colored:          true,
		Input:            ColorAttr{colUndefined, AttrUndefined},
		Fg:               ColorAttr{colUndefined, AttrUndefined},
		Bg:               ColorAttr{colUndefined, AttrUndefined},
		ListFg:           ColorAttr{colUndefined, AttrUndefined},
		ListBg:           ColorAttr{colUndefined, AttrUndefined},
		SelectedFg:       ColorAttr{colUndefined, AttrUndefined},
		SelectedBg:       ColorAttr{colUndefined, AttrUndefined},
		SelectedMatch:    ColorAttr{colUndefined, AttrUndefined},
		DarkBg:           ColorAttr{colUndefined, AttrUndefined},
		Prompt:           ColorAttr{colUndefined, AttrUndefined},
		Match:            ColorAttr{colUndefined, AttrUndefined},
		Current:          ColorAttr{colUndefined, AttrUndefined},
		CurrentMatch:     ColorAttr{colUndefined, AttrUndefined},
		Spinner:          ColorAttr{colUndefined, AttrUndefined},
		Info:             ColorAttr{colUndefined, AttrUndefined},
		Cursor:           ColorAttr{colUndefined, AttrUndefined},
		Marker:           ColorAttr{colUndefined, AttrUndefined},
		Header:           ColorAttr{colUndefined, AttrUndefined},
		Border:           ColorAttr{colUndefined, AttrUndefined},
		BorderLabel:      ColorAttr{colUndefined, AttrUndefined},
		ListLabel:        ColorAttr{colUndefined, AttrUndefined},
		ListBorder:       ColorAttr{colUndefined, AttrUndefined},
		Disabled:         ColorAttr{colUndefined, AttrUndefined},
		PreviewFg:        ColorAttr{colUndefined, AttrUndefined},
		PreviewBg:        ColorAttr{colUndefined, AttrUndefined},
		Gutter:           ColorAttr{colUndefined, AttrUndefined},
		PreviewBorder:    ColorAttr{colUndefined, AttrUndefined},
		PreviewScrollbar: ColorAttr{colUndefined, AttrUndefined},
		PreviewLabel:     ColorAttr{colUndefined, AttrUndefined},
		Separator:        ColorAttr{colUndefined, AttrUndefined},
		Scrollbar:        ColorAttr{colUndefined, AttrUndefined},
	}
}

func NoColorTheme() *ColorTheme {
	return &ColorTheme{
		Colored:          false,
		Input:            ColorAttr{colDefault, AttrUndefined},
		Fg:               ColorAttr{colDefault, AttrUndefined},
		Bg:               ColorAttr{colDefault, AttrUndefined},
		ListFg:           ColorAttr{colDefault, AttrUndefined},
		ListBg:           ColorAttr{colDefault, AttrUndefined},
		SelectedFg:       ColorAttr{colDefault, AttrUndefined},
		SelectedBg:       ColorAttr{colDefault, AttrUndefined},
		SelectedMatch:    ColorAttr{colDefault, AttrUndefined},
		DarkBg:           ColorAttr{colDefault, AttrUndefined},
		Prompt:           ColorAttr{colDefault, AttrUndefined},
		Match:            ColorAttr{colDefault, Underline},
		Current:          ColorAttr{colDefault, Reverse},
		CurrentMatch:     ColorAttr{colDefault, Reverse | Underline},
		Spinner:          ColorAttr{colDefault, AttrUndefined},
		Info:             ColorAttr{colDefault, AttrUndefined},
		Cursor:           ColorAttr{colDefault, AttrUndefined},
		Marker:           ColorAttr{colDefault, AttrUndefined},
		Header:           ColorAttr{colDefault, AttrUndefined},
		Border:           ColorAttr{colDefault, AttrUndefined},
		BorderLabel:      ColorAttr{colDefault, AttrUndefined},
		Disabled:         ColorAttr{colDefault, AttrUndefined},
		PreviewFg:        ColorAttr{colDefault, AttrUndefined},
		PreviewBg:        ColorAttr{colDefault, AttrUndefined},
		Gutter:           ColorAttr{colDefault, AttrUndefined},
		PreviewBorder:    ColorAttr{colDefault, AttrUndefined},
		PreviewScrollbar: ColorAttr{colDefault, AttrUndefined},
		PreviewLabel:     ColorAttr{colDefault, AttrUndefined},
		ListLabel:        ColorAttr{colDefault, AttrUndefined},
		ListBorder:       ColorAttr{colDefault, AttrUndefined},
		Separator:        ColorAttr{colDefault, AttrUndefined},
		Scrollbar:        ColorAttr{colDefault, AttrUndefined},
	}
}

func init() {
	Default16 = &ColorTheme{
		Colored:          true,
		Input:            ColorAttr{colDefault, AttrUndefined},
		Fg:               ColorAttr{colDefault, AttrUndefined},
		Bg:               ColorAttr{colDefault, AttrUndefined},
		ListFg:           ColorAttr{colUndefined, AttrUndefined},
		ListBg:           ColorAttr{colUndefined, AttrUndefined},
		SelectedFg:       ColorAttr{colUndefined, AttrUndefined},
		SelectedBg:       ColorAttr{colUndefined, AttrUndefined},
		SelectedMatch:    ColorAttr{colUndefined, AttrUndefined},
		DarkBg:           ColorAttr{colBlack, AttrUndefined},
		Prompt:           ColorAttr{colBlue, AttrUndefined},
		Match:            ColorAttr{colGreen, AttrUndefined},
		Current:          ColorAttr{colYellow, AttrUndefined},
		CurrentMatch:     ColorAttr{colGreen, AttrUndefined},
		Spinner:          ColorAttr{colGreen, AttrUndefined},
		Info:             ColorAttr{colWhite, AttrUndefined},
		Cursor:           ColorAttr{colRed, AttrUndefined},
		Marker:           ColorAttr{colMagenta, AttrUndefined},
		Header:           ColorAttr{colCyan, AttrUndefined},
		Border:           ColorAttr{colBlack, AttrUndefined},
		BorderLabel:      ColorAttr{colWhite, AttrUndefined},
		Disabled:         ColorAttr{colUndefined, AttrUndefined},
		PreviewFg:        ColorAttr{colUndefined, AttrUndefined},
		PreviewBg:        ColorAttr{colUndefined, AttrUndefined},
		Gutter:           ColorAttr{colUndefined, AttrUndefined},
		PreviewBorder:    ColorAttr{colUndefined, AttrUndefined},
		PreviewScrollbar: ColorAttr{colUndefined, AttrUndefined},
		PreviewLabel:     ColorAttr{colUndefined, AttrUndefined},
		ListLabel:        ColorAttr{colUndefined, AttrUndefined},
		ListBorder:       ColorAttr{colUndefined, AttrUndefined},
		Separator:        ColorAttr{colUndefined, AttrUndefined},
		Scrollbar:        ColorAttr{colUndefined, AttrUndefined},
	}
	Dark256 = &ColorTheme{
		Colored:          true,
		Input:            ColorAttr{colDefault, AttrUndefined},
		Fg:               ColorAttr{colDefault, AttrUndefined},
		Bg:               ColorAttr{colDefault, AttrUndefined},
		ListFg:           ColorAttr{colUndefined, AttrUndefined},
		ListBg:           ColorAttr{colUndefined, AttrUndefined},
		SelectedFg:       ColorAttr{colUndefined, AttrUndefined},
		SelectedBg:       ColorAttr{colUndefined, AttrUndefined},
		SelectedMatch:    ColorAttr{colUndefined, AttrUndefined},
		DarkBg:           ColorAttr{236, AttrUndefined},
		Prompt:           ColorAttr{110, AttrUndefined},
		Match:            ColorAttr{108, AttrUndefined},
		Current:          ColorAttr{254, AttrUndefined},
		CurrentMatch:     ColorAttr{151, AttrUndefined},
		Spinner:          ColorAttr{148, AttrUndefined},
		Info:             ColorAttr{144, AttrUndefined},
		Cursor:           ColorAttr{161, AttrUndefined},
		Marker:           ColorAttr{168, AttrUndefined},
		Header:           ColorAttr{109, AttrUndefined},
		Border:           ColorAttr{59, AttrUndefined},
		BorderLabel:      ColorAttr{145, AttrUndefined},
		Disabled:         ColorAttr{colUndefined, AttrUndefined},
		PreviewFg:        ColorAttr{colUndefined, AttrUndefined},
		PreviewBg:        ColorAttr{colUndefined, AttrUndefined},
		Gutter:           ColorAttr{colUndefined, AttrUndefined},
		PreviewBorder:    ColorAttr{colUndefined, AttrUndefined},
		PreviewScrollbar: ColorAttr{colUndefined, AttrUndefined},
		PreviewLabel:     ColorAttr{colUndefined, AttrUndefined},
		ListLabel:        ColorAttr{colUndefined, AttrUndefined},
		ListBorder:       ColorAttr{colUndefined, AttrUndefined},
		Separator:        ColorAttr{colUndefined, AttrUndefined},
		Scrollbar:        ColorAttr{colUndefined, AttrUndefined},
	}
	Light256 = &ColorTheme{
		Colored:          true,
		Input:            ColorAttr{colDefault, AttrUndefined},
		Fg:               ColorAttr{colDefault, AttrUndefined},
		Bg:               ColorAttr{colDefault, AttrUndefined},
		ListFg:           ColorAttr{colUndefined, AttrUndefined},
		ListBg:           ColorAttr{colUndefined, AttrUndefined},
		SelectedFg:       ColorAttr{colUndefined, AttrUndefined},
		SelectedBg:       ColorAttr{colUndefined, AttrUndefined},
		SelectedMatch:    ColorAttr{colUndefined, AttrUndefined},
		DarkBg:           ColorAttr{251, AttrUndefined},
		Prompt:           ColorAttr{25, AttrUndefined},
		Match:            ColorAttr{66, AttrUndefined},
		Current:          ColorAttr{237, AttrUndefined},
		CurrentMatch:     ColorAttr{23, AttrUndefined},
		Spinner:          ColorAttr{65, AttrUndefined},
		Info:             ColorAttr{101, AttrUndefined},
		Cursor:           ColorAttr{161, AttrUndefined},
		Marker:           ColorAttr{168, AttrUndefined},
		Header:           ColorAttr{31, AttrUndefined},
		Border:           ColorAttr{145, AttrUndefined},
		BorderLabel:      ColorAttr{59, AttrUndefined},
		Disabled:         ColorAttr{colUndefined, AttrUndefined},
		PreviewFg:        ColorAttr{colUndefined, AttrUndefined},
		PreviewBg:        ColorAttr{colUndefined, AttrUndefined},
		Gutter:           ColorAttr{colUndefined, AttrUndefined},
		PreviewBorder:    ColorAttr{colUndefined, AttrUndefined},
		PreviewScrollbar: ColorAttr{colUndefined, AttrUndefined},
		PreviewLabel:     ColorAttr{colUndefined, AttrUndefined},
		ListLabel:        ColorAttr{colUndefined, AttrUndefined},
		ListBorder:       ColorAttr{colUndefined, AttrUndefined},
		Separator:        ColorAttr{colUndefined, AttrUndefined},
		Scrollbar:        ColorAttr{colUndefined, AttrUndefined},
	}
}

func InitTheme(theme *ColorTheme, baseTheme *ColorTheme, forceBlack bool) {
	if forceBlack {
		theme.Bg = ColorAttr{colBlack, AttrUndefined}
	}

	o := func(a ColorAttr, b ColorAttr) ColorAttr {
		c := a
		if b.Color != colUndefined {
			c.Color = b.Color
		}
		if b.Attr != AttrUndefined {
			c.Attr = b.Attr
		}
		return c
	}
	theme.Input = o(baseTheme.Input, theme.Input)
	theme.Fg = o(baseTheme.Fg, theme.Fg)
	theme.Bg = o(baseTheme.Bg, theme.Bg)
	theme.DarkBg = o(baseTheme.DarkBg, theme.DarkBg)
	theme.Prompt = o(baseTheme.Prompt, theme.Prompt)
	theme.Match = o(baseTheme.Match, theme.Match)
	theme.Current = o(baseTheme.Current, theme.Current)
	theme.CurrentMatch = o(baseTheme.CurrentMatch, theme.CurrentMatch)
	theme.Spinner = o(baseTheme.Spinner, theme.Spinner)
	theme.Info = o(baseTheme.Info, theme.Info)
	theme.Cursor = o(baseTheme.Cursor, theme.Cursor)
	theme.Marker = o(baseTheme.Marker, theme.Marker)
	theme.Header = o(baseTheme.Header, theme.Header)
	theme.Border = o(baseTheme.Border, theme.Border)
	theme.BorderLabel = o(baseTheme.BorderLabel, theme.BorderLabel)

	// These colors are not defined in the base themes
	theme.ListFg = o(theme.Fg, theme.ListFg)
	theme.ListBg = o(theme.Bg, theme.ListBg)
	theme.SelectedFg = o(theme.ListFg, theme.SelectedFg)
	theme.SelectedBg = o(theme.ListBg, theme.SelectedBg)
	theme.SelectedMatch = o(theme.Match, theme.SelectedMatch)
	theme.Disabled = o(theme.Input, theme.Disabled)
	theme.Gutter = o(theme.DarkBg, theme.Gutter)
	theme.PreviewFg = o(theme.Fg, theme.PreviewFg)
	theme.PreviewBg = o(theme.Bg, theme.PreviewBg)
	theme.PreviewLabel = o(theme.BorderLabel, theme.PreviewLabel)
	theme.PreviewBorder = o(theme.Border, theme.PreviewBorder)
	theme.ListLabel = o(theme.BorderLabel, theme.ListLabel)
	theme.ListBorder = o(theme.Border, theme.ListBorder)
	theme.Separator = o(theme.ListBorder, theme.Separator)
	theme.Scrollbar = o(theme.ListBorder, theme.Scrollbar)
	theme.PreviewScrollbar = o(theme.PreviewBorder, theme.PreviewScrollbar)

	initPalette(theme)
}

func initPalette(theme *ColorTheme) {
	pair := func(fg, bg ColorAttr) ColorPair {
		if fg.Color == colDefault && (fg.Attr&Reverse) > 0 {
			bg.Color = colDefault
		}
		return ColorPair{fg.Color, bg.Color, fg.Attr}
	}
	blank := theme.ListFg
	blank.Attr = AttrRegular

	ColPrompt = pair(theme.Prompt, theme.ListBg)
	ColNormal = pair(theme.ListFg, theme.ListBg)
	ColSelected = pair(theme.SelectedFg, theme.SelectedBg)
	ColInput = pair(theme.Input, theme.ListBg)
	ColDisabled = pair(theme.Disabled, theme.ListBg)
	ColMatch = pair(theme.Match, theme.ListBg)
	ColSelectedMatch = pair(theme.SelectedMatch, theme.SelectedBg)
	ColCursor = pair(theme.Cursor, theme.Gutter)
	ColCursorEmpty = pair(blank, theme.Gutter)
	if theme.SelectedBg.Color != theme.ListBg.Color {
		ColMarker = pair(theme.Marker, theme.SelectedBg)
	} else {
		ColMarker = pair(theme.Marker, theme.Gutter)
	}
	ColCurrent = pair(theme.Current, theme.DarkBg)
	ColCurrentMatch = pair(theme.CurrentMatch, theme.DarkBg)
	ColCurrentCursor = pair(theme.Cursor, theme.DarkBg)
	ColCurrentCursorEmpty = pair(blank, theme.DarkBg)
	ColCurrentMarker = pair(theme.Marker, theme.DarkBg)
	ColCurrentSelectedEmpty = pair(blank, theme.DarkBg)
	ColSpinner = pair(theme.Spinner, theme.ListBg)
	ColInfo = pair(theme.Info, theme.ListBg)
	ColHeader = pair(theme.Header, theme.ListBg)
	ColSeparator = pair(theme.Separator, theme.ListBg)
	ColScrollbar = pair(theme.Scrollbar, theme.ListBg)
	ColBorder = pair(theme.Border, theme.Bg)
	ColBorderLabel = pair(theme.BorderLabel, theme.Bg)
	ColPreviewLabel = pair(theme.PreviewLabel, theme.PreviewBg)
	ColPreview = pair(theme.PreviewFg, theme.PreviewBg)
	ColPreviewBorder = pair(theme.PreviewBorder, theme.PreviewBg)
	ColPreviewScrollbar = pair(theme.PreviewScrollbar, theme.PreviewBg)
	ColPreviewSpinner = pair(theme.Spinner, theme.PreviewBg)
	ColListLabel = pair(theme.ListLabel, theme.ListBg)
	ColListBorder = pair(theme.ListBorder, theme.ListBg)
}

func runeWidth(r rune) int {
	return uniseg.StringWidth(string(r))
}
