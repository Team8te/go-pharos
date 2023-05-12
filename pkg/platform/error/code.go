package error

// Code ...
type Code string

const (
	ErrorCodeUndefined = Code("undefined")
	ErrorInternal      = Code("internal")
	ErrorCodeDB        = Code("db_error")
)
