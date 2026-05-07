package main

import (
	"log"

	"github.com/chopic82region/tz-junior.git/internal/config"
	"github.com/chopic82region/tz-junior.git/internal/repository/db"
	"github.com/chopic82region/tz-junior.git/internal/service/postgres"
	"github.com/chopic82region/tz-junior.git/internal/transport/server"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error of load config")
	}

	// бд
	dbConn, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Error of connect to db: %v", err)
	}
	if err := db.RunMigrations(*cfg, dbConn); err != nil {
		log.Fatalf("Error of run migrations: %v", err)
	}

	// Инициализация репозиториев, сервисов и сервера
	URepo := postgres.NewUserRepo(dbConn)
	SRepo := postgres.NewSubscriptionRepo(dbConn)
	FRepo := postgres.NewFilterRepo(dbConn)

	srv := server.NewServer(URepo, SRepo, FRepo, cfg, dbConn)

	if err := srv.StartServer(); err != nil {
		log.Fatalf("Error of start server: %v", err)
	}

}
