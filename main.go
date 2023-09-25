package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"myits-gate-api/handler"
	"myits-gate-api/repository"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var (
		dbhost     string = os.Getenv("DB_HOST")
		dbport     string = os.Getenv("DB_PORT")
		dbname     string = os.Getenv("DB_NAME")
		dbuser     string = os.Getenv("DB_USER")
		dbpassword string = os.Getenv("DB_PASSWORD")
	)

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", dbhost, dbuser, dbpassword, dbport, dbname)

	// Create connection pool
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	ctx := context.Background()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")

	gateAccessRepo := repository.NewSqlServerGateAccessRepository(db, ctx)
	gateLogRepo := repository.NewSqlServerGateLogRepository(db, ctx)

	inquiryHandler := handler.NewInquiryHandler(gateAccessRepo, gateLogRepo)

	mode := os.Getenv("GIN_MODE")

	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.POST("/inquiry", inquiryHandler.PostInquiry)
	router.POST("/inquiry-complete", inquiryHandler.PostInquiryComplete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "33000"
		log.Printf("Defaulting to port %s", port)
	}

	log.Fatal(router.Run(":" + port))
}
