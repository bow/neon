package store

import (
	"context"
	"database/sql"
	"fmt"
)

// SetEntryFields updates fields of an entry.
func (s *SQLite) SetEntryFields(
	ctx context.Context,
	entryDBID DBID,
	isRead *bool,
) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("sqliteStore.SetEntryFields")

	var entry *Entry
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		var err error

		if isRead != nil {
			err = s.updateEntryIsRead(ctx, tx, entryDBID, *isRead)
			if err != nil {
				return fail(err)
			}
		}

		entry, err = s.getEntry(ctx, tx, entryDBID)
		if err != nil {
			return fail(err)
		}

		return nil
	}

	err := s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *SQLite) updateEntryIsRead(
	ctx context.Context,
	tx *sql.Tx,
	entryDBID DBID,
	isRead bool,
) error {

	sql1 := `UPDATE entries SET is_read = $2 WHERE id = $1 RETURNING id`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()

	var updatedID DBID
	err = stmt1.QueryRowContext(ctx, entryDBID, isRead).Scan(&updatedID)
	if err != nil {
		return err
	}
	if updatedID == 0 {
		// TODO: Wrap in proper gRPC errors.
		return fmt.Errorf("entry id %d does not exist", updatedID)
	}
	return nil
}

func (s *SQLite) getEntry(ctx context.Context, tx *sql.Tx, entryDBID DBID) (*Entry, error) {

	sql1 := `
		SELECT
			e.id AS id,
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
