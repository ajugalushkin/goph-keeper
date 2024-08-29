package models

type Item struct {
	ID              string `json:"id" db:"id"`
	Name            string `json:"name" db:"name"`
	Type            string `json:"type" db:"type"`
	Value           []byte `json:"value" db:"value"`
	ServerUpdatedAt string `json:"server_updated_at" db:"server_updated_at"`
	IsDeleted       bool   `json:"is_deleted" db:"is_deleted"`
}

type ListItem []Item
