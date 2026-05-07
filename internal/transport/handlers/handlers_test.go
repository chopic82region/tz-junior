package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/chopic82region/tz-junior.git/internal/config"
	"github.com/chopic82region/tz-junior.git/internal/models"
	"github.com/chopic82region/tz-junior.git/internal/service/apperrors"
	"github.com/chopic82region/tz-junior.git/internal/transport/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userRepoStub struct {
	createErr error
	getErr    error
	updateErr error
	deleteErr error
	showErr   error
}

func (s userRepoStub) Create_user(ctx context.Context, user *models.User) error { return s.createErr }
func (s userRepoStub) GetById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return &models.User{ID: id, Name: "n", Email: "e"}, nil
}
func (s userRepoStub) Update_user(ctx context.Context, id uuid.UUID, user *models.User) error { return s.updateErr }
func (s userRepoStub) Delete(ctx context.Context, id uuid.UUID) error                        { return s.deleteErr }
func (s userRepoStub) Show_Users(ctx context.Context) ([]models.User, error) {
	if s.showErr != nil {
		return nil, s.showErr
	}
	return []models.User{}, nil
}

type subRepoStub struct {
	createErr error
	getErr    error
	cancelErr error
	showErr   error
}

func (s subRepoStub) Create_subscription(ctx context.Context, sub *models.Subscription) error {
	return s.createErr
}
func (s subRepoStub) GetById(ctx context.Context, id int) (*models.Subscription, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return &models.Subscription{ID: id}, nil
}
func (s subRepoStub) Cancel_subscription(ctx context.Context, id int) error { return s.cancelErr }
func (s subRepoStub) Show_subscription(ctx context.Context) ([]models.Subscription, error) {
	if s.showErr != nil {
		return nil, s.showErr
	}
	return []models.Subscription{}, nil
}

type filterRepoStub struct {
	err error
}

func (s filterRepoStub) GetTotalCost(ctx context.Context, userID uuid.UUID, start, end time.Time, name string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return "0", nil
}

func newTestRouter(u userRepoStub, s subRepoStub, f filterRepoStub) *gin.Engine {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{ServerPort: "0", MigrationsPath: "migrations"}
	h := handlers.NewHandler(u, s, f, cfg)

	r := gin.New()
	r.POST("/subscriptions", h.Create_subscription)
	r.GET("/subscriptions/:id", h.Get_subscription_by_id)
	return r
}

func TestCreateSubscription_EndDateOptional_DefaultApplied(t *testing.T) {
	t.Parallel()

	r := newTestRouter(userRepoStub{}, subRepoStub{}, filterRepoStub{})

	body := `{
		"user_id":"` + uuid.New().String() + `",
		"service_name":"netflix",
		"price":"10",
		"payment_time":"2026-05-07T12:00:00Z"
	}`

	req := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", w.Code, w.Body.String())
	}

	var got models.Subscription
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if got.End_date.IsZero() {
		t.Fatalf("expected end_date default to be set, got zero")
	}
	if got.Payment_time.IsZero() {
		t.Fatalf("expected payment_time to be set")
	}
	wantEnd := got.Payment_time.AddDate(0, 1, 0)
	if !got.End_date.Equal(wantEnd) {
		t.Fatalf("expected end_date=%s, got %s", wantEnd, got.End_date)
	}
}

func TestHandlers_ErrorMapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    error
		want   int
		method string
		url    string
		body   string
	}{
		{
			name:   "NilField => 400",
			err:    apperrors.NilField,
			want:   http.StatusBadRequest,
			method: http.MethodPost,
			url:    "/subscriptions",
			body:   `{"user_id":"` + uuid.New().String() + `","service_name":"x","price":"10","payment_time":"2026-05-07T12:00:00Z"}`,
		},
		{
			name:   "NotFound => 404",
			err:    apperrors.NotFound,
			want:   http.StatusNotFound,
			method: http.MethodGet,
			url:    "/subscriptions/1",
			body:   "",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := newTestRouter(userRepoStub{}, subRepoStub{createErr: tc.err, getErr: tc.err}, filterRepoStub{})

			var req *http.Request
			if tc.body == "" {
				req = httptest.NewRequest(tc.method, tc.url, nil)
			} else {
				req = httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.want {
				t.Fatalf("expected %d, got %d, body=%s", tc.want, w.Code, w.Body.String())
			}
		})
	}
}

