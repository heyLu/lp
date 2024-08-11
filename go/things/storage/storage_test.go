package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ourEpoch = time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC)

func TestQueryV2(t *testing.T) {
	st, err := NewDBStorage(context.Background(), ":memory:")
	require.NoError(t, err)

	res, err := st.(*dbStorage).db.Exec(`INSERT INTO things_v2 (namespace, kind, id, summary, content, ref, number, float, bool, time, fields_json, tags, date_created, date_modified) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"test", "test_thing", 1, "this is a summary", nil, nil, nil, nil, nil, nil, nil, "#test,#hello-world", ourEpoch.Unix(), ourEpoch.Unix(),
	)
	require.NoError(t, err)

	n, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), n)

	rows, err := st.QueryV2(context.Background(), "test")
	require.NoError(t, err)

	numResults := 0
	for rows.Next() {
		numResults += 1

		var row Row
		err := rows.ScanV2(&row)
		assert.NoError(t, err)
		if err != nil {
			continue
		}

		assert.Equal(t,
			Row{
				Metadata: Metadata{
					Namespace:    "test",
					Kind:         "test_thing",
					ID:           1,
					DateCreated:  ourEpoch,
					DateModified: ourEpoch,
					Tags:         []string{"#test", "#hello-world"},
				},
				Summary: "this is a summary",
			},
			row)
	}

	assert.Equal(t, 1, numResults)
}

func TestInsertV2(t *testing.T) {
	st, err := NewDBStorage(context.Background(), ":memory:")
	require.NoError(t, err)

	expectedRow := Row{
		Metadata: Metadata{
			Namespace:    "test",
			Kind:         "test_thing",
			ID:           1,
			DateCreated:  ourEpoch,
			DateModified: time.Now().Round(time.Second).UTC(),
			Tags:         []string{"#test", "#hello-world"},
		},
		Summary: "this is a summary",
	}

	err = st.InsertV2(context.Background(), &expectedRow)
	require.NoError(t, err)

	rows, err := st.QueryV2(context.Background(), "test")
	require.NoError(t, err)

	require.True(t, rows.Next())
	var actualRow Row

	err = rows.ScanV2(&actualRow)
	require.NoError(t, err)
	require.Equal(t, expectedRow, actualRow)

	require.False(t, rows.Next())
}
