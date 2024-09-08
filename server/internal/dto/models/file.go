package models

import (
	"io"
)

type Data interface {
	io.Reader
	io.Closer
}

type File struct {
	Name        string
	NameWithExt string
	Size        int64
	UserID      int64
	Data        Data
}
