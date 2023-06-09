package main

import (
	"testing"
	"time"

	"snippetbox.jaswanthp.com/internal/assert"
)

func TestHumanDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	for _, tc := range tests {

		//t.Run creates a subtest for every test case in the table
		t.Run(tc.name, func(t *testing.T) {
			hd := humanDate(tc.tm)
			assert.Equal(t, hd, tc.want)
		})

	}
}
