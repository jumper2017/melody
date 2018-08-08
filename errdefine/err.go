package errdefine

import "errors"

var (
	ERR_CTX_PARSE error = errors.New("parse context failed.")

	ERR_ETCD_INVALID_PARAM     error = errors.New("invalid param.")
	ERR_ETCD_SERVICE_NOT_FOUND error = errors.New("couldn't find service specified.")
)
