package main

import (
	"testing"
	"time"
)

func Test_humanDate(t *testing.T) {
	tests := []struct {
		name string
		arg  time.Time
		want string
	}{
		{
			name: "UTC",
			arg:  time.Date(2021, 02, 18, 16, 0, 0, 0, time.UTC),
			want: "18 Feb 2021 at 16:00",
		},
		{
			name: "Empty",
			arg:  time.Time{},
			want: "",
		},
		{
			name: "CET",
			// CET (Central European Time) as the time zone, which is one hour ahead of UTC
			arg:  time.Date(2021, 02, 18, 16, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "18 Feb 2021 at 15:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := humanDate(tt.arg); got != tt.want {
				t.Errorf("humanDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
