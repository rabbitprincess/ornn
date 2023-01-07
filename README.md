## ornn: Object Relation No Need

Generate database interaction code from schema and sql

built in atlas


오른이 되세요

## Quick Start

```go
package main

import (
	"os"

	"github.com/gokch/ornn/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
```

## Build

    go get github.com/gokch/ornn
    cd $GOPATH/src/github.com/gokch/ornn
    go build .


## Run

```
Usage:
  ornn [flags]

Flags:
      --address string        database server address (default "127.0.0.1")
      --class_name string     class name (default "Gen")
      --config_path string    config json file path (default "./output/config.json")
      --db_name string        database name (default "test")
      --db_type string        database type ( mysql, mariadb, postgres, sqlite, tidb, cockroachdb ) (default "mysql")
      --do_not_edit string    do not edit comment (default "// Code generated - DO NOT EDIT.\n// This file is a generated and any changes will be lost.\n")
      --gen_file string       generate golang file path (default "./output/gen.go")
  -h, --help                  help for ornn
      --id string             database server id (default "root")
      --load_config_file      load config from existing file
      --load_schema_file      load schema from existing file and migrate database
      --package_name string   package name (default "gen")
      --path string           path for save db files. sqlite only (default "./output")
      --port string           database server port (default "3306")
      --pw string             database server password (default "1234")
      --schema_path string    schema hcl file path (default "./output/schema.hcl")
```

Example

    ./ornn --id root --pw 1234
