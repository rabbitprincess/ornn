// Example usage:
//
//	package main
//
//	import (
//		"os"
//
//	    ornn "https://github.com/gokch/ornn/cli"
//	)
//
//	func main() {
//		if err := ornn.Run(os.Args[1:]); err != nil {
//			os.Exit(1)
//		}
//	}
package cli

import (
	"github.com/rs/zerolog/log"

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
	"github.com/spf13/cobra"
)

func Run(args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

var (
	rootCmd = &cobra.Command{
		Use:   "ornn",
		Short: "ornn is a code generator for golang",
		Long:  "ornn is a code generator for golang db access",
		Run:   rootRun,
	}
	logger = log.Logger

	loadExistSchemaFile bool // 기존 스키마 파일에서 로딩, 스키마 파일대로 db migrate
	loadExistConfigFile bool // 기존 설정 파일에서 로딩

	// db config
	dbType string
	dbPath string
	addr   string
	port   string
	id     string
	pw     string
	dbName string

	schemaFilePath string
	configFilePath string
	genFilePath    string
	packageName    string
	className      string
	doNotEdit      string
)

func init() {
	fs := rootCmd.PersistentFlags()

	fs.StringVarP(&dbType, "db_type", "D", "mysql", "database type ( mysql, mariadb, postgres, sqlite, tidb, cockroachdb )")
	fs.StringVarP(&addr, "db_addr", "A", "127.0.0.1", "database server address")
	fs.StringVarP(&port, "db_port", "P", "3306", "database server port")
	fs.StringVarP(&id, "db_id", "i", "root", "database server id")
	fs.StringVarP(&pw, "db_pw", "p", "1234", "database server password")
	fs.StringVar(&dbName, "db_name", "test", "database name")
	fs.StringVar(&dbPath, "db_path", "./output", "path for save db files. sqlite only")

	fs.BoolVar(&loadExistSchemaFile, "file_schema_load", false, "load schema from existing file and migrate database")
	fs.BoolVar(&loadExistConfigFile, "file_config_load", false, "load config from existing file")
	fs.StringVar(&schemaFilePath, "file_schema_path", "./output/schema.hcl", "schema hcl file path")
	fs.StringVar(&configFilePath, "file_config_path", "./output/config.json", "config json file path")
	fs.StringVar(&genFilePath, "file_gen_path", "./output/gen.go", "generate golang file path")
	fs.StringVar(&packageName, "gen_package_name", "gen", "package name")
	fs.StringVar(&className, "gen_class_name", "Gen", "class name")
	fs.StringVar(&doNotEdit, "gen_do_not_edit", "// Code generated - DO NOT EDIT.\n// This file is a generated and any changes will be lost.\n", "do not edit comment")
}

func rootRun(cmd *cobra.Command, args []string) {
	var err error
	atlasDbType := atlas.DbTypeStrReverse[dbType]

	// 1. connect db
	var conn *db.Conn
	switch atlasDbType {
	case atlas.DbTypeMySQL, atlas.DbTypeMaria, atlas.DbTypeTiDB:
		conn, err = db_mysql.New(addr, port, id, pw, dbName)
	case atlas.DbTypePostgre, atlas.DbTypeCockroachDB:
		conn, err = db_postgres.New(addr, port, id, pw, dbName)
	case atlas.DbTypeSQLite:
		conn, err = db_sqlite.New(dbPath)
	default:
		logger.Fatal().Msgf("invalid db type: %s", dbType)
	}
	if err != nil {
		logger.Fatal().Err(err).Msg("db connect error")
	}

	// 2. init schema from atl
	var sch *schema.Schema
	atl := atlas.New(atlasDbType, conn)
	if loadExistSchemaFile { // load from existing schema file
		if sch, err = atl.Load(schemaFilePath); err != nil {
			logger.Fatal().Err(err).Msg("schema load error")
		}
		// migrate db from file
		if err = atl.MigrateSchema(sch); err != nil {
			logger.Fatal().Err(err).Msg("atlas migrate error")
		}
		// inspect schema fron migrated db
		if sch, err = atl.InspectSchema(); err != nil {
			logger.Fatal().Err(err).Msg("atlas inspect error")
		}
	} else {
		if sch, err = atl.InspectSchema(); err != nil {
			logger.Fatal().Err(err).Msg("atlas inspect error")
		}
		if err = atl.Save(schemaFilePath, sch); err != nil {
			logger.Fatal().Err(err).Msg("schema save error")
		}
	}

	// 3. set config
	var conf *config.Config = &config.Config{}
	if loadExistConfigFile { // load from existing config file
		if err = conf.Load(configFilePath); err != nil { // load
			logger.Fatal().Err(err).Msg("config load error")
		}
		if err = conf.Init(atlasDbType, sch, packageName, className, doNotEdit); err != nil { // init
			logger.Fatal().Err(err).Msg("config init error")
		}
	} else {
		if err = conf.Init(atlasDbType, sch, packageName, className, doNotEdit); err != nil { // init
			logger.Fatal().Err(err).Msg("config init error")
		}
		if err = conf.Save(configFilePath); err != nil { // save
			logger.Fatal().Err(err).Msg("config save error")
		}
	}

	// 4. set parser
	var psr parser.Parser
	switch atlasDbType {
	case atlas.DbTypeMySQL, atlas.DbTypeMaria, atlas.DbTypeTiDB:
		psr = parser_mysql.New(&conf.Schema)
	case atlas.DbTypePostgre, atlas.DbTypeCockroachDB:
		psr = parser_postgres.New(&conf.Schema)
	case atlas.DbTypeSQLite:
		psr = parser_sqlite.New(&conf.Schema)
	default:
		logger.Fatal().Msgf("invalid db type: %s", dbType)
	}

	// 5. init ornn
	var ornn *ornn.ORNN = &ornn.ORNN{}
	{
		ornn.Init(conf, psr)
		if err = ornn.GenCode(genFilePath); err != nil { // code generate
			logger.Fatal().Err(err).Msg("code generate error")
		}
	}
}
