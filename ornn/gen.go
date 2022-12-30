package ornn

import (
	"os"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/db"
)

type Gen struct {
	data *GenData
	code *GenCode
}

func (t *Gen) Gen(conn *db.Conn, conf *config.Config, path string) (code string, err error) {
	// json -> 코드 생성을 위한 gen 데이터 준비
	t.data = &GenData{}
	t.data.Init(conn)
	err = t.data.SetData(conf)
	if err != nil {
		return "", err
	}

	// check error
	/*
		for _, table := range conf.Schema.Tables {
			for _, query := range table.Queries {
				if query.ErrParser != "" {
					fmt.Printf("parser err - table : %s | query : %s | err : %s\n", table.Name, query.Name, query.ErrParser)
					err = fmt.Errorf("query error")
				}
				if query.ErrQuery != "" {
					fmt.Printf("query err - table : %s | query : %s | err : %s\n", table.Name, query.Name, query.ErrQuery)
					err = fmt.Errorf("query error")
				}
			}
		}
		if err != nil {
			return "", err
		}
	*/
	// gen code
	t.code = &GenCode{}
	code, err = t.code.gen(conf, t.data)
	if err != nil {
		return "", err
	}

	// write file
	err = os.WriteFile(path, []byte(code), 0700)
	if err != nil {
		return "", err
	}

	return code, nil
}
