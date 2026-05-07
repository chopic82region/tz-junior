package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/chopic82region/tz-junior.git/internal/service/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/chopic82region/tz-junior.git/internal/config"
	"github.com/chopic82region/tz-junior.git/internal/models"
	"github.com/chopic82region/tz-junior.git/internal/repository/repository"
)

type Handler struct {
	filter repository.Filter_interface
	users  repository.User_interface
	subs   repository.Subscription_interface
	cfg    *config.Config
}

func NewHandler(users repository.User_interface, subs repository.Subscription_interface, filter repository.Filter_interface, cfg *config.Config) *Handler {
	return &Handler{
		users:  users,
		subs:   subs,
		filter: filter,
		cfg:    cfg,
	}
}

func writeError(c *gin.Context, err error) {
	switch {
	case err == nil:
		return
	case errors.Is(err, apperrors.NilField), errors.Is(err, apperrors.BadPayload), errors.Is(err, apperrors.InvalidID):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, apperrors.NotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, apperrors.Duplicate):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

//-------------------User handlers------------------

// POST /users
func (h *Handler) Create_user(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.users.Create_user(c.Request.Context(), &user); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(201, user)
}

// GET /users/:id
func (h *Handler) Get_user_by_id(c *gin.Context) {
	id := c.Param("id")

	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.users.GetById(c.Request.Context(), idUUID)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(200, user)
}

// PATCH /users/:id
func (h *Handler) Update_user(c *gin.Context) {
	id := c.Param("id")

	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.users.Update_user(c.Request.Context(), idUUID, &user); err != nil {
		writeError(c, err)
		return
	}

	c.JSON(200, user)
}

// DELETE /users/:id
func (h *Handler) Delete_user(c *gin.Context) {
	id := c.Param("id")

	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.users.Delete(c.Request.Context(), idUUID); err != nil {
		writeError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "user deleted successfully"})
}

// GET /users
func (h *Handler) Show_users(c *gin.Context) {
	users, err := h.users.Show_Users(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(200, users)
}

//-------------------Subscription handlers------------------

func (h *Handler) Create_subscription(c *gin.Context) {
	var subscription models.Subscription

	if err := c.BindJSON(&subscription); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	// end_date is optional: if not provided, set default in a predictable way.
	if subscription.End_date.IsZero() && !subscription.Payment_time.IsZero() {
		subscription.End_date = subscription.Payment_time.AddDate(0, 1, 0)
	}
	if err := h.subs.Create_subscription(c.Request.Context(), &subscription); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(201, subscription)
}

func (h *Handler) Get_subscription_by_id(c *gin.Context) {

	id := c.Param("id")

	idint, err := strconv.Atoi(id)
	if err != nil || idint == 0 {
		c.JSON(400, gin.H{"error": "invalid subscription ID"})
		return
	}

	subscription, err := h.subs.GetById(c.Request.Context(), idint)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(200, subscription)
}

func (h *Handler) Cancel_subscription(c *gin.Context) {
	id := c.Param("id")

	idint, err := strconv.Atoi(id)
	if err != nil || idint == 0 {
		c.JSON(400, gin.H{"error": "invalid subscription ID"})
		return
	}
	// Return the cancelled subscription data (including end_date) for better client UX.
	subscription, err := h.subs.GetById(c.Request.Context(), idint)
	if err != nil {
		writeError(c, err)
		return
	}
	if err := h.subs.Cancel_subscription(c.Request.Context(), idint); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"message":      "subscription canceled successfully",
		"subscription": subscription,
	})
}

func (h *Handler) Show_subscription(c *gin.Context) {
	subscriptions, err := h.subs.Show_subscription(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(200, subscriptions)
}

//-------------------Filter handlers------------------

func (h *Handler) Get_total_cost(c *gin.Context) {
	userID := c.Query("user_id")
	start := c.Query("created_at")
	end := c.Query("end")
	name := c.Query("name")

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid start time"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid end time"})
		return
	}

	totalCost, err := h.filter.GetTotalCost(c.Request.Context(), userIDUUID, startTime, endTime, name)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(200, gin.H{"total_cost": totalCost})
}
