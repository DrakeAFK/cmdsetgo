package pick

import (
	"reflect"
	"testing"
)

func TestParseSelection(t *testing.T) {
	tests := []struct {
		input    string
		maxIndex int
		want     []int
		wantErr  bool
	}{
		{"1 3 5", 10, []int{1, 3, 5}, false},
		{"1-3 5", 10, []int{1, 2, 3, 5}, false},
		{"all", 5, []int{1, 2, 3, 4, 5}, false},
		{"1,2,3", 10, []int{1, 2, 3}, false},
		{"1-3 2-4", 10, []int{1, 2, 3, 4}, false}, // Deduplication
		{"10", 5, nil, true},                      // Out of range
		{"abc", 10, nil, true},                    // Invalid index
		{"1-", 10, nil, true},                     // Invalid range
		{"5-1", 10, nil, true},                    // Reversed range
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseSelection(tt.input, tt.maxIndex)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSelection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSelection() = %v, want %v", got, tt.want)
			}
		})
	}
}
