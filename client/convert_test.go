package client

import "testing"

func TestSelectImageURL(t *testing.T) {
	tests := []struct {
		name     string
		imageSet *ImageSet
		minWidth int
		want     string
	}{
		{
			name:     "nil image set",
			imageSet: nil,
			minWidth: 1024,
			want:     "",
		},
		{
			name:     "empty variants",
			imageSet: &ImageSet{Variants: []ImageVariant{}},
			minWidth: 1024,
			want:     "",
		},
		{
			name: "single image above threshold",
			imageSet: &ImageSet{Variants: []ImageVariant{
				{URL: "https://example.com/2048.jpg", Width: 2048},
			}},
			minWidth: 1024,
			want:     "https://example.com/2048.jpg",
		},
		{
			name: "single image below threshold",
			imageSet: &ImageSet{Variants: []ImageVariant{
				{URL: "https://example.com/512.jpg", Width: 512},
			}},
			minWidth: 1024,
			want:     "https://example.com/512.jpg",
		},
		{
			name: "picks smallest above 1024",
			imageSet: &ImageSet{Variants: []ImageVariant{
				{URL: "https://example.com/2048.jpg", Width: 2048},
				{URL: "https://example.com/512.jpg", Width: 512},
				{URL: "https://example.com/1024.jpg", Width: 1024},
				{URL: "https://example.com/1400.jpg", Width: 1400},
			}},
			minWidth: 1024,
			want:     "https://example.com/1024.jpg",
		},
		{
			name: "all below threshold picks largest",
			imageSet: &ImageSet{Variants: []ImageVariant{
				{URL: "https://example.com/200.jpg", Width: 200},
				{URL: "https://example.com/800.jpg", Width: 800},
				{URL: "https://example.com/500.jpg", Width: 500},
			}},
			minWidth: 1024,
			want:     "https://example.com/800.jpg",
		},
		{
			name: "custom threshold selects accordingly",
			imageSet: &ImageSet{Variants: []ImageVariant{
				{URL: "https://example.com/200.jpg", Width: 200},
				{URL: "https://example.com/800.jpg", Width: 800},
				{URL: "https://example.com/500.jpg", Width: 500},
			}},
			minWidth: 500,
			want:     "https://example.com/500.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := selectImageURL(tt.imageSet, tt.minWidth)
			if got != tt.want {
				t.Errorf("selectImageURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
