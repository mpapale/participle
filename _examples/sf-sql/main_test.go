package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShow(t *testing.T) {
	stmt := &Statement{}
	err := parser.ParseString("", `SHOW DATABASES`, stmt)
	require.NoError(t, err)
	require.Equal(t, *stmt.Show.ObjectTypePlural, ObjectTypePluralDatabases)

	stmt = &Statement{}
	err = parser.ParseString("", `SHOW SCHEMAS`, stmt)
	require.NoError(t, err)
	require.Equal(t, *stmt.Show.ObjectTypePlural, ObjectTypePluralSchemas)

	stmt = &Statement{}
	err = parser.ParseString("",
		`SHOW TERSE DATABASES`, stmt)
	require.NoError(t, err)
	require.Equal(t, *stmt.Show.ObjectTypePlural, ObjectTypePluralDatabases)
	require.True(t, stmt.Show.Terse)

	stmt = &Statement{}
	err = parser.ParseString("",
		`SHOW DATABASES LIKE '%foo%' LIMIT 10 FROM 'foo2'`, stmt)
	require.NoError(t, err)
	require.Equal(t, *stmt.Show.ObjectTypePlural, ObjectTypePluralDatabases)
	require.Equal(t, *stmt.Show.Like, "%foo%")
	require.Equal(t, *stmt.Show.Limit.Rows, 10)
	require.Equal(t, *stmt.Show.Limit.From, "foo2")
}

func TestShowErrors(t *testing.T) {
	stmt := &Statement{}
	err := parser.ParseString("",
		`SHOW DATABASES LIKE "%foo%"`, stmt)
	require.Error(t, err, "Snowflake does not allow double quoted LIKE string values")
}

func TestDescribe(t *testing.T) {
	stmt := &Statement{}
	err := parser.ParseString("", `DESCRIBE TABLE foo`, stmt)
	require.NoError(t, err)
	require.Equal(t, stmt.Describe.Name, "foo")

	stmt = &Statement{}
	err = parser.ParseString("", `DESC TABLE foo`, stmt)
	require.NoError(t, err)
	require.Equal(t, stmt.Describe.Name, "foo")

	stmt = &Statement{}
	err = parser.ParseString("", `DESCRIBE VIEW foo`, stmt)
	require.NoError(t, err)
	require.Equal(t, stmt.Describe.Name, "foo")

	stmt = &Statement{}
	err = parser.ParseString("", `DESC VIEW foo`, stmt)
	require.NoError(t, err)
	require.Equal(t, stmt.Describe.Name, "foo")
}
