package st

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOption_Scan(t *testing.T) {
	t.Run("scan int", func(t *testing.T) {
		o := Some(fake.Int())

		val := fake.Int()
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some(val), o)
	})

	t.Run("scan float", func(t *testing.T) {
		o := Some(fake.Float64(2, 0, 100000))

		val := fake.Float64(2, 0, 100000)
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some[float64](val), o)
	})

	t.Run("scan string", func(t *testing.T) {
		o := Some(fake.RandomStringWithLength(8))

		val := fake.RandomStringWithLength(8)
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some[string](val), o)
	})

	t.Run("scan bool", func(t *testing.T) {
		o := Some(fake.Bool())

		val := fake.Bool()
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some[bool](val), o)
	})

	t.Run("scan to pointer", func(t *testing.T) {
		o := Some(ptr(fake.Int()))

		val := fake.Int()
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some[*int](ptr(val)), o)
	})

	t.Run("type conversion", func(t *testing.T) {
		// We should be able to assign an int32 to an int64
		o := Some[int64](0)

		val := fake.Int32()
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some[int64](int64(val)), o)
	})

	t.Run("nil value scans to none", func(t *testing.T) {
		o := Some("")

		err := o.Scan(nil)
		require.NoError(t, err)

		assert.Equal(t, None[string](), o)
	})

	t.Run("incompatible types", func(t *testing.T) {
		o := Some(0)

		err := o.Scan(fake.RandomStringWithLength(8))
		require.ErrorContains(t, err, "invalid type for *st.Option[int]: string")
	})

	t.Run("can't assign to nil ptr", func(t *testing.T) {
		o := Some[*string](nil)

		err := o.Scan(fake.RandomStringWithLength(8))
		require.ErrorContains(t, err, "can't assign value to nil pointer of type *string")
	})
}

func TestOption_Value(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		o := Some(val)

		res, err := o.Value()
		require.NoError(t, err)

		assert.Equal(t, val, res)
	})

	t.Run("none", func(t *testing.T) {
		o := None[int]()

		res, err := o.Value()
		require.NoError(t, err)

		assert.Nil(t, res)
	})
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func TestOption_SQLInsert(t *testing.T) {
	createTableStmt := "CREATE TABLE test (i INTEGER, s TEXT)"
	insertStmt := "INSERT INTO test VALUES (?, ?)"
	selectStmt := "SELECT i, s FROM test"

	tests := []struct {
		name      string
		i         Option[int64]
		s         Option[string]
		expectedI sql.NullInt64
		expectedS sql.NullString
	}{
		{
			name:      "none",
			i:         None[int64](),
			s:         None[string](),
			expectedI: sql.NullInt64{Valid: false, Int64: 0},
			expectedS: sql.NullString{Valid: false, String: ""},
		},
		{
			name:      "some with zero values",
			i:         Some[int64](0),
			s:         Some(""),
			expectedI: sql.NullInt64{Valid: true, Int64: 0},
			expectedS: sql.NullString{Valid: true, String: ""},
		},
		{
			name:      "some with values",
			i:         Some[int64](42),
			s:         Some("forty two"),
			expectedI: sql.NullInt64{Valid: true, Int64: 42},
			expectedS: sql.NullString{Valid: true, String: "forty two"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			_, err := db.ExecContext(t.Context(), createTableStmt)
			require.NoError(t, err)

			stmt, err := db.PrepareContext(t.Context(), insertStmt)
			require.NoError(t, err)

			defer stmt.Close()

			// Save the values
			res, err := stmt.ExecContext(t.Context(), tt.i, tt.s)
			require.NoError(t, err)

			rowsCount, err := res.RowsAffected()
			require.NoError(t, err)

			assert.EqualValues(t, 1, rowsCount)

			// Check what was inserted
			var i sql.NullInt64
			var s sql.NullString

			rows, err := db.QueryContext(t.Context(), selectStmt)
			require.NoError(t, err)
			require.NoError(t, rows.Err())

			defer rows.Close()

			require.True(t, rows.Next())
			err = rows.Scan(&i, &s)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedI, i)
			assert.Equal(t, tt.expectedS, s)
		})
	}
}

func TestOption_SQLFetch(t *testing.T) {
	createTableStmt := "CREATE TABLE test (i INTEGER, s TEXT)"
	insertStmt := "INSERT INTO test VALUES (?, ?)"
	selectStmt := "SELECT i, s FROM test"

	tests := []struct {
		name      string
		i         sql.NullInt64
		s         sql.NullString
		expectedI Option[int64]
		expectedS Option[string]
	}{
		{
			name:      "none",
			i:         sql.NullInt64{Valid: false, Int64: 0},
			s:         sql.NullString{Valid: false, String: ""},
			expectedI: None[int64](),
			expectedS: None[string](),
		},
		{
			name:      "some with zero values",
			i:         sql.NullInt64{Valid: true, Int64: 0},
			s:         sql.NullString{Valid: true, String: ""},
			expectedI: Some[int64](0),
			expectedS: Some(""),
		},
		{
			name:      "some with values",
			i:         sql.NullInt64{Valid: true, Int64: 42},
			s:         sql.NullString{Valid: true, String: "forty two"},
			expectedI: Some[int64](42),
			expectedS: Some("forty two"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			_, err := db.ExecContext(t.Context(), createTableStmt)
			require.NoError(t, err)

			stmt, err := db.PrepareContext(t.Context(), insertStmt)
			require.NoError(t, err)

			defer stmt.Close()

			// Save the values
			res, err := stmt.ExecContext(t.Context(), tt.i, tt.s)
			require.NoError(t, err)

			rowsCount, err := res.RowsAffected()
			require.NoError(t, err)

			assert.EqualValues(t, 1, rowsCount)

			// Check what was inserted
			var i Option[int64]
			var s Option[string]

			rows, err := db.QueryContext(t.Context(), selectStmt)
			require.NoError(t, err)
			require.NoError(t, rows.Err())

			defer rows.Close()

			require.True(t, rows.Next())
			err = rows.Scan(&i, &s)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedI, i)
			assert.Equal(t, tt.expectedS, s)
		})
	}
}

func TestGorm(t *testing.T) {
	type TestModel struct {
		ID uuid.UUID `gorm:"primaryKey;type:uuid"`
		I  Option[int]
		S  Option[string]
	}

	db, err := gorm.Open(sqlite.Open(":memory:"))
	require.NoError(t, err)

	err = db.AutoMigrate(&TestModel{})
	require.NoError(t, err)

	data := TestModel{ID: uuid.New(), I: Some[int](42), S: Some("forty two")}

	err = db.Create(&data).Error
	require.NoError(t, err)

	var dbData TestModel
	err = db.First(&dbData).Error
	require.NoError(t, err)

	assert.Equal(t, data, dbData)
}
