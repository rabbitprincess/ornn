package go_orm_gen

import (
	"fmt"
	"io/ioutil"

	"github.com/gokch/go-orm-gen/config"
	"github.com/gokch/go-orm-gen/db"
)

type Gen struct {
	data *GenData
	code *GenCode
}

func (t *Gen) Gen(_pc_db *db.DB, config *config.Config, _maps_gen_config map[string]string, _s_filepath string) (code string, err error) {
	// json -> 코드 생성을 위한 gen 데이터 준비
	t.data = &GenData{}
	t.data.Init(_pc_db)
	err = t.data.prepare_db_schema(config)
	if err != nil {
		return "", err
	}

	// check error
	for _, table := range config.Schema.Tables {
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

	// gen code
	t.code = &GenCode{}
	code, err = t.code.gen_source_code(config, _maps_gen_config, t.data)
	if err != nil {
		return "", err
	}

	// write file
	err = ioutil.WriteFile(_s_filepath, []byte(code), 0700)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (t *Gen) printError(config *config.Config) {

}
