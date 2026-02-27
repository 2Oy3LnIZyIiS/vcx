package datatype

type DataType int

const (
	UNKNOWN DataType = iota
	TEXT
	INT
	BINARY
	BYTES
	HEX
)
