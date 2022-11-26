package month

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
			x:    firstDay(now),
			y:    firstDay(now),
			want: true,
		},
		{
			x:    firstDay(now),
			y:    lastDay(now),
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
			y:    firstDay(now),
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
