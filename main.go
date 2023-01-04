package main

import (
	"os"

	"ariga.io/atlas/sql/schema"
	"github.com/gokch/ornn/atlas"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db"
	"github.com/gokch/ornn/db/db_mysql"
	"github.com/gokch/ornn/db/db_postgres"
	"github.com/gokch/ornn/db/db_sqlite"
	"github.com/gokch/ornn/ornn"
	"github.com/gokch/ornn/parser"
	"github.com/gokch/ornn/parser/parser_mysql"
	"github.com/gokch/ornn/parser/parser_postgres"
	"github.com/gokch/ornn/parser/parser_sqlite"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ornn",
		Short: "ornn is a code generator for golang",
		Long:  "ornn is a code generator for golang",
		Run:   rootRun,
	}
	logger = log.Logger

	// db config
	DbType string
	DbPath string
	Addr   string
	Port   string
	Id     string
	Pw     string
	DbName string

	// gen config
	ConfigFilePath string
	GenFilePath    string
	PackageName    string
	ClassName      string
	DoNotEdit      string
)

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVar(&DbType, "db_type", "mysql", "database type ( mysql, mariadb, postgres, sqlite, tidb, cockroachdb )")
	fs.StringVar(&DbPath, "path", "./", "path for save db files. sqlite only")
	fs.StringVar(&Addr, "address", "127.0.0.1", "database server address")
	fs.StringVar(&Port, "port", "3306", "database server port")
	fs.StringVar(&Id, "id", "root", "database server id")
	fs.StringVar(&Pw, "pw", "", "database server password")
	fs.StringVar(&DbName, "db_name", "test", "database name")

	fs.StringVar(&ConfigFilePath, "config", "./output/config.json", "config json file path")
	fs.StringVar(&GenFilePath, "gen_file", "./output/gen.go", "generate golang file path")
	fs.StringVar(&PackageName, "package_name", "gen", "package name")
	fs.StringVar(&ClassName, "class_name", "Gen", "class name")
	fs.StringVar(&DoNotEdit, "do_not_edit", "// Code generated - DO NOT EDIT.\n// This file is a generated and any changes will be lost.\n", "do not edit comment")

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	var err error
	dbType := atlas.DbTypeStrReverse[DbType]

	// 1. connect db
	var conn *db.Conn
	switch dbType {
	case atlas.DbTypeMySQL, atlas.DbTypeMaria, atlas.DbTypeTiDB:
		conn, err = db_mysql.New(Addr, Port, Id, Pw, DbName)
	case atlas.DbTypePostgre, atlas.DbTypeCockroachDB:
		conn, err = db_postgres.New(Addr, Port, Id, Pw, DbName)
	case atlas.DbTypeSQLite:
		conn, err = db_sqlite.New(DbPath)
	default:
		logger.Fatal().Msgf("invalid db type: %s", DbType)
	}
	if err != nil {
		logger.Fatal().Err(err).Msg("db connect error")
	}

	// 2. init schema from atl ( TODO : Migrate schema from atl )
	var sch *schema.Schema
	{
		atl := &atlas.Atlas{}
		atl.Init(dbType, conn)
		sch, err = atl.InspectSchema()
		if err != nil {
			logger.Fatal().Err(err).Msg("atlas inspect error")
		}
	}

	// 3. set config
	var conf *config.Config = &config.Config{}
	{
		if err = conf.Load(ConfigFilePath); err != nil { // load
			logger.Fatal().Err(err).Msg("config load error")
		}
		if err = conf.Init(dbType, sch, DoNotEdit, PackageName, ClassName); err != nil { // init
			logger.Fatal().Err(err).Msg("config init error")
		}
		if err = conf.Save("./output/gen.json"); err != nil { // save
			logger.Fatal().Err(err).Msg("config save error")
		}
	}

	// 4. set parser
	var psr parser.Parser
	switch dbType {
	case atlas.DbTypeMySQL, atlas.DbTypeMaria, atlas.DbTypeTiDB:
		psr = parser_mysql.New(sch)
	case atlas.DbTypePostgre, atlas.DbTypeCockroachDB:
		psr = parser_postgres.New(sch)
	case atlas.DbTypeSQLite:
		psr = parser_sqlite.New(sch)
	default:
		logger.Fatal().Msgf("invalid db type: %s", DbType)
	}

	// 5. init ornn
	var ornn *ornn.ORNN = &ornn.ORNN{}
	{
		ornn.Init(conf, psr)
		if err = ornn.GenCode("./output/gen.go"); err != nil { // code generate
			logger.Fatal().Err(err).Msg("code generate error")
		}
	}
}
