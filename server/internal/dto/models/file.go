package models

import "io"

type File struct {
	Name   string
	Size   int64
	Type   string
	Bucket string
	Data   io.Reader
}
