package store

import (
	"context"
	"database/sql"
)

// EditEntries updates fields of an entry.
func (s *SQLite) EditEntries(
	ctx context.Context,
	ops []*EntryEditOp,
) ([]*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("SQLite.EditEntries")

	updateFunc := func(ctx context.Context, tx *sql.Tx, op *EntryEditOp) (*Entry, error) {
		if err := setEntryIsRead(ctx, tx, op.DBID, op.IsRead); err != nil {
			return nil, err
		}
		return getEntry(ctx, tx, op.DBID)
	}

	var entries = make([]*Entry, len(ops))
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		for i, op := range ops {
			entry, err := updateFunc(ctx, tx, op)
			if err != nil {
				return fail(err)
			}
			entries[i] = entry
		}
		return nil
	}

	err := s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func getEntry(ctx context.Context, tx *sql.Tx, entryDBID DBID) (*Entry, error) {

	sql1 := `
		SELECT
			e.id AS id,
			e.feed_id AS feed_id,
			e.title AS title,
			e.is_read AS is_read,
			e.external_id AS ext_id,
			e.description AS description,
			e.content AS content,
			e.url AS url,
			e.update_time AS update_time,
			e.publication_time AS publication_time
		FROM
			entries e
		WHERE
			e.id = $1
		ORDER BY
			COALESCE(e.update_time, e.publication_time) DESC
`
	scanRow := func(row *sql.Row) (*Entry, error) {
		var entry Entry
		if err := row.Scan(
			&entry.DBID,
			&entry.FeedDBID,
			&entry.Title,
			&entry.IsRead,
			&entry.ExtID,
			&entry.Description,
			&entry.Content,
			&entry.URL,
			&entry.Updated,
			&entry.Published,
		); err != nil {
			return nil, err
		}
		return &entry, nil
	}

	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return nil, err
	}
	defer stmt1.Close()

	return scanRow(stmt1.QueryRowContext(ctx, entryDBID))
}

var setEntryIsRead = setTableField[bool]("entries", "is_read")
