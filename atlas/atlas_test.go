package atlas

import (
	"os"
	"testing"

	"github.com/gokch/ornn/db/db_mysql"
	"github.com/stretchr/testify/require"
)

func TestMysql(t *testing.T) {
	db, err := db_mysql.New("127.0.0.1", "3306", "root", "951753ck", "test")
	require.NoError(t, err)

	atlas := &Atlas{}
	err = atlas.Init(DbTypeMySQL, db)
	require.NoError(t, err)

	schema, err := atlas.InspectSchema()
	require.NoError(t, err)

	bt, err := atlas.MarshalHCL(schema)
	require.NoError(t, err)

	err = os.WriteFile("./schema.hcl", bt, 0700)
	require.NoError(t, err)

	schemaNew, err := atlas.UnmarshalHCL(bt)
	require.NoError(t, err)

	require.EqualValues(t, schema, schemaNew)
}
