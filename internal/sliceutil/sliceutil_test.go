// Copyright (c) 2023-2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []uint32
		want  []uint32
	}{
		{name: "empty input", input: []uint32{}, want: []uint32{}},
		{
			"one item",
			[]uint32{3},
			[]uint32{3},
		},
		{
			"multiple items, no duplicates",
			[]uint32{2, 7, 3},
			[]uint32{2, 7, 3},
		},
		{
			"multiple items, duplicates at start",
			[]uint32{2, 2, 7, 3},
			[]uint32{2, 7, 3},
		},
		{
			"multiple items, duplicates in the middle",
			[]uint32{2, 7, 7, 3},
			[]uint32{2, 7, 3},
		},
		{
			"multiple items, duplicates at end",
			[]uint32{2, 7, 3, 3, 3, 3},
			[]uint32{2, 7, 3},
		},
		{
			"multiple items, duplicates in several places",
			[]uint32{1, 2, 5, 5, 7, 3, 3, 3, 3},
			[]uint32{1, 2, 5, 7, 3},
		},
		{
			"multiple items, duplicates across several places",
			[]uint32{1, 2, 5, 5, 7, 3, 3, 2, 2, 3, 3},
			[]uint32{1, 2, 5, 7, 3},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()
				got := Dedup(tt.input)
				assert.Equal(t, tt.want, got)
			},
		)
	}

}
