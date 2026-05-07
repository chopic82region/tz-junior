package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/chopic82region/tz-junior.git/internal/models"
	"github.com/chopic82region/tz-junior.git/internal/service/apperrors"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

// функционал для пользователя

type UserRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create_user(ctx context.Context, user *models.User) error {

	if user == nil {
		return apperrors.NilField
	}
	if user.Name == "" || user.Email == "" {
		return apperrors.NilField
	}

	query := `
		INSERT INTO users (name, email, created_at)
		VALUES ($1, $2, COALESCE($3, NOW()))
		RETURNING id, created_at
	`

	var createdAt time.Time
	if err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Cteate_at).Scan(&user.ID, &createdAt); err != nil {
		return err
	}
	user.Cteate_at = createdAt
	return nil
}

func (r *UserRepository) GetById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if id == uuid.Nil {
		return nil, apperrors.InvalidID
	}
	var u models.User
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Cteate_at)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update_user(ctx context.Context, id uuid.UUID, user *models.User) error {
	if id == uuid.Nil {
		return apperrors.InvalidID
	}
	if user == nil {
		return apperrors.NilField
	}
	if user.Name == "" || user.Email == "" {
		return apperrors.NilField
	}

	query := `
		UPDATE users
		SET name = $1, email = $2
		WHERE id = $3
	`
	res, err := r.db.ExecContext(ctx, query, user.Name, user.Email, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err == nil && aff == 0 {
		return apperrors.NotFound
	}
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return apperrors.InvalidID
	}

	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err == nil && aff == 0 {
		return apperrors.NotFound
	}
	return err
}
func (r *UserRepository) Show_Users(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, email, created_at FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Cteate_at); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// функционал для подписок пользователя

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) Create_subscription(ctx context.Context, subscription *models.Subscription) error {

	if subscription == nil {
		return apperrors.NilField
	}
	if subscription.UserID == uuid.Nil || subscription.Service == "" || subscription.Price == "" || subscription.Payment_time.IsZero() {
		return apperrors.NilField
	}

	query := `
		INSERT INTO subscriptions (user_id, service_name, price, payment_time)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var subscriptionID int
	if err := r.db.QueryRowContext(ctx, query, subscription.UserID, subscription.Service, subscription.Price, subscription.Payment_time).Scan(&subscriptionID); err != nil {
		return err
	}
	subscription.ID = subscriptionID
	return nil
}

func (r *SubscriptionRepository) GetById(ctx context.Context, id int) (*models.Subscription, error) {
	if id <= 0 {
		return nil, apperrors.InvalidID
	}
	var s models.Subscription
	query := `SELECT id, user_id, service_name, price, payment_time FROM subscriptions WHERE id = $1`
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.UserID, &s.Service, &s.Price, &s.Payment_time); err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *SubscriptionRepository) Update_subscription(ctx context.Context, id int, subscription *models.Subscription) error {
	if id <= 0 {
		return apperrors.InvalidID
	}
	if subscription == nil {
		return apperrors.NilField
	}
	if subscription.Service == "" || subscription.Price == "" || subscription.Payment_time.IsZero() {
		return apperrors.NilField
	}

	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, payment_time = $3
		WHERE id = $4
	`
	res, err := r.db.ExecContext(ctx, query, subscription.Service, subscription.Price, subscription.Payment_time, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err == nil && aff == 0 {
		return apperrors.NotFound
	}
	return err
}

func (r *SubscriptionRepository) Cancel_subscription(ctx context.Context, id int) error {
	if id <= 0 {
		return apperrors.InvalidID
	}

	// Отмена подписки: удаляем запись. При необходимости можно заменить на "status/cancelled_at".
	res, err := r.db.ExecContext(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err == nil && aff == 0 {
		return apperrors.NotFound
	}
	return err
}

func (r *SubscriptionRepository) Show_subscription(ctx context.Context) ([]models.Subscription, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, service_name, price, payment_time FROM subscriptions ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Subscription, 0)
	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.Service, &s.Price, &s.Payment_time); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
