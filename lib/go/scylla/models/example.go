package models

import "github.com/google/uuid"

type Example struct {
	ID uuid.UUID
	A  int64
	B  string
}
