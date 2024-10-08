package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/hiroshijp/try-clean-arch/handler"
	"github.com/hiroshijp/try-clean-arch/handler/middleware"
	"github.com/hiroshijp/try-clean-arch/handler/public"
	postgresRepo "github.com/hiroshijp/try-clean-arch/repository/postgres"
	"github.com/hiroshijp/try-clean-arch/usecase"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

// set deafult
const (
	defaultAddress = ":8080"
)

func main() {
	//  prepare database source name
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// prepare echo
	e := echo.New()

	// prepare repositry
	txRepo := postgresRepo.NewTxRepository(dbConn)
	historyRepo := postgresRepo.NewHistoryRepository(dbConn)
	visitorRepo := postgresRepo.NewVisitorRepository(dbConn)
	userRepo := postgresRepo.NewUserRepository(dbConn)

	// prepare usecase
	historyUsecase := usecase.NewHistoryUsecase(txRepo, historyRepo, visitorRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// prepare middleware and handler
	middleware.NewCORSMiddleware(e, os.Getenv("ALLOWED_ORIGIN"))
	public.NewVisitedHandler(e, historyUsecase)
	public.NewSigninHandler(e, userUsecase)

	api := e.Group("/api")
	middleware.NewJWTMiddleware(api)
	handler.NewHistoryHandler(api, historyUsecase)

	// start server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}

	log.Fatal(e.Start(address))
}
