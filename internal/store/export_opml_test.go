package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportOPMLEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	payload, err := st.ExportOPML(context.Background())
	r.EqualError(err, "SQLite.ExportOPML: unimplemented")

	a.Nil(payload)
}
