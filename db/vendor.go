package db

type Vendor interface {
	ConvType(dbType string) (genType string)
	CreateTable() (sql []string, err error)
}
