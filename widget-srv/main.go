package main

import (
	"log"
	"os"

	"github.com/adewaleo/widgetcorp/widget-srv/internal/db"
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

	// 1. Bring the schema up to date before doing anything else.
	log.Println("running migrations...")
	if err := db.RunMigrations(dsn); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	log.Println("migrations applied (schema is current)")

	// TODO (next stages):
	//   - open a pgxpool connection using dsn
	//   - build the gin router with our handlers
	//   - listen on PORT
}
