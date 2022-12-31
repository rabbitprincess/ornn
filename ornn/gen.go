package ornn

import (
	"fmt"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db"
	"github.com/gokch/ornn/db/db_mysql"
)

type Gen struct {
	data *GenData
	code *GenCode
}

func (t *Gen) Gen(conn *db.Conn, conf *config.Config, path string) (code string, err error) {
	vendor := db_mysql.NewVendor(conn)

	// 코드 생성을 위한 gen 데이터 준비
	t.data = &GenData{}
	t.data.Init(conf, conn, vendor)
	err = t.data.SetData()
	if err != nil {
		return "", err
	}

	// check error
	for tableName, def := range conf.Queries.Default {
		for _, query := range def {
			if query.ErrParser != "" {
				fmt.Printf("parser err - table : %s | query : %s | err : %s\n", tableName, query.Name, query.ErrParser)
				err = fmt.Errorf("query error")
			}
			if query.ErrQuery != "" {
				fmt.Printf("query err - table : %s | query : %s | err : %s\n", tableName, query.Name, query.ErrQuery)
				err = fmt.Errorf("query error")
			}
		}
	}

	// gen code
	t.code = &GenCode{}
	code, err = t.code.code(conf, t.data)
	if err != nil {
		return "", err
	}

	return code, nil
}
