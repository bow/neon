package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/bow/courier/internal/store/opml"
)

func (s *SQLite) ImportOPML(ctx context.Context, payload []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("SQLite.ImportOPML")

	doc, err := opml.Parse(payload)
	if err != nil {
		return 0, fail(err)
	}

	if doc.Empty() {
		return 0, nil
	}

	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		now := time.Now()

		for _, outl := range doc.Body.Outlines {
			feedDBID, ierr := upsertFeed(
				ctx,
				tx,
				outl.XMLURL,
				nullIfTextEmpty(outl.Text),
				outl.Description,
				outl.HTMLURL,
				nil, // TODO: Set isStarred.
				nil,
				&now,
			)
			if ierr != nil {
				return ierr
			}

			if ierr = addFeedTags(ctx, tx, feedDBID, outl.Categories); ierr != nil {
				return ierr
			}
		}

		return nil
	}

	err = s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return 0, fail(err)
	}
	return doc.Length(), nil
}
