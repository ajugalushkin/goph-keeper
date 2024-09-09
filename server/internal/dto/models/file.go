package models

import (
	"io"
)

type Data interface {
	io.Reader
	io.Closer
}

type File struct {
	Item Item
	Size int64
	Data Data
}
