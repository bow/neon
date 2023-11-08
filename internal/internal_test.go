// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []ID
		want  []ID
	}{
		{name: "empty input", input: []ID{}, want: []ID{}},
		{
			"one item",
			[]ID{3},
			[]ID{3},
		},
		{
			"multiple items, no duplicates",
			[]ID{2, 7, 3},
			[]ID{2, 7, 3},
		},
		{
			"multiple items, duplicates at start",
			[]ID{2, 2, 7, 3},
			[]ID{2, 7, 3},
		},
		{
			"multiple items, duplicates in the middle",
			[]ID{2, 7, 7, 3},
			[]ID{2, 7, 3},
		},
		{
			"multiple items, duplicates at end",
			[]ID{2, 7, 3, 3, 3, 3},
			[]ID{2, 7, 3},
		},
		{
			"multiple items, duplicates in several places",
			[]ID{1, 2, 5, 5, 7, 3, 3, 3, 3},
			[]ID{1, 2, 5, 7, 3},
		},
		{
			"multiple items, duplicates across several places",
			[]ID{1, 2, 5, 5, 7, 3, 3, 2, 2, 3, 3},
			[]ID{1, 2, 5, 7, 3},
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
