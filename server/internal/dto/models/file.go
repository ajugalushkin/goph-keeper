package models

import "io"

type File struct {
	Name        string
	NameWithExt string
	Type        string
	Size        int64
	UserID      int64
	Data        io.Reader
}
