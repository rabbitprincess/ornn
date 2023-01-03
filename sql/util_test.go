package sql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUtil_ClearStrInQuota(t *testing.T) {
	require.Equal(t, fmt.Sprintf("%s", strings.Repeat(" ", len("'will delete'"))), Util_ClearInQuot("'will delete'"))
	require.Equal(t, fmt.Sprintf("%s", strings.Repeat(" ", len(`"will delete"`))), Util_ClearInQuot(`"will delete"`))
	require.Equal(t, fmt.Sprintf("%s", strings.Repeat(" ", len("`will delete`"))), Util_ClearInQuot("`will delete`"))
}

func TestUtil_SplitByDelimiter(t *testing.T) {
	for _, test := range []struct {
		input        string
		expectFront  string
		expectBehind string
	}{
		{"'?' \"?\" `?` ?", "'?' \"?\" `?` ?", ""},
		{"??????????", "??????????", ""},
		{"?'?'?'?'?'?'?'?'?'?'", "?'?'?'?'?'?'?'?'?'?'", ""},
		{
			"SELECT * FROM table_name WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
			"SELECT * FROM table_name ",
			"WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
		},
		{
			"INSERT INTO table_name (field_name_1,field_name_2,field_name_3,field_name_4) VALUES (?,?,?,?);",
			"INSERT INTO table_name (field_name_1,field_name_2,field_name_3,field_name_4) VALUES (?,?,?,?);",
			"",
		},
		{
			"UPDATE table_name SET field_name_1=?, field_name_2=?, field_name_3=?, field_name_4=? WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
			"UPDATE table_name SET field_name_1=?, field_name_2=?, field_name_3=?, field_name_4=? ",
			"WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
		},
		{
			"DELETE FROM table_name WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
			"DELETE FROM table_name ",
			"WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
		},
		{
			"SELECT * FROM table_name WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
			"SELECT * FROM table_name ",
			"WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
		},
		{
			"INSERT INTO table_name (field_name_1,field_name_2,field_name_3,field_name_4) VALUES (?,'?', \"?\",`?`);",
			"INSERT INTO table_name (field_name_1,field_name_2,field_name_3,field_name_4) VALUES (?,'?', \"?\",`?`);",
			"",
		},
		{
			"UPDATE table_name SET field_name_1=?, field_name_2='?', field_name_3=\"?\", field_name_4=`?` WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
			"UPDATE table_name SET field_name_1=?, field_name_2='?', field_name_3=\"?\", field_name_4=`?` ",
			"WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
		},
		{
			"DELETE FROM table_name WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
			"DELETE FROM table_name ",
			"WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
		},
	} {
		retFront, retBehind := Util_SplitByDelimiter(test.input, "where")
		require.Equal(t, test.expectFront, retFront)
		require.Equal(t, test.expectBehind, retBehind)
	}
}

func TestUtil_ExportBetweenDelimiter(t *testing.T) {
	for _, test := range []struct {
		input   string
		expects []string
	}{
		{"%a%", []string{"a"}},
		{"%a%b", []string{"a"}},
		{"a%b%", []string{"b"}},
		{"a%b%c", []string{"b"}},
		{"%a%b%c%", []string{"a", "c"}},
		{"a%b%c%d%e", []string{"b", "d"}},
		{"%a%b%c%d%e%", []string{"a", "c", "e"}},
		{"%a%`%b%`%c%d%e%", []string{"a", "c", "e"}},
		{"%a`%b%`c%d%e%", []string{"a`%b%`c", "e"}},
		{"%a%'%b%c%d%'%e%", []string{"a", "e"}},
		{"%a'%b%c%d%'e%", []string{"a'%b%c%d%'e"}},
		{"'%a'%b%c%d%'e%'", []string{"b", "d"}},
		{"'%a%b%c%d%e%'", []string{}},
	} {
		rets, err := Util_ExportBetweenDelimiter(test.input, "%")
		require.NoError(t, err)
		require.Equal(t, test.expects, rets)
	}
}

func TestUtil_ReplaceBetweenDelimiter(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect string
	}{
		{"%a%", " "},
		{"%a%b", " b"},
		{"a%b%", "a "},
		{"a%b%c", "a c"},
		{"%a%b%c%", " b "},
		{"a%b%c%d%e", "a c e"},
		{"%a%b%c%d%e%", " b d "},
		{"%a%`%b%`%c%d%e%", " `%b%` d "},
		{"%a`%b%`c%d%e%", " d "},
		{"%a%'%b%c%d%'%e%", " '%b%c%d%' "},
		{"%a'%b%c%d%'e%", " "},
		{"'%a'%b%c%d%'e%'", "'%a' c 'e%'"},
		{"'%a%b%c%d%e%'", "'%a%b%c%d%e%'"},
	} {
		ret := Util_ReplaceBetweenDelimiter(test.input, "%", " ")
		require.Equal(t, ret, test.expect)
	}
}

func TestUtil_ClearDelimiter(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect string
	}{
		{"%a%", "a"},
		{"%a%b", "ab"},
		{"a%b%", "ab"},
		{"a%b%c", "abc"},
		{"%%%", ""},
		{"a%a'%a'a%a", "aa'%a'aa"},
	} {
		ret := Util_ClearDelimiter(test.input, "%")
		require.Equal(t, test.expect, ret)
	}
}

func TestUtil_ReplaceInDelimiter(t *testing.T) {
	for _, test := range []struct {
		input     string
		expect    string
		delimiter string
		spliter   string
	}{
		{"%a%", "a", "%", ""},
		{"%a%", "a", "%", "/"},
		{"%a%b%", "ab", "%", ""},
		{"%%%", "", "%", "/"},
		{"a%a'%a'a%a", "aa'%a'aa", "%", "/"},
		{"%a/b%", "a/b", "%", ""},
		{"%a/b%", "b", "%", "/"},
		{"%ab%", "ab", "%", "/"},
		{"%a/b%c/d%e/f%", "bc/df", "%", "/"},
	} {
		ret := Util_ReplaceInDelimiter(test.input, test.delimiter, test.spliter)
		require.Equal(t, test.expect, ret)
	}
}

func TestUtil_ExportInsertQueryValues(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect string
	}{
		{
			"INSERT INTO `test table` VALUES (?, ?, ?, ?)",
			"?, ?, ?, ?",
		},
		{
			"INSERT INTO `test table` VALUES (?, ?, ?, ?);",
			"?, ?, ?, ?",
		},
		{
			"INSERT INTO `test table` (`seq`, `str`, `num`, `dtn`) VALUES (?, ?, ?, ?);",
			"?, ?, ?, ?",
		},
		{
			"INSERT INTO `test table` (`seq`, `str`, `num`, `dtn`) VALUES (`1`, `안녕하세요`, `12345`, ?);",
			"`1`, `안녕하세요`, `12345`, ?",
		},
		{
			"INSERT INTO Customers (CustomerName, ContactName, Address, City, PostalCode, Country) VALUES (`SupplierName`, `ContactName`, ?, ?, ?, `Country`)",
			"`SupplierName`, `ContactName`, ?, ?, ?, `Country`",
		},
		{
			"INSERT INTO Customers (CustomerName, ContactName, Address, City, PostalCode, Country) SELECT SupplierName, ContactName, Address, City, PostalCode, Country FROM Suppliers;",
			"",
		},
		{
			"INSERT INTO test_table values (?, ?, ?, ?);",
			"?, ?, ?, ?",
		},
	} {
		ret := Util_ExportInsertQueryValues(test.input)
		require.Equal(t, test.expect, ret)
	}
}
