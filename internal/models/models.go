package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID      int       `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	Service string    `json:"service_name"`
	Price   string    `json:"price"`

	Payment_time time.Time `json:"payment_time"`
	End_date     time.Time `json:"end_date"`
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`

	Create_at time.Time `json:"created_at"`
}
