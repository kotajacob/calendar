// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package preview

import (
	"bytes"
	"strings"

	"git.sr.ht/~kota/calendar/config"
	"git.sr.ht/~kota/calendar/month"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

const BorderThickness = 2
const MaxHeight = month.MonthHeight*3 - BorderThickness/2

// Model is the Bubble Tea model for this preview element.
type Model struct {
	config    *config.Config
	style     lipgloss.Style
	content   string
	lines     []string
	height    int
	yoffset   int
	width     int
	isFocused bool
}

// New creates a new preview model.
func New(content string, width, height int, conf *config.Config) Model {
	return Model{
		style: lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder(), true).
			PaddingLeft(conf.PreviewPadding).
			PaddingRight(conf.PreviewPadding).
			MarginLeft(conf.PreviewLeftMargin),
		config: conf,
	}.SetContent(content)
}

// Init the preview in Bubble Tea.
func (m Model) Init() tea.Cmd {
	return nil
}

// Updates the preview in the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.isFocused {
			return m, nil
		}
		switch {
		case m.config.KeySelectDown.Contains(msg.String()):
			m.LineDown(1)
		case m.config.KeySelectUp.Contains(msg.String()):
			m.LineUp(1)
		}
	case tea.WindowSizeMsg:
		m.setWidth(msg.Width)
		m.setHeight(msg.Height)
	}
	return m, nil
}

// SetWidth of the preview window.
// Padding, Margin, MinWidth, and MaxWidth are all taken into account.
func (m *Model) setWidth(width int) {
	width = width - month.MonthWidth
	width = width - m.config.PreviewLeftMargin
	width = width - 2*m.config.PreviewPadding
	width = width - m.config.LeftPadding
	width = width - m.config.RightPadding
	width = width - BorderThickness

	// Handle min and max width.
	if width > m.config.PreviewMaxWidth {
		width = m.config.PreviewMaxWidth
	} else if width < m.config.PreviewMinWidth {
		width = 0
	}

	m.width = width
	m.lines = lines(m.content, m.width)
}

// setHeight of the preview window.
// Padding and MaxHeight are both taken into account.
func (m *Model) setHeight(height int) {
	height = height - BorderThickness
	if height > MaxHeight {
		m.height = MaxHeight
		return
	}
	m.height = height
}

// SetContent is used to change the content displayed in the preview window.
func (m Model) SetContent(s string) Model {
	// Remove hard-line breaks so we can re-wrap to the current width later.
	var b bytes.Buffer
	var last rune
	for _, r := range s {
		if r == '\n' {
			if last != '\n' {
				b.WriteRune(' ')
			} else {
				b.WriteRune('\n')
				b.WriteRune('\n')
			}
		} else {
			b.WriteRune(r)
		}
		last = r
	}

	m.content = b.String()
	m.yoffset = 0
	m.lines = lines(m.content, m.width)
	return m
}

// Focus the preview.
func (m *Model) Focus() {
	m.style.Border(lipgloss.RoundedBorder(), true)
	m.isFocused = true
}

// Unfocus the preview.
func (m *Model) Unfocus() {
	m.style.Border(lipgloss.HiddenBorder(), true)
	m.isFocused = false
}

// AtTop returns whether or not the viewport is in the very top position.
func (m Model) AtTop() bool {
	return m.yoffset <= 0
}

// AtBottom returns whether or not the viewport is at or past the very bottom
// position.
func (m Model) AtBottom() bool {
	return m.yoffset >= m.maxYOffset()
}

// maxYOffset returns the maximum possible value of the y-offset based on the
// viewport's content and set height.
func (m Model) maxYOffset() int {
	return max(0, len(m.lines)-m.height)
}

// LineUp moves the view down by the given number of lines.
func (m *Model) LineUp(n int) {
	if m.AtTop() || n == 0 {
		return
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we are from the top.
	m.SetYOffset(m.yoffset - n)
	return
}

// LineDown moves the view down by the given number of lines.
func (m *Model) LineDown(n int) {
	if m.AtBottom() || n == 0 {
		return
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we actually have left before we reach
	// the bottom.
	m.SetYOffset(m.yoffset + n)
	return
}

// SetYOffset sets the Y offset.
func (m *Model) SetYOffset(n int) {
	m.yoffset = clamp(n, 0, m.maxYOffset())
}

// View renders the preview in its current state.
func (m Model) View() string {
	if m.width < m.config.PreviewMinWidth {
		return ""
	}

	// Show empty preview window even when there's no content.
	visible := m.visibleLines()

	// Fill empty space with newlines
	extraLines := ""
	if len(visible) < m.height {
		extraLines = strings.Repeat("\n", max(0, m.height-len(visible)))
	}

	return m.style.Render(
		strings.Join(visible, "\n") + extraLines,
	)
}

// lines splits a long string into multiple lines using a max width and some
// padding.
func lines(s string, width int) []string {
	if width == 0 {
		return nil
	}

	// Empty documents should be the same "width".
	if len(s) == 0 {
		return []string{strings.Repeat(" ", width)}
	}

	s = wordwrap.String(s, width-2)
	s = wrap.String(s, width)
	s = padding.String(s, uint(width))
	return strings.Split(s, "\n")
}

// visibleLines returns the lines that should currently be visible in the
// viewport.
func (m Model) visibleLines() (lines []string) {
	if len(m.lines) > 0 {
		top := max(0, m.yoffset)
		bottom := clamp(m.yoffset+m.height, top, len(m.lines))
		lines = m.lines[top:bottom]
	}
	return lines
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
