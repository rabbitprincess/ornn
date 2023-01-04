package ornn

import (
	"fmt"
	"os"

	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/sql/parser"
)

type ORNN struct {
	conf *config.Config
	psr  parser.Parser
}

func (t *ORNN) Init(conf *config.Config, psr parser.Parser) {
	t.conf = conf
	t.psr = psr
}

func (t *ORNN) GenCode(path string) (err error) {
	if t.conf == nil {
		return fmt.Errorf("config is emtpy")
	}

	// gen code
	gen := &Gen{}
	code, err := gen.Gen(t.conf, t.psr, path)
	if err != nil {
		return err
	}

	// write code to file
	err = os.WriteFile(path, []byte(code), 0700)
	if err != nil {
		return err
	}
	return nil
}
