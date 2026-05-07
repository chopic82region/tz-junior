package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID     int       `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Service string   `json:"service_name"`
	Price   string   `json:"price"`

	Payment_time time.Time `json:"start_data"`
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`

	Cteate_at time.Time `json:"created_at"`
}
