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

const (
	BorderThickness = 2
	MaxHeight       = month.MonthHeight*3 - BorderThickness/2
)

// Preview is the Bubble Tea model for this preview element.
type Preview struct {
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
func New(content string, conf *config.Config) Preview {
	return Preview{
		style: lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder(), true).
			PaddingLeft(conf.PreviewPadding).
			PaddingRight(conf.PreviewPadding).
			MarginLeft(conf.PreviewLeftMargin),
		config: conf,
	}.SetContent(content)
}

// Init the preview in Bubble Tea.
func (p Preview) Init() tea.Cmd {
	return nil
}

// Updates the preview in the Bubble Tea update loop.
func (p Preview) Update(msg tea.Msg) (Preview, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !p.isFocused {
			return p, nil
		}
		switch {
		case p.config.KeySelectUp.Contains(msg.String()):
			p.LineUp(1)
		case p.config.KeySelectDown.Contains(msg.String()):
			p.LineDown(1)
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if p.isFocused {
				p.LineUp(1)
			}
		case tea.MouseWheelDown:
			if p.isFocused {
				p.LineDown(1)
			}
		}
	case tea.WindowSizeMsg:
		p.setWidth(msg.Width)
		p.setHeight(msg.Height)
	}
	return p, nil
}

// SetWidth of the preview window.
// Padding, Margin, MinWidth, and MaxWidth are all taken into account.
func (p *Preview) setWidth(width int) {
	width = width - month.MonthWidth
	width = width - p.config.PreviewLeftMargin
	width = width - 2*p.config.PreviewPadding
	width = width - p.config.LeftPadding
	width = width - p.config.RightPadding
	width = width - BorderThickness

	// Handle min and max width.
	if width > p.config.PreviewMaxWidth {
		width = p.config.PreviewMaxWidth
	} else if width < p.config.PreviewMinWidth {
		width = 0
	}

	p.width = width
	p.lines = lines(p.content, p.width)
}

// setHeight of the preview window.
// Padding and MaxHeight are both taken into account.
func (p *Preview) setHeight(height int) {
	height = height - BorderThickness
	if height > MaxHeight {
		p.height = MaxHeight
		return
	}
	p.height = height
}

// SetContent is used to change the content displayed in the preview window.
func (p Preview) SetContent(s string) Preview {
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

	p.content = b.String()
	p.yoffset = 0
	p.lines = lines(p.content, p.width)
	return p
}

// Focus the preview.
func (p *Preview) Focus() {
	p.style.Border(lipgloss.RoundedBorder(), true)
	p.isFocused = true
}

// Unfocus the preview.
func (p *Preview) Unfocus() {
	p.style.Border(lipgloss.HiddenBorder(), true)
	p.isFocused = false
}

// AtTop returns whether or not the viewport is in the very top position.
func (p Preview) AtTop() bool {
	return p.yoffset <= 0
}

// AtBottom returns whether or not the viewport is at or past the very bottom
// position.
func (p Preview) AtBottom() bool {
	return p.yoffset >= p.maxYOffset()
}

// maxYOffset returns the maximum possible value of the y-offset based on the
// viewport's content and set height.
func (p Preview) maxYOffset() int {
	return max(0, len(p.lines)-p.height)
}

// LineUp moves the view down by the given number of lines.
func (p *Preview) LineUp(n int) {
	if p.AtTop() || n == 0 {
		return
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we are from the top.
	p.SetYOffset(p.yoffset - n)
	return
}

// LineDown moves the view down by the given number of lines.
func (p *Preview) LineDown(n int) {
	if p.AtBottom() || n == 0 {
		return
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we actually have left before we reach
	// the bottom.
	p.SetYOffset(p.yoffset + n)
	return
}

// SetYOffset sets the Y offset.
func (p *Preview) SetYOffset(n int) {
	p.yoffset = clamp(n, 0, p.maxYOffset())
}

// View renders the preview in its current state.
func (p Preview) View() string {
	if p.width < p.config.PreviewMinWidth {
		return ""
	}

	// Show empty preview window even when there's no content.
	visible := p.visibleLines()

	// Fill empty space with newlines
	extraLines := ""
	if len(visible) < p.height {
		extraLines = strings.Repeat("\n", max(0, p.height-len(visible)))
	}

	return p.style.Render(
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
func (p Preview) visibleLines() (lines []string) {
	if len(p.lines) > 0 {
		top := max(0, p.yoffset)
		bottom := clamp(p.yoffset+p.height, top, len(p.lines))
		lines = p.lines[top:bottom]
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
