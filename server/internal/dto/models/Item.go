package models

import "github.com/google/uuid"

type Item struct {
	Name    string
	Content []byte
	Version uuid.UUID
	OwnerID int64
}
