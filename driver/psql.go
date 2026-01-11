package driver

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPsql は環境変数 flavor に応じて適切なPostgreSQLに接続し、*gorm.DB を返します。
// - flavor=dev: ローカル環境（DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME を使用）
// - flavor=prd: Neon（DATABASE_URL を使用）
func NewPsql() *gorm.DB {
	dsn := buildDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return db
}

func buildDSN() string {
	flavor := os.Getenv("FLAVOR")

	switch flavor {
	case "dev":
		return buildDevDSN()
	case "prd":
		return buildPrdDSN()
	default:
		log.Fatalf("unknown flavor: %q (expected 'dev' or 'prd')", flavor)
		return ""
	}
}

func buildDevDSN() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}

func buildPrdDSN() string {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required for prd flavor")
	}

	return ensureSSLMode(databaseURL, "require")
}

func ensureSSLMode(dsn, defaultMode string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return dsn
	}

	q := u.Query()
	if q.Get("sslmode") == "" {
		q.Set("sslmode", defaultMode)
		u.RawQuery = q.Encode()
	}

	return u.String()
}
