package bp

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	DEF_s_prefix = "\t"
	DEF_s_crlf   = "\n"
)

type Gen_writer struct {
	bt_buf   *bytes.Buffer
	i_write  io.Writer
	n_indent int
}

func (t *Gen_writer) Init() {
	t.bt_buf = &bytes.Buffer{}
	t.i_write = t.bt_buf
}

func (t *Gen_writer) String() string {
	return t.bt_buf.String()
}

func (t *Gen_writer) N(_s_format string, _arri ...interface{}) (err error) {
	s_new := fmt.Sprintf(_s_format, _arri...)
	_, err = t.i_write.Write([]byte(s_new))
	if err != nil {
		return err
	}

	return nil
}

func (t *Gen_writer) Get_indent() string {
	return strings.Repeat(DEF_s_prefix, t.n_indent)
}

func (t *Gen_writer) W(_s_format string, _arri ...interface{}) (err error) {
	t.i_write.Write([]byte(t.Get_indent()))

	s_new := fmt.Sprintf(_s_format, _arri...)
	_, err = t.i_write.Write([]byte(s_new))
	if err != nil {
		return err
	}

	return nil
}

func (t *Gen_writer) Indent__in() {
	t.n_indent++
}

func (t *Gen_writer) Indent__out() {
	t.n_indent--
}
