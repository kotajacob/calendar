// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package date

import (
	"testing"
	"time"
)

func TestSameMonth(t *testing.T) {
	now := time.Now()
	type test struct {
		x    time.Time
		y    time.Time
		want bool
	}

	tests := []test{
		{
			x:    now,
			y:    now,
			want: true,
		},
		{
			x:    FirstDay(now),
			y:    FirstDay(now),
			want: true,
		},
		{
			x:    FirstDay(now),
			y:    LastDay(now),
			want: true,
		},
		{
			x: time.Date(
				now.Year(),
				now.Month(),
				0,
				0,
				0,
				0,
				0,
				now.Location(),
			),
			y:    FirstDay(now),
			want: false,
		},
	}

	for _, tc := range tests {
		got := SameMonth(tc.x, tc.y)
		if got != tc.want {
			t.Fatalf(
				"input:\n%v\n%v\nwant: %v, got: %v\n",
				tc.x,
				tc.y,
				tc.want,
				got,
			)
		}
	}
}
