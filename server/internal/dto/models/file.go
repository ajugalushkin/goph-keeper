package models

import "io"

type File struct {
	Name   string
	Size   int64
	Type   string
	UserID int64
	Data   io.Reader
}
