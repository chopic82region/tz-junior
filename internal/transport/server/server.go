package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chopic82region/tz-junior.git/internal/config"
	"github.com/chopic82region/tz-junior.git/internal/repository/repository"
	"github.com/chopic82region/tz-junior.git/internal/transport/handlers"
	"github.com/gin-gonic/gin"
)

type Server struct {
	handlers *handlers.Handler
	cfg      *config.Config
	db       *sql.DB
}

func NewServer(
	users repository.User_interface,
	subs repository.Subscription_interface,
	filter repository.Filter_interface,
	cfg *config.Config,
	db *sql.DB,
) *Server {
	return &Server{
		handlers: handlers.NewHandler(users, subs, filter, cfg),
		cfg:      cfg,
		db:       db,
	}
}

func (s *Server) NewRouter(db *sql.DB) *gin.Engine {

	router := gin.Default()
	//User routes
	router.POST("/users", s.handlers.Create_user)
	router.GET("/users/:id", s.handlers.Get_user_by_id)
	router.PATCH("/users/:id", s.handlers.Update_user)
	router.DELETE("/users/:id", s.handlers.Delete_user)
	router.GET("/users", s.handlers.Show_users)

	//Subscription routes
	router.POST("/subscriptions", s.handlers.Create_subscription)
	router.GET("/subscriptions/:id", s.handlers.Get_subscription_by_id)
	router.DELETE("/subscriptions/:id", s.handlers.Cancel_subscription)
	router.GET("/subscriptions", s.handlers.Show_subscription)

	//Filter routes
	router.GET("/filter/total_cost", s.handlers.Get_total_cost)

	return router
}

func (s *Server) StartServer() error {
	port := s.cfg.ServerPort
	if port == "" {
		port = os.Getenv("SERVER_PORT")
	}
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	// Проверка соединения, если БД передана
	if s.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	if err := s.db.Ping(); err != nil {
		return err
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           s.NewRouter(s.db),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", addr)
		errCh <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-stop:
		log.Printf("shutdown signal: %s", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	case err := <-errCh:
		if err == nil || err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
