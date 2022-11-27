// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package preview

import "testing"

func TestVisibleLines(t *testing.T) {
	var m = Model{
		height:  10,
		yoffset: 0,
	}

	lines := m.visibleLines()
	if len(lines) != 0 {
		t.Errorf("len(line) = %v\n", len(lines))
	}
}
