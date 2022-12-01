// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package help

import (
	tea "github.com/charmbracelet/bubbletea"
)

const Content = `

Select             = hjkl, arrow keys, or mouse
Focus preview      = tab                       
Toggle preview     = p                         
Scroll preview     = jk, up/down (if focused)  
Edit note          = enter                     
Copy date          = y                         
Goto last Sunday   = b, H                      
Goto next Sunday   = w                         
Goto next Saturday = e, L                      
Go up one month    = ctrl+u                    
Go down one month  = ctrl+d                    
`

// Help is the Bubble Tea model for this help element.
type Help struct {
	version string
}

// New creates a new help model.
func New(version string) Help {
	return Help{version: version}
}

// Init the help window in Bubble Tea.
func (h Help) Init() tea.Cmd {
	return nil
}

// Updates the help window in the Bubble Tea update loop.
func (h Help) Update(msg tea.Msg) (Help, tea.Cmd) {
	return h, nil
}

// View renders the preview in its current state.
func (h Help) View() string {
	return "Calendar " + h.version + Content
}
