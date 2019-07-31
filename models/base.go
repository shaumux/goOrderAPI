package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type ModelInterface interface {
	Create() (interface{}, error)
	Update(map[string]interface{}) (interface{}, error)
	Delete() error
	Get() (interface{}, error)
	All(int, int) (interface{}, error)
	Filter() (interface{}, error)
}

type Base struct {
	ModelInterface `gorm:"-" json:"-"`
	ID             uuid.UUID `sql:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
