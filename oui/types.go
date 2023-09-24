package oui

type VendorDef struct {
	Prefix   string
	Length   int
	Org      string
	Registry string
}

func (v *VendorDef) PrefixString() string {
	if v == nil {
		return "<nil>"
	}
	return v.Prefix
}

type LoggerType interface {
	Success(s string, f ...interface{})
	Info(s string, f ...interface{})
	Warn(s string, f ...interface{})
	Error(s string, f ...interface{})
	Err(err error, strs ...string)
}

const (
	dialectSqlite int = iota
	dialectPsql
)

const (
	maxVarsSqlite int = 999
	maxVarsPsql   int = 65535
)
