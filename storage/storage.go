package storage

import "errors"

var(
	ErrBadRequest = errors.New("Bad Request")
	ErrNotFound   = errors.New("Not Found")
	ErrInternal   = errors.New("Internal Error")
	ErrUnknown    = errors.New("Unknown Error")
)

type Storage interface {
	Query(q string) (results []byte, err error)
	// lp, Line Protocol, contains the time series data that you want to store. Its components are measurement, tags, fields and timestamp.
	Write(db string, lp []byte)(err error)
	Ping() (version string, err error)
	Close() (err error)
}