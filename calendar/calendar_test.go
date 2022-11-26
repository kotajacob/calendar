package calendar

import (
	"reflect"
	"testing"
	"time"

	"git.sr.ht/~kota/calendar/config"
	"git.sr.ht/~kota/calendar/month"
)

func TestResize(t *testing.T) {
	now := time.Now()
	type test struct {
		count      int
		wantCount  int
		height     int
		today      time.Time
		selected   time.Time
		months     []month.Month
		wantMonths []month.Month
	}

	defaultConfig := config.Default()

	tests := []test{
		{
			// Basic resize from 1 to 3.
			count:     1,
			wantCount: 3,
			height:    3*month.MonthHeight + 1,
			today:     now,
			selected:  now,
			months: []month.Month{
				month.New(
					now,
					now,
					now,
					true,
					defaultConfig,
				),
			},
			wantMonths: []month.Month{
				month.New(
					lastMonth(now),
					now,
					now,
					true,
					defaultConfig,
				),
				month.New(
					now,
					now,
					now,
					true,
					defaultConfig,
				),
				month.New(
					nextMonth(now),
					now,
					now,
					true,
					defaultConfig,
				),
			},
		},
		{
			// Not quite enough height to resize.
			count:     1,
			wantCount: 1,
			height:    3 * month.MonthHeight,
			today:     now,
			selected:  now,
			months: []month.Month{
				month.New(
					now,
					now,
					now,
					true,
					defaultConfig,
				),
			},
			wantMonths: []month.Month{
				month.New(
					now,
					now,
					now,
					true,
					defaultConfig,
				),
			},
		},
		{
			// Basic resize from 3 to 1.
			count:     3,
			wantCount: 1,
			height:    1 * month.MonthHeight,
			today:     now,
			selected:  now,
			months: []month.Month{
				month.New(
					now,
					now,
					now,
					true,
					defaultConfig,
				),
				month.New(
					now,
					now,
					now,
					true,
					defaultConfig,
				),
				month.New(
					now,
					now,
					now,
					true,
					defaultConfig,
				),
			},
			wantMonths: []month.Month{
				month.New(
					now,
					now,
					now,
					true,
					defaultConfig,
				),
			},
		},
		{
			// Resize from 3 to 1, last month selected.
			count:     3,
			wantCount: 1,
			height:    1 * month.MonthHeight,
			today:     now,
			selected:  nextMonth(now),
			months: []month.Month{
				month.New(
					now,
					now,
					nextMonth(now),
					true,
					defaultConfig,
				),
				month.New(
					now,
					now,
					nextMonth(now),
					true,
					defaultConfig,
				),
				month.New(
					now,
					now,
					nextMonth(now),
					true,
					defaultConfig,
				),
			},
			wantMonths: []month.Month{
				month.New(
					nextMonth(now),
					now,
					nextMonth(now),
					true,
					defaultConfig,
				),
			},
		},
	}

	for _, tc := range tests {
		m := New(defaultConfig)
		m.today = tc.today
		m.selected = tc.selected
		m.months = tc.months
		m.height = tc.height
		got := m.resize()
		if len(got.months) != tc.wantCount {
			t.Fatalf(
				"wanted %v months, got %v months",
				tc.wantCount, len(got.months),
			)
		}
		for i, month := range got.months {
			if !reflect.DeepEqual(month, tc.wantMonths[i]) {
				t.Errorf(
					"month mismatch: error on month %v:\nwanted: %v\ngot: %v",
					i, tc.wantMonths[i], month,
				)
			}
		}
	}
}
