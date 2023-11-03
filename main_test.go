// License: GPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseArgs(t *testing.T) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatalf("failed loading UTC timezone: %v", err)
	}
	now := time.Date(
		2006,
		time.Month(0o1),
		2,
		0, 0, 0, 0,
		utc,
	)

	type test struct {
		input []string
		want  time.Time
	}

	tests := []test{
		{
			input: []string{
				"calendar",
			},
			want: time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				0, 0, 0, 0,
				now.Location(),
			),
		},
		{
			input: []string{
				"calendar",
				"10",
			},
			want: time.Date(
				now.Year(),
				now.Month(),
				10,
				0, 0, 0, 0,
				now.Location(),
			),
		},
		{
			input: []string{
				"calendar",
				"2008-05-04",
			},
			want: time.Date(
				2008,
				time.Month(5),
				4,
				0, 0, 0, 0,
				now.Location(),
			),
		},
		{
			input: []string{
				"calendar",
				"6",
				"4",
			},
			want: time.Date(
				now.Year(),
				time.Month(4),
				6,
				0, 0, 0, 0,
				now.Location(),
			),
		},
		{
			input: []string{
				"calendar",
				"6",
				"4",
				"2104",
			},
			want: time.Date(
				2104,
				time.Month(4),
				6,
				0, 0, 0, 0,
				now.Location(),
			),
		},
		{
			input: []string{
				"calendar",
				"February",
			},
			want: time.Date(
				now.Year(),
				time.Month(2),
				now.Day(),
				0, 0, 0, 0,
				now.Location(),
			),
		},
		{
			input: []string{
				"calendar",
				"february",
			},
			want: time.Date(
				now.Year(),
				time.Month(2),
				now.Day(),
				0, 0, 0, 0,
				now.Location(),
			),
		},
		{
			input: []string{
				"calendar",
				"FEBRUARY",
			},
			want: time.Date(
				now.Year(),
				time.Month(2),
				now.Day(),
				0, 0, 0, 0,
				now.Location(),
			),
		},
	}

	for _, tc := range tests {
		got := parseArgs(tc.input, now)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
