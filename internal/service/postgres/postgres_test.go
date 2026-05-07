package postgres

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chopic82region/tz-junior.git/internal/models"
	"github.com/chopic82region/tz-junior.git/internal/service/apperrors"
	"github.com/google/uuid"
)

func TestSubscriptionRepository_CreateSubscription_Validation(t *testing.T) {
	t.Parallel()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewSubscriptionRepo(db)

	t.Run("nil subscription", func(t *testing.T) {
		t.Parallel()
		if err := repo.Create_subscription(context.Background(), nil); err == nil || err != apperrors.NilField {
			t.Fatalf("expected NilField, got %v", err)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		t.Parallel()
		sub := &models.Subscription{}
		if err := repo.Create_subscription(context.Background(), sub); err == nil || err != apperrors.NilField {
			t.Fatalf("expected NilField, got %v", err)
		}
	})

	t.Run("bad price", func(t *testing.T) {
		t.Parallel()
		sub := &models.Subscription{
			UserID:       uuid.New(),
			Service:      "netflix",
			Price:        "abc",
			Payment_time: time.Now(),
			End_date:     time.Now().AddDate(0, 1, 0),
		}
		if err := repo.Create_subscription(context.Background(), sub); err == nil || err != apperrors.BadPayload {
			t.Fatalf("expected BadPayload, got %v", err)
		}
	})
}

func TestSubscriptionRepository_CreateSubscription_DefaultEndDate(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewSubscriptionRepo(db)

	payment := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	sub := &models.Subscription{
		UserID:       uuid.New(),
		Service:      "spotify",
		Price:        "199.99",
		Payment_time: payment,
		// End_date intentionally omitted (optional)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO subscriptions (user_id, service_name, price, payment_time, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`)).
		WithArgs(sub.UserID, sub.Service, sub.Price, sub.Payment_time, payment.AddDate(0, 1, 0)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	if err := repo.Create_subscription(context.Background(), sub); err != nil {
		t.Fatalf("Create_subscription error: %v", err)
	}
	if sub.ID != 42 {
		t.Fatalf("expected ID=42, got %d", sub.ID)
	}
	if sub.End_date.IsZero() {
		t.Fatalf("expected End_date to be set")
	}
	if !sub.End_date.Equal(payment.AddDate(0, 1, 0)) {
		t.Fatalf("expected End_date=%s, got %s", payment.AddDate(0, 1, 0), sub.End_date)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

