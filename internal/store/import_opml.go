// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bow/iris/internal/opml"
)

func (s *SQLite) ImportOPML(
	ctx context.Context,
	payload []byte,
) (processed int, imported int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(payload) == 0 {
		return 0, 0, ErrEmptyPayload
	}

	fail := failF("SQLite.ImportOPML")

	doc, err := opml.Parse(payload)
	if err != nil {
		return 0, 0, fail(err)
	}

	if doc.Empty() {
		return 0, 0, nil
	}

	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		now := time.Now()

		for _, outl := range doc.Body.Outlines {
			if outl.Text == "" {
				return fmt.Errorf("missing title for feed with URL=%s in OPML document", outl.XMLURL)
			}
			feedID, isAdded, ierr := upsertFeed(
				ctx,
				tx,
				outl.XMLURL,
				pointerOrNil(outl.Text),
				outl.Description,
				outl.HTMLURL,
				outl.IsStarred,
				nil,
				&now,
			)
			if ierr != nil {
				return ierr
			}

			if ierr = addFeedTags(ctx, tx, feedID, outl.Categories); ierr != nil {
				return ierr
			}
			processed++
			if isAdded {
				imported++
			}
		}

		return nil
	}

	err = s.withTx(ctx, dbFunc)
	if err != nil {
		return 0, 0, fail(err)
	}
	return processed, imported, nil
}
