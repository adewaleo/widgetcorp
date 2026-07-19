package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/adewaleo/widgetcorp/widget-srv/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// getenv returns the env var named key, or fallback if it is unset/empty.
// A tiny helper so config has sensible local-dev defaults.
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	// Default DSN matches the credentials in compose.yml. sslmode=disable
	// because the local container serves plaintext.
	dsn := getenv("DATABASE_URL", "postgres://widget:widget@localhost:5432/widgetdb?sslmode=disable")
	srvPort := getenv("PORT", "8080")

	// 1. Bring the schema up to date before doing anything else.
	log.Println("running migrations...")
	if err := db.RunMigrations(dsn); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	log.Println("migrations applied (schema is current)")

	// TODO (next stages):
	//   - open a pgxpool connection using dsn and ping to verify connectivity

	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer dbPool.Close()

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	//   - build the gin router with our handlers
	srv := gin.Default()
	srv.GET(
		"/healthz",
		func(c *gin.Context) {
			err := dbPool.Ping(c)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		},
	)

	//   - listen on PORT
	err = srv.Run(":" + srvPort)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
