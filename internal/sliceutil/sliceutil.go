// Copyright (c) 2023-2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package sliceutil

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
