package repository

import (
	"context"

	"github.com/chopic82region/tz-junior.git/internal/models"
	"github.com/google/uuid"
)

type User_interface interface {
	Create_user(ctx context.Context, user *models.User) error
	GetById(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update_user(ctx context.Context, id uuid.UUID, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Show_Users(ctx context.Context) ([]models.User, error)
}

type Subscrioption_intarface interface {
	Create_subscription(ctx context.Context, subscription *models.Subscription) error
	GetById(ctx context.Context, id int) (*models.Subscription, error)
	Update_subscription(ctx context.Context, id int, subscription *models.Subscription) error
	Cancel_subscription(ctx context.Context, id int) error
	Show_subscription(ctx context.Context) ([]models.Subscription, error)
}

type Filter_interface interface {
	Filter_by_user_id(ctx context.Context, user_id int)
}
