// Copyright (c) 2023-2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package sliceutil

import "sort"

func Dedup[T comparable](values []T) []T {
	seen := make(map[T]struct{})
	nodup := make([]T, 0)

	for _, val := range values {
		if _, exists := seen[val]; exists {
			continue
		}
		seen[val] = struct{}{}
		nodup = append(nodup, val)
	}

	return nodup
}

type SortFunc[T any] func(v1, v2 T) int

type Sorter[T any] struct {
	items  []T
	compfs []SortFunc[T]
}

func Ordered[T any]() *Sorter[T] {
	return &Sorter[T]{compfs: make([]SortFunc[T], 0)}
}

func (s *Sorter[T]) By(compf ...SortFunc[T]) *Sorter[T] {
	s.compfs = append(s.compfs, compf...)
	return s
}

func (s *Sorter[T]) Len() int {
	return len(s.items)
}

func (s *Sorter[T]) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s *Sorter[T]) Less(i, j int) bool {
	p, q := s.items[i], s.items[j]
	var k int
	for k = 0; k < len(s.compfs)-1; k++ {
		comp := s.compfs[k](p, q)
		if comp < 0 {
			return true
		}
		if comp > 0 {
			return false
		}
	}
	return s.compfs[k](p, q) < 0
}

func (s *Sorter[T]) Sort(items []T) {
	s.items = items
	sort.Sort(s)
}
