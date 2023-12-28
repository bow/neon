// Copyright (c) 2022-2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package entity

import "fmt"

type FeedNotFoundError struct{ ID any }

func (e FeedNotFoundError) Error() string {
	return fmt.Sprintf("feed with ID=%v not found", e.ID)
}

type EntryNotFoundError struct{ ID any }

func (e EntryNotFoundError) Error() string {
	return fmt.Sprintf("entry with ID=%v not found", e.ID)
}
