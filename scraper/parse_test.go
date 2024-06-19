package scraper

import (
	"testing"
	"time"
)

func Test_parseDuring(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    time.Duration
		wantErr bool
	}{
		{
			name: "many minutes",
			raw:  "88 minuter",
			want: time.Minute * 88,
		},
		{
			name: "below 1 minute",
			raw:  "0:55 minuter",
			want: time.Second * 55,
		},
		{
			name: "mix",
			raw:  "3:55 minuter",
			want: time.Minute*3 + time.Second*55,
		},
		{
			name:    "47 secunder",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDuring(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDuring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseDuring() = %v, want %v", got, tt.want)
			}
		})
	}
}
