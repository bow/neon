// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bow/neon/internal/entity"
)

func TestDedup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []entity.ID
		want  []entity.ID
	}{
		{name: "empty input", input: []entity.ID{}, want: []entity.ID{}},
		{
			"one item",
			[]entity.ID{3},
			[]entity.ID{3},
		},
		{
			"multiple items, no duplicates",
			[]entity.ID{2, 7, 3},
			[]entity.ID{2, 7, 3},
		},
		{
			"multiple items, duplicates at start",
			[]entity.ID{2, 2, 7, 3},
			[]entity.ID{2, 7, 3},
		},
		{
			"multiple items, duplicates in the middle",
			[]entity.ID{2, 7, 7, 3},
			[]entity.ID{2, 7, 3},
		},
		{
			"multiple items, duplicates at end",
			[]entity.ID{2, 7, 3, 3, 3, 3},
			[]entity.ID{2, 7, 3},
		},
		{
			"multiple items, duplicates in several places",
			[]entity.ID{1, 2, 5, 5, 7, 3, 3, 3, 3},
			[]entity.ID{1, 2, 5, 7, 3},
		},
		{
			"multiple items, duplicates across several places",
			[]entity.ID{1, 2, 5, 5, 7, 3, 3, 2, 2, 3, 3},
			[]entity.ID{1, 2, 5, 7, 3},
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
