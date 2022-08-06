package main

import (
	"reflect"
	"testing"
	"time"
)

func TestMonth(t *testing.T) {
	type test struct {
		date time.Time
		year bool
		want string
	}

	tests := []test{
		{
			date: time.Unix(1659745680, 0),
			year: true,
			want: "    August 2022     \n",
		},
		{
			date: time.Unix(1659745680, 0),
			year: false,
			want: "       August       \n",
		},
		{
			date: time.Unix(1662337680, 0),
			year: true,
			want: "   September 2022   \n",
		},
		{
			date: time.Unix(15552000, 0),
			year: true,
			want: "     June 1970      \n",
		},
	}

	for _, tc := range tests {
		got := month(tc.date, tc.year)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("\n\texpected:\n%v\tgot:\n%v", tc.want, got)
		}
	}
}
