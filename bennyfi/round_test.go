package bennyfi

import (
	"reflect"
	"testing"
)

func TestSplitEntriesIntoBatchSizeEntryIds(t *testing.T) {
	tests := []struct {
		name      string
		entries   []*Entry
		batchSize int
		expected  [][]uint64
	}{
		{
			name:      "Nil entries",
			entries:   nil,
			batchSize: 2,
			expected:  [][]uint64{},
		},
		{
			name:      "Empty entries",
			entries:   []*Entry{},
			batchSize: 2,
			expected:  [][]uint64{},
		},
		{
			name: "Less than batch size",
			entries: []*Entry{
				{EntryID: 1},
				{EntryID: 2},
			},
			batchSize: 3,
			expected:  [][]uint64{{1, 2}},
		},
		{
			name: "Exactly batch size",
			entries: []*Entry{
				{EntryID: 1},
				{EntryID: 2},
				{EntryID: 3},
			},
			batchSize: 3,
			expected:  [][]uint64{{1, 2, 3}},
		},
		{
			name: "Multiple batches",
			entries: []*Entry{
				{EntryID: 1},
				{EntryID: 2},
				{EntryID: 3},
				{EntryID: 4},
			},
			batchSize: 2,
			expected:  [][]uint64{{1, 2}, {3, 4}},
		},
		{
			name: "Multiple batches with remainder",
			entries: []*Entry{
				{EntryID: 1},
				{EntryID: 2},
				{EntryID: 3},
				{EntryID: 4},
				{EntryID: 5},
			},
			batchSize: 2,
			expected:  [][]uint64{{1, 2}, {3, 4}, {5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitEntriesIntoBatchSizeEntryIds(tt.entries, tt.batchSize)
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
